/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangGetPolicyByName
 *@date    2024/11/28 21:09
 */

package SetRangeOneStop

import (
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"encoding/json"
	"fmt"
	"strings"
)

// getPolicyByUser
// 获取指定服务下的策略信息列表
func getPolicyByUser(a *util.SetPolicy) ([]string, error) {
	r, err := Post(
		&util.RangerHost{
			Server: util.Read.Ranger.Host,
			Access: util.Read.Ranger.User,
			Secret: util.Read.Ranger.Password,
		}, "GET", fmt.Sprintf("/service/public/v2/api/service/%s/policy", a.App), nil)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil, err
	}
	var b policyStructsAll
	err = json.Unmarshal(r, &b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil, err
	}

	var policys []string
	var i int
	for _, c := range b {
		//从PolicyItems - users中获取目前账号是否存在列表中
	loop:
		for _, item := range c.PolicyItems {
			//账号
			for _, user := range strings.Split(a.User, ",") {
				if tools.StrInSlice(user, item.Users) {
					break
				} else {
					continue loop
				}
			}
			//库表
			for _, table := range strings.Split(a.Table, ",") {
				split := strings.Split(table, ".")
				if len(split) < 2 {
					continue loop
				}
				if split[1] != "*" {
					if !tools.StrInSlice(split[1], c.Resources.Table.Values) {
						continue loop
					}
				} else {
					break
				}
			}
			//权限
			var accessList []string
			for _, access := range item.Accesses {
				accessList = append(accessList, strings.ToLower(access.Type))
			}
			for _, permit := range strings.Split(a.Permission, ",") {
				if !tools.StrInSlice(strings.ToLower(permit), accessList) {
					continue loop
				}
			}

			marshal, err := json.Marshal(&c)
			if err != nil {
				util.Loggrs.Error(err.Error())
				continue
			}
			//剔除无效的信息
			policy := strings.NewReplacer("},}}", "}}}").Replace(ReplaceResource(c, string(marshal)))

			util.Loggrs.Info("GET POLICY ->:", i, " ", policy)
			policys = append(policys, fmt.Sprintf("%d>%s", c.ID, policy))
			i++
		}
	}
	return policys, nil
}
