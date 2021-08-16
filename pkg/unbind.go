package pkg

import (
	"bytes"
	"fmt"
)

const (
	SgipUnbindReqPktLen  = HeaderPktLen
	SgipUnbindRespPktLen = HeaderPktLen
)

type SgipUnbindReqPkt struct {
	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}
type SgipUnbindRespPkt struct {
	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}

func (p *SgipUnbindReqPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var w = newPkgWriter(SgipUnbindReqPktLen)

	// header
	w.WriteHeader(SgipUnbindReqPktLen, seqNum, SGIP_UNBIND)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

	return w.Bytes()
}

func (p *SgipUnbindReqPkt) Unpack(data []byte) error {
	return nil
}

func (p *SgipUnbindReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Unbind Req ---")
	return b.String()
}

func (p *SgipUnbindRespPkt) Pack(seqNum [3]uint32) ([]byte, error) {
	var w = newPkgWriter(SgipUnbindRespPktLen)

	// header
	w.WriteHeader(SgipUnbindRespPktLen, seqNum, SGIP_UNBIND_RESP)

	return w.Bytes()
}

func (p *SgipUnbindRespPkt) Unpack(data []byte) error {
	return nil
}

func (p *SgipUnbindRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Unbind Resp ---")
	return b.String()
}
