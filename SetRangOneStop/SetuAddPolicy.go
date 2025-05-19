package SetRangeOneStop

import (
	"StarRocksDict/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"strings"
)

// ApiaddPolicy 创建策略
func ApiaddPolicy(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		util.Loggrs.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	if len(data) == 0 {
		util.Loggrs.Error("json data is nil")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": "json data is nil."})
		return
	}
	var a util.SetPolicy
	err = json.Unmarshal(data, &a)
	if err != nil {
		util.Loggrs.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	r, err := addPolicy(&a)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	if strings.Contains(r.(string), "Fail") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": r.(string)})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok", "message": r.(string)})
}

// addPolicy 创建策略
func addPolicy(a *util.SetPolicy) (interface{}, error) {
	var err error
	var dict string

	switch a.Revoke {
	case true:
		//回收权限
		err = PermissionRevoke(a)
		dict = "REVOKE"
	case false:
		//授予权限
		err = PermissionGrant(a)
		dict = "GRANT"
	}

	if err != nil {
		return fmt.Sprintf("Fail %s %s: CATALOG=%s,USER=%s,TABLE=%s,PERMISSION=%s %s", dict, a.App, a.Catalog, a.User, a.Table, a.Permission, err.Error()), err
	}
	return fmt.Sprintf("Success %s %s: CATALOG=%s,USER=%s,TABLE=%s,PERMISSION=%s %s", dict, a.App, a.Catalog, a.User, a.Table, a.Permission, ""), nil
}
