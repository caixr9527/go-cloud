package common

type component interface {
	StartUp()
	Order() int
}

type ComponentSort []component

var Components = make([]component, 0)

func RegisterComponent(c ...component) {
	Components = append(Components, c...)
}

func (s ComponentSort) Len() int {
	return len(s)
}

func (s ComponentSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ComponentSort) Less(i, j int) bool {
	return s[i].Order() < s[j].Order()
}
