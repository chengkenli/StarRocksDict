/*
 *@author  chengkenli
 *@project deepseek_ai
 *@package lark
 *@file    init
 *@date    2025/4/10 14:12
 */

package lark

import (
	"StarRocksDict/conn"
	"StarRocksDict/util"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

var larkmeta []map[string]interface{}
var larkclient *resty.Client
var larkcache *cache.Cache
var appName, appid, appSecret, larkproxy string

type ttoken struct {
	Code   int    `json:"code"`
	Expire int    `json:"expire"`
	Msg    string `json:"msg"`
	Token  string `json:"tenant_access_token"`
}

func init() {
	appName = util.MetaConf["slow_query_lark_app"].(string)
	appid = util.MetaConf["slow_query_lark_appid"].(string)
	appSecret = util.MetaConf["slow_query_lark_appsecret"].(string)
	larkproxy = util.MetaConf["slow_query_proxy_feishu"].(string)

	larkclient = resty.New().SetProxy(larkproxy)
	larkcache = cache.New(20*time.Minute, 20*time.Minute)

	larkdata := util.Config.GetString("Schema.larkMeta")
	if larkdata == "" {
		return
	}
	db, err := conn.StarRocks("sr-adhoc")
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
	r := db.Raw(fmt.Sprintf("select * from %s", larkdata)).Scan(&larkmeta)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return
	}
	if len(larkmeta) > 0 {
		util.Loggrs.Info("loading lark_information is success. ", len(larkmeta))
	}
}

// Lark2userid
// Lark2userid 通过userid获取lark信息
func Lark2userid(userid string) map[string]interface{} {
	for _, information := range larkmeta {
		if information["user_id"].(string) == strings.NewReplacer(" ", "").Replace(userid) {
			return information
		}
	}
	return nil
}

// Lark2openid 通过openid获取lark信息
func lark2openid(openid string) map[string]interface{} {
	for _, information := range larkmeta {
		if information["open_id"].(string) == strings.NewReplacer(" ", "").Replace(openid) {
			return information
		}
	}
	return nil
}
