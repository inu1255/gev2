package main

import (
	"github.com/inu1255/gev2"
	"github.com/inu1255/gev2/config"
	. "github.com/inu1255/gev2/models"
	"github.com/inu1255/go-swagger/core"
)

func main() {
	config.SetDb("mysql", "root:199337@/gev2?parseTime=true")
	core.CopySwagger()
	gev2.App.Static("api", "api")
	gev2.App.Use(UserMW(new(UserRoleModel)))
	gev2.App.Use(gev2.CrossDomainMW())

	UserVerify = NewModel(new(VerifyModel)).(IVerifyModel)
	gev2.Bind("address", new(AddressModel))
	gev2.Bind("file", new(FileModel))
	gev2.Bind("item", new(ItemModel))
	gev2.Bind("item.owner", new(ItemOwnerModel))
	gev2.Bind("item.role", new(ItemRoleModel))
	gev2.Bind("user", new(UserModel))
	gev2.Bind("user.regist", new(UserRegistModel))
	gev2.Bind("user.role", new(UserRoleModel))
	gev2.Bind("verify", new(VerifyModel))
	gev2.Bind("verify.mail", new(VerifyMailModel))
	gev2.Run(":8019")
}
