/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangExecPolicy
 *@date    2024/11/28 21:04
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"errors"
	"fmt"
	"strings"
)

// 执行policy操作，create || update
func execPolicy(id int, action, body string) (interface{}, error) {
	var (
		r   []byte
		err error
	)
	switch action {
	case "create":
		/*提交create授权策略*/
		util.Loggrs.Info("CREATE POLICY ->:", body)
		r, err = Post(
			&util.RangerHost{
				Server: util.Read.Ranger.Host,
				Access: util.Read.Ranger.User,
				Secret: util.Read.Ranger.Password,
			}, "POST", "/service/public/v2/api/policy", strings.NewReader(body))
	case "update":
		/*提交update授权策略*/
		util.Loggrs.Info("UPDATE POLICY ->:", body)
		r, err = Post(
			&util.RangerHost{
				Server: util.Read.Ranger.Host,
				Access: util.Read.Ranger.User,
				Secret: util.Read.Ranger.Password,
			}, "PUT", fmt.Sprintf(`/service/public/v2/api/policy/%d`, id), strings.NewReader(body))
	case "delete":
		/*提交update授权策略*/
		util.Loggrs.Info("DELETE POLICY ->:", body)
		r, err = Post(
			&util.RangerHost{
				Server: util.Read.Ranger.Host,
				Access: util.Read.Ranger.User,
				Secret: util.Read.Ranger.Password,
			}, "DELETE", fmt.Sprintf(`/service/public/v2/api/policy/%d`, id), nil)
		r = []byte("Ok")
	}
	if len(r) == 0 || err != nil || strings.Contains(string(r), "Invalid") || strings.Contains(string(r), "Validation") || strings.Contains(string(r), "Failed") {
		return nil, errors.New("EXEC FAILED ->:" + string(r))
	}
	util.Loggrs.Info("POLICY RESULT ->:", string(r))

	return string(r), nil
}
