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

type dictDataBase struct {
	StarRocks  string `json:"starrocks"`
	User       string `json:"user"`
	Database   string `json:"database"`
	Option     string `json:"option"`
	QuotaSize  string `json:"quota"`
	Policyname string `json:"policyname"`
}

func SetAddDataBase(c *gin.Context) {
	util.Loggrs.Info("数据库管理")

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

	var dictData dictDataBase
	err = json.Unmarshal(data, &dictData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": err.Error()})
		return
	}

	if dictData.StarRocks == "" || dictData.User == "" || dictData.Option == "" {
		c.AbortWithStatusJSON(http.StatusAccepted, gin.H{"status": "Fail", "message": "集群名称，用户名，指标缺失！"})
		return
	}

	//连接集群
	db, err := conn.StarRocks(dictData.StarRocks)
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

	var dbinfo []util.SetDataBase
	defer func() {
		result, err := json.MarshalIndent(dbinfo, "", "  ")
		if err != nil {
			util.Loggrs.Warn(err.Error())
			return
		}
		util.Loggrs.Info("Job result ->:\n" + string(result))
	}()

	var wg sync.WaitGroup
	for _, database := range strings.Split(dictData.Database, ",") {
		wg.Add(1)
		go func(database string) {
			defer wg.Done()

			exist := false
			var m []map[string]interface{}
			r := db.Raw(fmt.Sprintf("select SCHEMA_NAME from information_schema.schemata")).Scan(&m)
			if r.Error != nil {
				dbinfo = append(dbinfo, util.SetDataBase{
					App:        dictData.StarRocks,
					User:       dictData.User,
					Database:   database,
					Option:     dictData.Option,
					QuotaSize:  dictData.QuotaSize,
					Policyname: dictData.Policyname,
					State:      "Fail",
					Comment:    r.Error.Error(),
				})
			}
			for _, m2 := range m {
				if m2["SCHEMA_NAME"] == database {
					exist = true
				}
			}

			option := tools.DictDB{
				Db:       db,
				App:      dictData.StarRocks,
				User:     dictData.User,
				Database: database,
				Option:   dictData.Option,
				Policy:   dictData.Policyname,
			}

			//如库存在，不支持创建
			if exist && dictData.Option == "create" {
				dbinfo = append(dbinfo, util.SetDataBase{
					App:        dictData.StarRocks,
					User:       dictData.User,
					Database:   database,
					Option:     dictData.Option,
					QuotaSize:  dictData.QuotaSize,
					Policyname: dictData.Policyname,
					State:      "Ok",
					Comment:    "数据库存在，不允许创建！",
				})
				return
			}

			//如库不存在，支持创建
			if !exist && dictData.Option == "create" {
				dbinfo = append(dbinfo, tools.OnExecDatabase(option, fmt.Sprintf("CREATE DATABASE %s", database)))
				dbinfo = append(dbinfo, tools.OnExecDatabase(option, fmt.Sprintf("ALTER DATABASE %s SET DATA QUOTA 3T", database)))
			}

			//如库存在，支持drop
			if exist && dictData.Option == "drop" {
				dbinfo = append(dbinfo, util.SetDataBase{
					App:        dictData.StarRocks,
					User:       dictData.User,
					Database:   database,
					Option:     dictData.Option,
					QuotaSize:  dictData.QuotaSize,
					Policyname: dictData.Policyname,
					State:      "Fail",
					Comment:    "销毁库暂时被约束，无法执行销毁！",
				})
			}

			//如库不存在，不支持drop
			if !exist && dictData.Option == "drop" {
				dbinfo = append(dbinfo, util.SetDataBase{
					App:        dictData.StarRocks,
					User:       dictData.User,
					Database:   database,
					Option:     dictData.Option,
					QuotaSize:  dictData.QuotaSize,
					Policyname: dictData.Policyname,
					State:      "Fail",
					Comment:    "数据库不存在，无法执行销毁！",
				})
				return
			}

			//如库存在，支持调整库容量大小
			if exist && dictData.Option == "setquota" {
				dbinfo = append(dbinfo, tools.OnExecDatabase(option, fmt.Sprintf("ALTER DATABASE %s SET DATA QUOTA %s", database, dictData.QuotaSize)))
			}

			//如库不存在，不支持调整库容量大小
			if !exist && dictData.Option == "setquota" {
				dbinfo = append(dbinfo, util.SetDataBase{
					App:        dictData.StarRocks,
					User:       dictData.User,
					Database:   database,
					Option:     dictData.Option,
					QuotaSize:  dictData.QuotaSize,
					Policyname: dictData.Policyname,
					State:      "Fail",
					Comment:    "数据库不存在，无法执行重置库容量大小！",
				})
				return
			}

		}(database)
	}
	wg.Wait()

	defer func() {
		marshal, _ := json.Marshal(dbinfo)
		if strings.Contains(string(marshal), "Ok") {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Ok", "message": dbinfo})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": dbinfo})
	}()
	return
}
