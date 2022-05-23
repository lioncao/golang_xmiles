package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func RuntimeInfo() (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(2)
}

// http发送简单文本
func HttpSendSimplePage(w *http.ResponseWriter, code int, content string) {
	if w != nil {
		(*w).WriteHeader(code)
		(*w).Write([]byte(content))
	}
}

// 发送空返回
func HttpSend200Empty(w *http.ResponseWriter) {
	HttpSendSimplePage(w, 200, "")
}

// 发送404页面
func HttpSend404NotFound(w *http.ResponseWriter) {
	HttpSendSimplePage(w, 404, "NOT FOUND")
}

func Error(fmtStr string, a ...interface{}) error {
	if a == nil {
		return errors.New(fmtStr)
	}
	return errors.New(fmt.Sprintf(fmtStr, a...))
}

const (
	// 拼接session字符串时用到的时间格式
	SESSION_STR_TIME_FMT = "060102150405"
)

// 获取一个随机的sessionid
func MakeSessionId(title string) string {
	return fmt.Sprintf("%s_%s_%s", title, TimeStringFmt(EMPTY_TIME, SESSION_STR_TIME_FMT), RandomString(16))
}

// 统一的redisKey生成函数
// dbInde：			数据库id（实际功能为给数据划分区段）
// serviceKey：		服务编号,对应唯一的一个service, 例如: auth_001, guobao_003 , balala_100
// userId：			用户的唯一id
// srcKey:			原始key（即功能自己用来区分的key）
func MakeRedisKey(dbIndex int64, serviceKey string, userId int64, srcKey string) string {
	return fmt.Sprintf("%d/%s/%d/%s", dbIndex, serviceKey, userId, srcKey)
}

// 将lua中的table的字符串形式转化为go中的数据结构
func DecodeLuaTableString(name string, tableStr string, data interface{}) error {
	str := strings.TrimSpace(tableStr)
	if str == "" || str == "{}" {
		return nil
	}

	str = strings.Replace(str, "{", "[", -1)
	str = strings.Replace(str, "}", "]", -1)
	str = fmt.Sprintf("{\"%s\":%s}", name, str)
	err := json.Unmarshal([]byte(str), data)
	if err != nil {
		ShowError(name, err.Error(), tableStr, str)
	}
	return err
}

type MapString map[string]string

func (this MapString) EnsureBool(key string, defaultValue bool) bool {
	v, ok := this[key]
	if !ok {
		return defaultValue
	}

	r, e := strconv.ParseBool(v)
	if e != nil {
		return defaultValue
	}
	return r
}

func (this MapString) EnsureString(key string, defaultValue string) string {
	v, ok := this[key]
	if !ok {
		return defaultValue
	}
	return v
}

func (this MapString) EnsureInt64(key string, defaultValue int64) int64 {
	v, ok := this[key]
	if !ok {
		return defaultValue
	}

	r, e := strconv.ParseInt(v, 0, 64)
	if e != nil {
		return defaultValue
	}
	return r
}

func StrIsEmpty(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return false
}

func SliptStr(str string, sep string) []string {
	strs := strings.Split(str, sep)
	if strs == nil {
		return nil
	}

	for i, s := range strs {
		strs[i] = strings.TrimSpace(s)
	}
	return strs
}

func Unix_time_ms(t time.Time) int64 {
	return t.UnixNano() / 1000000
}
func Unix_time_ms_now() int64 {
	return time.Now().UnixNano() / 1000000
}

/******************************************************************************
	slice 相关
******************************************************************************/
func Slice_copy_int64_1(src []int64) []int64 {
	if src == nil {
		return nil
	}

	cnt := len(src)
	dst := make([]int64, cnt, cnt)
	if cnt > 0 {
		copy(dst, src)
	}
	return dst
}

func Slice_copy_int64_2(src [][]int64) [][]int64 {
	if src == nil {
		return nil
	}

	cnt := len(src)
	dst := make([][]int64, 0, cnt)

	for i := 0; i < cnt; i++ {
		dst = append(dst, Slice_copy_int64_1(src[i]))
	}
	return dst
}

func Slice_copy_string_1(src []string) []string {
	if src == nil {
		return nil
	}

	cnt := len(src)
	dst := make([]string, cnt, cnt)
	if cnt > 0 {
		copy(dst, src)
	}
	return dst
}

func Slice_copy_string_2(src [][]string) [][]string {
	if src == nil {
		return nil
	}

	cnt := len(src)
	dst := make([][]string, 0, cnt)

	for i := 0; i < cnt; i++ {
		dst = append(dst, Slice_copy_string_1(src[i]))
	}
	return dst
}

func Int2Str(i int64) string {
	return fmt.Sprintf("%d", i)
}
