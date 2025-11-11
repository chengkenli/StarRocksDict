package SetRangeOneStop

import (
	"StarRocksDict/util"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func Post(a *util.RangerHost, method, u string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, a.Server+u, body)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	request.Header.Set("Accept", "application/json;charset=utf-8")
	request.SetBasicAuth(a.Access, a.Secret)

	//飞书代理地址
	larkproxy := "http://xxx"

	proxy, _ := url.Parse(larkproxy)
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	respone, err := client.Do(request)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil, err
	}
	defer respone.Body.Close()
	b, err := ioutil.ReadAll(respone.Body)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil, err
	}
	return b, nil
}
