package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex

	data       map[string]string
	requestIDs map[string]string
}

func (kv *KVServer) deleteEntry(r RequestID) {
	delete(kv.requestIDs, r.GetString())
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.deleteEntry(RequestID{
		ClientID: args.ID.ClientID,
		RPCCount: args.ID.RPCCount - 1,
	})

	value := kv.data[args.Key]
	reply.Value = value
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.deleteEntry(RequestID{
		ClientID: args.ID.ClientID,
		RPCCount: args.ID.RPCCount - 1,
	})

	kv.data[args.Key] = args.Value
	reply.Value = args.Value
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.deleteEntry(RequestID{
		ClientID: args.ID.ClientID,
		RPCCount: args.ID.RPCCount - 1,
	})

	if value, exists := kv.requestIDs[args.ID.GetString()]; exists {
		reply.Value = value
		return
	}

	initialValue := kv.data[args.Key]
	kv.data[args.Key] = initialValue + args.Value
	kv.requestIDs[args.ID.GetString()] = initialValue

	reply.Value = initialValue
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	kv.data = make(map[string]string)
	kv.requestIDs = make(map[string]string)

	return kv
}
