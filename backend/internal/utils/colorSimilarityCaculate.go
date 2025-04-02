package utils

import (
	"math"
	"strconv"
)

// 计算两个颜色的相似度
// 参数为颜色的RGB值，十六进制表示，例如0xFF0000表示红色或#FF0000
// 返回值为相似度，范围为0-1，1表示完全相同，0表示完全不同
func ColorSimilarity(color1, color2 string) float64 {
	// 将十六进制颜色转换为RGB
	r1, g1, b1 := hexToRGB(color1)
	r2, g2, b2 := hexToRGB(color2)
	// 计算欧几里得距离
	distance := math.Sqrt(math.Pow(float64(r1-r2), 2) + math.Pow(float64(g1-g2), 2) + math.Pow(float64(b1-b2), 2))

	// 最大可能距离（颜色空间对角线长度）
	maxDistance := math.Sqrt(3 * math.Pow(255, 2))

	// 归一化相似度
	return 1 - (distance / maxDistance)
}

// 辅助函数：将十六进制颜色转换为RGB
func hexToRGB(hex string) (int, int, int) {
	if len(hex) <= 6 {
		return 0, 0, 0 // 返回黑色，表示无效颜色
	}
	var r, g, b int64
	if len(hex) == 7 {
		//"#FF0000"格式
		r, _ = strconv.ParseInt(hex[1:3], 16, 32)
		g, _ = strconv.ParseInt(hex[3:5], 16, 32)
		b, _ = strconv.ParseInt(hex[5:7], 16, 32)
	} else {
		//"0xFF0000"格式
		r, _ = strconv.ParseInt(hex[2:4], 16, 32)
		g, _ = strconv.ParseInt(hex[4:6], 16, 32)
		b, _ = strconv.ParseInt(hex[6:8], 16, 32)
	}

	return int(r), int(g), int(b)
}
