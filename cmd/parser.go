package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	accessorTag  = "vstruct"
	ignoreTag    = "-"
	tagKeyRepeat = "repeat"
	tagKeyRefer  = "refer"
	tagKeyCrc16  = "crc16"
	tagKeyAccess = "access"
	tagKeySize   = "size"
	tagKeyPoint  = "point"
)

const (
	tagSep         = ","
	tagKeyValueSep = ":"
)

type Package struct {
	Dir   string
	Name  string
	Files []*File
}

type File struct {
	File    *ast.File
	Structs []*Struct
}

type Struct struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name     string
	DataType string
	Tag      *Tag
}

type Tag struct {
	Repeat    *string
	Refer     *string
	Crc16     *string
	Access    *string
	FiledSize *string
	Point     *string
}

// ParsePackage parses the specified directory's package.
func ParsePackage(dir string) (*Package, error) {
	const mode = packages.NeedName | packages.NeedFiles |
		packages.NeedImports | packages.NeedTypes | packages.NeedSyntax

	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	cfg := &packages.Config{
		Mode:  mode,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(pkgs))
	}

	return &Package{Dir: dir, Name: pkgs[0].Name, Files: parseFiles(pkgs[0])}, nil
}

func parseFiles(pkg *packages.Package) []*File {
	files := make([]*File, len(pkg.Syntax))
	for i := range pkg.Syntax {
		files[i] = &File{
			File:    pkg.Syntax[i],
			Structs: parseStructs(pkg.Fset, pkg.Syntax[i]),
		}
	}

	return files
}

func parseStructs(fileSet *token.FileSet, file *ast.File) []*Struct {
	structs := make([]*Struct, 0)

	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structs = append(structs, &Struct{
			Name:   ts.Name.Name,
			Fields: parseFields(s, fileSet),
		})

		return false
	})

	return structs
}

func parseFields(st *ast.StructType, fileSet *token.FileSet) []*Field {
	fields := make([]*Field, 0)
	for _, field := range st.Fields.List {
		name := field.Names[0].Name
		buf := new(bytes.Buffer)
		printer.Fprint(buf, fileSet, field.Type)
		sf := &Field{
			Name:     name,
			DataType: buf.String(),
			Tag:      parseTag(field.Tag),
		}
		fields = append(fields, sf)
	}

	return fields
}

func parseTag(tag *ast.BasicLit) *Tag {
	if tag == nil {
		return nil
	}

	tagStr, ok := reflect.StructTag(strings.Trim(tag.Value, "`")).Lookup(accessorTag)
	if !ok {
		return nil
	}

	var repeat, refer, crc16, access, filedSize, point *string

	tags := strings.Split(tagStr, tagSep)
	for _, tag := range tags {
		keyValue := strings.Split(tag, tagKeyValueSep)

		var value string
		if len(keyValue) == 2 {
			if v := strings.TrimSpace(keyValue[1]); v != ignoreTag {
				value = v
			}
		}
		switch strings.TrimSpace(keyValue[0]) {
		case tagKeyRepeat:
			repeat = &value
		case tagKeyRefer:
			refer = &value
		case tagKeyCrc16:
			crc16 = &value
		case tagKeyAccess:
			access = &value
		case tagKeySize:
			filedSize = &value
		case tagKeyPoint:
			point = &value
		}
	}

	return &Tag{Repeat: repeat, Refer: refer, Crc16: crc16,
		Access: access, FiledSize: filedSize,
		Point: point,
	}
}
