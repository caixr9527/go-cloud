package component

import "github.com/caixr9527/go-cloud/common"

type Pool map[string]any

type Singleton struct {
	cPool Pool
}

var SinglePool = &Singleton{cPool: make(Pool, 0)}

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

type component interface {
	common.Bean
	Create()
	Refresh()
	Destroy()
}

var Components = make([]component, 0)

var beans = make([]common.Bean, 0)

func RegisterComponent(c ...component) {
	Components = append(Components, c...)
}

func RegisterBean(b ...common.Bean) {
	beans = append(beans, b...)
}

type Sort []component

func (s Sort) Len() int {
	return len(s)
}

func (s Sort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sort) Less(i, j int) bool {
	return s[i].Order() < s[j].Order()
}
