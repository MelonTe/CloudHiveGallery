package main

import (
	"chg/cmd"
	_ "chg/docs"
	_ "chg/pkg/rds"
)

// @title           CloudHiveGallery
// @version         1.0
// @description		云巢画廊接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cmd.Main()
}

func countGoodNumbers(n int64) int {
	ans := int64(1)
	const mod = 1e9 + 7
	once := int64(20) //每两个数的组合
	times := n / 2    //once ^ times
	//快速幂运算
	pow := func(x, y int64) int64 {
		res := int64(1)
		for y != 0 {
			if y&1 == 1 {
				res = res * x % mod
			}
			x = x * x % mod
			y >>= 1
		}
		return res
	}
	ans = pow(once, times)
	if n%2 == 1 {
		ans = ans * 5 % mod
	}
	return int(ans)
}
