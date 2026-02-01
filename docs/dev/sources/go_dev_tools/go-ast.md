# Go AST Package

> Source: https://pkg.go.dev/go/ast
> Fetched: 2026-02-01T11:54:46.998417+00:00
> Content-Hash: 174867119c362055
> Type: html

---

### Overview ¶

Package ast declares the types used to represent syntax trees for Go packages.

Syntax trees may be constructed directly, but they are typically produced from Go source code by the parser; see the ParseFile function in package [go/parser](/go/parser).

### Index ¶

- func FileExports(src *File) bool
- func FilterDecl(decl Decl, f Filter) bool
- func FilterFile(src *File, f Filter) bool
- func FilterPackage(pkg *Package, f Filter) booldeprecated
- func Fprint(w io.Writer, fset *token.FileSet, x any, f FieldFilter) error
- func Inspect(node Node, f func(Node) bool)
- func IsExported(name string) bool
- func IsGenerated(file *File) bool
- func NotNilFilter(_ string, v reflect.Value) bool
- func PackageExports(pkg *Package) booldeprecated
- func Preorder(root Node) iter.Seq[Node]
- func PreorderStack(root Node, stack []Node, f func(n Node, stack []Node) bool)
- func Print(fset *token.FileSet, x any) error
- func SortImports(fset *token.FileSet, f*File)
- func Walk(v Visitor, node Node)
- type ArrayType
-     * func (x *ArrayType) End() token.Pos
  - func (x *ArrayType) Pos() token.Pos
- type AssignStmt
-     * func (s *AssignStmt) End() token.Pos
  - func (s *AssignStmt) Pos() token.Pos
- type BadDecl
-     * func (d *BadDecl) End() token.Pos
  - func (d *BadDecl) Pos() token.Pos
- type BadExpr
-     * func (x *BadExpr) End() token.Pos
  - func (x *BadExpr) Pos() token.Pos
- type BadStmt
-     * func (s *BadStmt) End() token.Pos
  - func (s *BadStmt) Pos() token.Pos
- type BasicLit
-     * func (x *BasicLit) End() token.Pos
  - func (x *BasicLit) Pos() token.Pos
- type BinaryExpr
-     * func (x *BinaryExpr) End() token.Pos
  - func (x *BinaryExpr) Pos() token.Pos
- type BlockStmt
-     * func (s *BlockStmt) End() token.Pos
  - func (s *BlockStmt) Pos() token.Pos
- type BranchStmt
-     * func (s *BranchStmt) End() token.Pos
  - func (s *BranchStmt) Pos() token.Pos
- type CallExpr
-     * func (x *CallExpr) End() token.Pos
  - func (x *CallExpr) Pos() token.Pos
- type CaseClause
-     * func (s *CaseClause) End() token.Pos
  - func (s *CaseClause) Pos() token.Pos
- type ChanDir
- type ChanType
-     * func (x *ChanType) End() token.Pos
  - func (x *ChanType) Pos() token.Pos
- type CommClause
-     * func (s *CommClause) End() token.Pos
  - func (s *CommClause) Pos() token.Pos
- type Comment
-     * func (c *Comment) End() token.Pos
  - func (c *Comment) Pos() token.Pos
- type CommentGroup
-     * func (g *CommentGroup) End() token.Pos
  - func (g *CommentGroup) Pos() token.Pos
  - func (g *CommentGroup) Text() string
- type CommentMap
-     * func NewCommentMap(fset *token.FileSet, node Node, comments []*CommentGroup) CommentMap
-     * func (cmap CommentMap) Comments() []*CommentGroup
  - func (cmap CommentMap) Filter(node Node) CommentMap
  - func (cmap CommentMap) String() string
  - func (cmap CommentMap) Update(old, new Node) Node
- type CompositeLit
-     * func (x *CompositeLit) End() token.Pos
  - func (x *CompositeLit) Pos() token.Pos
- type Decl
- type DeclStmt
-     * func (s *DeclStmt) End() token.Pos
  - func (s *DeclStmt) Pos() token.Pos
- type DeferStmt
-     * func (s *DeferStmt) End() token.Pos
  - func (s *DeferStmt) Pos() token.Pos
- type Ellipsis
-     * func (x *Ellipsis) End() token.Pos
  - func (x *Ellipsis) Pos() token.Pos
- type EmptyStmt
-     * func (s *EmptyStmt) End() token.Pos
  - func (s *EmptyStmt) Pos() token.Pos
- type Expr
-     * func Unparen(e Expr) Expr
- type ExprStmt
-     * func (s *ExprStmt) End() token.Pos
  - func (s *ExprStmt) Pos() token.Pos
- type Field
-     * func (f *Field) End() token.Pos
  - func (f *Field) Pos() token.Pos
- type FieldFilter
- type FieldList
-     * func (f *FieldList) End() token.Pos
  - func (f *FieldList) NumFields() int
  - func (f *FieldList) Pos() token.Pos
- type File
-     * func MergePackageFiles(pkg *Package, mode MergeMode) *Filedeprecated
-     * func (f *File) End() token.Pos
  - func (f *File) Pos() token.Pos
- type Filter
- type ForStmt
-     * func (s *ForStmt) End() token.Pos
  - func (s *ForStmt) Pos() token.Pos
- type FuncDecl
-     * func (d *FuncDecl) End() token.Pos
  - func (d *FuncDecl) Pos() token.Pos
- type FuncLit
-     * func (x *FuncLit) End() token.Pos
  - func (x *FuncLit) Pos() token.Pos
- type FuncType
-     * func (x *FuncType) End() token.Pos
  - func (x *FuncType) Pos() token.Pos
- type GenDecl
-     * func (d *GenDecl) End() token.Pos
  - func (d *GenDecl) Pos() token.Pos
- type GoStmt
-     * func (s *GoStmt) End() token.Pos
  - func (s *GoStmt) Pos() token.Pos
- type Ident
-     * func NewIdent(name string) *Ident
-     * func (x *Ident) End() token.Pos
  - func (id *Ident) IsExported() bool
  - func (x *Ident) Pos() token.Pos
  - func (id *Ident) String() string
- type IfStmt
-     * func (s *IfStmt) End() token.Pos
  - func (s *IfStmt) Pos() token.Pos
- type ImportSpec
-     * func (s *ImportSpec) End() token.Pos
  - func (s *ImportSpec) Pos() token.Pos
- type Importerdeprecated
- type IncDecStmt
-     * func (s *IncDecStmt) End() token.Pos
  - func (s *IncDecStmt) Pos() token.Pos
- type IndexExpr
-     * func (x *IndexExpr) End() token.Pos
  - func (x *IndexExpr) Pos() token.Pos
- type IndexListExpr
-     * func (x *IndexListExpr) End() token.Pos
  - func (x *IndexListExpr) Pos() token.Pos
- type InterfaceType
-     * func (x *InterfaceType) End() token.Pos
  - func (x *InterfaceType) Pos() token.Pos
- type KeyValueExpr
-     * func (x *KeyValueExpr) End() token.Pos
  - func (x *KeyValueExpr) Pos() token.Pos
- type LabeledStmt
-     * func (s *LabeledStmt) End() token.Pos
  - func (s *LabeledStmt) Pos() token.Pos
- type MapType
-     * func (x *MapType) End() token.Pos
  - func (x *MapType) Pos() token.Pos
- type MergeModedeprecated
- type Node
- type ObjKind
-     * func (kind ObjKind) String() string
- type Objectdeprecated
-     * func NewObj(kind ObjKind, name string) *Object
-     * func (obj *Object) Pos() token.Pos
- type Packagedeprecated
-     * func NewPackage(fset *token.FileSet, files map[string]*File, importer Importer, ...) (*Package, error)deprecated
-     * func (p *Package) End() token.Pos
  - func (p *Package) Pos() token.Pos
- type ParenExpr
-     * func (x *ParenExpr) End() token.Pos
  - func (x *ParenExpr) Pos() token.Pos
- type RangeStmt
-     * func (s *RangeStmt) End() token.Pos
  - func (s *RangeStmt) Pos() token.Pos
