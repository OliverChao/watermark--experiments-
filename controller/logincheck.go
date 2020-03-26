package controller

import (
	"WaterMasking/util/shadow"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoginCheck(c *gin.Context) {
	//baseConfig := baseCon.LoadBaseConfig()
	ipAddress := "http://127.0.0.1:8080"

	session := sessions.Default(c)
	token, ok := session.Get("token").(string)
	if !ok {

		c.Redirect(http.StatusSeeOther, ipAddress)
		c.Abort()
		return
	} else {
		logrus.Info("[LoginCheck] get token...")
	}
	data, sign := shadow.UnParseToken(token)
	e := shadow.RsaVerify(data, sign)
	if e != nil {
		c.Redirect(http.StatusSeeOther, ipAddress)
		c.Abort()
		return
	} else {
		logrus.Debug("login check correct")
	}
	c.Next()
}
