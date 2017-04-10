package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inu1255/gev2/config"
)

type IUserModel interface {
	IItemModel
	JudgeChpwdCode(code string) error
	GetByTelphone(telphone string) bool
	EncodePwd(string) string
	Exist(telphone string) bool
}

type UserModel struct {
	ItemRoleModel `xorm:"extends"`
	Nickname      string `gev:"用户昵称" json:"nickname" xorm:""`
	Telphone      string `gev:"电话号码" json:"telphone" xorm:"varchar(32) unique(telphone) not null"`
	Password      string `gev:"密码" json:"-" xorm:""`
}

func (this *UserModel) TableName() string {
	return "user"
}
func (this *UserModel) Init() {
	this.ItemRoleModel.Init()
	config.Db.Sync2(new(AccessToken))
}

// 自己可以查看自己
func (this *UserModel) CanRead(user IUserModel) bool {
	if user == nil {
		return false
	}
	if this.Id == user.GetId() {
		return true
	}
	return this.ItemRoleModel.CanWrite(user)
}

// 自己可以修改自己
func (this *UserModel) CanWrite(user IUserModel) bool {
	if user == nil {
		return false
	}
	if this.Id == user.GetId() {
		return true
	}
	return this.ItemRoleModel.CanWrite(user)
}

// 添加者/管理员 才能删除
func (this *UserModel) CanDelete(user IUserModel) bool {
	return this.ItemRoleModel.CanWrite(user)
}

//  save时对密码进行加密
func (this *UserModel) CopyTo(user IUserModel, bean interface{}) error {
	data := bean.(*UserModel)
	data.Nickname = this.Nickname
	data.Telphone = this.Telphone
	data.Password = this.Self().(IUserModel).EncodePwd(this.Password)
	return nil
}

// 加密算法，@overable
func (this *UserModel) EncodePwd(password string) string {
	h := md5.New()
	h.Write([]byte(this.Telphone + password))
	hexText := make([]byte, 32)
	hex.Encode(hexText, h.Sum(nil))
	return string(hexText)
}

// 用户是否存在，@overable
func (this *UserModel) Exist(telphone string) bool {
	if telphone == "" {
		return false
	}
	ok, _ := this.Db.Where("telphone=?", telphone).Get(this.Self())
	return ok
}

// 通过telphone获取取用户,@overable
func (this *UserModel) GetByTelphone(telphone string) bool {
	if telphone == "" {
		return false
	}
	bean := this.Self().(IUserModel)
	ok, _ := this.Db.Where("telphone=?", telphone).Get(bean)
	return ok
}

// 比较密码是否正确,@overable
func (this *UserModel) JudgeChpwdCode(code string) error {
	if this.Password != this.Self().(IUserModel).EncodePwd(code) {
		return errors.New("旧密码错误")
	}
	return nil
}

// @path
func (this *UserModel) MakeLoginData(action string) (*LoginData, error) {
	access := NewAccessToken(this.Id, this.Ctx.Context)
	access.Action = action
	if data, ok := this.Self().(IDataDetail); ok {
		data.GetDetail(this.Self().(IUserModel))
		return &LoginData{access, data}, nil
	}
	if _, err := this.Db.InsertOne(access); err != nil {
		return nil, err
	}
	switch action {
	case "login":
		// 下线同种设备
		this.Db.Exec("update access_token set expired_at='1993-03-07' where id!=? and user_id=? and device=?", access.Id, access.UserId, access.Device)
	case "chpwd":
		// 注销所有
		this.Db.Exec("update access_token set expired_at='1993-03-07' where id!=? and user_id=?", access.Id, access.UserId)
	}
	return &LoginData{access, this.Self()}, nil
}

// 注销
func (this *UserModel) Logout() (sql.Result, error) {
	user_id := this.GetCurrentUserId(this.Id)
	if user_id > 0 {
		return this.Db.Exec("update access_token set expired_at='1993-03-07' where user_id=?", user_id)
	}
	return nil, errors.New("需要user_id")
}

// 登录
func (this *UserModel) Login(telphone, password string) (*LoginData, error) {
	bean := this.Self().(IUserModel)
	// 通过手机号查用户
	ok, err := this.Db.Where("telphone=?", telphone).Get(bean)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("用户不存在")
	}
	// 匹配密码
	if this.Password == bean.EncodePwd(password) {
		return this.MakeLoginData("login")
	}
	return nil, errors.New("密码不正确")
}

// 修改密码
func (this *UserModel) ChangePassword(rbody *RegistorBody) (*LoginData, error) {
	bean := this.Self().(IUserModel)
	if len(rbody.Password) < 6 || len(rbody.Password) > 32 {
		return nil, errors.New("请输入6~32位密码")
	}
	if ok := bean.GetByTelphone(rbody.Telphone); !ok {
		return nil, errors.New("用户不存在")
	}
	if err := this.Self().(IUserModel).JudgeChpwdCode(rbody.Code); err != nil {
		return nil, err
	}
	this.Password = bean.EncodePwd(rbody.Password)
	_, err := this.Db.ID(this.Id).Cols("password").Update(bean)
	if err != nil {
		return nil, err
	}
	return this.MakeLoginData("chpwd")
}

// 我的信息
func (this *UserModel) MineInfo() (interface{}, error) {
	if data, ok := this.GetUser(); ok {
		if bean, ok := data.(IDataDetail); ok {
			bean.GetDetail(data)
		}
		return data, nil
	}
	return nil, NeedAuthError
}

func UserMW(user IUserModel) gin.HandlerFunc {
	user.SetSelf(user)
	return func(c *gin.Context) {
		// 当前登录用户数据
		token := c.Query("access_token")
		if token != "" {
			now := time.Now()
			bean := user.New()
			ok, _ := config.Db.Where("id in (select user_id from access_token where token=? and expired_at>?)", token, now).Get(bean)
			if ok {
				c.Set("user", bean)
			}
		}
	}
}

// 登录时POST数据结构
type LoginBody struct {
	Telphone string `gev:"电话号码" json:"telphone" xorm:"varchar(32) unique(telphone) not null"`
	Password string `gev:"密码" json:"password" xorm:"varchar(64)"`
}

// 登录返回数据结构
type LoginData struct {
	Access *AccessToken `json:"access,omitempty" xorm:""`
	User   interface{}  `json:"user,omitempty" xorm:""`
}
