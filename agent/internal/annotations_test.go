package internal

import (
	"github.com/open-policy-agent/opa/ast"
	"reflect"
	"testing"
)

func TestExtractAnnotations(t *testing.T) {
	tests := []struct {
		name     string
		comments []*ast.Comment
		expected map[string]interface{}
	}{
		{
			name:     "No Comments",
			comments: []*ast.Comment{},
			expected: map[string]interface{}{},
		},
		{
			name: "Comments but No Metadata",
			comments: []*ast.Comment{
				{Text: []byte("Just a comment"), Location: &ast.Location{Row: 1}},
				{Text: []byte("Another comment"), Location: &ast.Location{Row: 2}},
			},
			expected: map[string]interface{}{},
		},
		{
			name: "Only Metadata",
			comments: []*ast.Comment{
				{Text: []byte("METADATA"), Location: &ast.Location{Row: 1}},
				{Text: []byte("name: test-policy"), Location: &ast.Location{Row: 2}},
				{Text: []byte("description: A test policy"), Location: &ast.Location{Row: 3}},
			},
			expected: map[string]interface{}{
				"name":        "test-policy",
				"description": "A test policy",
			},
		},
		{
			name: "Comments Before Metadata",
			comments: []*ast.Comment{
				{Text: []byte("This is a general comment"), Location: &ast.Location{Row: 1}},
				{Text: []byte("METADATA"), Location: &ast.Location{Row: 2}},
				{Text: []byte("name: test-policy"), Location: &ast.Location{Row: 3}},
				{Text: []byte("description: A test policy with preceding comments"), Location: &ast.Location{Row: 4}},
			},
			expected: map[string]interface{}{
				"name":        "test-policy",
				"description": "A test policy with preceding comments",
			},
		},
		{
			name: "Metadata with Additional Text After",
			comments: []*ast.Comment{
				{Text: []byte("METADATA"), Location: &ast.Location{Row: 1}},
				{Text: []byte("name: test-policy"), Location: &ast.Location{Row: 2}},
				{Text: []byte("description: A test policy with additional comments after metadata"), Location: &ast.Location{Row: 3}},
				{Text: []byte(""), Location: &ast.Location{Row: 4}}, // Simulate an empty line
				{Text: []byte("Additional comment after metadata"), Location: &ast.Location{Row: 5}},
			},
			expected: map[string]interface{}{
				"name":        "test-policy",
				"description": "A test policy with additional comments after metadata",
			},
		},
		{
			name: "Metadata with other comments in file",
			comments: []*ast.Comment{
				{Text: []byte("METADATA"), Location: &ast.Location{Row: 1}},
				{Text: []byte("name: test-policy"), Location: &ast.Location{Row: 2}},
				{Text: []byte("description: A test policy with additional comments after metadata"), Location: &ast.Location{Row: 3}},
				{Text: []byte("Additional comment after metadata"), Location: &ast.Location{Row: 5}},
			},
			expected: map[string]interface{}{
				"name":        "test-policy",
				"description": "A test policy with additional comments after metadata",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractAnnotations(tt.comments)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractAnnotations() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
