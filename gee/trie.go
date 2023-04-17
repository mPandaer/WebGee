package gee

import (
	"strings"
)

//利用Trie前缀树实现动态路由。 /:name/name
// 不支持 /:name

// 树节点
type node struct {
	pattern  string  //主要用于叶子节点，存储完整的URL路径
	part     string  //将域名切割后每个节点部分 比如 /name/abc 那么part = name 或者 abc
	children []*node //子节点
	isWild   bool    //该节点是否是模糊匹配
}

// 实现一个辅助方法，辅助插入节点
func (node *node) matchChild(part string) *node {
	for _, child := range node.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 实现另外一个辅助方法，方便查询节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 实现插入节点序列
func (n *node) insert(pattern string, parts []string, height int) {
	//base case 插入节点如果插入的深度已经和节点序列的长度一样的就应该停止递归。因为不需要在插入了 因为深度是 从零开始
	if height >= len(parts) {
		n.pattern = pattern
		return
	}

	//如果没有满足停止插入的条件。
	part := parts[height]       //取出要匹配的部分
	child := n.matchChild(part) //查看当前节点的子孩子是否有一个满足匹配条件的。

	//如果没有就创建一个字节点。
	if child == nil {
		child = &node{
			part: part, isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child) //并加入到当前节点的子集合中
	}
	//递归插入下一个part部分。
	child.insert(pattern, parts, height+1)
}

// 实现搜索指定节点序列 返回叶子节点。
func (n *node) search(pattern string, parts []string, height int) *node {
	// base case 搜索的深度 >= 字节序列的长度就应该停止搜索，或者匹配到了* 因为规定 /*/c 是无效URL 因为 *匹配一切。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	//同理，如果不满足停止条件 取出要匹配的部分
	part := parts[height]
	//查看他的孩子们是有匹配的。 为什么会返回一个切片呢？或者说是一个集合，因为存在模糊匹配。
	children := n.matchChildren(part)

	//遍历匹配的集合，看下一个需要匹配的部分是否匹配成功
	for _, child := range children {
		res := child.search(pattern, parts, height+1)
		if res != nil {
			return res
		}
	}

	return nil

}
