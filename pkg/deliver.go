package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SgipDeliverReqPktLen = HeaderPktLen + 21 + 21 + 1 + 1 + 1 + 4 + 8 //77d, 0x4d
)

type SgipDeliverReqPkt struct {
	UserNumber     string
	SPNumber       string
	TP_pid         uint8
	TP_udhi        uint8
	MessageCoding  uint8
	MessageLength  uint32
	MessageContent string
	Reserve        string

	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}

func (p *SgipDeliverReqPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	pktLen := SgipDeliverReqPktLen + p.MessageLength
	var w = newPkgWriter(pktLen)
	// header
	w.WriteHeader(pktLen, seqNum, SGIP_DELIVER)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	// body
	w.WriteFixedSizeString(p.UserNumber, 21)
	w.WriteFixedSizeString(p.SPNumber, 21)
	w.WriteByte(p.TP_pid)
	w.WriteByte(p.TP_udhi)
	w.WriteByte(p.MessageCoding)
	w.WriteInt(binary.BigEndian, p.MessageLength)
	w.WriteString(p.MessageContent)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipDeliverReqPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)
	p.UserNumber = string(r.ReadCString(21))
	p.SPNumber = string(r.ReadCString(21))
	p.TP_pid = r.ReadByte()
	p.TP_udhi = r.ReadByte()
	p.MessageCoding = r.ReadByte()
	r.ReadInt(binary.BigEndian, &p.MessageLength)
	msgContent := make([]byte, p.MessageLength)
	r.ReadBytes(msgContent)
	p.MessageContent = string(msgContent)
	p.Reserve = string(r.ReadCString(8))

	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)
	return r.Error()
}

func (p *SgipDeliverReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Deliver Req ---")
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "SPNumber: ", p.SPNumber)
	fmt.Fprintln(&b, "TP_pid: ", p.TP_pid)
	fmt.Fprintln(&b, "TP_udhi: ", p.TP_udhi)
	fmt.Fprintln(&b, "MessageCoding: ", p.MessageCoding)
	fmt.Fprintln(&b, "MessageLength: ", p.MessageLength)
	fmt.Fprintln(&b, "MessageContent: ", p.MessageContent)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)

	return b.String()
}

type SgipDeliverRespPkt struct {
	SgipRespPkt
}

func (p *SgipDeliverRespPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	return p.SgipRespPkt.Pack(seqNum, SGIP_DELIVER_RESP)
}

func (p *SgipDeliverRespPkt) Unpack(data []byte) error {
	return p.SgipRespPkt.Unpack(data)
}

func (p *SgipDeliverRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Deliver Resp ---")
	fmt.Fprintln(&b, p.SgipRespPkt.String())

	return b.String()
}
