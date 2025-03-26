package consts

// SpaceLevel 表示空间等级及其属性。
type SpaceLevel struct {
	Text     string // 等级名称
	Value    int    // 等级值
	MaxCount int64  // 最大图片数量
	MaxSize  int64  // 最大图片总大小（字节）
}

// 定义每个空间等级及其属性。
var (
	COMMON = SpaceLevel{
		Text:     "普通版",             // 普通版
		Value:    0,                 // 等级值为 0
		MaxCount: 100,               // 最大图片数量为 100
		MaxSize:  100 * 1024 * 1024, // 最大图片总大小为 100MB
	}
	PROFESSIONAL = SpaceLevel{
		Text:     "专业版",              // 专业版
		Value:    1,                  // 等级值为 1
		MaxCount: 1000,               // 最大图片数量为 1000
		MaxSize:  1000 * 1024 * 1024, // 最大图片总大小为 1000MB
	}
	FLAGSHIP = SpaceLevel{
		Text:     "旗舰版",               // 旗舰版
		Value:    2,                   // 等级值为 2
		MaxCount: 10000,               // 最大图片数量为 10000
		MaxSize:  10000 * 1024 * 1024, // 最大图片总大小为 10000MB
	}
)

// GetSpaceLevelByValue 根据等级值获取对应的 SpaceLevel。
func GetSpaceLevelByValue(value int) *SpaceLevel {
	switch value {
	case COMMON.Value:
		return &COMMON
	case PROFESSIONAL.Value:
		return &PROFESSIONAL
	case FLAGSHIP.Value:
		return &FLAGSHIP
	default:
		return nil // 如果没有匹配的等级值，返回 nil
	}
}
