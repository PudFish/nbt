// Package nbt enables robust reading and writing of Minecraft named binary tags (NBT) files.
package nbt

import "testing"

func TestTagType(t *testing.T) {
	successCases := []struct {
		name string
		wantTagType string
		t tag
	}{
		{"tagEnd (0)", "tagEnd", tag{0, "", nil}},
		{"tagByte (1)", "tagByte", tag{1, "", nil}},
		{"tagShort (2)", "tagShort", tag{2, "", nil}},
		{"tagInt (3)", "tagInt", tag{3, "", nil}},
		{"tagLong (4)", "tagLong", tag{4, "", nil}},
		{"tagFloat (5)", "tagFloat", tag{5, "", nil}},
		{"tagDouble (6)", "tagDouble", tag{6, "", nil}},
		{"tagByteArray (7)", "tagByteArray", tag{7, "", nil}},
		{"tagString (8)", "tagString", tag{8, "", nil}},
		{"tagList (9)", "tagList", tag{9, "", nil}},
		{"tagCompound (10)", "tagCompound", tag{10, "", nil}},
		{"tagIntArray (11)", "tagIntArray", tag{11, "", nil}},
		{"tagLongArray (12)", "tagLongArray", tag{12, "", nil}},
	}
	for _, successCase := range successCases {
		t.Run("Test success case: "+successCase.name, func(t *testing.T) {
			gotTagType, gotErr := successCase.t.tagType()
			if gotTagType != successCase.wantTagType {
				t.Errorf("got %v, want %v", gotTagType, successCase.wantTagType)
			}
			if gotErr != nil {
				t.Errorf("got %v, want nil", gotErr)
			}
		})
	}

	t.Run("Test failure case: tag id out of range", func(t *testing.T) {
		failTag := tag{13, "", nil}
		_, gotErr := failTag.tagType()
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}