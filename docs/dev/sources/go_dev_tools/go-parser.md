# Go Parser Package

> Source: https://pkg.go.dev/go/parser
> Fetched: 2026-02-01T11:54:49.167076+00:00
> Content-Hash: eb13825d8cb3cc80
> Type: html

---

### Overview ¶

Package parser implements a parser for Go source files.

The ParseFile function reads file input from a string, []byte, or io.Reader, and produces an [ast.File](/go/ast#File) representing the complete abstract syntax tree of the file.

The ParseExprFrom function reads a single source-level expression and produces an [ast.Expr](/go/ast#Expr), the syntax tree of the expression.

The parser accepts a larger language than is syntactically permitted by the Go spec, for simplicity, and for improved robustness in the presence of syntax errors. For instance, in method declarations, the receiver is treated like an ordinary parameter list and thus may contain multiple entries where the spec permits exactly one. Consequently, the corresponding field in the AST (ast.FuncDecl.Recv) field is not restricted to one entry.

Applications that need to parse one or more complete packages of Go source code may find it more convenient not to interact directly with the parser but instead to use the Load function in package [golang.org/x/tools/go/packages](/golang.org/x/tools/go/packages).

### Index ¶

- func ParseDir(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode Mode) (pkgs map[string]*ast.Package, first error)deprecated
- func ParseExpr(x string) (ast.Expr, error)
- func ParseExprFrom(fset *token.FileSet, filename string, src any, mode Mode) (expr ast.Expr, err error)
- func ParseFile(fset *token.FileSet, filename string, src any, mode Mode) (f*ast.File, err error)
- type Mode

### Examples ¶

- ParseFile

### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

#### func [ParseDir](https://cs.opensource.google/go/go/+/go1.25.6:src/go/parser/interface.go;l=153) deprecated

    func ParseDir(fset *[token](/go/token).[FileSet](/go/token#FileSet), path [string](/builtin#string), filter func([fs](/io/fs).[FileInfo](/io/fs#FileInfo)) [bool](/builtin#bool), mode Mode) (pkgs map[[string](/builtin#string)]*[ast](/go/ast).[Package](/go/ast#Package), first [error](/builtin#error))

ParseDir calls ParseFile for all files with names ending in ".go" in the directory specified by path and returns a map of package name -> package AST with all the packages found.

If filter != nil, only the files with [fs.FileInfo](/io/fs#FileInfo) entries passing through the filter (and ending in ".go") are considered. The mode bits are passed to ParseFile unchanged. Position information is recorded in fset, which must not be nil.

If the directory couldn't be read, a nil map and the respective error are returned. If a parse error occurred, a non-nil but incomplete map and the first error encountered are returned.

Deprecated: ParseDir does not consider build tags when associating files with packages. For precise information about the relationship between packages and files, use golang.org/x/tools/go/packages, which can also optionally parse and type-check the files too.

#### func [ParseExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/parser/interface.go;l=251) ¶

    func ParseExpr(x [string](/builtin#string)) ([ast](/go/ast).[Expr](/go/ast#Expr), [error](/builtin#error))

ParseExpr is a convenience function for obtaining the AST of an expression x. The position information recorded in the AST is undefined. The filename used in error messages is the empty string.

If syntax errors were found, the result is a partial AST (with [ast.Bad](/go/ast#Bad)* nodes representing the fragments of erroneous source code). Multiple errors are returned via a scanner.ErrorList which is sorted by source position.

#### func [ParseExprFrom](https://cs.opensource.google/go/go/+/go1.25.6:src/go/parser/interface.go;l=203) ¶ added in go1.5

    func ParseExprFrom(fset *[token](/go/token).[FileSet](/go/token#FileSet), filename [string](/builtin#string), src [any](/builtin#any), mode Mode) (expr [ast](/go/ast).[Expr](/go/ast#Expr), err [error](/builtin#error))

ParseExprFrom is a convenience function for parsing an expression. The arguments have the same meaning as for ParseFile, but the source must be a valid Go (type or value) expression. Specifically, fset must not be nil.

If the source couldn't be read, the returned AST is nil and the error indicates the specific failure. If the source was read but syntax errors were found, the result is a partial AST (with [ast.Bad](/go/ast#Bad)* nodes representing the fragments of erroneous source code). Multiple errors are returned via a scanner.ErrorList which is sorted by source position.

#### func [ParseFile](https://cs.opensource.google/go/go/+/go1.25.6:src/go/parser/interface.go;l=84) ¶

    func ParseFile(fset *[token](/go/token).[FileSet](/go/token#FileSet), filename [string](/builtin#string), src [any](/builtin#any), mode Mode) (f *[ast](/go/ast).[File](/go/ast#File), err [error](/builtin#error))

ParseFile parses the source code of a single Go source file and returns the corresponding [ast.File](/go/ast#File) node. The source code may be provided via the filename of the source file, or via the src parameter.

If src != nil, ParseFile parses the source from src and the filename is only used when recording position information. The type of the argument for the src parameter must be string, []byte, or [io.Reader](/io#Reader). If src == nil, ParseFile parses the file specified by filename.

The mode parameter controls the amount of source text parsed and other optional parser functionality. If the SkipObjectResolution mode bit is set (recommended), the object resolution phase of parsing will be skipped, causing File.Scope, File.Unresolved, and all Ident.Obj fields to be nil. Those fields are deprecated; see [ast.Object](/go/ast#Object) for details.

Position information is recorded in the file set fset, which must not be nil.

If the source couldn't be read, the returned AST is nil and the error indicates the specific failure. If the source was read but syntax errors were found, the result is a partial AST (with [ast.Bad](/go/ast#Bad)* nodes representing the fragments of erroneous source code). Multiple errors are returned via a scanner.ErrorList which is sorted by source position.

Example ¶

    package main
    
    import (
     "fmt"
     "go/parser"
     "go/token"
    )
    
    func main() {
     fset := token.NewFileSet() // positions are relative to fset
    
     src := `package foo
    
    import (
     "fmt"
     "time"
    )
    
    func bar() {
     fmt.Println(time.Now())
    }`
    
     // Parse src but stop after processing the imports.
     f, err := parser.ParseFile(fset, "", src, parser.ImportsOnly)
     if err != nil {
      fmt.Println(err)
      return
     }
    
     // Print the imports from the file's AST.
     for _, s := range f.Imports {
      fmt.Println(s.Path.Value)
     }
    
    }
    
    
    
    Output:
    
    
    "fmt"
    "time"
    

Share Format Run

### Types ¶

#### type [Mode](https://cs.opensource.google/go/go/+/go1.25.6:src/go/parser/interface.go;l=47) ¶

    type Mode [uint](/builtin#uint)

A Mode value is a set of flags (or 0). They control the amount of source code parsed and other optional parser functionality.

    const (
     PackageClauseOnly    Mode             = 1 << [iota](/builtin#iota) // stop parsing after package clause
     ImportsOnly                                       // stop parsing after import declarations
     ParseComments                                     // parse comments and add them to AST
     Trace                                             // print a trace of parsed productions
     DeclarationErrors                                 // report declaration errors
     SpuriousErrors                                    // same as AllErrors, for backward-compatibility
     SkipObjectResolution                              // skip deprecated identifier resolution; see ParseFile
     AllErrors            = SpuriousErrors             // report all errors (not just the first 10 on different lines)
    )
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
