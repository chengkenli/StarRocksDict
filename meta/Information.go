/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package meta
 *@file    MetaInformation
 *@date    2025/5/16 9:40
 */

package meta

import (
	"StarRocksDict/conn"
	"StarRocksDict/util"
	"fmt"
	"time"
)

var infomation_cache []map[string]string

func Information() {
	_init()
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			_init()
		}
	}
}

// Information
// 加载information_schema数据
func _init() {
	for _, m := range util.MetaLink {
		db, err := conn.StarRocks(m["app"].(string))
		if err != nil {
			util.Loggrs.Warn(err.Error())
			continue
		}
		var result []map[string]interface{}
		r := db.Raw("select TABLE_SCHEMA,TABLE_NAME,TABLE_TYPE from information_schema.tables").Scan(&result)
		if r.Error != nil {
			util.Loggrs.Error(r.Error.Error())
			continue
		}

		cache := make(map[string]string)
		for _, m2 := range result {
			schema := fmt.Sprintf("%s.%s.%s", m["app"].(string), m2["TABLE_SCHEMA"].(string), m2["TABLE_NAME"].(string))
			engine := m2["TABLE_TYPE"].(string)
			cache[schema] = engine
		}
		infomation_cache = append(infomation_cache, cache)
		util.Loggrs.Info(fmt.Sprintf("loading meta_information is success. [%s]:[%d]", m["app"].(string), len(result)))
	}
}

func Getcache(key string) string {
	for _, m := range infomation_cache {
		if m[key] != "" {
			return m[key]
		}
	}
	return ""
}
