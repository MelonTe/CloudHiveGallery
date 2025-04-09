package consts

//团队空间权限枚举
const (
	SPACEROLE_VIEWER = "viewer" // 只读权限
	SPACEROLE_EDITOR = "editor" // 编辑权限
	SPACEROLE_ADMIN  = "admin"  // 管理员权限
)

func IsSpaceUserRoleExist(role string) bool {
	switch role {
	case SPACEROLE_VIEWER, SPACEROLE_EDITOR, SPACEROLE_ADMIN:
		return true
	default:
		return false
	}
}
