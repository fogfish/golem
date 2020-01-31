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
// You can use any type (e.g. `interface{}`) to implement valid, buildable and testable
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
//   src/
//     github.com/fogfish/golem/
//       seq/                     # package generic pattern (e.g. sequence)
//         doc.go                 # documentation of generic pattern and go:generate
//         seq.go                 # generic implementation with parametric type variable
//         seq_test.go            # testing of implementation
//         int.go                 # generated code for built-in type int
//         string.go              # generated code for built-in type string
//         ..
//
// Application
//
// As a application developer I want to parametrise a generic types with my own
// application specific types so that the application benefits from re-use of
// generic implementations
// The mode implies a following rules
//
// ↣ one package defines a type and parametrization of various generic algorithms
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
//   src/
//     github.com/fogfish/myapp/
//       main.go
//       foobar/                  # a custom type definition
//         foobar.go              # type definition and go:generate
//         seq.go                 # parametrization of sequence pattern with foobar type
//         stream.go              # parametrization of stream pattern with foobar type
//         ...
//
// Package names
//
// These workflows and source code structure refers to widely use Go package naming.
//
//   // When we are using a generic library then the name of
//   // generic pattern is visible as package name
//   seq.AnyT{}
//   seq.Int{1, 2, 3, 4}
//   seq.String{"a", "b", "c", "d"}
//
//   // When we are using a custom parametrization then the name of
//   // generic pattern is visible as type name
//   foobar.Seq{FooBar{1}, FooBar{2}, FooBar{3}}
//   foobar.Stream{}
//   foobar.Stack{}
//
// How It works
//
// Use special labels in your code to define generic types:
// mandatory `AnyT`, `generic.T`
//
//   type AnyT struct {
//     t generic.T
//   }
//
//   func foo(x AnyT) generic.T
//   func bar(x generic.T) generic.T
//
//
// The label `generic.T` is replaced with values supplied
// via corresponding command line parameter `-T`.
// The labels `AnyT` is either replaced with a name of type `T` or generic
// pattern depending on the used workflow.  Insert the following comment
// in your source code file:
//
//
//   // Library workflow
//   //go:generate golem -lib -T int -generic github.com/fogfish/golem/seq/seq.go
//
//   // Application workflow
//   //go:generate golem -T FooBar -generic github.com/fogfish/golem/seq/seq.go
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//
type opts struct {
	generic  *string
	genericT *string
	lib      *bool
}

func parseOpts() opts {
	spec := opts{
		flag.String("generic", "", "locates a path to generic type definition."),
		flag.String("T", "", "defines a parametrization to generic.T"),
		flag.Bool("lib", false, "use library declaration schema."),
	}
	flag.Parse()
	return spec
}

func isGenericType(node ast.Node) bool {
	switch spec := node.(type) {
	case *ast.SelectorExpr:
		switch v := spec.X.(type) {
		case *ast.Ident:
			if v.Name == "generic" {
				return true
			}
		}
	}
	return false
}

//
func transform(pkg, tGenT, tAnyT string) func(ast.Node) bool {
	return func(node ast.Node) bool {
		switch spec := node.(type) {
		case *ast.File:
			spec.Name.Name = pkg
		case *ast.ImportSpec:
			if strings.HasSuffix(spec.Path.Value, "golem/generic\"") {
				spec.Path.Value = ""
			}
		case *ast.TypeSpec:
			if isGenericType(spec.Type) {
				spec.Type = &ast.Ident{Name: tGenT}
			}
		case *ast.Field:
			if isGenericType(spec.Type) {
				spec.Type = &ast.Ident{Name: tGenT}
			}
		case *ast.ValueSpec:
			if isGenericType(spec.Type) {
				spec.Type = &ast.Ident{Name: tGenT}
			}
		case *ast.Ident:
			if spec.Name == "AnyT" {
				spec.Name = tAnyT
			}
		}
		return true
	}
}

// The code generator assumes two major development workflows.
//
// Library development.
// One package defines one generic type.
//
// As a generic library developer I want to define a generic type and supply its
// parametrized variants of standard Go type so that my generic is ready for application
// development.
//
// App development
// One package defines a type and parametrization of various generic algorithms
//
// As a application developer I want to parametrize a generic types with my own
// application specific types so that the application benefits from re-use of generic
// implementations.
//
func typeGenericT(source string, opt *opts) string {
	return *opt.genericT
}

func typeAnyT(source string, opt *opts) string {
	generic := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
	typename := strings.Title(generic)
	if *opt.lib {
		typename = strings.Title(*opt.genericT)
	}
	return typename
}

func genericFileName(source string, opt *opts) string {
	generic := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))

	subtype := strings.Split(generic, "_")
	suffix := ""
	if len(subtype) > 1 {
		suffix = "_" + subtype[1]
	}

	filename := fmt.Sprintf("g_%s%s.go", generic, suffix)
	if *opt.lib {
		filename = fmt.Sprintf("g_%s%s.go", *opt.genericT, suffix)
	}
	return strings.ToLower(filename)
}

//
func disclaimer(generic string) *ast.CommentGroup {
	comment := &ast.Comment{
		Text:  fmt.Sprintf("//\n// Code generated by `golem` package\n// Source: %s\n// Time: %s\n//\n\n", generic, time.Now().UTC()),
		Slash: 1,
	}
	return &ast.CommentGroup{
		List: []*ast.Comment{comment},
	}
}

//
func main() {
	log.SetFlags(0)
	log.SetPrefix("==> golem: ")
	opt := parseOpts()

	// detect spec of building package
	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.Fatal(err)
	}

	// generic spec
	source := filepath.Join(build.Default.GOPATH, "src", *opt.generic)
	generic := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))

	// parse generic definition
	fgroup := token.NewFileSet()
	file, err := parser.ParseFile(fgroup, source, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// re-write generic definition
	tGenT := typeGenericT(source, &opt)
	tAnyT := typeAnyT(source, &opt)

	ast.Inspect(file,
		transform(pkg.Name, tGenT, tAnyT))

	// append comment
	doc := disclaimer(*opt.generic)
	file.Doc = doc
	file.Comments = append([]*ast.CommentGroup{doc}, file.Comments...)

	// writes destination file
	filename := filepath.Join(pkg.PkgRoot, genericFileName(source, &opt))
	fd, err := os.Create(filename)
	defer fd.Close()

	printer.Fprint(fd, fgroup, file)
	log.Printf("%s.%s", generic, tAnyT)
}
