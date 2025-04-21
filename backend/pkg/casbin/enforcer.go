package casbin

import (
	"bufio"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	_ "os"
	"strings"

	// 引入 embed 包（必须导入）
	_ "embed"
)

// 将模型和策略文件嵌入到变量中
//
//go:embed rbac_model.conf
var embeddedRBACModelConf string

//go:embed rbac_policy.csv
var embeddedRBACPolicyCsv string

// 定义一个结构体，包含 Enforcer 和 Adapter
type CasbinMethod struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}

var Casbin *CasbinMethod

func LoadCasbinMethod() *CasbinMethod {
	return Casbin
}

// InitCasbinGorm 初始化 Casbin Gorm适配器，并从嵌入的文件加载模型和策略
func InitCasbinGorm(db *gorm.DB) (*CasbinMethod, error) {
	// 创建 Gorm 适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	// 通过嵌入的模型字符串创建 Casbin 模型
	m, err := model.NewModelFromString(embeddedRBACModelConf)
	if err != nil {
		return nil, err
	}
	// 初始化 Enforcer，使用模型和适配器
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	// 从嵌入的策略 CSV 字符串加载策略
	if err := loadCsvPolicy(enforcer, embeddedRBACPolicyCsv); err != nil {
		return nil, err
	}
	Casbin = &CasbinMethod{
		Enforcer: enforcer,
		Adapter:  a,
	}
	return Casbin, nil
}

// loadCsvPolicy 从 CSV 字符串加载策略
func loadCsvPolicy(e *casbin.Enforcer, csvContent string) error {
	scanner := bufio.NewScanner(strings.NewReader(csvContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		switch parts[0] {
		case "p":
			// 这里假设 p 策略格式为：p, sub, obj, act
			if len(parts) < 4 {
				continue
			}
			_, _ = e.AddPolicy(parts[1], parts[2], parts[3])
		case "g":
			// 如果存在分组策略格式为：g, sub, obj, dom
			if len(parts) == 4 {
				_, _ = e.AddGroupingPolicy(parts[1], parts[2], parts[3])
			}
		}
	}
	e.BuildRoleLinks()
	// 将策略保存到数据库
	return e.SavePolicy()
}

// UpdateUserRoleInDomain 更新用户在某个域（如space_123或global）下的角色，不存在则插入
func UpdateUserRoleInDomain(c *CasbinMethod, userID uint64, role string, domain string) error {
	sub := fmt.Sprintf("user_%d", userID)
	// 删除该用户在该域下的所有角色绑定（g, sub, _, dom）
	oldRoles := c.Enforcer.GetRolesForUserInDomain(sub, domain)
	for _, oldRole := range oldRoles {
		_, err := c.Enforcer.DeleteRoleForUserInDomain(sub, oldRole, domain)
		if err != nil {
			return fmt.Errorf("删除旧角色失败: %v", err)
		}
	}
	// 添加新的角色绑定
	ok, err := c.Enforcer.AddRoleForUserInDomain(sub, role, domain)
	if err != nil || !ok {
		return fmt.Errorf("添加角色失败: %v", err)
	}
	// 持久化策略
	c.Enforcer.BuildRoleLinks()
	err = c.Enforcer.SavePolicy()
	if err != nil {
		return fmt.Errorf("持久化角色策略失败: %v", err)
	}
	return nil
}
