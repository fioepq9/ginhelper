package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fioepq9/ginhelper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"go.uber.org/multierr"
)

// http localhost:8080/echo/1234 message==hello token:1234
type EchoRequest struct {
	ID      int    `uri:"id" binding:"required"`
	Message string `form:"message" binding:"required"`
	Token   string `header:"token" binding:"required"`
}

func EchoController(c *gin.Context, req *EchoRequest) error {
	fmt.Printf("%+v\n", req)
	return NewResponse(ResponseCodeNotFound, "not found id")
}

// success: http localhost:8080/create username=foo@bar.com password=qwer
// fail: http localhost:8080/create username=foo password=qwer
type CreateRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateResponse struct {
	ID string `json:"id"`
}

func CreateController(c *gin.Context, req *CreateRequest) (*CreateResponse, error) {
	fmt.Printf("%+v\n", req)
	resp := &CreateResponse{
		ID: "1234",
	}
	return resp, nil
}

type Validator struct {
	validate *validator.Validate
	trans    *ut.Translator
}

func NewValidator() *Validator {
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	validate.SetTagName("binding")
	zh_translations.RegisterDefaultTranslations(validate, trans)
	return &Validator{
		validate: validate,
		trans:    &trans,
	}
}

func (v *Validator) ValidateStruct(s any) error {
	err := v.validate.Struct(s)
	if err != nil {
		var merr error
		for _, err := range err.(validator.ValidationErrors) {
			merr = multierr.Append(merr, eris.New(err.Translate(*v.trans)))
		}
		return merr
	}
	return nil
}

func (v *Validator) Engine() any {
	return v.validate
}

func main() {
	w := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.DateTime
	})
	log := zerolog.New(w).With().Timestamp().Logger()

	ginhelper.
		NewZerologWriter(log, zerolog.InfoLevel).
		SetAll()

	ginhelper.H.WithBindingErrorHandler(func(c *gin.Context, err error) {
		c.JSON(http.StatusOK, NewResponse(ResponseCodeBadRequest, err.Error()))
	}).WithErrorHandler(func(c *gin.Context, err error) {
		c.Error(err)
	}).WithSuccessHandler(func(c *gin.Context, resp any) {
		c.JSON(http.StatusOK, NewResponse(ResponseCodeSuccess, resp))
	}).WithBindingValidator(NewValidator())

	app := gin.New()
	app.Use(gin.Recovery())

	app.Use(func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			if r := new(Response); eris.As(err, r) || eris.As(err, &r) {
				c.JSON(http.StatusOK, r)
			} else {
				c.JSON(200, NewResponse(ResponseCodeInternalError, err.JSON()))
			}
		}
	})

	v1g := app.Group("v1")

	ginhelper.H.Router(v1g).
		GET("/echo/:id", EchoController).
		POST("/create", CreateController)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
