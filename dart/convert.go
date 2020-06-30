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

// src is the Go code to be converted to Dart. "package" is mandatory. "bt" means backtick "`"
const src = `
package main
func NewImageUploadReq() *ImageUploadReq {
	r := &ImageUploadReq{}
	fillNmpReq(r, NMP_OP_WRITE, NMP_GROUP_IMAGE, NMP_ID_IMAGE_UPLOAD)
	return r
}
`

/* Objective: Convert the above Go code to Dart:
ImageUploadReq NewImageUploadReq {
  var r = ImageUploadReq();
  fillNmpReq(r, NMP_OP_WRITE, NMP_GROUP_IMAGE, NMP_ID_IMAGE_UPLOAD);
  return r;
} */

const src2 = `
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
    builder.writeInt(ImageNum);
    builder.writeString("off");
    builder.writeInt(Off);
    builder.writeString("len");
    builder.writeInt(Len);
    builder.writeString("sha");
    builder.writeArray(DataSha);
    builder.writeString("upgrade");
    builder.writeBool(Upgrade);
    builder.writeString("data");
    builder.writeArray(Data);
  }
} */

const bt = "`" //  Backtick character

// Inspect the Abstract Syntax Tree of our Go code and convert to Dart
func convertGoToDart() {
	fmt.Printf("//  Go Code...%s\n", src)
	fmt.Println("//  Converted To Dart...")

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
			// Convert Go Struct to Dart
			convertStruct(fileset, decl)
		case *ast.FuncDecl:
			// Convert Go Function to Dart
			convertFunction(fileset, decl)
		default:
			fmt.Println("*** Unknown Decl:")
			ast.Print(fileset, decl)
		}
		// fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
	}
}

// Convert Go Function to Dart
func convertFunction(fileset *token.FileSet, decl *ast.FuncDecl) {
	// ast.Print(fileset, decl)
	name := decl.Name                                               // "NewImageUploadReq"
	returnType := fmt.Sprintf("%v", decl.Type.Results.List[0].Type) // "&{40 ImageUploadReq}"
	// Convert the return type "&{40 ImageUploadReq}" to "ImageUploadReq"
	if strings.HasPrefix(returnType, "&{") && strings.HasSuffix(returnType, "}") {
		returnType = strings.Split(returnType, " ")[1]
		returnType = strings.Replace(returnType, "}", "", 1)
	}
	fmt.Printf("%s %s() {\n", returnType, name)
	// Convert the body
	body := decl.Body.List
	for _, stmt := range body {
		// ast.Print(fileset, stmt)
		// Convert the statement to a string
		var buf bytes.Buffer
		if err := format.Node(&buf, fileset, stmt); err != nil {
			panic(err)
		}
		dartStmt := fmt.Sprintf("%s", buf.Bytes())
		// Convert specific kinds of statements
		switch stmt.(type) {
		case *ast.AssignStmt:
			// For Go assignment "r := &ImageUploadReq{}", rewrite to "var r = ImageUploadReq()"
			dartStmt = strings.Replace(dartStmt, ":=", "=", 1)
			dartStmt = strings.Replace(dartStmt, "&", "", 1)
			dartStmt = strings.Replace(dartStmt, "{}", "()", 1)
			dartStmt = "var " + dartStmt
		}
		fmt.Printf("  %s;\n", dartStmt)
	}
	fmt.Println("}\n")
}

