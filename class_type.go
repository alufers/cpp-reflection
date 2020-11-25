package main

import "fmt"

type ClassField struct {
	Name string
	Type GeneratableType
}

type ClassType struct {
	Name           string
	Fields         []ClassField
	AdditionalCode string
}

func (ct *ClassType) IdentifierName() string {
	return "Class" + ct.Name
}

func (ct *ClassType) CppType() string {
	return ct.Name
}

func (et *ClassType) DeclarationOrder() int {
	return -10
}

func (ct *ClassType) ForwardDeclaration() string {
	return fmt.Sprintf("class %v;", ct.Name)
}

func (ct *ClassType) WriteDeclarations(gen *CppGenerator) {
	gen.OutputClass(ct.Name, func() {
		for _, t := range ct.Fields {
			gen.OutputClassField(t.Type.CppType(), t.Name)
		}
		gen.OutputClassTypeID(ct.IdentifierName())
		fmt.Fprint(gen.Body, ct.AdditionalCode)
	})
}
func (ct *ClassType) WriteReflection(gen *CppGenerator) {
	gen.AddLibraryInclude("vector")
	fieldsContents := ""
	for _, f := range ct.Fields {
		fieldsContents += fmt.Sprintf("ReflectField( /* typeID */ ReflectTypeID::%v, /* name */ %v, /* offset */ offsetof(%v, %v)),\n",
			f.Type.IdentifierName(), gen.EscapeCppString(f.Name), ct.CppType(), f.Name)
	}
	fmt.Fprintf(gen.Body, `ReflectType::ofClass(
	/* mine type id */ ReflectTypeID::%v, 
	/* name */ %v, 
	/* fields */ std::move(std::vector<ReflectField>{%v}), 
	/* size */ sizeof(%v), 
	__reflectConstruct<%v>,
	__reflectDestruct<%v>)`,
		ct.IdentifierName(), gen.EscapeCppString(ct.CppType()), fieldsContents, ct.CppType(), ct.CppType(), ct.CppType())
}
