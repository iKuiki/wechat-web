package main

import (
	"fmt"
	"github.com/yinhui87/wechat-web"
	"testing"
)

func TestContact(t *testing.T) {
	wx := wxweb.NewWechatWeb()
	err := wx.Login()
	if err != nil {
		panic("WxWeb Login error: " + err.Error())
	}
	contacts := wx.GetContactList()
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s 	Alias	: %v\n", v.NickName, v.Alias)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s AppAccountFlag: %v\n", v.NickName, v.AppAccountFlag)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s AttrStatus: %v\n", v.NickName, v.AttrStatus)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s ChatRoomID: %v\n", v.NickName, v.ChatRoomID)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s City: %v\n", v.NickName, v.City)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s ContactFlag: %v\n", v.NickName, v.ContactFlag)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s DisplayName: %v\n", v.NickName, v.DisplayName)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s EncryChatRoomID: %v\n", v.NickName, v.EncryChatRoomID)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s HeadImgURL: %v\n", v.NickName, v.HeadImgURL)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s HideInputBarFlag: %v\n", v.NickName, v.HideInputBarFlag)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s IsOwner: %v\n", v.NickName, v.IsOwner)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s KeyWord: %v\n", v.NickName, v.KeyWord)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s MemberCount: %v\n", v.NickName, v.MemberCount)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s MemberList: %v\n", v.NickName, v.MemberList)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s NickName: %v\n", v.NickName, v.NickName)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s OwnerUin: %v\n", v.NickName, v.OwnerUin)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s PYInitial: %v\n", v.NickName, v.PYInitial)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s PYQuanPin: %v\n", v.NickName, v.PYQuanPin)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s Province: %v\n", v.NickName, v.Province)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s RemarkName: %v\n", v.NickName, v.RemarkName)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s RemarkPYInitial: %v\n", v.NickName, v.RemarkPYInitial)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s RemarkPYQuanPin: %v\n", v.NickName, v.RemarkPYQuanPin)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s Sex: %v\n", v.NickName, v.Sex)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s Signature: %v\n", v.NickName, v.Signature)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s SnsFlag: %v\n", v.NickName, v.SnsFlag)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s StarFriend: %v\n", v.NickName, v.StarFriend)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s Statues: %v\n", v.NickName, v.Statues)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s Uin: %v\n", v.NickName, v.Uin)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s UniFriend: %v\n", v.NickName, v.UniFriend)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s UserName: %v\n", v.NickName, v.UserName)
	}
	fmt.Println("")
	for _, v := range contacts {
		fmt.Printf("%s VerifyFlag: %v\n", v.NickName, v.VerifyFlag)
	}
}