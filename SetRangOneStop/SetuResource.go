/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangReplaceResource
 *@date    2024/11/29 10:10
 */

package SetRangeOneStop

import (
	"fmt"
	"strings"
)

func ReplaceResource(c policyStructSingle, policy string) string {
	//重新组合
	if c.Resources.Catalog.Values == nil {
		policy = sourceRep("catalog", policy)
	}
	if c.Resources.Database.Values == nil {
		policy = sourceRep("database", policy)
	}
	if c.Resources.Table.Values == nil {
		policy = sourceRep("table", policy)
	}
	if c.Resources.MaterializedView.Values == nil {
		policy = sourceRep("materialized_view", policy)
	}
	if c.Resources.View.Values == nil {
		policy = sourceRep("view", policy)
	}
	if c.Resources.Column.Values == nil {
		policy = sourceRep("column", policy)
	}
	if c.Resources.GlobalFunction.Values == nil {
		policy = sourceRep("global_function", policy)
	}
	if c.Resources.Function.Values == nil {
		policy = sourceRep("function", policy)
	}
	if c.Resources.MaskingPolicy.Values == nil {
		policy = sourceRep("masking_policy", policy)
	}
	if c.Resources.StorageVolume.Values == nil {
		policy = sourceRep("storage_volume", policy)
	}
	if c.Resources.FailoverGroup.Values == nil {
		policy = sourceRep("failover_group", policy)
	}
	if c.Resources.ResourceGroup.Values == nil {
		policy = sourceRep("resource_group", policy)
	}
	if c.Resources.System.Values == nil {
		policy = sourceRep("system", policy)
	}
	if c.Resources.RowAccessPolicy.Values == nil {
		policy = sourceRep("row_access_policy", policy)
	}
	if c.Resources.Warehouse.Values == nil {
		policy = sourceRep("warehouse", policy)
	}
	return policy
}

func sourceRep(sign, policy string) string {
	length := len(strings.Split(policy, `{"values":null,"isExcludes":false,"isRecursive":false}`))
	if length == 2 {
		policy = strings.NewReplacer(
			fmt.Sprintf(`"%s":{"values":null,"isExcludes":false,"isRecursive":false}`, sign), "",
		).Replace(policy)
	} else {
		policy = strings.NewReplacer(
			fmt.Sprintf(`"%s":{"values":null,"isExcludes":false,"isRecursive":false},`, sign), "",
		).Replace(policy)
	}
	return policy
}
