//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

// Golem is a tool to instantiate a specific type from generic definition.
// The absence of generics in Go causes the usage of `go generate` to re-write
// abstract definition at build time, like this:
//
//   //go:generate golem -T Foo -generic github.com/fogfish/golem/seq/seq.go
//
// The command takes few arguments:
//
//   -T string   defines a parametrization to `generic.T`, required
//   -A string   defines a parametrization to `generic.A`, optional
//   -B string   defines a parametrization to `generic.B`, optional
//   -C string   defines a parametrization to `generic.C`, optional
//
//   -generic path  locates a path to generic algorithm.
//
// The command creates a file in same directory containing a parametrized definition
// of generic type.
//
// Install
//
//   go get -u github.com/fogfish/golem/cmd/golem
//
// Generics
//
// The library uses any type `interface{}` to implement valid, buildable and testable
// generic Go code. Any other language uses a type variables to express generic types,
// e.g. `Stack[T]`. This Go library uses `generic.T` literal as type aliases instead
// of variable for this purpose.
//
//   package stack
//
//   import "github.com/fogfish/golem/generic"
//
//   type AnyT struct {
//	   elements []generic.T
//   }
//
//   func (s AnyT) push(x generic.T) {/* ... */}
//   func (s AnyT) pop() generic.T {/* ... */}
//
// Any one is able to use this generic types directly or its its parametrized version.
//
//   stack.AnyT{}
//   stack.Int{}
//   stack.String{}
//
// The unsafe type definitions are replaced with a specied type, each literal `generic.T`
// and `AnyT` is substitute with value derived from specified type. A few replacement
// modes are supported, these modes takes an advantage of Go package naming schema
// and provides intuitive approach to reference generated types, e.g.
//
//   stack.Int{}     // generics in library
//   foobar.Stack{}  // generics in application
//
// Library
//
// As a generic library developer I want to define a generic type and supply its
// parametrized variants of standard Go type so that my generic is ready for
// application development.
// The mode implies a following rules
//
// ↣ one package defines one generic type.
//
// ↣ concrete types are named after the type, `AnyT` is replaced with `Type`
// (e.g `AnyT` -> `Int`).
//
// ↣ parametrized type value `generic.T` is repaced with defined type
// (e.g `generic.T` -> `int`).
//
// ↣ file type.go is created in the package
// (e.g. `int.go`)
//
// Application
//
// As a application developer I want to parametrise a generic types with my own
// application specific types so that the application benefits from re-use of
// generic implementations
// The mode implies a following rules
//
// ↣ one package implements various generic variants for the custom type
//
// ↣ concrete types are named after the generic, `AnyT` is replaced with `Generic`
// (e.g `AnyT` -> `Stack`).
//
// ↣ type alias `generic.T` is repaced with with defined type
// (e.g `generic.T` -> `FooBar`).
//
// ↣ file generic.go is created in the package
// (e.g. `stack.go`)
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//
type opts struct {
	generic  *string
	genericT *string
	genericA *string
	genericB *string
	genericC *string
	lib      *bool
}

func parseOpts() opts {
	spec := opts{
		flag.String("generic", "", "locates a path to generic type definition."),
		flag.String("T", "", "defines a parametrization to generic.T"),
		flag.String("A", "", "defines a parametrization to generic.A"),
		flag.String("B", "", "defines a parametrization to generic.B"),
		flag.String("C", "", "defines a parametrization to generic.C"),
		flag.Bool("lib", false, "use library declaration schema."),
	}
	flag.Parse()
	return spec
}

//
func declareTypes(x []byte, spec opts) []byte {
	t := declareType(x, "T", *spec.genericT)
	a := declareType(t, "A", *spec.genericA)
	b := declareType(a, "B", *spec.genericB)
	c := declareType(b, "C", *spec.genericC)
	return c
}

func declareType(content []byte, t string, kind string) []byte {
	if kind == "" {
		return content
	}
	return bytes.ReplaceAll(content, []byte("generic."+t), []byte(kind))
}

//
func referenceType(content []byte, kind string) []byte {
	return bytes.ReplaceAll(content,
		[]byte("AnyT"),
		[]byte(kind),
	)
}

//
func repackage(content []byte, pkg string) []byte {
	re := regexp.MustCompile(`package (.*)\n`)
	return re.ReplaceAll(content, []byte("package "+pkg+"\n"))
}

//
func unimport(content []byte) []byte {
	a := bytes.ReplaceAll(content,
		[]byte("import \"github.com/fogfish/golem/generic\"\n"),
		[]byte(""),
	)

	b := bytes.ReplaceAll(a,
		[]byte("\"github.com/fogfish/golem/generic\"\n"),
		[]byte(""),
	)

	re := regexp.MustCompile(`import \(\s+\)\n`)
	return re.ReplaceAll(b, []byte{})
}

//
func main() {
	var err error
	log.SetFlags(0)
	log.SetPrefix("==> golem: ")
	opt := parseOpts()

	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	source := filepath.Join(build.Default.GOPATH, "src", *opt.generic)
	generic := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))

	filename := fmt.Sprintf("%s.go", generic)
	typename := strings.Title(generic)
	if *opt.lib {
		filename = fmt.Sprintf("%s.go", *opt.genericT)
		typename = strings.Title(*opt.genericT)
	}

	input, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatal(err)
	}

	a := declareTypes(input, opt)
	b := referenceType(a, typename)
	c := repackage(b, pkg.Name)
	d := unimport(c)

	output := bytes.NewBuffer([]byte{})
	output.Write([]byte("// Code generated by `golem` package\n"))
	output.Write([]byte(fmt.Sprintf("// Source: %s\n", *opt.generic)))
	output.Write([]byte(fmt.Sprintf("// Time: %s\n\n", time.Now().UTC())))

	output.Write(d)

	ioutil.WriteFile(filepath.Join(pkg.PkgRoot, filename), output.Bytes(), 0777)
	log.Printf("%s.%s", generic, typename)
}
