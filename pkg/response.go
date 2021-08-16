package pkg

import (
	"bytes"
	"fmt"
)

type SgipRespPkt struct {
	Result  Status // 请求返回结果，0：执行成功 其它：错误码
	Reserve string // 保留

	// used in session
	SequenceNum string
}

func (p *SgipRespPkt) Pack(seqNum string) ([]byte, error) {
	var w = newPkgWriter(SgipBindRespPktLen)
	// header
	w.WriteHeader(SgipBindRespPktLen, seqNum, SGIP_BIND_RESP)
	p.SequenceNum = seqNum

	// body
	w.WriteByte(p.Result.Data())
	w.WriteFixedSizeString(p.Reserve, 8)

	return w.Bytes()
}

func (p *SgipRespPkt) Unpack(data []byte) error {
	var r = newPkgReader(data)

	// Body: Result
	p.Result = Status(r.ReadByte())
	// Body: Reserve
	p.Reserve = string(r.ReadCString(8))

	return r.Error()
}

func (p *SgipRespPkt) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- SGIP Resp ---")
	fmt.Fprintln(&b, "Result: ", p.Result)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	return b.String()
}
