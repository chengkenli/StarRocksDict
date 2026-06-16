package SetDictOneStop

import (
	"StarRocksDict/conn"
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

type dictJsonRole struct {
	Starrocks  string `bson:"starrocks"`
	Role       string `bson:"role"`
	Table      string `bson:"table"`
	Permission string `bson:"permission"`
	Revoke     bool   `bson:"revoke"`
	Catalog    string `bson:"catalog"`
	Policyname string `bson:"policyname"`
}

func SetuDictRole(c *gin.Context) {
	util.Loggrs.Info("角色授权机制")

	data, err := c.GetRawData()
	util.Loggrs.Info(string(data))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}
	if len(data) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": "结构体异常"})
		return
	}
	var e dictJsonRole
	err = json.Unmarshal(data, &e)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}

	if e.Starrocks == "" || e.Role == "" || e.Permission == "" || e.Table == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": "参数缺失：集群名称，用户名，权限，表名"})
		return
	}

	/*匹配授权，回收标识*/
	var dictH, dictE string
	switch e.Revoke {
	case true:
		dictH = "REVOKE"
		dictE = "FROM ROLE"
	case false:
		dictH = "GRANT"
		dictE = "TO ROLE"
	default:
		dictH = "GRANT"
		dictE = "TO ROLE"
	}

	if e.Catalog == "" {
		e.Catalog = "default_catalog"
	}

	//连接集群
	db, err := conn.StarRocks(e.Starrocks)
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

	var wg sync.WaitGroup
	//执行授权逻辑
	var roleInfo []util.SetRole
	defer func() {
		result, err := json.MarshalIndent(roleInfo, "", "  ")
		if err != nil {
			util.Loggrs.Warn(err.Error())
			return
		}
		util.Loggrs.Info("Job result ->:\n" + string(result))
	}()

	for _, table := range strings.Split(strings.ReplaceAll(e.Table, " ", ""), ",") {
		wg.Add(1)
		go func(table string) {
			defer wg.Done()
			for _, role := range strings.Split(e.Role, ",") {
				/*****************************************/
				//判断账号状态是否异常
				if version < 2.5 {
					var m map[string]interface{}
					r := db.Raw("show grants for role " + role).Scan(&m)
					if r.Error != nil {
						roleInfo = append(roleInfo, util.SetRole{
							App:        e.Starrocks,
							Role:       role,
							Table:      table,
							Permission: e.Permission,
							Revoke:     e.Revoke,
							Catalog:    e.Catalog,
							Policyname: e.Policyname,
							State:      "Fail",
							Comment:    r.Error.Error(),
						})
						continue
					}
					if m["AuthPlugin"] == nil && m["Password"] == nil {
						roleInfo = append(roleInfo, util.SetRole{
							App:        e.Starrocks,
							Role:       role,
							Table:      table,
							Permission: e.Permission,
							Revoke:     e.Revoke,
							Catalog:    e.Catalog,
							Policyname: e.Policyname,
							State:      "Fail",
							Comment:    fmt.Sprintf("%s 角色异常，或不存在", role),
						})
						continue
					}
				}
				var m map[string]interface{}
				r := db.Raw("show grants for role " + role).Scan(&m)
				if r.Error != nil {
					roleInfo = append(roleInfo, util.SetRole{
						App:        e.Starrocks,
						Role:       role,
						Table:      table,
						Permission: e.Permission,
						Revoke:     e.Revoke,
						Catalog:    e.Catalog,
						Policyname: e.Policyname,
						State:      "Fail",
						Comment:    r.Error.Error(),
					})
					continue
				}
				/*****************************************/
				//开始执行权限步骤
				edict := dictStructRole{
					dictH: dictH,
					dictE: dictE,
					db:    db,
					e: dictJsonRole{
						Starrocks:  e.Starrocks,
						Permission: e.Permission,
						Revoke:     e.Revoke,
						Catalog:    e.Catalog,
					},
					table:  table,
					role:   role,
					policy: e.Policyname,
				}

				if version >= 3.0 {
					roleInfo = append(roleInfo, SetuAddDictRole3(edict)...)
				}
			}
		}(table)
	}
	wg.Wait()

	defer func() {
		marshal, _ := json.Marshal(roleInfo)
		if strings.Contains(string(marshal), "Ok") {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok", "message": roleInfo})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": roleInfo})
	}()
	return
}
