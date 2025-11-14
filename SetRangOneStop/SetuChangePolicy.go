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

	var dblist, tblist, total []string
	for _, item := range strings.Split(strings.NewReplacer(" ", "").Replace(a.Table), ",") {
		tag := strings.Split(item, ".")
		if len(strings.Split(item, ".")) == 2 {
			tbname := tag[1]
			if tbname == "*" {
				dblist = append(dblist, item)
			} else {
				tblist = append(tblist, item)
			}
		} else {
			dblist = append(dblist, item)
		}
	}

	if dblist != nil {
		util.Loggrs.Info("database schema.")
		util.Loggrs.Debug(strings.Join(dblist, ","))
		total = append(total, splitJson2db(a, strings.Join(dblist, ","))...)
	}
	if tblist != nil {
		util.Loggrs.Info("table|materialized_view|view schema.")
		util.Loggrs.Debug(strings.Join(tblist, ","))
		total = append(total, splitjson2tb(a, strings.Join(tblist, ","))...)
	}

	//System := []string{"GRANT", "NODE", "OPERATE", "PLUGIN", "FILE", "BLACKLIST", "REPOSITORY", "CREATE GLOBAL FUNCTION", "CREATE RESOURCE", "CREATE RESOURCE GROUP", "CREATE EXTERNAL CATALOG", "CREATE STORAGE VOLUME", "CREATE WAREHOUSE", "SECURITY", "CREATE FAILOVER GROUP"}
	////判断权限类型,管理员||普通
	//var auth_system, auth_ordinary, bodys []string
	//for _, item := range strings.Split(a.Permission, ",") {
	//	if tools.StrInSlice(strings.ToUpper(item), System) {
	//		auth_system = append(auth_system, item)
	//	} else {
	//		auth_ordinary = append(auth_ordinary, item)
	//	}
	//}
	//
	//var label string
	//if a.Policyname == "" {
	//	label = fmt.Sprintf("api-%s-%s", time.Now().Format("20060102"), tools.GenRandomStr(7))
	//} else {
	//	label = fmt.Sprintf("api-%s-%s", time.Now().Format("20060102"), a.Policyname)
	//}
	////管理员方式
	//if len(auth_system) > 0 {
	//	body := `"system": {"values": ["*"],"isExcludes": false,"isRecursive": false}`
	//	policy := strings.NewReplacer(
	//		"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
	//		"STARROCKS.SERVICE.NAME", a.App,
	//		"STARROCKS.PARMS.POLICYNAME", label,
	//		"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
	//		"STARROCKS.PARMS.USER", policyJson(strings.Split(a.User, ",")),
	//		"STARROCKS.PARMS.PERMISSION", formatColData(a.Permission),
	//		"STARROCKS.PARMS.RESOURCES", "{"+body+"}",
	//	).Replace(PerJson)
	//	//policy
	//	bodys = append(bodys, policy)
	//}
	//
	////普通用户方式
	//if len(auth_ordinary) > 0 {
	//	sign := arrJsonDt(a.App, a.Table)
	//
	//	var body []string
	//	//catalog
	//	if len(a.Catalog) != 0 {
	//		body = append(body, fmt.Sprintf(`"catalog": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(strings.Split(a.Catalog, ","))))
	//	}
	//	//database
	//	if len(sign.Database) != 0 && len(sign.Table) == 0 && len(sign.Materialized) == 0 && len(sign.View) == 0 {
	//		//if len(sign.Table) == 0 && a.Permission == "select" {
	//		//	body = append(body, fmt.Sprintf(`"database": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Database)))
	//		//	body = append(body, fmt.Sprintf(`"table": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//		//	body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//		//} else {
	//		body = append(body, fmt.Sprintf(`"database": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Database)))
	//		body = append(body, fmt.Sprintf(`"table": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//		body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//		//}
	//	}
	//	//table
	//	if len(sign.Table) != 0 {
	//		body = append(body, fmt.Sprintf(`"table": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Table)))
	//		body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//	}
	//	//materialized
	//	if len(sign.Materialized) != 0 {
	//		body = append(body, fmt.Sprintf(`"materialized_view": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Materialized)))
	//		body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//	}
	//	//view
	//	if len(sign.View) != 0 {
	//		body = append(body, fmt.Sprintf(`"view": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.View)))
	//		body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
	//	}
	//	policy := strings.NewReplacer(
	//		"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
	//		"STARROCKS.SERVICE.NAME", a.App,
	//		"STARROCKS.PARMS.POLICYNAME", label,
	//		"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
	//		"STARROCKS.PARMS.USER", policyJson(strings.Split(a.User, ",")),
	//		"STARROCKS.PARMS.PERMISSION", formatColData(a.Permission),
	//		"STARROCKS.PARMS.RESOURCES", "{"+strings.Join(body, ",")+"}",
	//	).Replace(PerJson)
	//	bodys = append(bodys, policy)
	//}
	//return bodys
	return total
}

