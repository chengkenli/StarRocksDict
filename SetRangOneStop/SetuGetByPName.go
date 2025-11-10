/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangGetPolicyByPolicyName
 *@date    2024/11/28 21:10
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// getPolicyByName
// 获取指定服务下的策略信息列表,根据policy name获取到内容返回
func getPolicyByName(app, name string) ([]byte, error) {
	r, err := Post(
		&util.RangerHost{
			Server: util.Read.Ranger.Host,
			Access: util.Read.Ranger.User,
			Secret: util.Read.Ranger.Password,
		}, "GET", fmt.Sprintf("/service/public/v2/api/service/%s/policy", app), nil)
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
	for _, c := range b {
		if c.Name != name {
			continue
		}

		marshal, err := json.Marshal(&c)
		if err != nil {
			util.Loggrs.Error(err.Error())
			return nil, err
		}
		policy := string(marshal)
		util.Loggrs.Info("SS->:", policy)

		//重新组合
		if c.Resources.Catalog.Values == nil {
			policy = strings.NewReplacer(`"catalog":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.Database.Values == nil {
			policy = strings.NewReplacer(`"database":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.Table.Values == nil {
			policy = strings.NewReplacer(`"table":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.MaterializedView.Values == nil {
			policy = strings.NewReplacer(`"materialized_view":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.View.Values == nil {
			policy = strings.NewReplacer(`"view":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.Column.Values == nil {
			policy = strings.NewReplacer(`"column":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.GlobalFunction.Values == nil {
			policy = strings.NewReplacer(`"global_function":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.Function.Values == nil {
			policy = strings.NewReplacer(`"function":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.MaskingPolicy.Values == nil {
			policy = strings.NewReplacer(`"masking_policy":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.StorageVolume.Values == nil {
			policy = strings.NewReplacer(`"storage_volume":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.FailoverGroup.Values == nil {
			policy = strings.NewReplacer(`"failover_group":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.ResourceGroup.Values == nil {
			policy = strings.NewReplacer(`"resource_group":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.System.Values == nil {
			policy = strings.NewReplacer(`"system":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.RowAccessPolicy.Values == nil {
			policy = strings.NewReplacer(`"row_access_policy":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		if c.Resources.Warehouse.Values == nil {
			policy = strings.NewReplacer(`"warehouse":{"values":null,"isExcludes":false,"isRecursive":false},`, "").Replace(policy)
		}
		util.Loggrs.Info("GET POLICY ->:", policy)
		return []byte(policy), nil
	}
	return nil, errors.New("POLICY IS NULL")
}
