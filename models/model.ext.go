package models

import "github.com/go-xorm/xorm"

type ModelExt struct {
	Model `json:"-" xorm:"-"`
	ok    bool
}

func (this *ModelExt) AfterInsert() {
	this.ok = true
}

func (this *ModelExt) AfterSet(name string, cell xorm.Cell) {
	this.ok = true
}

func (this *ModelExt) Save() (err error) {
	if this.ok {
		_, err = this.Db.Id(GetPk(this.Self())).Update(this.Self())
	} else if _, err = this.Db.InsertOne(this.Self()); err != nil {
		_, err = this.Db.Id(GetPk(this.Self())).Update(this.Self())
	}
	return
}

func (this *ModelExt) GetUser() (IUserModel, bool) {
	if u, ok := this.Ctx.Get("user"); ok {
		return u.(IUserModel), ok
	}
	return nil, false
}

// @path save
func (this *ModelExt) Add(self interface{}) (interface{}, error) {
	var err error
	if _, err = this.Db.InsertOne(self); err != nil {
		_, err = this.Db.Id(GetPk(self)).Update(self)
	}
	return self, err
}

func (this *ModelExt) Search(condition *SearchPage) (interface{}, error) {
	var user IUserModel
	if u, ok := this.Ctx.Get("user"); ok {
		user = u.(IUserModel)
	}
	bean := this.Self().New()
	return GetSearchData(this.Db, user, bean, condition, func(session *xorm.Session) {})
}
