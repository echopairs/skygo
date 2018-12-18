package session

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

var _ SessionStore = (*MemoryStore)(nil)
var _ SessionDriver = (*MemoryDriver)(nil)

type MemoryDriver struct {
	sync.RWMutex
	sessions    map[string]*list.Element
	list        *list.List // for gc
	maxLifetime int64
}

type MemoryStore struct {
	sync.RWMutex
	values       map[interface{}]interface{}
	id           string
	timeAccessed time.Time
}

func (mem *MemoryStore) Set(key, value interface{}) error {
	mem.Lock()
	defer mem.Unlock()
	mem.values[key] = value
	mem.timeAccessed = time.Now()
	return nil
}

func (mem *MemoryStore) Get(key interface{}) interface{} {
	mem.Lock()
	defer mem.Unlock()
	if value, ok := mem.values[key]; ok {
		mem.timeAccessed = time.Now()
		return value
	}
	return nil
}

func (mem *MemoryStore) Delete(key interface{}) {
	mem.Lock()
	defer mem.Unlock()
	delete(mem.values, key)
}

func (mem *MemoryStore) SessionID() string {
	mem.timeAccessed = time.Now()
	return mem.id
}

func (md *MemoryDriver) SessionCreate(sid string) (SessionStore, error) {
	md.Lock()
	defer md.Unlock()
	sess := &MemoryStore{
		id:           sid,
		timeAccessed: time.Now(),
		values:       make(map[interface{}]interface{}),
	}
	element := md.list.PushFront(sess)
	md.sessions[sid] = element
	return sess, nil
}

func (md *MemoryDriver) SessionDestroy(sid string) error {
	md.Lock()
	defer md.Unlock()
	if element, ok := md.sessions[sid]; ok {
		md.list.Remove(element)
		delete(md.sessions, sid)
	}
	return nil
}

func (md *MemoryDriver) SessionRead(sid string) (SessionStore, error) {
	md.Lock()
	defer md.Unlock()
	if element, ok := md.sessions[sid]; ok {
		element.Value.(*MemoryStore).timeAccessed = time.Now()
		return element.Value.(*MemoryStore), nil
	}
	return nil, errors.New(fmt.Sprintf("sessionsid %s in not exit ", sid))
}

func (md *MemoryDriver) SessionGC(tickTime int64) {
	md.Lock()
	defer md.Unlock()
	element := md.list.Back()
	if element != nil {
		if (element.Value.(*MemoryStore).timeAccessed.Unix())+md.maxLifetime < time.Now().Unix() {
			md.list.Remove(element)
			delete(md.sessions, element.Value.(*MemoryStore).SessionID())
		}
	}
	time.AfterFunc(time.Duration(md.maxLifetime), func() {
		md.SessionGC(tickTime)
	})
}

func (md *MemoryDriver) SessionExist(sid string) bool {
	md.Lock()
	defer md.Unlock()
	if _, ok := md.sessions[sid]; ok {
		return true
	}
	return false
}

func (md *MemoryDriver) SessionAll() []string {
	md.Lock()
	defer md.Unlock()
	var list []string
	list = make([]string, 0, len(md.sessions))
	index := 0
	for key := range md.sessions {
		list[index] = key
		index++
	}
	return list
}

func (md *MemoryDriver) SessionInit(cfg *ManagerConfig) error {
	md.maxLifetime = cfg.MaxLifeTime
	md.list = list.New()
	md.sessions = make(map[string]*list.Element)
	go md.SessionGC(cfg.GcLifeTime)
	return nil
}

func init() {
	Register("memory", &MemoryDriver{})
}
