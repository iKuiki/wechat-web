package wxweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/mdp/qrterminal"
	"github.com/yinhui87/wechat-web/conf"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"
)

func (this *WechatWeb) getUUID() (uuid string, err error) {
	req := httplib.Get("https://login.weixin.qq.com/jslogin")
	req.Param("appid", conf.APP_ID)
	req.Param("fun", "new")
	req.Param("lang", "zh_CN")
	req.Param("_", tool.GetWxTimeStamp())
	resp, err := req.String()
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	ret := tool.AnalysisWxWindowRespond(resp)
	if ret["window.QRLogin.code"] != "200" {
		return "", errors.New("window.QRLogin.code = " + ret["window.QRLogin.code"])
	}
	return ret["window.QRLogin.uuid"], nil
}

func (this *WechatWeb) getQrCode(uuid string) (err error) {
	req := httplib.Post("https://login.weixin.qq.com/qrcode/" + uuid)
	req.Param("t", "webwx")
	req.Param("_", tool.GetWxTimeStamp())
	_, err = req.String()
	if err != nil {
		return err
	}
	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)
	return nil
}

func (this *WechatWeb) waitForScan(uuid string) (redirectUrl string, err error) {
	var ret map[string]string
	req := httplib.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login")
	req.Param("tip", "1")
	req.Param("uuid", uuid)
	req.Param("_", tool.GetWxTimeStamp())
	_, err = req.String()
	if err != nil {
		return "", errors.New("waitForScan request error: " + err.Error())
	}
	log.Println("Scan success, waiting for login")
	for true {
		req := httplib.Get("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login")
		req.Param("tip", "1")
		req.Param("uuid", uuid)
		req.Param("_", tool.GetWxTimeStamp())
		resp, err := req.String()
		if err != nil {
			log.Println("waitForScan request error: " + err.Error())
			continue
		}
		ret = tool.AnalysisWxWindowRespond(resp)
		if ret["window.code"] != "200" {
			time.Sleep(500 * time.Microsecond)
			continue
		}
		break
	}
	return ret["window.redirect_uri"], nil
}

func (this *WechatWeb) getCookie(redirectUrl, userAgent string) (err error) {
	u, err := url.Parse(redirectUrl)
	if err != nil {
		return errors.New("redirect_url parse fail: " + err.Error())
	}
	query := u.Query()
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	req.Param("ticket", query.Get("ticket"))
	req.Param("uuid", query.Get("uuid"))
	req.Param("lang", "zh_CN")
	req.Param("scan", query.Get("scan"))
	req.Param("fun", "new")
	req.SetUserAgent(userAgent)
	resp, err := req.Response()
	if err != nil {
		return errors.New("getCookie request error: " + err.Error())
	}
	cookies := make(map[string]string)
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Read respond body error: " + err.Error())
	}
	var bodyResp datastruct.GetCookieRespond
	err = xml.Unmarshal(body, &bodyResp)
	if err != nil {
		return errors.New("Unmarshal respond xml error: " + err.Error())
	}
	this.cookie = &wechatCookie{
		Wxuin:      cookies["wxuin"],
		Wxsid:      cookies["wxsid"],
		Uvid:       cookies["webwxuvid"],
		DataTicket: cookies["webwx_data_ticket"],
		AuthTicket: cookies["webwx_auth_ticket"],
		PassTicket: bodyResp.PassTicket,
	}
	return nil
}

func (this *WechatWeb) wxInit() (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxinit")
	body := datastruct.WxInitRequestBody{
		BaseRequest: getBaseRequest(this.cookie, this.deviceId),
	}
	req.Header("Content-Type", "application/json")
	req.Header("charset", "UTF-8")
	req.Param("r", tool.GetWxTimeStamp())
	setWechatCookie(req, this.cookie)
	resp := datastruct.WxInitRespond{}
	data, err := json.Marshal(body)
	if err != nil {
		return errors.New("json.Marshal error: " + err.Error())
	}
	req.Body(data)
	// err = req.ToJSON(&resp)
	r, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if resp.BaseResponse.Ret != 0 {
		return errors.New(fmt.Sprintf("respond ret error: %d", resp.BaseResponse.Ret))
	}
	this.user = resp.User
	this.syncKey = resp.SyncKey
	this.sKey = resp.SKey
	return nil
}

func (this *WechatWeb) getContactList() (err error) {
	req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact")
	req.Param("r", tool.GetWxTimeStamp())
	setWechatCookie(req, this.cookie)
	req.Body([]byte("{}"))
	resp := datastruct.GetContactRespond{}
	r, err := req.Bytes()
	if err != nil {
		return errors.New("request error: " + err.Error())
	}
	err = json.Unmarshal(r, &resp)
	if err != nil {
		return errors.New("respond json Unmarshal to struct fail: " + err.Error())
	}
	if resp.BaseResponse.Ret != 0 {
		return errors.New(fmt.Sprintf("respond ret error: %d", resp.BaseResponse.Ret))
	}
	this.contactList = resp.MemberList
	return nil
}

func (this *WechatWeb) Login() (err error) {
	uuid, err := this.getUUID()
	if err != nil {
		return errors.New("Get UUID fail: " + err.Error())
	}
	err = this.getQrCode(uuid)
	if err != nil {
		return errors.New("Get QrCode fail: " + err.Error())
	}
	redirectUrl, err := this.waitForScan(uuid)
	if err != nil {
		return errors.New("waitForScan error: " + err.Error())
	}
	// panic(redirectUrl)
	err = this.getCookie(redirectUrl, this.userAgent)
	if err != nil {
		return errors.New("getCookie error: " + err.Error())
	}
	err = this.wxInit()
	if err != nil {
		return errors.New("wxInit error: " + err.Error())
	}
	err = this.getContactList()
	if err != nil {
		return errors.New("getContactList error: " + err.Error())
	}
	log.Printf("User %s has Login Success, total %d contacts\n", this.user.NickName, len(this.contactList))
	return nil
}
