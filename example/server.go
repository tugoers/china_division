/**
 * Copyright 2017 orivil.com. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found at https://mit-license.org.
 */

package main

import (
	"gopkg.in/tugoers/china_division.v0"
	"net/http"
	"log"
	"encoding/json"
	"runtime"
	"path/filepath"
)

func main() {
	// 静态文件服务
	_, filename, _, _ := runtime.Caller(0)
	public := filepath.Join(filepath.Dir(filename), "public")
	http.Handle("/", http.FileServer(http.Dir(public)))

	// 返回所有省级 name 及 code 的 json 对象
	http.HandleFunc("/provinces", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json;charset=UTF-8")
		writer.Write(china_division.GetJsonProvinces())
	})

	// 返回对应省的市级 name 及 code 的 json 对象
	http.HandleFunc("/cities", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("province")
		writer.Header().Add("Content-Type", "application/json;charset=UTF-8")
		city := china_division.GetJsonCities(code)
		writer.Write(city)
	})

	// 返回对应市的区县级 name 及 code 的 json 对象
	http.HandleFunc("/counties", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("city")
		writer.Header().Add("Content-Type", "application/json;charset=UTF-8")
		county := china_division.GetJsonCounties(code)
		writer.Write(county)
	})

	// 获取并返回 post 数据
	http.HandleFunc("/location", func(writer http.ResponseWriter, request *http.Request) {
		code := request.PostFormValue("county")

		// 只需区县级代码就可取得所有信息，保存数据时也只需保存区县级代码
		province, city, county := china_division.GetName(code)
		writer.Header().Add("Content-Type", "application/json;charset=UTF-8")
		json.NewEncoder(writer).Encode(map[string]string{
			"province": province,
			"city":     city,
			"county":   county,
		})
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
