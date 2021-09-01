package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/boxtsecond/gosgip/client"
	"github.com/boxtsecond/gosgip/pkg"
)

var (
	addr      = flag.String("addr", ":8810", "sgip addr(运营商地址)")
	user      = flag.String("user", "10000001", "登陆账号")
	pwd       = flag.String("pwd", "12345678", "登陆密码")
	loginType = flag.String("loginType", "1", "登陆类型")
	nodeId    = flag.String("nodeId", "123456", "节点编号")
	spCode    = flag.String("spCode", "12345", "SP的接入号码")
	corpId    = flag.String("corpId", "12345", "企业代码")
	phone     = flag.String("phone", "8618012345678", "接收手机号码, 86..., 多个使用,分割")
	msg       = flag.String("msg", "验证码：1234", "短信内容")
)

func startAClient(idx int) {
	nodeIdInt, _ := strconv.ParseUint(*nodeId, 10, 32)
	c := client.NewClient(uint32(nodeIdInt), *spCode)
	defer wg.Done()
	defer c.Disconnect()

	login, _ := strconv.Atoi(*loginType)
	err := c.Connect(*addr, *user, *pwd, uint8(login), 3*time.Second)
	if err != nil {
		log.Printf("client %d: connect error: %s.", idx, err)
		return
	}
	log.Printf("client %d: connect and auth ok", idx)

	t := time.NewTicker(time.Second * 3)
	defer t.Stop()

	maxSubmit := 5
	count := 0
	for {
		select {
		case <-t.C:
			if count >= maxSubmit {
				break
			}

			cont, err := pkg.Utf8ToGB18030(*msg)
			if err != nil {
				log.Printf("client %d: utf8 to gb18030 transform err: %s.", idx, err)
				return
			}
			destStrArr := strings.Split(*phone, ",")
			p := &pkg.SgipSubmitReqPkt{
				SPNumber:         *spCode,
				ChargeNumber:     "",
				UserCount:        uint8(len(destStrArr)),
				UserNumber:       destStrArr,
				CorpId:           *corpId,
				ServiceType:      "",
				FeeType:          0,
				FeeValue:         "",
				GivenValue:       "",
				AgentFlag:        0,
				MorelatetoMTFlag: 0,
				Priority:         pkg.DEFAULT_PRIORITY,
				ExpireTime:       "",
				ScheduleTime:     "",
				ReportFlag:       pkg.NEED_REPORT,
				TP_pid:           0,
				TP_udhi:          0,
				MessageCoding:    pkg.GB18030,
				MessageType:      pkg.MSG,
				MessageLength:    uint32(len(cont)),
				MessageContent:   cont,
				Reserve:          "",
			}

			_, err = c.SendReqPkt(p)
			if err != nil {
				log.Printf("client %d: send a sgip submit request error: %s.", idx, err)
				return
			} else {
				log.Printf("client %d: send a sgip submit request ok", idx)
			}
			count += 1

			go func() {
				tracePkg := &pkg.SgipTraceReqPkt{
					SubmitSequenceNum: p.SequenceNum,
					UserNumber:        p.UserNumber[0],
					Reserve:           "",
				}
				_, err := c.SendReqPkt(tracePkg)
				if err != nil {
					log.Printf("client %d: send a sgip trace request error: %s.", idx, err)
					return
				} else {
					log.Printf("client %d: send a sgip trace request ok", idx)
				}
			}()
			//default:
		}

		// recv packets
		i, err := c.RecvAndUnpackPkt(0)
		if err != nil {
			log.Printf("client %d: client read and unpack pkt error: %s.", idx, err)
			break
		}

		switch p := i.(type) {
		case *pkg.SgipSubmitRespPkt:
			log.Printf("client %d: receive a sgip submit response.", idx)

		case *pkg.SgipBindRespPkt:
			log.Printf("client %d: receive a sgip bind response.", idx)

		case *pkg.SgipDeliverReqPkt:
			log.Printf("client %d: receive a sgip deliver request ok.", idx)
			rsp := &pkg.SgipDeliverRespPkt{SgipRespPkt: pkg.SgipRespPkt{Result: pkg.Status(0)}}
			err := c.SendRspPkt(rsp, p.SequenceNum)
			if err != nil {
				log.Printf("client %d: send sgip deliver response error: %s.", idx, err)
				break
			} else {
				log.Printf("client %d: send sgip deliver response ok.", idx)
			}

		case *pkg.SgipReportReqPkt:
			log.Printf("client %d: receive a sgip report request ok.", idx)
			rsp := &pkg.SgipReportRespPkt{SgipRespPkt: pkg.SgipRespPkt{Result: pkg.Status(0)}}
			err := c.SendRspPkt(rsp, p.SequenceNum)
			if err != nil {
				log.Printf("client %d: send sgip report response error: %s.", idx, err)
				break
			} else {
				log.Printf("client %d: send sgip report response ok.", idx)
			}

		case *pkg.SgipUnbindReqPkt:
			log.Printf("client %d: receive a sgip exit request.", idx)
			rsp := &pkg.SgipUnbindRespPkt{}
			err := c.SendRspPkt(rsp, p.SequenceNum)
			if err != nil {
				log.Printf("client %d: send sgip exit response error: %s.", idx, err)
				break
			}

		case *pkg.SgipUnbindRespPkt:
			log.Printf("client %d: receive a sgip exit response.", idx)

		case *pkg.SgipUserRptReqPkt:
			log.Printf("client %d: receive a sgip userRpt request.", idx)
			rsp := &pkg.SgipUserRptRespPkt{}
			err := c.SendRspPkt(rsp, p.SequenceNum)
			if err != nil {
				log.Printf("client %d: send sgip userRpt response error: %s.", idx, err)
				break
			}

		case *pkg.SgipUserRptRespPkt:
			log.Printf("client %d: receive a sgip userRpt response.", idx)

		case *pkg.SgipTraceReqPkt:
			log.Printf("client %d: receive a sgip trace request.", idx)
			rsp := &pkg.SgipUnbindRespPkt{}
			err := c.SendRspPkt(rsp, p.SequenceNum)
			if err != nil {
				log.Printf("client %d: send sgip trace response error: %s.", idx, err)
				break
			}

		case *pkg.SgipTraceRespPkt:
			log.Printf("client %d: receive a sgip trace response.", idx)
		}
	}
}

var wg sync.WaitGroup

func init() {
	flag.Parse()
}

func main() {
	log.Println("Client example start!")
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go startAClient(i + 1)
	}
	wg.Wait()
	log.Println("Client example ends!")
}