/*
     0  *ast.FuncDecl {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: src.go:3:6
     3  .  .  Name: "NewImageUploadReq"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: func
     6  .  .  .  Name: "NewImageUploadReq"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.FuncType {
    11  .  .  Func: src.go:3:1
    12  .  .  Params: *ast.FieldList {
    13  .  .  .  Opening: src.go:3:23
    14  .  .  .  Closing: src.go:3:24
    15  .  .  }
    16  .  .  Results: *ast.FieldList {
    17  .  .  .  Opening: -
    18  .  .  .  List: []*ast.Field (len = 1) {
    19  .  .  .  .  0: *ast.Field {
    20  .  .  .  .  .  Type: *ast.StarExpr {
    21  .  .  .  .  .  .  Star: src.go:3:26
    22  .  .  .  .  .  .  X: *ast.Ident {
    23  .  .  .  .  .  .  .  NamePos: src.go:3:27
    24  .  .  .  .  .  .  .  Name: "ImageUploadReq"
    25  .  .  .  .  .  .  }
    26  .  .  .  .  .  }
    27  .  .  .  .  }
    28  .  .  .  }
    29  .  .  .  Closing: -
    30  .  .  }
    31  .  }
    32  .  Body: *ast.BlockStmt {
    33  .  .  Lbrace: src.go:3:42
    34  .  .  List: []ast.Stmt (len = 3) {
    35  .  .  .  0: *ast.AssignStmt {
    36  .  .  .  .  Lhs: []ast.Expr (len = 1) {
    37  .  .  .  .  .  0: *ast.Ident {
    38  .  .  .  .  .  .  NamePos: src.go:4:2
    39  .  .  .  .  .  .  Name: "r"
    40  .  .  .  .  .  .  Obj: *ast.Object {
    41  .  .  .  .  .  .  .  Kind: var
    42  .  .  .  .  .  .  .  Name: "r"
    43  .  .  .  .  .  .  .  Decl: *(obj @ 35)
    44  .  .  .  .  .  .  }
    45  .  .  .  .  .  }
    46  .  .  .  .  }
    47  .  .  .  .  TokPos: src.go:4:4
    48  .  .  .  .  Tok: :=
    49  .  .  .  .  Rhs: []ast.Expr (len = 1) {
    50  .  .  .  .  .  0: *ast.UnaryExpr {
    51  .  .  .  .  .  .  OpPos: src.go:4:7
    52  .  .  .  .  .  .  Op: &
    53  .  .  .  .  .  .  X: *ast.CompositeLit {
    54  .  .  .  .  .  .  .  Type: *ast.Ident {
    55  .  .  .  .  .  .  .  .  NamePos: src.go:4:8
    56  .  .  .  .  .  .  .  .  Name: "ImageUploadReq"
    57  .  .  .  .  .  .  .  }
    58  .  .  .  .  .  .  .  Lbrace: src.go:4:22
    59  .  .  .  .  .  .  .  Rbrace: src.go:4:23
    60  .  .  .  .  .  .  .  Incomplete: false
    61  .  .  .  .  .  .  }
    62  .  .  .  .  .  }
    63  .  .  .  .  }
    64  .  .  .  }
    65  .  .  .  1: *ast.ExprStmt {
    66  .  .  .  .  X: *ast.CallExpr {
    67  .  .  .  .  .  Fun: *ast.Ident {
    68  .  .  .  .  .  .  NamePos: src.go:5:2
    69  .  .  .  .  .  .  Name: "fillNmpReq"
    70  .  .  .  .  .  }
    71  .  .  .  .  .  Lparen: src.go:5:12
    72  .  .  .  .  .  Args: []ast.Expr (len = 4) {
    73  .  .  .  .  .  .  0: *ast.Ident {
    74  .  .  .  .  .  .  .  NamePos: src.go:5:13
    75  .  .  .  .  .  .  .  Name: "r"
    76  .  .  .  .  .  .  .  Obj: *(obj @ 40)
    77  .  .  .  .  .  .  }
    78  .  .  .  .  .  .  1: *ast.Ident {
    79  .  .  .  .  .  .  .  NamePos: src.go:5:16
    80  .  .  .  .  .  .  .  Name: "NMP_OP_WRITE"
    81  .  .  .  .  .  .  }
    82  .  .  .  .  .  .  2: *ast.Ident {
    83  .  .  .  .  .  .  .  NamePos: src.go:5:30
    84  .  .  .  .  .  .  .  Name: "NMP_GROUP_IMAGE"
    85  .  .  .  .  .  .  }
    86  .  .  .  .  .  .  3: *ast.Ident {
    87  .  .  .  .  .  .  .  NamePos: src.go:5:47
    88  .  .  .  .  .  .  .  Name: "NMP_ID_IMAGE_UPLOAD"
    89  .  .  .  .  .  .  }
    90  .  .  .  .  .  }
    91  .  .  .  .  .  Ellipsis: -
    92  .  .  .  .  .  Rparen: src.go:5:66
    93  .  .  .  .  }
    94  .  .  .  }
    95  .  .  .  2: *ast.ReturnStmt {
    96  .  .  .  .  Return: src.go:6:2
    97  .  .  .  .  Results: []ast.Expr (len = 1) {
    98  .  .  .  .  .  0: *ast.Ident {
    99  .  .  .  .  .  .  NamePos: src.go:6:9
   100  .  .  .  .  .  .  Name: "r"
   101  .  .  .  .  .  .  Obj: *(obj @ 40)
   102  .  .  .  .  .  }
   103  .  .  .  .  }
   104  .  .  .  }
   105  .  .  }
   106  .  .  Rbrace: src.go:7:1
   107  .  }
   108  }
*/

