package pkg

import (
	"errors"
	"strconv"
)

type Status uint8

func (s *Status) Data() uint8 {
	return uint8(*s)
}

func (s *Status) Error() error {
	return errors.New(strconv.Itoa(int(*s)) + " : " + s.String())
}

func (s Status) String() string {

	var msg string
	switch s {
	case 0:
		msg = "成功"

	// 1-20所指错误一般在各类命令的应答中用到
	case 1:
		msg = "非法登录"
	case 2:
		msg = "重复登录"
	case 3:
		msg = "连接过多"
	case 4:
		msg = "登录类型错"
	case 5:
		msg = "参数格式错"
	case 6:
		msg = "非法手机号码"
	case 7:
		msg = "消息ID错"
	case 8:
		msg = "信息长度错"
	case 9:
		msg = "非法序列号"
	case 10:
		msg = "非法操作GNS"
	case 11:
		msg = "节点忙"

	// 21-32 所指错误一般在report命令中用到
	case 21:
		msg = "目的地址不可达"
	case 22:
		msg = "路由错"
	case 23:
		msg = "路由不存在"
	case 24:
		msg = "计费号码无效"
	case 25:
		msg = "用户不能通信"
	case 26:
		msg = "手机内存不足"
	case 27:
		msg = "手机不支持短消息"
	case 28:
		msg = "手机接收短消息出现错误"
	case 29:
		msg = "不知道的用户"
	case 30:
		msg = "不提供此功能"
	case 31:
		msg = "非法设备"
	case 32:
		msg = "系统失败"
	case 33:
		msg = "短信中心队列满"

	default:
		msg = "Status Unknown: " + strconv.Itoa(int(s))
	}

	return msg
}

const (
	STAT_OK Status = iota
)
