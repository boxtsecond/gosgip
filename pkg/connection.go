package pkg

import (
	"encoding/binary"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"
)

type State uint8

const (
	CONNECTION_CLOSED State = iota
	CONNECTION_CONNECTED
	CONNECTION_AUTHOK
)

type Conn struct {
	net.Conn
	State State

	// for SequenceNum generator goroutine
	SequenceID <-chan uint32
	done       chan<- struct{}
}

func newRandomSequenceIDGenerator() (<-chan uint32, chan<- struct{}) {
	out := make(chan uint32)
	done := make(chan struct{})
	rand.Seed(time.Now().UnixNano())

	go func() {
		var i = uint32(rand.Intn(100000))

		for {
			select {
			case out <- i:
				i++
			case <-done:
				close(out)
				return
			}
		}
	}()
	return out, done
}

func newSequenceIDGenerator() (<-chan uint32, chan<- struct{}) {
	out := make(chan uint32)
	done := make(chan struct{})

	go func() {
		var i uint32
		for {
			select {
			case out <- i:
				i++
			case <-done:
				close(out)
				return
			}
		}
	}()
	return out, done
}

func NewConnection(conn net.Conn) *Conn {
	sequenceID, done := newRandomSequenceIDGenerator()
	c := &Conn{
		Conn:       conn,
		SequenceID: sequenceID,
		done:       done,
	}
	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true) //Keepalive as default
	return c
}

func (c *Conn) Close() {
	if c != nil {
		if c.State == CONNECTION_CLOSED {
			return
		}
		close(c.done)  // let the SeqId goroutine exit.
		c.Conn.Close() // close the underlying net.Conn
		c.State = CONNECTION_CLOSED
	}
}

func (c *Conn) SetState(state State) {
	c.State = state
}

func (c *Conn) SendPkt(packet Packer, seqNum [3]uint32) error {
	if c.State == CONNECTION_CLOSED {
		return ErrConnIsClosed
	}

	data, err := packet.Pack(seqNum)
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(data) //block write
	if err != nil {
		return err
	}

	return nil
}

const (
	defaultReadBufferSize = 4096
)

type CommandIDHeader struct {
	PacketLength uint32   // 数据包长度
	CommandID    uint32   // 请求标识
	SequenceID   [12]byte // 请求标识
}

type readBuffer struct {
	CommandIDHeader
	leftData [defaultReadBufferSize]byte
}

var readBufferPool = sync.Pool{
	New: func() interface{} {
		return &readBuffer{}
	},
}

func (c *Conn) RecvAndUnpackPkt(timeout time.Duration) (Packer, error) {
	if c.State == CONNECTION_CLOSED {
		return nil, ErrConnIsClosed
	}
	rb := readBufferPool.Get().(*readBuffer)
	defer func() {
		readBufferPool.Put(rb)
		c.SetReadDeadline(time.Time{})
	}()

	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}

	// packet header
	err := binary.Read(c.Conn, binary.BigEndian, &rb.CommandIDHeader)
	if err != nil {
		return nil, err
	}

	if rb.CommandIDHeader.PacketLength < SGIP_PACKET_MIN || rb.CommandIDHeader.PacketLength > SGIP_PACKET_MAX {
		return nil, ErrTotalLengthInvalid
	}

	if !((CommandID(rb.CommandIDHeader.CommandID) > SGIP_REQUEST_MIN && CommandID(rb.CommandIDHeader.CommandID) < SGIP_REQUEST_MAX) ||
		(CommandID(rb.CommandIDHeader.CommandID) > SGIP_RESPONSE_MIN && CommandID(rb.CommandIDHeader.CommandID) < SGIP_RESPONSE_MAX)) {
		return nil, ErrCommandIDInvalid
	}

	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}

	sequenceID := UnpackSequenceNum(rb.CommandIDHeader.SequenceID)

	// packet body
	var leftData = rb.leftData[0:(rb.CommandIDHeader.PacketLength - 20)]

	_, err = io.ReadFull(c.Conn, leftData)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok {
			if netErr.Timeout() {
				return nil, ErrReadPktBodyTimeout
			}
		}
		return nil, err
	}

	var p Packer
	switch CommandID(rb.CommandIDHeader.CommandID) {
	case SGIP_BIND:
		p = &SgipBindReqPkt{SequenceNum: sequenceID}
	case SGIP_BIND_RESP:
		p = &SgipBindRespPkt{SgipRespPkt{SequenceNum: sequenceID}}
	case SGIP_SUBMIT:
		p = &SgipSubmitReqPkt{SequenceNum: sequenceID}
	case SGIP_SUBMIT_RESP:
		p = &SgipSubmitRespPkt{SgipRespPkt{SequenceNum: sequenceID}}
	case SGIP_DELIVER:
		p = &SgipDeliverReqPkt{SequenceNum: sequenceID}
	case SGIP_DELIVER_RESP:
		p = &SgipDeliverRespPkt{SgipRespPkt{SequenceNum: sequenceID}}
	case SGIP_REPORT:
		p = &SgipReportReqPkt{SequenceNum: sequenceID}
	case SGIP_REPORT_RESP:
		p = &SgipReportRespPkt{SgipRespPkt{SequenceNum: sequenceID}}
	case SGIP_UNBIND:
		p = &SgipUnbindReqPkt{SequenceNum: sequenceID}
	case SGIP_UNBIND_RESP:
		p = &SgipUnbindRespPkt{SequenceNum: sequenceID}
	case SGIP_USERRPT:
		p = &SgipUserRptReqPkt{SequenceNum: sequenceID}
	case SGIP_USERRPT_RESP:
		p = &SgipUserRptRespPkt{SgipRespPkt{SequenceNum: sequenceID}}
	case SGIP_TRACE:
		p = &SgipTraceReqPkt{SequenceNum: sequenceID}
	case SGIP_TRACE_RESP:
		p = &SgipTraceRespPkt{SequenceNum: sequenceID}

	default:
		return nil, ErrCommandIDNotSupported
	}

	err = p.Unpack(leftData)
	if err != nil {
		return nil, err
	}
	return p, nil
}
