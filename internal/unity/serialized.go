package unity

import "fmt"

// Common Unity class IDs.
const (
	ClassTextAsset     = 49
	ClassAudioClip     = 83
	ClassMonoBehaviour = 114
	ClassMonoScript    = 115
	ClassTexture2D     = 28
	ClassSprite        = 213
	ClassAssetBundle   = 142
)

type SerializedType struct {
	ClassID         int32
	IsStrippedType  bool
	ScriptTypeIndex int16
	ScriptID        []byte
	OldTypeHash     []byte
	Node            *TypeTreeNode

	// ref types
	ClassName    string
	NameSpace    string
	AssemblyName string
	TypeDeps     []int32
}

type ObjectInfo struct {
	PathID    int64
	ByteStart int64
	ByteSize  int64
	TypeID    int32
	ClassID   int32
}

type SerializedFile struct {
	MetadataSize int64
	FileSize     int64
	Version      int32
	DataOffset   int64
	Endian       string

	UnityVersion   string
	TargetPlatform int32
	EnableTypeTree bool

	Types     []*SerializedType
	Objects   []*ObjectInfo
	BigID     int32
	Externals []string
}

// ParseSerializedFile parses a SerializedFile (the "CAB-..." payload of a bundle).
func ParseSerializedFile(r *Reader) (*SerializedFile, error) {
	sf := &SerializedFile{Endian: "<"}
	r.SetBE(true)
	sf.MetadataSize = int64(r.U32())
	sf.FileSize = int64(r.U32())
	sf.Version = r.I32()
	sf.DataOffset = int64(r.U32())

	if sf.Version >= 9 {
		if r.U8() != 0 {
			sf.Endian = ">"
		}
		r.Read(3) // reserved
		if sf.Version >= 22 {
			sf.MetadataSize = int64(r.U32())
			sf.FileSize = int64(r.U64())
			sf.DataOffset = int64(r.U64())
			r.U64() // unknown
		}
	} else {
		r.Seek(int(sf.FileSize - sf.MetadataSize))
		if r.U8() != 0 {
			sf.Endian = ">"
		}
	}
	r.SetBE(sf.Endian == ">")

	if sf.Version >= 7 {
		sf.UnityVersion = r.CString()
	}
	if sf.Version >= 8 {
		sf.TargetPlatform = r.I32()
	}
	if sf.Version >= 13 {
		sf.EnableTypeTree = r.Bool()
	}

	typeCount := int(r.I32())
	for i := 0; i < typeCount; i++ {
		st, err := parseSerializedType(r, sf, false)
		if err != nil {
			return nil, err
		}
		sf.Types = append(sf.Types, st)
	}

	if sf.Version >= 7 && sf.Version < 14 {
		sf.BigID = r.I32()
	}

	objectCount := int(r.I32())
	for i := 0; i < objectCount; i++ {
		obj, err := parseObject(r, sf)
		if err != nil {
			return nil, err
		}
		sf.Objects = append(sf.Objects, obj)
	}

	if sf.Version >= 11 {
		scriptCount := int(r.I32())
		for i := 0; i < scriptCount; i++ {
			r.I32()
			if sf.Version < 14 {
				r.I32()
			} else {
				r.Align(4)
				r.I64()
			}
		}
	}

	externalCount := int(r.I32())
	for i := 0; i < externalCount; i++ {
		if sf.Version >= 6 {
			r.CString()
		}
		if sf.Version >= 5 {
			r.Read(16)
			r.I32()
		}
		sf.Externals = append(sf.Externals, r.CString())
	}

	if sf.Version >= 20 {
		refTypeCount := int(r.I32())
		for i := 0; i < refTypeCount; i++ {
			if _, err := parseSerializedType(r, sf, true); err != nil {
				return nil, err
			}
		}
	}

	if sf.Version >= 5 {
		r.CString() // userInformation
	}
	return sf, nil
}

func parseSerializedType(r *Reader, sf *SerializedFile, isRefType bool) (*SerializedType, error) {
	st := &SerializedType{ScriptTypeIndex: -1}
	st.ClassID = r.I32()
	if sf.Version >= 16 {
		st.IsStrippedType = r.Bool()
	}
	if sf.Version >= 17 {
		st.ScriptTypeIndex = r.I16()
	}
	if sf.Version >= 13 {
		if (isRefType && st.ScriptTypeIndex >= 0) ||
			(sf.Version < 16 && st.ClassID < 0) ||
			(sf.Version >= 16 && st.ClassID == ClassMonoBehaviour) {
			st.ScriptID = append([]byte(nil), r.Read(16)...)
		}
		st.OldTypeHash = append([]byte(nil), r.Read(16)...)
	}

	if sf.EnableTypeTree {
		if sf.Version >= 23 {
			r.Read(16) // type tree content hash
			ttSize := r.I32()
			if ttSize != 0 {
				st.Node = ParseBlob(r, int(sf.Version))
			}
		} else if sf.Version >= 12 || sf.Version == 10 {
			st.Node = ParseBlob(r, int(sf.Version))
		} else {
			return nil, fmt.Errorf("unity: inline type tree (version %d) not supported", sf.Version)
		}
		if sf.Version >= 21 {
			if isRefType {
				st.ClassName = r.CString()
				st.NameSpace = r.CString()
				st.AssemblyName = r.CString()
			} else {
				n := int(r.I32())
				for i := 0; i < n; i++ {
					st.TypeDeps = append(st.TypeDeps, r.I32())
				}
			}
		}
	}
	return st, nil
}

func parseObject(r *Reader, sf *SerializedFile) (*ObjectInfo, error) {
	o := &ObjectInfo{}
	if sf.BigID != 0 {
		o.PathID = r.I64()
	} else if sf.Version < 14 {
		o.PathID = int64(r.I32())
	} else {
		r.Align(4)
		o.PathID = r.I64()
	}
	if sf.Version >= 22 {
		o.ByteStart = r.I64()
	} else {
		o.ByteStart = int64(r.U32())
	}
	o.ByteStart += sf.DataOffset
	o.ByteSize = int64(r.U32())
	o.TypeID = r.I32()
	if sf.Version < 16 {
		o.ClassID = int32(r.U16())
	} else {
		if int(o.TypeID) >= len(sf.Types) {
			return nil, fmt.Errorf("unity: object type id %d out of range", o.TypeID)
		}
		o.ClassID = sf.Types[o.TypeID].ClassID
	}
	if sf.Version < 11 {
		r.U16() // isDestroyed
	}
	return o, nil
}

// ObjectBytes returns the raw serialized bytes for an object.
func (sf *SerializedFile) ObjectBytes(data []byte, o *ObjectInfo) []byte {
	return data[o.ByteStart : o.ByteStart+o.ByteSize]
}
