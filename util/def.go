package util

import (
	"github.com/sirupsen/logrus"
)

var (
	Loggrs   *logrus.Logger
	MetaLink []map[string]interface{}
	MetaConf map[string]interface{}
	Read     ReadInConf
)

const KeyRequestId = "requestId"
const KeySingle = "api_dict"

type SrAvgs struct {
	Host string
	Port int
	User string
	Pass string
}

type ReadInConf struct {
	Metadb struct {
		Host       string `json:"host"`
		Port       int    `json:"port"`
		User       string `json:"user"`
		Password   string `json:"password"`
		Base       string `json:"base"`
		Slowconfig string `json:"slowconfig"`
	} `json:"metadb"`
	StarRocks struct {
		ServicePort       int    `json:"service_port"`
		ServiceEnckey     string `json:"service_enckey"`
		ServiceXStarRocks string `json:"service_X-StarRocks"`
	} `json:"StarRocks"`
	Ranger struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"Ranger"`
	Schema struct {
		Lkremind []string `json:"lkremind"`
		LarkMeta string   `json:"larkMeta"`
		GroupURI string   `json:"GroupUri"`
		Connects []string `json:"Connects"`
	} `json:"Schema"`
	Log struct {
		Line    bool   `json:"line"`
		Level   int    `json:"level"`
		Path    string `json:"path"`
		Console bool   `json:"console"`
	} `json:"log"`
}

type Grants struct {
	UserIdentity      string `bson:"UserIdentity"`
	Password          string `bson:"Password"`
	AuthPlugin        string `bson:"AuthPlugin"`
	UserForAuthPlugin string `bson:"UserForAuthPlugin"`
	GlobalPrivs       string `bson:"GlobalPrivs"`
	DatabasePrivs     string `bson:"DatabasePrivs"`
	TablePrivs        string `bson:"TablePrivs"`
	ResourcePrivs     string `bson:"ResourcePrivs"`
}
type Grants2 struct {
	UserIdentity string `bson:"UserIdentity"`
	Grants       string `bson:"Grants"`
}
type Grants3 []struct {
	UserIdentity string `bson:"UserIdentity"`
	Catalog      string `bson:"Catalog"`
	Grants       string `bson:"Grants"`
}
type Auth struct {
	UserIdentity      string `bson:"UserIdentity"`
	Password          string `bson:"Password"`
	AuthPlugin        string `bson:"AuthPlugin"`
	UserForAuthPlugin string `bson:"UserForAuthPlugin"`
}

type Emailinfo struct {
	Subject string
	To      string
	From    string
	Cc      string
	Bc      string
	Attach  string
	Emsg    string
}
type EmailBody struct {
	Username string
	Password string
	App      string
	Web      string
	Other    string
}

type DictOption struct {
	StarRocks string `json:"starrocks"`
	User      string `json:"user"`
	Owner     string `json:"owner"`
	Password  string `json:"password"`
	Option    string `json:"option"`
	Web       string
	Other     string
	Policy    string
}

type RangerHost struct {
	Server string
	Access string
	Secret string
}

type SetPolicy struct {
	App        string `json:"starrocks"`
	Catalog    string `json:"catalog"`
	User       string `json:"user"`
	Table      string `json:"table"`
	Permission string `json:"permission"`
	Revoke     bool   `json:"revoke"`
	Comment    string `json:"comment"`
}

type SetResult struct {
	App     string
	Catalog string
	User    string
	Table   string
	Permit  string
	State   string
	Demand  string
	Comment string
}
type SetUserId struct {
	App        string
	User       string
	Owner      string
	Ldap       bool
	Option     string
	Policyname string
	State      string
	Comment    string
}
type SetDataBase struct {
	App        string
	User       string
	Database   string
	Option     string
	QuotaSize  string
	Policyname string
	State      string
	Comment    string
}
type SetRole struct {
	App        string
	Role       string
	Table      string
	Permission string
	Revoke     bool
	Catalog    string
	Policyname string
	State      string
	Comment    string
}
