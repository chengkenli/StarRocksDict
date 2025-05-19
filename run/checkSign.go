package run

import (
	"StarRocksDict/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func checkSign(c *gin.Context) {
	sign := c.GetHeader("X-StarRocks")
	if len(sign) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "密钥HEAD为空，您无权使用管理员POST功能！！")
		return
	}
	token, err := util.AesDecrypt2(sign, util.Config.GetString("StarRocks.service.enckey"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "密钥HEAD为空，您无权使用管理员POST功能！！")
		return
	}
	env := util.Config.GetString("StarRocks.service.X-StarRocks")
	if token != env {
		c.AbortWithStatusJSON(http.StatusBadRequest, "密钥HEAD为空，您无权使用管理员POST功能！！")
		return
	}
	c.Next()
}
