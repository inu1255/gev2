package models

import (
	"database/sql"
	"errors"

	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev2/config"
)

// Address Entity
type AddressModel struct {
	Model    `xorm:"-"`
	Id       int    `json:"id,omitempty" xorm:"pk autoincr"`
	Center   string `json:"center,omitempty" xorm:"" gev:"中心经纬度"`
	Citycode string `json:"citycode,omitempty" xorm:"" gev:"城市区号"`
	Level    string `json:"level,omitempty" xorm:"" gev:"级别"`
	Name     string `json:"name,omitempty" xorm:"" gev:"城市名"`
	ParentId int    `json:"parent_id,omitempty" xorm:"" gev:"父地址"`
	Value    string `json:"value,omitempty" xorm:"" gev:"三级地址名|分隔"`
}

func (this *AddressModel) TableName() string {
	return "address"
}

// 查找地址
func (this *AddressModel) Search(search *SearchAddress) (*SearchData, error) {
	var user IUserModel
	if u, ok := this.Ctx.Get("user"); ok {
		user, _ = u.(IUserModel)
	}
	bean := &AddressModel{}
	return GetSearchData(this.Db, user, bean, search, func(session *xorm.Session) {
		if search.Keyword != "" {
			session.Where("value like ?", search.Keyword+"%")
		}
		if search.ParentId != 0 {
			session.Where("parent_id=?", search.ParentId)
		}
	})
}

var _pkg_path = config.GetPkgPath()

// 导入地址数据
func (this *AddressModel) LoadSql() ([]sql.Result, error) {
	Db := config.Db
	if ok, err := Db.IsTableEmpty(this); err == nil {
		if ok {
			res, err := Db.ImportFile(_pkg_path + "/address.sql")
			return res, err
		} else {
			return nil, errors.New("数据已经导入address表")
		}
	} else {
		return nil, err
	}
}

// 查找地地
type SearchAddress struct {
	SearchPage
	Keyword  string `json:"keyword,omitempty" gev:"地地 如:'北京|'"`
	ParentId int    `json:"parent_id,omitempty" gev:"父地址id"`
}
