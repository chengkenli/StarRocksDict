package SetRangeOneStop

import (
	"StarRocksDict/util"
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
	r, err := AddPolicy(&a)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok", "message": r})
}

// AddPolicy 创建策略
func AddPolicy(a *util.SetPolicy) ([]util.SetResult, error) {
	var err error
	var dict string
	var result []util.SetResult

	switch a.Revoke {
	case true:
		//回收权限
		dict = "revoke"
		err = PermissionRevoke(a)
	case false:
		//授予权限
		dict = "grant"
		err = PermissionGrant(a)
	}

	if err != nil {
		for _, table := range strings.Split(a.Table, ",") {
			table = strings.NewReplacer(" ", "").Replace(table)
			result = append(result, util.SetResult{
				App:     a.App,
				Catalog: a.Catalog,
				User:    a.User,
				Table:   table,
				Permit:  a.Permission,
				State:   "Fail",
				Demand:  dict,
				Comment: err.Error(),
				Mark:    "ranger",
			})
		}
		return result, nil
	}

	for _, table := range strings.Split(a.Table, ",") {
		table = strings.NewReplacer(" ", "").Replace(table)
		result = append(result, util.SetResult{
			App:     a.App,
			Catalog: a.Catalog,
			User:    a.User,
			Table:   table,
			Permit:  a.Permission,
			State:   "Ok",
			Demand:  dict,
			Comment: "",
			Mark:    "ranger",
		})
	}
	return result, nil
}
