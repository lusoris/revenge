# golang.org/x/tools/go/packages

> Source: https://pkg.go.dev/golang.org/x/tools/go/packages
> Fetched: 2026-01-31T16:07:25.637187+00:00
> Content-Hash: d00044935d10fdff
> Type: html

---

### Overview ¶

  * The driver protocol



Package packages loads Go packages for inspection and analysis. 

The Load function takes as input a list of patterns and returns a list of Package values describing individual packages matched by those patterns. A Config specifies configuration options, the most important of which is the LoadMode, which controls the amount of detail in the loaded packages. 

Load passes most patterns directly to the underlying build tool. The default build tool is the go command. Its supported patterns are described at <https://pkg.go.dev/cmd/go#hdr-Package_lists_and_patterns>. Other build systems may be supported by providing a "driver"; see [The driver protocol]. 

All patterns with the prefix "query=", where query is a non-empty string of letters from [a-z], are reserved and may be interpreted as query operators. 

Two query operators are currently supported: "file" and "pattern". 

The query "file=path/to/file.go" matches the package or packages enclosing the Go source file path/to/file.go. For example "file=~/go/src/fmt/print.go" might return the packages "fmt" and "fmt [fmt.test]". 

The query "pattern=string" causes "string" to be passed directly to the underlying build tool. In most cases this is unnecessary, but an application can use Load("pattern=" + x) as an escaping mechanism to ensure that x is not interpreted as a query operator if it contains '='. 

All other query operators are reserved for future use and currently cause Load to report an error. 

The Package struct provides basic information about the package, including 

  * ID, a unique identifier for the package in the returned set;
  * GoFiles, the names of the package's Go source files;
  * Imports, a map from source import strings to the Packages they name;
  * Types, the type information for the package's exported symbols;
  * Syntax, the parsed syntax trees for the package's source code; and
  * TypesInfo, the result of a complete type-check of the package syntax trees.



(See the documentation for type Package for the complete list of fields and more detailed descriptions.) 

For example, 
    
    
    Load(nil, "bytes", "unicode...")
    

returns four Package structs describing the standard library packages bytes, unicode, unicode/utf16, and unicode/utf8. Note that one pattern can match multiple packages and that a package might be matched by multiple patterns: in general it is not possible to determine which packages correspond to which patterns. 

Note that the list returned by Load contains only the packages matched by the patterns. Their dependencies can be found by walking the import graph using the Imports fields. 

The Load function can be configured by passing a pointer to a Config as the first argument. A nil Config is equivalent to the zero Config, which causes Load to run in LoadFiles mode, collecting minimal information. See the documentation for type Config for details. 

As noted earlier, the Config.Mode controls the amount of detail reported about the loaded packages. See the documentation for type LoadMode for details. 

Most tools should pass their command-line arguments (after any flags) uninterpreted to Load, so that it can interpret them according to the conventions of the underlying build system. 

See the Example function for typical usage. See also [golang.org/x/tools/go/packages/internal/linecount](/golang.org/x/tools@v0.41.0/go/packages/internal/linecount) for an example application. 

#### The driver protocol ¶

Load may be used to load Go packages even in Go projects that use alternative build systems, by installing an appropriate "driver" program for the build system and specifying its location in the GOPACKAGESDRIVER environment variable. For example, <https://github.com/bazelbuild/rules_go/wiki/Editor-and-tool-integration> explains how to use the driver for Bazel. 

The driver program is responsible for interpreting patterns in its preferred notation and reporting information about the packages that those patterns identify. Drivers must also support the special "file=" and "pattern=" patterns described above. 

The patterns are provided as positional command-line arguments. A JSON-encoded DriverRequest message providing additional information is written to the driver's standard input. The driver must write a JSON-encoded DriverResponse message to its standard output. (This message differs from the JSON schema produced by 'go list'.) 

