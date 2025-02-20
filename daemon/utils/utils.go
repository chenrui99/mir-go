// Copyright [2022] [MIN-Group -- Peking University Shenzhen Graduate School Multi-Identifier Network Development Group]
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package utils
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/12/21 4:32 PM
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//

package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"minlib/encoding"
	"time"
)

// 时间模块部分
const (
	BINano  = "2006-01-02 15:04:05.000000000"
	BIMicro = "2006-01-02 15:04:05.000000"
	BIMil   = "2006-01-02 15:04:05.000"
	BISec   = "2006-01-02 15:04:05"
	BICST   = "2006-01-02 15:04:05 +0800 CST"
	BIUTC   = "2006-01-02 15:04:05 +0000 UTC"
	BIDate  = "2006-01-02"
	BITime  = "15:04:05"
)

func main() {

	// snapchatAPI 获取到的时间格式是这种格式的字符串：2020-08-21T10:59:53.850Z
	timeStr := "2020-08-21T10:59:53.850Z"
	// 字符串转时间 得到的是CST 中国标准时间
	ret1, _ := Timestr2Time(timeStr)
	fmt.Printf("ret1>>>  %v, %T \n", ret1, ret1) // 2020-08-21 10:59:53 +0800 CST, time.Time
	// 字符串转时间戳
	ret2, _ := Timestr2Timestamp(timeStr)
	fmt.Printf("ret2>>> %v, %T \n", ret2, ret2) //1597978793, int64

	// 时间戳转时间
	ret3 := Timestamp2TimeSec(1597978793)
	fmt.Printf("ret3>>> %v, %T \n", ret3, ret3) //2020-08-21 10:59:53 +0800 CST, time.Time

	// 时间转字符串 —— ret1 是CST格式的时间
	ret4 := ret1.Format(BICST)
	fmt.Printf("ret4>>> %v, %T \n", ret4, ret4) //2020-08-21 10:59:53 +0800 CST, string
}

// Timestr2Time 时间字符串转时间
func Timestr2Time(str string) (time.Time, error) {
	return Timestr2TimeBasic(str, "", nil)
}

// Timestr2Timestamp 时间字符串转时间戳
func Timestr2Timestamp(str string) (int64, error) {
	return Timestr2TimestampBasic(str, "", nil)
}

// Timestamp2TimeSec 时间戳转时间 秒
func Timestamp2TimeSec(stamp int64) time.Time {
	return Timestamp2Time(stamp, 0)
}

// Timestr2TimeBasic base...
func Timestr2TimeBasic(value string, resultFormat string, resultLoc *time.Location) (time.Time, error) {
	/**
	  - params
	      value:             转换内容字符串
	      resultFormat:    结果时间格式
	      resultLoc:        结果时区
	*/
	resultLoc = getLocationDefault(resultLoc)
	useFormat := []string{ // 可能的转换格式
		BINano, BIMicro, BIMil, BISec, BICST, BIUTC, BIDate, BITime,
		time.RFC3339,
		time.RFC3339Nano,
	}
	var t time.Time
	for _, usef := range useFormat {
		tt, error := time.ParseInLocation(usef, value, resultLoc)
		t = tt
		if error != nil {
			continue
		}
		break
	}
	if t == getTimeDefault(resultLoc) {
		return t, errors.New("时间字符串格式错误")
	}

	if resultFormat == "" {
		resultFormat = "2006-01-02 15:04:05"
	}
	st := t.Format(resultFormat)
	fixedt, _ := time.ParseInLocation(resultFormat, st, resultLoc)

	return fixedt, nil
}

func Timestr2TimestampBasic(str string, format string, loc *time.Location) (int64, error) {
	t, err := Timestr2TimeBasic(str, format, loc)
	if err != nil {
		return -1., err
	}
	return (int64(t.UnixNano()) * 1) / 1e9, nil
}

func Timestamp2Time(stamp int64, nsec int64) time.Time {
	return time.Unix(stamp, nsec)
}

// 获取time默认值, 造一个错误
func getTimeDefault(loc *time.Location) time.Time {
	loc = getLocationDefault(loc)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "", loc)
	return t
}

func getLocationDefault(loc *time.Location) *time.Location {
	if loc == nil {
		loc, _ = time.LoadLocation("Local")
	}
	return loc
}

// 计算模块部分

func Min(x, y encoding.SizeT) (min encoding.SizeT) {
	if x < y {
		return x
	}
	return y
}

// RandomBytes
// @Description:  生成随机用的Byte数组，因为go中没有默认参数的设置方法，所以这里使用可变长参数的方式来实现
//                使用默认的种子0  -> RandomBytes(1000)
//                使用时间作为种子 -> RandomBytes(1000, time.Now().Unix())
// @param n		  要生成的数组长度
// @param seed	  随机种子
// @return []byte 返回随机的Byte数组
//
func RandomBytes(n int, seed ...int64) []byte {
	if len(seed) > 1 {
		panic("输入参数错误, 仅接受一个参数")
	}
	var s int64 = 0
	if len(seed) != 0 {
		s = seed[0]
	}
	rand := rand.New(rand.NewSource(s))
	r := make([]byte, n)
	if _, err := rand.Read(r); err != nil {
		panic("rand.Read failed: " + err.Error())
	}
	return r
}

// RandomString
// @Description:  生成随机用的字符串，字符串组成由letterBytes组成
//                使用默认的种子0  -> RandomString(1000)
//                使用时间作为种子 -> RandomString(1000, time.Now().Unix())
// @param n		  要生成的字符串长度
// @param seed	  随机种子
// @return []byte 返回随机字符串
//
func RandomString(n int, seed ...int64) string {
	if len(seed) > 1 {
		panic("输入参数错误, 仅接受一个参数")
	}
	var s int64 = 0
	if len(seed) != 0 {
		s = seed[0]
	}
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand := rand.New(rand.NewSource(s))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// RandomMINName
// @Description:  生成随机的MIN name，格式类似于：/min/pku/edu
// @param min	  生成的name的最小component数量，如/min/pku/edu的component数量为3个
// @param max	  生成的name的最大component数量
// @param componentLength 生成名字的component的随机大小，范围为：[1, componentLength]
// @param seed	  随机种子
// @return string 返回的是随机生成的MIN name
//
func RandomMINName(min, max, componentLength int, seed int64) string {
	nameComponentCnt := rand.Intn(max-min) + min
	var name string
	rand := rand.New(rand.NewSource(seed))
	for i := 0; i < nameComponentCnt; i++ {
		name += "/"
		name += RandomString(rand.Intn(componentLength)+1, seed+int64(i))
	}
	return name
}
