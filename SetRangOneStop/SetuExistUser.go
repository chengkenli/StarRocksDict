/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetRangeOneStop
 *@file    SetuAddRangExistUser
 *@date    2024/11/28 21:05
 */

package SetRangeOneStop

import (
	"StarRocksDict/util"
	"strings"
)

// arrSplitUser
// 判断ranger用户是否存在,如果不存在,那么新建
func existUser(user string) error {
	userList := strings.ReplaceAll(user, " ", "")

	if strings.Contains(userList, ",") {
		for _, user := range strings.Split(userList, ",") {
			r, err := getByName(user)
			if err != nil {
				return err
			}
			if strings.Contains(string(r), "is Not Found") {
				util.Loggrs.Warn("user not exist -> ", string(r))
				name, err := addByName(user)
				if err != nil {
					util.Loggrs.Error(err.Error())
					return err
				}
				util.Loggrs.Info("user add done -> ", string(name))
				return nil
			}
			return nil
		}
		return nil
	}

	r, err := getByName(user)
	if err != nil {
		return err
	}
	if strings.Contains(string(r), "is Not Found") {
		util.Loggrs.Warn("user not exist -> ", string(r))
		name, err := addByName(user)
		if err != nil {
			util.Loggrs.Error(err.Error())
			return err
		}
		util.Loggrs.Info("user add done  -> ", string(name))
		return nil
	}
	return nil
}
