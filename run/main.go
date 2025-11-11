package run

import (
	"StarRocksDict/SetDictOneStop"
	SetRangeOneStop "StarRocksDict/SetRangOneStop"
	"StarRocksDict/meta"
	"StarRocksDict/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	_ "net/http/pprof"
)

func Run() {
	go meta.Information()

	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(util.KeyRequestId, xid.New().String())
		c.Next()
	})

	//go func() {
	//	// 使用默认的http.ServeMux，pprof路由已经注册
	//	if err := http.ListenAndServe(":6060", nil); err != nil {
	//		fmt.Printf("Failed to start server: %v\n", err)
	//	}
	//}()

	r.Use(checkSign)
	r.POST("/setui", SetDictOneStop.SetuDict)
	r.POST("/setri", SetDictOneStop.SetuDictRole)
	r.POST("/reuser", SetDictOneStop.SetuAddUserID)
	r.POST("/addbase", SetDictOneStop.SetAddDataBase)
	r.POST("/ranger/setui", SetRangeOneStop.ApiaddPolicy)

	port := fmt.Sprintf(":%d", util.Read.StarRocks.ServicePort)
	err := r.Run(port)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
}

// pprofHandler 将请求代理到pprof的默认处理程序
func pprofHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/debug/pprof/":
		http.Redirect(w, r, "/debug/pprof/", http.StatusTemporaryRedirect)
	default:
		http.DefaultServeMux.ServeHTTP(w, r)
	}
}
