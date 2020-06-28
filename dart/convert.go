// Convert Go code to Dart by generating the Abstract Syntax Tree for Go code.
// Why convert Go to Dart? See this: https://lupyuen.github.io/pinetime-rust-mynewt/articles/companion
// Based on https://golang.org/src/go/ast/example_test.go
// and https://zupzup.org/go-ast-traversal/
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

// src is the input for which we want to generate the Abstract Syntax Tree. "package" is mandatory.
const src = `
package main
type NmpHdr struct {
	Op    uint8 /* 3 bits of opcode */
	Flags uint8
	Len   uint16
	Group uint16
	Seq   uint8
	Id    uint8
}
`

/* Objective: Convert the above Go code to Dart:
class NmpHdr {
  int Op;    //  uint8: 3 bits of opcode
  int Flags; //  uint8
  int Len;   //  uint16
  int Group; //  uint16
  int Seq;   //  uint8
  int Id;    //  uint8

  /// Construct an SMP Header
  NmpHdr(
    this.Op,    //  uint8: 3 bits of opcode
    this.Flags, //  uint8
    this.Len,   //  uint16
    this.Group, //  uint16
    this.Seq,   //  uint8
    this.Id     //  uint8
  );
} */

// Inspect the Abstract Syntax Tree of our Go code
func inspectAST() {
	fmt.Printf("//  Convert Go Code...%s\n", src)
	fmt.Println("//  To Dart...")

	// Create the AST by parsing src
	fileset := token.NewFileSet()                            // Positions are relative to fileset
	node, err := parser.ParseFile(fileset, "src.go", src, 0) // Change src to nil to parse a file instead of string
	if err != nil {
		panic(err)
	}

	// Process all declarations
	for _, decl := range node.Decls {
		// Process a declaration
		// fmt.Println("Decl:")
		// ast.Print(fileset, decl)
		switch decl := decl.(type) {
		case *ast.GenDecl:
			// fmt.Printf("Tok: %s\n", decl.Tok) // "type"
			switch decl.Tok.String() {
			case "type":
				// Process a type declaration
				for _, spec := range decl.Specs {
					// ast.Print(fileset, spec)
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						typeName := spec.Name.Name // "NmpHdr"
						// fmt.Printf("typeName: %s\n", typeName)
						fmt.Printf("class %s {\n", typeName)
						switch structType := spec.Type.(type) {
						case *ast.StructType: // "struct {"
							// Process a struct declaration
							// ast.Print(fileset, structType)
							for _, field := range structType.Fields.List {
								// Process a struct field and type
								// ast.Print(fileset, field)
								fieldName := field.Names[0].Name                      // "Op"
								fieldType := field.Type                               // "uint8"
								dartType := convertType(fmt.Sprintf("%s", fieldType)) // "int"
								// fmt.Printf("field: %s,\ttype: %s\n", fieldName, fieldType)
								fmt.Printf("  %s %s;\t//  %s\n", dartType, fieldName, fieldType)
							}

						default:
							fmt.Println("*** Unknown Spec Type:")
							ast.Print(fileset, spec.Type)
						}

					default:
						fmt.Println("*** Unknown Spec:")
						ast.Print(fileset, spec)
					}
				}
				fmt.Println("}")

			default:
				fmt.Println("*** Unknown Tok:")
				ast.Print(fileset, decl.Tok)
			}

		default:
			fmt.Println("*** Unknown Decl:")
			ast.Print(fileset, decl)
		}
		// fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
	}
}

// Convert Go type to Dart type
func convertType(typeName string) string {
	switch typeName {
	case "uint8":
		return "int"
	case "uint16":
		return "int"
	default:
		return "Unknown"
	}
}

// This example shows what an AST looks like when printed for debugging.
func ExamplePrint() {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}

