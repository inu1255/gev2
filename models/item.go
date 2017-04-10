package models

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev2"
)

type IBody interface {
	CopyTo(user IUserModel, data interface{}) error
}
type IDataDetail interface {
	GetDetail(user IUserModel)
}

type IItemModel interface {
	IModel
	IBody
	CanRead(user IUserModel) bool
	CanWrite(user IUserModel) bool
	CanDelete(user IUserModel) bool
}

type ItemModel struct {
	Model `xorm:"extends"`
}

func (this *ItemModel) CanRead(user IUserModel) bool {
	return true
}
func (this *ItemModel) CanWrite(user IUserModel) bool {
	return true
}
func (this *ItemModel) CanDelete(user IUserModel) bool {
	return false
}
func (this *ItemModel) CopyTo(user IUserModel, bean interface{}) error {
	t1 := reflect.TypeOf(this.Self())
	t2 := reflect.TypeOf(bean)
	if t1.String() == t2.String() {
		cloneValue(reflect.ValueOf(this), reflect.ValueOf(bean))
		return nil
	} else {
		return errors.New(fmt.Sprintf("%v -> %v 复制错误", t1, t2))
	}
}
func (this *ItemModel) GetUser() (IUserModel, bool) {
	if u, ok := this.Ctx.Get("user"); ok {
		return u.(IUserModel), ok
	}
	return nil, false
}
func (this *ItemModel) GetCurrentUserId(user_id int) int {
	if user_id > 0 {
		return user_id
	}
	if user, ok := this.GetUser(); ok {
		return user.GetId()
	}
	return 0
}

// 查看详情
func (this *ItemModel) Info(id int) (interface{}, error) {
	data := this.Self().New().(IItemModel)
	user, _ := this.GetUser()
	return ItemGetInfo(this.Db, user, data, id)
}

// 添加/修改
func (this *ItemModel) Save(self IBody) (interface{}, error) {
	if v, ok := self.(gev2.IClass); ok {
		v.SetSelf(v)
	}
	bean := this.Self().New().(IItemModel)
	user, _ := this.GetUser()
	return ItemSave(this.Db, user, bean, self)
}

// 删除
func (this *ItemModel) Delete(id int) (int64, error) {
	bean := this.Self().New().(IItemModel)
	user, _ := this.GetUser()
	return ItemDelete(this.Db, user, bean, id)
}

// 查找
func (this *ItemModel) Search(condition *SearchPage) (interface{}, error) {
	bean := this.Self().New()
	user, _ := this.GetUser()
	return GetSearchData(this.Db, user, bean, condition, func(session *xorm.Session) {})
}

func ItemGetInfo(Db *xorm.Session, user IUserModel, bean IItemModel, id interface{}) (interface{}, error) {
	ok, err := Db.Table(bean).Where("id=?", id).Get(bean)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("不存在")
	}
	if !bean.CanRead(user) {
		return nil, errors.New("没有权限")
	}
	if data, ok := bean.(IDataDetail); ok {
		data.GetDetail(user)
		return data, err
	}
	return bean, nil
}

func ItemSave(Db *xorm.Session, user IUserModel, bean IItemModel, schema IBody) (interface{}, error) {
	var err error
	if err = schema.CopyTo(user, bean); err != nil {
		return nil, err
	}
	// 更新或插入
	if bean.IsNew() {
		if !bean.CanWrite(user) {
			return nil, errors.New("没有权限")
		}
		_, err = Db.InsertOne(bean)
		if err != nil {
			return nil, err
		}
		if data, ok := bean.(IDataDetail); ok {
			return data, nil
		}
		return bean, nil
	} else {
		item := bean.New().(IItemModel)
		var ok bool
		ok, err = Db.Id(bean.GetId()).Get(item)
		if !ok {
			return nil, errors.New("不存在")
		}
		if err != nil {
			return nil, err
		}
		if !item.CanWrite(user) {
			return nil, errors.New("没有修改权限")
		}
		_, err = Db.ID(bean.GetId()).Update(bean)
		if err != nil {
			return nil, err
		}
		return ItemGetInfo(Db, user, bean, bean.GetId())
	}
}

func ItemDelete(Db *xorm.Session, user IUserModel, bean IItemModel, id interface{}) (int64, error) {
	ok, err := Db.Id(id).Get(bean)
	if !ok {
		return 0, errors.New("不存在")
	}
	if err != nil {
		return 0, err
	}
	if !bean.CanDelete(user) {
		return 0, errors.New("没有权限")
	}
	return Db.ID(id).Delete(bean.New())
}
