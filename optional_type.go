package main

import (
	"fmt"
	"text/template"
)

type OptionalType struct {
	InnerType GeneratableType
}

func (vt *OptionalType) IdentifierName() string {
	return "OptionalOf" + vt.InnerType.IdentifierName()
}
func (vt *OptionalType) CppType() string {
	return fmt.Sprintf("std::vector<%v>", vt.InnerType.CppType())
}

func (vt *OptionalType) WriteDeclarations(gen *CppGenerator) {
	gen.AddLibraryInclude("vector")
}
func (vt *OptionalType) WriteReflection(gen *CppGenerator) {

	template.Must(template.New("any.cpp").Parse(`
	ReflectType::ofOptional(
		/* mine typeId */ ReflectTypeID::{{ .IdentifierName }},
		/* inner type id */  ReflectTypeID::{{ .InnerType.IdentifierName }},
		/* size */ sizeof({{ .CppType }}),
		OptionalOperations{
			.push_back = __OptionalManipulator<{{ .InnerType.CppType }}>::push_back,
			.at = __OptionalManipulator<{{ .InnerType.CppType }}>::at,
			.size = __OptionalManipulator<{{ .InnerType.CppType }}>::size,
		},
		__reflectConstruct<{{ .CppType }}>,
		__reflectDestruct<{{ .CppType }}>
	)
	
	
	
	`)).Execute(gen.Body, vt)

}
