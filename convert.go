package gev2

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/inu1255/go-swagger/core"
)

type convertFunc func(*Context) reflect.Value
type ccF func(paramIn, paramName string, kind reflect.Kind) (convertFunc, bool)

// only new
func newInstCall(t reflect.Type) convertFunc {
	return func(*Context) reflect.Value {
		return reflect.New(t)
	}
}

// new and copy
func copyInstCall(src interface{}) convertFunc {
	return func(*Context) reflect.Value {
		return CloneValue(reflect.ValueOf(src))
	}
}
func CloneValue(src reflect.Value) reflect.Value {
	var dst reflect.Value
	if src.Kind() == reflect.Ptr || src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	dst = reflect.New(src.Type())
	cloneValue(src, dst.Elem())
	return dst
}
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

func newString(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		return reflect.ValueOf(ctx.Param(index))
	}
}
func newInt(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.Atoi(ctx.Param(index)); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newInt64(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.ParseInt(ctx.Param(index), 10, 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newFloat32(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Param(index), 32); e == nil {
			return reflect.ValueOf(float32(r))
		}
		return reflect.ValueOf(float32(0))
	}
}
func newFloat64(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Param(index), 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0.0)
	}
}
func newQueryString(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		return reflect.ValueOf(ctx.Query(index))
	}
}
func newQueryInt(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.Atoi(ctx.Query(index)); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newQueryInt64(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.ParseInt(ctx.Query(index), 10, 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0)
	}
}
func newQueryFloat32(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Query(index), 32); e == nil {
			return reflect.ValueOf(float32(r))
		}
		return reflect.ValueOf(float32(0))
	}
}
func newQueryFloat64(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		if r, e := strconv.ParseFloat(ctx.Query(index), 64); e == nil {
			return reflect.ValueOf(r)
		}
		return reflect.ValueOf(0.0)
	}
}
func newMultiFile(index string) convertFunc {
	return func(ctx *Context) reflect.Value {
		_, header, _ := ctx.Request.FormFile("file")
		return reflect.ValueOf(header)
	}
}
func newJsonCall(t reflect.Type) convertFunc {
	return func(ctx *Context) reflect.Value {
		v := reflect.New(t)
		if err := ctx.BindJSON(v.Interface()); err == nil {
			return v.Elem()
		} else {
			panic(err)
		}
	}
}
func newJsonArrayCall(t reflect.Type, t2 reflect.Type) convertFunc {
	return func(ctx *Context) reflect.Value {
		v := reflect.New(t)
		if err := ctx.BindJSON(v.Interface()); err == nil {
			v = v.Elem()
			v2 := reflect.MakeSlice(t2, v.Len(), v.Cap())
			for i := 0; i < v.Len(); i++ {
				v2.Index(i).Set(v.Index(i))
			}
			return v2
		} else {
			panic(err)
		}
	}
}
func newNilCall(t reflect.Type) convertFunc {
	return func(ctx *Context) reflect.Value {
		return reflect.Zero(t)
	}
}

var (
	name2route = regexp.MustCompile(`([a-z]|^)[A-Z]`)
	_e         = errors.New("")
	errorType  = reflect.TypeOf(_e).Elem()
)

func nameToRoute(from string) string {
	if len(from) == 2 {
		return from[:1] + "/" + strings.ToLower(from[1:])
	}
	return strings.ToLower(from)
}
func convertMethodParams(router ISwagRouter, m reflect.Method) (int, string, []convertFunc) {
	numOut := m.Type.NumOut()
	// Log.Println(m.Func, m.Name, m.Type)
	if numOut != 2 || m.Type.Out(1).Kind() != reflect.Interface || errorType.Implements(m.Type.Out(1)) {
		return -1, "", nil
	}
	numIn := m.Type.NumIn()
	parser := NewMethodParser(m)
	parser.ReadComment()
	path := ""
	if parser.path == "" {
		path = "/" + name2route.ReplaceAllStringFunc(m.Name, nameToRoute)
	} else if parser.path == " " {
		return -1, "", nil
	} else {
		path = parser.path
	}
	methodParams := parser.GetMethodParams()
	parser.ParseToRouter(router)
	call := make([]convertFunc, numIn)
	params := make([]*core.Param, 0, numIn)
	flag := 0 // -1:"ignore current" 0:"GET" 1:"POST"
	for i := 1; i < numIn; i++ {
		var paramName string
		if methodParams[i].Name == "" {
			paramName = fmt.Sprintf("p%d", i)
		} else {
			paramName = methodParams[i].Name
		}
		// custom convertFunc
		var cf convertFunc
		for j := len(ccFs) - 1; j >= 0; j-- {
			if cf, ok := ccFs[j](methodParams[i].In, paramName, m.Type.In(i).Kind()); ok {
				call[i] = cf
				break
			}
		}
		// default convertFunc
		if cf == nil {
			switch m.Type.In(i).Kind() {
			case reflect.String:
				if methodParams[i].In == "path" {
					path += "/:" + paramName
					call[i] = newString(paramName)
				} else {
					call[i] = newQueryString(paramName)
				}
				params = append(params, methodParams[i])
			case reflect.Int:
				if methodParams[i].In == "path" {
					path += "/:" + paramName
					call[i] = newInt(paramName)
				} else {
					call[i] = newQueryInt(paramName)
				}
				params = append(params, methodParams[i])
			case reflect.Int64:
				if methodParams[i].In == "path" {
					path += "/:" + paramName
					call[i] = newInt64(paramName)
				} else {
					call[i] = newQueryInt64(paramName)
				}
				params = append(params, methodParams[i])
			case reflect.Float32:
				if methodParams[i].In == "path" {
					path += "/:" + paramName
					call[i] = newFloat32(paramName)
				} else {
					call[i] = newQueryFloat32(paramName)
				}
				params = append(params, methodParams[i])
			case reflect.Float64:
				if methodParams[i].In == "path" {
					path += "/:" + paramName
					call[i] = newFloat64(paramName)
				} else {
					call[i] = newQueryFloat64(paramName)
				}
				params = append(params, methodParams[i])
			case reflect.Struct, reflect.Ptr, reflect.Slice, reflect.Map:
				if flag == 1 {
					flag = -1
					break
				} else if m.Type.In(i).Kind() == reflect.Slice && m.Type.In(i).Elem().Kind() == reflect.Interface && methodParams[i].Name == "self" {
					typ := reflect.SliceOf(m.Type.In(0))
					call[i] = newJsonArrayCall(typ, m.Type.In(i))
					router.Body(reflect.New(typ).Interface())
				} else if m.Type.In(i).String() == "*multipart.FileHeader" {
					methodParams[i].In = "formData"
					methodParams[i].Type = "file"
					methodParams[i].AllowMultiple = true
					params = append(params, methodParams[i])
					call[i] = newMultiFile(paramName)
				} else {
					typ := m.Type.In(i)
					call[i] = newJsonCall(typ)
					router.Body(reflect.New(typ).Interface())
				}
				flag = 1
			case reflect.Interface:
				if methodParams[i].Name == "self" {
					if flag == 1 {
						flag = -1
						break
					} else {
						typ := m.Type.In(0)
						call[i] = newJsonCall(typ)
						router.Body(reflect.New(typ).Interface())
					}
					flag = 1
				}
			case reflect.Func:
				call[i] = newNilCall(m.Type.In(i))
			default:
				Log.Println("default", i, m.Type.In(i).Kind())
				flag = -1
			}
		}
		if flag == -1 {
			break
		}
	}
	router.Params(params...)
	return flag, path, call
}
