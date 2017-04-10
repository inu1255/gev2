package models

import "errors"

var (
	UserVerify IVerifyModel = NewModel(new(VerifyModel)).(IVerifyModel)
)

type IUserRegistModel interface {
	IUserRoleModel
	RegistorJudge(rbody *RegistorBody) error
}

type UserRegistModel struct {
	UserRoleModel `xorm:"extends"`
}

// 默认注册数据结构
type RegistorBody struct {
	Code     string `json:"code,omitempty" xorm:""`
	Telphone string `json:"telphone,omitempty" xorm:""`
	Password string `json:"password,omitempty" xorm:""`
}

func (this *UserRegistModel) GetVerify() IVerifyModel {
	return this.NewModel(UserVerify.New().(IModel)).(IVerifyModel)
}

func (this *UserRegistModel) JudgeChpwdCode2(code string) error {
	if err := this.GetVerify().JudgeCode(this.Telphone, code); err != nil {
		return err
	}
	return nil
}

func (this *UserRegistModel) RegistorJudge(rbody *RegistorBody) error {
	bean := this.Self()
	ok, _ := this.Db.Where("telphone=?", rbody.Telphone).Get(bean)
	if ok {
		return errors.New("账号已注册")
	}
	if len(rbody.Password) < 6 || len(rbody.Password) > 32 {
		return errors.New("请输入6~32位密码")
	}
	return this.GetVerify().JudgeCode(rbody.Telphone, rbody.Code)
}

// 注册
func (this *UserRegistModel) Registor(rbody *RegistorBody) (*LoginData, error) {
	if UserVerify == nil {
		return nil, ApiErr("系统禁止注册", 0)
	}
	bean := this.Self().(IUserRegistModel)
	if err := bean.RegistorJudge(rbody); err != nil {
		return nil, err
	}
	this.Telphone = rbody.Telphone
	this.Password = bean.EncodePwd(rbody.Password)
	_, err := this.Db.InsertOne(bean)
	if err != nil {
		return nil, err
	}
	return this.MakeLoginData("regist")
}
