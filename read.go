// Package nbt enables robust reading and writing of Minecraft named binary tags (NBT) files.
package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
)

// ReadTag reads the next tags worth of bytes on the buffer, undertakes basic structure checks,
func ReadTag(buffer io.Reader, order binary.ByteOrder) (t tag, err error) {
	t.id, err = readTagID(buffer, order)
	if err != nil {
		return tag{}, fmt.Errorf("Unable to read tag: %w", err)
	}

	// tagEnd is used to mark the end of compound tags. This tag does not have a name, so it is only ever a single byte
	// 0. It may also be the type of empty List tags.
	if t.id == tagEnd {
		return t, nil
	}

	t.name, err = readTagName(buffer, order)
	if err != nil {
		return tag{}, fmt.Errorf("Unable to read tag: %w", err)
	}

	t.payload, err = readTagPayload(buffer, order, t.id)
	if err != nil {
		return tag{}, fmt.Errorf("Unable to read tag: %w", err)
	}

	return t, nil
}

// readTagID is intended to read the ID of a tag. The ID is the first byte in a tag. The tag ID is also known as the tag
// type. In this implementation, tag ID refers to the uint8 number (0 -> 12), and tag Type refers to the type name
// associated with that ID (ID 0 == type tagEnd, ID 12 == type tagLongArray).
func readTagID(buffer io.Reader, order binary.ByteOrder) (id uint8, err error) {
	err = binary.Read(buffer, order, &id)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tag ID: %w", err)
	}

	if id > tagLongArray {
		return 0, fmt.Errorf("ID %v not between 0 (tagEnd) and 12 (tagLongArray)", id)
	}

	return id, nil
}

// readTagName is intended to read the name of a tag. The second and third byte of a tag is an unsigned integer length
// of the tag name. The following 'length' amount of bytes is the name as a string in UTF-8 format. TagEnd is an
// exception, as it never has a name, therefore is only one byte. That is, tagEnd does not have a second and third byte
// for name length nor a series of bytes for the name.
func readTagName(buffer io.Reader, order binary.ByteOrder) (name string, err error) {
	var length int16
	err = binary.Read(buffer, order, &length)
	if err != nil {
		return "", fmt.Errorf("Unable to read tag name length for: %w", err)
	}

	nameBytes := make([]byte, length)
	err = binary.Read(buffer, order, nameBytes)
	if err != nil {
		return "", fmt.Errorf("Unable to read tag name: %w", err)
	}

	name = string(nameBytes)

	if !utf8.ValidString(name) {
		return "", fmt.Errorf("Unable to read tag name: \"%v\" contains non UTF-8 charters", name)
	}

	return name, nil
}

// readTagPayload is intended to read the variable number of subsequent bytes after the tag ID and tag Name. The number
// of bytes in the payload is dependant on the type of tag. A tagEnd does not have a payload, so expect an error if a
// tagEnd is passed as the ID.
func readTagPayload(buffer io.Reader, order binary.ByteOrder, tagID uint8) (payload any, err error) {
	switch tagID {
	case tagEnd:
		err = fmt.Errorf("Not expecting to read a tagEnd in the payload")
	case tagByte:
		payload, err = readTagBytePayload(buffer, order)
	case tagShort:
		payload, err = readTagShortPayload(buffer, order)
	case tagInt:
		payload, err = readTagIntPayload(buffer, order)
	case tagLong:
		payload, err = readTagLongPayload(buffer, order)
	case tagFloat:
		payload, err = readTagFloatPayload(buffer, order)
	case tagDouble:
		payload, err = readTagDoublePayload(buffer, order)
	case tagByteArray:
		payload, err = readTagByteArrayPayload(buffer, order)
	case tagString:
		payload, err = readTagStringPayload(buffer, order)
	case tagList:
		payload, err = readTagListPayload(buffer, order)
	case tagCompound:
		payload, err = readTagCompoundPayload(buffer, order)
	case tagIntArray:
		payload, err = readTagIntArrayPayload(buffer, order)
	case tagLongArray:
		payload, err = readTagLongArrayPayload(buffer, order)
	default:
		err = fmt.Errorf("tag ID %v not between 0 (tagEnd) and 12 (tagLongArray)", tagID)
	}
	return payload, err
}

