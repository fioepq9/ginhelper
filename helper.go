package ginhelper

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Helper struct {
	bindings            map[string]binding.Binding
	bindingUri          binding.BindingUri
	bindingValidator    binding.StructValidator
	bindingErrorHandler func(*gin.Context, error)
	successHandler      func(*gin.Context, any)
	errorHandler        func(*gin.Context, error)
}

// New returns a new Engine instance
//
//	Notes: binding.Validator will be set to nil
func New(options ...func(*Helper)) *Helper {
	e := &Helper{
		bindings: map[string]binding.Binding{
			"json":   binding.JSON,
			"form":   binding.Form,
			"header": binding.Header,
		},
		bindingUri:       binding.Uri,
		bindingValidator: binding.Validator,
		bindingErrorHandler: func(c *gin.Context, err error) {
			c.AbortWithError(http.StatusBadRequest, err)
		},
		errorHandler: func(c *gin.Context, err error) {
			c.AbortWithError(http.StatusInternalServerError, err)
		},
		successHandler: func(c *gin.Context, resp any) {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data":    resp,
			})
		},
	}

	// disable gin binding validator
	binding.Validator = nil

	for _, o := range options {
		o(e)
	}

	return e
}

func WithBinding(tag string, binding binding.Binding) func(*Helper) {
	return func(e *Helper) {
		e.bindings[tag] = binding
	}
}

func WithBindingUri(binding binding.BindingUri) func(*Helper) {
	return func(e *Helper) {
		e.bindingUri = binding
	}
}

func WithBindingValidator(binding binding.StructValidator) func(*Helper) {
	return func(e *Helper) {
		e.bindingValidator = binding
	}
}

func WithBindingErrorHandler(handler func(*gin.Context, error)) func(*Helper) {
	return func(e *Helper) {
		e.bindingErrorHandler = handler
	}
}

func WithErrorHandler(handler func(*gin.Context, error)) func(*Helper) {
	return func(e *Helper) {
		e.errorHandler = handler
	}
}

func WithSuccessHandler(handler func(*gin.Context, any)) func(*Helper) {
	return func(e *Helper) {
		e.successHandler = handler
	}
}

func (e *Helper) GET(router gin.IRoutes, path string, handler any) gin.IRoutes {
	return e.handle(router, http.MethodGet, path, handler)
}

func (e *Helper) POST(router gin.IRoutes, path string, handler any) gin.IRoutes {
	return e.handle(router, http.MethodPost, path, handler)
}

func (e *Helper) handle(router gin.IRoutes, method string, path string, handler any) gin.IRoutes {
	checkHandler(handler)
	v := reflect.ValueOf(handler)
	t := v.Type()

	request := func(c *gin.Context) ([]reflect.Value, error) {
		in := make([]reflect.Value, 0, t.NumIn())
		in = append(in, reflect.ValueOf(c))
		if t.NumIn() == 2 {
			hasTags := make(map[string]bool)
			hasUriTag := false
			reqV := reflect.New(t.In(1).Elem())
			reqT := reqV.Elem().Type()
			for i := 0; i < reqT.NumField(); i++ {
				for tag := range e.bindings {
					if hasTags[tag] {
						continue
					}
					if _, ok := reqT.Field(i).Tag.Lookup(tag); ok {
						hasTags[tag] = true
					}
				}
				if _, ok := reqT.Field(i).Tag.Lookup("uri"); ok {
					hasUriTag = true
				}
			}
			// bind uri
			if hasUriTag {
				m := make(map[string][]string)
				for _, v := range c.Params {
					m[v.Key] = []string{v.Value}
				}
				err := e.bindingUri.BindUri(m, reqV.Interface())
				if err != nil {
					return nil, err
				}
			}
			// bind other tags
			for tag := range hasTags {
				err := c.ShouldBindWith(reqV.Interface(), e.bindings[tag])
				if err != nil {
					return nil, err
				}
			}
			err := e.bindingValidator.ValidateStruct(reqV.Elem().Interface())
			if err != nil {
				return nil, err
			}
			in = append(in, reqV)
		}
		return in, nil
	}

	return router.Handle(method, path, func(c *gin.Context) {
		in, err := request(c)
		if err != nil {
			e.bindingErrorHandler(c, err)
			return
		}
		out := v.Call(in)
		var resp any
		if len(out) == 1 {
			if errVal := out[0].Interface(); errVal != nil {
				err = errVal.(error)
			}
		} else {
			resp = out[0].Interface()
			if errVal := out[1].Interface(); errVal != nil {
				err = errVal.(error)
			}
		}
		if err != nil {
			e.errorHandler(c, err)
			return
		}
		e.successHandler(c, resp)
	})
}
