package utils

import (
	"reflect"
	"testing"
)

func TestTags(t *testing.T) {

	t.Run("merge", func(t *testing.T) {
		testCases := []struct {
			name         string
			tags1, tags2 []string
			expectedTags []string
		}{
			{
				name:         "merge two non-empty slices",
				tags1:        []string{"tag1", "tag2", "tag3"},
				tags2:        []string{"tag4", "tag5", "tag6"},
				expectedTags: []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"},
			},
			{
				name:         "merge an empty slice with a non-empty slice",
				tags1:        []string{},
				tags2:        []string{"tag1", "tag2", "tag3"},
				expectedTags: []string{"tag1", "tag2", "tag3"},
			},
			{
				name:         "merge two empty slices",
				tags1:        []string{},
				tags2:        []string{},
				expectedTags: []string{},
			},
		}

		for _, tc := range testCases {
			if !reflect.DeepEqual(Merge(tc.tags1, tc.tags2...), tc.expectedTags) {
				t.Errorf("unexpected result: got %v, want %v", Merge(tc.tags1, tc.tags2...), tc.expectedTags)
			}
		}
	})
}
