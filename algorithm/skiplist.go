package algorithm

import (
	"fmt"
	"math/rand"
)

// Node 跳表节点
type Node struct {
	pre    *Node
	next   *Node
	down   *Node
	isHead bool
	key    Key
	value  interface{}
}

func (n *Node) GetValue() (value interface{}) {
	if n == nil {
		return
	}
	value = n.value
	return
}

func (n *Node) Next() (next *Node, exist bool) {
	if n == nil || n.next == nil {
		return
	}
	next, exist = n.pre, true
	return
}

func (n *Node) Prev() (pre *Node, exist bool) {
	if n == nil || n.pre == nil || n.pre.isHead {
		return
	}
	pre, exist = n.pre, true
	return
}

func (n *Node) Down() (down *Node, exist bool) {
	if n == nil || n.down == nil {
		return
	}
	down, exist = n.down, true
	return
}

// Key 索引
type Key interface {
	Equal(then Key) bool
	Less(then Key) bool
	Check(then Key) bool
}

// SkipList 跳表
type SkipList struct {
	top   *Node
	layer int32
	len   int32
}

// Print 打印整个跳表
func (l *SkipList) Print() {
	fmt.Println("layer: ", l.layer)
	fmt.Println("len: ", l.len)
	layerHead := l.top
	for layerHead != nil {
		curNode := layerHead.next
		for curNode != nil {
			str := fmt.Sprintf("[%#v:%#v]", curNode.key, curNode.value)
			if curNode.next != nil {
				str += " -->"
			}
			fmt.Printf(str)
			curNode = curNode.next
		}
		fmt.Println("---layer:---")
		layerHead = layerHead.down
	}
}

// Len 只输出底层的结点个数
func (l *SkipList) Len() int32 {
	return l.len
}

// Sort 将最底层的链表转化为切片输出
func (l *SkipList) Sort() []interface{} {
	var ret []interface{}
	head := l.top
	if head == nil {
		return nil
	}

	for head.down != nil {
		head = head.down
	}
	node := head.next
	for node != nil {
		ret = append(ret, node.value)
		node = node.next
	}
	return ret
}

// GetMin 获得最小值
func (l *SkipList) GetMin() (*Node, bool) {
	if l.top == nil {
		return nil, false
	}
	head := l.top
	for head.down != nil {
		head = head.down
	}

	if head == nil {
		return nil, false
	}
	node := head.next
	if node == nil {
		return nil, false
	}
	return node, true
}

// GetMax 获得最大值
func (l *SkipList) GetMax() (*Node, bool) {
	node := l.top
	if node == nil {
		return nil, false
	}
	for node.down != nil || node.next != nil {
		for node.next != nil {
			node = node.next
		}
		if node.down != nil {
			node = node.down
		}
	}
	if node.isHead {
		return nil, false
	}
	return node, true
}

// Insert 插入节点
func (l *SkipList) Insert(key Key, value interface{}) {
	preNodes, node, ok := l.find(key, 0, false)
	if ok {
		foundNode := node
		for foundNode != nil {
			foundNode.value = value
			foundNode = foundNode.down
		}
	} else {
		isBreak := false
		var downNode *Node
		for i := len(preNodes) - 1; i >= 0; i-- {
			preNode := preNodes[i]
			nextNode := preNode.next
			aNode := &Node{key: key, value: value}

			preNode.next = aNode
			aNode.pre = preNode
			aNode.next = nextNode
			aNode.down = downNode
			if nextNode != nil {
				nextNode.pre = aNode
			}

			downNode = aNode
			if !l.needCreatUpNode() {
				isBreak = true
				break
			}
		}

		//建立顶层节点
		if !isBreak && l.needCreatUpNode() {
			aNode := &Node{key: key, value: value}
			head := &Node{down: l.top, isHead: true}
			aNode.down = downNode
			aNode.pre = head
			head.next = aNode
			l.top = head
			l.layer++
		}

		// 更新结点数
		l.len++
	}
}

// Delete 删除结点
func (l *SkipList) Delete(tarKey Key) {
	_, node, ok := l.find(tarKey, 0, false)
	if !ok {
		return
	}
	foundNode := node
	for foundNode != nil {
		pre := foundNode.pre
		next := foundNode.next
		pre.next = next
		if next != nil {
			next.pre = pre
		}

		foundNode = foundNode.down
	}

	l.len--
	// 删除头结点
	head := l.top
	for head != nil && head.next == nil {
		head = head.down
		l.layer--
	}
	l.top = head
}

// Find 查找结点
// 只用于搜索匹配玩家
func (l *SkipList) Find(key Key, searchLimit int) (*Node, bool) {
	_, node, ok := l.find(key, searchLimit, true)
	if ok {
		return node, true
	}

	return nil, false
}

//FindLeft 返回左侧结点
func (l *SkipList) FindLeft(key Key) (*Node, bool) {
	_, node, ok := l.find(key, 0, false)
	if ok {
		node = node.pre
	}

	if node != nil && !node.isHead {
		return node, true
	}

	return nil, false
}

//FindRight 返回右侧侧结点
func (l *SkipList) FindRight(key Key) (*Node, bool) {
	_, node, _ := l.find(key, 0, false)
	if node != nil && node.next != nil {
		return node.next, true
	}

	return nil, false
}

// ClearAll 清空节点
func (l *SkipList) ClearAll() {
	l.top = nil
	l.len = 0
	l.layer = 0
}

func (l *SkipList) find(tarKey Key, searchLimit int, isSearch bool) ([]*Node, *Node, bool) {
	var cmpCnt int
	var preNodeList []*Node
	preNode, curNode := l.top, l.top
	for curNode != nil {
		if curNode.isHead {
			if curNode.next != nil {
				preNode = curNode
				curNode = curNode.next
			} else {
				preNode = curNode
				curNode = curNode.down
				preNodeList = append(preNodeList, preNode)
			}
			continue
		}

		cmpCnt++
		if searchLimit > 0 && cmpCnt > searchLimit {
			return preNodeList, preNode, false
		}

		if tarKey.Equal(curNode.key) {
			return nil, curNode, true
		}

		if isSearch && tarKey.Check(curNode.key) {
			return nil, curNode, true
		}

		if tarKey.Less(curNode.key) {
			curNode = preNode.down
			preNodeList = append(preNodeList, preNode)
		} else if curNode.next != nil {
			preNode = curNode
			curNode = curNode.next
		} else {
			preNode = curNode
			curNode = curNode.down
			preNodeList = append(preNodeList, preNode)
		}
	}

	return preNodeList, preNode, false
}

// 是否需要创建上层结点, 随机决定，概率可以设置
func (l *SkipList) needCreatUpNode() bool {
	if l.top == nil {
		return true
	}

	base := 3
	if rand.Intn(2*base) < base {
		return true
	}

	return false
}