The value of the PWD environment variable seen by the driver process is the preferred name of its working directory. (The working directory may have other aliases due to symbolic links; see the comment on the Dir field of [exec.Cmd](/os/exec#Cmd) for related information.) When the driver process emits in its response the name of a file that is a descendant of this directory, it must use an absolute path that has the value of PWD as a prefix, to ensure that the returned filenames satisfy the original query. 

Example ¶

Example demonstrates how to load the packages specified on the command line from source syntax. 
    
    
    package main
    
    import (
    	"flag"
    	"fmt"
    	"os"
    
    	"golang.org/x/tools/go/packages"
    )
    
    func main() {
    	flag.Parse()
    
    	// Many tools pass their command-line arguments (after any flags)
    	// uninterpreted to packages.Load so that it can interpret them
    	// according to the conventions of the underlying build system.
    	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
    	pkgs, err := packages.Load(cfg, flag.Args()...)
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "load: %v\n", err)
    		os.Exit(1)
    	}
    	if packages.PrintErrors(pkgs) > 0 {
    		os.Exit(1)
    	}
    
    	// Print the names of the source files
    	// for each package listed on the command line.
    	for _, pkg := range pkgs {
    		fmt.Println(pkg.ID, pkg.GoFiles)
    	}
    }
    

Share Format Run

### Index ¶

  * Constants
  * func Postorder(pkgs []*Package) iter.Seq[*Package]
  * func PrintErrors(pkgs []*Package) int
  * func Visit(pkgs []*Package, pre func(*Package) bool, post func(*Package))
  * type Config
  * type DriverRequest
  * type DriverResponse
  * type Error
  *     * func (err Error) Error() string
  * type ErrorKind
  * type LoadMode
  *     * func (mode LoadMode) String() string
  * type Module
  * type ModuleError
  * type Package
  *     * func Load(cfg *Config, patterns ...string) ([]*Package, error)
  *     * func (p *Package) MarshalJSON() ([]byte, error)
    * func (p *Package) String() string
    * func (p *Package) UnmarshalJSON(b []byte) error



### Examples ¶

  * Package



### Constants ¶

[View Source](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=127)
    
    
    const (
    	// LoadFiles loads the name and file names for the initial packages.
    	LoadFiles = NeedName | NeedFiles | NeedCompiledGoFiles
    
    	// LoadImports loads the name, file names, and import mapping for the initial packages.
    	LoadImports = LoadFiles | NeedImports
    
    	// LoadTypes loads exported type information for the initial packages.
    	LoadTypes = LoadImports | NeedTypes | NeedTypesSizes
    
    	// LoadSyntax loads typed syntax for the initial packages.
    	LoadSyntax = LoadTypes | NeedSyntax | NeedTypesInfo
    
    	// LoadAllSyntax loads typed syntax for the initial packages and all dependencies.
    	LoadAllSyntax = LoadSyntax | NeedDeps
    
    	// Deprecated: NeedExportsFile is a historical misspelling of NeedExportFile.
    	//
    	//go:fix inline
    	NeedExportsFile = NeedExportFile
    )

### Variables ¶

This section is empty.

### Functions ¶

