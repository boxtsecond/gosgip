package pkg

import (
	"bytes"
	"fmt"
)

const (
	SgipUserRptReqPktLen = HeaderPktLen + 21 + 21 + 1 + 8 //71d, 0x47
)

type SgipUserRptReqPkt struct {
	SPNumber      string
	UserNumber    string
	UserCondition uint8
	Reserve       string

	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}

func (p *SgipUserRptReqPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var w = newPkgWriter(SgipUserRptReqPktLen)
	// header
	w.WriteHeader(SgipUserRptReqPktLen, seqNum, SGIP_USERRPT)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	// body
	w.WriteFixedSizeString(p.SPNumber, 21)
	w.WriteFixedSizeString(p.UserNumber, 21)
	w.WriteByte(p.UserCondition)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipUserRptReqPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)

	p.SPNumber = string(r.ReadCString(21))
	p.UserNumber = string(r.ReadCString(21))
	p.UserCondition = r.ReadByte()
	p.Reserve = string(r.ReadCString(8))

	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)

	return r.Error()
}

func (p *SgipUserRptReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP UserRpt Req ---")
	fmt.Fprintln(&b, "SPNumber: ", p.SPNumber)
	fmt.Fprintln(&b, "UserNumber: ", p.UserNumber)
	fmt.Fprintln(&b, "UserCondition: ", p.UserCondition)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)

	return b.String()
}

type SgipUserRptRespPkt struct {
	SgipRespPkt
}

func (p *SgipUserRptRespPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	return p.SgipRespPkt.Pack(seqNum, SGIP_USERRPT_RESP)
}

func (p *SgipUserRptRespPkt) Unpack(data []byte) error {
	return p.SgipRespPkt.Unpack(data)
}

func (p *SgipUserRptRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP UserRpt Resp ---")
	fmt.Fprintln(&b, p.SgipRespPkt.String())

	return b.String()
}
