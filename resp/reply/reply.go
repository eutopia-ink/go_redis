package reply

import (
	"Go_redis/interface/resp"
	"bytes"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// 字符串回复
type BulkReply struct {
	Arg []byte
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + string(b.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

// 多个字符串回复
type MultiBulkReply struct {
	Args [][]byte
}

func (r *MultiBulkReply) ToBytes() []byte {
	argLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

// 状态回复
type StatusReply struct {
	Status string
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

// 数字回复
type IntReply struct {
	Code int64
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(r.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

// 通用错误回复
type StandardErrReply struct {
	Status string
}

func (r *StandardErrReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

func (r *StandardErrReply) Error() string {
	return r.Status
}

func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
