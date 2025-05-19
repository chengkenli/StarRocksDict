package SetDictOneStop

import (
	"StarRocksDict/Send"
	"StarRocksDict/conn"
	"StarRocksDict/lark"
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

type dictUser struct {
	StarRocks  string `json:"starrocks"`
	User       string `json:"user"`
	Owner      string `json:"owner"`
	Ldap       bool   `json:"ldap"`
	Option     string `json:"option"`
	Policyname string `json:"policyname"`
}

func SetuAddUserID(c *gin.Context) {
	util.Loggrs.Info("账户管理")

	data, err := c.GetRawData()
	util.Loggrs.Info(string(data))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	var add dictUser
	err = json.Unmarshal(data, &add)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}

	if add.StarRocks == "" || add.User == "" {
		c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"status": "Fail", "message": "集群名称，或用户名为空。"})
		return
	}

	//连接集群
	db, err := conn.StarRocks(add.StarRocks)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
	/*每次使用完，主动关闭连接数*/
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			util.Loggrs.Error(err.Error())
			return
		}
		sqlDB.SetMaxOpenConns(30)                  //最大连接数
		sqlDB.SetMaxIdleConns(30)                  //最大空闲连接数
		sqlDB.SetConnMaxLifetime(30 * time.Second) //空闲连接最多存活时间
		sqlDB.Close()
	}()
	//匹配集群版本
	version := tools.CurrentVersion(db)
	util.Loggrs.Info(fmt.Sprintf("%0.2f", version))

	var userIds []util.SetUserId
	defer func() {
		result, err := json.MarshalIndent(userIds, "", "  ")
		if err != nil {
			util.Loggrs.Warn(err.Error())
			return
		}
		util.Loggrs.Info("Job result ->:\n" + string(result))
	}()

	var wg sync.WaitGroup
	for _, user := range strings.Split(add.User, ",") {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()

			exist := true
			if version < 2.5 {
				var m map[string]interface{}
				r := db.Raw("show grants for " + user).Scan(&m)
				if r.Error != nil {
					userIds = append(userIds, util.SetUserId{
						App:        add.StarRocks,
						User:       user,
						Owner:      add.Owner,
						Ldap:       add.Ldap,
						Option:     add.Option,
						Policyname: add.Policyname,
						State:      "Fail",
						Comment:    r.Error.Error(),
					})
					return
				}
				if m["AuthPlugin"] == nil && m["Password"] == nil {
					exist = false
				}
			}
			var m map[string]interface{}
			r := db.Raw("show grants for " + user).Scan(&m)
			if r.Error != nil {
				util.Loggrs.Warn(r.Error.Error())
				if strings.Contains(strings.ToLower(r.Error.Error()), "cannot find user") || strings.Contains(strings.ToLower(r.Error.Error()), "not exist") {
					exist = false
				} else {
					userIds = append(userIds, util.SetUserId{
						App:        add.StarRocks,
						User:       user,
						Owner:      add.Owner,
						Ldap:       add.Ldap,
						Option:     add.Option,
						Policyname: add.Policyname,
						State:      "Fail",
						Comment:    r.Error.Error(),
					})
					return
				}
			}

			dict := tools.Dicts{
				Db:     db,
				App:    add.StarRocks,
				User:   user,
				Ldap:   add.Ldap,
				Owner:  add.Owner,
				Option: add.Option,
				Policy: add.Policyname,
			}

			/*主体,密钥生成函数*/
			password := tools.RandomPassWord(32)
			//如果用户存在，不支持创建
			if exist && add.Option == "create" {
				userIds = append(userIds, util.SetUserId{
					App:        add.StarRocks,
					User:       user,
					Owner:      add.Owner,
					Ldap:       add.Ldap,
					Option:     add.Option,
					Policyname: add.Policyname,
					State:      "Ok",
					Comment:    "用户已经存在！",
				})
				return
			}

			//如果用户不存在，支持创建
			if !exist && add.Option == "create" {
				if add.Ldap {
					password = "AD PASSWORD(Office Computer Unlock Password)"
					userIds = append(userIds, tools.OnExecAdd(dict, fmt.Sprintf("CREATE USER %s IDENTIFIED WITH authentication_ldap_simple", user)))
				} else {
					userIds = append(userIds, tools.OnExecAdd(dict, fmt.Sprintf("CREATE USER %s", user)))
					userIds = append(userIds, tools.OnExecAdd(dict, fmt.Sprintf("ALTER USER '%s' IDENTIFIED BY '%s'", user, password)))
				}
				if strings.ReplaceAll(add.StarRocks, " ", "") == "sr-adhoc" {
					db, err := conn.StarRocks(add.StarRocks)
					if err != nil {
						util.Loggrs.Error(err.Error())
						return
					}
					/*每次使用完，主动关闭连接数*/
					defer func() {
						sqlDB, err := db.DB()
						if err != nil {
							util.Loggrs.Error(err.Error())
							return
						}
						sqlDB.SetMaxOpenConns(30)                  //最大连接数
						sqlDB.SetMaxIdleConns(30)                  //最大空闲连接数
						sqlDB.SetConnMaxLifetime(30 * time.Second) //空闲连接最多存活时间
						sqlDB.Close()
					}()
					go ResourceGroup(db, user, add.Ldap)
				}
			}

			//如果用户存在，不支持重置密码
			if !exist && add.Option == "reset" {
				userIds = append(userIds, util.SetUserId{
					App:        add.StarRocks,
					User:       user,
					Owner:      add.Owner,
					Ldap:       add.Ldap,
					Option:     add.Option,
					Policyname: add.Policyname,
					State:      "Fail",
					Comment:    "用户不存在，不支持重置密码！",
				})
				return
			}

			//如果用户存在，可以支持重置密码
			if exist && add.Option == "reset" {
				if add.Ldap {
					userIds = append(userIds, util.SetUserId{
						App:        add.StarRocks,
						User:       user,
						Owner:      add.Owner,
						Ldap:       add.Ldap,
						Option:     add.Option,
						Policyname: add.Policyname,
						State:      "Fail",
						Comment:    "LDAP 账号无法重置密码！",
					})
					return
				} else {
					userIds = append(userIds, tools.OnExecAdd(dict, fmt.Sprintf("ALTER USER '%s' IDENTIFIED BY '%s'", user, password)))
				}
			}

			//如果用户不存在，不支持销毁账号
			if !exist && add.Option == "drop" {
				userIds = append(userIds, util.SetUserId{
					App:        add.StarRocks,
					User:       user,
					Owner:      add.Owner,
					Ldap:       add.Ldap,
					Option:     add.Option,
					Policyname: add.Policyname,
					State:      "Fail",
					Comment:    "用户不存在，不支持销毁！",
				})
				return
			}

			//如果用户存在，可以支持销毁账号
			if exist && add.Option == "drop" {
				userIds = append(userIds, tools.OnExecAdd(dict, fmt.Sprintf("DROP USER '%s'", user)))
			}

			go func() {

				web, jdbc := tools.JdbcHost(add.StarRocks)
				option := util.DictOption{
					StarRocks: add.StarRocks,
					User:      user,
					Owner:     add.Owner,
					Password:  password,
					Option:    add.Option,
					Web:       web,
					Other:     jdbc,
					Policy:    add.Policyname,
				}

				to, msg := Send.HtmlEmailBody(&option)
				util.Loggrs.Info(fmt.Sprintf("#发送邮件 %s", to))
				Send.SendEmail2(
					&util.Emailinfo{
						Subject: fmt.Sprintf("starrocks %s notification", add.StarRocks),
						To:      strings.Join(to, ","),
						From:    util.MetaConf["slow_query_email_from"].(string),
						Cc:      "",
						Bc:      util.MetaConf["slow_query_email_bc"].(string),
						Attach:  "",
						Emsg:    msg,
					})
				util.Loggrs.Info("send success.")
			}()

		}(user)
	}
	wg.Wait()

	defer func() {
		marshal, _ := json.Marshal(userIds)
		if strings.Contains(string(marshal), "Ok") {
			go func() {
				u := userIds[0]
				var username string
				meta := lark.Lark2userid(u.Owner)
				if meta != nil {
					username = meta["user_name"].(string)
				} else {
					username = u.Owner
				}
				var option string
				switch u.Option {
				case "create":
					option = "登录时账号严格区分大小写，请您使用小写进行登录"
				case "drop":
					option = "账号已经被删除"
				case "reset":
					option = "密码已经重置，请关注StarRocks邮件获取新密码"
				}
				var other string
				if u.Ldap {
					other = "登录密码是您的AD密码！"
				} else {
					other = "登录密码存放在您的邮箱中，请关注StarRocks邮件！"
				}
				var message string
				if u.Policyname == "" {
					message = fmt.Sprintf(`
集群：%s
账号：(%s)%s
%s
`, u.App, u.State, u.User, strings.Join(util.Config.GetStringSlice("Schema.Connects"), "\n"))
				} else {
					message = fmt.Sprintf(`
单号：%s
集群：%s
账号：(%s)%s
%s
`, u.Policyname, u.State, u.App, u.User, strings.Join(util.Config.GetStringSlice("Schema.Connects"), "\n"))
				}

				messages := fmt.Sprintf("Hello！%s 😊\n您申请的StarRocks的账号流程已经结束。\n%s，%s\n如有其他疑问或需要帮助可加入[StarRocks总群](%s)提问，感谢您使用StarRocks！\n%s",
					username, option, other, util.Config.GetString("Schema.GroupUri"), message)

				util.Loggrs.Info(messages)
				err := lark.Send2userid(u.Owner, messages)
				if err != nil {
					util.Loggrs.Warn(err.Error())
					return
				}
				if util.Config.GetStringSlice("Schema.lkremind") != nil {
					go lark.Send2Markdown(fmt.Sprintf("%s向%s发送了一条信息", time.Now().Format("2006-01-02 15:04:05"), add.Owner), messages, util.Config.GetStringSlice("Schema.lkremind"))
				}
			}()
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok", "message": userIds})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": userIds})
	}()
	return
}
