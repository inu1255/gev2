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
	TokenExpire = 86400
)

func SetDb(driverName, dataSourceName string) {
	Db, _ = xorm.NewEngine(driverName, dataSourceName)
}

// 返回只读 session
func Rdb() *xorm.Session {
	session := Db.NewSession()
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
