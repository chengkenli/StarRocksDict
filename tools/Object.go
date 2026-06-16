package tools

import (
	"StarRocksDict/util"
	"crypto/rand"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

type SrAvgs struct {
	Host string
	Port int
	User string
	Pass string
	Base string
}
type Schem struct {
	Db         *gorm.DB
	Sql        string
	Catalog    string
	User       string
	Table      string
	Permission string
	App        string
	Policy     string
}
type SchemRole struct {
	Db         *gorm.DB
	Sql        string
	Catalog    string
	Role       string
	Table      string
	Permission string
	App        string
	Policy     string
}

type Dicts struct {
	Db     *gorm.DB
	App    string
	User   string
	Owner  string
	Ldap   bool
	Option string
	Policy string
}

type DictDB struct {
	Db       *gorm.DB
	App      string
	User     string
	Database string
	Option   string
	Policy   string
}

// CurrentVersion 获取集群版本
func CurrentVersion(db *gorm.DB) float64 {
	/*匹配starrocks版本*/
	var v map[string]interface{}
	db.Raw("select current_version() as version").Scan(&v)
	if v["version"].(string) == "branch-3.1-base-3.1.10-ee-7ec8955" {
		return 3.1
	}
	vs, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", strings.Split(strings.Split(v["version"].(string), " ")[0], ".")[0], strings.Split(strings.Split(v["version"].(string), " ")[0], ".")[1]), 64)
	return vs
}

// 判断某个视图是普通视图还是物化视图
func Materialized(db *gorm.DB, table string) (string, bool) {
	var m map[string]interface{}
	r := db.Raw("show create table " + table).Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return "", false
	}
	var materialized bool
	var str string
	if m["Create Materialized View"] != nil {
		str = "MATERIALIZED VIEW"
		materialized = true
	}
	if m["Create View"] != nil {
		str = "VIEW"
		materialized = false
	}
	return str, materialized
}

// OnExec20 执行授权ddl
func OnExec20(s *Schem, sql string, comment ...string) util.SetResult {
	var dict string
	if strings.Contains(strings.ToLower(sql), "grant") {
		dict = "授予"
	}
	if strings.Contains(strings.ToLower(sql), "revoke") {
		dict = "回收"
	}

	stmt := fmt.Sprintf(`set @policyName="%s";%s`, s.Policy, sql)
	var err error
	for i := 0; i < 3; i++ {
		r := s.Db.Exec(stmt)
		if r.Error != nil {
			err = r.Error
			util.Loggrs.Info("Task Failed ->: " + stmt)
			if strings.Contains(r.Error.Error(), "map[]") {
				util.Loggrs.Error(fmt.Sprintf("#%d %v    ->  ", i, s) + stmt + " -> " + r.Error.Error())
			}
			continue
		}
		util.Loggrs.Info("Task Success ->: " + stmt)
		break
	}
	if err != nil {
		return util.SetResult{
			App:     s.App,
			Catalog: s.Catalog,
			User:    s.User,
			Table:   s.Table,
			Permit:  s.Permission,
			State:   "Fail " + strings.Join(comment, ","),
			Demand:  dict,
			Comment: err.Error(),
		}
	}

	return util.SetResult{
		App:     s.App,
		Catalog: s.Catalog,
		User:    s.User,
		Table:   s.Table,
		Permit:  s.Permission,
		State:   "Ok " + strings.Join(comment, ","),
		Demand:  dict,
		Comment: "",
	}
}

// Onexec 执行授权ddl
func Onexec(mutex *sync.Mutex, s *Schem, sql string, comment ...string) util.SetResult {
	mutex.Lock()
	defer mutex.Unlock()
	var dict string
	if strings.Contains(strings.ToLower(sql), "grant") {
		dict = "授予"
	}
	if strings.Contains(strings.ToLower(sql), "revoke") {
		dict = "回收"
	}

	stmt := fmt.Sprintf(`set @policyName="%s";%s`, s.Policy, sql)
	var err error
	for i := 0; i < 3; i++ {
		util.Loggrs.Info("Task Execute ->: " + stmt)
		r := s.Db.Exec(stmt)
		if r.Error != nil {
			err = r.Error
			util.Loggrs.Info("Task Failed ->: " + stmt)
			if strings.Contains(r.Error.Error(), "map[]") {
				util.Loggrs.Error(fmt.Sprintf("#%d %v    ->  ", i, s) + stmt + " -> " + r.Error.Error())
			}
			continue
		}
		util.Loggrs.Info("Task Success ->: " + stmt)
		break
	}
	if err != nil {
		return util.SetResult{
			App:     s.App,
			Catalog: s.Catalog,
			User:    s.User,
			Table:   s.Table,
			Permit:  s.Permission,
			State:   "Fail " + strings.Join(comment, ","),
			Demand:  dict,
			Comment: err.Error(),
		}
	}
	return util.SetResult{
		App:     s.App,
		Catalog: s.Catalog,
		User:    s.User,
		Table:   s.Table,
		Permit:  s.Permission,
		State:   "Ok " + strings.Join(comment, ","),
		Demand:  dict,
		Comment: "",
	}
}