// 表级别的拆分
// @setpolicy
// @table list
// @result slice
func splitjson2tb(e *util.SetPolicy, tbmap string) []string {

	System := []string{"GRANT", "NODE", "OPERATE", "PLUGIN", "FILE", "BLACKLIST", "REPOSITORY", "CREATE GLOBAL FUNCTION", "CREATE RESOURCE", "CREATE RESOURCE GROUP", "CREATE EXTERNAL CATALOG", "CREATE STORAGE VOLUME", "CREATE WAREHOUSE", "SECURITY", "CREATE FAILOVER GROUP"}
	//判断权限类型,管理员||普通
	var auth_system, auth_ordinary, bodys []string
	for _, item := range strings.Split(e.Permission, ",") {
		if tools.StrInSlice(strings.ToUpper(item), System) {
			auth_system = append(auth_system, item)
		} else {
			auth_ordinary = append(auth_ordinary, item)
		}
	}

	var label string
	if e.Policyname == "" {
		label = fmt.Sprintf("api-%s-%s", time.Now().Format("20060102"), tools.GenRandomStr(7))
	} else {
		label = fmt.Sprintf("api-%s-%s", time.Now().Format("20060102"), e.Policyname)
	}
	//管理员方式
	if len(auth_system) > 0 {
		body := `"system": {"values": ["*"],"isExcludes": false,"isRecursive": false}`
		policy := strings.NewReplacer(
			"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
			"STARROCKS.SERVICE.NAME", e.App,
			"STARROCKS.PARMS.POLICYNAME", label,
			"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
			"STARROCKS.PARMS.USER", policyJson(strings.Split(e.User, ",")),
			"STARROCKS.PARMS.PERMISSION", formatColData(e.Permission),
			"STARROCKS.PARMS.RESOURCES", "{"+body+"}",
		).Replace(PerJson)
		//policy
		bodys = append(bodys, policy)
	}

	//普通用户方式
	if len(auth_ordinary) > 0 {
		sign := arrJsonDt(e.App, tbmap)

		var body []string
		//catalog
		if len(e.Catalog) != 0 {
			body = append(body, fmt.Sprintf(`"catalog": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(strings.Split(e.Catalog, ","))))
		}
		//database
		if len(sign.Database) != 0 {
			body = append(body, fmt.Sprintf(`"database": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Database)))
		}
		//table
		if len(sign.Table) != 0 {
			body = append(body, fmt.Sprintf(`"table": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Table)))
			body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
		}
		//materialized
		if len(sign.Materialized) != 0 {
			body = append(body, fmt.Sprintf(`"materialized_view": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Materialized)))
			body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
		}
		//view
		if len(sign.View) != 0 {
			body = append(body, fmt.Sprintf(`"view": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.View)))
			body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
		}
		policy := strings.NewReplacer(
			"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
			"STARROCKS.SERVICE.NAME", e.App,
			"STARROCKS.PARMS.POLICYNAME", label,
			"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
			"STARROCKS.PARMS.USER", policyJson(strings.Split(e.User, ",")),
			"STARROCKS.PARMS.PERMISSION", formatColData(e.Permission),
			"STARROCKS.PARMS.RESOURCES", "{"+strings.Join(body, ",")+"}",
		).Replace(PerJson)
		bodys = append(bodys, policy)
	}
	return bodys

}

// 库级别的拆分，库级别下面分为表、物化视图、视图
// @parm setpolicy
// @parm database list
// @result slice
func splitJson2db(e *util.SetPolicy, dbmap string) []string {

	var rand string
	if e.Policyname == "" {
		rand = tools.GenRandomStr(7)
	} else {
		rand = e.Policyname
	}

	var bodys []string
	for _, tag := range []string{"table", "materialized_view", "view"} {
		label := fmt.Sprintf("api-%s-%s-%s", time.Now().Format("20060102"), rand, tag)

		System := []string{"GRANT", "NODE", "OPERATE", "PLUGIN", "FILE", "BLACKLIST", "REPOSITORY", "CREATE GLOBAL FUNCTION", "CREATE RESOURCE", "CREATE RESOURCE GROUP", "CREATE EXTERNAL CATALOG", "CREATE STORAGE VOLUME", "CREATE WAREHOUSE", "SECURITY", "CREATE FAILOVER GROUP"}
		//判断权限类型,管理员||普通
		var auth_system, auth_ordinary []string
		for _, item := range strings.Split(e.Permission, ",") {
			if tools.StrInSlice(strings.ToUpper(item), System) {
				auth_system = append(auth_system, item)
			} else {
				auth_ordinary = append(auth_ordinary, item)
			}
		}

		//管理员方式
		if len(auth_system) > 0 {
			body := `"system": {"values": ["*"],"isExcludes": false,"isRecursive": false}`
			policy := strings.NewReplacer(
				"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
				"STARROCKS.SERVICE.NAME", e.App,
				"STARROCKS.PARMS.POLICYNAME", label,
				"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
				"STARROCKS.PARMS.USER", policyJson(strings.Split(e.User, ",")),
				"STARROCKS.PARMS.PERMISSION", formatColData(e.Permission),
				"STARROCKS.PARMS.RESOURCES", "{"+body+"}",
			).Replace(PerJson)
			//policy
			bodys = append(bodys, policy)
		}
		//普通用户方式
		if len(auth_ordinary) > 0 {
			sign := arrJsonDt(e.App, dbmap)

			var body []string
			//catalog
			if len(e.Catalog) != 0 {
				body = append(body, fmt.Sprintf(`"catalog": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(strings.Split(e.Catalog, ","))))
			}
			//database
			if len(sign.Database) != 0 {
				body = append(body, fmt.Sprintf(`"database": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson(sign.Database)))
				body = append(body, fmt.Sprintf(`"%s": {"values": %s,"isExcludes": false,"isRecursive": false}`, tag, policyJson([]string{"*"})))
				if tag == "table" {
					body = append(body, fmt.Sprintf(`"column": {"values": %s,"isExcludes": false,"isRecursive": false}`, policyJson([]string{"*"})))
				}
			}
			policy := strings.NewReplacer(
				"STARROCKS.PARMS.DESCRIPTION", "api-"+time.Now().Format("2006.01.02 15:04:05"),
				"STARROCKS.SERVICE.NAME", e.App,
				"STARROCKS.PARMS.POLICYNAME", label,
				"STARROCKS.PARMS.POLICYLABELS", "starrocks@"+time.Now().Format("20060102"),
				"STARROCKS.PARMS.USER", policyJson(strings.Split(e.User, ",")),
				"STARROCKS.PARMS.PERMISSION", formatColData(e.Permission),
				"STARROCKS.PARMS.RESOURCES", "{"+strings.Join(body, ",")+"}",
			).Replace(PerJson)
			bodys = append(bodys, policy)
		}
	}

	return bodys
}
