/**
 * Copyright 2017 orivil.com. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found at https://mit-license.org.
 */

package china_division

import (
	"encoding/json"
	"golang.org/x/net/html"
	"runtime"
	"path/filepath"
	"os"
)

type Code int

const (
	UnknownCode Code = iota
	ProvinceCode
	CityCode
	CountyCode
)

// 所有省
// e.g. {{11, 北京市}, {12, 天津市}, ...}
var provinces [][2]string

// 所有省 Json 字符串
// e.g. {{11, 北京市}, {12, 天津市}, ...}
var provinceJson []byte

// e.g. {11: 北京市, 12: 天津市, ...}
var provinceMap = make(map[string]string, 40)

// 所有市
// e.g. {13: {{1301, 石家庄市}, {1302, 唐山市}}, ...}
var cities = make(map[string][][2]string, 40)

// 所有市 Json 字符串
// e.g. {13: {{1301, 石家庄市}, {1302, 唐山市}}, ...}
var cityJson = make(map[string][]byte, 40)

// e.g. {1301: 石家庄市, 1302: 唐山市, ...}
var cityMap = make(map[string]string, 300)

// 所有区县
// e.g. {1301: {{130101, 市辖区}, {130102, 长安区}}, ...}
var counties = make(map[string][][2]string, 300)

// 所有区县 Json 字符串
// e.g. {1301: {{130101, 市辖区}, {130102, 长安区}}, ...}
var countyJson = make(map[string][]byte, 300)

// e.g. {130104: 桥西区, 130102: 长安区, ...}
var countyMap = make(map[string]string, 1000)

// 返回 code 对应的名称. code 可以是 2 位数的省级代码, 可以是 4 位数的市级代码,
// 也可以是 6 位数的所有级别代码
func GetName(code string) (province, city, county string) {
	ln := len(code)
	if ln >= 2 {
		province = provinceMap[code[:2]]
	}

	if ln >= 4 {
		city = cityMap[code[:4]]
	}

	if ln >= 6 {
		county = countyMap[code[:6]]
	}
	return
}

// 获得 json 格式的所有省名及其代码
func GetJsonProvinces() []byte {

	return provinceJson
}

// 获得 slice 格式的所有省名及其代码
func GetProvinces() [][2]string {

	return provinces
}

// 获得 json 格式的所有市名及其代码
func GetJsonCities(code string) []byte {
	// 截取 province code
	if len(code) > 2 {
		code = code[:2]
	}
	if CodeType(code) != ProvinceCode {
		return []byte("{}")
	}
	return cityJson[code]
}

// 获得 slice 格式的所有市名及其代码
func GetCities(code string) [][2]string {
	// 截取 province code
	if len(code) > 2 {
		code = code[:2]
	}
	if CodeType(code) != ProvinceCode {
		return nil
	}
	return cities[code]
}

// 获得 json 格式的所有区县名及其代码
func GetJsonCounties(code string) []byte {
	// 截取 city code
	if len(code) > 4 {
		code = code[:4]
	}
	if CodeType(code) != CityCode {
		return []byte("{}")
	}
	return countyJson[code]
}

// 获得 slice 格式的所有区县名及其代码
func GetCounties(code string) [][2]string {
	// 截取 city code
	if len(code) > 4 {
		code = code[:4]
	}
	if CodeType(code) != CityCode {
		return nil
	}
	return counties[code]
}

// 返回当前代码的类型
func CodeType(code string) Code {
	switch len(code) {
	case 2:
		if _, ok := provinceMap[code]; ok {
			return ProvinceCode
		}
	case 4:
		if _, ok := cityMap[code]; ok {
			return CityCode
		}
	case 6:
		if _, ok := countyMap[code]; ok {
			return CountyCode
		}
	default:
		return UnknownCode
	}
	return UnknownCode
}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	file, err := os.Open(filepath.Join(filepath.Dir(filename), "data", "code.html"))
	if err != nil {
		panic(err)
	}
	doc, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	var call func(*html.Node)
	call = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			codeNode := n.FirstChild
			nameNode := codeNode.NextSibling
			c := codeNode.FirstChild.Data
			name := nameNode.FirstChild.Data
			provinceCode := c[:2]
			cityCode := c[2:4]
			countyCode := c[4:]
			provinceCityCode := c[:4]

			// 当前名称视为省名
			if cityCode == "00" {
				if countyCode == "00" {
					// {11, 北京市}
					province := [2]string{provinceCode, name}

					// {{11, 北京市}, {12, 天津市}, ...}
					provinces = append(provinces, province)

					// {11: 北京市, 12: 天津市, ...}
					provinceMap[provinceCode] = name
				}
				// 当前名称视为市名
			} else if countyCode == "00" {
				// {1301, 石家庄市}
				city := [2]string{provinceCityCode, name}

				// {13: {{1301, 石家庄市}, {1302, 唐山市}}, ...}
				cities[provinceCode] = append(cities[provinceCode], city)

				// {1301: 石家庄市, 1302: 唐山市, ...}
				cityMap[provinceCityCode] = name
				// 当前名称视为区县名
			} else {
				// {130102, 长安区}
				county := [2]string{c, name}

				// {1301: {{130101, 市辖区}, {130102, 长安区}}, ...}
				counties[provinceCityCode] = append(counties[provinceCityCode], county)

				// {130104: 桥西区, 130102: 长安区, ...}
				countyMap[c] = name
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			call(c)
		}
	}

	call(doc)

	// 初始化省 json 对象
	provinceJson, err = json.Marshal(provinces)
	if err != nil {
		panic(err)
	}

	// 初始化市 json 对象
	for provinceCode, city := range cities {
		cityJson[provinceCode], err = json.Marshal(city)
		if err != nil {
			panic(err)
		}
	}

	// 初始化区县 json 对象
	for provinceCityCode, county := range counties {
		countyJson[provinceCityCode], err = json.Marshal(county)
		if err != nil {
			panic(err)
		}
	}
}