// OnExecRole 执行授权ddl
func OnExecRole(s SchemRole, sql string, comment ...string) util.SetRole {
	stmt := fmt.Sprintf(`set @policyName="%s";%s`, s.Policy, sql)
	fmt.Println(stmt)
	r := s.Db.Exec(stmt)
	if r.Error != nil {
		return util.SetRole{
			App:        s.App,
			Role:       s.Role,
			Table:      s.Table,
			Permission: s.Permission,
			Revoke:     false,
			Catalog:    s.Catalog,
			Policyname: s.Policy,
			State:      "Fail " + strings.Join(comment, ","),
			Comment:    r.Error.Error(),
		}
	}

	return util.SetRole{
		App:        s.App,
		Role:       s.Role,
		Table:      s.Table,
		Permission: s.Permission,
		Revoke:     false,
		Catalog:    s.Catalog,
		Policyname: s.Policy,
		State:      "Ok " + strings.Join(comment, ","),
		Comment:    "",
	}
}

// OnExecAdd 执行账号管理ddl
func OnExecAdd(s Dicts, sql string, comment ...string) util.SetUserId {
	/*匹配授权，回收标识*/
	stmt := fmt.Sprintf(`set @policyName="%s";%s`, s.Policy, sql)
	util.Loggrs.Info(stmt)

	r := s.Db.Exec(stmt)
	if r.Error != nil {
		return util.SetUserId{
			App:        s.App,
			User:       s.User,
			Owner:      s.Owner,
			Ldap:       s.Ldap,
			Option:     s.Option,
			Policyname: s.Policy,
			State:      "Fail " + strings.Join(comment, ","),
			Comment:    r.Error.Error(),
		}
	}
	return util.SetUserId{
		App:        s.App,
		User:       s.User,
		Owner:      s.Owner,
		Ldap:       s.Ldap,
		Option:     s.Option,
		Policyname: s.Policy,
		State:      "Ok " + strings.Join(comment, ","),
		Comment:    sql,
	}
}

// OnExecDatabase 执行管理数据库
func OnExecDatabase(s DictDB, sql string, comment ...string) util.SetDataBase {
	stmt := fmt.Sprintf(`set @policyName="%s";%s`, s.Policy, sql)
	util.Loggrs.Info(stmt)

	r := s.Db.Exec(stmt)
	if r.Error != nil {
		return util.SetDataBase{
			App:        s.App,
			User:       s.User,
			Database:   s.Database,
			Option:     s.Option,
			QuotaSize:  "",
			Policyname: s.Policy,
			State:      "Fail " + strings.Join(comment, ","),
			Comment:    r.Error.Error(),
		}
	}

	return util.SetDataBase{
		App:        s.App,
		User:       s.User,
		Database:   s.Database,
		Option:     s.Option,
		QuotaSize:  "",
		Policyname: s.Policy,
		State:      "Ok " + strings.Join(comment, ","),
		Comment:    sql,
	}
}

// RandomPassWord 根据指定长度的生产高敏感度字符串
func RandomPassWord(n int) string {
	allowedChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()-_=+"

	b := make([]byte, n)
	for i := range b {
		// 生成一个随机索引
		ri, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
		if err != nil {
			return ""
		}
		// 使用随机索引获取一个字符
		b[i] = allowedChars[ri.Int64()]
	}

	return string(b)
}

// StrInSlice 检查数组中是否存在某个元素
func StrInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// RemoveInSlice 从数组中移除某个元素
func RemoveInSlice(slice []string, element string) []string {
	// 创建一个新的切片来保存结果
	result := make([]string, 0)
	// 遍历原切片，并将不等于element的元素添加到结果切片中
	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}

// RmSliceStr /*数组去重*/
func RmSliceStr(strs []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复字符串
	for _, e := range strs {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// FindDiffSlice
// 找到与元素a不同的其他元素
func FindDiffSlice(slice []string, a string) []string {
	// 创建一个新的切片来存储不同的元素
	different := make([]string, 0)

	// 遍历原始切片
	for _, element := range slice {
		// 如果当前元素不等于a，则添加到结果切片中
		if element != a {
			different = append(different, element)
		}
	}
	return different
}

// FindDiffSlices 返回两个字符串切片的不同元素。
func FindDiffSlices(slice1, slice2 []string) []string {
	diffMap := make(map[string]int)
	var differences []string

	// 增加slice1中元素的计数
	for _, str := range slice1 {
		diffMap[str]++
	}
	// 增加slice2中元素的计数
	for _, str := range slice2 {
		diffMap[str]++
	}
	// 遍历map来找到只出现一次的元素
	for str, count := range diffMap {
		if count == 1 {
			differences = append(differences, str)
		}
	}
	return differences
}

func GenRandomStr(n int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

// AccessControl
// 检查当前集群，属于native鉴权，还是ranger鉴权
func AccessControl(db *gorm.DB) (string, error) {
	var ac map[string]interface{}
	r := db.Raw("admin show frontend config like 'access_control'").Scan(&ac)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return "native", r.Error
	}
	if ac == nil {
		return "native", nil
	}
	return ac["Value"].(string), nil
}
