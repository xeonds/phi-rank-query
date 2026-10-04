package unity

// Generic type-tree deserializer, mirroring UnityPy's TypeTreeHelper.read_value
// with as_dict semantics. Values are Go maps/slices/primitives.

const metaAlign = int32(0x4000)

func isPrimitive(t string) bool {
	switch t {
	case "SInt8", "UInt8", "char",
		"short", "SInt16", "unsigned short", "UInt16",
		"int", "SInt32", "unsigned int", "UInt32", "Type*",
		"long long", "SInt64", "unsigned long long", "UInt64", "FileSize",
		"float", "double", "bool":
		return true
	}
	return false
}

func readPrimitive(node *TypeTreeNode, r *Reader) any {
	switch node.Type {
	case "SInt8":
		return int64(r.I8())
	case "UInt8", "char":
		return int64(r.U8())
	case "short", "SInt16":
		return int64(r.I16())
	case "unsigned short", "UInt16":
		return int64(r.U16())
	case "int", "SInt32":
		return int64(r.I32())
	case "unsigned int", "UInt32", "Type*":
		return int64(r.U32())
	case "long long", "SInt64":
		return r.I64()
	case "unsigned long long", "UInt64", "FileSize":
		return r.U64()
	case "float":
		return float64(r.F32())
	case "double":
		return r.F64()
	case "bool":
		return r.Bool()
	}
	panic("unity: not a primitive: " + node.Type)
}

func readAlignedString(r *Reader) string {
	length := int(r.I32())
	if length > 0 && length <= r.Remaining() {
		s := string(r.Read(length))
		r.Align(4)
		return s
	}
	return ""
}

// ReadTypeTree deserializes one object using its type tree root node.
func ReadTypeTree(root *TypeTreeNode, r *Reader) any {
	return readValue(root, r)
}

func readValue(node *TypeTreeNode, r *Reader) any {
	align := node.MetaFlag&metaAlign != 0
	var v any

	switch {
	case isPrimitive(node.Type):
		v = readPrimitive(node, r)
	case node.Type == "string":
		v = readAlignedString(r)
	case node.Type == "TypelessData":
		n := int(r.I32())
		v = append([]byte(nil), r.Read(n)...)
	case node.Type == "pair":
		v = []any{readValue(node.Children[0], r), readValue(node.Children[1], r)}
	case len(node.Children) > 0 && node.Children[0].Type == "Array":
		if node.Children[0].MetaFlag&metaAlign != 0 {
			align = true
		}
		size := int(r.I32())
		subtype := node.Children[0].Children[1]
		if isPrimitive(subtype.Type) || subtype.Type == "string" || subtype.Type == "TypelessData" {
			v = readValueArray(subtype, r, size)
		} else {
			arr := make([]any, size)
			for i := range arr {
				arr[i] = readValue(subtype, r)
			}
			v = arr
		}
	default:
		m := make(map[string]any, len(node.Children))
		for _, c := range node.Children {
			m[c.Name] = readValue(c, r)
		}
		v = m
	}

	if align {
		r.Align(4)
	}
	return v
}

func readValueArray(node *TypeTreeNode, r *Reader, size int) any {
	switch node.Type {
	case "SInt8", "UInt8", "char":
		return append([]byte(nil), r.Read(size)...)
	case "string":
		out := make([]string, size)
		for i := range out {
			out[i] = readAlignedString(r)
		}
		return out
	case "TypelessData":
		out := make([][]byte, size)
		for i := range out {
			n := int(r.I32())
			out[i] = append([]byte(nil), r.Read(n)...)
		}
		return out
	}
	out := make([]any, size)
	for i := range out {
		out[i] = readPrimitive(node, r)
	}
	return out
}

// asInt coerces a deserialized numeric value to int64.
func asInt(v any) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case uint64:
		return int64(x)
	case float64:
		return int64(x)
	}
	return 0
}
