package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SgipTraceReqPktLen  = HeaderPktLen + 12 + 21 + 8             //61d, 0x3d
	SgipTraceRespPktLen = HeaderPktLen + 1 + 1 + 6 + 16 + 16 + 8 //68d, 0x44
)

type SgipTraceReqPkt struct {
	SubmitSequenceNum [3]uint32
	UserNumber        string
	Reserve           string

	// used in session
	SequenceNum          [3]uint32
	SequenceNumStr       string
	SubmitSequenceNumStr string
}

func (p *SgipTraceReqPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var w = newPkgWriter(SgipTraceReqPktLen)
	// header
	w.WriteHeader(SgipTraceReqPktLen, seqNum, SGIP_TRACE)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	// body
	w.WriteInt(binary.BigEndian, p.SubmitSequenceNum[0])
	w.WriteInt(binary.BigEndian, p.SubmitSequenceNum[1])
	w.WriteInt(binary.BigEndian, p.SubmitSequenceNum[2])
	w.WriteFixedSizeString(p.UserNumber, 21)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipTraceReqPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)
	s := make([]byte, 12)
	r.ReadBytes(s)
	p.SubmitSequenceNum = [3]uint32{
		binary.BigEndian.Uint32(s[:4]),
		binary.BigEndian.Uint32(s[4:8]),
		binary.BigEndian.Uint32(s[8:]),
	}
	p.UserNumber = string(r.ReadCString(21))
	p.Reserve = string(r.ReadCString(8))

	p.SubmitSequenceNumStr = GenSequenceNumStr(p.SubmitSequenceNum)
	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)
	return r.Error()
}

func (p *SgipTraceReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Trace Req ---")
	fmt.Fprintln(&b, "SubmitSequenceNumStr", p.SubmitSequenceNumStr)
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)

	return b.String()
}

type SgipTraceRespPkt struct {
	Count       uint8
	Result      uint8
	NodeId      string
	ReceiveTime string
	SendTime    string
	Reserve     string

	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}

func (p *SgipTraceRespPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var w = newPkgWriter(SgipTraceRespPktLen)
	// header
	w.WriteHeader(SgipTraceRespPktLen, seqNum, SGIP_TRACE_RESP)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	// body
	w.WriteByte(p.Count)
	w.WriteByte(p.Result)
	w.WriteFixedSizeString(p.NodeId, 6)
	w.WriteFixedSizeString(p.ReceiveTime, 16)
	w.WriteFixedSizeString(p.SendTime, 16)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipTraceRespPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)
	p.Count = r.ReadByte()
	p.Result = r.ReadByte()
	p.NodeId = string(r.ReadCString(6))
	p.ReceiveTime = string(r.ReadCString(16))
	p.SendTime = string(r.ReadCString(16))
	p.Reserve = string(r.ReadCString(8))

	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)
	return r.Error()
}

func (p *SgipTraceRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Trace Resp ---")
	fmt.Fprintln(&b, "Count", p.Count)
	fmt.Fprintln(&b, "Result", p.Result)
	fmt.Fprintln(&b, "NodeId: ", p.NodeId)
	fmt.Fprintln(&b, "ReceiveTime: ", p.ReceiveTime)
	fmt.Fprintln(&b, "SendTime: ", p.SendTime)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)

	return b.String()
}
