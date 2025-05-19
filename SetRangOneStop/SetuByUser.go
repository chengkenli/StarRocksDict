/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangeUpdatePolicyByUser
 *@date    2024/11/28 21:25
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"encoding/json"
	"strings"
)

// updatePolicy
// 更新policy内容（根据用户更新它得权限）
func updatePolicyByUser(a *util.SetPolicy, item []byte) (int, []byte, error) {
	var b policyStruct
	err := json.Unmarshal(item, &b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return -1, nil, err
	}
	//根据已有的授权策略变更新成新的policy
	//策略
	if !a.Revoke {
		util.Loggrs.Info("AUTHORIZE")
		//权限
		for _, permission := range strings.Split(a.Permission, ",") {
			b.PolicyItems[0].Accesses = append(b.PolicyItems[0].Accesses, accesses{
				Type:      permission,
				IsAllowed: true,
			})
		}
	} else {
		util.Loggrs.Info("REVOKE")
		// 创建一个新的切片来保存结果
		filtered := make(
			[]struct {
				Type      string `json:"type"`
				IsAllowed bool   `json:"isAllowed"`
			}, 0)
		for _, access := range b.PolicyItems[0].Accesses {
			// 如果Type字段不等于typ，则保留这个元素
			if access.Type != a.Permission {
				filtered = append(filtered, access)
			}
		}
		b.PolicyItems[0].Accesses = filtered
	}

	marshal, err := json.Marshal(&b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return b.ID, marshal, err
	}
	return b.ID, marshal, nil
}
