package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func (cg *CppGenerator) WriteToWriter(w io.Writer) {
	fmt.Fprintf(w, "// THIS FILE IS GENERATED. DO NOT EDIT!\n")
	fmt.Fprintf(w, "#ifndef __ALU_CODEGEN\n")
	fmt.Fprintf(w, "#define __ALU_CODEGEN\n")
	for _, a := range cg.includes {
		fmt.Fprintf(w, "%v\n", a)
	}
	io.Copy(w, cg.Body)
	fmt.Fprintf(w, "#endif\n")
}

type GeneratableType interface {
	IdentifierName() string
	CppType() string
	WriteDeclarations(gen *CppGenerator)
	WriteReflection(gen *CppGenerator)
}

type ForwardDeclarable interface {
	ForwardDeclaration() string
}

type DelcarationOrderable interface {
	DeclarationOrder() int
}

///////////////////////////////// DDD
var UserDefinedTypes = []GeneratableType{
	&ClassType{
		Name: "Foo",
		Fields: []ClassField{
			{"alpha", PrimitiveTypes["int"]},
			{"beta", PrimitiveTypes["bool"]},
			{"gamma", PrimitiveTypes["bool"]},
		},
	},
	&ClassType{
		Name: "Foo2",
		Fields: []ClassField{
			{"alpha", PrimitiveTypes["int"]},
			{"beta", PrimitiveTypes["bool"]},
			{"gamma", PrimitiveTypes["bool"]},
		},
	},
}

//////////////////////

var PrimitiveTypes = map[string]GeneratableType{
	"int":           &PrimitiveType{identifierName: "Int", cppType: "int"},
	"unsigned int":  &PrimitiveType{identifierName: "UnsignedInt", cppType: "unsigned int"},
	"char":          &PrimitiveType{identifierName: "Char", cppType: "char"},
	"unsigned char": &PrimitiveType{identifierName: "UnsignedChar", cppType: "unsigned char"},
	"double":        &PrimitiveType{identifierName: "Double", cppType: "double"},
	"float":         &PrimitiveType{identifierName: "Float", cppType: "float"},
	"bool":          &PrimitiveType{identifierName: "Bool", cppType: "bool"},
	"std::string":   &PrimitiveType{identifierName: "String", cppType: "std::string"},
	"size_t":        &PrimitiveType{identifierName: "SizeT", cppType: "size_t"},
}

