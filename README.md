# 中国行政区划

## 介绍

最新中国行政区划，数据发布时间：2017-03-10。
数据来源：http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/201703/t20170310_1471429.html

## 安装

go get -v gopkg.in/tugoers/china_division.v0


> **Note:** 如果中国大陆地区用户出现 golang.org/x/net 包不能下载的情况，可选择从 GitHub 下载(go get -v github.com/golang/net)，
然后将 $GOPATH/src/github.com/golang/net 目录复制到 $GOPATH/src/golang.org/x/net 目录下。

## 示例

`cmd.go`:

```GO
package main

import (
	"gopkg.in/tugoers/china_division.v0"
	"fmt"
)

func main() {
	// 获取所有省级名称及省级代码
	provinces := china_division.GetProvinces()

	for _, province := range provinces {
		// 省级代码
		provinceCode := province[0]
		// 省级名称
		provinceName := province[1]

		fmt.Println("┝", provinceName)

		// 根据省级代码获取市级名称及市级代码
		cities := china_division.GetCities(provinceCode)
		for _, city := range cities {
			// 市级代码
			cityCode := city[0]
			// 市级名称
			cityName := city[1]
			fmt.Println("│ ┝", cityName)

			// 根据市级代码获取区县级名称及区县级代码
			counties := china_division.GetCounties(cityCode)
			for _, county := range counties {
				countyCode := county[0]
				countyName := county[1]
				fmt.Println("│ │ ┝", countyCode, ":", countyName)
			}
		}
	}
}
```

更多示例参考 [example](https://github.com/tugoers/china_division/tree/master/example)

## License

Released under the [MIT License](https://mit-license.org).