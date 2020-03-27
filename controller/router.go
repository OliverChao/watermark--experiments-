package controller

import (
	"WaterMasking/controller/handlerFuncs"
	"WaterMasking/service"
	"github.com/b3log/gulu"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func MapRoutes() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		Secure:   strings.HasPrefix("http://127.0.0.1", "https"),
		HttpOnly: true,
	})

	engine.Use(sessions.Sessions("database_watermarking", store))
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/pan", "./pan")

	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})
	engine.POST("/", handlerFuncs.Login)

	engine.GET("/test", func(c *gin.Context) {
		result := gulu.Ret.NewResult()
		defer c.JSON(http.StatusOK, result)
		data := map[string]interface{}{}
		data["some_data"] = service.SourceCache.GetSourceData()[:10]
		result.Data = data
	})

	exp := engine.Group("/exp")
	exp.Use(LoginCheck)
	exp.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "test token ... ")
	})

	exp.GET("/test", func(c *gin.Context) {
		c.HTML(200, "showdataHint.html", gin.H{
			"total": 100000,
		})
		//c.String(200, "<h1>oliver<h1>")
		//c.HTML(200,"")
	})

	exp.POST("/test", func(c *gin.Context) {
		c.String(200, "<h1>oliver<h1>")
	})

	exp.GET("/data", handlerFuncs.ShowData)
	exp.POST("/data", handlerFuncs.ShowDataPost)

	exp.GET("/insert", func(c *gin.Context) {
		var table []string
		switch service.MarkedCache.IsCached() {
		case true:
			table = []string{"students", "student_backs"}
		case false:
			table = []string{"students"}
		}
		c.HTML(200, "insert.html", gin.H{
			"table": table,
		})
	})
	exp.POST("/insert", handlerFuncs.InsertWaterMarking)

	exp.GET("/backup", handlerFuncs.BackupData)

	exp.GET("/resumption", func(c *gin.Context) {
		service.ExchangeTableName()
		c.String(200, "resume successfully...")
	})

	exp.GET("/verify", handlerFuncs.VerifyWaterMarking)

	alg := engine.Group("/algorithms")
	alg.Use(LoginCheck)

	return engine
}
