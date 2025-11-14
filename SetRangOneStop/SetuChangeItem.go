/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangChangeItem
 *@date    2024/12/2 18:30
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"encoding/json"
	"strings"
)

func changeItem(b policyStructSingle, i int, permission, user string, Access []struct {
	Type      string `json:"type"`
	IsAllowed bool   `json:"isAllowed"`
}) policyItem {

	//遍历
	for _, account := range b.PolicyItems[i].Users {
		//定义结构体
		var policyItem struct {
			Accesses []struct {
				Type      string `json:"type"`
				IsAllowed bool   `json:"isAllowed"`
			} `json:"accesses"`
			Users         []string      `json:"users"`
			Groups        []interface{} `json:"groups"`
			Roles         []interface{} `json:"roles"`
			Conditions    []interface{} `json:"conditions"`
			DelegateAdmin bool          `json:"delegateAdmin"`
		}

		policyItem.Users = append(policyItem.Users, account)
		//定义权限
		for _, permit := range strings.Split(strings.ToLower(permission), ",") {
			for y := 0; y < len(b.PolicyItems[i].Accesses); y++ {
				if strings.ToLower(b.PolicyItems[i].Accesses[y].Type) == permit {
					Access = append(Access, b.PolicyItems[i].Accesses[y])
				}
			}
		}
		if user == account {
			policyItem.Accesses = accessSlice(Access, b.PolicyItems[i].Accesses)
		} else {
			policyItem.Accesses = b.PolicyItems[i].Accesses
		}
		//policy item载入重构
		marshal, err := json.Marshal(&policyItem)
		if err != nil {
			util.Loggrs.Warn(err.Error())
		}
		util.Loggrs.Info("ADDED  POLICY ->:", string(marshal))

		policyItems = append(policyItems, policyItem)
	}
	return policyItems
}
