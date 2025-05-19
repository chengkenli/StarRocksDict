/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangUpdatePolicyByRevoke
 *@date    2024/11/29 14:01
 */

package SetRangeOneStop

import (
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"encoding/json"
	"fmt"
	"strings"
)

func updatePolicyByRevoke(result string, a *util.SetPolicy) string {
	open := strings.Split(result, ">")

	var b policyStructSingle
	err := json.Unmarshal([]byte(open[1]), &b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return ""
	}

	var policyItems []struct {
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

loop:
	for i := 0; i < len(b.PolicyItems); i++ {
		//权限收集
		var accessList, userList, dbList, tableList []string
		for _, access := range b.PolicyItems[i].Accesses {
			accessList = append(accessList, access.Type)
		}
		//用户收集
		userList = b.PolicyItems[i].Users
		userListAgent := SplitSlice(a.User)
		//库表收集
		for _, value := range b.Resources.Database.Values {
			if strings.Contains(value, util.KeySingle) {
				continue
			}
			dbList = append(dbList, value)
		}
		for _, value := range b.Resources.Table.Values {
			if strings.Contains(value, util.KeySingle) {
				continue
			}
			tableList = append(tableList, value)
		}

		var dbListAgent, tableListAgent []string
		for _, schema := range SplitSlice(a.Table) {
			data := strings.Split(schema, ".")
			dbListAgent = append(dbListAgent, data[0])
			tableListAgent = append(tableListAgent, data[1])
		}

		userList = tools.RmSliceStr(userList)
		userListAgent = tools.RmSliceStr(userListAgent)

		dbList = tools.RmSliceStr(dbList)
		dbListAgent = tools.RmSliceStr(dbListAgent)

		tableList = tools.RmSliceStr(tableList)
		tableListAgent = tools.RmSliceStr(tableListAgent)

		//初步判断
		for _, user := range userListAgent {
			if !tools.StrInSlice(user, userList) {
				util.Loggrs.Info(b.ID, " USER CONTINUE")
				continue loop
			}
		}
		for _, db := range dbListAgent {
			if !tools.StrInSlice(db, dbList) {
				util.Loggrs.Info(b.ID, " DB CONTINUE")
				continue loop
			}
		}
		for _, table := range tableListAgent {
			if !tools.StrInSlice(table, tableList) {
				util.Loggrs.Info(b.ID, " TABLE CONTINUE")
				continue loop
			}
		}

		util.Loggrs.Info(b.ID, " RUN")

		//当Allow Conditions中只有1个用户，只有1个权限, item扔掉（回收）
		if len(userList) == 1 && len(accessList) == 1 && len(tableList) < 2 {
			util.Loggrs.Info("单个用户，单个权限，单个库表")
			for _, permit := range SplitSlice(a.Permission) {
				if strings.ToLower(accessList[0]) == permit {
					continue loop
				}
			}
		}

		//当Allow Conditions中有2个以上用户，只有一个权限, 用户移除（回收）
		if len(userList) >= 2 && len(accessList) == 1 {
			util.Loggrs.Info("多个用户，单个权限，不管库表")
			var acconts []string
			for _, user := range SplitSlice(a.User) {
				for _, accont := range b.PolicyItems[i].Users {
					if user != accont {
						acconts = append(acconts, accont)
					}
				}
			}
			b.PolicyItems[i].Users = acconts
		}

		//当Allow Conditions中只有1个用户，有2个以上权限, 权限替换（回收）
		if len(userList) == 1 && len(accessList) >= 2 && len(tableList) < 2 {
			util.Loggrs.Info("单个用户，多个权限，单个库表")
			var Access []struct {
				Type      string `json:"type"`
				IsAllowed bool   `json:"isAllowed"`
			}
			for _, permit := range SplitSlice(a.Permission) {
				for y := 0; y < len(b.PolicyItems[i].Accesses); y++ {
					if strings.ToLower(b.PolicyItems[i].Accesses[y].Type) == permit {
						Access = append(Access, b.PolicyItems[i].Accesses[y])
					}
				}
			}
			b.PolicyItems[i].Accesses = accessSlice(Access, b.PolicyItems[i].Accesses)
		}

		//当Allow Conditions中有2个以上用户，有2个以上权限, item拆分（回收）
		if len(userList) >= 2 && len(accessList) >= 2 && len(tableList) < 2 {
			util.Loggrs.Info("多个用户，多个权限，单个库表")
			var Access []struct {
				Type      string `json:"type"`
				IsAllowed bool   `json:"isAllowed"`
			}
			items := changeItem(b, i, a.Permission, a.User, Access)
			policyItems = append(policyItems, items...)
		}

		//当Allow Conditions中只有1个用户，有2个以上权限，表名有多个时，拆分+新增policy struct
		if len(userList) == 1 && len(accessList) >= 2 && len(tableList) >= 2 {
			util.Loggrs.Info("单个用户，多个权限，多个库表")
			//【重构】发起一个新的新增请求
			err := ChangeSplit2Add(a, b, i)
			if err != nil {
				util.Loggrs.Warn(err.Error())
			}
			//【忽略】- 不要后面的access了
			continue loop
		}

		//当有多个用户，多个权限，多个库表
		if len(userList) >= 2 && len(accessList) >= 2 && len(tableList) >= 2 {
			util.Loggrs.Info("多个用户，多个权限，多个库表")
			//【重构】发起一个新的新增请求
			err := ChangeSplit2Add2(a, b, i)
			if err != nil {
				util.Loggrs.Warn(err.Error())
			}
			//【替换】将原有的信息进行replace
			var acconts []string
			for _, user := range SplitSlice(a.User) {
				for _, accont := range b.PolicyItems[i].Users {
					if user != accont {
						acconts = append(acconts, accont)
					}
				}
			}
			b.PolicyItems[i].Users = acconts
		}

		policyItems = append(policyItems, b.PolicyItems[i])
	}

	b.PolicyItems = rmDuplicates(policyItems)

	var policy string
	if len(b.PolicyItems) >= 1 {
		util.Loggrs.Info(b.ID, "【命中】 ", len(b.PolicyItems))

		marshal, err := json.Marshal(&b)
		if err != nil {
			util.Loggrs.Error(err.Error())
		}
		//剔除无效的信息
		policy = strings.NewReplacer("},}}", "}}}").Replace(ReplaceResource(b, string(marshal)))
	}
	return policy
}

type accessStruct []struct {
	Type      string `json:"type"`
	IsAllowed bool   `json:"isAllowed"`
}

// 不同的切片找到不同
func accessSlice(slice1, slice2 accessStruct) accessStruct {
	// 创建一个map用于快速查找
	map1 := make(map[string]bool)
	for _, item := range slice1 {
		key := fmt.Sprintf("%s-%t", item.Type, item.IsAllowed)
		map1[key] = true
	}
	//
	var m accessStruct
	// 遍历slice2并查找不同的元素
	for _, item := range slice2 {
		key := fmt.Sprintf("%s-%t", item.Type, item.IsAllowed)
		if !map1[key] {
			// 如果元素不在slice1的map中，或者元素不同，打印出来
			m = append(m, item)
		}
	}
	return m
}
