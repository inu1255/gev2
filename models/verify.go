package models

import (
	"errors"
	"math/rand"
	"time"
)

type IVerifyModel interface {
	IModel
	SendCode(title, code string) error
	RandCode() string
	JudgeCode(title, code string) error
	CanSend() bool
}

// 验证码模型
type VerifyModel struct {
	Model `xorm:"extends"`
	Title string `json:"title,omitempty" xorm:"" gev:"手机号/邮箱等"`
	Code  string `json:"code,omitempty" xorm:"not null" gev:"验证码"`
	Rest  int    `json:"rest,omitempty" xorm:"not null default 10" gev:"剩余验证次数"`
}

func (this *VerifyModel) Disable() {
	if this.Id > 0 {
		this.Rest = 0
		this.Db.ID(this.Id).Cols("rest").Update(this.Self())
	}
}

func (this *VerifyModel) JudgeCode(title, code string) error {
	Db := this.Db
	bean := this.Self().(IVerifyModel)
	ok, _ := Db.Where("title=?", title).Get(bean)
	if !ok {
		return errors.New("尚未发送验证码")
	}
	if this.UpdateAt.Add(10 * time.Minute).Before(time.Now()) {
		return errors.New("验证码已过期")
	}
	if this.Rest < 1 {
		return errors.New("验证码已失效")
	}
	if this.Code != code {
		this.Rest--
		Db.ID(this.Id).Cols("rest").Update(bean)
		return errors.New("验证码错误")
	}
	return nil
}

func (this *VerifyModel) SendCode(title, code string) error {
	Log.Println(title, "=>", code)
	this.Title = title
	this.Code = code
	this.Rest = 10
	return nil
}

func (this *VerifyModel) RandCode() string {
	code := make([]byte, 4)
	for i := 0; i < 4; i++ {
		code[i] = byte('0' + rand.Intn(10))
	}
	return string(code)
}

func (this *VerifyModel) CanSend() bool {
	return this.UpdateAt.Add(time.Minute).After(time.Now())
}

// 发送验证码
// title:要发送的手机号或邮箱
// @path /send
func (this *VerifyModel) NewVerifyCode(title string) (interface{}, error) {
	Db := this.Db
	bean := this.Self().(IVerifyModel)
	ok, _ := Db.Where("title=?", title).Get(bean)
	if bean.CanSend() {
		return nil, errors.New("发送太频繁")
	}
	err := bean.SendCode(title, bean.RandCode())
	if err != nil {
		return nil, err
	}
	// Log.Printf("class:%T - %v", bean, bean)
	var number int64
	if ok {
		number, err = Db.Where("title=?", title).Update(bean)
	} else {
		number, err = Db.InsertOne(bean)
	}
	return number, err
}

func init() {
	rand.Seed(time.Now().Unix())
}
