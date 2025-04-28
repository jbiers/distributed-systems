package kvsrv

import "fmt"

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	ID    RequestID
	// You'll have to add definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
}

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	Key string
	ID  RequestID
}

type GetReply struct {
	Value string
}

type RequestID struct {
	ClientID int64
	RPCCount int
}

func (r *RequestID) GetString() string {
	return fmt.Sprintf("%d-%d", r.ClientID, r.RPCCount)
}
