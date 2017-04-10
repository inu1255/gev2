package gev2

import (
	"errors"
	"fmt"
)

type IService interface {
	Before(ctx *Context) bool
	Finish(err interface{})
	After(data interface{}, err error)
}

type BaseService struct {
	Class `json:"-" xorm:"-"`
	Ctx   *Context `json:"-" xorm:"-"`
}

func (this *BaseService) Before(ctx *Context) bool {
	this.Ctx = ctx
	return true
}
func (this *BaseService) After(data interface{}, err error) {
	Api(this.Ctx.Context, data, err)
}
func (this *BaseService) Finish(err interface{}) {
	if err != nil {
		Log.Printf("%v\n\033[31m%s\033[0m", err, string(stack()))
		Err(this.Ctx.Context, 500, errors.New(fmt.Sprintf("系统错误 : %v", err)))
	}
}
