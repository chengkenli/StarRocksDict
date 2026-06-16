/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuChangeSplit
 *@date    2024/12/3 10:46
 */

package SetRangeOneStop

import (
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"strings"
)

func ChangeSplit2Add2(a *util.SetPolicy, b policyStructSingle, i int) error {
	//拿到json传过来权限
	var rangAgent []string
	userAgent := SplitSlice(a.Permission)
	//拿到ranger返回的权限
	for y := 0; y < len(b.PolicyItems[i].Accesses); y++ {
		rangAgent = append(rangAgent, b.PolicyItems[i].Accesses[y].Type)
	}
	/////////////////////////////////
	//拿到库表
	var schema, unschema, diffschema []string

	for _, db := range b.Resources.Database.Values {
		//筛选库同表不同
		diffschema = append(diffschema, sliceStrdtb(b.Resources.Table.Values, db, a.Table)...)
		diffschema = append(diffschema, sliceStrdtb(b.Resources.MaterializedView.Values, db, a.Table)...)
		diffschema = append(diffschema, sliceStrdtb(b.Resources.View.Values, db, a.Table)...)

		//筛出不同库的
		unschema = append(unschema, sliceStrdb(b.Resources.Table.Values, db, a.Table)...)
		unschema = append(unschema, sliceStrdb(b.Resources.MaterializedView.Values, db, a.Table)...)
		unschema = append(unschema, sliceStrdb(b.Resources.View.Values, db, a.Table)...)

		//筛出需要隔离的库表
		schema = append(schema, sliceStrtb(b.Resources.Table.Values, db, a.Table)...)
		schema = append(schema, sliceStrtb(b.Resources.MaterializedView.Values, db, a.Table)...)
		schema = append(schema, sliceStrtb(b.Resources.View.Values, db, a.Table)...)
	}

	err := PermissionGrant(
		&util.SetPolicy{
			App:        a.App,
			Catalog:    a.Catalog,
			User:       a.User,
			Table:      strings.Join(tools.RmSliceStr(schema), ","),
			Permission: strings.Join(tools.FindDiffSlices(rangAgent, userAgent), ","),
			Revoke:     false,
		})
	if err != nil {
		util.Loggrs.Error(err.Error())
		return err
	}

	err = PermissionGrant(
		&util.SetPolicy{
			App:        a.App,
			Catalog:    a.Catalog,
			User:       a.User,
			Table:      strings.Join(tools.RmSliceStr(unschema), ","),
			Permission: strings.Join(rangAgent, ","),
			Revoke:     false,
		})
	if err != nil {
		util.Loggrs.Error(err.Error())
		return err
	}

	err = PermissionGrant(
		&util.SetPolicy{
			App:        a.App,
			Catalog:    a.Catalog,
			User:       a.User,
			Table:      strings.Join(tools.RmSliceStr(diffschema), ","),
			Permission: strings.Join(rangAgent, ","),
			Revoke:     false,
		})
	if err != nil {
		util.Loggrs.Error(err.Error())
		return err
	}

	return nil
}
