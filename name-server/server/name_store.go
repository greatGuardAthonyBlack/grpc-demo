package server

import (
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

// 数据平台输入和输出的标准化，是她想听到的。发送的时间收到的响应，平台保留，为下次处置提供基础，作为沉淀信息。
// 大量接入传感器
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
	//平台的功能，大道理。被训了，肖总可怜。这个女人看上去很务实，其实很务虚，没有把开会的内容东西说清楚。没准备会议，这些人太搞了。人家要数据需求、财政预算，人员需求。
	//这个女人其实是算钱的，平台建设需要的资源，向领导展示，并说出预算。
	//他们其实需要平台的每一个功能所消耗的资源，最好把各个模块拆碎了（模块需要的数据资源，他的数量，强度）跟女人汇报。把各个部门拉过来整合自己的需求。
	//需求、钱
	//各个部门需要监控的风险，采集这些风险需要只能加的设备，需要提供给风险预警系统的接口。风险预警平台展现了
	//数据通道打通。
	//从技术上来说，其实是采样的数据来源说明，跟业务功能的联系
	//肖总重复女人的意思。
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
