package SetRangeOneStop

import (
	"StarRocksDict/util"
	"io"
	"io/ioutil"
	"net/http"
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

	client := &http.Client{
		Timeout:   time.Second * 30,
		Transport: &http.Transport{},
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
