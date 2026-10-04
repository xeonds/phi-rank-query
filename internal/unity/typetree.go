package unity

import "fmt"

// TypeTreeNode mirrors UnityPy's TypeTreeNode for the blob format (version >= 12).
// Only fields needed to deserialize objects and to resolve names are kept.
type TypeTreeNode struct {
	Level       uint8
	Type        string
	Name        string
	ByteSize    int32
	Version     int16
	TypeFlags   uint8
	Index       int32
	MetaFlag    int32
	RefTypeHash uint64
	Children    []*TypeTreeNode
}

// NodeSpec is a single node in a flat (level-ordered) type tree definition.
type NodeSpec struct {
	Type     string `json:"m_Type"`
	Name     string `json:"m_Name"`
	MetaFlag int32  `json:"m_MetaFlag"`
	Level    int    `json:"m_Level"`
	ByteSize int32  `json:"m_ByteSize"`
	Version  int16  `json:"m_Version"`
}

// BuildTypeTree reconstructs a tree from a flat level-ordered node list
// (e.g. UnityPy's typetree.json).
func BuildTypeTree(specs []NodeSpec) *TypeTreeNode {
	nodes := make([]*TypeTreeNode, len(specs))
	for i, s := range specs {
		nodes[i] = &TypeTreeNode{
			Level:    uint8(s.Level),
			Type:     s.Type,
			Name:     s.Name,
			MetaFlag: s.MetaFlag,
			ByteSize: s.ByteSize,
			Version:  s.Version,
		}
	}
	if len(nodes) == 0 {
		return nil
	}
	root := nodes[0]
	stack := []*TypeTreeNode{root}
	for _, n := range nodes[1:] {
		for len(stack) > 1 && stack[len(stack)-1].Level >= n.Level {
			stack = stack[:len(stack)-1]
		}
		p := stack[len(stack)-1]
		p.Children = append(p.Children, n)
		stack = append(stack, n)
	}
	return root
}

// Walk visits the tree depth-first, calling f with the current level.
func (n *TypeTreeNode) Walk(level int, f func(*TypeTreeNode, int)) {
	f(n, level)
	for _, c := range n.Children {
		c.Walk(level+1, f)
	}
}

// ParseBlob reads a type tree in the "mhtt" blob format used by Unity 2019.4+/2020+.
func ParseBlob(r *Reader, version int) *TypeTreeNode {
	if version >= 23 {
		if string(r.Read(4)) != "mhtt" {
			panic("unity: invalid type tree blob magic")
		}
		formatVersion := r.I32()
		if int(formatVersion) != version {
			panic("unity: inconsistent type tree format version")
		}
	}
	nodeCount := int(r.I32())
	stringBufferSize := int(r.I32())

	structSize := 24
	if version >= 19 {
		structSize += 8
	}
	structData := r.Read(structSize * nodeCount)
	sb := NewReader(r.Read(stringBufferSize))

	readString := func(value uint32) string {
		if value&0x80000000 == 0 {
			sb.Seek(int(value))
			return sb.CString()
		}
		id := value & 0x7FFFFFFF
		if s, ok := commonStrings[id]; ok {
			return s
		}
		return fmt.Sprintf("#%d", id)
	}

	nodes := make([]*TypeTreeNode, 0, nodeCount)
	dr := NewReader(structData)
	for i := 0; i < nodeCount; i++ {
		n := &TypeTreeNode{}
		n.Version = dr.I16()
		n.Level = dr.U8()
		n.TypeFlags = dr.U8()
		typeOff := dr.U32()
		nameOff := dr.U32()
		n.ByteSize = dr.I32()
		n.Index = dr.I32()
		n.MetaFlag = dr.I32()
		if version >= 19 {
			n.RefTypeHash = dr.U64()
		}
		n.Type = readString(typeOff)
		n.Name = readString(nameOff)
		nodes = append(nodes, n)
	}

	// Rebuild the tree from the flat level-ordered list.
	root := nodes[0]
	stack := []*TypeTreeNode{root}
	for _, n := range nodes[1:] {
		for len(stack) > 1 && stack[len(stack)-1].Level >= n.Level {
			stack = stack[:len(stack)-1]
		}
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, n)
		stack = append(stack, n)
	}
	return root
}