// readTagBytePayload reads a tag payload defined as: "1 byte / 8 bits, signed. A signed integral type. Sometimes used
// for booleans." While the definition says signed, it is just a byte, use it as you will.
func readTagBytePayload(buffer io.Reader, order binary.ByteOrder) (payload byte, err error) {
	err = binary.Read(buffer, order, &payload)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tagByte payload: %w", err)
	}

	return payload, nil
}

// readTagShortPayload reads a tag payload defined as: "2 bytes / 16 bits, signed. A signed integral type."
func readTagShortPayload(buffer io.Reader, order binary.ByteOrder) (payload int16, err error) {
	err = binary.Read(buffer, order, &payload)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tagShort payload: %w", err)
	}

	return payload, nil
}

// readTagIntPayload reads a tag payload defined as: "4 bytes / 32 bits, signed. A signed integral type."
func readTagIntPayload(buffer io.Reader, order binary.ByteOrder) (payload int32, err error) {
	err = binary.Read(buffer, order, &payload)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tagInt payload: %w", err)
	}

	return payload, nil
}

// readTagLongPayload reads a tag payload defined as: "8 bytes / 64 bits, signed. A signed integral type."
func readTagLongPayload(buffer io.Reader, order binary.ByteOrder) (payload int64, err error) {
	err = binary.Read(buffer, order, &payload)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tagLong payload: %w", err)
	}

	return payload, nil
}

// readTagFloatPayload reads a tag payload defined as: "4 bytes / 32 bits, signed, IEEE 754-2008, binary32. A signed
// floating point type."
func readTagFloatPayload(buffer io.Reader, order binary.ByteOrder) (payload float32, err error) {
	err = binary.Read(buffer, order, &payload)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tagFloat payload: %w", err)
	}

	return payload, nil
}

// readTagDoublePayload reads a tag payload defined as: "8 bytes / 64 bits, signed, IEEE 754-2008, binary64. A signed
// floating point type."
func readTagDoublePayload(buffer io.Reader, order binary.ByteOrder) (payload float64, err error) {
	err = binary.Read(buffer, order, &payload)
	if err != nil {
		return 0, fmt.Errorf("Unable to read tagDouble payload: %w", err)
	}

	return payload, nil
}

// readTagByteArrayPayload reads a tag payload defined as: "A signed integer (4 bytes) size, then the bytes comprising
// an array of length size. An array of bytes." While the definition says the size is signed, that makes no sense,
// going to keep with the definition to maintain compatibility, but throw an error on negative size.
func readTagByteArrayPayload(buffer io.Reader, order binary.ByteOrder) (payload []byte, err error) {
	var size int32
	err = binary.Read(buffer, order, &size)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tagByteArray payload size: %w", err)
	}

	if size < 0 {
		return nil, fmt.Errorf("Unable to read tagByteArray payload size: size %v is negative", size)
	}

	for i := 0; i < int(size); i++ {
		var p byte
		err = binary.Read(buffer, order, &p)
		if err != nil {
			return nil, fmt.Errorf("Unable to read tagByteArray payload element %v: %w", i, err)
		}
		payload = append(payload, p)
	}

	return payload, nil
}

// readTagStringPayload reads a tag payload defined as: "An unsigned short (2 bytes) payload length, then a UTF-8 string
// resembled by length bytes. A UTF-8 string. It has a size, rather than being null terminated."
func readTagStringPayload(buffer io.Reader, order binary.ByteOrder) (payload string, err error) {
	var length uint16
	err = binary.Read(buffer, order, &length)
	if err != nil {
		return "", fmt.Errorf("Unable to read tagString payload length: %w", err)
	}

	stringPayloadBytes := make([]byte, length)
	err = binary.Read(buffer, order, stringPayloadBytes)
	if err != nil {
		return "", fmt.Errorf("Unable to read tagString payload: %w", err)
	}
	payload = string(stringPayloadBytes)

	if !utf8.ValidString(payload) {
		return "", fmt.Errorf("Unable to read tagString payload: \"%v\" contains non UTF-8 charters", payload)
	}

	return payload, nil
}

