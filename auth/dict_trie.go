package auth

type dictTrie struct {
	name     string
	children []*dictTrie
	end      bool
}

type DictTrie struct {
	root        *dictTrie
	initialized bool
}

var Allow = New()

func New() *DictTrie {
	return &DictTrie{root: &dictTrie{name: "", children: make([]*dictTrie, 0)}}
}

func (d *DictTrie) Insert(word string) {

}

func (d *DictTrie) Search(word string) bool {
	node := d.root
	for _, v := range node.children {
		if (v.name == word && v.end) || v.name == "**" {
			return true
		}
		node = v
	}
	return false
}
