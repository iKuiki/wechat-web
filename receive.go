package wxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// getMessage 同步消息
// 如果同步状态接口返回有新消息需要同步，通过此接口从服务器中获取新消息
func (wxwb *WechatWeb) getMessage() (gmResp datastruct.GetMessageRespond, err error) {
	gmResp = datastruct.GetMessageRespond{}
	data, err := json.Marshal(datastruct.GetMessageRequest{
		BaseRequest: wxwb.baseRequest(),
		SyncKey:     wxwb.syncKey,
		Rr:          ^time.Now().Unix() + 1,
	})
	if err != nil {
		return gmResp, errors.New("Marshal request body to json fail: " + err.Error())
	}
	params := url.Values{}
	params.Set("sid", wxwb.cookie.Wxsid)
	params.Set("skey", wxwb.sKey)
	// params.Set("pass_ticket", wxwb.PassTicket)
	resp, err := wxwb.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync?"+params.Encode(),
		"application/json;charset=UTF-8",
		bytes.NewReader(data))
	if err != nil {
		return gmResp, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &gmResp)
	if err != nil {
		return gmResp, errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if gmResp.BaseResponse.Ret != 0 {
		return gmResp, errors.New("respond error ret: " + strconv.FormatInt(gmResp.BaseResponse.Ret, 10))
	}
	// if gmResp.AddMsgCount > 0 {
	// 	fmt.Println(string(resp))
	// 	panic(nil)
	// }
	return gmResp, nil
}

// SaveMessageImage 保存消息图片到指定位置
func (wxwb *WechatWeb) SaveMessageImage(msg datastruct.Message) (filename string, err error) {
	params := url.Values{}
	params.Set("MsgID", msg.MsgID)
	params.Set("skey", wxwb.sKey)
	// params.Set("type", "slave")
	resp, err := wxwb.client.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg?" + params.Encode())
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	filename = msg.MsgID + ".jpg"
	_, err = tool.WriteToFile(filename, resp.Body)
	if err != nil {
		return "", errors.New("WriteToFile error: " + err.Error())
	}
	return filename, nil
}

// SaveMessageVoice 保存消息声音到指定位置
func (wxwb *WechatWeb) SaveMessageVoice(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.VoiceMsg {
		return "", errors.New("Message type wrong")
	}
	params := url.Values{}
	params.Set("MsgID", msg.MsgID)
	params.Set("skey", wxwb.sKey)
	// params.Set("type", "slave")
	resp, err := wxwb.client.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvoice?" + params.Encode())
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	filename = msg.MsgID + ".mp3"
	_, err = tool.WriteToFile(filename, resp.Body)
	if err != nil {
		return "", errors.New("WriteToFile error: " + err.Error())
	}
	return filename, nil
}

// SaveMessageVideo 保存消息视频到指定位置
func (wxwb *WechatWeb) SaveMessageVideo(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.LittleVideoMsg {
		return "", errors.New("Message type wrong")
	}
	params := url.Values{}
	params.Set("msgid", msg.MsgID)
	params.Set("skey", wxwb.sKey)
	req, err := http.NewRequest("GET", "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvideo?"+params.Encode(), strings.NewReader(""))
	if err != nil {
		return "", errors.New("create request error: " + err.Error())
	}
	req.Header.Set("Range", "bytes=0-")
	resp, err := wxwb.client.Do(req)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	filename = msg.MsgID + ".mp4"
	n, err := tool.WriteToFile(filename, resp.Body)
	if err != nil {
		return "", errors.New("WriteToFile error: " + err.Error())
	}
	if int64(n) != resp.ContentLength {
		return filename, errors.New("File size wrong")
	}
	return filename, nil
}