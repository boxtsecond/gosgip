package pkg

import (
	"bytes"
	"fmt"
)

const (
	SgipRespPktLen = HeaderPktLen + 1 + 8 //29d, 0x1d
)

type SgipRespPkt struct {
	Result  Status // 请求返回结果，0：执行成功 其它：错误码
	Reserve string // 保留

	// used in session
	SequenceNum    [3]uint32
	SequenceNumStr string
}

func (p *SgipRespPkt) Pack(seqNum [3]uint32, commandId CommandID) ([]byte, error) {
	var w = newPkgWriter(SgipRespPktLen)
	// header
	w.WriteHeader(SgipRespPktLen, seqNum, commandId)
	p.SequenceNumStr = GenSequenceNumStr(seqNum)

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

	p.SequenceNumStr = GenSequenceNumStr(p.SequenceNum)

	return r.Error()
}

func (p *SgipRespPkt) String() string {
	var b bytes.Buffer
	//fmt.Fprintln(&b, "--- SGIP Resp ---")
	fmt.Fprintln(&b, "Result: ", p.Result)
	fmt.Fprintln(&b, "Reserve: ", p.Reserve)
	fmt.Fprintln(&b, "SequenceNumStr", p.SequenceNumStr)
	return b.String()
}
