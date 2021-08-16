package client

import (
	"errors"
	"net"
	"time"

	"github.com/boxtsecond/gosgip/pkg"
)

var ErrRespNotMatch = errors.New("the response is not matched with the request")

type Client struct {
	conn *pkg.Conn

	nodeId uint32
	CorpId string
}

func NewClient(nodeId uint32, corpId string) *Client {
	//nodeId, err := pkg.GenNodeId(areaCode, corpId)
	//if err != nil {
	//	return nil, err
	//}
	return &Client{
		nodeId: nodeId,
		CorpId: corpId,
	}
}

func (cli *Client) Connect(serverAddr, user, pwd string, loginType uint8, timeout time.Duration) error {
	var err error
	conn, err := net.DialTimeout("tcp", serverAddr, timeout)
	if err != nil {
		return err
	}
	cli.conn = pkg.NewConnection(conn)
	defer func() {
		if err != nil {
			if cli.conn != nil {
				cli.conn.Close()
			}
		}
	}()

	cli.conn.SetState(pkg.CONNECTION_CONNECTED)

	// Login to the server.
	if err != nil {
		return err
	}
	req := &pkg.SgipBindReqPkt{
		LoginType:     loginType,
		LoginName:     user,
		LoginPassword: pwd,
		Reserve:       "",
	}

	err = cli.SendReqPkt(req)
	if err != nil {
		return err
	}

	p, err := cli.conn.RecvAndUnpackPkt(timeout)
	if err != nil {
		return err
	}

	rsp, ok := p.(*pkg.SgipRespPkt)
	if !ok {
		err = ErrRespNotMatch
		return err
	}

	if rsp.Result.Data() != 0 {
		return rsp.Result.Error()
	}

	cli.conn.SetState(pkg.CONNECTION_AUTHOK)
	return nil
}

func (cli *Client) Disconnect() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}

func (cli *Client) SendReqPkt(packet pkg.Packer) error {
	seqNum, err := pkg.GenSequenceNum(cli.nodeId, <-cli.conn.SequenceID)

	if err != nil {
		return err
	}
	return cli.conn.SendPkt(packet, seqNum)
}

func (cli *Client) SendRspPkt(packet pkg.Packer, seqNum string) error {
	return cli.conn.SendPkt(packet, seqNum)
}

func (cli *Client) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	return cli.conn.RecvAndUnpackPkt(timeout)
}
