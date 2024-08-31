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
	successCases := []struct {
		name   string
		wantID byte
		order  binary.ByteOrder
		input  []byte
	}{
		{"lower bound value", 0, binary.LittleEndian, []byte{0x00}},
		{"in bound value", 8, binary.LittleEndian, []byte{0x08}},
		{"upper bound value", 12, binary.LittleEndian, []byte{0x0C}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotID, gotErr := readTagID(buffer, successCase.order)
			if gotID != successCase.wantID {
				t.Errorf("got %v, want %v", gotID, successCase.wantID)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"out of bound value", binary.LittleEndian, []byte{0x0D}},
		{"empty buffer", binary.LittleEndian, []byte{}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagID(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagID(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagName(t *testing.T) {
	successCases := []struct {
		name     string
		wantName string
		order    binary.ByteOrder
		input    []byte
	}{
		{"typical tag name", "BiomeOverride", binary.LittleEndian, []byte{0x0D, 0x00, 0x42, 0x69, 0x6F, 0x6D, 0x65,
			0x4F, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65}},
		{"another typical tag name", "saved_with_toggled_experiments", binary.LittleEndian, []byte{0x1E, 0x00, 0x73,
			0x61, 0x76, 0x65, 0x64, 0x5F, 0x77, 0x69, 0x74, 0x68, 0x5F, 0x74, 0x6F, 0x67, 0x67, 0x6C, 0x65, 0x64, 0x5F,
			0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6D, 0x65, 0x6E, 0x74, 0x73}},
		{"empty tag name", "", binary.LittleEndian, []byte{0x00, 0x00}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotName, gotErr := readTagName(buffer, successCase.order)
			if gotName != successCase.wantName {
				t.Errorf("got %v, want %v", gotName, successCase.wantName)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0x01}},
		{"empty tag name with non-zero length", binary.LittleEndian, []byte{0x0D, 0x00}},
		{"typical tag name with incorrect longer length", binary.LittleEndian, []byte{0xA2, 0x00, 0x47, 0x61, 0x6D,
			0x65, 0x54, 0x79, 0x70, 0x65}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagName(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagName(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagBytePayload(t *testing.T) {
	successCases := []struct {
		name     string
		wantByte byte
		order    binary.ByteOrder
		input    []byte
	}{
		{"zero", 0, binary.LittleEndian, []byte{0x00}},
		{"in bound value", 123, binary.LittleEndian, []byte{0x7B}},
		{"max uint8", 255, binary.LittleEndian, []byte{0xFF}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotByte, gotErr := readTagBytePayload(buffer, successCase.order)
			if gotByte != successCase.wantByte {
				t.Errorf("got %v, want %v", gotByte, successCase.wantByte)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagBytePayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagBytePayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagShortPayload(t *testing.T) {
	successCases := []struct {
		name      string
		wantShort int16
		order     binary.ByteOrder
		input     []byte
	}{
		{"zero", 0, binary.LittleEndian, []byte{0x00, 0x00}},
		{"positive", 12345, binary.LittleEndian, []byte{0x39, 0x30}},
		{"negative", -6789, binary.LittleEndian, []byte{0x7B, 0xE5}},
		{"max int16", 32767, binary.LittleEndian, []byte{0xFF, 0x7F}},
		{"min int16", -32768, binary.LittleEndian, []byte{0x00, 0x80}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotShort, gotErr := readTagShortPayload(buffer, successCase.order)
			if gotShort != successCase.wantShort {
				t.Errorf("got %v, want %v", gotShort, successCase.wantShort)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0xD4}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagShortPayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagShortPayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagIntPayload(t *testing.T) {
	successCases := []struct {
		name    string
		wantInt int32
		order   binary.ByteOrder
		input   []byte
	}{
		{"zero", 0, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00}},
		{"positive", 1234567, binary.LittleEndian, []byte{0x87, 0xD6, 0x12, 0x00}},
		{"negative", -89012345, binary.LittleEndian, []byte{0x87, 0xC7, 0xB1, 0xFA}},
		{"max int32", 2147483647, binary.LittleEndian, []byte{0xFF, 0xFF, 0xFF, 0x7F}},
		{"min int32", -2147483648, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x80}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotInt, gotErr := readTagIntPayload(buffer, successCase.order)
			if gotInt != successCase.wantInt {
				t.Errorf("got %v, want %v", gotInt, successCase.wantInt)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0xD4, 0xF2}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagIntPayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagIntPayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagLongPayload(t *testing.T) {
	successCases := []struct {
		name     string
		wantLong int64
		order    binary.ByteOrder
		input    []byte
	}{
		{"zero", 0, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"positive", 123456789012, binary.LittleEndian, []byte{0x14, 0x1A, 0x99, 0xBE, 0x1C, 0x00, 0x00, 0x00}},
		{"negative", -345678901234, binary.LittleEndian, []byte{0x0E, 0x90, 0xEE, 0x83, 0xAF, 0xFF, 0xFF, 0xFF}},
		{"max int16", 9223372036854775807, binary.LittleEndian, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}},
		{"min int16", -9223372036854775808, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x80}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotLong, gotErr := readTagLongPayload(buffer, successCase.order)
			if gotLong != successCase.wantLong {
				t.Errorf("got %v, want %v", gotLong, successCase.wantLong)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0xD4, 0xF2, 0x12, 0xBB}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagLongPayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagLongPayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagFloatPayload(t *testing.T) {
	successCases := []struct {
		name      string
		wantFloat float32
		order     binary.ByteOrder
		input     []byte
	}{
		{"positive zero", 0.0, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00}},
		{"negative zero", 0.0, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x80}},
		{"positive", 1.23, binary.LittleEndian, []byte{0xA4, 0x70, 0x9D, 0x3F}},
		{"negative", -4.56e+6, binary.LittleEndian, []byte{0x00, 0x29, 0x8B, 0xCA}},
		{"pi", 3.1415927, binary.LittleEndian, []byte{0xDB, 0x0F, 0x49, 0x40}},
		{"smallest non zero float32", 1e-45, binary.LittleEndian, []byte{0x01, 0x00, 0x00, 0x00}},
		{"max float32", 3.4028235e+38, binary.LittleEndian, []byte{0xFF, 0xFF, 0x7F, 0x7F}},
		{"positive infinity", float32(math.Inf(1)), binary.LittleEndian, []byte{0x00, 0x00, 0x80, 0x7F}},
		{"negative infinity", float32(math.Inf(-1)), binary.LittleEndian, []byte{0x00, 0x00, 0x80, 0xFF}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotFloat, gotErr := readTagFloatPayload(buffer, successCase.order)
			if gotFloat != successCase.wantFloat {
				t.Errorf("got %v, want %v", gotFloat, successCase.wantFloat)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	t.Run("Test success case: NaN", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{0x00, 0x00, 0xC0, 0x7F})
		gotFloat, gotErr := readTagFloatPayload(buffer, binary.LittleEndian)
		if !math.IsNaN(float64(gotFloat)) {
			t.Errorf("got %v, want NaN", gotFloat)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0xD4, 0xC3, 0x13}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagFloatPayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagFloatPayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagDoublePayload(t *testing.T) {
	successCases := []struct {
		name       string
		wantDouble float64
		order      binary.ByteOrder
		input      []byte
	}{
		{"positive zero", 0.0, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"negative zero", 0.0, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}},
		{"positive", 1.23456, binary.LittleEndian, []byte{0x38, 0x32, 0x8F, 0xFC, 0xC1, 0xC0, 0xF3, 0x3F}},
		{"negative", -7.89012e3, binary.LittleEndian, []byte{0x85, 0xEB, 0x51, 0xB8, 0x1E, 0xD2, 0xBE, 0xC0}},
		{"pi", 3.141592653589793, binary.LittleEndian, []byte{0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40}},
		{"smallest non zero float64", 5e-324, binary.LittleEndian, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00}},
		{"max float64", 1.7976931348623157e+308, binary.LittleEndian, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xEF,
			0x7F}},
		{"positive infinity", math.Inf(1), binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0, 0x7F}},
		{"negative infinity", math.Inf(-1), binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF0,
			0xFF}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotDouble, gotErr := readTagDoublePayload(buffer, successCase.order)
			if gotDouble != successCase.wantDouble {
				t.Errorf("got %v, want %v", gotDouble, successCase.wantDouble)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	t.Run("Test success case: NaN", func(t *testing.T) {
		buffer := bytes.NewBuffer([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF8, 0xFF})
		gotDouble, gotErr := readTagDoublePayload(buffer, binary.LittleEndian)
		if !math.IsNaN(gotDouble) {
			t.Errorf("got %v, want NaN", gotDouble)
		}
		if gotErr != nil {
			t.Errorf("got %v, want nil", gotErr)
		}
	})

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0xD4, 0xC3, 0x13, 0xAA, 0x43}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagDoublePayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagDoublePayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagByteArrayPayload(t *testing.T) {
	successCases := []struct {
		name          string
		wantByteArray []byte
		order         binary.ByteOrder
		input         []byte
	}{
		{"empty byte array", []byte{}, binary.LittleEndian, []byte{0x00, 0x00, 0x00, 0x00}},
		{"single byte array", []byte{45}, binary.LittleEndian, []byte{0x01, 0x00, 0x00, 0x00, 0x2D}},
		{"typical byte array", []byte{0, 255, 1, 50, 48, 0, 0, 74}, binary.LittleEndian, []byte{0x08, 0x00, 0x00, 0x00,
			0x00, 0xFF, 0x01, 0x32, 0x30, 0x00, 0x00, 0x4A}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotByteArray, gotErr := readTagByteArrayPayload(buffer, successCase.order)

			gotLength := len(gotByteArray)
			wantLength := len(successCase.input) - 4
			if gotLength != wantLength {
				t.Errorf("got length=%v, want length=%v", gotLength, wantLength)
			}

			for i, gotByte := range gotByteArray {
				if gotByte != successCase.wantByteArray[i] {
					t.Errorf("got %v, want %v, i=%v", gotByte, successCase.wantByteArray[i], i)
				}
			}

			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0x01, 0x00}},
		{"empty array with non-zero size", binary.LittleEndian, []byte{0x08, 0x00, 0x00, 0x00}},
		{"typical array with incorrect larger size", binary.LittleEndian, []byte{0x04, 0x00, 0x00, 0x00, 0x12, 0x34}},
		{"negative size array", binary.LittleEndian, []byte{0xFD, 0xFF, 0xFF, 0xFF, 0x12, 0x34, 0x56}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagByteArrayPayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagByteArrayPayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}

func TestReadTagStringPayload(t *testing.T) {
	successCases := []struct {
		name       string
		wantString string
		order      binary.ByteOrder
		input      []byte
	}{
		{"empty string", "", binary.LittleEndian, []byte{0x00, 0x00}},
		{"string with single byte UTF-8 characters", "Dummy string used for testing the TagString payload read",
			binary.LittleEndian, []byte{0x38, 0x00, 0x44, 0x75, 0x6D, 0x6D, 0x79, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6E,
				0x67, 0x20, 0x75, 0x73, 0x65, 0x64, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6E,
				0x67, 0x20, 0x74, 0x68, 0x65, 0x20, 0x54, 0x61, 0x67, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x20, 0x70,
				0x61, 0x79, 0x6C, 0x6F, 0x61, 0x64, 0x20, 0x72, 0x65, 0x61, 0x64, 0x2C}},
		{"string with multi-byte UTF-8 characters", "你好世界", binary.LittleEndian, []byte{0x0C, 0x00, 0xE4, 0xBD, 0xA0,
			0xE5, 0xA5, 0xBD, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C}},
		{"string with single and multi-byte UTF-8 characters", "你好 hello 世界 world", binary.LittleEndian, []byte{0x19,
			0x00, 0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD, 0x20, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0xE4, 0xB8, 0x96, 0xE7,
			0x95, 0x8C, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(successCase.input)
			gotString, gotErr := readTagStringPayload(buffer, successCase.order)
			if gotString != successCase.wantString {
				t.Errorf("got %v, want %v", gotString, successCase.wantString)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	failureCases := []struct {
		name  string
		order binary.ByteOrder
		input []byte
	}{
		{"empty buffer", binary.LittleEndian, []byte{}},
		{"partial buffer", binary.LittleEndian, []byte{0x34}},
		{"empty string with non-zero length", binary.LittleEndian, []byte{0x0D, 0x00}},
		{"typical string with incorrect longer length", binary.LittleEndian, []byte{0xA2, 0x00, 0x47, 0x61, 0x6D, 0x65,
			0x54, 0x79, 0x70, 0x65}},
		{"string with invalid UTF-8 characters", binary.LittleEndian, []byte{0x04, 0x00, 0x41, 0xc0, 0xff, 0x61}},
	}
	for _, failureCase := range failureCases {
		t.Run("Test failure case: "+failureCase.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(failureCase.input)
			_, gotErr := readTagStringPayload(buffer, failureCase.order)
			if gotErr == nil {
				t.Errorf("got nil, want non-nil")
			}
		})
	}

	t.Run("Test failure case: broken io.Reader", func(t *testing.T) {
		errBuffer := iotest.ErrReader(fmt.Errorf("mock broken io.reader"))
		_, gotErr := readTagStringPayload(errBuffer, binary.LittleEndian)
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
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
		binary.LittleEndian.PutUint32(b, 4)
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

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, 4)
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

		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, 4)
		buffer.Write(b)
		// Do not write the list

		_, gotErr := readTagListPayload(buffer, order)
		if gotErr == nil {
			t.Errorf("got %v, want non-nil", gotErr)
		}
	})
}

func TestReadTagCompoundPayload(t *testing.T) {
	_ = t
}

func TestReadTagIntArrayPayload(t *testing.T) {
	_ = t
}

func TestReadTagLongArrayPayload(t *testing.T) {
	_ = t
}
