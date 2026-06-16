/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangUpdatePolicyByAccess
 *@date    2024/11/28 21:24
 */

package SetRangeOneStop

import (
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"encoding/json"
	"strings"
)

// updatePolicyByAccess
// 更新policy内容（根据policy更新它的用户） - system级别
func updatePolicyBySystem(a *util.SetPolicy, item []byte) (int, []byte, error) {
	var b policyStructBySystem
	err := json.Unmarshal(item, &b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return -1, nil, err
	}
	//根据已有的授权策略变更新成新的policy
	//策略
	if !a.Revoke {
		util.Loggrs.Info("AUTHORIZE")
		var acces []struct {
			Type      string `json:"type"`
			IsAllowed bool   `json:"isAllowed"`
		}
		for _, item := range strings.Split(a.Permission, ",") {
			acces = append(acces, accesses{
				Type:      item,
				IsAllowed: true,
			})
		}

		type ssm struct {
			Accesses []struct {
				Type      string `json:"type"`
				IsAllowed bool   `json:"isAllowed"`
			} `json:"accesses"`
			Users         []string      `json:"users"`
			Groups        []string      `json:"groups"`
			Roles         []interface{} `json:"roles"`
			Conditions    []interface{} `json:"conditions"`
			DelegateAdmin bool          `json:"delegateAdmin"`
		}
		b.PolicyItems = append(b.PolicyItems, ssm{
			Accesses:      acces,
			Users:         strings.Split(a.User, ","),
			Groups:        nil,
			Roles:         nil,
			Conditions:    nil,
			DelegateAdmin: true,
		})

	} else {
		util.Loggrs.Info("REVOKE")
		// 创建一个新的切片来保存结果
		filtered := make(
			[]struct {
				Accesses []struct {
					Type      string `json:"type"`
					IsAllowed bool   `json:"isAllowed"`
				} `json:"accesses"`
				Users         []string      `json:"users"`
				Groups        []string      `json:"groups"`
				Roles         []interface{} `json:"roles"`
				Conditions    []interface{} `json:"conditions"`
				DelegateAdmin bool          `json:"delegateAdmin"`
			}, 0)
		for _, item := range b.PolicyItems {
			for _, user_sign := range strings.Split(a.User, ",") {
				// 如果Type字段不等于typ，则保留这个元素
				if !tools.StrInSlice(user_sign, item.Users) {
					filtered = append(filtered, item)
				}
			}
		}
		b.PolicyItems = filtered
	}

	marshal, err := json.Marshal(&b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return b.ID, marshal, err
	}
	return b.ID, marshal, nil
}
