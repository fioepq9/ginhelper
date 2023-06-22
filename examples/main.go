package main

import (
	"fmt"
	"net/http"

	"github.com/fioepq9/ginhelper"
	"github.com/gin-gonic/gin"
)

// http localhost:8080/echo/1234 message==hello token:1234
type EchoRequest struct {
	ID      int    `uri:"id" binding:"required"`
	Message string `form:"message" binding:"required"`
	Token   string `header:"token" binding:"required"`
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

func main() {
	app := gin.Default()
	gu := ginhelper.New(
		ginhelper.WithBindingErrorHandler(func(c *gin.Context, err error) {
			c.JSON(http.StatusOK, NewResponse(ResponseCodeBadRequest, err))
		}),
		ginhelper.WithErrorHandler(func(c *gin.Context, err error) {
			c.Error(err)
		}),
		ginhelper.WithSuccessHandler(func(c *gin.Context, resp any) {
			c.JSON(http.StatusOK, NewResponse(ResponseCodeSuccess, resp))
		}),
	)

	app.Use(func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			c.JSON(200, NewResponse(ResponseCodeInternalError, err))
		}
	})

	gu.GET(app, "/echo/:id", func(c *gin.Context, req *EchoRequest) error {
		fmt.Printf("%+v\n", req)
		return fmt.Errorf("what is the problem? 42 is the answer")
	})

	gu.POST(app, "/create", func(c *gin.Context, req *CreateRequest) (*CreateResponse, error) {
		fmt.Printf("%+v\n", req)
		resp := &CreateResponse{
			ID: "1234",
		}
		return resp, nil
	})
	if err := app.Run(); err != nil {
		panic(err)
	}
}
