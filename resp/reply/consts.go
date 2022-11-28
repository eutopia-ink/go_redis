package reply

// 返回PONG
type PongReply struct {
}

var pongbytes = []byte("+PONG\n\n")

func (r *PongReply) ToBytes() []byte {
	return pongbytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

// 返回OK
type OkReply struct {
}

var okbytes = []byte("+OK\n\n")

func (r *OkReply) ToBytes() []byte {
	return okbytes
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}

// 无字符回复
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n\r\n")

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// 空数组回复
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// 空回复
type NoReply struct {
}

var noBytes = []byte("*0\r\n")

func (r *NoReply) ToBytes() []byte {
	return noBytes
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}
