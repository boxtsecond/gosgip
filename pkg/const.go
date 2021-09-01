package pkg

const (
	BIND_CLIENT = 1
	BIND_SERVER = 2
)

// MsgType
const (
	MSG = 0 // 短信息信息
)

// MsgFormat
// 短消息内容体的编码格式
// 对于文字短消息，要求MsgFormat＝15, 对于回执消息，要求MsgFormat＝0
const (
	ASCII   = 0  // ASCII编码
	BINARY  = 4  // 二进制短消息
	UCS2    = 8  // UCS2编码
	GB18030 = 15 // GB18030编码
)

const (
	SUBMIT_REPORT  = 0 // 不是状态报告
	DELIVER_REPORT = 1 // 是状态报告
)

// 是否要求返回状态报告
const (
	FAILED_REPORT  = 0 // 只有最后出错时要返回状态报告
	NEED_REPORT    = 1 // 无论最后是否成功都要返回状态报告
	NO_NEED_REPORT = 2
)

// 短消息发送优先级
const (
	DEFAULT_PRIORITY = 0
)

// Report 所涉及的短消息的当前执行状态
const (
	SUCCESS = 0 // 发送成功
	WAITING = 1 // 等待发送
	FAILED  = 2 // 发送失败
)
