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
	requestIDs map[int64]string
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if value, exists := kv.requestIDs[args.ID]; exists {
		reply.Value = value
		//delete(kv.requestIDs, args.ID)
		return
	}

	value := kv.data[args.Key]

	kv.requestIDs[args.ID] = value
	reply.Value = value
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if value, exists := kv.requestIDs[args.ID]; exists {
		reply.Value = value
		//delete(kv.requestIDs, args.ID)
		return
	}

	kv.data[args.Key] = args.Value
	kv.requestIDs[args.ID] = args.Value

	reply.Value = args.Value
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if value, exists := kv.requestIDs[args.ID]; exists {
		reply.Value = value
		//delete(kv.requestIDs, args.ID)
		return
	}

	initialValue := kv.data[args.Key]
	kv.data[args.Key] = initialValue + args.Value
	kv.requestIDs[args.ID] = initialValue

	reply.Value = initialValue
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	kv.data = make(map[string]string)
	kv.requestIDs = make(map[int64]string)

	return kv
}
