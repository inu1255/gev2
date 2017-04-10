package models

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/inu1255/gev2/config"
)

var Log = config.NewLogger("gev2.models")

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func WordLike(key string) string {
	return strings.Join([]string{"%", key, "%"}, "")
}

func CharLike(key string) string {
	s := strings.Split(key, "")
	return WordLike(strings.Join(s, "%"))
}

// 深度克隆结构体
func cloneValue(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		n := src.NumField()
		for i := 0; i < n; i++ {
			f := src.Field(i)
			cloneValue(f, dst.Field(i))
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
}

// API 错误相关
var (
	NeedAuthError = ApiErr("需要登录", 1255)
)

type ApiError struct {
	code int
	msg  string
}

func (this *ApiError) Error() string {
	return this.msg
}

func (this *ApiError) Code() int {
	return this.code
}

func ApiErr(msg string, code int) error {
	return &ApiError{code: code, msg: msg}
}
