package kvsrv

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"6.5840/labrpc"
)

type Clerk struct {
	server   *labrpc.ClientEnd
	clientID int64
	rpcCount int
}

func getRPCID(client int64, rpc int) string {
	return fmt.Sprintf("%d-%d", client, rpc)
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	ck.clientID = nrand()
	ck.rpcCount = 0

	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	args := GetArgs{
		Key: key,
		ID: RequestID{
			ClientID: ck.clientID,
			RPCCount: ck.rpcCount,
		},
	}
	reply := GetReply{}
	ck.rpcCount++

	var ok bool
	for {
		ok = ck.server.Call("KVServer.Get", &args, &reply)
		if ok {
			break
		}
	}

	return reply.Value
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer."+op, &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) string {
	args := PutAppendArgs{
		Key:   key,
		Value: value,
		ID: RequestID{
			ClientID: ck.clientID,
			RPCCount: ck.rpcCount,
		},
	}
	reply := PutAppendReply{}
	ck.rpcCount++

	var ok bool
	for {
		ok = ck.server.Call("KVServer."+op, &args, &reply)
		if ok {
			break
		}
	}

	return reply.Value
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
