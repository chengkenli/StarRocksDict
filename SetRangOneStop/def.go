/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    def
 *@date    2024/11/25 18:11
 */

package SetRangeOneStop

import (
	"StarRocksDict/conn"
	"StarRocksDict/meta"
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type policyStructsAll []struct {
	ID             int    `json:"id"`
	GUID           string `json:"guid"`
	IsEnabled      bool   `json:"isEnabled"`
	Version        int    `json:"version"`
	Service        string `json:"service"`
	Name           string `json:"name"`
	PolicyType     int    `json:"policyType"`
	PolicyPriority int    `json:"policyPriority"`
	Description    string `json:"description"`
	IsAuditEnabled bool   `json:"isAuditEnabled"`
	PolicyItems    []struct {
		Accesses []struct {
			Type      string `json:"type"`
			IsAllowed bool   `json:"isAllowed"`
		} `json:"accesses"`
		Users         []string      `json:"users"`
		Groups        []interface{} `json:"groups"`
		Roles         []interface{} `json:"roles"`
		Conditions    []interface{} `json:"conditions"`
		DelegateAdmin bool          `json:"delegateAdmin"`
	} `json:"policyItems"`
	DenyPolicyItems      []interface{} `json:"denyPolicyItems"`
	AllowExceptions      []interface{} `json:"allowExceptions"`
	DenyExceptions       []interface{} `json:"denyExceptions"`
	DataMaskPolicyItems  []interface{} `json:"dataMaskPolicyItems"`
	RowFilterPolicyItems []interface{} `json:"rowFilterPolicyItems"`
	ServiceType          string        `json:"serviceType"`
	Options              struct {
	} `json:"options"`
	ValiditySchedules []interface{} `json:"validitySchedules"`
	PolicyLabels      []interface{} `json:"policyLabels"`
	ZoneName          string        `json:"zoneName"`
	IsDenyAllElse     bool          `json:"isDenyAllElse"`
	Resources         struct {
		//Catalog
		Catalog struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"catalog"`
		//Database
		Database struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"database"`
		//Table
		Table struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"table"`
		//MaterializedView
		MaterializedView struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"materialized_view"`
		//View
		View struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"view"`
		//Column
		Column struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"column"`
		//GlobalFunction
		GlobalFunction struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"global_function"`
		//Function
		Function struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"function"`
		//MaskingPolicy
		MaskingPolicy struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"masking_policy"`
		//StorageVolume
		StorageVolume struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"storage_volume"`
		//FailoverGroup
		FailoverGroup struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"failover_group"`
		//ResourceGroup
		ResourceGroup struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"resource_group"`
		//System
		System struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"system"`
		//RowAccessPolicy
		RowAccessPolicy struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"row_access_policy"`
		//Warehouse
		Warehouse struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"warehouse"`
	} `json:"resources"`
}

type policyStructSingle struct {
	ID             int    `json:"id"`
	GUID           string `json:"guid"`
	IsEnabled      bool   `json:"isEnabled"`
	Version        int    `json:"version"`
	Service        string `json:"service"`
	Name           string `json:"name"`
	PolicyType     int    `json:"policyType"`
	PolicyPriority int    `json:"policyPriority"`
	Description    string `json:"description"`
	IsAuditEnabled bool   `json:"isAuditEnabled"`
	PolicyItems    []struct {
		Accesses []struct {
			Type      string `json:"type"`
			IsAllowed bool   `json:"isAllowed"`
		} `json:"accesses"`
		Users         []string      `json:"users"`
		Groups        []interface{} `json:"groups"`
		Roles         []interface{} `json:"roles"`
		Conditions    []interface{} `json:"conditions"`
		DelegateAdmin bool          `json:"delegateAdmin"`
	} `json:"policyItems"`
	DenyPolicyItems      []interface{} `json:"denyPolicyItems"`
	AllowExceptions      []interface{} `json:"allowExceptions"`
	DenyExceptions       []interface{} `json:"denyExceptions"`
	DataMaskPolicyItems  []interface{} `json:"dataMaskPolicyItems"`
	RowFilterPolicyItems []interface{} `json:"rowFilterPolicyItems"`
	ServiceType          string        `json:"serviceType"`
	Options              struct {
	} `json:"options"`
	ValiditySchedules []interface{} `json:"validitySchedules"`
	PolicyLabels      []interface{} `json:"policyLabels"`
	ZoneName          string        `json:"zoneName"`
	IsDenyAllElse     bool          `json:"isDenyAllElse"`
	Resources         struct {
		//Catalog
		Catalog struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"catalog"`
		//Database
		Database struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"database"`
		//Table
		Table struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"table"`
		//MaterializedView
		MaterializedView struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"materialized_view"`
		//View
		View struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"view"`
		//Column
		Column struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"column"`
		//GlobalFunction
		GlobalFunction struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"global_function"`
		//Function
		Function struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"function"`
		//MaskingPolicy
		MaskingPolicy struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"masking_policy"`
		//StorageVolume
		StorageVolume struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"storage_volume"`
		//FailoverGroup
		FailoverGroup struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"failover_group"`
		//ResourceGroup
		ResourceGroup struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"resource_group"`
		//System
		System struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"system"`
		//RowAccessPolicy
		RowAccessPolicy struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"row_access_policy"`
		//Warehouse
		Warehouse struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"warehouse"`
	} `json:"resources"`
}

type policyItem []struct {
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

type policyStruct struct {
	ID             int    `json:"id"`
	GUID           string `json:"guid"`
	IsEnabled      bool   `json:"isEnabled"`
	Version        int    `json:"version"`
	Service        string `json:"service"`
	Name           string `json:"name"`
	PolicyType     int    `json:"policyType"`
	PolicyPriority int    `json:"policyPriority"`
	Description    string `json:"description"`
	IsAuditEnabled bool   `json:"isAuditEnabled"`
	Resources      struct {
		Catalog struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"catalog"`
		Database struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"database"`
		Table struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"table"`
	} `json:"resources"`
	PolicyItems []struct {
		Accesses []struct {
			Type      string `json:"type"`
			IsAllowed bool   `json:"isAllowed"`
		} `json:"accesses"`
		Users         []string      `json:"users"`
		Groups        []string      `json:"groups"`
		Roles         []interface{} `json:"roles"`
		Conditions    []interface{} `json:"conditions"`
		DelegateAdmin bool          `json:"delegateAdmin"`
	} `json:"policyItems"`
	DenyPolicyItems      []interface{} `json:"denyPolicyItems"`
	AllowExceptions      []interface{} `json:"allowExceptions"`
	DenyExceptions       []interface{} `json:"denyExceptions"`
	DataMaskPolicyItems  []interface{} `json:"dataMaskPolicyItems"`
	RowFilterPolicyItems []interface{} `json:"rowFilterPolicyItems"`
	ServiceType          string        `json:"serviceType"`
	Options              struct {
	} `json:"options"`
	ValiditySchedules []interface{} `json:"validitySchedules"`
	PolicyLabels      []interface{} `json:"policyLabels"`
	ZoneName          string        `json:"zoneName"`
	IsDenyAllElse     bool          `json:"isDenyAllElse"`
}

type policyStructBySystem struct {
	ID             int    `json:"id"`
	GUID           string `json:"guid"`
	IsEnabled      bool   `json:"isEnabled"`
	Version        int    `json:"version"`
	Service        string `json:"service"`
	Name           string `json:"name"`
	PolicyType     int    `json:"policyType"`
	PolicyPriority int    `json:"policyPriority"`
	Description    string `json:"description"`
	IsAuditEnabled bool   `json:"isAuditEnabled"`
	Resources      struct {
		//System
		System struct {
			Values      []string `json:"values"`
			IsExcludes  bool     `json:"isExcludes"`
			IsRecursive bool     `json:"isRecursive"`
		} `json:"system"`
	} `json:"resources"`
	PolicyItems []struct {
		Accesses []struct {
			Type      string `json:"type"`
			IsAllowed bool   `json:"isAllowed"`
		} `json:"accesses"`
		Users         []string      `json:"users"`
		Groups        []string      `json:"groups"`
		Roles         []interface{} `json:"roles"`
		Conditions    []interface{} `json:"conditions"`
		DelegateAdmin bool          `json:"delegateAdmin"`
	} `json:"policyItems"`
	DenyPolicyItems      []interface{} `json:"denyPolicyItems"`
	AllowExceptions      []interface{} `json:"allowExceptions"`
	DenyExceptions       []interface{} `json:"denyExceptions"`
	DataMaskPolicyItems  []interface{} `json:"dataMaskPolicyItems"`
	RowFilterPolicyItems []interface{} `json:"rowFilterPolicyItems"`
	ServiceType          string        `json:"serviceType"`
	Options              struct {
	} `json:"options"`
	ValiditySchedules []interface{} `json:"validitySchedules"`
	PolicyLabels      []interface{} `json:"policyLabels"`
	ZoneName          string        `json:"zoneName"`
	IsDenyAllElse     bool          `json:"isDenyAllElse"`
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

type accesses struct {
	Type      string `json:"type"`
	IsAllowed bool   `json:"isAllowed"`
}

type signs struct {
	Database, Table, Materialized, View []string
}

func Code(body []byte, key string) (int64, error) {
	// 创建一个map来存储解析后的数据
	data := make(map[string]interface{})
	// 解析JSON字符串到map中
	err := json.Unmarshal(body, &data)
	if err != nil {
		return -1, err
	}
	// 获取code字段的值
	code, ok := data[key].(float64) // JSON中的数字默认解析为float64
	if !ok {
		return -1, errors.New("code failed")
	}
	return int64(code), nil
}

// arrJsonDt
// 根据用户提交的json串将table字段拆分,返回database,table两个不同的数组
func arrJsonDt(app, value string) signs {

	value = strings.ReplaceAll(value, " ", "")
	var database, table, materialized, view []string

	for _, schema := range strings.Split(value, ",") {
		splt := strings.Split(schema, ".")
		//split
		database = append(database, splt[0])
		switch juct(app, schema) {
		case "VIEW":
			view = append(view, splt[1])
		case "MATERIALIZED VIEW":
			materialized = append(materialized, splt[1])
		case "BASE TABLE", "UNKNOWN":
			table = append(table, splt[1])
		case "DATA BASE":

		default:
			table = append(table, splt[1])
		}
	}
	database = append(database, fmt.Sprintf("%s_%s", util.KeySingle, uuid.New().String()))
	return signs{
		Database:     database,
		Table:        table,
		Materialized: materialized,
		View:         view,
	}
}

// juct
// 判断数据表是内表，还是视图
func juct(app, table string) string {
	schema := strings.Split(strings.NewReplacer(" ", "").Replace(table), ".")

	if len(schema) == 1 || schema[1] == "*" {
		return "DATA BASE"
	}

	db, err := conn.StarRocks(app)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return ""
	}
	/*每次使用完，主动关闭连接数*/
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			util.Loggrs.Error(err.Error())
			return
		}
		sqlDB.SetMaxOpenConns(30)                  //最大连接数
		sqlDB.SetMaxIdleConns(30)                  //最大空闲连接数
		sqlDB.SetConnMaxLifetime(30 * time.Second) //空闲连接最多存活时间
		sqlDB.Close()
	}()

	engine := meta.Getcache(fmt.Sprintf("%s.%s.%s", app, schema[0], schema[1]))
	util.Loggrs.Info(fmt.Sprintf("getdata meta_information [%s],[%s.%s],[%s]", app, schema[0], schema[1], engine))

	tableType := engine
	switch tableType {
	case "VIEW":
		var m map[string]interface{}
		r := db.Raw("show create table " + table).Scan(&m)
		if r.Error != nil {
			util.Loggrs.Error(err.Error())
		}
		if _, ok := m["Create Materialized View"]; ok {
			tableType = "MATERIALIZED VIEW"
		} else {
			tableType = "VIEW"
		}
	}
	util.Loggrs.Info(tableType, " ->:", table)
	return tableType
}

// policyJson
// 格式化将数组打印成带双引号格式
func policyJson(data []string) string {
	marshal, err := json.Marshal(&data)
	if err != nil {
		util.Loggrs.Error(err.Error())
	}
	return string(marshal)
}

// formatColData
// 格式化授权字段,columns格式 [{"type": "select","isAllowed": true},{"type": "insert","isAllowed": true}]
func formatColData(data string) string {
	var d []string
	for _, p := range strings.Split(data, ",") {
		d = append(d, fmt.Sprintf(`{"type": "%s","isAllowed": true}`, p))
	}
	b := strings.NewReplacer(
		`["`, `[`,
		`"]`, `]`,
		`\`, ``,
		`}","{`, `},{`,
	).Replace(policyJson(d))

	return b
}

func rmDuplicates(policyItems []struct {
	Accesses []struct {
		Type      string `json:"type"`
		IsAllowed bool   `json:"isAllowed"`
	} `json:"accesses"`
	Users         []string      `json:"users"`
	Groups        []interface{} `json:"groups"`
	Roles         []interface{} `json:"roles"`
	Conditions    []interface{} `json:"conditions"`
	DelegateAdmin bool          `json:"delegateAdmin"`
}) []struct {
	Accesses []struct {
		Type      string `json:"type"`
		IsAllowed bool   `json:"isAllowed"`
	} `json:"accesses"`
	Users         []string      `json:"users"`
	Groups        []interface{} `json:"groups"`
	Roles         []interface{} `json:"roles"`
	Conditions    []interface{} `json:"conditions"`
	DelegateAdmin bool          `json:"delegateAdmin"`
} {
	uniqueItems := make([]struct {
		Accesses []struct {
			Type      string `json:"type"`
			IsAllowed bool   `json:"isAllowed"`
		} `json:"accesses"`
		Users         []string      `json:"users"`
		Groups        []interface{} `json:"groups"`
		Roles         []interface{} `json:"roles"`
		Conditions    []interface{} `json:"conditions"`
		DelegateAdmin bool          `json:"delegateAdmin"`
	}, 0)
	itemMap := make(map[string]bool)

	for _, item := range policyItems {
		itemJSON, _ := json.Marshal(item)
		if _, exists := itemMap[string(itemJSON)]; !exists {
			uniqueItems = append(uniqueItems, item)
			itemMap[string(itemJSON)] = true
		}
	}
	return uniqueItems
}

func SplitSlice(slice string) []string {
	return tools.RmSliceStr(strings.Split(strings.ToLower(slice), ","))
}