// Convert Go Struct to Dart
func convertStruct(fileset *token.FileSet, decl *ast.GenDecl) {
	// ast.Print(fileset, decl)
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
					fmt.Println("")
					fmt.Println("  with NmpBase       //  Get and set SMP Message Header")
					fmt.Println("  implements NmpReq  //  SMP Request Message")
				}
				fmt.Println("{")

				switch structType := spec.Type.(type) {
				case *ast.StructType: // "struct {"
					// Process a struct declaration
					// ast.Print(fileset, structType)
					fields := structType.Fields.List
					convertFields(fileset, fields)
					fmt.Println("")

					// Handle request messages
					if strings.HasSuffix(typeName, "Req") {
						fmt.Println("  NmpMsg Msg() { return MsgFromReq(this); }\n")
					}

					// Generate CBOR encoder
					generateCborEncoder(fileset, fields)

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
}

/*
     0  *ast.GenDecl {
     1  .  TokPos: src.go:3:1
     2  .  Tok: type
     3  .  Lparen: -
     4  .  Specs: []ast.Spec (len = 1) {
     5  .  .  0: *ast.TypeSpec {
     6  .  .  .  Name: *ast.Ident {
     7  .  .  .  .  NamePos: src.go:3:6
     8  .  .  .  .  Name: "ImageUploadReq"
     9  .  .  .  .  Obj: *ast.Object {
    10  .  .  .  .  .  Kind: type
    11  .  .  .  .  .  Name: "ImageUploadReq"
    12  .  .  .  .  .  Decl: *(obj @ 5)
    13  .  .  .  .  }
    14  .  .  .  }
    15  .  .  .  Assign: -
    16  .  .  .  Type: *ast.StructType {
    17  .  .  .  .  Struct: src.go:3:21
    18  .  .  .  .  Fields: *ast.FieldList {
    19  .  .  .  .  .  Opening: src.go:3:28
    20  .  .  .  .  .  List: []*ast.Field (len = 7) {
    21  .  .  .  .  .  .  0: *ast.Field {
    22  .  .  .  .  .  .  .  Type: *ast.Ident {
    23  .  .  .  .  .  .  .  .  NamePos: src.go:4:2
    24  .  .  .  .  .  .  .  .  Name: "NmpBase"
    25  .  .  .  .  .  .  .  }
    26  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
    27  .  .  .  .  .  .  .  .  ValuePos: src.go:4:11
    28  .  .  .  .  .  .  .  .  Kind: STRING
    29  .  .  .  .  .  .  .  .  Value: "`codec:\"-\"`"
    30  .  .  .  .  .  .  .  }
    31  .  .  .  .  .  .  }
    32  .  .  .  .  .  .  1: *ast.Field {
    33  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    34  .  .  .  .  .  .  .  .  0: *ast.Ident {
    35  .  .  .  .  .  .  .  .  .  NamePos: src.go:5:2
    36  .  .  .  .  .  .  .  .  .  Name: "ImageNum"
    37  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    38  .  .  .  .  .  .  .  .  .  .  Kind: var
    39  .  .  .  .  .  .  .  .  .  .  Name: "ImageNum"
    40  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 32)
    41  .  .  .  .  .  .  .  .  .  }
    42  .  .  .  .  .  .  .  .  }
    43  .  .  .  .  .  .  .  }
    44  .  .  .  .  .  .  .  Type: *ast.Ident {
    45  .  .  .  .  .  .  .  .  NamePos: src.go:5:11
    46  .  .  .  .  .  .  .  .  Name: "uint8"
    47  .  .  .  .  .  .  .  }
    48  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
    49  .  .  .  .  .  .  .  .  ValuePos: src.go:5:18
    50  .  .  .  .  .  .  .  .  Kind: STRING
    51  .  .  .  .  .  .  .  .  Value: "`codec:\"image\"`"
    52  .  .  .  .  .  .  .  }
    53  .  .  .  .  .  .  }
    54  .  .  .  .  .  .  2: *ast.Field {
    55  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    56  .  .  .  .  .  .  .  .  0: *ast.Ident {
    57  .  .  .  .  .  .  .  .  .  NamePos: src.go:6:2
    58  .  .  .  .  .  .  .  .  .  Name: "Off"
    59  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    60  .  .  .  .  .  .  .  .  .  .  Kind: var
    61  .  .  .  .  .  .  .  .  .  .  Name: "Off"
    62  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 54)
    63  .  .  .  .  .  .  .  .  .  }
    64  .  .  .  .  .  .  .  .  }
    65  .  .  .  .  .  .  .  }
    66  .  .  .  .  .  .  .  Type: *ast.Ident {
    67  .  .  .  .  .  .  .  .  NamePos: src.go:6:11
    68  .  .  .  .  .  .  .  .  Name: "uint32"
    69  .  .  .  .  .  .  .  }
    70  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
    71  .  .  .  .  .  .  .  .  ValuePos: src.go:6:18
    72  .  .  .  .  .  .  .  .  Kind: STRING
    73  .  .  .  .  .  .  .  .  Value: "`codec:\"off\"`"
    74  .  .  .  .  .  .  .  }
    75  .  .  .  .  .  .  }
    76  .  .  .  .  .  .  3: *ast.Field {
    77  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    78  .  .  .  .  .  .  .  .  0: *ast.Ident {
    79  .  .  .  .  .  .  .  .  .  NamePos: src.go:7:2
    80  .  .  .  .  .  .  .  .  .  Name: "Len"
    81  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    82  .  .  .  .  .  .  .  .  .  .  Kind: var
    83  .  .  .  .  .  .  .  .  .  .  Name: "Len"
    84  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 76)
    85  .  .  .  .  .  .  .  .  .  }
    86  .  .  .  .  .  .  .  .  }
    87  .  .  .  .  .  .  .  }
    88  .  .  .  .  .  .  .  Type: *ast.Ident {
    89  .  .  .  .  .  .  .  .  NamePos: src.go:7:11
    90  .  .  .  .  .  .  .  .  Name: "uint32"
    91  .  .  .  .  .  .  .  }
    92  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
    93  .  .  .  .  .  .  .  .  ValuePos: src.go:7:18
    94  .  .  .  .  .  .  .  .  Kind: STRING
    95  .  .  .  .  .  .  .  .  Value: "`codec:\"len,omitempty\"`"
    96  .  .  .  .  .  .  .  }
    97  .  .  .  .  .  .  }
    98  .  .  .  .  .  .  4: *ast.Field {
    99  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   100  .  .  .  .  .  .  .  .  0: *ast.Ident {
   101  .  .  .  .  .  .  .  .  .  NamePos: src.go:8:2
   102  .  .  .  .  .  .  .  .  .  Name: "DataSha"
   103  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   104  .  .  .  .  .  .  .  .  .  .  Kind: var
   105  .  .  .  .  .  .  .  .  .  .  Name: "DataSha"
   106  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 98)
   107  .  .  .  .  .  .  .  .  .  }
   108  .  .  .  .  .  .  .  .  }
   109  .  .  .  .  .  .  .  }
   110  .  .  .  .  .  .  .  Type: *ast.ArrayType {
   111  .  .  .  .  .  .  .  .  Lbrack: src.go:8:11
   112  .  .  .  .  .  .  .  .  Elt: *ast.Ident {
   113  .  .  .  .  .  .  .  .  .  NamePos: src.go:8:13
   114  .  .  .  .  .  .  .  .  .  Name: "byte"
   115  .  .  .  .  .  .  .  .  }
   116  .  .  .  .  .  .  .  }
   117  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
   118  .  .  .  .  .  .  .  .  ValuePos: src.go:8:18
   119  .  .  .  .  .  .  .  .  Kind: STRING
   120  .  .  .  .  .  .  .  .  Value: "`codec:\"sha,omitempty\"`"
   121  .  .  .  .  .  .  .  }
   122  .  .  .  .  .  .  }
   123  .  .  .  .  .  .  5: *ast.Field {
   124  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   125  .  .  .  .  .  .  .  .  0: *ast.Ident {
   126  .  .  .  .  .  .  .  .  .  NamePos: src.go:9:2
   127  .  .  .  .  .  .  .  .  .  Name: "Upgrade"
   128  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   129  .  .  .  .  .  .  .  .  .  .  Kind: var
   130  .  .  .  .  .  .  .  .  .  .  Name: "Upgrade"
   131  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 123)
   132  .  .  .  .  .  .  .  .  .  }
   133  .  .  .  .  .  .  .  .  }
   134  .  .  .  .  .  .  .  }
   135  .  .  .  .  .  .  .  Type: *ast.Ident {
   136  .  .  .  .  .  .  .  .  NamePos: src.go:9:11
   137  .  .  .  .  .  .  .  .  Name: "bool"
   138  .  .  .  .  .  .  .  }
   139  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
   140  .  .  .  .  .  .  .  .  ValuePos: src.go:9:18
   141  .  .  .  .  .  .  .  .  Kind: STRING
   142  .  .  .  .  .  .  .  .  Value: "`codec:\"upgrade,omitempty\"`"
   143  .  .  .  .  .  .  .  }
   144  .  .  .  .  .  .  }
   145  .  .  .  .  .  .  6: *ast.Field {
   146  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
   147  .  .  .  .  .  .  .  .  0: *ast.Ident {
   148  .  .  .  .  .  .  .  .  .  NamePos: src.go:10:2
   149  .  .  .  .  .  .  .  .  .  Name: "Data"
   150  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   151  .  .  .  .  .  .  .  .  .  .  Kind: var
   152  .  .  .  .  .  .  .  .  .  .  Name: "Data"
   153  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 145)
   154  .  .  .  .  .  .  .  .  .  }
   155  .  .  .  .  .  .  .  .  }
   156  .  .  .  .  .  .  .  }
   157  .  .  .  .  .  .  .  Type: *ast.ArrayType {
   158  .  .  .  .  .  .  .  .  Lbrack: src.go:10:11
   159  .  .  .  .  .  .  .  .  Elt: *ast.Ident {
   160  .  .  .  .  .  .  .  .  .  NamePos: src.go:10:13
   161  .  .  .  .  .  .  .  .  .  Name: "byte"
   162  .  .  .  .  .  .  .  .  }
   163  .  .  .  .  .  .  .  }
   164  .  .  .  .  .  .  .  Tag: *ast.BasicLit {
   165  .  .  .  .  .  .  .  .  ValuePos: src.go:10:18
   166  .  .  .  .  .  .  .  .  Kind: STRING
   167  .  .  .  .  .  .  .  .  Value: "`codec:\"data\"`"
   168  .  .  .  .  .  .  .  }
   169  .  .  .  .  .  .  }
   170  .  .  .  .  .  }
   171  .  .  .  .  .  Closing: src.go:11:1
   172  .  .  .  .  }
   173  .  .  .  .  Incomplete: false
   174  .  .  .  }
   175  .  .  }
   176  .  }
   177  .  Rparen: -
   178  }
*/

