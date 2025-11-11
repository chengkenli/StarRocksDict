/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangAddChangePolicy
 *@date    2024/11/28 21:01
 */

package SetRangeOneStop

import (
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"fmt"
	"strings"
	"time"
)

// changePolicy
// 用于判断权限是system？还是databases？table？view？
func changePolicy(a *util.SetPolicy) []string {
	System := []string{"GRANT", "NODE", "OPERATE", "PLUGIN", "FILE", "BLACKLIST", "REPOSITORY", "CREATE GLOBAL FUNCTION", "CREATE RESOURCE", "CREATE RESOURCE GROUP", "CREATE EXTERNAL CATALOG", "CREATE STORAGE VOLUME", "CREATE WAREHOUSE", "SECURITY", "CREATE FAILOVER GROUP"}
	//判断权限类型,管理员||普通
	var auth_system, auth_ordinary, bodys []string
	for _, item := range strings.Split(a.Permission, ",") {
		if tools.StrInSlice(strings.ToUpper(item), System) {
			auth_system = append(auth_system, item)
		} else {
			auth_ordinary = append(auth_ordinary, item)
		}
	}

	//label := fmt.Sprintf("api - %s_%d", time.Now().Format("200601021504"), time.Now().UnixNano())
	label := fmt.Sprintf("api-%s-%s", time.Now().Format("20060102"), tools.GenRandomStr(7))
	//管理员方式
	if len(auth_system) > 0 {
		body := `"system": {"values": ["*"],"isExcludes": false,"isRecursive": false}`
		policy := strings.NewReplacer(
			"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
			"STARROCKS.SERVICE.NAME", a.App,
			"STARROCKS.PARMS.POLICYNAME", label,
			"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
			"STARROCKS.PARMS.USER", policyJson(strings.Split(a.User, ",")),
			"STARROCKS.PARMS.PERMISSION", formatColData(a.Permission),
			"STARROCKS.PARMS.RESOURCES", "{"+body+"}",
		).Replace(PerJson)
		//policy
		bodys = append(bodys, policy)
	}

	//普通用户方式
	if len(auth_ordinary) > 0 {
		sign := arrJsonDt(a.App, a.Table)

		var body []string
		//catalog
		if len(a.Catalog) != 0 {
			body = append(body, fmt.Sprintf(`"catalog": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(strings.Split(a.Catalog, ","))))
		}
		//database
		if len(sign.Database) != 0 {
			strings.NewReplacer(" ", "").Replace(strings.ToLower(a.Permission))
			if len(sign.Table) == 0 && a.Permission == "select" {
				body = append(body, fmt.Sprintf(`"database": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Database)))
				body = append(body, fmt.Sprintf(`"table": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
				body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
			} else {
				body = append(body, fmt.Sprintf(`"database": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Database)))
			}
		}
		//table
		if len(sign.Table) != 0 {
			body = append(body, fmt.Sprintf(`"table": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Table)))
		}
		//materialized
		if len(sign.Materialized) != 0 {
			body = append(body, fmt.Sprintf(`"materialized_view": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Materialized)))
		}
		//view
		if len(sign.View) != 0 {
			body = append(body, fmt.Sprintf(`"view": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.View)))
		}
		policy := strings.NewReplacer(
			"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
			"STARROCKS.SERVICE.NAME", a.App,
			"STARROCKS.PARMS.POLICYNAME", label,
			"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
			"STARROCKS.PARMS.USER", policyJson(strings.Split(a.User, ",")),
			"STARROCKS.PARMS.PERMISSION", formatColData(a.Permission),
			"STARROCKS.PARMS.RESOURCES", "{"+strings.Join(body, ",")+"}",
		).Replace(PerJson)
		bodys = append(bodys, policy)
	}
	return bodys
}
