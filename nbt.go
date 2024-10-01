// Package nbt enables robust reading and writing of Minecraft named binary tags (NBT) files.
package nbt

import "fmt"

// tag types and IDs: source https://minecraft.fandom.com/wiki/NBT_format.
const (
	tagEnd       uint8 = 0
	tagByte      uint8 = 1
	tagShort     uint8 = 2
	tagInt       uint8 = 3
	tagLong      uint8 = 4
	tagFloat     uint8 = 5
	tagDouble    uint8 = 6
	tagByteArray uint8 = 7
	tagString    uint8 = 8
	tagList      uint8 = 9
	tagCompound  uint8 = 10
	tagIntArray  uint8 = 11
	tagLongArray uint8 = 12
)

// globalByteOrder set the endian'ness for all the binary functions to allow it to be modified easily. In minecraft Java
// edition, all numbers are encoded in big-endian. In Minecraft Bedrock edition, all numbers are encoded in
// little-endian.
// var globalByteOrder binary.ByteOrder = binary.LittleEndian

// tag is the custom type to hold common information of each tag type, with a generic payload capacity. Most tag
// payloads are the expected type.
// tagEnd: N/A, no payload
// tagByte: byte
// tagShort: int16
// tagInt: int32
// tagLong: int64
// tagFloat: float32
// tagDouble: float64
// tagByteArray: []byte
// tagString: string
// tagList: []any, assumes the type of listed tags
// tagCompound: []*tag, representing child tags an omitting the tagEnd
// tagIntArray: []int32
// tagLongArray: []int64
type tag struct {
	id      uint8
	name    string
	payload any
}

// tagType returns the name associated with the tag ID
func (t *tag) tagType() (tagType string, err error) {
	switch t.id {
	case tagEnd:
		tagType = "tagEnd"
	case tagByte:
		tagType = "tagByte"
	case tagShort:
		tagType = "tagShort"
	case tagInt:
		tagType = "tagInt"
	case tagLong:
		tagType = "tagLong"
	case tagFloat:
		tagType = "tagFloat"
	case tagDouble:
		tagType = "tagDouble"
	case tagByteArray:
		tagType = "tagByteArray"
	case tagString:
		tagType = "tagString"
	case tagList:
		tagType = "tagList"
	case tagCompound:
		tagType = "tagCompound"
	case tagIntArray:
		tagType = "tagIntArray"
	case tagLongArray:
		tagType = "tagLongArray"
	default:
		err = fmt.Errorf("tag ID %v not between 0 (tagEnd) and 12 (tagLongArray)", t.id)
	}
	return tagType, err
}
