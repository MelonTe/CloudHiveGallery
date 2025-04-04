package consts

//团队空间权限枚举
const (
	VIEWER = "viewer" // 只读权限
	EDITOR = "editor" // 编辑权限
	ADMIN  = "admin"  // 管理员权限
)

func IsSpaceUserRoleExist(role string) bool {
	switch role {
	case VIEWER, EDITOR, ADMIN:
		return true
	default:
		return false
	}
}
