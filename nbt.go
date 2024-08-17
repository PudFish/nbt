// Package nbt enables robust reading and writing of Minecraft named binary tags (NBT) files.
package nbt

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