- type ReturnStmt
-     * func (s *ReturnStmt) End() token.Pos
  - func (s *ReturnStmt) Pos() token.Pos
- type Scopedeprecated
-     * func NewScope(outer *Scope) *Scope
-     * func (s *Scope) Insert(obj *Object) (alt *Object)
  - func (s *Scope) Lookup(name string)*Object
  - func (s *Scope) String() string
- type SelectStmt
-     * func (s *SelectStmt) End() token.Pos
  - func (s *SelectStmt) Pos() token.Pos
- type SelectorExpr
-     * func (x *SelectorExpr) End() token.Pos
  - func (x *SelectorExpr) Pos() token.Pos
- type SendStmt
-     * func (s *SendStmt) End() token.Pos
  - func (s *SendStmt) Pos() token.Pos
- type SliceExpr
-     * func (x *SliceExpr) End() token.Pos
  - func (x *SliceExpr) Pos() token.Pos
- type Spec
- type StarExpr
-     * func (x *StarExpr) End() token.Pos
  - func (x *StarExpr) Pos() token.Pos
- type Stmt
- type StructType
-     * func (x *StructType) End() token.Pos
  - func (x *StructType) Pos() token.Pos
- type SwitchStmt
-     * func (s *SwitchStmt) End() token.Pos
  - func (s *SwitchStmt) Pos() token.Pos
- type TypeAssertExpr
-     * func (x *TypeAssertExpr) End() token.Pos
  - func (x *TypeAssertExpr) Pos() token.Pos
- type TypeSpec
-     * func (s *TypeSpec) End() token.Pos
  - func (s *TypeSpec) Pos() token.Pos
- type TypeSwitchStmt
-     * func (s *TypeSwitchStmt) End() token.Pos
  - func (s *TypeSwitchStmt) Pos() token.Pos
- type UnaryExpr
-     * func (x *UnaryExpr) End() token.Pos
  - func (x *UnaryExpr) Pos() token.Pos
- type ValueSpec
-     * func (s *ValueSpec) End() token.Pos
  - func (s *ValueSpec) Pos() token.Pos
- type Visitor

### Examples ¶

- CommentMap
- Inspect
- Preorder
- Print

### Constants ¶

This section is empty.

### Variables ¶

This section is empty.

### Functions ¶

