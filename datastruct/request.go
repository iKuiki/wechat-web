package datastruct

// BaseRequest 基本请求结构，包含通用的请求信息
type BaseRequest struct {
	DeviceID string `json:"DeviceID"`
	Sid      string `json:"Sid"`
	Skey     string `json:"Skey"`
	Uin      string `json:"Uin"`
}

// WxInitRequestBody 微信初始化请求
type WxInitRequestBody struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
}

// BatchGetContactRequestListItem 获取群组联系人的请求的列表元素
type BatchGetContactRequestListItem struct {
	ChatRoomID string `json:"ChatRoomId"`
	UserName   string `json:"UserName"`
}

// BatchGetContactRequest 获取群组联系人的请求
type BatchGetContactRequest struct {
	BaseRequest *BaseRequest                     `json:"BaseRequest"`
	Count       int64                            `json:"Count"`
	List        []BatchGetContactRequestListItem `json:"List"`
}

// WebwxSyncRequest 获取新消息的轮询请求
type WebwxSyncRequest struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
	SyncKey     *SyncKey     `json:"SyncKey"`
	Rr          int64        `json:"rr"`
}

// StatusNotifyRequest 状态通知请求
type StatusNotifyRequest struct {
	BaseRequest  *BaseRequest `json:"BaseRequest"`
	ClientMsgID  int64        `json:"ClientMsgId"`
	Code         int64        `json:"Code"`
	FromUserName string       `json:"FromUserName"`
	ToUserName   string       `json:"ToUserName"`
}

// TextMessage 发送纯文本消息，用SendMessage也一样
type TextMessage struct {
	ClientMsgID  string      `json:"ClientMsgId"`
	Content      string      `json:"Content"`
	FromUserName string      `json:"FromUserName"`
	LocalID      string      `json:"LocalID"`
	ToUserName   string      `json:"ToUserName"`
	Type         MessageType `json:"Type"`
}

// SendMessage 发送消息，可发送带媒体的消息
type SendMessage struct {
	ClientMsgID  string      `json:"ClientMsgId"`
	Content      string      `json:"Content"`
	FromUserName string      `json:"FromUserName"`
	LocalID      string      `json:"LocalID"`
	MediaID      string      `json:"MediaId"`
	ToUserName   string      `json:"ToUserName"`
	Type         MessageType `json:"Type"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
	Msg         *SendMessage `json:"Msg"`
	Scene       int64        `json:"Scene"`
}

// RevokeMessageRequest 撤回消息请求，需要附带要撤回消息的客户端、服务端消息ID
type RevokeMessageRequest struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
	ClientMsgID string       `json:"ClientMsgId"`
	SvrMsgID    string       `json:"SvrMsgId"`
	ToUserName  string       `json:"ToUserName"`
}

// ModifyRemarkRequest 修改用户备注的请求
type ModifyRemarkRequest struct {
	BaseRequest *BaseRequest `json:"BaseRequest"`
	CmdID       int64        `json:"CmdId"`
	RemarkName  string       `json:"RemarkName"`
	UserName    string       `json:"UserName"`
}

// ModifyChatRoomTopicRequest 修改群名的请求
type ModifyChatRoomTopicRequest struct {
	BaseRequest  *BaseRequest `json:"BaseRequest"`
	ChatRoomName string       `json:"ChatRoomName"`
	NewTopic     string       `json:"NewTopic"`
}
