package ginhelper

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// checkHandler checks if handler is valid
// handler must be a function
// handler's first argument must be *gin.Context
// handler's second argument must be a struct
// handler's last return value must be error
// handler's first return value must be a pointer
// example:
//   - func(c *gin.Context) error
//   - func(c *gin.Context, req any) error
//   - func(c *gin.Context) (*resp, error)
//   - func(c *gin.Context, req any) (*resp, error)
func checkHandler(handler any) {
	v := reflect.ValueOf(handler)
	t := v.Type()

	if t.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	if t.NumIn() == 0 || t.NumIn() > 2 {
		panic("handler must have 1 or 2 arguments")
	}
	if t.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic("handler's first argument must be *gin.Context")
	}
	if t.NumIn() == 2 && t.In(1).Kind() != reflect.Struct {
		panic("handler's second argument must be a struct")
	}

	if t.NumOut() == 0 || t.NumOut() > 2 {
		panic("handler must have 1 or 2 return values")
	}
	if !t.Out(t.NumOut() - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("handler's last return value must be error")
	}
	if t.NumOut() == 2 && t.Out(0).Kind() != reflect.Ptr {
		panic("handler's first return value must be a pointer")
	}
}