#### func [FileExports](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=27) ¶

    func FileExports(src *File) [bool](/builtin#bool)

FileExports trims the AST for a Go source file in place such that only exported nodes remain: all top-level identifiers which are not exported and their associated information (such as type, initial value, or function body) are removed. Non-exported fields and methods of exported types are stripped. The [File.Comments] list is not changed.

FileExports reports whether there are exported declarations.

#### func [FilterDecl](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=232) ¶

    func FilterDecl(decl Decl, f Filter) [bool](/builtin#bool)

FilterDecl trims the AST for a Go declaration in place by removing all names (including struct field and interface method names, but not from parameter lists) that don't pass through the filter f.

FilterDecl reports whether there are any declared names left after filtering.

#### func [FilterFile](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=256) ¶

    func FilterFile(src *File, f Filter) [bool](/builtin#bool)

FilterFile trims the AST for a Go file in place by removing all names from top-level declarations (including struct field and interface method names, but not from parameter lists) that don't pass through the filter f. If the declaration is empty afterwards, the declaration is removed from the AST. Import declarations are always removed. The [File.Comments] list is not changed.

FilterFile reports whether there are any top-level declarations left after filtering.

#### func [FilterPackage](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=285) deprecated

    func FilterPackage(pkg *Package, f Filter) [bool](/builtin#bool)

FilterPackage trims the AST for a Go package in place by removing all names from top-level declarations (including struct field and interface method names, but not from parameter lists) that don't pass through the filter f. If the declaration is empty afterwards, the declaration is removed from the AST. The pkg.Files list is not changed, so that file names and top-level package comments don't get lost.

FilterPackage reports whether there are any top-level declarations left after filtering.

Deprecated: use the type checker [go/types](/go/types) instead of Package; see Object. Alternatively, use FilterFile.

#### func [Fprint](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/print.go;l=39) ¶

    func Fprint(w [io](/io).[Writer](/io#Writer), fset *[token](/go/token).[FileSet](/go/token#FileSet), x [any](/builtin#any), f FieldFilter) [error](/builtin#error)

Fprint prints the (sub-)tree starting at AST node x to w. If fset != nil, position information is interpreted relative to that file set. Otherwise positions are printed as integer values (file set specific offsets).

A non-nil FieldFilter f may be provided to control the output: struct fields for which f(fieldname, fieldvalue) is true are printed; all others are filtered from the output. Unexported struct fields are never printed.

#### func [Inspect](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/walk.go;l=376) ¶

    func Inspect(node Node, f func(Node) [bool](/builtin#bool))

Inspect traverses an AST in depth-first order: It starts by calling f(node); node must not be nil. If f returns true, Inspect invokes f recursively for each of the non-nil children of node, followed by a call of f(nil).

In many cases it may be more convenient to use Preorder, which returns an iterator over the sqeuence of nodes, or PreorderStack, which (like Inspect) provides control over descent into subtrees, but additionally reports the stack of enclosing nodes.

Example ¶

This example demonstrates how to inspect the AST of a Go program.

    package main
    
    import (
     "fmt"
     "go/ast"
     "go/parser"
     "go/token"
    )
    
    func main() {
     // src is the input for which we want to inspect the AST.
     src := `
    package p
    const c = 1.0
    var X = f(3.14)*2 + c
    `
    
     // Create the AST by parsing src.
     fset := token.NewFileSet() // positions are relative to fset
     f, err := parser.ParseFile(fset, "src.go", src, 0)
     if err != nil {
      panic(err)
     }
    
     // Inspect the AST and print all identifiers and literals.
     ast.Inspect(f, func(n ast.Node) bool {
      var s string
      switch x := n.(type) {
      case *ast.BasicLit:
       s = x.Value
      case *ast.Ident:
       s = x.Name
      }
      if s != "" {
       fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
      }
      return true
     })
    
    }
    
    
    
    Output:
    
    src.go:2:9: p
    src.go:3:7: c
    src.go:3:11: 1.0
    src.go:4:5: X
    src.go:4:9: f
    src.go:4:11: 3.14
    src.go:4:17: 2
    src.go:4:21: c
    

Share Format Run

#### func [IsExported](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=599) ¶

    func IsExported(name [string](/builtin#string)) [bool](/builtin#bool)

IsExported reports whether name starts with an upper-case letter.

#### func [IsGenerated](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1112) ¶ added in go1.21.0

    func IsGenerated(file *File) [bool](/builtin#bool)

IsGenerated reports whether the file was generated by a program, not handwritten, by detecting the special comment described at <https://go.dev/s/generatedcode>.

The syntax tree must have been parsed with the [parser.ParseComments] flag. Example:

    f, err := parser.ParseFile(fset, filename, src, parser.ParseComments|parser.PackageClauseOnly)
    if err != nil { ... }
    gen := ast.IsGenerated(f)
    

#### func [NotNilFilter](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/print.go;l=22) ¶

    func NotNilFilter(_ [string](/builtin#string), v [reflect](/reflect).[Value](/reflect#Value)) [bool](/builtin#bool)

NotNilFilter is a FieldFilter that returns true for field values that are not nil; it returns false otherwise.

#### func [PackageExports](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=40) deprecated

    func PackageExports(pkg *Package) [bool](/builtin#bool)

PackageExports trims the AST for a Go package in place such that only exported nodes remain. The pkg.Files list is not changed, so that file names and top-level package comments don't get lost.

PackageExports reports whether there are exported declarations; it returns false otherwise.

Deprecated: use the type checker [go/types](/go/types) instead of Package; see Object. Alternatively, use FileExports.

#### func [Preorder](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/walk.go;l=386) ¶ added in go1.23.0

    func Preorder(root Node) [iter](/iter).[Seq](/iter#Seq)[Node]

Preorder returns an iterator over all the nodes of the syntax tree beneath (and including) the specified root, in depth-first preorder.

For greater control over the traversal of each subtree, use Inspect or PreorderStack.

Example ¶

    package main
    
    import (
     "fmt"
     "go/ast"
     "go/parser"
     "go/token"
    )
    
    func main() {
     src := `
    package p
    
    func f(x, y int) {
     print(x + y)
    }
    `
    
     fset := token.NewFileSet()
     f, err := parser.ParseFile(fset, "", src, 0)
     if err != nil {
      panic(err)
     }
    
     // Print identifiers in order
     for n := range ast.Preorder(f) {
      id, ok := n.(*ast.Ident)
      if !ok {
       continue
      }
      fmt.Println(id.Name)
     }
    
    }
    
    
    
    Output:
    
    p
    f
    x
    y
    int
    print
    x
    y
    

Share Format Run

#### func [PreorderStack](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/walk.go;l=411) ¶ added in go1.25.0

    func PreorderStack(root Node, stack []Node, f func(n Node, stack []Node) [bool](/builtin#bool))

PreorderStack traverses the tree rooted at root, calling f before visiting each node.

Each call to f provides the current node and traversal stack, consisting of the original value of stack appended with all nodes from root to n, excluding n itself. (This design allows calls to PreorderStack to be nested without double counting.)

If f returns false, the traversal skips over that subtree. Unlike Inspect, no second call to f is made after visiting node n. (In practice, the second call is nearly always used only to pop the stack, and it is surprisingly tricky to do this correctly.)

#### func [Print](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/print.go;l=73) ¶

    func Print(fset *[token](/go/token).[FileSet](/go/token#FileSet), x [any](/builtin#any)) [error](/builtin#error)

Print prints x to standard output, skipping nil fields. Print(fset, x) is the same as Fprint(os.Stdout, fset, x, NotNilFilter).

Example ¶

This example shows what an AST looks like when printed for debugging.

    package main
    
    import (
     "go/ast"
     "go/parser"
     "go/token"
    )
    
    func main() {
     // src is the input for which we want to print the AST.
     src := `
    package main
    func main() {
     println("Hello, World!")
    }
    `
    
     // Create the AST by parsing src.
     fset := token.NewFileSet() // positions are relative to fset
     f, err := parser.ParseFile(fset, "", src, 0)
     if err != nil {
      panic(err)
     }
    
     // Print the AST.
     ast.Print(fset, f)
    
    }
    
    
    
    Output:
    
         0  *ast.File {
         1  .  Package: 2:1
         2  .  Name: *ast.Ident {
         3  .  .  NamePos: 2:9
         4  .  .  Name: "main"
         5  .  }
         6  .  Decls: []ast.Decl (len = 1) {
         7  .  .  0: *ast.FuncDecl {
         8  .  .  .  Name: *ast.Ident {
         9  .  .  .  .  NamePos: 3:6
        10  .  .  .  .  Name: "main"
        11  .  .  .  .  Obj: *ast.Object {
        12  .  .  .  .  .  Kind: func
        13  .  .  .  .  .  Name: "main"
        14  .  .  .  .  .  Decl: *(obj @ 7)
        15  .  .  .  .  }
        16  .  .  .  }
        17  .  .  .  Type: *ast.FuncType {
        18  .  .  .  .  Func: 3:1
        19  .  .  .  .  Params: *ast.FieldList {
        20  .  .  .  .  .  Opening: 3:10
        21  .  .  .  .  .  Closing: 3:11
        22  .  .  .  .  }
        23  .  .  .  }
        24  .  .  .  Body: *ast.BlockStmt {
        25  .  .  .  .  Lbrace: 3:13
        26  .  .  .  .  List: []ast.Stmt (len = 1) {
        27  .  .  .  .  .  0: *ast.ExprStmt {
        28  .  .  .  .  .  .  X: *ast.CallExpr {
        29  .  .  .  .  .  .  .  Fun: *ast.Ident {
        30  .  .  .  .  .  .  .  .  NamePos: 4:2
        31  .  .  .  .  .  .  .  .  Name: "println"
        32  .  .  .  .  .  .  .  }
        33  .  .  .  .  .  .  .  Lparen: 4:9
        34  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
        35  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
        36  .  .  .  .  .  .  .  .  .  ValuePos: 4:10
        37  .  .  .  .  .  .  .  .  .  Kind: STRING
        38  .  .  .  .  .  .  .  .  .  Value: "\"Hello, World!\""
        39  .  .  .  .  .  .  .  .  }
        40  .  .  .  .  .  .  .  }
        41  .  .  .  .  .  .  .  Ellipsis: -
        42  .  .  .  .  .  .  .  Rparen: 4:25
        43  .  .  .  .  .  .  }
        44  .  .  .  .  .  }
        45  .  .  .  .  }
        46  .  .  .  .  Rbrace: 5:1
        47  .  .  .  }
        48  .  .  }
        49  .  }
        50  .  FileStart: 1:1
        51  .  FileEnd: 5:3
        52  .  Scope: *ast.Scope {
        53  .  .  Objects: map[string]*ast.Object (len = 1) {
        54  .  .  .  "main": *(obj @ 11)
        55  .  .  }
        56  .  }
        57  .  Unresolved: []*ast.Ident (len = 1) {
        58  .  .  0: *(obj @ 29)
        59  .  }
        60  .  GoVersion: ""
        61  }
    

Share Format Run

#### func [SortImports](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/import.go;l=16) ¶

    func SortImports(fset *[token](/go/token).[FileSet](/go/token#FileSet), f *File)

SortImports sorts runs of consecutive import lines in import blocks in f. It also removes duplicate imports when it is possible to do so without data loss.

#### func [Walk](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/walk.go;l=33) ¶

    func Walk(v Visitor, node Node)

Walk traverses an AST in depth-first order: It starts by calling v.Visit(node); node must not be nil. If the visitor w returned by v.Visit(node) is not nil, Walk is invoked recursively with visitor w for each of the non-nil children of node, followed by a call of w.Visit(nil).

### Types ¶

#### type [ArrayType](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=448) ¶

    type ArrayType struct {
     Lbrack [token](/go/token).[Pos](/go/token#Pos) // position of "["
     Len    Expr      // Ellipsis node for [...]T array types, nil for slice types
     Elt    Expr      // element type
    }

An ArrayType node represents an array or slice type.

#### func (*ArrayType) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=552) ¶

    func (x *ArrayType) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ArrayType) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=518) ¶

    func (x *ArrayType) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [AssignStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=670) ¶

    type AssignStmt struct {
     Lhs    []Expr
     TokPos [token](/go/token).[Pos](/go/token#Pos)   // position of Tok
     Tok    [token](/go/token).[Token](/go/token#Token) // assignment token, DEFINE
     Rhs    []Expr
    }

An AssignStmt node represents an assignment or a short variable declaration.

#### func (*AssignStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=817) ¶

    func (s *AssignStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*AssignStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=788) ¶

    func (s *AssignStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BadDecl](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=973) ¶

    type BadDecl struct {
     From, To [token](/go/token).[Pos](/go/token#Pos) // position range of bad declaration
    }

A BadDecl node is a placeholder for a declaration containing syntax errors for which a correct declaration node cannot be created.

#### func (*BadDecl) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1013) ¶

    func (d *BadDecl) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BadDecl) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1009) ¶

    func (d *BadDecl) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BadExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=286) ¶

    type BadExpr struct {
     From, To [token](/go/token).[Pos](/go/token#Pos) // position range of bad expression
    }

A BadExpr node is a placeholder for an expression containing syntax errors for which a correct expression node cannot be created.

#### func (*BadExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=530) ¶

    func (x *BadExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BadExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=496) ¶

    func (x *BadExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BadStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=621) ¶

    type BadStmt struct {
     From, To [token](/go/token).[Pos](/go/token#Pos) // position range of bad statement
    }

A BadStmt node is a placeholder for statements containing syntax errors for which no correct statement nodes can be created.

#### func (*BadStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=803) ¶

    func (s *BadStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BadStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=781) ¶

    func (s *BadStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BasicLit](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=318) ¶

    type BasicLit struct {
     ValuePos [token](/go/token).[Pos](/go/token#Pos)   // literal position
     Kind     [token](/go/token).[Token](/go/token#Token) // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
     Value    [string](/builtin#string)      // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
    }

A BasicLit node represents a literal of basic type.

Note that for the CHAR and STRING kinds, the literal is stored with its quotes. For example, for a double-quoted STRING, the first and the last rune in the Value field will be ". The [strconv.Unquote](/strconv#Unquote) and [strconv.UnquoteChar](/strconv#UnquoteChar) functions can be used to unquote STRING and CHAR values, respectively.

For raw string literals (Kind == token.STRING && Value[0] == '`'), the Value field contains the string text without carriage returns (\r) that may have been present in the source. Because the end position is computed using len(Value), the position reported by BasicLit.End does not match the true source end position for raw string literals containing carriage returns.

#### func (*BasicLit) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=538) ¶

    func (x *BasicLit) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BasicLit) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=499) ¶

    func (x *BasicLit) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BinaryExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=417) ¶

    type BinaryExpr struct {
     X     Expr        // left operand
     OpPos [token](/go/token).[Pos](/go/token#Pos)   // position of Op
     Op    [token](/go/token).[Token](/go/token#Token) // operator
     Y     Expr        // right operand
    }

A BinaryExpr node represents a binary expression.

#### func (*BinaryExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=550) ¶

    func (x *BinaryExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BinaryExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=516) ¶

    func (x *BinaryExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BlockStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=705) ¶

    type BlockStmt struct {
     Lbrace [token](/go/token).[Pos](/go/token#Pos) // position of "{"
     List   []Stmt
     Rbrace [token](/go/token).[Pos](/go/token#Pos) // position of "}", if any (may be absent due to syntax error)
    }

A BlockStmt node represents a braced statement list.

#### func (*BlockStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=832) ¶

    func (s *BlockStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BlockStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=793) ¶

    func (s *BlockStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [BranchStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=698) ¶

    type BranchStmt struct {
     TokPos [token](/go/token).[Pos](/go/token#Pos)   // position of Tok
     Tok    [token](/go/token).[Token](/go/token#Token) // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
     Label  *Ident      // label name; or nil
    }

A BranchStmt node represents a break, continue, goto, or fallthrough statement.

#### func (*BranchStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=826) ¶

    func (s *BranchStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*BranchStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=792) ¶

    func (s *BranchStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [CallExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=391) ¶

    type CallExpr struct {
     Fun      Expr      // function expression
     Lparen   [token](/go/token).[Pos](/go/token#Pos) // position of "("
     Args     []Expr    // function arguments; or nil
     Ellipsis [token](/go/token).[Pos](/go/token#Pos) // position of "..." (token.NoPos if there is no "...")
     Rparen   [token](/go/token).[Pos](/go/token#Pos) // position of ")"
    }

A CallExpr node represents an expression followed by an argument list.

#### func (*CallExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=547) ¶

    func (x *CallExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*CallExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=513) ¶

    func (x *CallExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [CaseClause](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=721) ¶

    type CaseClause struct {
     Case  [token](/go/token).[Pos](/go/token#Pos) // position of "case" or "default" keyword
     List  []Expr    // list of expressions or types; nil means default case
     Colon [token](/go/token).[Pos](/go/token#Pos) // position of ":"
     Body  []Stmt    // statement list; or nil
    }

A CaseClause represents a case of an expression or type switch statement.

#### func (*CaseClause) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=847) ¶

    func (s *CaseClause) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*CaseClause) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=795) ¶

    func (s *CaseClause) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [ChanDir](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=436) ¶

    type ChanDir [int](/builtin#int)

The direction of a channel type is indicated by a bit mask including one or both of the following constants.

    const (
     SEND ChanDir = 1 << [iota](/builtin#iota)
     RECV
    )

#### type [ChanType](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=486) ¶

    type ChanType struct {
     Begin [token](/go/token).[Pos](/go/token#Pos) // position of "chan" keyword or "<-" (whichever comes first)
     Arrow [token](/go/token).[Pos](/go/token#Pos) // position of "<-" (token.NoPos if there is no "<-")
     Dir   ChanDir   // channel direction
     Value Expr      // value type
    }

A ChanType node represents a channel type.

#### func (*ChanType) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=562) ¶

    func (x *ChanType) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ChanType) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=528) ¶

    func (x *ChanType) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [CommClause](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=745) ¶

    type CommClause struct {
     Case  [token](/go/token).[Pos](/go/token#Pos) // position of "case" or "default" keyword
     Comm  Stmt      // send or receive statement; nil means default case
     Colon [token](/go/token).[Pos](/go/token#Pos) // position of ":"
     Body  []Stmt    // statement list; or nil
    }

A CommClause node represents a case of a select statement.

#### func (*CommClause) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=855) ¶

    func (s *CommClause) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*CommClause) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=798) ¶

    func (s *CommClause) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Comment](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=68) ¶

    type Comment struct {
     Slash [token](/go/token).[Pos](/go/token#Pos) // position of "/" starting the comment
     Text  [string](/builtin#string)    // comment text (excluding '\n' for //-style comments)
    }

A Comment node represents a single //-style or /*-style comment.

The Text field contains the comment text without carriage returns (\r) that may have been present in the source. Because a comment's end position is computed using len(Text), the position reported by Comment.End does not match the true source end position for comments containing carriage returns.

#### func (*Comment) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=74) ¶

    func (c *Comment) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*Comment) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=73) ¶

    func (c *Comment) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [CommentGroup](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=78) ¶

    type CommentGroup struct {
     List []*Comment // len(List) > 0
    }

A CommentGroup represents a sequence of comments with no other tokens and no empty lines between.

#### func (*CommentGroup) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=83) ¶

    func (g *CommentGroup) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*CommentGroup) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=82) ¶

    func (g *CommentGroup) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### func (*CommentGroup) [Text](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=101) ¶

    func (g *CommentGroup) Text() [string](/builtin#string)

Text returns the text of the comment. Comment markers (//, /*, and*/), the first space of a line comment, and leading and trailing empty lines are removed. Comment directives like "//line" and "//go:noinline" are also removed. Multiple empty lines are reduced to one, and trailing space on lines is trimmed. Unless the result is empty, it is newline-terminated.

#### type [CommentMap](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/commentmap.go;l=26) ¶ added in go1.1

    type CommentMap map[Node][]*CommentGroup

A CommentMap maps an AST node to a list of comment groups associated with it. See NewCommentMap for a description of the association.

Example ¶

This example illustrates how to remove a variable declaration in a Go program while maintaining correct comment association using an ast.CommentMap.

    package main
    
    import (
     "fmt"
     "go/ast"
     "go/format"
     "go/parser"
     "go/token"
     "strings"
    )
    
    func main() {
     // src is the input for which we create the AST that we
     // are going to manipulate.
     src := `
    // This is the package comment.
    package main
    
    // This comment is associated with the hello constant.
    const hello = "Hello, World!" // line comment 1
    
    // This comment is associated with the foo variable.
    var foo = hello // line comment 2
    
    // This comment is associated with the main function.
    func main() {
     fmt.Println(hello) // line comment 3
    }
    `
    
     // Create the AST by parsing src.
     fset := token.NewFileSet() // positions are relative to fset
     f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
     if err != nil {
      panic(err)
     }
    
     // Create an ast.CommentMap from the ast.File's comments.
     // This helps keeping the association between comments
     // and AST nodes.
     cmap := ast.NewCommentMap(fset, f, f.Comments)
    
     // Remove the first variable declaration from the list of declarations.
     for i, decl := range f.Decls {
      if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.VAR {
       copy(f.Decls[i:], f.Decls[i+1:])
       f.Decls = f.Decls[:len(f.Decls)-1]
       break
      }
     }
    
     // Use the comment map to filter comments that don't belong anymore
     // (the comments associated with the variable declaration), and create
     // the new comments list.
     f.Comments = cmap.Filter(f).Comments()
    
     // Print the modified AST.
     var buf strings.Builder
     if err := format.Node(&buf, fset, f); err != nil {
      panic(err)
     }
     fmt.Printf("%s", buf.String())
    
    }
    
    
    
    Output:
    
    // This is the package comment.
    package main
    
    // This comment is associated with the hello constant.
    const hello = "Hello, World!" // line comment 1
    
    // This comment is associated with the main function.
    func main() {
     fmt.Println(hello) // line comment 3
    }
    

Share Format Run

#### func [NewCommentMap](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/commentmap.go;l=127) ¶ added in go1.1

    func NewCommentMap(fset *[token](/go/token).[FileSet](/go/token#FileSet), node Node, comments []*CommentGroup) CommentMap

NewCommentMap creates a new comment map by associating comment groups of the comments list with the nodes of the AST specified by node.

A comment group g is associated with a node n if:

- g starts on the same line as n ends
- g starts on the line immediately following n, and there is at least one empty line after g and before the next node
- g starts before n and is not associated to the node before n via the previous rules

NewCommentMap tries to associate a comment group to the "largest" node possible: For instance, if the comment is a line comment trailing an assignment, the comment is associated with the entire assignment rather than just the last operand in the assignment.

#### func (CommentMap) [Comments](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/commentmap.go;l=253) ¶ added in go1.1

    func (cmap CommentMap) Comments() []*CommentGroup

Comments returns the list of comment groups in the comment map. The result is sorted in source order.

#### func (CommentMap) [Filter](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/commentmap.go;l=240) ¶ added in go1.1

    func (cmap CommentMap) Filter(node Node) CommentMap

Filter returns a new comment map consisting of only those entries of cmap for which a corresponding node exists in the AST specified by node.

#### func (CommentMap) [String](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/commentmap.go;l=298) ¶ added in go1.1

    func (cmap CommentMap) String() [string](/builtin#string)

#### func (CommentMap) [Update](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/commentmap.go;l=229) ¶ added in go1.1

    func (cmap CommentMap) Update(old, new Node) Node

Update replaces an old node in the comment map with the new node and returns the new node. Comments that were associated with the old node are associated with the new node.

#### type [CompositeLit](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=331) ¶

    type CompositeLit struct {
     Type       Expr      // literal type; or nil
     Lbrace     [token](/go/token).[Pos](/go/token#Pos) // position of "{"
     Elts       []Expr    // list of composite elements; or nil
     Rbrace     [token](/go/token).[Pos](/go/token#Pos) // position of "}"
     Incomplete [bool](/builtin#bool)      // true if (source) expressions are missing in the Elts list
    }

A CompositeLit node represents a composite literal.

#### func (*CompositeLit) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=540) ¶

    func (x *CompositeLit) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*CompositeLit) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=501) ¶

    func (x *CompositeLit) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Decl](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=54) ¶

    type Decl interface {
     Node
     // contains filtered or unexported methods
    }

All declaration nodes implement the Decl interface.

#### type [DeclStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=626) ¶

    type DeclStmt struct {
     Decl Decl // *GenDecl with CONST, TYPE, or VAR token
    }

A DeclStmt node represents a declaration in a statement list.

#### func (*DeclStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=804) ¶

    func (s *DeclStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*DeclStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=782) ¶

    func (s *DeclStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [DeferStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=684) ¶

    type DeferStmt struct {
     Defer [token](/go/token).[Pos](/go/token#Pos) // position of "defer" keyword
     Call  *CallExpr
    }

A DeferStmt node represents a defer statement.

#### func (*DeferStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=819) ¶

    func (s *DeferStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*DeferStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=790) ¶

    func (s *DeferStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Ellipsis](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=300) ¶

    type Ellipsis struct {
     Ellipsis [token](/go/token).[Pos](/go/token#Pos) // position of "..."
     Elt      Expr      // ellipsis element type (parameter lists only); or nil
    }

An Ellipsis node stands for the "..." type in a parameter list or the "..." length in an array type.

#### func (*Ellipsis) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=532) ¶

    func (x *Ellipsis) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*Ellipsis) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=498) ¶

    func (x *Ellipsis) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [EmptyStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=634) ¶

    type EmptyStmt struct {
     Semicolon [token](/go/token).[Pos](/go/token#Pos) // position of following ";"
     Implicit  [bool](/builtin#bool)      // if set, ";" was omitted in the source
    }

An EmptyStmt node represents an empty statement. The "position" of the empty statement is the position of the immediately following (explicit or implicit) semicolon.

#### func (*EmptyStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=805) ¶

    func (s *EmptyStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*EmptyStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=783) ¶

    func (s *EmptyStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Expr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=42) ¶

    type Expr interface {
     Node
     // contains filtered or unexported methods
    }

All expression nodes implement the Expr interface.

#### func [Unparen](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1140) ¶ added in go1.22.0

    func Unparen(e Expr) Expr

Unparen returns the expression with any enclosing parentheses removed.

#### type [ExprStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=649) ¶

    type ExprStmt struct {
     X Expr // expression
    }

An ExprStmt node represents a (stand-alone) expression in a statement list.

#### func (*ExprStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=812) ¶

    func (s *ExprStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ExprStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=785) ¶

    func (s *ExprStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Field](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=201) ¶

    type Field struct {
     Doc     *CommentGroup // associated documentation; or nil
     Names   []*Ident      // field/method/(type) parameter names; or nil
     Type    Expr          // field/method/parameter type; or nil
     Tag     *BasicLit     // field tag; or nil
     Comment *CommentGroup // line comments; or nil
    }

A Field represents a Field declaration list in a struct type, a method list in an interface type, or a parameter/result declaration in a signature. [Field.Names] is nil for unnamed parameters (parameter lists which only contain types) and embedded struct fields. In the latter case, the field name is the type name.

#### func (*Field) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=219) ¶

    func (f *Field) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*Field) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=209) ¶

    func (f *Field) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [FieldFilter](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/print.go;l=18) ¶

    type FieldFilter func(name [string](/builtin#string), value [reflect](/reflect).[Value](/reflect#Value)) [bool](/builtin#bool)

A FieldFilter may be provided to Fprint to control the output.

#### type [FieldList](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=234) ¶

    type FieldList struct {
     Opening [token](/go/token).[Pos](/go/token#Pos) // position of opening parenthesis/brace/bracket, if any
     List    []*Field  // field list; or nil
     Closing [token](/go/token).[Pos](/go/token#Pos) // position of closing parenthesis/brace/bracket, if any
    }

A FieldList represents a list of Fields, enclosed by parentheses, curly braces, or square brackets.

#### func (*FieldList) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=252) ¶

    func (f *FieldList) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*FieldList) [NumFields](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=265) ¶

    func (f *FieldList) NumFields() [int](/builtin#int)

NumFields returns the number of parameters or struct fields represented by a FieldList.

#### func (*FieldList) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=240) ¶

    func (f *FieldList) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [File](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1057) ¶

    type File struct {
     Doc     *CommentGroup // associated documentation; or nil
     Package [token](/go/token).[Pos](/go/token#Pos)     // position of "package" keyword
     Name    *Ident        // package name
     Decls   []Decl        // top-level declarations; or nil
    
     FileStart, FileEnd [token](/go/token).[Pos](/go/token#Pos)       // start and end of entire file
     Scope              *Scope          // package scope (this file only). Deprecated: see Object
     Imports            []*ImportSpec   // imports in this file
     Unresolved         []*Ident        // unresolved identifiers in this file. Deprecated: see Object
     Comments           []*CommentGroup // list of all comments in the source file
     GoVersion          [string](/builtin#string)          // minimum Go version required by //go:build or // +build directives
    }

A File node represents a Go source file.

The Comments list contains all comments in the source file in order of appearance, including the comments that are pointed to from other nodes via Doc and Comment fields.

For correct printing of source code containing comments (using packages go/format and go/printer), special care must be taken to update comments when a File's syntax tree is modified: For printing, comments are interspersed between tokens based on their position. If syntax tree nodes are removed or moved, relevant comments in their vicinity must also be removed (from the [File.Comments] list) or moved accordingly (by updating their positions). A CommentMap may be used to facilitate some of these operations.

Whether and how a comment is associated with a node depends on the interpretation of the syntax tree by the manipulating program: except for Doc and Comment comments directly associated with nodes, the remaining comments are "free-floating" (see also issues [#18593](https://go.dev/issue/18593), [#20744](https://go.dev/issue/20744)).

#### func [MergePackageFiles](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=349) deprecated

    func MergePackageFiles(pkg *Package, mode MergeMode) *File

MergePackageFiles creates a file AST by merging the ASTs of the files belonging to a package. The mode flags control merging behavior.

Deprecated: this function is poorly specified and has unfixable bugs; also Package is deprecated.

#### func (*File) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1081) ¶

    func (f *File) End() [token](/go/token).[Pos](/go/token#Pos)

End returns the end of the last declaration in the file. It may be invalid, for example in an empty file.

(Use FileEnd for the end of the entire file. It is always valid.)

#### func (*File) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1075) ¶

    func (f *File) Pos() [token](/go/token).[Pos](/go/token#Pos)

Pos returns the position of the package declaration. It may be invalid, for example in an empty file.

(Use FileStart for the start of the entire file. It is always valid.)

#### type [Filter](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=47) ¶

    type Filter func([string](/builtin#string)) [bool](/builtin#bool)

#### type [ForStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=759) ¶

    type ForStmt struct {
     For  [token](/go/token).[Pos](/go/token#Pos) // position of "for" keyword
     Init Stmt      // initialization statement; or nil
     Cond Expr      // condition; or nil
     Post Stmt      // post iteration statement; or nil
     Body *BlockStmt
    }

A ForStmt represents a for statement.

#### func (*ForStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=862) ¶

    func (s *ForStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ForStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=800) ¶

    func (s *ForStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [FuncDecl](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=998) ¶

    type FuncDecl struct {
     Doc  *CommentGroup // associated documentation; or nil
     Recv *FieldList    // receiver (methods); or nil (functions)
     Name *Ident        // function/method name
     Type *FuncType     // function signature: type and value parameters, results, and position of "func" keyword
     Body *BlockStmt    // function body; or nil for external (non-Go) function
    }

A FuncDecl node represents a function declaration.

#### func (*FuncDecl) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1020) ¶

    func (d *FuncDecl) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*FuncDecl) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1011) ¶

    func (d *FuncDecl) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [FuncLit](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=325) ¶

    type FuncLit struct {
     Type *FuncType  // function type
     Body *BlockStmt // function body
    }

A FuncLit node represents a function literal.

#### func (*FuncLit) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=539) ¶

    func (x *FuncLit) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*FuncLit) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=500) ¶

    func (x *FuncLit) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [FuncType](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=464) ¶

    type FuncType struct {
     Func       [token](/go/token).[Pos](/go/token#Pos)  // position of "func" keyword (token.NoPos if there is no "func")
     TypeParams *FieldList // type parameters; or nil
     Params     *FieldList // (incoming) parameters; non-nil
     Results    *FieldList // (outgoing) results; or nil
    }

A FuncType node represents a function type.

#### func (*FuncType) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=554) ¶

    func (x *FuncType) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*FuncType) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=520) ¶

    func (x *FuncType) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [GenDecl](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=988) ¶

    type GenDecl struct {
     Doc    *CommentGroup // associated documentation; or nil
     TokPos [token](/go/token).[Pos](/go/token#Pos)     // position of Tok
     Tok    [token](/go/token).[Token](/go/token#Token)   // IMPORT, CONST, TYPE, or VAR
     Lparen [token](/go/token).[Pos](/go/token#Pos)     // position of '(', if any
     Specs  []Spec
     Rparen [token](/go/token).[Pos](/go/token#Pos) // position of ')', if any
    }

A GenDecl node (generic declaration node) represents an import, constant, type or variable declaration. A valid Lparen position (Lparen.IsValid()) indicates a parenthesized declaration.

Relationship between Tok value and Specs element type:

    token.IMPORT  *ImportSpec
    token.CONST   *ValueSpec
    token.TYPE    *TypeSpec
    token.VAR     *ValueSpec
    

#### func (*GenDecl) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1014) ¶

    func (d *GenDecl) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*GenDecl) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1010) ¶

    func (d *GenDecl) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [GoStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=678) ¶

    type GoStmt struct {
     Go   [token](/go/token).[Pos](/go/token#Pos) // position of "go" keyword
     Call *CallExpr
    }

A GoStmt node represents a go statement.

#### func (*GoStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=818) ¶

    func (s *GoStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*GoStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=789) ¶

    func (s *GoStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Ident](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=291) ¶

    type Ident struct {
     NamePos [token](/go/token).[Pos](/go/token#Pos) // identifier position
     Name    [string](/builtin#string)    // identifier name
     Obj     *Object   // denoted object, or nil. Deprecated: see Object.
    }

An Ident node represents an identifier.

#### func [NewIdent](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=596) ¶

    func NewIdent(name [string](/builtin#string)) *Ident

NewIdent creates a new Ident without position. Useful for ASTs generated by code other than the Go parser.

#### func (*Ident) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=531) ¶

    func (x *Ident) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*Ident) [IsExported](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=602) ¶

    func (id *Ident) IsExported() [bool](/builtin#bool)

IsExported reports whether id starts with an upper-case letter.

#### func (*Ident) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=497) ¶

    func (x *Ident) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### func (*Ident) [String](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=604) ¶

    func (id *Ident) String() [string](/builtin#string)

#### type [IfStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=712) ¶

    type IfStmt struct {
     If   [token](/go/token).[Pos](/go/token#Pos) // position of "if" keyword
     Init Stmt      // initialization statement; or nil
     Cond Expr      // condition
     Body *BlockStmt
     Else Stmt // else branch; or nil
    }

An IfStmt node represents an if statement.

#### func (*IfStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=841) ¶

    func (s *IfStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*IfStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=794) ¶

    func (s *IfStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [ImportSpec](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=902) ¶

    type ImportSpec struct {
     Doc     *CommentGroup // associated documentation; or nil
     Name    *Ident        // local package name (including "."); or nil
     Path    *BasicLit     // import path
     Comment *CommentGroup // line comments; or nil
     EndPos  [token](/go/token).[Pos](/go/token#Pos)     // end of spec (overrides Path.Pos if nonzero)
    }

An ImportSpec node represents a single package import.

#### func (*ImportSpec) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=943) ¶

    func (s *ImportSpec) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ImportSpec) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=934) ¶

    func (s *ImportSpec) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Importer](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/resolve.go;l=65) deprecated

    type Importer func(imports map[[string](/builtin#string)]*Object, path [string](/builtin#string)) (pkg *Object, err [error](/builtin#error))

An Importer resolves import paths to package Objects. The imports map records the packages already imported, indexed by package id (canonical import path). An Importer must determine the canonical import path and check the map to see if it is already present in the imports map. If so, the Importer can return the map entry. Otherwise, the Importer should load the package data for the given path into a new *Object (pkg), record pkg in the imports map, and then return pkg.

Deprecated: use the type checker [go/types](/go/types) instead; see Object.

#### type [IncDecStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=661) ¶

    type IncDecStmt struct {
     X      Expr
     TokPos [token](/go/token).[Pos](/go/token#Pos)   // position of Tok
     Tok    [token](/go/token).[Token](/go/token#Token) // INC or DEC
    }

An IncDecStmt node represents an increment or decrement statement.

#### func (*IncDecStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=814) ¶

    func (s *IncDecStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*IncDecStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=787) ¶

    func (s *IncDecStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [IndexExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=353) ¶

    type IndexExpr struct {
     X      Expr      // expression
     Lbrack [token](/go/token).[Pos](/go/token#Pos) // position of "["
     Index  Expr      // index expression
     Rbrack [token](/go/token).[Pos](/go/token#Pos) // position of "]"
    }

An IndexExpr node represents an expression followed by an index.

#### func (*IndexExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=543) ¶

    func (x *IndexExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*IndexExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=509) ¶

    func (x *IndexExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [IndexListExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=362) ¶ added in go1.18

    type IndexListExpr struct {
     X       Expr      // expression
     Lbrack  [token](/go/token).[Pos](/go/token#Pos) // position of "["
     Indices []Expr    // index expressions
     Rbrack  [token](/go/token).[Pos](/go/token#Pos) // position of "]"
    }

An IndexListExpr node represents an expression followed by multiple indices.

#### func (*IndexListExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=544) ¶ added in go1.18

    func (x *IndexListExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*IndexListExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=510) ¶ added in go1.18

    func (x *IndexListExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [InterfaceType](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=472) ¶

    type InterfaceType struct {
     Interface  [token](/go/token).[Pos](/go/token#Pos)  // position of "interface" keyword
     Methods    *FieldList // list of embedded interfaces, methods, or types
     Incomplete [bool](/builtin#bool)       // true if (source) methods or types are missing in the Methods list
    }

An InterfaceType node represents an interface type.

#### func (*InterfaceType) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=560) ¶

    func (x *InterfaceType) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*InterfaceType) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=526) ¶

    func (x *InterfaceType) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [KeyValueExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=427) ¶

    type KeyValueExpr struct {
     Key   Expr
     Colon [token](/go/token).[Pos](/go/token#Pos) // position of ":"
     Value Expr
    }

A KeyValueExpr node represents (key : value) pairs in composite literals.

#### func (*KeyValueExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=551) ¶

    func (x *KeyValueExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*KeyValueExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=517) ¶

    func (x *KeyValueExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [LabeledStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=640) ¶

    type LabeledStmt struct {
     Label *Ident
     Colon [token](/go/token).[Pos](/go/token#Pos) // position of ":"
     Stmt  Stmt
    }

A LabeledStmt node represents a labeled statement.

#### func (*LabeledStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=811) ¶

    func (s *LabeledStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*LabeledStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=784) ¶

    func (s *LabeledStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [MapType](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=479) ¶

    type MapType struct {
     Map   [token](/go/token).[Pos](/go/token#Pos) // position of "map" keyword
     Key   Expr
     Value Expr
    }

A MapType node represents a map type.

#### func (*MapType) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=561) ¶

    func (x *MapType) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*MapType) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=527) ¶

    func (x *MapType) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [MergeMode](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/filter.go;l=306) deprecated

    type MergeMode [uint](/builtin#uint)

The MergeMode flags control the behavior of MergePackageFiles.

Deprecated: use the type checker [go/types](/go/types) instead of Package; see Object.

    const (
     // If set, duplicate function declarations are excluded.
     FilterFuncDuplicates MergeMode = 1 << [iota](/builtin#iota)
     // If set, comments that are not associated with a specific
     // AST node (as Doc or Comment) are excluded.
     FilterUnassociatedComments
     // If set, duplicate import declarations are excluded.
     FilterImportDuplicates
    )

Deprecated: use the type checker [go/types](/go/types) instead of Package; see Object.

#### type [Node](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=36) ¶

    type Node interface {
     Pos() [token](/go/token).[Pos](/go/token#Pos) // position of first character belonging to the node
     End() [token](/go/token).[Pos](/go/token#Pos) // position of first character immediately after the node
    }

All node types implement the Node interface.

#### type [ObjKind](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=148) ¶

    type ObjKind [int](/builtin#int)

ObjKind describes what an Object represents.

    const (
     Bad ObjKind = [iota](/builtin#iota) // for error handling
     Pkg                // package
     Con                // constant
     Typ                // type
     Var                // variable
     Fun                // function or method
     Lbl                // label
    )

The list of possible Object kinds.

#### func (ObjKind) [String](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=171) ¶

    func (kind ObjKind) String() [string](/builtin#string)

#### type [Object](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=87) deprecated

    type Object struct {
     Kind ObjKind
     Name [string](/builtin#string) // declared name
     Decl [any](/builtin#any)    // corresponding Field, XxxSpec, FuncDecl, LabeledStmt, AssignStmt, Scope; or nil
     Data [any](/builtin#any)    // object-specific data; or nil
     Type [any](/builtin#any)    // placeholder for type information; may be nil
    }

An Object describes a named language entity such as a package, constant, type, variable, function (incl. methods), or label.

The Data fields contains object-specific data:

    Kind    Data type         Data value
    Pkg     *Scope            package scope
    Con     int               iota for the respective declaration
    

Deprecated: The relationship between Idents and Objects cannot be correctly computed without type information. For example, the expression T{K: 0} may denote a struct, map, slice, or array literal, depending on the type of T. If T is a struct, then K refers to a field of T, whereas for the other types it refers to a value in the environment.

New programs should set the [parser.SkipObjectResolution] parser flag to disable syntactic object resolution (which also saves CPU and memory), and instead use the type checker [go/types](/go/types) if object resolution is desired. See the Defs, Uses, and Implicits fields of the [types.Info] struct for details.

#### func [NewObj](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=96) ¶

    func NewObj(kind ObjKind, name [string](/builtin#string)) *Object

NewObj creates a new object of a given kind and name.

#### func (*Object) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=103) ¶

    func (obj *Object) Pos() [token](/go/token).[Pos](/go/token#Pos)

Pos computes the source position of the declaration of an object name. The result may be an invalid position if it cannot be computed (obj.Decl may be nil or not correct).

#### type [Package](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1092) deprecated

    type Package struct {
     Name    [string](/builtin#string)             // package name
     Scope   *Scope             // package scope across all files
     Imports map[[string](/builtin#string)]*Object // map of package id -> package object
     Files   map[[string](/builtin#string)]*File   // Go source files by filename
    }

A Package node represents a set of source files collectively building a Go package.

Deprecated: use the type checker [go/types](/go/types) instead; see Object.

#### func [NewPackage](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/resolve.go;l=77) deprecated

    func NewPackage(fset *[token](/go/token).[FileSet](/go/token#FileSet), files map[[string](/builtin#string)]*File, importer Importer, universe *Scope) (*Package, [error](/builtin#error))

NewPackage creates a new Package node from a set of File nodes. It resolves unresolved identifiers across files and updates each file's Unresolved list accordingly. If a non-nil importer and universe scope are provided, they are used to resolve identifiers not declared in any of the package files. Any remaining unresolved identifiers are reported as undeclared. If the files belong to different packages, one package name is selected and files with different package names are reported and then ignored. The result is a package node and a [scanner.ErrorList](/go/scanner#ErrorList) if there were errors.

Deprecated: use the type checker [go/types](/go/types) instead; see Object.

#### func (*Package) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1100) ¶

    func (p *Package) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*Package) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=1099) ¶

    func (p *Package) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [ParenExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=340) ¶

    type ParenExpr struct {
     Lparen [token](/go/token).[Pos](/go/token#Pos) // position of "("
     X      Expr      // parenthesized expression
     Rparen [token](/go/token).[Pos](/go/token#Pos) // position of ")"
    }

A ParenExpr node represents a parenthesized expression.

#### func (*ParenExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=541) ¶

    func (x *ParenExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ParenExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=507) ¶

    func (x *ParenExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [RangeStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=768) ¶

    type RangeStmt struct {
     For        [token](/go/token).[Pos](/go/token#Pos)   // position of "for" keyword
     Key, Value Expr        // Key, Value may be nil
     TokPos     [token](/go/token).[Pos](/go/token#Pos)   // position of Tok; invalid if Key == nil
     Tok        [token](/go/token).[Token](/go/token#Token) // ILLEGAL if Key == nil, ASSIGN, DEFINE
     Range      [token](/go/token).[Pos](/go/token#Pos)   // position of "range" keyword
     X          Expr        // value to range over
     Body       *BlockStmt
    }

A RangeStmt represents a for statement with a range clause.

#### func (*RangeStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=863) ¶

    func (s *RangeStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*RangeStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=801) ¶

    func (s *RangeStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [ReturnStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=690) ¶

    type ReturnStmt struct {
     Return  [token](/go/token).[Pos](/go/token#Pos) // position of "return" keyword
     Results []Expr    // result expressions; or nil
    }

A ReturnStmt node represents a return statement.

#### func (*ReturnStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=820) ¶

    func (s *ReturnStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ReturnStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=791) ¶

    func (s *ReturnStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Scope](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=20) deprecated

    type Scope struct {
     Outer   *Scope
     Objects map[[string](/builtin#string)]*Object
    }

A Scope maintains the set of named language entities declared in the scope and a link to the immediately surrounding (outer) scope.

Deprecated: use the type checker [go/types](/go/types) instead; see Object.

#### func [NewScope](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=26) ¶

    func NewScope(outer *Scope) *Scope

NewScope creates a new scope nested in the outer scope.

#### func (*Scope) [Insert](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=42) ¶

    func (s *Scope) Insert(obj *Object) (alt *Object)

Insert attempts to insert a named object obj into the scope s. If the scope already contains an object alt with the same name, Insert leaves the scope unchanged and returns alt. Otherwise it inserts obj and returns nil.

#### func (*Scope) [Lookup](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=34) ¶

    func (s *Scope) Lookup(name [string](/builtin#string)) *Object

Lookup returns the object with the given name if it is found in scope s, otherwise it returns nil. Outer scopes are ignored.

#### func (*Scope) [String](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/scope.go;l=50) ¶

    func (s *Scope) String() [string](/builtin#string)

Debugging support

#### type [SelectStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=753) ¶

    type SelectStmt struct {
     Select [token](/go/token).[Pos](/go/token#Pos)  // position of "select" keyword
     Body   *BlockStmt // CommClauses only
    }

A SelectStmt node represents a select statement.

#### func (*SelectStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=861) ¶

    func (s *SelectStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*SelectStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=799) ¶

    func (s *SelectStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [SelectorExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=347) ¶

    type SelectorExpr struct {
     X   Expr   // expression
     Sel *Ident // field selector
    }

A SelectorExpr node represents an expression followed by a selector.

#### func (*SelectorExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=542) ¶

    func (x *SelectorExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*SelectorExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=508) ¶

    func (x *SelectorExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [SendStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=654) ¶

    type SendStmt struct {
     Chan  Expr
     Arrow [token](/go/token).[Pos](/go/token#Pos) // position of "<-"
     Value Expr
    }

A SendStmt node represents a send statement.

#### func (*SendStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=813) ¶

    func (s *SendStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*SendStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=786) ¶

    func (s *SendStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [SliceExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=370) ¶

    type SliceExpr struct {
     X      Expr      // expression
     Lbrack [token](/go/token).[Pos](/go/token#Pos) // position of "["
     Low    Expr      // begin of slice range; or nil
     High   Expr      // end of slice range; or nil
     Max    Expr      // maximum capacity of slice; or nil
     Slice3 [bool](/builtin#bool)      // true if 3-index slice (2 colons present)
     Rbrack [token](/go/token).[Pos](/go/token#Pos) // position of "]"
    }

A SliceExpr node represents an expression followed by slice indices.

#### func (*SliceExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=545) ¶

    func (x *SliceExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*SliceExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=511) ¶

    func (x *SliceExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Spec](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=896) ¶

    type Spec interface {
     Node
     // contains filtered or unexported methods
    }

The Spec type stands for any of *ImportSpec,*ValueSpec, and *TypeSpec.

#### type [StarExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=402) ¶

    type StarExpr struct {
     Star [token](/go/token).[Pos](/go/token#Pos) // position of "*"
     X    Expr      // operand
    }

A StarExpr node represents an expression of the form "*" Expression. Semantically it could be a unary "*" expression, or a pointer type.

#### func (*StarExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=548) ¶

    func (x *StarExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*StarExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=514) ¶

    func (x *StarExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Stmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=48) ¶

    type Stmt interface {
     Node
     // contains filtered or unexported methods
    }

All statement nodes implement the Stmt interface.

#### type [StructType](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=455) ¶

    type StructType struct {
     Struct     [token](/go/token).[Pos](/go/token#Pos)  // position of "struct" keyword
     Fields     *FieldList // list of field declarations
     Incomplete [bool](/builtin#bool)       // true if (source) fields are missing in the Fields list
    }

A StructType node represents a struct type.

#### func (*StructType) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=553) ¶

    func (x *StructType) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*StructType) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=519) ¶

    func (x *StructType) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [SwitchStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=729) ¶

    type SwitchStmt struct {
     Switch [token](/go/token).[Pos](/go/token#Pos)  // position of "switch" keyword
     Init   Stmt       // initialization statement; or nil
     Tag    Expr       // tag expression; or nil
     Body   *BlockStmt // CaseClauses only
    }

A SwitchStmt node represents an expression switch statement.

#### func (*SwitchStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=853) ¶

    func (s *SwitchStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*SwitchStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=796) ¶

    func (s *SwitchStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [TypeAssertExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=383) ¶

    type TypeAssertExpr struct {
     X      Expr      // expression
     Lparen [token](/go/token).[Pos](/go/token#Pos) // position of "("
     Type   Expr      // asserted type; nil means type switch X.(type)
     Rparen [token](/go/token).[Pos](/go/token#Pos) // position of ")"
    }

A TypeAssertExpr node represents an expression followed by a type assertion.

#### func (*TypeAssertExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=546) ¶

    func (x *TypeAssertExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*TypeAssertExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=512) ¶

    func (x *TypeAssertExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [TypeSpec](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=922) ¶

    type TypeSpec struct {
     Doc        *CommentGroup // associated documentation; or nil
     Name       *Ident        // type name
     TypeParams *FieldList    // type parameters; or nil
     Assign     [token](/go/token).[Pos](/go/token#Pos)     // position of '=', if any
     Type       Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
     Comment    *CommentGroup // line comments; or nil
    }

A TypeSpec node represents a type declaration (TypeSpec production).

#### func (*TypeSpec) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=959) ¶

    func (s *TypeSpec) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*TypeSpec) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=941) ¶

    func (s *TypeSpec) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [TypeSwitchStmt](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=737) ¶

    type TypeSwitchStmt struct {
     Switch [token](/go/token).[Pos](/go/token#Pos)  // position of "switch" keyword
     Init   Stmt       // initialization statement; or nil
     Assign Stmt       // x := y.(type) or y.(type)
     Body   *BlockStmt // CaseClauses only
    }

A TypeSwitchStmt node represents a type switch statement.

#### func (*TypeSwitchStmt) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=854) ¶

    func (s *TypeSwitchStmt) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*TypeSwitchStmt) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=797) ¶

    func (s *TypeSwitchStmt) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [UnaryExpr](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=410) ¶

    type UnaryExpr struct {
     OpPos [token](/go/token).[Pos](/go/token#Pos)   // position of Op
     Op    [token](/go/token).[Token](/go/token#Token) // operator
     X     Expr        // operand
    }

A UnaryExpr node represents a unary expression. Unary "*" expressions are represented via StarExpr nodes.

#### func (*UnaryExpr) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=549) ¶

    func (x *UnaryExpr) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*UnaryExpr) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=515) ¶

    func (x *UnaryExpr) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [ValueSpec](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=913) ¶

    type ValueSpec struct {
     Doc     *CommentGroup // associated documentation; or nil
     Names   []*Ident      // value names (len(Names) > 0)
     Type    Expr          // value type; or nil
     Values  []Expr        // initial values; or nil
     Comment *CommentGroup // line comments; or nil
    }

A ValueSpec node represents a constant or variable declaration (ConstSpec or VarSpec production).

#### func (*ValueSpec) [End](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=950) ¶

    func (s *ValueSpec) End() [token](/go/token).[Pos](/go/token#Pos)

#### func (*ValueSpec) [Pos](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/ast.go;l=940) ¶

    func (s *ValueSpec) Pos() [token](/go/token).[Pos](/go/token#Pos)

#### type [Visitor](https://cs.opensource.google/go/go/+/go1.25.6:src/go/ast/walk.go;l=15) ¶

    type Visitor interface {
     Visit(node Node) (w Visitor)
    }

A Visitor's Visit method is invoked for each node encountered by Walk. If the result visitor w is not nil, Walk visits each of the children of node with the visitor w, followed by a call of w.Visit(nil).
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
