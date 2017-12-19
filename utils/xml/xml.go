package xml

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

type XmlAttr struct {
	Key string
	Val string
}

type XmlNode struct {
	Key      string
	Content  string
	Attrs    []*XmlAttr //该节点的属性
	Next     []*XmlNode
	Parent   *XmlNode //父节点
	CloseTag bool     // default false
}

func XmlNodeNew() *XmlNode {
	return new(XmlNode)
}

// 初始化XML
func XmlInit(data []byte) *XmlNode {
	return nil
}

func (this *XmlNode) Parse(data []byte) {
	if data == nil || len(data) == 0 {
		panic("data is nil")
	}
	decoder := xml.NewDecoder(bytes.NewBuffer(data))
	addNode(this, decoder)
}

func addNode(current *XmlNode, decoder *xml.Decoder) {
	token, err := decoder.Token()
	if err != nil {
		if err == io.EOF {
			return
		}
		panic(err)
	}
	switch element := token.(type) {
	case xml.StartElement:
		// 标签开始
		name := strings.ToLower(element.Name.Local)
		if name == "xml" {
			addNode(current, decoder)
		} else {
			txn := XmlNodeNew()
			txn.Key = name
			txn.Parent = current
			current.Next = append(current.Next, txn)
			addNode(txn, decoder)
		}

	case xml.EndElement:
		name := strings.ToLower(element.Name.Local)
		if name == "xml" {
			addNode(current, decoder)
		} else {
			if current.Key != name {
				panic("not found end node!")
			}
			addNode(current.Parent, decoder)
		}

	case xml.CharData:
		current.Content = string([]byte(element))
		addNode(current, decoder)
	case xml.Comment:
		addNode(current, decoder)
	case xml.Directive:
		addNode(current, decoder)
	case xml.ProcInst:
		addNode(current, decoder)
	default:
		addNode(current, decoder)
	}
}

// 添加单个节点到XML上
func (this *XmlNode)AddNode(key, value string) {
	node := new(XmlNode)
	node.Key = key
	node.Content = value
	node.CloseTag = true
	node.Parent = this
	this.Next = append(this.Next, node)
}

// 添加一颗节点树到当前节点上
func (this *XmlNode)AddNodes(node *XmlNode) {
	this.Next = append(this.Next, node)
}

// 根据Key从下一级节点中找到匹配的节点对象
func (this *XmlNode) Gets(key string) []*XmlNode {
	var node []*XmlNode
	for _, v := range this.Next {
		if v.Key == key {
			node = append(node, v)
		}
	}
	return node
}

// 根据Key从下一级节点中找到匹配的单个节点对象
func (this *XmlNode) Get(key string) *XmlNode {
	var node *XmlNode
	for _, v := range this.Next {
		if v.Key == key {
			node = v
			break
		}
	}
	return node
}

// 根据正则表达式方式从下一级节点中找到匹配的节点对象
func (this *XmlNode) GetPattern(pattern string) []*XmlNode {
	return nil
}

// 根据多级路径按层级查找匹配的节点对象
func (this *XmlNode) GetManyPathes(path string) []*XmlNode {
	// person/age

	return nil
}

// 把树形结构转为JSON字符串
func (this *XmlNode) ToJson() string {
	if this == nil {
		return ""
	}

	return getNodes(this.Next)
}

// flag 表示是否是列表
func getNode(current *XmlNode, flag bool) string {
	var content string
	if current.Key != "" {
		if !flag {
			content += "\"" + current.Key + "\":"
		}

		if current.Next != nil {
			content += getNodes(current.Next)
		}else{
			content += "\"" + current.Content +"\""
		}
	}
	return content
}

type classNode struct {
	count int
	node []*XmlNode
}

func getNodes(current []*XmlNode) string {
	if current == nil {
		return ""
	}
	var content string

	sameflag := make(map[string]*classNode)
	for _, v := range current {
		if sameflag[v.Key] == nil {
			cn := new(classNode)
			cn.node = append(cn.node, v)
			cn.count++
			sameflag[v.Key] = cn
		}else{
			cn := sameflag[v.Key]
			cn.count++
			cn.node = append(cn.node, v)
		}
	}
	content += "{"
	for k, v := range sameflag {
		if v.count > 1 {
			// 含有超过一个相同标签
			content += "\"" + k +"\":["
			for _, v1 := range v.node {
				content += getNode(v1, true) + ","
			}
			content = strings.TrimSuffix(content,",")
			content += "],"
		}else if v.count == 1{
			// 相同标签只有一个
			content += getNode(v.node[0], false) + ","
		}else{
			for _, v := range current {
				content += getNode(v, false) + ","
			}
		}
	}
	content = strings.TrimSuffix(content,",")
	content += "}"

	return  content
}