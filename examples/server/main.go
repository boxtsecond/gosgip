package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/boxtsecond/gosgip/pkg"
	"github.com/boxtsecond/gosgip/server"
)

const (
	user     string = "10000001"
	password string = "12345678"
	nodeId   uint32 = 123456
)

func handleBind(r *server.Response, p *server.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*pkg.SgipBindReqPkt)
	if !ok {
		return true, nil
	}

	l.Println("remote addr:", p.Conn.Conn.RemoteAddr().(*net.TCPAddr).IP.String())
	resp := r.Packer.(*pkg.SgipBindRespPkt)

	if req.LoginName != user || req.LoginPassword != password {
		resp.Result = pkg.Status(1)
		l.Println("handleBind error:", resp.Result.Error())
		return false, resp.Result.Error()
	}

	if req.LoginType != pkg.BIND_CLIENT {
		resp.Result = pkg.Status(4)
		l.Println("handleBind error:", resp.Result.Error())
		return false, resp.Result.Error()
	}

	l.Printf("handleBind: %s login ok\n", req.LoginName)

	return false, nil
}

func handleSubmit(r *server.Response, p *server.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*pkg.SgipSubmitReqPkt)
	if !ok {
		return true, nil
	}

	//log.Printf("server sgip: receive a sgip submit request ok: %v.", req)

	resp := r.Packer.(*pkg.SgipSubmitRespPkt)
	resp.SequenceNum = pkg.GenSequenceNum(nodeId, <-p.Conn.SequenceID)
	deliverPkgs := make([]*pkg.SgipReportReqPkt, 0)
	userRptPkgs := make([]*pkg.SgipUserRptReqPkt, 0)
	for i, d := range req.UserNumber {
		l.Printf("handleSubmit: handle submit from %s ok! msgid[%s], destTerminalId[%s]\n",
			req.SPNumber, fmt.Sprintf("[%s]_%d", req.SequenceNumStr, i), d)

		seqNum := pkg.GenSequenceNum(nodeId, <-p.Conn.SequenceID)
		userRptSeqNum := pkg.GenSequenceNum(nodeId, <-p.Conn.SequenceID)

		deliverPkgs = append(deliverPkgs, &pkg.SgipReportReqPkt{
			SubmitSequenceNum: req.SequenceNum,
			ReportType:        1,
			UserNumber:        d,
			State:             0,
			ErrorCode:         0,
			Reserve:           "",
			SequenceNum:       seqNum,
		})
		userRptPkgs = append(userRptPkgs, &pkg.SgipUserRptReqPkt{
			SPNumber:      req.SPNumber,
			UserNumber:    d,
			UserCondition: 0,
			Reserve:       "",
			SequenceNum:   userRptSeqNum,
		})
	}
	go mockReport(deliverPkgs, p)
	go mockUserRpt(userRptPkgs, p)
	return false, nil
}

func mockReport(pkgs []*pkg.SgipReportReqPkt, s *server.Packet) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:

			for _, p := range pkgs {
				err := s.SendPkt(p, p.SequenceNum)
				if err != nil {
					log.Printf("server sgip: send a sgip report request error: %s.", err)
					return
				} else {
					log.Printf("server sgip: send a sgip report request ok.")
				}
			}
			return

		default:
		}

	}
}

func mockUserRpt(pkgs []*pkg.SgipUserRptReqPkt, s *server.Packet) {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:

			for _, p := range pkgs {
				err := s.SendPkt(p, p.SequenceNum)
				if err != nil {
					log.Printf("server sgip: send a sgip userRpt request error: %s.", err)
					return
				} else {
					log.Printf("server sgip: send a sgip userRpt request ok")
				}
			}
			return

		default:
		}

	}
}

func main() {
	var handlers = []server.Handler{
		server.HandlerFunc(handleBind),
		server.HandlerFunc(handleSubmit),
	}

	err := server.ListenAndServe(":8810",
		nodeId,
		5*time.Second,
		3,
		nil,
		handlers...,
	)
	if err != nil {
		log.Println("sgip Listen And Server error:", err)
	}
	return
}
