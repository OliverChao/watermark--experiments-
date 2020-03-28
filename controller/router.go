package controller

import (
	"WaterMasking/controller/handlerFuncs"
	"WaterMasking/service"
	"github.com/b3log/gulu"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	engine.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Options(sessions.Options{
			Path:   "/",
			MaxAge: -1,
		})
		session.Clear()
		if err := session.Save(); nil != err {
			logrus.Errorf("saves session failed: " + err.Error())
		}

		c.Redirect(http.StatusSeeOther, "http://127.0.0.1:8080/")
		c.Abort()
	})

	engine.GET("/test", func(c *gin.Context) {
		result := gulu.Ret.NewResult()
		defer c.JSON(http.StatusOK, result)
		data := map[string]interface{}{}
		data["some_data"] = service.SourceCache.GetSourceData()[:10]
		result.Data = data
	})

	// exp remains functional model
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
		if !service.ResumeController.IsResumed() {
			service.ExchangeTableName()
			service.ResumeController.EnResumed()
		}
		c.String(200, "resume successfully...")
	})

	exp.GET("/verify", func(c *gin.Context) {
		var tables []string
		switch service.MarkedCache.IsCached() {
		case true:
			tables = []string{"students", "student_backs"}
		case false:
			tables = []string{"students"}
		}
		c.HTML(http.StatusOK, "verify.html", gin.H{
			"table": tables,
		})
	})

	exp.POST("/verify", handlerFuncs.VerifyWaterMarking)

	//att is attacking model
	att := engine.Group("/attack")
	att.Use(LoginCheck)

	att.GET("/delete", handlerFuncs.DeleteAttack)

	att.GET("/addmore", handlerFuncs.AddMoreAttack)

	att.GET("/reverse", func(c *gin.Context) {
		c.HTML(200, "reverse_attack.html", gin.H{})
	})
	att.POST("/reverse", handlerFuncs.ReverseAttack)
	return engine
}
