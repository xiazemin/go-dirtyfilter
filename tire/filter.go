package tire

import (
	"fmt"
	"sort"
)

type TireFilter interface {
	//过滤敏感词列表
	FilterResult(texts []string) [][]string
	GetRoot() *node
	Print()
}

// NewTireFilter 创建节点过滤器，实现敏感词的过滤
// 从切片中读取敏感词数据
func NewTireFilter(text [][]string, ignoreOrder bool) TireFilter {
	tf := &tireFilter{
		root:        newNode(),
		ignoreOrder: ignoreOrder,
	}
	for i, l := 0, len(text); i < l; i++ {
		tf.addDirtyWords(text[i])
	}
	return tf
}

type tireFilter struct {
	ignoreOrder bool
	root        *node
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

func (tf *tireFilter) FilterResult(texts []string) [][]string {
	if tf.ignoreOrder {
		sort.Strings(texts)
		fmt.Println(texts)
	}

	var notEmptyTexts []string
	for _, text := range texts {
		if text == "" {
			continue
		}
		notEmptyTexts = append(notEmptyTexts, text)
	}
	var result [][]string
	if len(notEmptyTexts) > 0 {
		result = tf.doFilter(notEmptyTexts)
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
			result = append(result, []string{uchars[i]})
		}
		for j := i + 1; j < ul; j++ {
			if _, ok := n.child[uchars[j]]; !ok {
				if tf.ignoreOrder {
					continue
				}
				break
			}
			n = n.child[uchars[j]]
			buf = append(buf, uchars[j])
			if n.end {
				result = append(result, buf)
			}
		}
		buf = make([]string, 0)
		n = tf.root
	}
	return result
}
