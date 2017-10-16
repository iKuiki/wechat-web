package wxweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
	"github.com/yinhui87/wechat-web/tool"
	"html"
	"strconv"
	"strings"
)

// StatusNotify 消息通知
func (wxwb *WechatWeb) StatusNotify(fromUserName, toUserName string) (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify")
	req.Param("pass_ticket", wxwb.cookie.PassTicket)
	setWechatCookie(req, wxwb.cookie)
	msgID, _ := strconv.ParseInt(tool.GetWxTimeStamp(), 10, 64)
	reqBody := datastruct.StatusNotifyRequest{
		BaseRequest:  getBaseRequest(wxwb.cookie, wxwb.deviceID),
		ClientMsgID:  msgID,
		Code:         1,
		FromUserName: fromUserName,
		ToUserName:   toUserName,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return errors.New("Marshal request body to json fail: " + err.Error())
	}
	req.Body(body)
	resp, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	var snResp datastruct.StatusNotifyRespond
	err = json.Unmarshal(resp, &snResp)
	if err != nil {
		return errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if snResp.BaseResponse.Ret != 0 {
		return errors.New("respond error ret: " + strconv.FormatInt(snResp.BaseResponse.Ret, 10))
	}
	return nil
}

func (wxwb *WechatWeb) messageProcesser(msg *datastruct.Message) (err error) {
	context := Context{App: wxwb, hasStop: false}
	switch msg.MsgType {
	case datastruct.TextMsg:
		for _, v := range wxwb.messageHook[datastruct.TextMsg] {
			if f, ok := v.(TextMessageHook); ok {
				f(&context, *msg)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.ImageMsg:
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		var imgContent appmsg.ImageMsgContent
		err = xml.Unmarshal([]byte(msg.Content), &imgContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.ImageMsg] {
			if f, ok := v.(ImageMessageHook); ok {
				f(&context, *msg, imgContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.AnimationEmotionsMsg:
		msg.Content = html.UnescapeString(msg.Content)
		var emojiContent appmsg.EmotionMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &emojiContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.AnimationEmotionsMsg] {
			if f, ok := v.(EmotionMessageHook); ok {
				f(&context, *msg, emojiContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.RevokeMsg:
		msg.Content = html.UnescapeString(msg.Content)
		var revokeContent appmsg.RevokeMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &revokeContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.RevokeMsg] {
			if f, ok := v.(RevokeMessageHook); ok {
				f(&context, *msg, revokeContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.LittleVideoMsg:
		msg.Content = strings.Replace(html.UnescapeString(msg.Content), "<br/>", "", -1)
		var videoContent appmsg.VideoMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &videoContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.LittleVideoMsg] {
			if f, ok := v.(VideoMessageHook); ok {
				f(&context, *msg, videoContent)
			}
			if context.hasStop {
				break
			}
		}
	case datastruct.VoiceMsg:
		msg.Content = html.UnescapeString(msg.Content)
		var voiceContent appmsg.VoiceMsgContent
		err := xml.Unmarshal([]byte(msg.Content), &voiceContent)
		if err != nil {
			return errors.New("Unmarshal message content to struct: " + err.Error())
		}
		for _, v := range wxwb.messageHook[datastruct.VoiceMsg] {
			if f, ok := v.(VoiceMessageHook); ok {
				f(&context, *msg, voiceContent)
			}
			if context.hasStop {
				break
			}
		}
	default:
		return fmt.Errorf("Unknown MsgType %v: %#v", msg.MsgType, msg)
	}
	return nil
}
