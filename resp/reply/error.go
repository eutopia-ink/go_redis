package reply

// 未知错误
type UnknowErrReply struct {
}

var unkownErrBytes = []byte("-Err unknown\r\n")

func (u *UnknowErrReply) Error() string {
	return "Err unknown"
}

func (u *UnknowErrReply) ToBytes() []byte {
	return unkownErrBytes
}

func MakeUnknowErrReply() *UnknowErrReply {
	return &UnknowErrReply{}
}

// 参数个数错误
type ArgNumErrReply struct {
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return "-ERR wrong number of arguments for '" + a.Cmd + "' command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// 语法错误
type SyntaxErrReply struct {
}

var SyntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

func (u *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (u *SyntaxErrReply) ToBytes() []byte {
	return SyntaxErrBytes
}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return &SyntaxErrReply{}
}

// 数据类型错误
type WrongTypeReply struct {
}

var WrongTypeBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (u *WrongTypeReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (u *WrongTypeReply) ToBytes() []byte {
	return WrongTypeBytes
}

func MakeWrongTypeReply() *WrongTypeReply {
	return &WrongTypeReply{}
}

// 协议错误
type ProtocolErrReply struct {
	Msg string
}

func (u *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + u.Msg + "'"
}

func (u *ProtocolErrReply) ToBytes() []byte {
	return []byte("ERR Protocol error: '" + u.Msg + "'\r\n")
}

func MakeProtocolErrReply(msg string) *ProtocolErrReply {
	return &ProtocolErrReply{
		Msg: msg,
	}
}
