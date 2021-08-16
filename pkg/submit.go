package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SgipSubmitRespPktLen = HeaderPktLen + 10 + 4 //26d, 0x1a
)

type SgipSubmitReqPkt struct {
	SPNumber         string
	ChargeNumber     string
	UserCount        uint8
	UserNumber       []string // 接收该短消息的手机号，该字段重复UserCount指定的次数，手机号码前加“86”国别标志
	CorpId           string
	ServiceType      string
	FeeType          uint8
	FeeValue         string
	GivenValue       string
	AgentFlag        uint8
	MorelatetoMTFlag uint8
	Priority         uint8
	ExpireTime       string
	ScheduleTime     string
	ReportFlag       uint8
	TP_pid           uint8
	TP_udhi          uint8
	MessageCoding    uint8
	MessageType      uint8
	MessageLength    uint32
	MessageContent   string
	Reserve          string

	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}

func (p *SgipSubmitReqPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var pktLen = HeaderPktLen + 123 + uint32(p.UserCount)*21 + p.MessageLength
	var w = newPkgWriter(pktLen)
	// header
	w.WriteHeader(pktLen, seqNum, SGIP_SUBMIT)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	// body
	w.WriteFixedSizeString(p.SPNumber, 21)
	w.WriteFixedSizeString(p.ChargeNumber, 21)
	w.WriteByte(p.UserCount)
	for _, n := range p.UserNumber {
		w.WriteFixedSizeString(n, 21)
	}

	w.WriteFixedSizeString(p.CorpId, 5)
	w.WriteFixedSizeString(p.ServiceType, 10)
	w.WriteByte(p.FeeType)
	w.WriteFixedSizeString(p.FeeValue, 6)
	w.WriteFixedSizeString(p.GivenValue, 6)
	w.WriteByte(p.AgentFlag)
	w.WriteByte(p.MorelatetoMTFlag)
	w.WriteByte(p.Priority)
	w.WriteFixedSizeString(p.ExpireTime, 16)
	w.WriteFixedSizeString(p.ScheduleTime, 16)
	w.WriteByte(p.ReportFlag)
	w.WriteByte(p.TP_pid)
	w.WriteByte(p.TP_udhi)
	w.WriteByte(p.MessageCoding)
	w.WriteByte(p.MessageType)
	w.WriteInt(binary.BigEndian, p.MessageLength)
	w.WriteString(p.MessageContent)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipSubmitReqPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)

	p.SPNumber = string(r.ReadCString(21))
	p.ChargeNumber = string(r.ReadCString(21))
	p.UserCount = r.ReadByte()
	for i := 0; i < int(p.UserCount); i++ {
		p.UserNumber = append(p.UserNumber, string(r.ReadCString(21)))
	}
	p.CorpId = string(r.ReadCString(5))
	p.ServiceType = string(r.ReadCString(10))
	p.FeeType = r.ReadByte()
	p.FeeValue = string(r.ReadCString(6))
	p.GivenValue = string(r.ReadCString(6))
	p.AgentFlag = r.ReadByte()
	p.MorelatetoMTFlag = r.ReadByte()
	p.Priority = r.ReadByte()
	p.ExpireTime = string(r.ReadCString(16))
	p.ScheduleTime = string(r.ReadCString(16))
	p.ReportFlag = r.ReadByte()
	p.TP_pid = r.ReadByte()
	p.TP_udhi = r.ReadByte()
	p.MessageCoding = r.ReadByte()
	p.MessageType = r.ReadByte()
	r.ReadInt(binary.BigEndian, &p.MessageLength)
	msgContent := make([]byte, p.MessageLength)
	r.ReadBytes(msgContent)
	p.MessageContent = string(msgContent)
	p.Reserve = string(r.ReadCString(8))

	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)
	return r.Error()
}

func (p *SgipSubmitReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Submit Req ---")
	fmt.Fprintln(&b, "SPNumber: ", p.SPNumber)
	fmt.Fprintln(&b, "ChargeNumber: ", p.ChargeNumber)
	fmt.Fprintln(&b, "UserCount: ", p.UserCount)
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "CorpId: ", p.CorpId)
	fmt.Fprintln(&b, "ServiceType: ", p.ServiceType)
	fmt.Fprintln(&b, "FeeType: ", p.FeeType)
	fmt.Fprintln(&b, "FeeValue: ", p.FeeValue)
	fmt.Fprintln(&b, "GivenValue: ", p.GivenValue)
	fmt.Fprintln(&b, "AgentFlag: ", p.AgentFlag)
	fmt.Fprintln(&b, "MorelatetoMTFlag: ", p.MorelatetoMTFlag)
	fmt.Fprintln(&b, "Priority: ", p.Priority)
	fmt.Fprintln(&b, "ExpireTime: ", p.ExpireTime)
	fmt.Fprintln(&b, "ScheduleTime: ", p.ScheduleTime)
	fmt.Fprintln(&b, "ReportFlag: ", p.ReportFlag)
	fmt.Fprintln(&b, "TP_pid: ", p.TP_pid)
	fmt.Fprintln(&b, "TP_udhi: ", p.TP_udhi)
	fmt.Fprintln(&b, "MessageCoding: ", p.MessageCoding)
	fmt.Fprintln(&b, "MessageType: ", p.MessageType)
	fmt.Fprintln(&b, "MessageLength: ", p.MessageLength)
	fmt.Fprintln(&b, "MessageContent: ", p.MessageContent)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)

	return b.String()
}

type SgipSubmitRespPkt struct {
	SgipRespPkt
}

func (p *SgipSubmitRespPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	return p.SgipRespPkt.Pack(seqNum, SGIP_SUBMIT_RESP)
}

func (p *SgipSubmitRespPkt) Unpack(data []byte) error {
	return p.SgipRespPkt.Unpack(data)
}

func (p *SgipSubmitRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Submit Resp ---")
	fmt.Fprintln(&b, p.SgipRespPkt.String())

	return b.String()
}
