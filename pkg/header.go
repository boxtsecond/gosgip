package pkg

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const HeaderPktLen uint32 = 4 + 4 + 12

// 消息头(所有消息公共包头)
type Header struct {
	PacketLength uint32 // 数据包长度
	CommandID    uint32 // 请求标识
	SequenceNum  string // 消息流水号
}

func (p *Header) Pack(w *pkgWriter, pktLen, commandId uint32, seqNum string) *pkgWriter {
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, commandId)
	msgId, _ := hex.DecodeString(seqNum)
	w.WriteBytes(NewOctetString(fmt.Sprintf("%s", msgId)).Byte(12))
	return w
}

func (p *Header) Unpack(r *pkgReader) *Header {
	r.ReadInt(binary.BigEndian, &p.PacketLength)
	r.ReadInt(binary.BigEndian, &p.CommandID)
	var s = make([]byte, 12)
	r.ReadBytes(s)
	p.SequenceNum = hex.EncodeToString(s)
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
