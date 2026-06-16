/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangPermissionGrant
 *@date    2024/11/29 10:05
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"regexp"
	"strings"
)

// PermissionGrant
// 授权逻辑
func PermissionGrant(a *util.SetPolicy) error {
	util.Loggrs.Info("POLICY ->:NOT EXIST")
	//*合并组成授权策略json*/
	bodys := changePolicy(a)
	for i, body := range bodys {
		/*检查账户是否存在*/
		err := existUser(a.User)
		if err != nil {
			util.Loggrs.Error(err.Error())
			continue
		}
		result, err := execPolicy(-1, "create", body)
		if err != nil {
			util.Loggrs.Error(err.Error())
			//如果 Another policy already exists for matching resource: policy-name=[all - system]
			if strings.Contains(err.Error(), "policy already exists") {

				util.Loggrs.Info("POLICY ->:EXIST ", i, " ")
				//解析policy name
				regex := regexp.MustCompile(`policy-name=\[(.*?)\]`).FindStringSubmatch(err.Error())
				//根据policy name找到相关的policy json
				util.Loggrs.Info(regex[1], " ", i, " ")
				policy, err := getPolicyByName(a.App, regex[1])
				if err != nil {
					util.Loggrs.Error(err.Error())
					return err
				}
				//提交更新
				var body []byte
				var id int
				switch setuAddRangUpdateIsType(policy) {
				//system类型
				case "system":
					id, body, err = updatePolicyBySystem(a, policy)
					if err != nil {
						util.Loggrs.Error(err.Error())
						return err
					}
				}
				//执行policy json
				_, err = execPolicy(id, "update", string(body))
				if err != nil {
					util.Loggrs.Error(err.Error())
					return err
				}
			}
		}
		if result == nil {
			util.Loggrs.Error("Failed: ", i, " ")
		} else {
			util.Loggrs.Info("Success: ", i, " ", result)
		}
	}
	return nil
}
