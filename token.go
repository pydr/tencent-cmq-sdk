package cmq

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TencentSDK struct {

	// 公共参数
	secretId  string
	secretKey string
	action    string
	region    string
	serverUrl string
}

// 查询参数排序
func strSort(queryMap map[string]string) string {

	var qureys bytes.Buffer
	var keys []string
	for k := range queryMap {
		keys = append(keys, k)
	}

	// 按key排序
	sort.Strings(keys)

	for _, k := range keys {
		qureys.WriteString(k)
		qureys.WriteString("=")
		qureys.WriteString(queryMap[k])
		qureys.WriteString("&")
	}

	rs := []rune(qureys.String())
	rs_len := len(rs)

	return string(rs[0 : rs_len-1])
}

// 请求参数组装
func (sdk *TencentSDK) genQueryStr(action, queueName, msg, receiptHandle string) string {

	var queryMap = make(map[string]string)
	//queryMap["ProjectId"] = sdk.projectId
	//queryMap["AppId"] = sdk.appId
	queryMap["SecretId"] = sdk.secretId
	queryMap["SecretKey"] = sdk.secretKey
	queryMap["Action"] = action
	queryMap["Nonce"] = sdk.genNonce()
	queryMap["Region"] = sdk.region

	var t = int(time.Now().Unix())
	queryMap["Timestamp"] = strconv.Itoa(t)

	if queueName != "" {
		queryMap["queueName"] = queueName
	}

	if msg != "" {
		queryMap["msgBody"] = msg
	}

	if receiptHandle != "" {
		queryMap["receiptHandle"] = receiptHandle
	}

	queryStr := strSort(queryMap)

	//logs.Debug("queryStr is： ", queryStr)
	return queryStr

}

// 生成随机数
func (sdk *TencentSDK) genNonce() string {

	var t = time.Now().Unix()

	// 随机数生成器
	source := rand.NewSource(t)
	r := rand.New(source)

	num := r.Intn(10000)
	return strconv.Itoa(num)
}

// 生成签名串
func (sdk *TencentSDK) genSign(serverUrl, secretKey string) string {

	var url string
	rs := []rune(serverUrl)

	if strings.HasPrefix(serverUrl, "https") {
		url = "GET" + string(rs[8:])
	} else {
		url = "GET" + string(rs[7:])
	}

	hmac := hmac.New(sha1.New, []byte(secretKey))
	hmac.Write([]byte(url))
	encoded_byte_arr := hmac.Sum([]byte(""))

	var signature = base64.StdEncoding.EncodeToString(encoded_byte_arr)

	//logs.Debug("sign is: ", signature)
	return signature
}

// 生成请求url
func (sdk *TencentSDK) genQueryUrl(serverUrl string, queryStr string) string {

	var queryUrl bytes.Buffer
	queryUrl.WriteString(serverUrl)
	//queryUrl.WriteString(sdk.appId)
	queryUrl.WriteString("?")
	queryUrl.WriteString(queryStr)

	//logs.Debug("url is: ", queryUrl.String())
	return queryUrl.String()
}

func (sdk *TencentSDK) Pipeline(action, queueName, msg, receiptHandle string) string {
	queryStr := sdk.genQueryStr(action, queueName, msg, receiptHandle)

	queryUrl := sdk.genQueryUrl(sdk.serverUrl, queryStr)

	sign := sdk.genSign(queryUrl, sdk.secretKey)

	url := queryUrl + "&Signature=" + url.QueryEscape(sign)
	return url
}
