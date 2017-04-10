package models

import (
	"errors"

	"github.com/go-xorm/xorm"
)

type ItemOwnerModel struct {
	ItemModel `xorm:"extends"`
	OwnerId   int `gev:"-" json:"-" xorm:""`
}

func (this *ItemOwnerModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if this.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (this *ItemOwnerModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if this.Id < 1 && this.OwnerId == 0 {
		this.OwnerId = user.GetId()
		return true
	}
	if this.OwnerId == user.GetId() {
		return true
	}
	return false
}

func (this *ItemOwnerModel) CanDelete(user IUserModel) bool {
	return this.Self().(IItemModel).CanDelete(user)
}

// 搜索，只能查找owner_id为自己的
func (this *ItemOwnerModel) Search(condition *SearchPage) (interface{}, error) {
	bean := this.Self().New()
	user, _ := this.GetUser()
	return GetSearchData(this.Db, user, bean, condition, func(session *xorm.Session) {
		if user == nil {
			session.Where("owner_id=?", 0)
		} else {
			session.Where("owner_id=?", user.GetId())
		}
	})
}

// 批量添加
func (this *ItemOwnerModel) SaveAll(self []IBody) ([]interface{}, error) {
	err := this.Db.Begin()
	if err != nil {
		return nil, err
	}
	count := len(self)
	data := make([]interface{}, count)
	for i := 0; i < count; i++ {
		if bean, err := this.Save(self[i]); err == nil {
			data[i] = bean
		} else {
			err2 := this.Db.Rollback()
			if err2 != nil {
				return nil, err2
			}
			return nil, err
		}
	}
	err = this.Db.Commit()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 批量删除
func (this *ItemOwnerModel) DeleteIds(ids []int) (int64, error) {
	if len(ids) < 1 {
		return 0, errors.New("数组长度不能为0")
	}
	user, _ := this.GetUser()
	return this.Db.In("id", ids).Where("owner_id=?", user.GetId()).Delete(this.Self())
}
