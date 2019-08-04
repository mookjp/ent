package gen

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"text/template"
)

//go:generate go-bindata -pkg=gen ./template/...

var (
	// Templates holds the template information for a file that the graph is generating.
	Templates = []struct {
		Name   string
		Format func(*Type) string
	}{
		{
			Name:   "create",
			Format: pkgf("%s_create.go"),
		},
		{
			Name:   "update",
			Format: pkgf("%s_update.go"),
		},
		{
			Name:   "delete",
			Format: pkgf("%s_delete.go"),
		},
		{
			Name:   "query",
			Format: pkgf("%s_query.go"),
		},
		{
			Name:   "model",
			Format: pkgf("%s.go"),
		},
		{
			Name:   "where",
			Format: pkgf("%s/where.go"),
		},
		{
			Name: "meta",
			Format: func(t *Type) string {
				return fmt.Sprintf("%s/%s.go", t.Package(), t.Package())
			},
		},
	}
	// GraphTemplates holds the templates applied on the graph.
	GraphTemplates = []struct {
		Name   string
		Format string
		Skip   func(*Graph) bool
	}{
		{
			Name:   "base",
			Format: "ent.go",
		},
		{
			Name:   "client",
			Format: "client.go",
		},
		{
			Name:   "context",
			Format: "context.go",
		},
		{
			Name:   "tx",
			Format: "tx.go",
		},
		{
			Name:   "config",
			Format: "config.go",
		},
		{
			Name:   "migrate",
			Format: "migrate/migrate.go",
			Skip:   func(g *Graph) bool { return !g.migrateSupport() },
		},
		{
			Name:   "schema",
			Format: "migrate/schema.go",
			Skip:   func(g *Graph) bool { return !g.migrateSupport() },
		},
		{
			Name:   "predicate",
			Format: "predicate/predicate.go",
		},
		{
			Name:   "example",
			Format: "example_test.go",
		},
	}
	// templates holds the Go templates for the code generation.
	// the init function below initializes the templates and its
	// funcs to avoid initialization loop.
	templates = template.New("templates")
	// imports are the import packages used for code generation.
	imports = make(map[string]string)
)

func init() {
	templates.Funcs(funcs)
	for _, asset := range AssetNames() {
		templates = template.Must(templates.Parse(string(MustAsset(asset))))
	}
	b := bytes.NewBuffer([]byte("package main\n"))
	check(templates.ExecuteTemplate(b, "import", Type{}), "load imports")
	f, err := parser.ParseFile(token.NewFileSet(), "", b, parser.ImportsOnly)
	check(err, "parse imports")
	for _, spec := range f.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		check(err, "unquote import path")
		imports[filepath.Base(path)] = path
	}
	for _, s := range drivers {
		for _, path := range s.Imports {
			imports[filepath.Base(path)] = path
		}
	}
}

func pkgf(s string) func(t *Type) string {
	return func(t *Type) string { return fmt.Sprintf(s, t.Package()) }
}