/* Output:
     0  *ast.File {
     1  .  Package: 2:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 2:9
     4  .  .  Name: "dummy_package"
     5  .  }
     6  .  Decls: []ast.Decl (len = 1) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 3:1
     9  .  .  .  Tok: type
    10  .  .  .  Lparen: -
    11  .  .  .  Specs: []ast.Spec (len = 1) {
    12  .  .  .  .  0: *ast.TypeSpec {
    13  .  .  .  .  .  Name: *ast.Ident {
    14  .  .  .  .  .  .  NamePos: 3:6
    15  .  .  .  .  .  .  Name: "NmpHdr"
    16  .  .  .  .  .  .  Obj: *ast.Object {
    17  .  .  .  .  .  .  .  Kind: type
    18  .  .  .  .  .  .  .  Name: "NmpHdr"
    19  .  .  .  .  .  .  .  Decl: *(obj @ 12)
    20  .  .  .  .  .  .  }
    21  .  .  .  .  .  }
    22  .  .  .  .  .  Assign: -
    23  .  .  .  .  .  Type: *ast.StructType {
    24  .  .  .  .  .  .  Struct: 3:13
    25  .  .  .  .  .  .  Fields: *ast.FieldList {
    26  .  .  .  .  .  .  .  Opening: 3:20
    27  .  .  .  .  .  .  .  List: []*ast.Field (len = 6) {
    28  .  .  .  .  .  .  .  .  0: *ast.Field {
    29  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    30  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    31  .  .  .  .  .  .  .  .  .  .  .  NamePos: 4:2
    32  .  .  .  .  .  .  .  .  .  .  .  Name: "Op"
    33  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    34  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    35  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Op"
    36  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 28)
    37  .  .  .  .  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  .  .  .  .  }
    39  .  .  .  .  .  .  .  .  .  }
    40  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    41  .  .  .  .  .  .  .  .  .  .  NamePos: 4:8
    42  .  .  .  .  .  .  .  .  .  .  Name: "uint8"
    43  .  .  .  .  .  .  .  .  .  }
    44  .  .  .  .  .  .  .  .  }
    45  .  .  .  .  .  .  .  .  1: *ast.Field {
    46  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    47  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    48  .  .  .  .  .  .  .  .  .  .  .  NamePos: 5:2
    49  .  .  .  .  .  .  .  .  .  .  .  Name: "Flags"
    50  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    51  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    52  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Flags"
    53  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 45)
    54  .  .  .  .  .  .  .  .  .  .  .  }
    55  .  .  .  .  .  .  .  .  .  .  }
    56  .  .  .  .  .  .  .  .  .  }
    57  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    58  .  .  .  .  .  .  .  .  .  .  NamePos: 5:8
    59  .  .  .  .  .  .  .  .  .  .  Name: "uint8"
    60  .  .  .  .  .  .  .  .  .  }
    61  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  .  2: *ast.Field {
    63  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    64  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    65  .  .  .  .  .  .  .  .  .  .  .  NamePos: 6:2
    66  .  .  .  .  .  .  .  .  .  .  .  Name: "Len"
    67  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    68  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    69  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Len"
    70  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 62)
    71  .  .  .  .  .  .  .  .  .  .  .  }
    72  .  .  .  .  .  .  .  .  .  .  }
    73  .  .  .  .  .  .  .  .  .  }
    74  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    75  .  .  .  .  .  .  .  .  .  .  NamePos: 6:8
    76  .  .  .  .  .  .  .  .  .  .  Name: "uint16"
    77  .  .  .  .  .  .  .  .  .  }
    78  .  .  .  .  .  .  .  .  }
    79  .  .  .  .  .  .  .  .  3: *ast.Field {
    80  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    81  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    82  .  .  .  .  .  .  .  .  .  .  .  NamePos: 7:2
    83  .  .  .  .  .  .  .  .  .  .  .  Name: "Group"
    84  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    85  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    86  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Group"
    87  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 79)
    88  .  .  .  .  .  .  .  .  .  .  .  }
    89  .  .  .  .  .  .  .  .  .  .  }
    90  .  .  .  .  .  .  .  .  .  }
    91  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
    92  .  .  .  .  .  .  .  .  .  .  NamePos: 7:8
    93  .  .  .  .  .  .  .  .  .  .  Name: "uint16"
    94  .  .  .  .  .  .  .  .  .  }
    95  .  .  .  .  .  .  .  .  }
    96  .  .  .  .  .  .  .  .  4: *ast.Field {
    97  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    98  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    99  .  .  .  .  .  .  .  .  .  .  .  NamePos: 8:2
   100  .  .  .  .  .  .  .  .  .  .  .  Name: "Seq"
   101  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   102  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
   103  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Seq"
   104  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 96)
   105  .  .  .  .  .  .  .  .  .  .  .  }
   106  .  .  .  .  .  .  .  .  .  .  }
   107  .  .  .  .  .  .  .  .  .  }
   108  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
   109  .  .  .  .  .  .  .  .  .  .  NamePos: 8:8
   110  .  .  .  .  .  .  .  .  .  .  Name: "uint8"
   111  .  .  .  .  .  .  .  .  .  }
   112  .  .  .  .  .  .  .  .  }
   113  .  .  .  .  .  .  .  .  5: *ast.Field {
   114  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   115  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   116  .  .  .  .  .  .  .  .  .  .  .  NamePos: 9:2
   117  .  .  .  .  .  .  .  .  .  .  .  Name: "Id"
   118  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   119  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
   120  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Id"
   121  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 113)
   122  .  .  .  .  .  .  .  .  .  .  .  }
   123  .  .  .  .  .  .  .  .  .  .  }
   124  .  .  .  .  .  .  .  .  .  }
   125  .  .  .  .  .  .  .  .  .  Type: *ast.Ident {
   126  .  .  .  .  .  .  .  .  .  .  NamePos: 9:8
   127  .  .  .  .  .  .  .  .  .  .  Name: "uint8"
   128  .  .  .  .  .  .  .  .  .  }
   129  .  .  .  .  .  .  .  .  }
   130  .  .  .  .  .  .  .  }
   131  .  .  .  .  .  .  .  Closing: 10:1
   132  .  .  .  .  .  .  }
   133  .  .  .  .  .  .  Incomplete: false
   134  .  .  .  .  .  }
   135  .  .  .  .  }
   136  .  .  .  }
   137  .  .  .  Rparen: -
   138  .  .  }
   139  .  }
   140  .  Scope: *ast.Scope {
   141  .  .  Objects: map[string]*ast.Object (len = 1) {
   142  .  .  .  "NmpHdr": *(obj @ 16)
   143  .  .  }
   144  .  }
   145  .  Unresolved: []*ast.Ident (len = 6) {
   146  .  .  0: *(obj @ 40)
   147  .  .  1: *(obj @ 57)
   148  .  .  2: *(obj @ 74)
   149  .  .  3: *(obj @ 91)
   150  .  .  4: *(obj @ 108)
   151  .  .  5: *(obj @ 125)
   152  .  }
   153  }
*/

