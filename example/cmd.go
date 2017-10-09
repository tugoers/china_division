/**
 * Copyright 2017 orivil.com. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found at https://mit-license.org.
 */

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
