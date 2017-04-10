package gev2

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/inu1255/gev2/config"
	swaggin "github.com/inu1255/go-swagger/gin"
)

var (
	App       = gin.New()
	Swag      = swaggin.Swag
	Log       = config.NewLogger("gev2")
	_gev_path = ""
	// UserVerify   IVerifyModel
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	App.Use(gin.Logger())
	App.Use(gin.Recovery())
	_gev_path = config.GetPkgPath()
	CopySwagger()
}

func CopySwagger() {
	if info, err := os.Stat("api"); err != nil || !info.IsDir() {
		cmd := exec.Command("cp", "-R", _gev_path+"/api", "api")
		err := cmd.Start()
		if err != nil {
			Log.Println(err)
		}
	} else {
		Log.Println("swagger文件夹已经存在")
	}
}

func Description(info ...string) {
	Swag.Info.Add(info...)
}

func Run(host string) {
	if host == "" {
		host = ":8017"
	}
	Swag.WriteJson("api/swagger.json")

	config.Db.ShowSQL(true)
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		Log.Println("go run 不启动热更新", os.Args[0])
		App.Run(host)
	} else {
		AutoRestart()
		Server := endless.NewServer(host, App)
		Server.ListenAndServe()
	}
}