// This example demonstrates how to inspect the AST of a Go program.
func ExampleInspect() {
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

/*
	src.go:2:9:     dummy_package
	src.go:3:6:     NmpHdr
	src.go:4:2:     Op
	src.go:4:8:     uint8
	src.go:5:2:     Flags
	src.go:5:8:     uint8
	src.go:6:2:     Len
	src.go:6:8:     uint16
	src.go:7:2:     Group
	src.go:7:8:     uint16
	src.go:8:2:     Seq
	src.go:8:8:     uint8
	src.go:9:2:     Id
	src.go:9:8:     uint8
*/

// This example illustrates how to remove a variable declaration
// in a Go program while maintaining correct comment association
// using an ast.CommentMap.
func ExampleCommentMap() {
	/* // src is the input for which we create the AST that we
	// are going to manipulate.
	*/

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
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())

	// Output:
	// // This is the package comment.
	// package main
	//
	// // This comment is associated with the hello constant.
	// const hello = "Hello, World!" // line comment 1
	//
	// // This comment is associated with the main function.
	// func main() {
	// 	fmt.Println(hello) // line comment 3
	// }
}

func main() {
	inspectAST()
	// ExamplePrint()
	// ExampleInspect()
	// ExampleCommentMap()
}

/* Previously:
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
`
src.go:2:9:     p
src.go:3:7:     c
src.go:3:11:    1.0
src.go:4:5:     X
src.go:4:9:     f
src.go:4:11:    3.14
src.go:4:17:    2
src.go:4:21:    c

	src := `
package main
func main() {
	println("Hello, World!")
}
`
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
    50  .  Scope: *ast.Scope {
    51  .  .  Objects: map[string]*ast.Object (len = 1) {
    52  .  .  .  "main": *(obj @ 11)
    53  .  .  }
    54  .  }
    55  .  Unresolved: []*ast.Ident (len = 1) {
    56  .  .  0: *(obj @ 29)
    57  .  }
	58  }

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
// This is the package comment.
package main

// This comment is associated with the hello constant.
const hello = "Hello, World!" // line comment 1

// This comment is associated with the main function.
func main() {
        fmt.Println(hello) // line comment 3
}

*/
