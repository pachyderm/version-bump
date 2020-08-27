package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func lines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}

func TestEdit(t *testing.T) {
	testData := []struct {
		name        string
		input       string
		paths       []string
		replacement string
		want        string
	}{
		{
			name: "simple",
			input: lines(
				"# this is a test!",
				"apiVersion: v1",
				"kind: Foo",
			),
			paths:       []string{"kind"},
			replacement: "Bar",
			want: lines(
				"# this is a test!",
				"apiVersion: v1",
				"kind: Bar",
			),
		},
		{
			name: "missing",
			input: lines(
				"# this is a test!",
				"apiVersion: v1",
				"kind: Foo",
			),
			paths:       []string{"kindergarten"},
			replacement: "Bar",
			want: lines(
				"# this is a test!",
				"apiVersion: v1",
				"kind: Foo",
			),
		},
		{
			name: "nested",
			input: lines(
				"apiVersion: v1",
				"kind: Foo",
				"spec:",
				"  values:",
				"    a:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v1",
				"  somethingElse: true",
			),
			paths:       []string{"spec.values.a.image.tag"},
			replacement: "v2",
			want: lines(
				"apiVersion: v1",
				"kind: Foo",
				"spec:",
				"  values:",
				"    a:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v2",
				"  somethingElse: true",
			),
		},
		{
			name: "nested with lookalike sections",
			input: lines(
				"apiVersion: v1",
				"kind: Foo",
				"spec:",
				"  values:",
				"    a:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v1",
				"    b:",
				"      image:",
				"        repository: abc",
				"        # tag is not automatically updated!",
				"        tag: v1",
				"  somethingElse: true",
			),
			paths:       []string{"spec.values.a.image.tag"},
			replacement: "v2",
			want: lines(
				"apiVersion: v1",
				"kind: Foo",
				"spec:",
				"  values:",
				"    a:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v2",
				"    b:",
				"      image:",
				"        repository: abc",
				"        # tag is not automatically updated!",
				"        tag: v1",
				"  somethingElse: true",
			),
		},
		{
			name: "multiple edits",
			input: lines(
				"apiVersion: v1",
				"kind: Foo",
				"spec:",
				"  values:",
				"    a:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v1",
				"    b:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v1",
				"  somethingElse: true",
			),
			paths:       []string{"spec.values.a.image.tag", "spec.values.b.image.tag"},
			replacement: "v2",
			want: lines(
				"apiVersion: v1",
				"kind: Foo",
				"spec:",
				"  values:",
				"    a:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v2",
				"    b:",
				"      image:",
				"        repository: abc",
				"        # tag is automatically updated!",
				"        tag: v2",
				"  somethingElse: true",
			),
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			got, err := editYAML(test.input, test.paths, test.replacement)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("unexpected yaml generated:\n%s", diff)
			}
		})
	}

}
