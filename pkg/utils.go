package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func GenTimestamp() uint32 {
	s := time.Now().Format("0102150405")
	i, _ := strconv.Atoi(s)
	return uint32(i)
}

func GenNowTimeYYYYStr() string {
	s := time.Now().Format("20060102150405")
	return s
}

func GenNowTimeNoYearStr() string {
	return time.Unix(time.Now().Unix(), 0).Format("0102150405")
}

// 节点编号
// 通信节点编号规则
// 在整个网关系统中，所有的通信节点(SMG、GNS、SP和SMSC)都有一个唯一的数字编号，不同的SP或SMSC或SMG或GNS编号不能相同，编号由系统管理人员负责分配。编号规则如下：
// SMG的编号规则：1AAAAX
// SMSC的编号规则：	2AAAAX
// SP的编号规则：3AAAAQQQQQ
// GNS的编号规则：4AAAAX
// 其中, AAAA表示四位长途区号(不足四位的长途区号，左对齐，右补零),X表示1位序号,QQQQQ表示5位企业代码。
func GenNodeId(areaCode, corpId string) (uint32, error) {
	var (
		err error
		ac  int
		ci  int
	)

	// check arg
	if ac, err = strconv.Atoi(areaCode); err != nil {
		return 0, err
	}
	if ci, err = strconv.Atoi(corpId); err != nil {
		return 0, err
	}

	// 0XX 三位区号
	if ac < 100 {
		return uint32(3000000000 + ac*1000000 + ci), nil
	}
	// 0XXX 四位区号
	return uint32(3000000000 + ac*100000 + ci), nil
}

//SequenceNum 字段包含以下三部分内容：
//命令源节点的编号：4字节
//时间：4字节，格式为MMDDHHMMSS（月日时分秒）
//序列号：4字节
//返回16进制字符串，长度为：12*2
func GenSequenceNum(nodeId, sequenceId uint32) (string, error) {
	timeStr := GenNowTimeNoYearStr()
	timeInt, _ := strconv.ParseInt(timeStr, 10, 32)
	binaryStr := fmt.Sprintf("%032b%032b%032b", nodeId, timeInt, sequenceId)
	head, _ := strconv.ParseUint(binaryStr[:64], 2, 64)
	end, _ := strconv.ParseUint(binaryStr[64:], 2, 32)
	return fmt.Sprintf("%016x%08x", head, end), nil
}

func UnpackSequenceNum(sequenceNum string) string {
	spId, _ := strconv.ParseUint(sequenceNum[:32], 16, 32)
	timeStr, _ := strconv.ParseUint(sequenceNum[32:64], 16, 32)
	seqId, _ := strconv.ParseUint(sequenceNum[64:], 16, 32)
	return fmt.Sprintf("spId: %s, time: %s, seqId: %d, ", NewOctetString(strconv.Itoa(int(spId))).FixedString(10), NewOctetString(strconv.Itoa(int(timeStr))).FixedString(10), seqId)
}

func Utf8ToUcs2(in string) (string, error) {
	if !utf8.ValidString(in) {
		return "", errors.New("invalid utf8 runes")
	}

	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()) //UTF-16 bigendian, no-bom
	out, err := ioutil.ReadAll(t)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Ucs2ToUtf8(in string) (string, error) {
	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()) //UTF-16 bigendian, no-bom
	out, err := ioutil.ReadAll(t)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Utf8ToGB18030(in string) (string, error) {
	if !utf8.ValidString(in) {
		return "", errors.New("invalid utf8 runes")
	}

	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	out, err := ioutil.ReadAll(t)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func GB18030ToUtf8(in string) (string, error) {
	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	out, err := ioutil.ReadAll(t)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
