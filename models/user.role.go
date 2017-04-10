package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inu1255/gev2/config"
)

var (
	// 用于判断登录者类型
	user_roles = make([]IUserRoleModel, 0, 2)
)

type IUserRoleModel interface {
	IUserModel
	IsAdmin() bool
	GetRole() string
}

type UserRoleModel struct {
	UserModel `xorm:"extends"`
	Role      string `gev:"用户角色" json:"role,omitempty" xorm:"not null default '普通用户' VARCHAR(32)"`
}

func (this *UserRoleModel) BeforeInsert() {
	this.Role = this.Self().(IUserRoleModel).GetRole()
}
func (this *UserRoleModel) BeforeUpdate() {
	this.Role = ""
}
func (this *UserRoleModel) Init() {
	this.UserModel.Init()
	user_roles = append(user_roles, this.Self().(IUserRoleModel))
}
func (this *UserRoleModel) GetRole() string {
	return this.Role
}
func (this *UserRoleModel) IsAdmin() bool {
	if this.Role == "管理员" {
		return true
	}
	return false
}

func UserRoleMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 当前登录用户数据
		token := c.Query("access_token")
		if token != "" {
			now := time.Now()
			user := &UserRoleModel{}
			ok, _ := config.Db.Cols("id", "role").Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(user)
			if ok {
				// Log.Println(len(user_roles))
				// 判断登录者类型
				for _, item := range user_roles {
					// Log.Println(item.GetRole())
					if user.Role == item.GetRole() {
						bean := item.New()
						config.Db.ID(user.Id).Get(bean)
						c.Set("user", bean)
						break
					}
				}
			}
		}
	}
}
