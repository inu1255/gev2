package gev2

import "reflect"

type IClass interface {
	Self() IClass
	SetSelf(IClass)
	New() IClass
}

type Class struct {
	self IClass `xorm:"-"`
}

// IClass 接口
func (this *Class) Self() IClass {
	// if this.self == nil {
	// 	return this
	// }
	return this.self
}
func (this *Class) SetSelf(self IClass) {
	this.self = self
}
func (this *Class) New() IClass {
	// if this.self == nil {
	// 	this.self = this
	// }
	class := reflect.New(reflect.TypeOf(this.self).Elem()).Interface().(IClass)
	class.SetSelf(class)
	// Log.Printf("class: %T - %v", class, class)
	return class
}
