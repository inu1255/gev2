package models

import (
	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev2"
	"github.com/inu1255/gev2/config"
)

// 基础数据模型
// 父类可以通过Self()获取实例
type IModel interface {
	gev2.ITable
	gev2.IClass
	gev2.IService
	GetModel() *Model
	NewModel(model IModel) IModel
}

type Model struct {
	gev2.BaseService `json:"-" xorm:"-"`
	Db               *xorm.Session `json:"-" xorm:"-"`
}

func (this *Model) Before(ctx *gev2.Context) bool {
	ok := this.BaseService.Before(ctx)
	this.Db = config.Db.NewSession()
	return ok
}
func (this *Model) Finish(err interface{}) {
	this.Db.Close()
	this.BaseService.Finish(err)
}
func (this *Model) GetModel() *Model {
	return this
}
func (this *Model) NewModel(model IModel) IModel {
	that := model.GetModel()
	that.Ctx = this.Ctx
	that.Db = this.Db
	return NewModel(model)
}

func (this *Model) Init() {
	config.Db.Sync2(this.Self())
}

func NewModel(model IModel) IModel {
	model.SetSelf(model)
	return model
}
