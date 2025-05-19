/*
 *@author  chengkenli
 *@project deepseek_ai
 *@package lark
 *@file    lark_token
 *@date    2025/4/10 13:55
 */

package lark

import (
	"StarRocksDict/util"
	"encoding/json"
	"errors"
	"github.com/patrickmn/go-cache"
)

// GetTenantAccessToken 获取访问凭证 tenant_access_token
func getToken() (string, error) {
	v, Ok := larkcache.Get(appName)
	if Ok {
		util.Loggrs.Info("return [cache] token.")
		return v.(string), nil
	}
	//发送POST请求并处理响应
	respones, err := larkclient.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(map[string]string{
			"app_id":     appid,
			"app_secret": appSecret,
		}).Post("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal")
	if err != nil {
		return "", err
	}
	code, _ := getCode(respones.Body())
	if code != 0 {
		util.Loggrs.Warn(string(respones.Body()))
		return "", errors.New(string(respones.Body()))
	}

	var token ttoken
	err = json.Unmarshal(respones.Body(), &token)
	if err != nil {
		return "", err
	}
	larkcache.Set(appName, token.Token, cache.DefaultExpiration)

	util.Loggrs.Info("return [real] token.")
	return token.Token, nil
}

func getCode(body []byte) (int64, error) {
	// 创建一个map来存储解析后的数据
	data := make(map[string]interface{})
	// 解析JSON字符串到map中
	err := json.Unmarshal(body, &data)
	if err != nil {
		return -1, err
	}
	// 获取code字段的值
	code, ok := data["code"].(float64) // JSON中的数字默认解析为float64
	if !ok {
		return -1, errors.New("code failed")
	}
	return int64(code), nil
}
