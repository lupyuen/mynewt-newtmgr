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
	"strings"
)

// src is the input for which we want to generate the Abstract Syntax Tree. "package" is mandatory.
// bt means backtick "`"
const src = `
package main
type ImageUploadReq struct {
	NmpBase  ` + bt + `codec:"-"` + bt + `
	ImageNum uint8  ` + bt + `codec:"image"` + bt + `
	Off      uint32 ` + bt + `codec:"off"` + bt + `
	Len      uint32 ` + bt + `codec:"len,omitempty"` + bt + `
	DataSha  []byte ` + bt + `codec:"sha,omitempty"` + bt + `
	Upgrade  bool   ` + bt + `codec:"upgrade,omitempty"` + bt + `
	Data     []byte ` + bt + `codec:"data"` + bt + `
}
`

const bt = "`" //  Backtick character

/* Objective: Convert the above Go code to Dart:
class ImageUploadReq
  with NmpBase       //  Get and set SMP Message Header
  implements NmpReq  //  SMP Request Message
{
  int ImageNum; //  uint8
  int Off;      //  uint32
  int Len;      //  uint32
  typed.Uint8Buffer DataSha;    //  []byte
  bool Upgrade; //  bool
  typed.Uint8Buffer Data;       //  []byte

  NmpMsg Msg() { return MsgFromReq(this); }

  /// Encode the SMP Request fields to CBOR
  void Encode(cbor.MapBuilder builder) {
    builder.writeString("image");
    builder.writeString(ImageNum);
    builder.writeString("off");
    builder.writeString(Off);
    builder.writeString("len");
    builder.writeString(Len);
    builder.writeString("sha");
    builder.writeString(DataSha);
    builder.writeString("upgrade");
    builder.writeString(Upgrade);
    builder.writeString("data");
    builder.writeString(Data);
  }
} */

// Inspect the Abstract Syntax Tree of our Go code and convert to Dart
func convertGoToDart() {
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
						fmt.Printf("class %s ", typeName)
						// Handle request messages
						if strings.HasSuffix(typeName, "Req") {
							fmt.Println("\n  with NmpBase       //  Get and set SMP Message Header\n  implements NmpReq  //  SMP Request Message")
						}
						fmt.Println("{")

						switch structType := spec.Type.(type) {
						case *ast.StructType: // "struct {"
							// Process a struct declaration
							// ast.Print(fileset, structType)
							fields := structType.Fields.List
							convertFields(fields)
							fmt.Println("")

							// Handle request messages
							if strings.HasSuffix(typeName, "Req") {
								fmt.Println("  NmpMsg Msg() { return MsgFromReq(this); }\n")
							}

							// Generate CBOR encoder
							generateCborEncoder(fields)

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

// DartField represents a Go Struct Field converted to Dart
type DartField struct {
	Name     string // "Op"
	CborName string
	GoType   string // "uint8"
	DartType string // "int"
}

// Generate the CBOR Encoder function
func generateCborEncoder(astFields []*ast.Field) {
	fmt.Println("  /// Encode the SMP Request fields to CBOR")
	fmt.Println("  void Encode(cbor.MapBuilder builder) {")
	for _, field := range astFields {
		// Process a struct field and type
		// ast.Print(fileset, field)
		dartField := convertField(field)
		if dartField.CborName != "-" {
			fmt.Printf("    builder.writeString(\"%s\");\n", dartField.CborName)
			// TODO: Handle type
			fmt.Printf("    builder.writeString(%s);\n", dartField.Name)
		}
	}
	fmt.Println("  }")
}

// Convert Go Struct Fields to Dart
func convertFields(astFields []*ast.Field) {
	for _, field := range astFields {
		// Process a struct field and type
		// ast.Print(fileset, field)
		dartField := convertField(field)
		if dartField.Name != "" {
			fmt.Printf("  %s %s;\t//  %s\n", dartField.DartType, dartField.Name, dartField.GoType)
		}
	}
}

// Convert Go Struct Field to Dart
func convertField(astField *ast.Field) DartField {
	dartField := DartField{}
	if len(astField.Names) > 0 {
		dartField.Name = astField.Names[0].Name // "Op"
	}
	dartField.GoType = fmt.Sprintf("%v", astField.Type) // "uint8"
	// Handle "&{181 <nil> byte}" as "[]byte"
	if strings.HasPrefix(dartField.GoType, "&{") && strings.HasSuffix(dartField.GoType, " byte}") {
		dartField.GoType = "[]byte"
	}
	dartField.DartType = convertType(fmt.Sprintf("%s", dartField.GoType)) // "int"

	// Convert a Field Tag like `codec:"len,omitempty"`. CborName will be set to "len".
	if astField.Tag != nil {
		dartField.CborName = strings.Split(astField.Tag.Value, ",")[0]
		dartField.CborName = strings.Replace(dartField.CborName, "codec:", "", 1)
		dartField.CborName = strings.Replace(dartField.CborName, `"`, "", 2)
		dartField.CborName = strings.Replace(dartField.CborName, "`", "", 2)
	}
	// fmt.Printf("field: %s,\tcbor: %s,\ttype: %s,\tdart: %s\n", dartField.Name, dartField.CborName, dartField.GoType, dartField.DartType)
	return dartField
}

// Convert Go type to Dart type
func convertType(typeName string) string {
	switch typeName {
	case "bool":
		return "bool"
	case "uint8":
		return "int"
	case "uint16":
		return "int"
	case "uint32":
		return "int"
	case "[]byte":
		return "typed.Uint8Buffer"
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
	convertGoToDart()
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
