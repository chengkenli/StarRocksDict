package SetDictOneStop

import (
	"StarRocksDict/meta"
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type dictStructRole struct {
	dictH  string
	dictE  string
	db     *gorm.DB
	e      dictJsonRole
	table  string
	role   string
	policy string
}

func SetuAddDictRole3(t dictStructRole) []util.SetRole {
	fmt.Println(t)
	//匹配库级别
	var roleInfo []util.SetRole
	if strings.Contains(t.table, ".*") {
		for _, p := range strings.Split(t.e.Permission, ",") {
			schemrole := tools.SchemRole{
				Db:         t.db,
				Catalog:    t.e.Catalog,
				Role:       t.role,
				Table:      t.table,
				Permission: p,
				App:        t.e.Starrocks,
				Policy:     t.policy,
			}
			switch strings.ToUpper(p) {
			case "SELECT":
				roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON TABLE %s %s %s`, t.e.Catalog, t.dictH, p, t.table, t.dictE, t.role)))
				if t.e.Catalog == "default_catalog" {
					roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL VIEWS IN DATABASE %s %s %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))
					roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL MATERIALIZED VIEWS IN DATABASE %s %s %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))
				}
			case "INSERT", "EXPORT", "UPDATE", "DELETE", "ALTER", "DROP":
				roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON ALL TABLES IN DATABASE %s %s %s`, t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))
				if t.e.Catalog == "default_catalog" {
					if strings.ToUpper(p) == "DROP" || strings.ToUpper(p) == "ALTER" {
						roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL VIEWS IN DATABASE %s %s %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))
						roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL MATERIALIZED VIEWS IN DATABASE %s %s %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))
					}
				}
			case "CREATE TABLE", "CREATE VIEW", "CREATE FUNCTION", "CREATE MATERIALIZED VIEW":
				roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON DATABASE %s %s %s`, t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))

			case "CREATE_TABLE", "CREATE_VIEW", "CREATE_FUNCTION", "CREATE_MATERIALIZED_VIEW":
				roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON DATABASE %s %s %s`, t.e.Catalog, t.dictH, strings.ReplaceAll(p, "_", " "), strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))

			case "REFRESH":
				roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL MATERIALIZED VIEWS IN DATABASE %s %s %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.role)))

			case "USAGE":
				roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON ALL CATALOGS %s %s`, t.e.Catalog, t.dictH, p, t.dictE, t.role)))
			}
		}
	} else {
		//匹配表级别
		schemrole := tools.SchemRole{
			Db:         t.db,
			Catalog:    t.e.Catalog,
			Role:       t.role,
			Table:      t.table,
			Permission: t.e.Permission,
			App:        t.e.Starrocks,
			Policy:     t.policy,
		}

		var TableType string
		if strings.ToLower(t.e.Catalog) == "default_catalog" {
			schema := strings.Split(t.table, ".")
			if len(schema) <= 1 {
				roleInfo = append(roleInfo, util.SetRole{
					App:        t.e.Starrocks,
					Role:       t.role,
					Table:      t.table,
					Permission: t.e.Permission,
					Revoke:     t.e.Revoke,
					Catalog:    t.e.Catalog,
					Policyname: t.e.Policyname,
					State:      "Fail",
					Comment:    "授权库表名不明确，如授权库级别请使用<db.*>，如授权表级别请<db.table>!",
				})
				return roleInfo
			}

			engine := meta.Getcache(fmt.Sprintf("%s.%s.%s", t.e.Starrocks, schema[0], schema[1]))
			util.Loggrs.Info(fmt.Sprintf("getdata meta_information [%s],[%s.%s],[%s]", t.e.Starrocks, schema[0], schema[1], engine))

			if strings.ToLower(t.e.Catalog) == "default_catalog" && engine == "" {
				roleInfo = append(roleInfo, util.SetRole{
					App:        t.e.Starrocks,
					Role:       t.role,
					Table:      t.table,
					Permission: t.e.Permission,
					Revoke:     t.e.Revoke,
					Catalog:    t.e.Catalog,
					Policyname: t.e.Policyname,
					State:      "Fail",
					Comment:    "数据表不存在",
				})
				return roleInfo
			}

			TableType = engine
		} else {
			TableType = "BASE TABLE"
		}
		switch TableType {
		case "VIEW":
			roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON VIEW %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.role)))
			roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON MATERIALIZED VIEW %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.role)))

		case "BASE TABLE":
			roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON TABLE %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.role)))

		default:
			roleInfo = append(roleInfo, tools.OnExecRole(schemrole, fmt.Sprintf(`SET CATALOG %s;%s %s ON TABLE %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.role)))

		}
	}
	return roleInfo
}
