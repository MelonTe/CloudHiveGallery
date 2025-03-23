package consts

const (
	REVIEWING = 0 //待审核
	PASS      = 1 //审核通过
	REJECT    = 2 //审核拒绝
	ALL       = 3 //任意状态
)

//校验审核参数是否存在，用于写数据库的WRAPPER
func ReviewValueExist(value int) bool {
	switch value {
	case REVIEWING, PASS, REJECT:
		return true
	default:
		return false
	}
}
