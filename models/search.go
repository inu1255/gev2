package models

import (
	"strings"

	"github.com/go-xorm/xorm"
)

type IDataSearch interface {
	GetSearch(user IUserModel)
}

func searchData(session *xorm.Session, bean interface{}, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	sessionFunc(session)
	total, _ := session.Count(bean)
	sessionFunc(session)
	session.Limit(condition.GetSize(), condition.GetBegin())
	condition.GetOrderDefault(session, "-id")
	data := make([]interface{}, condition.GetSize())
	n := 0
	err := session.Iterate(bean, func(i int, item interface{}) error {
		data[i] = item
		n++
		return nil
	})
	return &SearchData{Content: data[:n], Total: total}, err
}

func GetSearchData(session *xorm.Session, user IUserModel, bean interface{}, condition ISearch, sessionFunc func(session *xorm.Session)) (*SearchData, error) {
	if data, ok := bean.(IDataSearch); ok {
		sessionFunc(session)
		total, _ := session.Count(data)
		sessionFunc(session)
		session.Limit(condition.GetSize(), condition.GetBegin())
		condition.GetOrderDefault(session, "-id")
		data := make([]interface{}, condition.GetSize())
		n := 0
		err := session.Iterate(bean, func(i int, item interface{}) error {
			item.(IDataSearch).GetSearch(user)
			data[i] = item
			n++
			return nil
		})
		return &SearchData{Content: data[:n], Total: total}, err
	}
	return searchData(session, bean, condition, sessionFunc)
}

type SearchData struct {
	Content []interface{} `json:"content" xorm:"" gev:"数据数组"`
	Total   int64         `json:"total" xorm:"" gev:"数据总量"`
	Ext     interface{}   `json:"ext,omitempty" xorm:"" gev:"附加数据"`
}

type ISearch interface {
	GetBegin() int
	GetSize() int
	SetDefaultOrder(o string)
	GetOrder(session *xorm.Session)
	GetOrderDefault(session *xorm.Session, default_order string)
}

// 分页查询
type SearchPage struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	OrderBy string `json:"order_by,omitempty" gev:"排序规则:-id"`
}

func (this *SearchPage) GetSize() int {
	if this.Size < 1 {
		return 10
	}
	return this.Size
}

func (this *SearchPage) GetBegin() int {
	return this.Page * this.GetSize()
}

func (this *SearchPage) SetDefaultOrder(o string) {
	if this.OrderBy == "" {
		this.OrderBy = o
	}
}

func (this *SearchPage) GetOrder(session *xorm.Session) {
	if this.OrderBy != "" {
		orders := strings.Split(this.OrderBy, ",")
		for _, item := range orders {
			if item != "" {
				if item[:1] == "-" && item[:1] != "" {
					session.Desc(item[1:])
				} else {
					session.Asc(item)
				}
			}
		}
	} else {
		session.Desc("id")
	}
}

func (this *SearchPage) GetOrderDefault(session *xorm.Session, default_order string) {
	if this.OrderBy != "" {
		default_order = this.OrderBy
	}
	if default_order != "" {
		orders := strings.Split(default_order, ",")
		for _, item := range orders {
			if item != "" {
				if item[:1] == "-" && item[:1] != "" {
					session.Desc(item[1:])
				} else {
					session.Asc(item)
				}
			}
		}
	}
}
