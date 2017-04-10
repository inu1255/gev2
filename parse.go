package gev2

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/inu1255/go-swagger/core"
	swaggin "github.com/inu1255/go-swagger/gin"
)

type MethodParser struct {
	method   reflect.Method
	funcDecl *ast.FuncDecl
	mmap     map[string]*core.Param `json:"mmap,omitempty" gev:"从注释中解析的参数"`
	info     []string               `json:"info,omitempty" gev:"从注释中解析的函数介绍"`
	path     string                 `json:"path,omitempty" gev:"接口uri"`
}

func (this *MethodParser) ReadComment() {
	if this.funcDecl == nil {
		return
	}
	doc := this.funcDecl.Doc
	this.mmap = make(map[string]*core.Param)
	this.info = make([]string, 0)
	if doc == nil {
		return
	}
	lines := strings.Split(doc.Text(), "\n")
	for _, line := range lines {
		var ss []string
		ss = strings.Split(line, "@param")
		if len(ss) > 1 {
			ss = strings.Fields(ss[1])
			count := len(ss)
			if count > 0 {
				param := new(core.Param)
				param.Name = ss[0]
				if count > 1 {
					param.Description = ss[1]
				}
				if count > 2 {
					param.In = ss[2]
				}
				this.mmap[param.Name] = param
				// Log.Println(param)
			}
			continue
		}
		ss = strings.Split(line, "@path")
		if len(ss) > 1 {
			s := strings.Trim(ss[1], " /")
			if s == "" {
				this.path = " "
			} else {
				this.path = "/" + s
			}
			continue
		}
		s := strings.Trim(line, " *")
		if s != "" {
			this.info = append(this.info, s)
		}
	}
}

func (this *MethodParser) GetMethodParams() []*core.Param {
	count := this.method.Type.NumIn()
	params := make([]*core.Param, count)
	if this.funcDecl != nil {
		if this.mmap == nil {
			this.ReadComment()
		}
		mmap := this.mmap
		n := 1
		for _, param := range this.funcDecl.Type.Params.List {
			// Log.Println(param.Names, param.Tag, param.Type)
			for _, item := range param.Names {
				var ok bool
				if params[n], ok = mmap[item.Name]; !ok {
					params[n] = new(core.Param)
				}
				if params[n].In == "" {
					if count < 3 {
						params[n].In = "path"
						params[n].Required = true
					} else {
						params[n].In = "query"
					}
				}
				params[n].Name = item.Name
				params[n].Type = "string"
				n++
			}
		}
	}
	return params
}

func (this *MethodParser) ParseToRouter(router swaggin.ISwagRouter) {
	if this.info == nil {
		this.ReadComment()
	}
	info := this.info
	typ := this.method.Type.Out(0)
	switch typ.Kind() {
	case reflect.Interface:
		info = append(info, "<b>返回数据结构不确定，以下仅作参考</b>")
		router.Data(reflect.New(this.method.Type.In(0)).Elem().Interface())
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Interface {
			info = append(info, "<b>返回数据结构不确定，以下仅作参考</b>")
			router.Data(reflect.New(reflect.SliceOf(this.method.Type.In(0))).Elem().Interface())
		} else {
			router.Data(reflect.New(typ).Elem().Interface())
		}
	default:
		router.Data(reflect.New(typ).Elem().Interface())
	}
	router.Info(info...)
}

func NewMethodParser(m reflect.Method) *MethodParser {
	this := &MethodParser{}
	this.method = m
	this.funcDecl = GetFuncDecl(m)
	return this
}

func GetFuncDecl(m reflect.Method) (funcDecl *ast.FuncDecl) {
	typ := m.Type.In(0).Elem()
	structName := typ.Name()
	structPkg := typ.PkgPath()
	methodName := m.Name
	var mpkg map[string]*ast.Package
	var err error
	fset := token.NewFileSet()
	if strings.Contains(structPkg, "/") {
		mpkg, err = parser.ParseDir(fset, path.Join(os.Getenv("GOPATH"), "src", structPkg), nil, parser.ParseComments)
	} else {
		mpkg, err = parser.ParseDir(fset, ".", nil, parser.ParseComments)
	}
	if err != nil {
		return nil
	}
	for _, pkg := range mpkg {
		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				if fun, ok := decl.(*ast.FuncDecl); ok {
					if fun.Name.Name == methodName {
						if len(fun.Recv.List) > 0 {
							if starExpr, ok := fun.Recv.List[0].Type.(*ast.StarExpr); ok {
								if index, ok := starExpr.X.(*ast.Ident); ok {
									if index.Name == structName {
										return fun
									}
								}
							}
						}
					}
				}
			}
		}
	}
	WalkFields(typ, func(field reflect.StructField) bool {
		if field.Anonymous {
			WalkMethods(field.Type, func(method reflect.Method) bool {
				if method.Name == methodName {
					funcDecl = GetFuncDecl(method)
					return true
				}
				return false
			})
		}
		return false
	})
	return
}

func WalkFields(typ reflect.Type, call func(reflect.StructField) bool) {
	switch typ.Kind() {
	case reflect.Interface, reflect.Ptr:
		WalkFields(typ.Elem(), call)
	case reflect.Struct:
		numField := typ.NumField()
		for i := 0; i < numField; i++ {
			field := typ.Field(i)
			// Log.Println("field", field.Name, field.PkgPath, field.Anonymous, field.Type)
			if call(field) {
				break
			}
		}
	}
}

func WalkMethods(typ reflect.Type, call func(reflect.Method) bool) {
	switch typ.Kind() {
	case reflect.Interface, reflect.Ptr:
		numMethod := typ.NumMethod()
		for i := 0; i < numMethod; i++ {
			method := typ.Method(i)
			// typ := method.Type.In(0).Elem()
			// Log.Println("method", typ.PkgPath(), typ.Name(), method.Name)
			if call(method) {
				break
			}
		}
	case reflect.Struct:
		WalkMethods(reflect.PtrTo(typ), call)
	}
}
