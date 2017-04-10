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
	Db, _       = xorm.NewEngine("mysql", "root:199337@/youyue?parseTime=true&loc=Asia%2FShanghai")
	TokenExpire = 86400
)

func SetDb(driverName, dataSourceName string) {
	Db, _ = xorm.NewEngine(driverName, dataSourceName)
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