####  func [Postorder](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/visit.go;l=84) ¶ added in v0.37.0
    
    
    func Postorder(pkgs []*Package) [iter](/iter).[Seq](/iter#Seq)[*Package]

Postorder returns an iterator over the packages in the import graph whose roots are pkg. Packages are enumerated in dependencies-first order. 

####  func [PrintErrors](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/visit.go;l=61) ¶
    
    
    func PrintErrors(pkgs []*Package) [int](/builtin#int)

PrintErrors prints to os.Stderr the accumulated errors of all packages in the import graph rooted at pkgs, dependencies first. PrintErrors returns the number of errors printed. 

####  func [Visit](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/visit.go;l=35) ¶
    
    
    func Visit(pkgs []*Package, pre func(*Package) [bool](/builtin#bool), post func(*Package))

Visit visits all the packages in the import graph whose roots are pkgs, calling the optional pre function the first time each package is encountered (preorder), and the optional post function after a package's dependencies have been visited (postorder). The boolean result of pre(pkg) determines whether the imports of package pkg are visited. 

Example: 
    
    
    pkgs, err := Load(...)
    if err != nil { ... }
    Visit(pkgs, nil, func(pkg *Package) {
    	log.Println(pkg)
    })
    

In most cases, it is more convenient to use Postorder: 
    
    
    for pkg := range Postorder(pkgs) {
    	log.Println(pkg)
    }
    

### Types ¶

####  type [Config](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=153) ¶
    
    
    type Config struct {
    	// Mode controls the level of information returned for each package.
    	Mode LoadMode
    
    	// Context specifies the context for the load operation.
    	// Cancelling the context may cause [Load] to abort and
    	// return an error.
    	Context [context](/context).[Context](/context#Context)
    
    	// Logf is the logger for the config.
    	// If the user provides a logger, debug logging is enabled.
    	// If the GOPACKAGESDEBUG environment variable is set to true,
    	// but the logger is nil, default to log.Printf.
    	Logf func(format [string](/builtin#string), args ...[any](/builtin#any))
    
    	// Dir is the directory in which to run the build system's query tool
    	// that provides information about the packages.
    	// If Dir is empty, the tool is run in the current directory.
    	Dir [string](/builtin#string)
    
    	// Env is the environment to use when invoking the build system's query tool.
    	// If Env is nil, the current environment is used.
    	// As in os/exec's Cmd, only the last value in the slice for
    	// each environment key is used. To specify the setting of only
    	// a few variables, append to the current environment, as in:
    	//
    	//	opt.Env = append(os.Environ(), "GOOS=plan9", "GOARCH=386")
    	//
    	Env [][string](/builtin#string)
    
    	// BuildFlags is a list of command-line flags to be passed through to
    	// the build system's query tool.
    	BuildFlags [][string](/builtin#string)
    
    	// Fset provides source position information for syntax trees and types.
    	// If Fset is nil, Load will use a new fileset, but preserve Fset's value.
    	Fset *[token](/go/token).[FileSet](/go/token#FileSet)
    
    	// ParseFile is called to read and parse each file
    	// when preparing a package's type-checked syntax tree.
    	// It must be safe to call ParseFile simultaneously from multiple goroutines.
    	// If ParseFile is nil, the loader will uses parser.ParseFile.
    	//
    	// ParseFile should parse the source from src and use filename only for
    	// recording position information.
    	//
    	// An application may supply a custom implementation of ParseFile
    	// to change the effective file contents or the behavior of the parser,
    	// or to modify the syntax tree. For example, selectively eliminating
    	// unwanted function bodies can significantly accelerate type checking.
    	ParseFile func(fset *[token](/go/token).[FileSet](/go/token#FileSet), filename [string](/builtin#string), src [][byte](/builtin#byte)) (*[ast](/go/ast).[File](/go/ast#File), [error](/builtin#error))
    
    	// If Tests is set, the loader includes not just the packages
    	// matching a particular pattern but also any related test packages,
    	// including test-only variants of the package and the test executable.
    	//
    	// For example, when using the go command, loading "fmt" with Tests=true
    	// returns four packages, with IDs "fmt" (the standard package),
    	// "fmt [fmt.test]" (the package as compiled for the test),
    	// "fmt_test" (the test functions from source files in package fmt_test),
    	// and "fmt.test" (the test binary).
    	//
    	// In build systems with explicit names for tests,
    	// setting Tests may have no effect.
    	Tests [bool](/builtin#bool)
    
    	// Overlay is a mapping from absolute file paths to file contents.
    	//
    	// For each map entry, [Load] uses the alternative file
    	// contents provided by the overlay mapping instead of reading
    	// from the file system. This mechanism can be used to enable
    	// editor-integrated tools to correctly analyze the contents
    	// of modified but unsaved buffers, for example.
    	//
    	// The overlay mapping is passed to the build system's driver
    	// (see "The driver protocol") so that it too can report
    	// consistent package metadata about unsaved files. However,
    	// drivers may vary in their level of support for overlays.
    	Overlay map[[string](/builtin#string)][][byte](/builtin#byte)
    }

A Config specifies details about how packages should be loaded. The zero value is a valid configuration. 

Calls to Load do not modify this struct. 

####  type [DriverRequest](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/external.go;l=26) ¶ added in v0.18.0
    
    
    type DriverRequest struct {
    	Mode LoadMode `json:"mode"`
    
    	// Env specifies the environment the underlying build system should be run in.
    	Env [][string](/builtin#string) `json:"env"`
    
    	// BuildFlags are flags that should be passed to the underlying build system.
    	BuildFlags [][string](/builtin#string) `json:"build_flags"`
    
    	// Tests specifies whether the patterns should also return test packages.
    	Tests [bool](/builtin#bool) `json:"tests"`
    
    	// Overlay maps file paths (relative to the driver's working directory)
    	// to the contents of overlay files (see Config.Overlay).
    	Overlay map[[string](/builtin#string)][][byte](/builtin#byte) `json:"overlay"`
    }

DriverRequest defines the schema of a request for package metadata from an external driver program. The JSON-encoded DriverRequest message is provided to the driver program's standard input. The query patterns are provided as command-line arguments. 

See the package documentation for an overview. 

####  type [DriverResponse](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/external.go;l=49) ¶ added in v0.18.0
    
    
    type DriverResponse struct {
    	// NotHandled is returned if the request can't be handled by the current
    	// driver. If an external driver returns a response with NotHandled, the
    	// rest of the DriverResponse is ignored, and go/packages will fallback
    	// to the next driver. If go/packages is extended in the future to support
    	// lists of multiple drivers, go/packages will fall back to the next driver.
    	NotHandled [bool](/builtin#bool)
    
    	// Compiler and Arch are the arguments pass of types.SizesFor
    	// to get a types.Sizes to use when type checking.
    	Compiler [string](/builtin#string)
    	Arch     [string](/builtin#string)
    
    	// Roots is the set of package IDs that make up the root packages.
    	// We have to encode this separately because when we encode a single package
    	// we cannot know if it is one of the roots as that requires knowledge of the
    	// graph it is part of.
    	Roots [][string](/builtin#string) `json:",omitempty"`
    
    	// Packages is the full set of packages in the graph.
    	// The packages are not connected into a graph.
    	// The Imports if populated will be stubs that only have their ID set.
    	// Imports will be connected and then type and syntax information added in a
    	// later pass (see refine).
    	Packages []*Package
    
    	// GoVersion is the minor version number used by the driver
    	// (e.g. the go command on the PATH) when selecting .go files.
    	// Zero means unknown.
    	GoVersion [int](/builtin#int)
    }

DriverResponse defines the schema of a response from an external driver program, providing the results of a query for package metadata. The driver program must write a JSON-encoded DriverResponse message to its standard output. 

See the package documentation for an overview. 

####  type [Error](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=569) ¶
    
    
    type Error struct {
    	Pos  [string](/builtin#string) // "file:line:col" or "file:line" or "" or "-"
    	Msg  [string](/builtin#string)
    	Kind ErrorKind
    }

An Error describes a problem with a package's metadata, syntax, or types. 

####  func (Error) [Error](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=587) ¶
    
    
    func (err Error) Error() [string](/builtin#string)

####  type [ErrorKind](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=578) ¶
    
    
    type ErrorKind [int](/builtin#int)

ErrorKind describes the source of the error, allowing the user to differentiate between errors generated by the driver, the parser, or the type-checker. 
    
    
    const (
    	UnknownError ErrorKind = [iota](/builtin#iota)
    	ListError
    	ParseError
    	TypeError
    )

####  type [LoadMode](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=66) ¶
    
    
    type LoadMode [int](/builtin#int)

A LoadMode controls the amount of detail to return when loading. The bits below can be combined to specify which fields should be filled in the result packages. 

The zero value is a special case, equivalent to combining the NeedName, NeedFiles, and NeedCompiledGoFiles bits. 

ID and Errors (if present) will always be filled. Load may return more information than requested. 

The Mode flag is a union of several bits named NeedName, NeedFiles, and so on, each of which determines whether a given field of Package (Name, Files, etc) should be populated. 

For convenience, we provide named constants for the most common combinations of Need flags: 
    
    
    [LoadFiles]     lists of files in each package
    [LoadImports]   ... plus imports
    [LoadTypes]     ... plus type information
    [LoadSyntax]    ... plus type-annotated syntax
    [LoadAllSyntax] ... for all dependencies
    

Unfortunately there are a number of open bugs related to interactions among the LoadMode bits: 

  * <https://go.dev/issue/56633>
  * <https://go.dev/issue/56677>
  * <https://go.dev/issue/58726>
  * <https://go.dev/issue/63517>


    
    
    const (
    	// NeedName adds Name and PkgPath.
    	NeedName LoadMode = 1 << [iota](/builtin#iota)
    
    	// NeedFiles adds Dir, GoFiles, OtherFiles, and IgnoredFiles
    	NeedFiles
    
    	// NeedCompiledGoFiles adds CompiledGoFiles.
    	NeedCompiledGoFiles
    
    	// NeedImports adds Imports. If NeedDeps is not set, the Imports field will contain
    	// "placeholder" Packages with only the ID set.
    	NeedImports
    
    	// NeedDeps adds the fields requested by the LoadMode in the packages in Imports.
    	NeedDeps
    
    	// NeedExportFile adds ExportFile.
    	NeedExportFile
    
    	// NeedTypes adds Types, Fset, and IllTyped.
    	NeedTypes
    
    	// NeedSyntax adds Syntax and Fset.
    	NeedSyntax
    
    	// NeedTypesInfo adds TypesInfo and Fset.
    	NeedTypesInfo
    
    	// NeedTypesSizes adds TypesSizes.
    	NeedTypesSizes
    
    	// NeedForTest adds ForTest.
    	//
    	// Tests must also be set on the context for this field to be populated.
    	NeedForTest
    
    	// NeedModule adds Module.
    	NeedModule
    
    	// NeedEmbedFiles adds EmbedFiles.
    	NeedEmbedFiles
    
    	// NeedEmbedPatterns adds EmbedPatterns.
    	NeedEmbedPatterns
    
    	// NeedTarget adds Target.
    	NeedTarget
    )

####  func (LoadMode) [String](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/loadmode_string.go;l=33) ¶
    
    
    func (mode LoadMode) String() [string](/builtin#string)

####  type [Module](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=542) ¶
    
    
    type Module struct {
    	Path      [string](/builtin#string)       // module path
    	Version   [string](/builtin#string)       // module version
    	Replace   *Module      // replaced by this module
    	Time      *[time](/time).[Time](/time#Time)   // time version was created
    	Main      [bool](/builtin#bool)         // is this the main module?
    	Indirect  [bool](/builtin#bool)         // is this module only an indirect dependency of main module?
    	Dir       [string](/builtin#string)       // directory holding files for this module, if any
    	GoMod     [string](/builtin#string)       // path to go.mod file used when loading this module, if any
    	GoVersion [string](/builtin#string)       // go version used in module
    	Error     *ModuleError // error loading module
    }

Module provides module information for a package. 

It also defines part of the JSON schema of DriverResponse. See the package documentation for an overview. 

####  type [ModuleError](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=556) ¶
    
    
    type ModuleError struct {
    	Err [string](/builtin#string) // the error itself
    }

ModuleError holds errors loading a module. 

####  type [Package](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=419) ¶
    
    
    type Package struct {
    	// ID is a unique identifier for a package,
    	// in a syntax provided by the underlying build system.
    	//
    	// Because the syntax varies based on the build system,
    	// clients should treat IDs as opaque and not attempt to
    	// interpret them.
    	ID [string](/builtin#string)
    
    	// Name is the package name as it appears in the package source code.
    	Name [string](/builtin#string)
    
    	// PkgPath is the package path as used by the go/types package.
    	PkgPath [string](/builtin#string)
    
    	// Dir is the directory associated with the package, if it exists.
    	//
    	// For packages listed by the go command, this is the directory containing
    	// the package files.
    	Dir [string](/builtin#string)
    
    	// Errors contains any errors encountered querying the metadata
    	// of the package, or while parsing or type-checking its files.
    	Errors []Error
    
    	// TypeErrors contains the subset of errors produced during type checking.
    	TypeErrors [][types](/go/types).[Error](/go/types#Error)
    
    	// GoFiles lists the absolute file paths of the package's Go source files.
    	// It may include files that should not be compiled, for example because
    	// they contain non-matching build tags, are documentary pseudo-files such as
    	// unsafe/unsafe.go or builtin/builtin.go, or are subject to cgo preprocessing.
    	GoFiles [][string](/builtin#string)
    
    	// CompiledGoFiles lists the absolute file paths of the package's source
    	// files that are suitable for type checking.
    	// This may differ from GoFiles if files are processed before compilation.
    	CompiledGoFiles [][string](/builtin#string)
    
    	// OtherFiles lists the absolute file paths of the package's non-Go source files,
    	// including assembly, C, C++, Fortran, Objective-C, SWIG, and so on.
    	OtherFiles [][string](/builtin#string)
    
    	// EmbedFiles lists the absolute file paths of the package's files
    	// embedded with go:embed.
    	EmbedFiles [][string](/builtin#string)
    
    	// EmbedPatterns lists the absolute file patterns of the package's
    	// files embedded with go:embed.
    	EmbedPatterns [][string](/builtin#string)
    
    	// IgnoredFiles lists source files that are not part of the package
    	// using the current build configuration but that might be part of
    	// the package using other build configurations.
    	IgnoredFiles [][string](/builtin#string)
    
    	// ExportFile is the absolute path to a file containing type
    	// information for the package as provided by the build system.
    	ExportFile [string](/builtin#string)
    
    	// Target is the absolute install path of the .a file, for libraries,
    	// and of the executable file, for binaries.
    	Target [string](/builtin#string)
    
    	// Imports maps import paths appearing in the package's Go source files
    	// to corresponding loaded Packages.
    	Imports map[[string](/builtin#string)]*Package
    
    	// Module is the module information for the package if it exists.
    	//
    	// Note: it may be missing for std and cmd; see Go issue #65816.
    	Module *Module
    
    	// Types provides type information for the package.
    	// The NeedTypes LoadMode bit sets this field for packages matching the
    	// patterns; type information for dependencies may be missing or incomplete,
    	// unless NeedDeps and NeedImports are also set.
    	//
    	// Each call to [Load] returns a consistent set of type
    	// symbols, as defined by the comment at [types.Identical].
    	// Avoid mixing type information from two or more calls to [Load].
    	Types *[types](/go/types).[Package](/go/types#Package) `json:"-"`
    
    	// Fset provides position information for Types, TypesInfo, and Syntax.
    	// It is set only when Types is set.
    	Fset *[token](/go/token).[FileSet](/go/token#FileSet) `json:"-"`
    
    	// IllTyped indicates whether the package or any dependency contains errors.
    	// It is set only when Types is set.
    	IllTyped [bool](/builtin#bool) `json:"-"`
    
    	// Syntax is the package's syntax trees, for the files listed in CompiledGoFiles.
    	//
    	// The NeedSyntax LoadMode bit populates this field for packages matching the patterns.
    	// If NeedDeps and NeedImports are also set, this field will also be populated
    	// for dependencies.
    	//
    	// Syntax is kept in the same order as CompiledGoFiles, with the caveat that nils are
    	// removed.  If parsing returned nil, Syntax may be shorter than CompiledGoFiles.
    	Syntax []*[ast](/go/ast).[File](/go/ast#File) `json:"-"`
    
    	// TypesInfo provides type information about the package's syntax trees.
    	// It is set only when Syntax is set.
    	TypesInfo *[types](/go/types).[Info](/go/types#Info) `json:"-"`
    
    	// TypesSizes provides the effective size function for types in TypesInfo.
    	TypesSizes [types](/go/types).[Sizes](/go/types#Sizes) `json:"-"`
    
    	// ForTest is the package under test, if any.
    	ForTest [string](/builtin#string)
    	// contains filtered or unexported fields
    }

A Package describes a loaded Go package. 

It also defines part of the JSON schema of DriverResponse. See the package documentation for an overview. 

####  func [Load](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=261) ¶
    
    
    func Load(cfg *Config, patterns ...[string](/builtin#string)) ([]*Package, [error](/builtin#error))

Load loads and returns the Go packages named by the given patterns. 

The cfg parameter specifies loading options; nil behaves the same as an empty Config. 

The [Config.Mode] field is a set of bits that determine what kinds of information should be computed and returned. Modes that require more information tend to be slower. See LoadMode for details and important caveats. Its zero value is equivalent to NeedName | NeedFiles | NeedCompiledGoFiles. 

Each call to Load returns a new set of Package instances. The Packages and their Imports form a directed acyclic graph. 

If the NeedTypes mode flag was set, each call to Load uses a new [types.Importer](/go/types#Importer), so [types.Object](/go/types#Object) and [types.Type](/go/types#Type) values from different calls to Load must not be mixed as they will have inconsistent notions of type identity. 

If any of the patterns was invalid as defined by the underlying build system, Load returns an error. It may return an empty list of packages without an error, for instance for an empty expansion of a valid wildcard. Errors associated with a particular package are recorded in the corresponding Package's Errors list, and do not cause Load to return an error. Clients may need to handle such errors before proceeding with further analysis. The PrintErrors function is provided for convenient display of all errors. 

####  func (*Package) [MarshalJSON](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=624) ¶
    
    
    func (p *Package) MarshalJSON() ([][byte](/builtin#byte), [error](/builtin#error))

MarshalJSON returns the Package in its JSON form. For the most part, the structure fields are written out unmodified, and the type and syntax fields are skipped. The imports are written out as just a map of path to package id. The errors are written using a custom type that tries to preserve the structure of error types we know about. 

This method exists to enable support for additional build systems. It is not intended for use by clients of the API and we may change the format. 

####  func (*Package) [String](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=676) ¶
    
    
    func (p *Package) String() [string](/builtin#string)

####  func (*Package) [UnmarshalJSON](https://cs.opensource.google/go/x/tools/+/v0.41.0:go/packages/packages.go;l=649) ¶
    
    
    func (p *Package) UnmarshalJSON(b [][byte](/builtin#byte)) [error](/builtin#error)

UnmarshalJSON reads in a Package from its JSON format. See MarshalJSON for details about the format accepted. 
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template
