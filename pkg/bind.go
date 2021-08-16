package pkg

import (
	"bytes"
	"fmt"
)

const (
	SgipBindReqPktLen  = HeaderPktLen + 1 + 16 + 16 + 8 //61d, 0x3d
	SgipBindRespPktLen = HeaderPktLen + 1 + 8           //29d, 0x1d
)

type SgipBindReqPkt struct {
	LoginType     uint8  // 登录类型 1 sp -> SMG, 2 SMG -> SP
	LoginName     string // 登陆名
	LoginPassword string // 登陆密码
	Reserve       string // 保留

	// used in session
	SequenceNum string
}

func (p *SgipBindReqPkt) Pack(seqNum string) ([]byte, error) {
	var w = newPkgWriter(SgipBindReqPktLen)
	// header
	w.WriteHeader(SgipBindReqPktLen, seqNum, SGIP_BIND)
	p.SequenceNum = seqNum

	// body
	w.WriteByte(p.LoginType)
	w.WriteFixedSizeString(p.LoginName, 16)
	w.WriteFixedSizeString(p.LoginPassword, 16)
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipBindReqPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)

	// Body: LoginType
	p.LoginType = r.ReadByte()
	// Body: LoginName
	p.LoginName = string(r.ReadCString(16))
	// Body: LoginPassword
	p.LoginPassword = string(r.ReadCString(16))
	// Body: Reserve
	p.Reserve = string(r.ReadCString(8))

	return r.Error()
}

func (p *SgipBindReqPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Bind Req ---")
	fmt.Fprintln(&b, "LoginType: ", p.LoginType)
	fmt.Fprintln(&b, "LoginName: ", p.LoginName)
	fmt.Fprintln(&b, "LoginPassword: ", p.LoginPassword)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	return b.String()
}
