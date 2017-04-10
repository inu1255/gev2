package models

import (
	"errors"

	"github.com/go-xorm/xorm"
)

type ItemRoleModel struct {
	ItemOwnerModel `xorm:"extends"`
}

func (this *ItemRoleModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if u, ok := user.(IUserRoleModel); ok && u.IsAdmin() {
		return true
	}
	return this.ItemOwnerModel.CanRead(user)
}
func (this *ItemRoleModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if u, ok := user.(IUserRoleModel); ok && u.IsAdmin() {
		return true
	}
	return this.ItemOwnerModel.CanWrite(user)
}
func (this *ItemRoleModel) GetUser() (IUserRoleModel, bool) {
	if u, ok := this.Ctx.Get("user"); ok {
		return u.(IUserRoleModel), ok
	}
	return nil, false
}

// 搜索，除了管理员，只能查找owner_id为自己的
func (this *ItemRoleModel) Search(condition *SearchPage) (interface{}, error) {
	bean := this.Self().New()
	user, _ := this.GetUser()
	return GetSearchData(this.Db, user, bean, condition, func(session *xorm.Session) {
		if user == nil {
			session.Where("owner_id=?", 0)
		} else if !user.IsAdmin() {
			session.Where("owner_id=?", user.GetId())
		}
	})
}

// 批量删除
func (this *ItemRoleModel) DeleteIds(ids []int) (int64, error) {
	if len(ids) < 1 {
		return 0, errors.New("数组长度不能为0")
	}
	user, _ := this.GetUser()
	if user.IsAdmin() {
		return this.Db.In("id", ids).Delete(this.Self().New())
	} else {
		return this.Db.In("id", ids).Where("owner_id=?", user.GetId()).Delete(this.Self().New())
	}
}
