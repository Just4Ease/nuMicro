package Idempotence

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/Just4Ease/nuMicro/store"
)

type Request struct {
	checksum   string
	store      string
	action     string
	cacheStore *store.Store
}

func New(d []byte, store, action string, cacheStore store.Store) *Request {
	hash := md5.Sum(d)
	hashChannel := make(chan []byte, 1)
	hashChannel <- hash[:]
	ch := <-hashChannel
	checksum := hex.EncodeToString(ch)
	//_ = cacheStore.Write(checksum, "", 0)
	return &Request{
		checksum: checksum,
		store:    "idempotent::" + store,
		action:   action,
		//cacheStore: cacheStore,
	}
}

func (r *Request) IsOngoing() bool {
	//checksum := Cache.Get(r.store)
	return false
}

func (r *Request) Cleanup() {

}
