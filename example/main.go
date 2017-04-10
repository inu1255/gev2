package main

import (
	"log"
	"mime/multipart"

	"github.com/inu1255/gev2"
	"github.com/inu1255/go-swagger/core"
)

type Data struct {
	A string `json:"a,omitempty" xorm:"" gev:""`
}

type Service struct {
	gev2.BaseService
	hide     int
	Show     int
	hideData Data
	ShowData Data
}

/**
 * 介绍TestA
 * 哈哈
 * @param  abc 介绍abc
 * @return map[string]string{"a": abc, "b": "2"}
 */
func (this *Service) TestA(file *multipart.FileHeader, abc string) (interface{}, error) {
	f, err := file.Open()
	log.Println(file.Filename, f, err)
	return map[string]string{"a": abc, "b": "2"}, nil
}

// 测试B
// 介绍
// @param self 	自动当成自己 	path
// @param b 	坏处 		query
func (this *Service) TestB(self []interface{}, b string) ([]interface{}, error) {
	log.Printf("%T,%v", self, self)
	for _, item := range self {
		log.Printf("%T,%v", item, item)
	}
	return []interface{}{b}, nil
}

// @path pp/ a/
// @param abc abc path
func (this *Service) PostA(self gev2.IService, abc string) (interface{}, error) {
	log.Printf("%T,%v", self, self)
	return self, nil
}

func (this *Service) PostB(data Data, a, b string) (Data, error) {
	return data, nil
}

type NotService struct{}

func (this *NotService) TestA(abc string) (interface{}, error) {
	return map[string]string{"a": abc, "b": "2"}, nil
}

func main() {
	core.CopySwagger()
	gev2.Bind("test", &Service{})
	gev2.Bind("serv", &NotService{})
	gev2.App.Static("api", "api")
	gev2.Run(":8019")
}
