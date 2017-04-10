package gev2

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context, data interface{}) {
	if data != nil {
		c.IndentedJSON(200, gin.H{"code": 0, "data": data})
	}
}

func Err(c *gin.Context, code int, err error) {
	msg := err.Error()
	if code == 0 {
		table := Str2bytes(msg)
		count := len(table)
		if count > 32 {
			count = 32
		}
		for i := 0; i < count; i++ {
			code += int(table[i])
		}
	}
	Log.Println("code:\033[41;37m", code, "\033[0m msg:", msg)
	c.IndentedJSON(200, gin.H{"code": code, "msg": msg})
}

type iApiError interface {
	error
	Code() int
}

func Api(c *gin.Context, data interface{}, err error) {
	if err != nil {
		if v, ok := err.(iApiError); ok {
			Err(c, v.Code(), err)
		} else {
			Err(c, 0, err)
		}
		return
	}
	Ok(c, data)
}

func NeedAuth(c *gin.Context) (interface{}, bool) {
	if user, ok := c.Get("user"); ok {
		return user, true
	} else {
		Err(c, 1255, errors.New("需要登录"))
		return nil, false
	}
}

// func NeedAuthRole(c *gin.Context, role string) (interface{}, bool) {
// 	if user, ok := c.Get("user"); ok {
// 		if v, ok := user.(IUserRoleModel); ok && v.GetRole() == role {
// 			return user, true
// 		}
// 		Err(c, 1256, errors.New("需要"+role+"权限"))
// 		return nil, false
// 	} else {
// 		Err(c, 1255, errors.New("需要登录"))
// 		return nil, false
// 	}
// }
