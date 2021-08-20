package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SgipReportReqPktLen = HeaderPktLen + 12 + 1 + 21 + 1 + 1 + 8 //64d, 0x40
)

type SgipReportReqPkt struct {
	SubmitSequenceNum [3]uint32
	ReportType        uint8
	UserNumber        string
	State             uint8
	ErrorCode         uint8
	Reserve           string

	// used in session
	SequenceNum          [3]uint32
	SequenceNumStr       string
	SubmitSequenceNumStr string
}

func (p *SgipReportReqPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var w = newPkgWriter(SgipReportReqPktLen)
	// header
	w.WriteHeader(SgipReportReqPktLen, seqNum, SGIP_REPORT)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	// body
	w.WriteInt(binary.BigEndian, p.SubmitSequenceNum[0])
	w.WriteInt(binary.BigEndian, p.SubmitSequenceNum[1])
	w.WriteInt(binary.BigEndian, p.SubmitSequenceNum[2])
	w.WriteByte(p.ReportType)
	w.WriteFixedSizeString(p.UserNumber, 21)
	w.WriteByte(p.State)
	w.WriteByte(p.ErrorCode)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipReportReqPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)
	s := make([]byte, 12)
	r.ReadBytes(s)
	p.SubmitSequenceNum = [3]uint32{
		binary.BigEndian.Uint32(s[:4]),
		binary.BigEndian.Uint32(s[4:8]),
		binary.BigEndian.Uint32(s[8:]),
	}
	p.ReportType = r.ReadByte()
	p.UserNumber = string(r.ReadCString(21))
	p.State = r.ReadByte()
	p.ErrorCode = r.ReadByte()
	p.Reserve = string(r.ReadCString(8))

	p.SubmitSequenceNumStr = GenSequenceNumStr(p.SubmitSequenceNum)
	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)
	return r.Error()
}

func (p *SgipReportReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Report Req ---")
	fmt.Fprintln(&b, "SubmitSequenceNumStr", p.SubmitSequenceNumStr)
	fmt.Fprintln(&b, "ReportType: ", p.ReportType)
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "State: ", p.State)
	fmt.Fprintln(&b, "ErrorCode: ", p.ErrorCode)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)

	return b.String()
}

type SgipReportRespPkt struct {
	SgipRespPkt
}

func (p *SgipReportRespPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	return p.SgipRespPkt.Pack(seqNum, SGIP_REPORT_RESP)
}

func (p *SgipReportRespPkt) Unpack(data []byte) error {
	return p.SgipRespPkt.Unpack(data)
}

func (p *SgipReportRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Report Resp ---")
	fmt.Fprintln(&b, p.SgipRespPkt.String())

	return b.String()
}
