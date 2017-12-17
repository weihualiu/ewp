package xml


type XmlNode struct {
	Key string
	Content string
	Next []*xmlNode
}

func XmlNodeNew() *XmlNode() {
	return new(XmlNode)
}

// 初始化XML
func XmlInit(data []byte) *XmlNode {
	
}

// 根据Key从下一级节点中找到匹配的节点对象
func (this *XmlNode)Get(key string) []*XmlNode {
	return nil
}

// 根据正则表达式方式从下一级节点中找到匹配的节点对象
func (this *XmlNode)GetPattern(pattern string) []*XmlNode {
	return nil
}

// 根据多级路径按层级查找匹配的节点对象
func (this *XmlNode)GetManyPathes(path string) []*XmlNode {
	return nil
}

// 把树形结构转为JSON字符串
func (this *XmlNode)ToJson() string {
	return ""
}

// 添加节点到当前节点中
func (this *XmlNode)Add(node *XmlNode) error {
	return nil
}
