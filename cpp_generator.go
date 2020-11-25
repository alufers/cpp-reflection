package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type CppGenerator struct {
	includes []string
	Body     *bytes.Buffer
}

func NewCppGenerator() *CppGenerator {
	return &CppGenerator{
		includes: []string{},
		Body:     &bytes.Buffer{},
	}
}

func (cg *CppGenerator) AddLibraryInclude(name string) {
	resultingLine := fmt.Sprintf("#include <%s>", name)
	for _, a := range cg.includes {
		if a == resultingLine {
			return
		}
	}
	cg.includes = append(cg.includes, resultingLine)
}

func (cg *CppGenerator) OutputClassField(theType string, name string) {
	fmt.Fprintf(cg.Body, "%v %v;\n", theType, name)
}

func (cg *CppGenerator) OutputClassTypeID(theID string) {
	fmt.Fprintf(cg.Body, "static constexpr ReflectTypeID _TYPE_ID = ReflectTypeID::%v;\n", theID)
}

func (cg *CppGenerator) OutputClass(name string, cb func()) {
	fmt.Fprintf(cg.Body, "class %v {\npublic:\n", name)
	cb()
	fmt.Fprintf(cg.Body, "};\n\n")
}

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
