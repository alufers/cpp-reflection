package main

import "fmt"

type VectorType struct {
	InnerType GeneratableType
}

func (vt *VectorType) IdentifierName() string {
	return "VectorOf" + vt.InnerType.IdentifierName()
}
func (vt *VectorType) CppType() string {
	return fmt.Sprintf("std::vector<%v>", vt.InnerType.CppType())
}

func (vt *VectorType) WriteDeclarations(gen *CppGenerator) {
	gen.AddLibraryInclude("vector")
}
func (vt *VectorType) WriteReflection(gen *CppGenerator) {
	fmt.Fprintf(gen.Body, `ReflectType::ofVector(
			/* mine typeId */ ReflectTypeID::%v,
			/* inner type id */  ReflectTypeID::%v,
			/* size */ sizeof(%v),
			__reflectConstruct<%v>,
			__reflectDestruct<%v>
		)`,
		vt.IdentifierName(), vt.InnerType.IdentifierName(), vt.CppType(), vt.CppType(), vt.CppType())
}
