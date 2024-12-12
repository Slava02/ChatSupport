package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const header = "// Code generated by my-super-puper-tool. DO NOT EDIT.\n"

var t = template.Must(template.New("generatorExample").Parse(header + `
package {{.Package}}

import (
	"database/sql/driver"
	"errors"
	"github.com/google/uuid"
)

type TypeSet = interface {
	{{.TypeSet}}
}

func Parse[T TypeSet](s string) (T, error) {
	v, err := uuid.Parse(s)
	return T(v), err
}

func MustParse[T TypeSet](s string) T {
	return T(uuid.MustParse(s))
}

{{ range .Types }} 
	// {{.}} type
	
	type {{.}} uuid.UUID
	
	func New{{.}}() {{.}} {
		return {{.}}(uuid.New())
	}
	
	var {{.}}Nil = {{.}}(uuid.Nil)
	
	func (t {{.}}) String() string { 
		return uuid.UUID(t).String() 
	}
	
	func (t {{.}}) Value() (driver.Value, error) { 
		return t.String(), nil 
	}
	
	func (t *{{.}}) Scan(src any) error { 
		return (*uuid.UUID)(t).Scan(src) 
	}
	
	func (t {{.}}) MarshalText() (text []byte, err error) { 
		return uuid.UUID(t).MarshalText() 
	}
	
	func (t *{{.}}) UnmarshalText(text []byte) error {
		id, err := uuid.ParseBytes(text)
		if err != nil {
			return err
		}
		*t = {{.}}(id)
		return nil
	}
	
	func (t {{.}}) IsZero() bool { 
		return t == {{.}}Nil 
	}
	
	func (t {{.}}) Validate() error {
		if t.IsZero() {
			return errors.New("zero {{.}}")
		}
		return nil
	}
	
	func (t *{{.}}) Matches(x interface{}) bool {
		other, ok := x.({{.}})
		if !ok {
			return false
		}
		return *t == other
	}
{{- end }}
`))

type TemplateData struct {
	Package string
	Types   []string
	TypeSet string
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("invalid args count: %d", len(os.Args)-1)
	}

	pkg, types, out := os.Args[1], strings.Split(os.Args[2], ","), os.Args[3]
	if err := run(pkg, types, out); err != nil {
		log.Fatal(err)
	}

	p, _ := os.Getwd()
	fmt.Printf("%v generated\n", filepath.Join(p, out))
}

func run(pkg string, types []string, outFile string) error {
	data := TemplateData{
		Package: pkg,
		Types:   types,
		TypeSet: strings.Join(types, " | "),
	}

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return fmt.Errorf("execute tmpl: %v", err)
	}

	generated, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("formate source: %v", err)
	}

	if err := os.WriteFile(outFile, generated, 0o644); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}