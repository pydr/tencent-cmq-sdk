package cmq

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

// 创建队列
func (sdk *TencentSDK) AddQueue(queueName string) bool {
	url := sdk.Pipeline("CreateQueue", queueName, "", "")
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return false
	}

	if buf.Code != 0 {
		logs.Error("create queue failed: ", buf.Message)
		return false
	}

	logs.Info("queue [%v]: created!", queueName)
	return true
}

// 获取队列列表
func (sdk *TencentSDK) GetQueueList() []queueRet {

	logs.Info("Get All queue...")
	url := sdk.Pipeline("ListQueue", "", "", "")
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return nil
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return nil
	}

	if buf.Code != 0 {
		logs.Error("get queue list failed: ", buf.Message)
		return nil
	}

	result := buf.QueueList

	return result
}

// 获取队列剩余未消费消息数
func (sdk *TencentSDK) GetMsgCount(queueName string) (msgCount, code int) {

	logs.Info("Get queue msg left...")

	url := sdk.Pipeline("GetQueueAttributes", queueName, "", "")
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return -1, -1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return -1, -1
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return -1, -1
	}

	if buf.Code != 0 {
		logs.Error("get queue list failed: ", buf.Message)
		return -1, buf.Code
	}

	result := buf.ActiveMsgNum

	return result, 0
}

// 删除队列
func (sdk *TencentSDK) DeleteQueue(queueName string) bool {
	url := sdk.Pipeline("DeleteQueue", queueName, "", "")
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return false
	}

	if buf.Code != 0 {
		logs.Error("delete queue failed: ", buf.Message)
		return false
	}

	logs.Info("queue [%v]: deleted!", queueName)
	return true
}

// 推送消息
func (sdk *TencentSDK) PushMsg(queueName, msg string) bool {
	url := sdk.Pipeline("SendMessage", queueName, msg, "")
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return false
	}

	if buf.Code != 0 {
		logs.Error("push message failed: ", buf.Message)
		return false
	}

	logs.Info("queue [%v]: push msg success.", queueName)
	return true
}

// 拉取消息
func (sdk *TencentSDK) GetMsg(queueName string) (msg, msgHandle string, code int) {

	url := sdk.Pipeline("ReceiveMessage", queueName, "", "")
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return "", "", -1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return "", "", -1
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return "", "", -1
	}

	if buf.Code != 0 {
		logs.Error("get receive message failed: ", buf.Message)
		return "", "", buf.Code
	}

	result := buf.MsgBody
	msgHandle = buf.ReceiptHandle

	return result, msgHandle, 0

}

// 删除消息
func (sdk *TencentSDK) DeleteMsg(queueName, msgHandle string) bool {
	url := sdk.Pipeline("DeleteMessage", queueName, "", msgHandle)
	//logs.Debug("request url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("API request failed: ", err)
		return false
	}

	var buf sdkRet
	err = json.Unmarshal(body, &buf)
	if err != nil {
		logs.Error("json unMarshal failed: ", err)
		return false
	}

	if buf.Code != 0 {
		logs.Error("push message failed: ", buf.Message)
		return false
	}

	logs.Info("queue [%v]: deleted msg.", queueName)
	return true
}

func (sdk *TencentSDK) PullMsg(queueName string) (msg string, code int) {
	msg, handle, code := sdk.GetMsg(queueName)

	if msg == "" {
		logs.Error("queue [%s]: there is no any message.", queueName)
		return msg, code
	}

	result := sdk.DeleteMsg(queueName, handle)
	if !result {
		logs.Error("queue [%s]: delete message failed.", queueName)
	}

	return msg, code

}

// 初始化参数
func (sdk *TencentSDK) Init(secretId, secretKey, serverUrl, region string) {
	logs.Debug("init tencent SDK")

	sdk.secretId = secretId
	sdk.secretKey = secretKey
	sdk.serverUrl = serverUrl
	sdk.region = region
}
