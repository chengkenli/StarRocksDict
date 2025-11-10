package SetDictOneStop

import (
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"fmt"
	"strings"
)

func SetuAddDict2(t dictStruct) []util.SetResult {
	//匹配库级别
	var tableResult []util.SetResult
	schem := tools.Schem{
		Db:         t.db,
		Catalog:    t.e.Catalog,
		User:       t.user,
		Table:      t.table,
		Permission: t.e.Permission,
		App:        t.e.Starrocks,
		Policy:     t.policy,
	}

	var pn []string
	for _, p := range strings.Split(t.e.Permission, ",") {
		pn = append(pn, strings.ToUpper(fmt.Sprintf("%s_priv", p)))
	}

	for _, table := range strings.Split(t.table, ",") {
		tableResult = append(tableResult, tools.OnExec20(&schem, fmt.Sprintf("%s %s ON %s %s '%s'", t.dictH, strings.Join(pn, ","), table, t.dictE, t.user)))
	}
	return tableResult
}
