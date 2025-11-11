package SetRangeOneStop

import (
	"StarRocksDict/util"
)

// getByName
// 获取用户信息
func getByName(user string) ([]byte, error) {
	r, err := Post(
		&util.RangerHost{
			Server: util.Read.Ranger.Host,
			Access: util.Read.Ranger.User,
			Secret: util.Read.Ranger.Password,
		}, "GET", "/service/xusers/users/userName/"+user, nil)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return r, err
	}
	return r, nil
}
