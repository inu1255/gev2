package gev2

import (
	"github.com/gin-gonic/gin"
	"github.com/inu1255/go-swagger/core"
)

type ISwagRouter interface {
	core.ISwagRouterBase
	Params(params ...*core.Param) ISwagRouter
	Body(body interface{}) ISwagRouter
	Data(data interface{}) ISwagRouter
	Info(info ...string) ISwagRouter
	AddPath(basePath, route, ms string)

	GET(route string, handler gin.HandlerFunc)
	POST(route string, handler gin.HandlerFunc)
	Handle(ms, route string, handler gin.HandlerFunc)
}

type GinSwaggerRouter struct {
	core.SwagRouter

	group *gin.RouterGroup
}

// 封装为链式调用
func (this *GinSwaggerRouter) Params(ps ...*core.Param) ISwagRouter {
	this.SwagRouter.Params(ps...)
	return this
}

func (this *GinSwaggerRouter) Body(body interface{}) ISwagRouter {
	this.SwagRouter.Body(body)
	return this
}

func (this *GinSwaggerRouter) Data(data interface{}) ISwagRouter {
	this.SwagRouter.Data(data)
	return this
}

func (this *GinSwaggerRouter) Info(info ...string) ISwagRouter {
	this.SwagRouter.Info(info...)
	return this
}

func (this *GinSwaggerRouter) Handle(ms, route string, handler gin.HandlerFunc) {
	this.AddPath(this.group.BasePath(), route, ms)
	this.group.Handle(ms, route, handler)
	this.Clear()
}

func (this *GinSwaggerRouter) GET(route string, handler gin.HandlerFunc) {
	this.Handle("GET", route, handler)
}

func (this *GinSwaggerRouter) POST(route string, handler gin.HandlerFunc) {
	this.Handle("POST", route, handler)
}

var Swag = core.NewSwagger()

func NewRouter(g *gin.RouterGroup, summary ...string) ISwagRouter {
	Swag.AddTag(g.BasePath()[1:], summary...)
	router := new(GinSwaggerRouter)
	router.Swagger = Swag
	router.group = g
	return router
}
