package models

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/inu1255/gev2/libs"
)

type IFileModel interface {
	Copy(file *FileModel)
}

type FileModel struct {
	ItemRoleModel `xorm:"extends"`
	Ext           string `json:"ext,omitempty" xorm:"" gev:"文件后缀"`
	Place         string `json:"-" xorm:""`
	Filename      string `json:"filename,omitempty" xorm:"" gev:""`
	MD5           string `json:"-" xorm:"" gev:""`
	Url           string `json:"url" xorm:"" gev:"文件地址,需加上host,如http://www.tederen.com:8017/"`
}

func (this *FileModel) TableName() string {
	return "file"
}
func (this *FileModel) Copy(file *FileModel) {
	this.ItemRoleModel = file.ItemRoleModel
	this.Ext = file.Ext
	this.Place = file.Place
	this.Filename = file.Filename
	this.MD5 = file.MD5
	this.Url = file.Url
}

// @path
func (this *FileModel) Save(self IBody) (interface{}, error) {
	return this.ItemRoleModel.Save(self)
}

// @path
func (this *FileModel) SaveAll(self []IBody) (interface{}, error) {
	return this.ItemRoleModel.SaveAll(self)
}

// 上传文件
func (this *FileModel) Upload(file *multipart.FileHeader) (interface{}, error) {
	user, _ := this.GetUser()
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	return FileUpload(this.Db, this, file.Filename, src, user)
}

/**
 * 上传base64文件
 * @param  filename 文件名 query
 */
func (this *FileModel) UploadBase64(filename string) (interface{}, error) {
	// 上传者
	user, _ := this.GetUser()
	src, err := ioutil.ReadAll(this.Ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	if filename == "" {
		filename = time.Now().Format("2006-01-02 03:04:05")
	}
	body := Bytes2str(src)
	index := strings.Index(body, ";base64,")
	if index < 4 {
		return nil, errors.New("没有找到;base64,")
	}
	if body[index-4:index] == "jpeg" {
		filename += ".jpg"
	} else if body[index-3:index] == "png" {
		filename += ".png"
	} else if i := strings.LastIndex(body[:index], "/"); i >= 0 {
		filename += "." + body[i+1:index]
	} else {
		return nil, errors.New(body[:index])
	}
	src = src[index+8:]
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, err = base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return FileUpload(this.Db, this, filename, bytes.NewReader(dst), user)
}

// 导出csv表格文件
func (this *FileModel) ExportCsv(tables [][]string) (interface{}, error) {
	this.Ctx.Header("Content-Type", "application/octet-stream")
	this.Ctx.Header("Content-Disposition", "attachment; filename=表格.csv")
	err := libs.SimpleWriteCsv(this.Ctx.Writer, tables)
	return nil, err
}

// 导出xlsx表格文件
func (this *FileModel) ExportXlsx(tables [][]string) (interface{}, error) {
	this.Ctx.Header("Content-Type", "application/octet-stream")
	this.Ctx.Header("Content-Disposition", "attachment; filename=表格.xlsx")
	err := libs.SimpleWriteExcel(this.Ctx.Writer, tables)
	return nil, err
}

func FileGetExt(filename string) string {
	index := strings.LastIndex(filename, ".")
	if index >= 0 {
		return strings.ToLower(filename[index+1:])
	}
	return ""
}

func FileUpload(Db *xorm.Session, bean IFileModel, filename string, src io.Reader, user IUserModel) (interface{}, error) {
	var err error
	f := new(FileModel)
	// 创建用户文件夹
	uid := "0"
	if user != nil {
		uid = strconv.Itoa(user.GetId())
		f.OwnerId = user.GetId()
	}
	dir := strings.Join([]string{"upload", uid}, "/")
	err = os.MkdirAll(dir, 0755)

	bs, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	h := md5.New()
	h.Write(bs)
	f.MD5 = hex.EncodeToString(h.Sum(nil))
	// 保存文件
	f.Place = strings.Join([]string{dir, "/", f.MD5}, "")
	if _, err = os.Stat(f.Place); err != nil {
		err = ioutil.WriteFile(f.Place, bs, 0644)
		if err != nil {
			return nil, err
		}
	}
	//  保存文件
	f.Ext = FileGetExt(filename)
	f.Filename = filename
	f.Url = f.Place
	if ok, _ := Db.Where("place=? and owner_id=?", f.Place, f.OwnerId).Get(bean); ok {
		return bean, nil
	}
	bean.Copy(f)
	Db.InsertOne(bean)
	return bean, nil
}
