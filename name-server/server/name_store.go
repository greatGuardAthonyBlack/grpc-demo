package server

import (
	"log"
	"sync"
	"time"
)

type Addr struct {
	serviceName string
	addr        string
	expireTime  int64
}

type Store struct {
	expireMap   map[int64][]*Addr
	expireLock  sync.Mutex
	instanceMap map[string]map[string]*Addr
	dataLock    sync.Mutex
}

func (r *Store) deleteNoBlock(serviceName string, addr string) {
	if instanceMap := r.instanceMap[serviceName]; instanceMap != nil {
		delete(instanceMap, addr)
	}
}

func (r *Store) ClearExpireInstance() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			list := make([]*Addr, 0)
			cache := GetNameServerCache()
			cache.expireLock.Lock()
			for expire, items := range r.expireMap {
				if expire < time.Now().Unix()-10 {
					list = append(list, items...)
					delete(r.expireMap, expire)
				}

			}
			cache.expireLock.Unlock()

			cache.dataLock.Lock()
			for _, addr := range list {
				if addr.expireTime < time.Now().Unix()-10 {
					log.Printf("timeout eliminate service :[%s ---:%s]", addr.serviceName, addr.addr)
					r.deleteNoBlock(addr.serviceName, addr.addr)
				}
			}
			cache.dataLock.Unlock()
		}
	}
}

var nameServerStore *Store

var timeout = time.Second * 10

func init() {
	nameServerStore = &Store{
		expireMap:   make(map[int64][]*Addr),
		instanceMap: make(map[string]map[string]*Addr),
	}
	go nameServerStore.ClearExpireInstance()
}

func GetNameServerCache() *Store {
	return nameServerStore
}

func Register(serviceName string, addr string) {
	cache := GetNameServerCache()

	instance := &Addr{
		serviceName: serviceName,
		addr:        addr,
		expireTime:  time.Now().Add(timeout).Unix(),
	}
	//todo: using timing wheel manage expired service instances
	cache.expireLock.Lock()
	if candidates := cache.expireMap[instance.expireTime]; candidates == nil {
		cache.expireMap[instance.expireTime] = make([]*Addr, 0)
	}
	cache.expireMap[instance.expireTime] = append(cache.expireMap[instance.expireTime], instance)
	cache.expireLock.Unlock()

	cache.dataLock.Lock()
	mp := cache.instanceMap[serviceName]
	if mp == nil {
		cache.instanceMap[serviceName] = make(map[string]*Addr)
	}
	cache.instanceMap[serviceName][addr] = instance
	cache.dataLock.Unlock()

}

func Delete(serviceName string, addr string) {
	cache := GetNameServerCache()
	cache.dataLock.Lock()
	cache.deleteNoBlock(serviceName, addr)
	cache.dataLock.Unlock()
}

func Keepalive(serviceName string, addr string) {
	cache := GetNameServerCache()
	mp, ok := cache.instanceMap[serviceName]
	if !ok {
		return
	}

	if mp == nil {
		return
	}

	cache.dataLock.Lock()
	instance := cache.instanceMap[serviceName][addr]
	instance.expireTime = time.Now().Add(timeout).Unix()
	cache.instanceMap[serviceName][addr] = instance
	cache.dataLock.Unlock()

	cache.expireLock.Lock()

	expireAddressList := cache.expireMap[instance.expireTime]
	if expireAddressList == nil {
		cache.expireMap[instance.expireTime] = make([]*Addr, 0)
	}
	cache.expireMap[instance.expireTime] = append(cache.expireMap[instance.expireTime], instance)
	cache.expireLock.Unlock()
}

func GetService(serviceName string) []string {
	cache := GetNameServerCache()
	cache.dataLock.Lock()
	defer cache.dataLock.Unlock()

	addrMap, ok := cache.instanceMap[serviceName]
	if !ok {
		return []string{}
	}
	if addrMap == nil || len(addrMap) == 0 {
		return []string{}
	}
	candidates := make([]string, 0)
	for _, instance := range addrMap {
		if instance.expireTime < time.Now().Unix() {
			continue
		}
		candidates = append(candidates, instance.addr)
	}
	return candidates
}
