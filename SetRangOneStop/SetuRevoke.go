/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangPermissionRev
 *@date    2024/11/29 10:05
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"errors"
	"strconv"
	"strings"
)

// PermissionRevoke
// 回收逻辑
func PermissionRevoke(a *util.SetPolicy) error {
	result, err := getPolicyByUser(a)
	if err != nil {
		util.Loggrs.Warn(err.Error())
	}
	if len(result) < 1 {
		return nil
	}
	util.Loggrs.Info("POLICY ->:EXIST")

	var errlist []string
	for _, item := range result {
		policy := updatePolicyByRevoke(item, a)
		if len(policy) == 0 {
			continue
		}
		util.Loggrs.Info("REVOKE POLICY ->:", policy)

		open := strings.Split(item, ">")
		id, _ := strconv.Atoi(open[0])
		//更新模式
		post, err := execPolicy(id, "update", policy)
		if err != nil {
			util.Loggrs.Error(err.Error())
			errlist = append(errlist, err.Error())
			continue
		}
		util.Loggrs.Info(post)
	}

	if errlist != nil {
		return errors.New(strings.Join(errlist, "\n"))
	}
	return nil
}
