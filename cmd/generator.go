package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type generator struct {
	receiver string
	hascrc16 bool
	headBuf  *bytes.Buffer
	buf      *bytes.Buffer
}

func (g *generator) writeHeader(pkgName string) {
	fmt.Fprintf(g.headBuf, "// Code generated by vstruct; DO NOT EDIT.\n")
	fmt.Fprintf(g.headBuf, "\n")
	fmt.Fprintf(g.headBuf, "package %s\n", pkgName)
	fmt.Fprintf(g.headBuf, "\n")

	fmt.Fprintf(g.headBuf, "import (\n")
	fmt.Fprintf(g.headBuf, "\t\"bytes\"\n")
	fmt.Fprintf(g.headBuf, "\t\"encoding/binary\"\n")
	if g.hascrc16 {
		fmt.Fprintf(g.headBuf, "\t\"github.com/yumm007/gohash\"\n")
	}
	fmt.Fprintf(g.headBuf, ")\n\n")
}

func (g *generator) filedEncodeGenerate(f *Field, r string, pre string) {
	defer func() {
		fmt.Fprintf(g.buf, "%s\t\treturn err\n", pre)
		fmt.Fprintf(g.buf, "%s\t}\n", pre)
	}()

	if f.Tag != nil {
		if f.Tag.Refer != nil {
			if f.Tag.Repeat != nil {
				fmt.Fprintf(g.buf, "%s\tif err := %s.%s[i].encodeToBuffer(buf); err != nil {\n", pre, r, f.Name)
			} else {
				fmt.Fprintf(g.buf, "%s\tif err := %s.%s.encodeToBuffer(buf); err != nil {\n", pre, r, f.Name)
			}
			return
		} else if f.Tag.Repeat != nil {
			fmt.Fprintf(g.buf, "%s\tif err := binary.Write(buf, binary.LittleEndian, &%s.%s[i]); err != nil {\n", pre, r, f.Name)
			return
		}
	}

	if f.Tag != nil && f.Tag.Crc16 != nil {
		fmt.Fprintf(g.buf, "%s\tif err := binary.Write(buf, binary.LittleEndian, gohash.Crc16ccitt(buf.Bytes())); err != nil {\n", pre)
		return
	}

	fmt.Fprintf(g.buf, "%s\tif err := binary.Write(buf, binary.LittleEndian, &%s.%s); err != nil {\n", pre, r, f.Name)
}

func (g *generator) structEncodeGenerate(st *Struct) {
	r := strings.ToLower(string(st.Name[0]))
	if len(g.receiver) > 0 {
		r = g.receiver
	}

	fmt.Fprintf(g.buf, "\nfunc (%s *%s)encodeToBuffer(buf *bytes.Buffer) error {\n", r, st.Name)

	for _, f := range st.Fields {
		if f.Tag != nil && f.Tag.Crc16 != nil {
			g.hascrc16 = true
		}
		if f.Tag != nil {
			if f.Tag.Repeat != nil {
				acc := "int"
				if f.Tag.Access != nil {
					acc = *f.Tag.Access
				}
				fmt.Fprintf(g.buf, "\tfor i := 0; i < int(%s(%s.%s)); i++ {\n", acc, r, *f.Tag.Repeat)
				g.filedEncodeGenerate(f, r, "\t")
				fmt.Fprintf(g.buf, "\t}\n")
				continue
			}
		}

		g.filedEncodeGenerate(f, r, "")
	}

	fmt.Fprintf(g.buf, "\n\treturn nil\n")
	fmt.Fprintf(g.buf, "}\n")

	fmt.Fprintf(g.buf, "\nfunc (%s *%s)Encode(buf *bytes.Buffer) ([]byte, error) {\n", r, st.Name)
	fmt.Fprintf(g.buf, "\tif buf == nil {\n")
	fmt.Fprintf(g.buf, "\t\tbuf = new(bytes.Buffer)\n")
	fmt.Fprintf(g.buf, "\t} else {\n")
	fmt.Fprintf(g.buf, "\t\tbuf.Reset()\n")
	fmt.Fprintf(g.buf, "\t}\n")
	fmt.Fprintf(g.buf, "\tif err := %s.encodeToBuffer(buf); err != nil {\n", r)
	fmt.Fprintf(g.buf, "\t\treturn nil, err\n")
	fmt.Fprintf(g.buf, "\t}\n")
	fmt.Fprintf(g.buf, "\treturn buf.Bytes(), nil\n")
	fmt.Fprintf(g.buf, "}\n")
}

