// Package nbt enables robust reading and writing of Minecraft named binary tags (NBT) files.
package nbt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"testing"
	"testing/iotest"
)

func TestReadTagID(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	t.Run("Test lower inbound value", func(t *testing.T) {
		var wantID byte
		buffer := bytes.NewBuffer([]byte{wantID})
		gotID, gotErr := readTagID(buffer, order)
		if gotID != wantID {
			t.Errorf("got %v, want %v", gotID, wantID)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	t.Run("Test upper inbound value", func(t *testing.T) {
		var wantID byte = 12
		buffer := bytes.NewBuffer([]byte{wantID})
		gotID, gotErr := readTagID(buffer, order)
		if gotID != wantID {
			t.Errorf("got %v, want %v", gotID, wantID)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check out of bounds tag ID", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{13})
		_, gotErr := readTagID(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagID(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of empty buffer", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{})
		_, gotErr := readTagID(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagName(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	t.Run("Test typical tag name", func(t *testing.T) {
		wantName := "BiomeOverride"
		nameLength := make([]byte, 2)
		binary.LittleEndian.PutUint16(nameLength, uint16(len(wantName)))

		buffer := bytes.NewBuffer([]byte{})
		buffer.Write(nameLength)
		buffer.Write([]byte(wantName))

		gotName, gotErr := readTagName(buffer, order)
		if gotName != wantName {
			t.Errorf("got %v, want %v", gotName, wantName)
		}
		if gotErr != nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Test empty tag name", func(t *testing.T) {
		wantName := ""
		buffer := bytes.NewBuffer([]byte(wantName))
		gotName, gotErr := readTagName(buffer, order)
		if gotName != wantName {
			t.Errorf("got %v, want empty name", gotName)
		}
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check missing tag name for non-zero length", func(t *testing.T) {
		nameLength := make([]byte, 2)
		binary.LittleEndian.PutUint16(nameLength, uint16(len("BiomeOverride")))

		buffer := bytes.NewBuffer([]byte{})
		buffer.Write(nameLength)
		// but do not write name

		_, gotErr := readTagName(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagName(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagBytePayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	wantBytes := []byte{
		0, 123, math.MaxUint8,
	}
	for _, wantByte := range wantBytes {
		t.Run("Test reading an inbound byte value", func(t *testing.T) {
			buffer := bytes.NewBuffer([]byte{wantByte})
			gotByte, gotErr := readTagBytePayload(buffer, order)
			if gotByte != wantByte {
				t.Errorf("got %v, want %v", gotByte, wantByte)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagBytePayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagShortPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	wantShorts := []int16{
		0, 12345, -6789, math.MinInt16, math.MaxInt16,
	}
	for _, wantShort := range wantShorts {
		t.Run("Test reading an inbound short value", func(t *testing.T) {
			b := make([]byte, 2)
			binary.LittleEndian.PutUint16(b, uint16(wantShort))
			buffer := bytes.NewBuffer(b)

			gotShort, gotErr := readTagShortPayload(buffer, order)
			if gotShort != wantShort {
				t.Errorf("got %v, want %v", gotShort, wantShort)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagShortPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{123})
		_, gotErr := readTagShortPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagIntPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	wantInts := []int32{
		0, 1234567, -89012345, math.MinInt32, math.MaxInt32,
	}
	for _, wantInt := range wantInts {
		t.Run("Test reading an inbound int value", func(t *testing.T) {
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, uint32(wantInt))
			buffer := bytes.NewBuffer(b)

			gotInt, gotErr := readTagIntPayload(buffer, order)
			if gotInt != wantInt {
				t.Errorf("got %v, want %v", gotInt, wantInt)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagIntPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(12345))
		buffer := bytes.NewBuffer(b)

		_, gotErr := readTagIntPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagLongPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	wantLongs := []int64{
		0, 123456789012, -345678901234, math.MinInt64, math.MaxInt64,
	}
	for _, wantLong := range wantLongs {
		t.Run("Test reading an inbound long value", func(t *testing.T) {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(wantLong))
			buffer := bytes.NewBuffer(b)

			gotLong, gotErr := readTagLongPayload(buffer, order)
			if gotLong != wantLong {
				t.Errorf("got %v, want %v", gotLong, wantLong)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagLongPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(1234567))
		buffer := bytes.NewBuffer(b)

		_, gotErr := readTagLongPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagFloatPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	wantFloats := []float32{
		0.0, 1.23, math.Pi, -4.5e+6, math.SmallestNonzeroFloat32, math.MaxFloat32,
	}
	for _, wantFloat := range wantFloats {
		t.Run("Test reading an inbound float value", func(t *testing.T) {
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, math.Float32bits(wantFloat))
			buffer := bytes.NewBuffer(b)

			gotFloat, gotErr := readTagFloatPayload(buffer, order)
			if gotFloat != wantFloat {
				t.Errorf("got %v, want %v", gotFloat, wantFloat)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	// edge case handling
	wantFloats = []float32{
		float32(math.Inf(1)), float32(math.Inf(-1)),
	}
	for _, wantFloat := range wantFloats {
		t.Run("Test reading an edge case float value", func(t *testing.T) {
			b := make([]byte, 4)
			binary.LittleEndian.PutUint32(b, math.Float32bits(wantFloat))
			buffer := bytes.NewBuffer(b)

			gotFloat, gotErr := readTagFloatPayload(buffer, order)
			if gotFloat != wantFloat {
				t.Errorf("got %v, want %v", gotFloat, wantFloat)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	t.Run("Test reading NaN float value", func(t *testing.T) {
		wantFloat := float32(math.NaN())
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, math.Float32bits(wantFloat))
		buffer := bytes.NewBuffer(b)

		gotFloat, gotErr := readTagFloatPayload(buffer, order)
		gotFloat64 := float64(gotFloat)
		if !math.IsNaN(gotFloat64) {
			t.Errorf("got %v, want %v", gotFloat, wantFloat)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagFloatPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(12345))
		buffer := bytes.NewBuffer(b)

		_, gotErr := readTagFloatPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagDoublePayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	wantDoubles := []float64{
		0.0, 1.23, math.Pi, -4.5e+6, math.SmallestNonzeroFloat64, math.MaxFloat64,
	}
	for _, wantDouble := range wantDoubles {
		t.Run("Test reading an inbound double value", func(t *testing.T) {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, math.Float64bits(wantDouble))
			buffer := bytes.NewBuffer(b)

			gotDouble, gotErr := readTagDoublePayload(buffer, order)
			if gotDouble != wantDouble {
				t.Errorf("got %v, want %v", gotDouble, wantDouble)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	// edge case handling
	wantDoubles = []float64{
		math.Inf(1), math.Inf(-1),
	}
	for _, wantDouble := range wantDoubles {
		t.Run("Test reading an edge case double value", func(t *testing.T) {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, math.Float64bits(wantDouble))
			buffer := bytes.NewBuffer(b)

			gotDouble, gotErr := readTagDoublePayload(buffer, order)
			if gotDouble != wantDouble {
				t.Errorf("got %v, want %v", gotDouble, wantDouble)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	t.Run("Test reading NaN double value", func(t *testing.T) {
		wantDouble := math.NaN()
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(wantDouble))
		buffer := bytes.NewBuffer(b)

		gotDouble, gotErr := readTagDoublePayload(buffer, order)
		if !math.IsNaN(gotDouble) {
			t.Errorf("got %v, want %v", gotDouble, wantDouble)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagDoublePayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(12345678))
		buffer := bytes.NewBuffer(b)

		_, gotErr := readTagDoublePayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagByteArrayPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	t.Run("Test typical byte array", func(t *testing.T) {
		wantByteArray := []byte{0, 255, 1, 50, 48, 0, 0, 74}
		size := int32(len(wantByteArray))

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(size))
		buffer := bytes.NewBuffer(b)
		buffer.Write(wantByteArray)

		gotByteArray, gotErr := readTagByteArrayPayload(buffer, order)
		for i, gotByte := range gotByteArray {
			if gotByte != wantByteArray[i] {
				t.Errorf("got %v, want %v", gotByte, wantByteArray[i])
			}
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})
	t.Run("Test empty byte array", func(t *testing.T) {
		wantByteArray := []byte{}
		size := int32(len(wantByteArray))

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(size))
		buffer := bytes.NewBuffer(b)
		buffer.Write(wantByteArray)

		gotByteArray, gotErr := readTagByteArrayPayload(buffer, order)
		for i, gotByte := range gotByteArray {
			if gotByte != wantByteArray[i] {
				t.Errorf("got %v, want %v", gotByte, wantByteArray[i])
			}
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagByteArrayPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		wantByteArray := []byte{1, 2, 3}
		size := int32(len(wantByteArray))

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(size))
		buffer := bytes.NewBuffer(b)
		// do not write the byte array

		_, gotErr := readTagByteArrayPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Test negative size array", func(t *testing.T) {
		wantByteArray := []byte{1, 2, 3}
		var size int32 = -3

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(size))
		buffer := bytes.NewBuffer(b)
		buffer.Write(wantByteArray)

		_, gotErr := readTagByteArrayPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagStringPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	t.Run("Test typical string", func(t *testing.T) {
		wantString := "Dummy string used for testing the TagString payload read, 你好世界"
		size := uint16(len(wantString))

		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, size)
		buffer := bytes.NewBuffer(b)
		buffer.WriteString(wantString)

		gotString, gotErr := readTagStringPayload(buffer, order)
		if gotString != wantString {
			t.Errorf("got %v, want %v", gotString, wantString)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	t.Run("Test empty string", func(t *testing.T) {
		wantString := ""
		size := uint16(len(wantString))

		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, size)
		buffer := bytes.NewBuffer(b)
		buffer.WriteString(wantString)

		gotString, gotErr := readTagStringPayload(buffer, order)
		if gotString != wantString {
			t.Errorf("got %v, want %v", gotString, wantString)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check handling of non-UTF-8 characters", func(t *testing.T) {
		// 0xc0 and 0xff are invalid
		var nonUTF8String = []byte{0x41, 0xc0, 0xff, 0x61}

		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(len(nonUTF8String)))
		buffer := bytes.NewBuffer(b)
		buffer.Write(nonUTF8String)

		_, gotErr := readTagStringPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagStringPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		wantString := "A short string"
		size := uint16(len(wantString))

		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, size)
		buffer := bytes.NewBuffer(b)
		// Do not write the string

		_, gotErr := readTagStringPayload(buffer, order)

		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagListPayload(t *testing.T) {
	var order binary.ByteOrder = binary.LittleEndian

	// success cases
	t.Run("Test typical tag list", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{tagInt})

		wantList := []int32{0, 1, 128, -4}

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(len(wantList)))
		buffer.Write(b)

		for _, wantInt := range wantList {
			wi := make([]byte, 4)
			binary.LittleEndian.PutUint32(wi, uint32(wantInt))
			buffer.Write(wi)
		}

		gotList, gotErr := readTagListPayload(buffer, order)
		for i, gotInt := range gotList {
			if gotInt != wantList[i] {
				t.Errorf("got %v, want %v", gotInt, wantList[i])
			}
		}

		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	t.Run("Test empty tag list", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{tagString})

		wantList := []string{}

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(len(wantList)))
		buffer.Write(b)

		gotList, gotErr := readTagListPayload(buffer, order)
		if len(gotList) > 0 {
			t.Errorf("got length %v list, want nil length", len(gotList))
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	// failure cases
	t.Run("Check handling of empty buffer", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf(""))
		_, gotErr := readTagListPayload(errBuffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of missing length", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{tagInt})

		_, gotErr := readTagListPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})

	t.Run("Check handling of partial buffer", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{tagInt})

		wantList := []int32{0, 1, 128, -4}

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(len(wantList)))
		buffer.Write(b)
		// Do not write the list

		_, gotErr := readTagListPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}
