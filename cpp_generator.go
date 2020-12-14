package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type CppGenerator struct {
	includes []string
	Body     *bytes.Buffer
}

// NewCppGenerator docsy bo ci wywali sie blad xd
func NewCppGenerator() *CppGenerator {
	return &CppGenerator{
		includes: []string{},
		Body:     &bytes.Buffer{},
	}
}

// AddLibraryInclude yes
func (cg *CppGenerator) AddLibraryInclude(name string) {
	resultingLine := fmt.Sprintf("#include <%s>", name)
	for _, a := range cg.includes {
		if a == resultingLine {
			return
		}
	}
	cg.includes = append(cg.includes, resultingLine)
}

// OutputClassField yes
func (cg *CppGenerator) OutputClassField(theType string, name string) {
	fmt.Fprintf(cg.Body, "%v %v;\n", theType, name)
}

// OutputClassTypeID yes
func (cg *CppGenerator) OutputClassTypeID(theID string) {
	fmt.Fprintf(cg.Body, "static constexpr ReflectTypeID _TYPE_ID = ReflectTypeID::%v;\n", theID)
}

// OutputClass yes
func (cg *CppGenerator) OutputClass(name string, cb func()) {
	fmt.Fprintf(cg.Body, "class %v {\npublic:\n", name)
	cb()
	fmt.Fprintf(cg.Body, "};\n\n")
}

// OutputEnumClass
func (cg *CppGenerator) OutputEnumClass(name string, cb func()) {
	fmt.Fprintf(cg.Body, "enum class %v {\n", name)
	cb()
	fmt.Fprintf(cg.Body, "};\n\n")
}

func (cg *CppGenerator) OutputArrayVariable(t string, name string, length int, cb func()) {
	fmt.Fprintf(cg.Body, "%v %v[%d] = {\n", t, name, length)
	cb()
	fmt.Fprintf(cg.Body, "};\n\n")
}

func (cg *CppGenerator) OutputEnumClassField(name string, value string) {
	fmt.Fprintf(cg.Body, "%v", name)
	if value != "" {
		fmt.Fprintf(cg.Body, " = %v", value)
	}
	fmt.Fprintf(cg.Body, ",\n")
}

func (cg *CppGenerator) EscapeCppString(str string) string {
	d, _ := json.Marshal(str)

	return string(d)
}

func (cg *CppGenerator) WriteToWriter(w io.Writer) {
	fmt.Fprintf(w, "// THIS CORNFILE IS GENERATED. DO NOT EDIT! 🌽\n")
	fmt.Fprintf(w, "#ifndef __ALU_CODEGEN\n")
	fmt.Fprintf(w, "#define __ALU_CODEGEN\n")
	for _, a := range cg.includes {
		fmt.Fprintf(w, "%v\n", a)
	}
	io.Copy(w, cg.Body)
	fmt.Fprintf(w, "#endif\n")
}