// DartField represents a Go Struct Field converted to Dart and CBOR
type DartField struct {
	Name     string // "Len"
	CborName string // "len"
	GoType   string // "uint32"
	DartType string // "int"
	CborType string // "Int"
}

// Generate the CBOR Encoder function
func generateCborEncoder(fileset *token.FileSet, astFields []*ast.Field) {
	fmt.Println("  /// Encode the SMP Request fields to CBOR")
	fmt.Println("  void Encode(cbor.MapBuilder builder) {")
	for _, field := range astFields {
		// ast.Print(fileset, field)
		dartField := convertField(fileset, field)
		if dartField.CborName != "-" {
			fmt.Printf("    builder.writeString(\"%s\");\n", dartField.CborName)
			fmt.Printf("    builder.write%s(%s);\n", dartField.CborType, dartField.Name)
		}
	}
	fmt.Println("  }")
}

// Convert Go Struct Fields to Dart
func convertFields(fileset *token.FileSet, astFields []*ast.Field) {
	for _, field := range astFields {
		// ast.Print(fileset, field)
		dartField := convertField(fileset, field)
		if dartField.Name != "" {
			fmt.Printf("  %s %s;\t//  %s\n", dartField.DartType, dartField.Name, dartField.GoType)
		}
	}
}

// Convert Go Struct Field to Dart
func convertField(fileset *token.FileSet, astField *ast.Field) DartField {
	dartField := DartField{}
	if len(astField.Names) > 0 {
		dartField.Name = astField.Names[0].Name // "Len"
	}
	dartField.GoType = fmt.Sprintf("%v", astField.Type) // "uint32"
	// Handle "&{181 <nil> byte}" as "[]byte"
	if strings.HasPrefix(dartField.GoType, "&{") && strings.HasSuffix(dartField.GoType, " byte}") {
		dartField.GoType = "[]byte"
	}
	dartField.DartType, dartField.CborType = convertType(dartField.GoType) // "int"

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

// Convert Go type to Dart type and CBOR type
func convertType(typeName string) (string, string) {
	switch typeName {
	case "bool":
		return "bool", "Bool"
	case "uint8":
		return "int", "Int"
	case "uint16":
		return "int", "Int"
	case "uint32":
		return "int", "Int"
	case "[]byte":
		return "typed.Uint8Buffer", "Array"
	default:
		return "UNKNOWN", "UNKNOWN"
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
    32  .  .  .  .  .  .  .  .  .  .  .  Name: "Len"
    33  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    34  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    35  .  .  .  .  .  .  .  .  .  .  .  .  Name: "Len"
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
