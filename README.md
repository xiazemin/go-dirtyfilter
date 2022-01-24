# Golang Dirty Filter

[![GoDoc](https://godoc.org/github.com/antlinker/go-dirtyfilter?status.svg)](https://godoc.org/github.com/antlinker/go-dirtyfilter)

> 基于DFA算法；
> 支持动态修改敏感词，同时支持特殊字符的筛选；
> 敏感词的存储支持内存存储及MongoDB存储。

## 获取

``` bash
$ go get -v github.com/antlinker/go-dirtyfilter
```

## 使用

``` go
package main

import (
  "fmt"

  "github.com/antlinker/go-dirtyfilter"
  "github.com/antlinker/go-dirtyfilter/store"
)

var (
  filterText = `我是需要过滤的内容，内容为：**文@@件，需要过滤。。。`
)

func main() {
  memStore, err := store.NewMemoryStore(store.MemoryConfig{
    DataSource: []string{"文件"},
  })
  if err != nil {
    panic(err)
  }
  filterManage := filter.NewDirtyManager(memStore)
  result, err := filterManage.Filter().Filter(filterText, '*', '@')
  if err != nil {
    panic(err)
  }
  fmt.Println(result)
}
```

## 输出结果

```
[文件]
```

## License

	Copyright 2016.All rights reserved.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

DAT
% git submodule add https://github.com/anknown/darts
% git submodule add https://github.com/awsong/go-darts

https://github.com/colin0000007/darts-go
https://zhuanlan.zhihu.com/p/113262718
https://github.com/heidawei/DoubleArrayTrie
https://linux.thai.net/~thep/datrie/datrie.html
https://www.co-ding.com/assets/pdf/dat.pdf


a DFA compression could be done using four linear arrays, namely default, base, next, and check. However, in a case simpler than the lexical analyzer, such as the mere trie for information retrieval, the default array could be omitted. Thus, a trie can be implemented using three arrays according to this scheme.

有限状态机可以用4个数组表示：default, base, next, and check
trie树不需要默认值

The tripple-array structure is composed of:

base. Each element in base corresponds to a node of the trie. For a trie node s, base[s] is the starting index within the next and check pool (to be explained later) for the row of the node s in the transition table.
next. This array, in coordination with check, provides a pool for the allocation of the sparse vectors for the rows in the trie transition table. The vector data, that is, the vector of transitions from every node, would be stored in this array.
check. This array works in parallel to next. It marks the owner of every cell in next. This allows the cells next to one another to be allocated to different trie nodes. That means the sparse vectors of transitions from more than one node are allowed to be overlapped.
Definition 1. For a transition from state s to t which takes character c as the input, the condition maintained in the tripple-array trie is:

  c 
s----->t

check[base[s] + c] = s
next[base[s] + c] = t

base：trie 树上的每一个节点，对于节点s，base[s]是状态转移表中next数组和check数组的开始index
next：记录状态转移表的稀疏向量，所有的状态转移向量的目标节点会被存储在这里
check：记录状态转移向量的开始节点

Walking
  t := base[s] + c;
  if check[t] = s then
      next state := next[t]
  else
      fail
  endif

Construction
Procedure Relocate(s : state; b : base_index)
{ Move base for state s to a new place beginning at b }
begin
    foreach input character c for the state s
    { i.e. foreach c such that check[base[s] + c]] = s }
    begin
        check[b + c] := s;     { mark owner }
        next[b + c] := next[base[s] + c];     { copy data }
        check[base[s] + c] := none     { free the cell }
    end;
    base[s] := b
end


Double-Array Trie
The next/check pool may be able to keep in a single array of integer couples, but the base array does not grow in parallel to the pool, and is therefore usually split.
next和check可以用pair存储，但是base不行

n the double-array structure, the base and next are merged, resulting in only two parallel arrays, namely, base and check.

check[base[s] + c] = s
base[s] + c = t

Walking
 t := base[s] + c;
  if check[t] = s then
      next state := t
  else
      fail
  endif

Construction
Procedure Relocate(s : state; b : base_index)
{ Move base for state s to a new place beginning at b }
begin
    foreach input character c for the state s
    { i.e. foreach c such that check[base[s] + c]] = s }
    begin
        check[b + c] := s;     { mark owner }
        base[b + c] := base[base[s] + c];     { copy data }
        { the node base[s] + c is to be moved to b + c;
          Hence, for any i for which check[i] = base[s] + c, update check[i] to b + c }
        foreach input character d for the node base[s] + c
        begin
            check[base[base[s] + c] + d] := b + c
        end;
        check[base[s] + c] := none     { free the cell }
    end;
    base[s] := b
end




 by splitting non-branching suffixes into single string storages, called tail,

 
https://blog.csdn.net/zzran/article/details/8462002
可以这么理解dat里面，check 存储了状态转移向量，value是箭尾（从何状态而来），index是箭尾（目标状态）
节点值和弧线上的值都存储在base数组里面，
https://www.cnblogs.com/zhangchaoyang/articles/4508266.html
   
   c 
s----->t

check[t]=s
base[s]+c=t

层次遍历来DFA的第二层

已知状态“阿”的下标是2，变量“根”、“胶”、“拉”的编号依次是4、5、6

给base[2]赋值：从小到大遍历所有的正整数，直到发现某个数正整k满足base[k+4]=base[k+5]=base[k+6]=check[k+4]=check[k+5]=check[k+6]=0。得到k=1，那么就把1赋给base[2]，同时也确定了状态“阿根”、“阿胶”、“阿拉”的下标依次是k+4、k+5、k+6，即5、6、7，而且check[5]=check[6]=check[7]=2。


最后遍历一次DFA，当某个节点已经是一个词的结尾时，按下列方法修改其base值。
if(base[i]==0)
    base[i]=-i
else
    base[i]=-base[i]

变量“阿”的编号是2，base[2]=1，变量“胶”的编号是5，base[2]+5=6，我们检查一下check[6]是否等于2。check[6]确实等于2，则继续看下一个状态转移。同时我们发现base[6]是负数，这说明“阿胶”已经是一个完整的词了。
继续看下一个状态转移，base[6]=-6，负数取其相反数，base[6]=6，变量“及”的编号是7，base[6]+7=13，我们检查一下check[13]是否等于6，发现不满足，则“阿胶及”不是一个词，甚至都是不是任意一个词的前缀。










