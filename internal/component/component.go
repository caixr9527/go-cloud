package component

import "github.com/caixr9527/go-cloud/common"

type pool map[string]any

type Singleton struct {
	cPool pool
}

var SinglePool = &Singleton{cPool: make(pool, 0)}

func (s *Singleton) Register(key string, obj any) {
	if key != "" && obj != nil {
		s.cPool[key] = obj
	}
}

func (s *Singleton) Del(key string) {
	if key != "" {
		delete(s.cPool, key)
	}
}

func (s *Singleton) Get(key string) (any, bool) {
	value, ok := s.cPool[key]
	return value, ok
}

type Component interface {
	common.Bean
	Create() bool
	Refresh() bool
	Destroy()
}

var Beans = make([]common.Bean, 0)

type Sort []common.Bean

func (s Sort) Len() int {
	return len(s)
}

func (s Sort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sort) Less(i, j int) bool {
	return s[i].Order() < s[j].Order()
}
