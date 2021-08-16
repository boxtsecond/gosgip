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
	SequenceNum string
}
type SgipUnbindRespPkt struct {
	// used in session
	SequenceNum string
}

func (p *SgipUnbindReqPkt) Pack(seqNum string) ([]byte, error) {
	var w = newPkgWriter(SgipUnbindReqPktLen)

	// header
	w.WriteHeader(SgipUnbindReqPktLen, seqNum, SGIP_UNBIND)
	p.SequenceNum = seqNum

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

func (p *SgipUnbindRespPkt) Pack(seqNum string) ([]byte, error) {
	var w = newPkgWriter(SgipUnbindRespPktLen)

	// header
	w.WriteHeader(SgipUnbindRespPktLen, seqNum, SGIP_UNBIND_RESP)
	p.SequenceNum = seqNum

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
