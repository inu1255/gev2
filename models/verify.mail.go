package models

import (
	"errors"
	"regexp"

	gomail "gopkg.in/gomail.v2"
)

var (
	_email_user     = "uniwise@aliyun.com"
	_email_password = "uniwise87"
	_email_title    = "云央科技"
	_email_dialer   = gomail.NewDialer("smtp.aliyun.com", 25, _email_user, _email_password)
	_email_regex    = regexp.MustCompile(`[.\d\w]+@[\d\w]+\.[\d\w]+`)
)

// VerifyMailModel Entity
type VerifyMailModel struct {
	VerifyModel `xorm:"extends"`
}

func (this *VerifyMailModel) SendCode(title, code string) error {
	if !_email_regex.Match(Str2bytes(title)) {
		return errors.New("邮箱格式不正确")
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", _email_user, _email_title)
	m.SetHeader("To", title)
	m.SetHeader("Subject", "<"+_email_title+">验证码")
	m.SetBody("text/html", "<html>您的验证码是[ "+code+" ]，请勿告诉他人</html>")
	this.VerifyModel.SendCode(title, code)
	return _email_dialer.DialAndSend(m)
}
