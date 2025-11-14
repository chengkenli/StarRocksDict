/*
 *@author  chengkenli
 *@project StarRocksSupport
 *@package etrics
 *@file    EtricsResourceGroupConstraint
 *@date    2024/10/22 9:46
 */

package SetDictOneStop

import (
	"StarRocksDict/util"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

var vclass []string

func ResourceGroup(db *gorm.DB, username string, ldap bool) {
	if vclass == nil {
		util.Loggrs.Warn("资源组为空，初始化！")
		vclass = constraint(db)
	}
	if strings.Contains(strings.Join(vclass, ","), username) {
		return
	}
	util.Loggrs.Info(fmt.Sprintf("开始 - 约束非白名单用户%s的并发度！", username))
	tf := time.Now().Format("060102150405")

	concurrency_limit := 3
	if ldap {
		concurrency_limit = 10
	}
	sql := fmt.Sprintf(`
		CREATE RESOURCE GROUP in%s
		TO(
			user='%s'
		)
		WITH(
			"cpu_core_limit"="10",
			"mem_limit"="50%%",
			"concurrency_limit"="%d"
		)`, tf, username, concurrency_limit)
	r := db.Exec(sql)
	if r.Error != nil {
		util.Loggrs.Info(sql)
		util.Loggrs.Warn(r.Error.Error())
		return
	}
	util.Loggrs.Info(fmt.Sprintf("结束 - 约束非白名单用户%s的并发度！[%s] - [%s]", username, tf, username))
}

func constraint(db *gorm.DB) []string {
	var m []map[string]interface{}
	r := db.Raw("SHOW RESOURCE GROUPS ALL").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Warn(r.Error.Error())
		return nil
	}
	var vc []string
	for _, m2 := range m {
		if v, ok := m2["classifiers"]; ok {
			vc = append(vc, v.(string))
		}
	}
	return vc
}
