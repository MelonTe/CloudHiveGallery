package casbin

import (
	"bufio"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"os"
	"strings"
)

// 定义一个结构体，包含Enforcer和Adapter
type CasbinMethod struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}

var Casbin *CasbinMethod

func LoadCasbinMethod() *CasbinMethod {
	return Casbin
}

// InitCasbinGorm 初始化Casbin Gorm适配器
func InitCasbinGorm(db *gorm.DB) (*CasbinMethod, error) {
	//创建 Gorm适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	// 创建 Casbin Enforcer , 指定自定义的model文件
	enforcer, err := casbin.NewEnforcer("./pkg/casbin/rbac_model.conf", a)
	if err != nil {
		return nil, err
	}
	// 初始化策略，导入csv文件
	if err := loadCsvPolicy(enforcer, "./pkg/casbin/rbac_policy.csv"); err != nil {
		return nil, err
	}
	Casbin = &CasbinMethod{
		Enforcer: enforcer,
		Adapter:  a,
	}
	return Casbin, nil
}

// 从csv导入策略
func loadCsvPolicy(e *casbin.Enforcer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

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
			_, _ = e.AddPolicy(parts[1], parts[2], parts[3])
		case "g":
			if len(parts) == 4 {
				_, _ = e.AddGroupingPolicy(parts[1], parts[2], parts[3])
			}
		}
	}
	e.BuildRoleLinks()
	// 保存到数据库
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
