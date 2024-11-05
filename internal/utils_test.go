package internal

import (
	"errors"
	"reflect"
	"testing"
)

func TestOnError(t *testing.T) {
	t.Run("Callback_On_Error", func(t *testing.T) {
		// Should run the callback on error
		err := errors.New("test")
		didRun := false
		OnError(err, func(err error) {
			didRun = true
		})
		if didRun == false {
			t.Errorf("OnError() unexpectedly skipped callback function on populated error.")
		}
	})

	t.Run("No_Callback_On_Empty", func(t *testing.T) {
		// Should not run the callback if the error is nil
		OnError(nil, func(err error) {
			t.Errorf("OnError() unexpectedly invoked callback function on empty error.")
		})
	})
}

func TestSubtractSlice(t *testing.T) {
	tests := []struct {
		name     string
		left     []string
		right    []string
		expected []string
	}{
		{
			name:     "All Present in Right",
			left:     []string{"one", "two", "three"},
			right:    []string{"one", "two", "three"},
			expected: []string{},
		},
		{
			name:     "None Present in Right",
			left:     []string{"one", "two", "three"},
			right:    []string{},
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "Some Present in Right",
			left:     []string{"one", "two", "three"},
			right:    []string{"one", "two"},
			expected: []string{"three"},
		},
		{
			name:     "Empty Left",
			left:     []string{},
			right:    []string{"one", "two"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SubtractSlice(tt.left, tt.right)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SubtractSlice() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
