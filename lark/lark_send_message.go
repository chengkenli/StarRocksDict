/*
 *@author  chengkenli
 *@project deepseek_ai
 *@package lark
 *@file    lark_send_message
 *@date    2025/4/10 13:54
 */

package lark

import (
	"StarRocksDict/util"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Send2userid 根据userid发送信息给用户
// @用户名 userid
func Send2userid(userid, message string) error {
	token, err := getToken()
	if err != nil {
		util.Loggrs.Error(err.Error())
		return err
	}
	meta := Lark2userid(userid)
	if meta == nil {
		return errors.New(fmt.Sprintf("%s openid is nil.", userid))
	}
	//发送POST请求并处理响应
	respones, err := larkclient.R().
		SetHeader("Authorization", "Bearer "+token).
		SetBody(template(message, meta["open_id"].(string))).
		Post("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id")
	if err != nil {
		util.Loggrs.Error(err.Error())
		return err
	}
	code, _ := getCode(respones.Body())
	if code != 0 {
		util.Loggrs.Warn(string(respones.Body()))
		return errors.New(string(respones.Body()))
	}

	util.Loggrs.Info(string(respones.Body()))
	util.Loggrs.Info("send done.")
	return nil
}

func template(content, openid string) string {
	request, _ := json.Marshal(interface{}(fmt.Sprintf(`{"text":"%s"}`, content)))
	//util.Loggrs.Info(string(request))
	body := strings.NewReplacer(`\n`, `\\n`).Replace(string(request))
	tmplmap := fmt.Sprintf(`
{
	"receive_id": "%s",
    "content": %v,
    "msg_type": "text"
}`, openid, body)
	return tmplmap
}
