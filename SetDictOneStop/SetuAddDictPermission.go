package SetDictOneStop

import (
	"StarRocksDict/conn"
	"StarRocksDict/lark"
	"StarRocksDict/tools"
	"StarRocksDict/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

type dictJson struct {
	Starrocks  string `bson:"starrocks"`
	User       string `bson:"user"`
	Submit     string `bson:"submit"`
	Table      string `bson:"table"`
	Permission string `bson:"permission"`
	Revoke     bool   `bson:"revoke"`
	Catalog    string `bson:"catalog"`
	Policyname string `bson:"policyname"`
}

func SetuDict(c *gin.Context) {
	util.Loggrs.Info("用户授权机制")

	data, err := c.GetRawData()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	if len(data) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": "结构体异常"})
		return
	}
	var e dictJson
	err = json.Unmarshal(data, &e)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	if e.Starrocks == "" || e.User == "" || e.Permission == "" || e.Table == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": "参数缺失：集群名称，用户名，权限，表名"})
		return
	}

	util.Loggrs.Info("Job start ->: " + string(data))

	/*匹配授权，回收标识*/
	var dictH, dictE, dict string
	switch e.Revoke {
	case true:
		dictH = "REVOKE"
		dictE = "FROM"
		dict = "回收"
	case false:
		dictH = "GRANT"
		dictE = "TO"
		dict = "授予"
	default:
		dictH = "GRANT"
		dictE = "TO"
		dict = "授予"
	}

	if e.Catalog == "" {
		e.Catalog = "default_catalog"
	}
	util.Loggrs.Info("Job start connect sr.")
	//连接集群
	var tableResult []util.SetResult
	defer func() {
		result, err := json.MarshalIndent(tableResult, "", "  ")
		if err != nil {
			util.Loggrs.Warn(err.Error())
			return
		}
		util.Loggrs.Info("Job result ->:\n" + string(result))

		var tablelist []string
		for i, item := range tableResult {
			if item.Table == "ALL CATALOGS" {
				continue
			}
			tablelist = append(tablelist, fmt.Sprintf("%02d. (%s)(%s) %s", i, item.State, item.Permit, item.Table))
		}
		var message string
		if e.Policyname == "" {
			message = fmt.Sprintf(`
集群：%s
目录：%s
账号：%s
库表：
%s`, e.Starrocks, e.Catalog, e.User, strings.Join(tablelist, "\n"))
		} else {
			message = fmt.Sprintf(`
单号：%s
集群：%s
目录：%s
账号：%s
库表：
%s`, e.Policyname, e.Starrocks, e.Catalog, e.User, strings.Join(tablelist, "\n"))
		}

		//给userid发送一个飞书信息
		go func() {
			//if e.Starrocks != "sr-adhoc" {
			//	return
			//}
			var submit_user string
			if e.Submit != "" {
				submit_user = e.Submit
			} else {
				submit_user = e.User
			}

			var username string
			meta := lark.Lark2userid(submit_user)
			if meta != nil {
				username = meta["user_name"].(string)
			} else {
				username = submit_user
			}
			message := fmt.Sprintf("Hello！%s 😊\n您申请的StarRocks的库表权限流程已经结束。\n如有其他疑问或需要帮助可加入[StarRocks总群](%s)提问，感谢您使用StarRocks！\n%s", username, util.Config.GetString("Schema.GroupUri"), message)
			util.Loggrs.Info(message)
			err := lark.Send2userid(submit_user, message)
			if err != nil {
				util.Loggrs.Warn(err.Error())
				return
			}
			if util.Config.GetStringSlice("Schema.lkremind") != nil {
				go lark.Send2Markdown(fmt.Sprintf("%s向%s发送了一条信息", time.Now().Format("2006-01-02 15:04:05"), submit_user), message, util.Config.GetStringSlice("Schema.lkremind"))
			}
		}()
		marshal, _ := json.Marshal(tableResult)
		if strings.Contains(string(marshal), "Ok") {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok", "message": tableResult})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": tableResult})
	}()

	// 重构超时重试函数(当fun(){}超过5秒没完成，断开重试。)
	retry := tools.ObjectRetry(
		tools.WithInterval(1*time.Second),
		tools.WithAttempts(3),
		tools.WithTimeout(10*time.Second),
	)
	err = retry.Do(func(ctx context.Context) error {
		done := make(chan struct{})
		// 这里放置你的长时间运行函数
		go func() {
			tableResult = seting(e, dictH, dictE, dict)
			util.Loggrs.Info("all tasks done.")
			done <- struct{}{}
		}()
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if err != nil {
		tableResult = append(tableResult, util.SetResult{
			App:     e.Starrocks,
			Catalog: e.Catalog,
			User:    e.User,
			Table:   e.Table,
			Permit:  e.Permission,
			State:   "Fail",
			Demand:  "",
			Comment: err.Error(),
		})
		util.Loggrs.Error(err.Error())
		return
	}
	// end
	if tableResult == nil {
		tableResult = append(tableResult, util.SetResult{
			App:     e.Starrocks,
			Catalog: e.Catalog,
			User:    e.User,
			Table:   e.Table,
			Permit:  e.Permission,
			State:   "Fail",
			Demand:  "",
			Comment: "result is nil",
		})
		return
	}
	util.Loggrs.Info("Job done ->: " + string(data))
	return
}

// 实现主体逻辑
func seting(e dictJson, dictH, dictE, dict string) []util.SetResult {
	db, err := conn.StarRocks(e.Starrocks)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil
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
	util.Loggrs.Info("Job connect sr done.")
	//匹配集群版本
	version := tools.CurrentVersion(db)
	util.Loggrs.Info(fmt.Sprintf("%0.2f", version))
	util.Loggrs.Info("Job get sr version.")

	//util.Loggrs.Info("sleep 10 second...")
	//time.Sleep(time.Second * 10)

	var once sync.Once
	var wg sync.WaitGroup
	//执行授权逻辑
	var tableResult []util.SetResult
	ch := make(chan struct{}, 3)
	for _, table := range strings.Split(strings.ReplaceAll(e.Table, " ", ""), ",") {
		wg.Add(1)
		go func(table string) {
			defer func() {
				<-ch
				wg.Done()
			}()
			ch <- struct{}{}

			for _, user := range strings.Split(e.User, ",") {
				/*****************************************/
				//判断账号状态是否异常
				if version < 2.5 {
					var m map[string]interface{}
					r := db.Raw("show grants for " + user).Scan(&m)
					if r.Error != nil {
						tableResult = append(tableResult, util.SetResult{
							App:     e.Starrocks,
							Catalog: e.Catalog,
							User:    user,
							Table:   table,
							Permit:  e.Permission,
							State:   "Fail",
							Demand:  dict,
							Comment: r.Error.Error(),
						})
						continue
					}
					if m["AuthPlugin"] == nil && m["Password"] == nil {
						tableResult = append(tableResult, util.SetResult{
							App:     e.Starrocks,
							Catalog: e.Catalog,
							User:    user,
							Table:   table,
							Permit:  e.Permission,
							State:   "Fail",
							Demand:  dict,
							Comment: fmt.Sprintf("%s 账号异常，或不存在", user),
						})
						continue
					}
				}
				var m map[string]interface{}
				r := db.Raw("show grants for " + user).Scan(&m)
				if r.Error != nil {
					tableResult = append(tableResult, util.SetResult{
						App:     e.Starrocks,
						Catalog: e.Catalog,
						User:    user,
						Table:   table,
						Permit:  e.Permission,
						State:   "Fail",
						Demand:  dict,
						Comment: r.Error.Error(),
					})
					continue
				}
				/*****************************************/
				//开始执行权限步骤
				edict := dictStruct{
					dictH: dictH,
					dictE: dictE,
					db:    db,
					e: dictJson{
						Starrocks:  e.Starrocks,
						Permission: e.Permission,
						Revoke:     e.Revoke,
						Catalog:    e.Catalog,
					},
					table:  table,
					user:   user,
					policy: e.Policyname,
				}

				switch version {
				case 3.0, 3.1, 3.2, 3.3, 3.4:
					tableResult = append(tableResult, SetuAddDict3(&edict, &once)...)

				case 2.1, 2.2, 2.5:
					tableResult = append(tableResult, SetuAddDict2(edict)...)

				}
			}
		}(table)
	}
	wg.Wait()

	return tableResult
}
