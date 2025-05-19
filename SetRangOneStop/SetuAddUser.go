package SetRangeOneStop

import (
	"StarRocksDict/util"
	"bytes"
	"github.com/goccy/go-json"
	"github.com/rs/xid"
)

type RangerUser struct {
	App  string `json:"app"`
	User string `json:"user"`
}

// addByName
// 创建ranger用户
func addByName(user string) ([]byte, error) {
	type User struct {
		//Id           int      `json:"id"`
		Name string `json:"name"`
		//CreateDate   string   `json:"createDate"`
		//UpdateDate   string   `json:"updateDate"`
		//Owner        string   `json:"owner"`
		//UpdateBy     string   `json:"updateBy"`
		FirstName string `json:"firstName"`
		//LastName     string   `json:"lastName"`
		//EmailAddress string   `json:"emailAddress"`
		Password     string   `json:"password"`
		Description  string   `json:"description"`
		Status       int      `json:"status"`
		IsVisible    int      `json:"isVisible"`
		UserSource   int      `json:"userSource"`
		UserRoleList []string `json:"userRoleList"`
	}
	u := User{
		Name:         user,
		FirstName:    user + "@starrocks.com",
		Password:     xid.New().String() + "@StarRocks0",
		Description:  user + "@starrocks.com",
		Status:       1,
		IsVisible:    1,
		UserSource:   0,
		UserRoleList: []string{"ROLE_USER"},
	}
	marshal, err := json.Marshal(&u)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil, err
	}
	r, err := Post(
		&util.RangerHost{
			Server: util.Config.GetString("Ranger.host"),
			Access: util.Config.GetString("Ranger.user"),
			Secret: util.Config.GetString("Ranger.password"),
		}, "POST", "/service/xusers/secure/users", bytes.NewReader(marshal))
	if err != nil {
		util.Loggrs.Error(err.Error())
		return r, err
	}
	return r, nil
}
