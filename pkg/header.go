package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const HeaderPktLen uint32 = 4 + 4 + 12

// 消息头(所有消息公共包头)
type Header struct {
	PacketLength uint32    // 数据包长度
	CommandID    uint32    // 请求标识
	SequenceNum  [3]uint32 // 消息流水号
}

func (p *Header) Pack(w *pkgWriter, pktLen, commandId uint32, seqNum [3]uint32) *pkgWriter {
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, commandId)
	w.WriteInt(binary.BigEndian, seqNum[0])
	w.WriteInt(binary.BigEndian, seqNum[1])
	w.WriteInt(binary.BigEndian, seqNum[2])
	return w
}

func (p *Header) Unpack(r *pkgReader) *Header {
	r.ReadInt(binary.BigEndian, &p.PacketLength)
	r.ReadInt(binary.BigEndian, &p.CommandID)
	var s = make([]byte, 12)
	r.ReadBytes(s)
	p.SequenceNum = [3]uint32{
		binary.BigEndian.Uint32(s[:4]),
		binary.BigEndian.Uint32(s[4:8]),
		binary.BigEndian.Uint32(s[8:]),
	}
	return p
}

func (p *Header) String() string {
	var b bytes.Buffer
	fmt.Fprintln(&b, "--- Header ---")
	fmt.Fprintln(&b, "Length: ", p.PacketLength)
	fmt.Fprintf(&b, "CommandID: 0x%x\n", p.CommandID)
	fmt.Fprintln(&b, "SequenceNum: ", p.SequenceNum)

	return b.String()

}
