/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangUpdate
 *@date    2024/11/29 9:40
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"encoding/json"
)

func setuAddRangUpdateIsType(item []byte) string {
	var b policyStructsAll
	err := json.Unmarshal(item, &b)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return ""
	}
	for _, c := range b {
		if c.Resources.Catalog.Values != nil && c.Resources.Database.Values == nil {
			return "catalog"
		}
		if c.Resources.Database.Values != nil && c.Resources.Table.Values == nil || c.Resources.MaterializedView.Values == nil || c.Resources.View.Values == nil {
			return "database"
		}
		if c.Resources.Table.Values != nil {
			return "table"
		}
		if c.Resources.MaterializedView.Values != nil {
			return "materialized_view"
		}
		if c.Resources.View.Values != nil {
			return "view"
		}
		if c.Resources.Column.Values != nil {
			return "column"
		}
		if c.Resources.GlobalFunction.Values != nil {
			return "global_function"
		}
		if c.Resources.Function.Values != nil {
			return "function"
		}
		if c.Resources.MaskingPolicy.Values != nil {
			return "masking_policy"
		}
		if c.Resources.StorageVolume.Values != nil {
			return "storage_volume"
		}
		if c.Resources.FailoverGroup.Values != nil {
			return "failover_group"
		}
		if c.Resources.ResourceGroup.Values != nil {
			return "resource_group"
		}
		if c.Resources.System.Values != nil {
			return "system"
		}
		if c.Resources.RowAccessPolicy.Values != nil {
			return "row_access_policy"
		}
		if c.Resources.Warehouse.Values != nil {
			return "warehouse"
		}
	}
	return ""
}
