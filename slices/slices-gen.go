// (see slices.go)

// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const toolName string = "slices-gen"

type record struct {
	ToolName      string
	PrimitiveType string
	SliceType     string
}

var codeGenerator = template.Must(template.New("Code").Parse(
	`// Code generated by {{ .ToolName }} -- DO NOT EDIT

package slices

// {{ .SliceType }} is a slice of {{ .PrimitiveType }}'s.
type {{ .SliceType }} []{{ .PrimitiveType }}

// Contains determines if the provided value is present in the slice.
// Runtime: O(n) -- appropriate for small-n slices.
func (s {{ .SliceType }}) Contains(value {{ .PrimitiveType }}) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}
`))

var testGenerator = template.Must(template.New("Test").Parse(
	`// Code generated by {{ .ToolName }} -- DO NOT EDIT

package slices

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{ .SliceType }}", func() {
	// More of a static assertion, but eh:
	It("is a slice of {{ .PrimitiveType }}'s", func() {
		var s {{ .SliceType }} = []{{ .PrimitiveType }}{1, 2, 3}
		s[0] = 0 // prevent unused error
	})

	_ = Describe(".Contains", func() {
		It("returns true if the provided value is in the slice", func() {
			s := {{ .SliceType }}{1, 2, 3, 4, 5}
			Expect(s.Contains(4)).To(BeTrue())
		})

		It("returns false if the provided value is not in the slice", func() {
			s := {{ .SliceType }}{1, 2, 3, 4, 5}
			Expect(s.Contains(42)).To(BeFalse())
		})
	})
})
`))

// generateFile uses the provided template and record to emit a file.
// 'generator' is the Template to use.
// 'rec' is the record to provide to the Template.
// 'filename' is the name of the file to generate.
func generateFile(generator *template.Template, rec *record, filename string) error {

	// here we're using the best practice of writing to a temp and then renaming to the target

	// On some platforms, os.Rename may fail if the system temp directory and the destination are on different
	// volumes, so instead we'll use a temporary file in the current directory

	tempFile, err := ioutil.TempFile(".", toolName) // create in the current directory
	if err != nil {
		return fmt.Errorf("failed to create a temporary file: %v", err)
	}

	defer tempFile.Close()           // okay to close more than once
	defer os.Remove(tempFile.Name()) // okay to delete if not existing (e.g. we successfully renamed)

	if err = generator.Execute(tempFile, rec); err != nil {
		return fmt.Errorf("failed to execute text template: %v", err)
	}

	tempFile.Close() // okay to close more than once

	os.Rename(tempFile.Name(), filename)

	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s: expected a type argument\n", toolName)
		os.Exit(2)
	} else if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "%s: expected a single type argument; called with %v\n", toolName, os.Args[1:])
		os.Exit(2)
	}

	theType := strings.ToLower(os.Args[1])                                   // e.g. 'int'
	sliceType := strings.Join([]string{strings.Title(theType), "Slice"}, "") // e.g. 'IntSlice'

	rec := &record{
		ToolName:      toolName,
		PrimitiveType: theType,
		SliceType:     sliceType,
	}

	if err := generateFile(codeGenerator, rec, fmt.Sprintf("%s_%s.go", toolName, theType)); err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to generate a file: %v\n", toolName, err)
		os.Exit(1)
	}

	if err := generateFile(testGenerator, rec, fmt.Sprintf("%s_%s_test.go", toolName, theType)); err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to generate a file: %v\n", toolName, err)
		os.Exit(1)
	}
}