// readTagListPayload reads a tag payload defined as: "A byte denoting the tag type of the list's contents, followed by
// the list's length as a signed integer (4 bytes), then length number of payloads that correspond to the given tag
// type. A list of tag payloads, without tag types or names, apart from the one before the length." While the definition
// says the size is signed, that makes no sense, keeping with the definition in case people use negative size values to
// indicate zero length or other novel meanings.
func readTagListPayload(buffer io.Reader, order binary.ByteOrder) (payload []any, err error) {
	var tagID uint8
	err = binary.Read(buffer, order, &tagID)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tagList type: %w", err)
	}

	var length int32
	err = binary.Read(buffer, order, &length)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tagList length: %w", err)
	}

	for i := 0; i < int(length); i++ {
		p, err := readTagPayload(buffer, order, tagID)
		if err != nil {
			return nil, fmt.Errorf("Unable to read tagList payload element %v: %w", i, err)
		}
		payload = append(payload, p)
	}

	return payload, nil
}

// readTagCompoundPayload reads a tag payload defined as: "Fully formed tags, followed by a tagEnd. A list of fully
// formed tags, including their IDs, names, and payloads. No two tags may have the same name." The payload for a
// compound is an array of pointers to child tags.
func readTagCompoundPayload(buffer io.Reader, order binary.ByteOrder) (payload []*tag, err error) {
	for i := 0; ; i++ {
		t, err := ReadTag(buffer, order)
		if err != nil {
			return nil, fmt.Errorf("Unable to read tagCompound payload element %v: %w", i, err)
		}

		if t.id == tagEnd {
			break
		}
		payload = append(payload, &t)
	}
	return payload, nil
}

// readTagIntArrayPayload reads a tag payload defined as: "A signed integer size, then size number of tagInt's payloads.
// An array of tagInt's payloads." While the definition says the size is signed, that makes no sense, keeping with the
// definition in case people use negative size values to indicate zero length or other novel meanings.
func readTagIntArrayPayload(buffer io.Reader, order binary.ByteOrder) (payload []int32, err error) {
	var size int32
	err = binary.Read(buffer, order, &size)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tagIntArray payload size: %w", err)
	}

	if size < 0 {
		return nil, fmt.Errorf("Unable to read tagIntArray payload size: size %v is negative", size)
	}

	for i := 0; i < int(size); i++ {
		var p int32
		err = binary.Read(buffer, order, &p)
		if err != nil {
			return nil, fmt.Errorf("Unable to read tagIntArray payload element %v: %w", i, err)
		}
		payload = append(payload, p)
	}

	return payload, nil
}

// readTagLongArrayPayload reads a tag payload defined as: "A signed integer size, then size number of tagLong's
// payloads. An array of tagLong's payloads." While the definition says the size is signed, that makes no sense, keeping
// with the definition in case people use negative size values to indicate zero length or other novel meanings.
func readTagLongArrayPayload(buffer io.Reader, order binary.ByteOrder) (payload []int64, err error) {
	var size int32
	err = binary.Read(buffer, order, &size)
	if err != nil {
		return nil, fmt.Errorf("Unable to read tagLongArray payload size: %w", err)
	}

	if size < 0 {
		return nil, fmt.Errorf("Unable to read tagLongArray payload size: size %v is negative", size)
	}

	for i := 0; i < int(size); i++ {
		var l int64
		err = binary.Read(buffer, order, &l)
		if err != nil {
			return nil, fmt.Errorf("Unable to read tagLongArray payload element %v: %w", i, err)
		}
		payload = append(payload, l)
	}

	return payload, nil
}
