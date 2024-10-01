// Package nbt enables robust reading and writing of Minecraft named binary tags (NBT) files.
package nbt

import "testing"

func TestTagType(t *testing.T) {
	successCases := []struct {
		name        string
		wantTagType string
		t           Tag
	}{
		{"tagEnd (0)", "tagEnd", Tag{0, "", nil}},
		{"tagByte (1)", "tagByte", Tag{1, "", nil}},
		{"tagShort (2)", "tagShort", Tag{2, "", nil}},
		{"tagInt (3)", "tagInt", Tag{3, "", nil}},
		{"tagLong (4)", "tagLong", Tag{4, "", nil}},
		{"tagFloat (5)", "tagFloat", Tag{5, "", nil}},
		{"tagDouble (6)", "tagDouble", Tag{6, "", nil}},
		{"tagByteArray (7)", "tagByteArray", Tag{7, "", nil}},
		{"tagString (8)", "tagString", Tag{8, "", nil}},
		{"tagList (9)", "tagList", Tag{9, "", nil}},
		{"tagCompound (10)", "tagCompound", Tag{10, "", nil}},
		{"tagIntArray (11)", "tagIntArray", Tag{11, "", nil}},
		{"tagLongArray (12)", "tagLongArray", Tag{12, "", nil}},
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
		failTag := Tag{13, "", nil}
		_, gotErr := failTag.tagType()
		if gotErr == nil {
			t.Errorf("got nil, want non-nil")
		}
	})
}
