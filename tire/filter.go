package tire

import (
	"sort"
	"strings"
	"sync"
)

type TireFilter interface {
	//过滤敏感词列表
	FilterResult(texts []string) [][]string
	FilterResultCount(texts []string) map[string]int
	GetRoot() *node
	Print()
}

// NewTireFilter 创建节点过滤器，实现敏感词的过滤
// 从切片中读取敏感词数据
func NewTireFilter(text [][]string, ignoreOrder bool) TireFilter {
	tf := &tireFilter{
		root:        newNode(),
		ignoreOrder: ignoreOrder,
		lock:        &sync.RWMutex{},
	}
	tf.lock.Lock()
	defer tf.lock.Unlock()
	for i, l := 0, len(text); i < l; i++ {
		tf.addDirtyWords(text[i])
	}
	return tf
}

type tireFilter struct {
	ignoreOrder bool
	root        *node
	lock        *sync.RWMutex
}

func (tf *tireFilter) GetRoot() *node {
	return tf.root
}
func (tf *tireFilter) Print() {
	tf.root.print(0)
}

func (tf *tireFilter) addDirtyWords(text []string) {
	if tf.ignoreOrder {
		sort.Strings(text)
	}
	n := tf.root
	for i, l := 0, len(text); i < l; i++ {
		if text[i] == "" {
			continue
		}

		if _, ok := n.child[text[i]]; !ok {
			n.child[text[i]] = newNode()
		}
		n = n.child[text[i]]
	}
	n.end = true
}

func (tf *tireFilter) FilterResultCount(texts []string) map[string]int {
	tf.lock.RLock()
	defer tf.lock.RUnlock()
	result := tf.FilterResult(texts)
	re := make(map[string]int)
	frequency := make(map[string]int)
	if tf.ignoreOrder {
		for _, text := range texts {
			frequency[text]++
		}
	}

	for _, li := range result {
		if tf.ignoreOrder {
			if _, ok := re[strings.Join(li, ",")]; ok {
				continue
			}
			min := frequency[li[0]]
			for _, v := range li {
				if frequency[v] < min {
					min = frequency[v]
					//不关心顺序的话，词组的频率是几个词中出现次数最少的词的频率
				}
			}
			re[strings.Join(li, ",")] = min
		} else {
			re[strings.Join(li, ",")]++
		}
	}
	return re
}

func (tf *tireFilter) FilterResult(texts []string) [][]string {
	if tf.ignoreOrder {
		sort.Strings(texts)
	}
	tf.lock.RLock()
	defer tf.lock.RUnlock()
	var notEmptyTexts []string
	for _, text := range texts {
		if text == "" {
			continue
		}
		notEmptyTexts = append(notEmptyTexts, text)
	}
	var result [][]string
	if len(notEmptyTexts) > 0 {
		if tf.ignoreOrder {
			result = tf.doFilterIgnoreOrder(notEmptyTexts)
		} else {
			result = tf.doFilter(notEmptyTexts)
		}
	}
	r1 := make(map[string]int)
	var res [][]string
	for _, r := range result {
		if _, ok := r1[strings.Join(r, ",")]; !ok {
			r1[strings.Join(r, ",")]++
			res = append(res, r)
		}
	}
	return res
}

type nodeWithIndex struct {
	n   *node
	i   int
	buf []string
}

func (tf *tireFilter) doFilterIgnoreOrder(uchars []string) [][]string {
	var result [][]string
	ul := len(uchars)
	n := tf.root
	var buf []string

	var q []*nodeWithIndex

	for i := 0; i < ul; i++ {
		if _, ok := n.child[uchars[i]]; !ok {
			continue
		}
		n1 := n.child[uchars[i]]
		buf = append(buf, uchars[i])
		q = append(q, &nodeWithIndex{
			n:   n1,
			i:   i,
			buf: buf[:len(buf):len(buf)],
		})
		if n1.end {
			result = append(result, buf[:len(buf):len(buf)])
		}
		buf = make([]string, 0)
	}

	for len(q) > 0 {
		no := q[0]
		for j := no.i + 1; j < ul; j++ {
			n = no.n
			if _, ok := n.child[uchars[j]]; !ok {
				continue
			}
			buf := no.buf
			n1 := n.child[uchars[j]]
			buf = append(buf, uchars[j])
			q = append(q, &nodeWithIndex{
				n:   n1,
				i:   j,
				buf: buf[:len(buf):len(buf)],
			})
			if n1.end {
				result = append(result, buf[:len(buf):len(buf)])
			}
		}
		q = q[1:]
	}
	return result
}

func (tf *tireFilter) doFilter(uchars []string) [][]string {
	var result [][]string
	ul := len(uchars)
	n := tf.root
	var buf []string
	for i := 0; i < ul; i++ {
		if _, ok := n.child[uchars[i]]; !ok {
			continue
		}
		n = n.child[uchars[i]]
		buf = append(buf, uchars[i])
		if n.end {
			result = append(result, buf[:len(buf):len(buf)])
		}
		for j := i + 1; j < ul; j++ {
			if _, ok := n.child[uchars[j]]; !ok {
				break
			}
			n = n.child[uchars[j]]
			buf = append(buf, uchars[j])
			if n.end {
				result = append(result, buf[:len(buf):len(buf)])
			}
		}
		buf = make([]string, 0)
		n = tf.root
	}
	return result
}
