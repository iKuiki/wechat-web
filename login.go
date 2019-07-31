package wwdk

import (
	"time"

	"github.com/pkg/errors"

	"github.com/ikuiki/wwdk/datastruct"
)

// LoginChannelItem 登录时channel内返回的东西
type LoginChannelItem struct {
	Err  error
	Code LoginStatus
	Msg  string
}

// LoginStatus 登录状态枚举
type LoginStatus int32

const (
	// LoginStatusErrorOccurred 发生异常
	LoginStatusErrorOccurred LoginStatus = -1
	// LoginStatusWaitForScan 等待扫码
	// 返回Msg: 待扫码url
	LoginStatusWaitForScan LoginStatus = 1
	// LoginStatusScanedWaitForLogin 用户已经扫码
	// 返回Msg: 用户头像的base64
	LoginStatusScanedWaitForLogin LoginStatus = 2
	// LoginStatusScanedFinish 用户已同意登陆
	LoginStatusScanedFinish LoginStatus = 3
	// LoginStatusGotCookie 已获取到Cookie
	LoginStatusGotCookie LoginStatus = 4
	// LoginStatusInitFinish 登陆初始化完成
	LoginStatusInitFinish LoginStatus = 5
	// LoginStatusGotContact 已获取到联系人
	LoginStatusGotContact LoginStatus = 6
	// LoginStatusBatchGotContact 已获取到群聊成员
	LoginStatusBatchGotContact LoginStatus = 7
)

// 获取uuid用于扫码
func (wxwb *WechatWeb) getUUID(loginChannel chan<- LoginChannelItem) (uuid string) {
	uuid, _, err := wxwb.api.JsLogin()
	if err != nil {
		panic(err)
	}
	loginChannel <- LoginChannelItem{
		Code: LoginStatusWaitForScan,
		Msg:  "https://login.weixin.qq.com/l/" + uuid,
	}
	return uuid
}

// waitForScan 等待用户扫描二维码登陆
// 当扫码超时、扫码失败时，应当从getUUID方法重新开始
func (wxwb *WechatWeb) waitForScan(uuid string, loginChannel chan<- LoginChannelItem) (redirectURL string) {
	tip := "1"
	for true {
		redirectURL = func() (redirectURL string) {
			code, avatar, redirectURL, _, err := wxwb.api.Login(uuid, tip)
			if err != nil {
				panic(err)
			}
			tip = "0" // 在第二次轮询的时候tip就为0了
			switch code {
			case "200": // 确认登陆
				wxwb.logger.Info("Login success\n")
				loginChannel <- LoginChannelItem{
					Code: LoginStatusScanedFinish,
				}
				return redirectURL
			case "201": // 用户已扫码
				wxwb.logger.Info("Scan success, waiting for login\n")
				loginChannel <- LoginChannelItem{
					Code: LoginStatusScanedWaitForLogin,
					Msg:  avatar,
				}
				return "" // continue
			case "400": // 登陆失败(二维码失效)
				panic(errors.New("Login fail: qrcode has run out"))
			case "408": // 等待登陆
				time.Sleep(500 * time.Microsecond)
			default:
				panic(errors.New("Login fail: unknown response code: " + code))
			}
			return
		}()
		if redirectURL != "" {
			break
		}
	}
	return redirectURL
}

// 完成登陆,获取登陆凭据
func (wxwb *WechatWeb) getCookie(redirectURL string, loginChannel chan<- LoginChannelItem) {
	_, err := wxwb.api.WebwxNewLoginPage(redirectURL)
	if err != nil {
		panic(err)
	}
	loginChannel <- LoginChannelItem{
		Code: LoginStatusGotCookie,
	}
}

// 初始化微信,获取当前登陆用户\部分联系人
func (wxwb *WechatWeb) wxInit(loginChannel chan<- LoginChannelItem) {
	user, contactList, _, err := wxwb.api.WebwxInit()
	if err != nil {
		panic(err)
	}
	loginChannel <- LoginChannelItem{
		Code: LoginStatusInitFinish,
	}
	wxwb.userInfo.user = user
	for _, contact := range contactList {
		wxwb.userInfo.contactList[contact.UserName] = contact
	}
}