func (g *generator) filedDecodeGenerate(f *Field, r string, pre string) {
	defer func() {
		fmt.Fprintf(g.buf, "%s\t\treturn err\n", pre)
		fmt.Fprintf(g.buf, "%s\t}\n", pre)
	}()

	if f.Tag != nil {
		if f.Tag.Refer != nil {
			if f.Tag.Repeat != nil {
				fmt.Fprintf(g.buf, "%s\tif err := ele.decodeFromBuffer(buf); err != nil {\n", pre)
			} else {
				fmt.Fprintf(g.buf, "%s\tif err := %s.%s.decodeFromBuffer(buf); err != nil {\n", pre, r, f.Name)
			}

			return
		} else if f.Tag.Repeat != nil {
			fmt.Fprintf(g.buf, "%s\tif err := binary.Read(buf, binary.LittleEndian, &ele); err != nil {\n", pre)
			return
		}
	}

	fmt.Fprintf(g.buf, "%s\tif err := binary.Read(buf, binary.LittleEndian, &%s.%s); err != nil {\n", pre, r, f.Name)
}

func (g *generator) structDecodeGenerate(st *Struct) {
	r := strings.ToLower(string(st.Name[0]))
	if len(g.receiver) > 0 {
		r = g.receiver
	}

	fmt.Fprintf(g.buf, "\nfunc (%s *%s)decodeFromBuffer(buf *bytes.Buffer) error {\n", r, st.Name)

	for _, f := range st.Fields {
		if f.Tag != nil && f.Tag.Crc16 != nil {
			g.hascrc16 = true
		}

		if f.Tag != nil {
			if f.Tag.Repeat != nil {
				acc := "int"
				if f.Tag.Access != nil {
					acc = *f.Tag.Access
				}
				fmt.Fprintf(g.buf, "\n\tele_len := int(%s(%s.%s))\n", acc, r, *f.Tag.Repeat)
				fmt.Fprintf(g.buf, "\t%s.%s = make(%s, 0, ele_len)\n", r, f.Name, f.DataType)
				fmt.Fprintf(g.buf, "\tfor i := 0; i < ele_len; i++ {\n")
				fmt.Fprintf(g.buf, "\t\tvar ele %s\n", f.DataType[2:])
				g.filedDecodeGenerate(f, r, "\t")
				fmt.Fprintf(g.buf, "\t\t%s.%s = append(%s.%s, ele)\n", r, f.Name, r, f.Name)
				fmt.Fprintf(g.buf, "\t}\n")
				continue
			}
		}
		g.filedDecodeGenerate(f, r, "")
	}

	fmt.Fprintf(g.buf, "\n\treturn nil\n")
	fmt.Fprintf(g.buf, "}\n")

	fmt.Fprintf(g.buf, "\nfunc (%s *%s)Decode(payload []byte) error {\n", r, st.Name)
	fmt.Fprintf(g.buf, "\tbuf := bytes.NewBuffer(payload)\n")
	fmt.Fprintf(g.buf, "\treturn %s.decodeFromBuffer(buf)\n", r)
	fmt.Fprintf(g.buf, "}\n")
}

// Generate a file and accessor methods in it.
func Generate(pkg *Package, typeList []string, output string, receiverName string) error {
	g := generator{buf: new(bytes.Buffer), headBuf: new(bytes.Buffer), receiver: receiverName}

	for _, file := range pkg.Files {
		for _, st := range file.Structs {

			for _, typeName := range typeList {
				if st.Name == typeName {
					g.structEncodeGenerate(st)
					g.structDecodeGenerate(st)
					break
				}
			}

		}
	}

	// 后写，因为涉及到动态import
	g.writeHeader(pkg.Name)
	_, _ = g.headBuf.ReadFrom(g.buf)

	pkgName := pkg.Name
	if len(typeList) == 1 {
		pkgName = typeList[0]
	}

	outputFile := g.outputFile(output, pkg.Name, pkgName, pkg.Dir)

	return os.WriteFile(outputFile, g.headBuf.Bytes(), 0644)
}

func (g *generator) outputFile(output, pkgName, typeName, dir string) string {
	if output == "" {
		// Use snake_case name of type as output file if output file is not specified.
		// type TestStruct will be test_struct_accessor.go
		var firstCapMatcher = regexp.MustCompile("(.)([A-Z][a-z]+)")
		var articleCapMatcher = regexp.MustCompile("([a-z0-9])([A-Z])")

		name := firstCapMatcher.ReplaceAllString(typeName, "${1}_${2}")
		name = articleCapMatcher.ReplaceAllString(name, "${1}_${2}")
		output = strings.ToLower(fmt.Sprintf("autogen_%s_vs.go", name))
	}

	return filepath.Join(dir, output)
}
