package gev2

import (
	"os"
	"strings"

	"github.com/inu1255/go-swagger/core"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/inu1255/gev2/config"
)

var (
	App       = gin.New()
	Log       = config.NewLogger("gev2")
	_gev_path = ""
	ccFs      = make([]ccF, 0, 2)
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	App.Use(gin.Logger())
	App.Use(gin.Recovery())
	_gev_path = config.GetPkgPath()
	core.CopySwagger()
}

func Description(info ...string) {
	Swag.Info.Add(info...)
}

func Run(host string) {
	if host == "" {
		host = ":8017"
	}
	Swag.WriteJson("api/swagger.json")

	// config.Db.ShowSQL(true)
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		Log.Println("go run 不启动热更新", os.Args[0])
		App.Run(host)
	} else {
		AutoRestart()
		Server := endless.NewServer(host, App)
		Server.ListenAndServe()
	}
}