// 获取联系人
// 注：坑！此处获取到的居然不是完整的联系人，必须和init中获取到的合并后才是完整的联系人列表
func (wxwb *WechatWeb) getContactList() (err error) {
	contactList, _, err := wxwb.api.GetContact()
	if err == nil {
		for _, contact := range contactList {
			wxwb.userInfo.contactList[contact.UserName] = contact
		}
	}
	return
}

// 获取群聊的成员
func (wxwb *WechatWeb) batchGetContact() (err error) {
	var itemList []datastruct.BatchGetContactRequestListItem
	for _, contact := range wxwb.userInfo.contactList {
		if contact.IsChatroom() {
			itemList = append(itemList, datastruct.BatchGetContactRequestListItem{
				UserName: contact.UserName,
			})
		}
	}
	contactList, _, err := wxwb.api.BatchGetContact(itemList)
	if err == nil {
		for _, contact := range contactList {
			if c, ok := wxwb.userInfo.contactList[contact.UserName]; ok {
				c.MemberCount = contact.MemberCount
				c.MemberList = contact.MemberList
				c.EncryChatRoomID = contact.EncryChatRoomID
				wxwb.userInfo.contactList[c.UserName] = c
			}
		}
	}
	return
}

// Login 登陆方法总成
// param loginChannel 登陆状态channel，从中可以读取到登录情况
func (wxwb *WechatWeb) Login(loginChannel chan<- LoginChannelItem) {
	// 尝试使用已存在的登录信息登录
	go func() {
		defer close(loginChannel)
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					// 发生了panic
					loginChannel <- LoginChannelItem{
						Code: LoginStatusErrorOccurred,
						Err:  err,
					}
				} else {
					wxwb.logger.Errorf("WechatWeb.Login panic: \n%+v\n", e)
				}
			}
		}()
		// 尝试使用已存在的登录信息登录
		logined := false
		readed, _ := wxwb.readLoginInfo()
		if readed {
			wxwb.logger.Info("loaded stored login info")
			if wxwb.getContactList() == nil {
				// 获取联系人成功，则为已登陆状态
				logined = true
				wxwb.logger.Infof("reuse loginInfo [%s] logined at %v\n", wxwb.userInfo.user.NickName, wxwb.runInfo.LoginAt.Format("2006-01-02 15:04:05"))
			}
		}
		if !logined {
			if readed {
				// 仅当成功读取了login信息并且登陆失败，才输出此log
				wxwb.logger.Info("stored login info not avaliable\n")
			}
			wxwb.resetLoginInfo()
			uuid := wxwb.getUUID(loginChannel)
			redirectURL := wxwb.waitForScan(uuid, loginChannel)
			wxwb.getCookie(redirectURL, loginChannel)
			wxwb.wxInit(loginChannel)
			err := wxwb.getContactList()
			if err != nil {
				loginChannel <- LoginChannelItem{
					Code: LoginStatusErrorOccurred,
					Err:  err,
				}
				return
			}
			// 此处即认为登陆成功
			wxwb.runInfo.LoginAt = time.Now()
		}
		loginChannel <- LoginChannelItem{
			Code: LoginStatusGotContact,
		}
		// err = wxwb.StatusNotify(wxwb.userInfo.user.UserName, wxwb.userInfo.user.UserName, 3)
		// if err != nil {
		// 	return errors.New("StatusNotify error: " + err.Error())
		// }
		err := wxwb.batchGetContact()
		if err != nil {
			loginChannel <- LoginChannelItem{
				Code: LoginStatusErrorOccurred,
				Err:  err,
			}
			return
		}
		loginChannel <- LoginChannelItem{
			Code: LoginStatusBatchGotContact,
		}
		wxwb.logger.Infof("User %s has Login Success, total %d contacts\n", wxwb.userInfo.user.NickName, len(wxwb.userInfo.contactList))
		// 如有必要，记录login信息到storer
		wxwb.writeLoginInfo()
	}()
}

// Logout 退出登录
func (wxwb *WechatWeb) Logout() (err error) {
	_, err = wxwb.api.Logout()
	return
}
