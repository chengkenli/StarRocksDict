package util

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Printf("\nUsage: %s [-c conf] [-h]\nStarRocks 权限API 2.0\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Println()
}

func init() {
	var (
		conf        string
		help        bool
		defaultConf string
	)

	execDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	defaultConf = fmt.Sprintf("%s/.%s.yaml", execDir, filepath.Base(os.Args[0]))

	flag.StringVar(&conf, "c", defaultConf, "conf file")
	flag.BoolVar(&help, "h", false, "show help information")

	flag.Parse()
	flag.Usage = usage

	if help {
		flag.Usage()
		os.Exit(-1)
	}

	paths, name := filepath.Split(conf)
	Config := viper.New()
	Config.SetConfigFile(fmt.Sprintf("%s%s", paths, name))
	if err := Config.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
	}
	// 读取配置文件所有内容放入内存
	marshal, err := json.Marshal(Config.AllSettings())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	err = json.Unmarshal(marshal, &Read)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	// log init
	Logrus()
	// auth metadb
	Instglobal.vaild <- struct{}{}
	// send begin channel
	Instglobal.Begin <- struct{}{}
}

func Parms() {

}
