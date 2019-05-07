package cmq

type queueRet struct {
	QueueId   string `json:"queueIdd"`
	QueueName string `json:"queueName"`
}

type msgId struct {
	MsgId string `json:"msgId"`
}

type msgInfo struct {
	MsgBody string `json:"msgBody"`
	MsgId   string `json:"msgId"`
}

type sdkRet struct {
	Code         int        `json:"code"`
	Message      string     `json:"message"`
	RequestId    string     `json:"requestId"`
	QueueId      string     `json:"queueId"`
	TotalCount   int        `json:"totalCount"`
	QueueList    []queueRet `json:"queueList"`
	ActiveMsgNum int        `json:"activeMsgNum"`

	MsgId         string    `json:"msgId"`
	MsgList       []msgId   `json:"msgList"`
	MsgBody       string    `json:"msgBody"`
	MsgInfoList   []msgInfo `json:"msgInfoList"`
	ReceiptHandle string    `json:"receiptHandle"`
}
