package config

import (
	"log"
	"os"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	Db          *xorm.Engine
	session     *xorm.Session
	TokenExpire = 86400
)

func SetDb(driverName, dataSourceName string) {
	var err error
	Db, err = xorm.NewEngine(driverName, dataSourceName)
	if err == nil {
		session = Db.NewSession()
	}
}

// 返回只读 session
func Rdb() *xorm.Session {
	return session
}

func NewLogger(name string) *log.Logger {
	return log.New(os.Stdout, "[ "+name+" ] ", log.Ltime|log.Lshortfile)
}

func GetPkgPath() string {
	_, file, _, _ := runtime.Caller(1)
	if index := strings.LastIndex(file, "/"); index > 0 {
		return file[:index]
	}
	return file
}
