package models

import (
	"time"

	"github.com/inu1255/gev2/config"

	"github.com/gin-gonic/gin"
	"github.com/golibs/uuid"
)

type AccessToken struct {
	Model     `xorm:"extends"`
	Token     string    `gev:"身份密钥" json:"token,omitempty" xorm:"index"`
	ExpiredAt time.Time `gev:"过期时间" json:"expired_at,omitempty" xorm:"index"`
	UserId    int       `json:"-" xorm:""`
	Ip        string    `json:"-" xorm:""`
	UA        string    `json:"-" xorm:""`
	Device    string    `json:"-" xorm:""`
	Uuid      string    `json:"-" xorm:""`
	Action    string    `json:"-" xorm:""`
}

func (this *AccessToken) ReadContextInfo(c *gin.Context) {
	UA := c.Request.Header.Get("User-Agent")
	this.Ip = c.ClientIP()
	this.UA = UA
	this.Device = c.Request.Header.Get("X-DEVICE")
	this.Uuid = c.Request.Header.Get("X-UUID")
}

func NewAccessToken(user_id int, c *gin.Context) *AccessToken {
	token := &AccessToken{
		UserId:    user_id,
		Token:     uuid.Rand().Hex(),
		ExpiredAt: time.Now().Add(time.Duration(config.TokenExpire) * time.Second),
	}
	token.ReadContextInfo(c)
	return token
}
