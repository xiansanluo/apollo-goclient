package agollo

import (
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"time"
)

type CallBack struct {
	SuccessCallBack   func([]byte) (interface{}, error)
	NotModifyCallBack func() error
}

func request(url string, callBack *CallBack) (interface{}, error) {
	client := &http.Client{
		Timeout: connect_timeout,
	}
	retry := 0
	var responseBody []byte
	var err error
	var res *http.Response
	for {
		retry++

		if retry > max_retries {
			break
		}

		res, err = client.Get(url)

		if res == nil || err != nil {
			seelog.Error("Connect Apollo Server Fail,Error:", err)
			continue
		}

		//not modified break
		switch res.StatusCode {
		case http.StatusOK:
			responseBody, err = ioutil.ReadAll(res.Body)
			if err != nil {
				seelog.Error("Connect Apollo Server Fail,Error:", err)
				continue
			}

			if callBack != nil && callBack.SuccessCallBack != nil {
				return callBack.SuccessCallBack(responseBody)
			} else {
				return nil, nil
			}
		case http.StatusNotModified:
			seelog.Warn("Config Not Modified:", err)
			if callBack != nil && callBack.NotModifyCallBack != nil {
				return nil, callBack.NotModifyCallBack()
			} else {
				return nil, nil
			}
		default:
			seelog.Error("Connect Apollo Server Fail,Error:", err)
			if res != nil {
				seelog.Error("Connect Apollo Server Fail,StatusCode:", res.StatusCode)
			}
			err = errors.New("Connect Apollo Server Fail!")
			// if error then sleep
			time.Sleep(on_error_retry_interval)
			continue
		}
	}

	seelog.Error("Over Max Retry Still Error,Error:", err)
	if err != nil {
		err = errors.New("Over Max Retry Still Error!")
	}
	return nil, err
}

func requestRecovery(appConfig *AppConfig,
	urlSuffix string,
	callBack *CallBack) (interface{}, error) {
	format := "%s%s"
	var err error
	var response interface{}

	for {
		host := appConfig.selectHost()
		if host == "" {
			return nil, err
		}

		requestUrl := fmt.Sprintf(format, host, urlSuffix)
		response, err = request(requestUrl, callBack)
		if err == nil {
			return response, err
		}

		setDownNode(host)
	}

	return nil, errors.New("Try all Nodes Still Error!")
}
