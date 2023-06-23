package ginhelper

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type router struct {
	routes gin.IRoutes
	helper *helper
}

func (r *router) GET(path string, handler any) *router {
	return r.Handle(http.MethodGet, path, handler)
}

func (r *router) POST(path string, handler any) *router {
	return r.Handle(http.MethodPost, path, handler)
}

func (r *router) Handle(method string, path string, handler any) *router {
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
				for tag := range r.helper.bindings {
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
				err := r.helper.bindingUri.BindUri(m, reqV.Interface())
				if err != nil {
					return nil, err
				}
			}
			// bind other tags
			for tag := range hasTags {
				err := c.ShouldBindWith(reqV.Interface(), r.helper.bindings[tag])
				if err != nil {
					return nil, err
				}
			}
			err := r.helper.bindingValidator.ValidateStruct(reqV.Elem().Interface())
			if err != nil {
				return nil, err
			}
			in = append(in, reqV)
		}
		return in, nil
	}

	r.routes.Handle(method, path, func(c *gin.Context) {
		in, err := request(c)
		if err != nil {
			r.helper.bindingErrorHandler(c, err)
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
			r.helper.errorHandler(c, err)
			return
		}
		r.helper.successHandler(c, resp)
	})

	return r
}
