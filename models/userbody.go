package models

import "errors"

type UserModelBody struct {
	Id       int    `json:"id,omitempty" xorm:"" gev:"id，修改数据时使用"`
	Telphone string `json:"telphone,omitempty" xorm:"" gev:"用户账号，需要管理员权限"`
	Password string `json:"password,omitempty" xorm:"" gev:"用户密码，需要管理员权限"`
	Nickname string `json:"nickname" xorm:"" gev:"用户昵称"`
}

func (this *UserModelBody) CopyTo(user IUserModel, bean interface{}) error {
	data := bean.(*UserModel)
	if this.Id < 1 {
		if role, ok := user.(IUserRoleModel); ok && role.IsAdmin() {
			if this.Telphone == "" || this.Password == "" {
				return errors.New("账号/密码不能为空")
			}
			data.Telphone = this.Telphone
			data.Password = data.EncodePwd(this.Password)
		} else {
			return errors.New("需要id")
		}
	}
	data.Nickname = this.Nickname
	data.Id = this.Id
	return nil
}
