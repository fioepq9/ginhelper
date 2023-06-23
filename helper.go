package ginhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var H *helper

type helper struct {
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
func init() {
	h := &helper{
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

	H = h
}

func (h *helper) WithBinding(tag string, binding binding.Binding) *helper {
	H.bindings[tag] = binding
	return H
}

func (h *helper) WithBindingUri(binding binding.BindingUri) *helper {
	H.bindingUri = binding
	return H
}

func (h *helper) WithBindingValidator(binding binding.StructValidator) *helper {
	H.bindingValidator = binding
	return H
}

func (h *helper) WithBindingErrorHandler(handler func(*gin.Context, error)) *helper {
	H.bindingErrorHandler = handler
	return H
}

func (h *helper) WithErrorHandler(handler func(*gin.Context, error)) *helper {
	H.errorHandler = handler
	return H
}

func (h *helper) WithSuccessHandler(handler func(*gin.Context, any)) *helper {
	H.successHandler = handler
	return H
}

func (h *helper) Router(routes gin.IRoutes) *router {
	return &router{
		helper: h,
		routes: routes,
	}
}