func main() {
	///////////////////////////////// DDD
	UserDefinedTypes = append(UserDefinedTypes, &ClassType{
		Name: "Bar",
		Fields: []ClassField{
			{"fooOne", UserDefinedTypes[0]},
			{"fooTwo", UserDefinedTypes[0]},
		},
	})

	//////////////////////

	TypeIDEnum := &EnumType{
		Name:   "ReflectTypeID",
		Values: []EnumValue{},
	}
	TypeKindEnum := &EnumType{
		Name: "ReflectTypeKind",
		Values: []EnumValue{
			{Name: "Primitive", Value: "0"},
			{Name: "Enum", Value: "1"},
			{Name: "Class", Value: "2"},
			{Name: "Vector", Value: "3"},
		},
	}
	ReflectField := &ClassType{
		Name: "ReflectField",
		Fields: []ClassField{
			{"typeID", TypeIDEnum},
			{"name", PrimitiveTypes["std::string"]},
			{"offset", PrimitiveTypes["size_t"]},
		},
		AdditionalCode: `
			ReflectField() {};
			ReflectField(ReflectTypeID typeID, std::string name, size_t offset) {
				this->typeID = typeID;
				this->name = name;
				this->offset = offset;
			}
		`,
	}
	ReflectEnumValue := &ClassType{
		Name: "ReflectEnumValue",
		Fields: []ClassField{

			{"name", PrimitiveTypes["std::string"]},
			{"value", PrimitiveTypes["int"]},
		},
		AdditionalCode: `
			ReflectEnumValue(){};
			ReflectEnumValue( std::string name, int value) {
				this->name = name;
				this->value = value;
			}
		`,
	}
	vectorOfReflectFields := &VectorType{
		InnerType: ReflectField,
	}
	vectorOfReflectEnumValues := &VectorType{
		InnerType: ReflectEnumValue,
	}
	ReflectType := &ClassType{
		Name: "ReflectType",
		Fields: []ClassField{
			{"typeID", TypeIDEnum},
			{"name", PrimitiveTypes["std::string"]},
			{"kind", TypeKindEnum},
			{"size", PrimitiveTypes["size_t"]},
			{"innerType", TypeIDEnum},
			{"fields", vectorOfReflectFields},
			{"enumValues", vectorOfReflectEnumValues},
		},
		AdditionalCode: `
		void (*_Construct)(void *mem);
		void (*_Destruct)(void *obj);
		// VectorOperations vectorOps;
		static ReflectType ofPrimitive(ReflectTypeID id, std::string name, size_t size) {
			ReflectType t;
			t.kind = ReflectTypeKind::Primitive;
			t.typeID = id;
			t.name = name;
			t.size = size;
			return t;
		}
		static ReflectType ofEnum(ReflectTypeID id, std::string name, std::vector<ReflectEnumValue> enumValues, size_t size) {
			ReflectType t;
			t.kind = ReflectTypeKind::Enum;
			t.typeID = id;
			t.name = name;
			t.size = size;
			t.enumValues = enumValues;
			return t;
		}
		static ReflectType ofVector(ReflectTypeID id, ReflectTypeID innerType, size_t size, 
			// VectorOperations vectorOps,
			void (*_Construct)(void *mem), void (*_Destruct)(void *obj)) {
			ReflectType t;
			t.kind = ReflectTypeKind::Vector;
			t.typeID = id;
			t.innerType = innerType;
			t.size = size;
			t._Construct = _Construct;
			t._Destruct = _Destruct;
			// this->vectorOps = vectorOps;
			return t;
		}
		static ReflectType ofClass(ReflectTypeID id, std::string name, std::vector<ReflectField> fields, size_t size, void (*_Construct)(void *mem), void (*_Destruct)(void *obj)) {
			ReflectType t;
			t.kind = ReflectTypeKind::Class;
			t.name = name;
			t.typeID = id;
			t.size = size;
			t.fields = std::move(fields);
			t._Construct = _Construct;
			t._Destruct = _Destruct;
			return t;
		}
		
		`,
	}
	AllTypes := []GeneratableType{
		TypeIDEnum,
		ReflectType,
		TypeKindEnum,
		vectorOfReflectFields,
		ReflectField,
		vectorOfReflectEnumValues,
		ReflectEnumValue,
	}
	for _, t := range PrimitiveTypes {
		AllTypes = append(AllTypes, t)
	}
	for _, t := range UserDefinedTypes {
		AllTypes = append(AllTypes, t)
	}

	for i, t := range AllTypes {
		TypeIDEnum.Values = append(TypeIDEnum.Values, EnumValue{Name: t.IdentifierName(), Value: fmt.Sprintf("%d", i)})
	}
	cg := NewCppGenerator()
	cg.AddLibraryInclude("utility")
	cg.AddLibraryInclude("vector")
	fmt.Fprintln(cg.Body, `
		template<class T>
		void __reflectConstruct(void *mem) {
			new(mem) T;
		}
		template<class T>
		void __reflectDestruct(void *obj) {
			((T*) obj)->~T();
		}
		class AnyRef;
		struct VectorOperations {
			void (*push_back)(AnyRef &vec, AnyRef &val);
			AnyRef (*at)(AnyRef &vec, size_t index);
			size_t (*size)(AnyRef vec, size_t index);
		};
		
	`)
	for _, t := range AllTypes {
		if fwd, ok := t.(ForwardDeclarable); ok {
			fmt.Fprintf(cg.Body, "%v\n", fwd.ForwardDeclaration())
		}

	}
	cg.AddLibraryInclude("type_traits")

	typesToDeclare := make([]GeneratableType, len(AllTypes))
	for i := range AllTypes {
		typesToDeclare[i] = AllTypes[i]
	}
	sort.Slice(typesToDeclare, func(i, j int) bool {
		var ival, jval int
		if orderable, ok := typesToDeclare[i].(DelcarationOrderable); ok {
			ival = orderable.DeclarationOrder()
		}
		if orderable, ok := typesToDeclare[j].(DelcarationOrderable); ok {
			jval = orderable.DeclarationOrder()
		}
		return ival < jval
	})
	for _, t := range typesToDeclare {
		t.WriteDeclarations(cg)
	}

	cg.OutputArrayVariable(ReflectType.CppType(), "reflectTypeInfo", len(AllTypes), func() {
		for _, t := range AllTypes {
			t.WriteReflection(cg)
			fmt.Fprintf(cg.Body, ",\n")
		}
	})

	primitiveList := make([]GeneratableType, 0, len(PrimitiveTypes))
	for _, pt := range PrimitiveTypes {
		primitiveList = append(primitiveList, pt)
	}
	GenerateAnyTypes(cg, primitiveList, AllTypes)

	fmt.Fprintln(cg.Body, `
	template<class T>
	class VectorManipulator {
		public:
			static void push_back(AnyRef &vec, AnyRef &val) {
				auto theVector = reinterpret_cast<std::vector<T>*>(vec.value.voidptr);
				auto theValue = reinterpret_cast<T*>(val.value.voidptr);
				theVector.push_back(theValue);
			};
			static AnyRef at(AnyRef &vec, size_t index) {
				auto theVector = reinterpret_cast<std::vector<T>*>(vec.value.voidptr);
				return theVector.at(index);
			};
			static size_t size(AnyRef &vec) {
				auto theVector = reinterpret_cast<std::vector<T>*>(vec.value.voidptr);
				return theVector.size();
			};
	};
	`)

	f, _ := os.Create("out/sranie.h")
	defer f.Close()
	cg.WriteToWriter(f)
}
