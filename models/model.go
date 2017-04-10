package models

import (
	"time"

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
	GetId() int
	IsNew() bool
	GetModel() *Model
	NewModel(model IModel) IModel
}

type Model struct {
	gev2.BaseService `json:"-" xorm:"-"`
	Db               *xorm.Session `json:"-" xorm:"-"`
	Id               int           `json:"id,omitempty" xorm:"pk autoincr"`
	CreateAt         time.Time     `json:"create_at,omitempty" xorm:"created"`
	UpdateAt         time.Time     `json:"-" xorm:"updated"`
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
func (this *Model) GetId() int {
	return this.Id
}
func (this *Model) IsNew() bool {
	return this.Id < 1
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
