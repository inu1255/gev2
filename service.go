package gev2

import (
	"reflect"

	"github.com/gin-gonic/gin"
	swaggin "github.com/inu1255/go-swagger/gin"
)

type Context struct {
	*gin.Context
}

type ITable interface {
	Init()
}

// 生成 IService 相关的 HanderFunc
func makeServiceHandlerFunc(m reflect.Method, call []convertFunc) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ctx := &Context{Context: gctx}
		n := len(call)
		params := make([]reflect.Value, n)
		params[0] = call[0](ctx)
		service := params[0].Interface().(IService)
		defer func() {
			err := recover()
			service.Finish(err)
		}()
		if service.Before(ctx) {
			for i := 1; i < n; i++ {
				params[i] = call[i](ctx)
			}
			// Log.Println(params[0].Type(), m.Name, params[1:])
			out := m.Func.Call(params)
			data := out[0].Interface()
			msg := out[1].Interface()
			if msg == nil {
				service.After(data, nil)
			} else {
				service.After(data, msg.(error))
			}
		}
	}
}

// 生成 interface{} 相关的 HanderFunc
func makeHandlerFunc(m reflect.Method, call []convertFunc) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		ctx := &Context{Context: gctx}
		n := len(call)
		params := make([]reflect.Value, n)
		params[0] = call[0](ctx)
		for i := 1; i < n; i++ {
			params[i] = call[i](ctx)
		}
		Log.Println(params[0].Type(), m.Name, params[1:])
		out := m.Func.Call(params)
		data := out[0].Interface()
		err := out[1].Interface()
		if err == nil {
			Ok(gctx, data)
		} else {
			Api(gctx, data, err.(error))
		}
	}
}

func Bind(prefix string, service interface{}, summary ...string) {
	router := swaggin.NewRouter(App.Group(prefix), summary...)
	t := reflect.TypeOf(service)
	numMethod := t.NumMethod()
	var instCall convertFunc
	if v, ok := service.(IClass); ok {
		v.SetSelf(v)
		instCall = func(ctx *Context) reflect.Value {
			return reflect.ValueOf(v.New())
		}
	} else {
		// instCall = newInstCall(t.Elem())
		instCall = copyInstCall(service)
	}
	if v, ok := service.(ITable); ok {
		v.Init()
	}
	_, isService := service.(IService)
	for i := 0; i < numMethod; i++ {
		m := t.Method(i)
		flag, path, call := convertMethodParams(router, m)
		if flag == -1 {
			continue
		}
		call[0] = instCall
		if flag == 1 {
			Log.Println("\x1b[34m post\x1b[0m", prefix+path)
			if isService {
				router.POST(path, makeServiceHandlerFunc(m, call))
			} else {
				router.POST(path, makeHandlerFunc(m, call))
			}
		} else {
			Log.Println("\x1b[32mget \x1b[0m", prefix+path)
			if isService {
				router.GET(path, makeServiceHandlerFunc(m, call))
			} else {
				router.GET(path, makeHandlerFunc(m, call))
			}
		}
	}
}
