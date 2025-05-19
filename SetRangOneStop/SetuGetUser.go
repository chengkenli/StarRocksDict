package SetRangeOneStop

import (
	"StarRocksDict/util"
)

// getByName
// 获取用户信息
func getByName(user string) ([]byte, error) {
	r, err := Post(
		&util.RangerHost{
			Server: util.Config.GetString("Ranger.host"),
			Access: util.Config.GetString("Ranger.user"),
			Secret: util.Config.GetString("Ranger.password"),
		}, "GET", "/service/xusers/users/userName/"+user, nil)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return r, err
	}
	return r, nil
}
