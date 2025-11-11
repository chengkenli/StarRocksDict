package SetDictOneStop

import (
	"StarRocksDict/meta"
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type dictStruct struct {
	dictH  string
	dictE  string
	db     *gorm.DB
	e      dictJson
	table  string
	user   string
	policy string
}

func SetuAddDict3(t *dictStruct, once *sync.Once) []util.SetResult {
	var mutex sync.Mutex
	//匹配库级别
	var tableResult []util.SetResult
	if strings.Contains(t.table, ".*") {
		for _, p := range strings.Split(t.e.Permission, ",") {
			schem := tools.Schem{
				Db:         t.db,
				Catalog:    t.e.Catalog,
				User:       t.user,
				Table:      t.table,
				Permission: p,
				App:        t.e.Starrocks,
				Policy:     t.policy,
			}
			once.Do(func() {
				schem := tools.Schem{
					Db:         t.db,
					Catalog:    t.e.Catalog,
					User:       t.user,
					Table:      "ALL CATALOGS",
					Permission: "USAGE",
					App:        t.e.Starrocks,
					Policy:     t.policy,
				}
				//默认赋予USAGE ALL CATALOGS
				if strings.Contains(strings.ToLower(t.dictH), "grant") {
					util.Loggrs.Info(t.dictH)
					tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`%s USAGE ON ALL CATALOGS %s USER %s`, t.dictH, t.dictE, t.user)))
				}
			})

			switch strings.ToUpper(p) {
			case "SELECT":
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON TABLE %s %s %s`, t.e.Catalog, t.dictH, p, t.table, t.dictE, t.user)))
				if t.e.Catalog == "default_catalog" {
					tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL VIEWS IN DATABASE %s %s USER %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))
					tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL MATERIALIZED VIEWS IN DATABASE %s %s USER %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))
				}
			case "INSERT", "EXPORT", "UPDATE", "DELETE", "ALTER", "DROP":
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON ALL TABLES IN DATABASE %s %s %s`, t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))
				//tableResult = append(tableResult, tools.Onexec(&mutex,&schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON DATABASE %s %s %s`, t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))
				if t.e.Catalog == "default_catalog" {
					if strings.ToUpper(p) == "DROP" || strings.ToUpper(p) == "ALTER" {
						tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL VIEWS IN DATABASE %s %s USER %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))
						tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL MATERIALIZED VIEWS IN DATABASE %s %s USER %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))
					}
				}
			case "CREATE TABLE", "CREATE VIEW", "CREATE FUNCTION", "CREATE MATERIALIZED VIEW", "CREATE ROW ACCESS POLICY", "CREATE MASKING POLICY":
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON DATABASE %s %s %s`, t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))

			case "CREATE_TABLE", "CREATE_VIEW", "CREATE_FUNCTION", "CREATE_MATERIALIZED_VIEW", "CREATE_ROW_ACCESS_POLICY", "CREATE_MASKING_POLICY":
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON DATABASE %s %s %s`, t.e.Catalog, t.dictH, strings.ReplaceAll(p, "_", " "), strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))

			case "REFRESH":
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s %s ON ALL MATERIALIZED VIEWS IN DATABASE %s %s USER %s", t.e.Catalog, t.dictH, p, strings.ReplaceAll(t.table, ".*", ""), t.dictE, t.user)))

			case "USAGE":
				if strings.Contains(strings.ToLower(t.dictH), "grant") {
					util.Loggrs.Info(t.dictH)
					tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON ALL CATALOGS %s USER %s`, t.e.Catalog, t.dictH, p, t.dictE, t.user)))
				}
			case "APPLY":
				PolicyName := strings.ReplaceAll(t.table, ".*", "")
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s APPLY ON ROW ACCESS POLICY %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s ALTER ON ROW ACCESS POLICY %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s DROP  ON ROW ACCESS POLICY %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s APPLY ON MASKING POLICY    %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s ALTER ON MASKING POLICY    %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s DROP  ON MASKING POLICY    %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			}
		}
	} else {
		//匹配表级别
		schem := tools.Schem{
			Db:         t.db,
			Catalog:    t.e.Catalog,
			User:       t.user,
			Table:      t.table,
			Permission: t.e.Permission,
			App:        t.e.Starrocks,
			Policy:     t.policy,
		}
		if t.e.Permission == "APPLY" {
			PolicyName := strings.ReplaceAll(t.table, ".*", "")
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s APPLY ON ROW ACCESS POLICY %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s ALTER ON ROW ACCESS POLICY %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s DROP  ON ROW ACCESS POLICY %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s APPLY ON MASKING POLICY    %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s ALTER ON MASKING POLICY    %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf("SET CATALOG %s;%s DROP  ON MASKING POLICY    %s %s USER %s", t.e.Catalog, t.dictH, PolicyName, t.dictE, t.user)))
			return tableResult
		}

		once.Do(func() {
			schem := tools.Schem{
				Db:         t.db,
				Catalog:    t.e.Catalog,
				User:       t.user,
				Table:      "ALL CATALOGS",
				Permission: "USAGE",
				App:        t.e.Starrocks,
				Policy:     t.policy,
			}
			//默认赋予USAGE ALL CATALOGS
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`%s USAGE ON ALL CATALOGS %s USER %s`, t.dictH, t.dictE, t.user)))
		})
		var TableType string
		if strings.ToLower(t.e.Catalog) == "default_catalog" {
			schema := strings.Split(t.table, ".")
			if len(schema) <= 1 {
				tableResult = append(tableResult, util.SetResult{
					App:     t.e.Starrocks,
					Catalog: t.e.Catalog,
					User:    t.user,
					Table:   t.table,
					Permit:  t.e.Permission,
					State:   "Fail",
					Demand:  "",
					Comment: "授权库表名不明确，如授权库级别请使用<db.*>，如授权表级别请<db.table>!",
				})
				return tableResult
			}

			engine := meta.Getcache(fmt.Sprintf("%s.%s.%s", t.e.Starrocks, schema[0], schema[1]))
			util.Loggrs.Info(fmt.Sprintf("getdata meta_information [%s],[%s.%s],[%s]", t.e.Starrocks, schema[0], schema[1], engine))

			if strings.ToLower(t.e.Catalog) == "default_catalog" && engine == "" {
				tableResult = append(tableResult, util.SetResult{
					App:     t.e.Starrocks,
					Catalog: t.e.Catalog,
					User:    t.user,
					Table:   t.table,
					Permit:  t.e.Permission,
					State:   "Fail",
					Demand:  "",
					Comment: "数据表不存在",
				})
				return tableResult
			}

			TableType = engine
		} else {
			TableType = "BASE TABLE"
		}

		switch TableType {
		case "VIEW":
			_, mv := tools.Materialized(t.db, t.table)
			if mv {
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON MATERIALIZED VIEW %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.user)))
			} else {
				tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON VIEW %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.user)))
			}
		case "BASE TABLE":
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON TABLE %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.user)))

		default:
			tableResult = append(tableResult, tools.Onexec(&mutex, &schem, fmt.Sprintf(`SET CATALOG %s;%s %s ON TABLE %s %s %s`, t.e.Catalog, t.dictH, t.e.Permission, t.table, t.dictE, t.user)))
		}
	}

	return tableResult
}
