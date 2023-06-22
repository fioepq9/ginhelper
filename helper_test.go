package ginhelper

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	code := m.Run()
	os.Exit(code)
}

type TestGETRequest struct {
	Name  string `uri:"name" binding:"required"`
	Age   int    `form:"age" binding:"required"`
	Token string `header:"token" binding:"required"`
}

func TestGET(t *testing.T) {
	r := gin.New()
	H.GET(r, "/test/:name", func(c *gin.Context, req *TestGETRequest) error {
		assert.Equal(t, "foo", req.Name)
		assert.Equal(t, 42, req.Age)
		assert.Equal(t, "1234", req.Token)
		return nil
	})
	go r.Run(":18080")

	req.R().SetHeader("token", "1234").Get("http://localhost:18080/test/foo?age=42")
}

type TestPOSTRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func TestPOST(t *testing.T) {
	r := gin.New()
	H.POST(r, "/test", func(c *gin.Context, req *TestPOSTRequest) error {
		assert.Equal(t, "foo", req.Username)
		assert.Equal(t, "bar", req.Password)
		return nil
	})
	go r.Run(":18081")

	req.R().SetHeader("Content-Type", "application/json").SetBody(`{"username":"foo","password":"bar"}`).Post("http://localhost:18081/test")
}
