package aliim

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type ConfigST struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	IsPro     bool   `json:"is_pro"`
	IsUseHTTP bool   `json:"is_use_http"`
}

var config ConfigST

// Init 初始化
func Init(cfg *ConfigST) {
	config = *cfg
}

// 阿里接口请求方法
func IMPost(m map[string]string) (success bool, resData []byte) {
	if config.AppKey == "" || config.AppSecret == "" {
		return false, []byte("appkey or appsecret is requierd!")
	}

	body, size := getHttpBody(m)
	client := &http.Client{}
	var req *http.Request
	var err error
	if !config.IsUseHTTP {
		log.Println("addr is " + GetHttpServerAddr(config.IsPro))
		req, err = http.NewRequest("POST", GetHttpServerAddr(config.IsPro), body)
	} else {
		log.Println("addr is " + GetHttpServerAddr(config.IsPro))
		req, err = http.NewRequest("POST", GetHttpsServerAddr(config.IsPro), body)
	}

	if err != nil {
		return false, []byte(err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.ContentLength = size
	resp, err := client.Do(req)
	if err != nil {
		return false, []byte(err.Error())
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return true, data
}

// 获取消息体
func getHttpBody(m map[string]string) (reader io.Reader, size int64) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	v := url.Values{}

	siginString := config.AppSecret
	for _, k := range keys {
		v.Set(k, m[k])
		siginString += k + m[k]
	}
	siginString += config.AppSecret
	signByte := md5.Sum([]byte(siginString))
	sign := strings.ToUpper(fmt.Sprintf("%x", signByte))
	v.Set("sign", sign)
	return ioutil.NopCloser(strings.NewReader(v.Encode())), int64(len(v.Encode()))
}
