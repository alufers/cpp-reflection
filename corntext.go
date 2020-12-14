package main

import (
	"fmt"
	"log"
	"strings"

	"os"
	"sort"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// czesc ;) 3===ddd
// jan pawel drugi lubil male dzieci

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

type Corntext struct {
	CG             *CppGenerator
	Request        *plugin.CodeGeneratorRequest  // The input.
	Response       *plugin.CodeGeneratorResponse // The output.
	AllTypes       []GeneratableType
	PrimitiveTypes map[string]GeneratableType

	// types required for reflection
	TypeIDEnum                GeneratableType
	ReflectType               GeneratableType
	ReflectField              GeneratableType
	ReflectEnumValue          GeneratableType
	TypeKindEnum              GeneratableType
	vectorOfReflectFields     GeneratableType
	vectorOfReflectEnumValues GeneratableType
}

func NewCorntext() *Corntext {
	return &Corntext{
		Request:  new(plugin.CodeGeneratorRequest),
		Response: new(plugin.CodeGeneratorResponse),
		CG:       NewCppGenerator(),
		AllTypes: []GeneratableType{},
		PrimitiveTypes: map[string]GeneratableType{
			"int":           &PrimitiveType{identifierName: "Int", cppType: "int"},
			"unsigned int":  &PrimitiveType{identifierName: "UnsignedInt", cppType: "unsigned int"},
			"char":          &PrimitiveType{identifierName: "Char", cppType: "char"},
			"unsigned char": &PrimitiveType{identifierName: "UnsignedChar", cppType: "unsigned char"},
			"double":        &PrimitiveType{identifierName: "Double", cppType: "double"},
			"float":         &PrimitiveType{identifierName: "Float", cppType: "float"},
			"bool":          &PrimitiveType{identifierName: "Bool", cppType: "bool"},
			"std::string":   &PrimitiveType{identifierName: "String", cppType: "std::string"},
			"size_t":        &PrimitiveType{identifierName: "SizeT", cppType: "size_t"},
			"int32_t":       &PrimitiveType{identifierName: "Int32", cppType: "int32_t"},
			"int64_t":       &PrimitiveType{identifierName: "Int64", cppType: "int64_t"},
			"uint32_t":      &PrimitiveType{identifierName: "Uint32", cppType: "uint32_t"},
			"uint64_t":      &PrimitiveType{identifierName: "Uint64", cppType: "uint64_t"},
			"uint8_t":       &PrimitiveType{identifierName: "Uint8", cppType: "uint8_t"},
		},
	}
}

func (c *Corntext) outputTypes() {
	c.CG.AddLibraryInclude("utility")
	c.CG.AddLibraryInclude("vector")
	vht := &VectorHelperTypes{}
	fmt.Fprintln(c.CG.Body, `
		template<class T>
		void __reflectConstruct(void *mem) {
			new(mem) T;
		}
		template<class T>
		void __reflectDestruct(void *obj) {
			((T*) obj)->~T();
		}
		
	`)
	vht.GenerateVectorOperationsStruct(c.CG)
	for _, t := range c.AllTypes {
		if fwd, ok := t.(ForwardDeclarable); ok {
			fmt.Fprintf(c.CG.Body, "%v\n", fwd.ForwardDeclaration())
		}

	}
	c.CG.AddLibraryInclude("type_traits")

	typesToDeclare := make([]GeneratableType, len(c.AllTypes))
	for i := range c.AllTypes {
		typesToDeclare[i] = c.AllTypes[i]
	}
	sort.SliceStable(typesToDeclare, func(i, j int) bool {
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
		t.WriteDeclarations(c.CG)
	}

	primitiveList := make([]GeneratableType, 0, len(c.PrimitiveTypes))
	for _, pt := range c.PrimitiveTypes {
		primitiveList = append(primitiveList, pt)
	}
	GenerateAnyTypes(c.CG, primitiveList, c.AllTypes)
	vht.GenerateVectorManipulator(c.CG)

	c.CG.OutputArrayVariable(c.ReflectType.CppType(), "reflectTypeInfo", len(c.AllTypes), func() {

		for _, t := range c.AllTypes {
			t.WriteReflection(c.CG)
			fmt.Fprintf(c.CG.Body, ",\n")
		}
	})

	GenerateAnyTypesImplementation(c.CG)

	f, _ := os.Create("out/sranie.h")
	defer f.Close()
	c.CG.WriteToWriter(f)
}

func (c *Corntext) buildAllTypes() {
	c.generateReflectionTypes()
	c.AllTypes = append(c.AllTypes,
		c.TypeIDEnum,
		c.ReflectType,
		c.ReflectField,
		c.ReflectEnumValue,
		c.TypeKindEnum,
		c.vectorOfReflectFields,
		c.vectorOfReflectEnumValues,
	)
	for _, t := range c.PrimitiveTypes {
		c.AllTypes = append(c.AllTypes, t)
	}
	c.generateProtobufTypes()
	for i, t := range c.AllTypes {
		c.TypeIDEnum.(*EnumType).Values = append(c.TypeIDEnum.(*EnumType).Values, EnumValue{Name: t.IdentifierName(), Value: fmt.Sprintf("%v", i)})
	}
}

func (c *Corntext) generateProtobufTypes() {
	var pbType2reflection = map[descriptor.FieldDescriptorProto_Type]GeneratableType{
		descriptor.FieldDescriptorProto_TYPE_INT32:  c.PrimitiveTypes["int32_t"],
		descriptor.FieldDescriptorProto_TYPE_SINT32: c.PrimitiveTypes["int32_t"],
		descriptor.FieldDescriptorProto_TYPE_SINT64: c.PrimitiveTypes["int64_t"],
		descriptor.FieldDescriptorProto_TYPE_INT64:  c.PrimitiveTypes["int64_t"],
		descriptor.FieldDescriptorProto_TYPE_UINT32: c.PrimitiveTypes["uint32_t"],
		descriptor.FieldDescriptorProto_TYPE_UINT64: c.PrimitiveTypes["uint64_t"],
		descriptor.FieldDescriptorProto_TYPE_BOOL:   c.PrimitiveTypes["bool"],
		descriptor.FieldDescriptorProto_TYPE_BYTES:  c.vectorOf(c.PrimitiveTypes["uint8_t"]),
		descriptor.FieldDescriptorProto_TYPE_STRING: c.PrimitiveTypes["std::string"],
	}
	for _, f := range c.Request.ProtoFile {
		log.Printf("Doing file %v", *f.Name)
		typeMappings := map[string]GeneratableType{}
		for _, e := range f.EnumType {
			values := make([]EnumValue, 0, len(e.Value))
			for _, v := range e.Value {
				values = append(values, EnumValue{
					Name:  *v.Name,
					Value: fmt.Sprint(*v.Number),
				})
			}
			et := &EnumType{
				Name:   *e.Name,
				Values: values,
			}
			typeMappings[*e.Name] = et
			c.AllTypes = append(c.AllTypes, et)
		}
		for _, m := range f.MessageType {

			fields := []ClassField{}
			for _, f := range m.Field {
				isMessage := *f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
				isEnum := *f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM
				var fieldType GeneratableType
				if isMessage || isEnum {
					fqn := strings.Split(*f.TypeName, ".")
					className := fqn[1]

					fieldType = typeMappings[className]
				} else {
					primitiveType, ok := pbType2reflection[*f.Type]
					if !ok {
						log.Fatal("unsupported proto type", (*f.Type).String())
					}
					fieldType = primitiveType
					// log.Printf("%#v == %v", primitiveType, *f.Type)

				}
				if f.Label != nil && *f.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED {
					fieldType = c.vectorOf(fieldType)
				}
				fields = append(fields, ClassField{
					Name: *f.Name,
					Type: fieldType,
				})
			}
			ct := &ClassType{
				Name:   *m.Name,
				Fields: fields,
			}
			typeMappings[*m.Name] = ct
			c.AllTypes = append(c.AllTypes, ct)
		}
	}
}

func (c *Corntext) vectorOf(inner GeneratableType) (ret GeneratableType) {
	for _, t := range c.AllTypes {
		if v, ok := t.(*VectorType); ok {
			if v.InnerType == inner {
				ret = v.InnerType

			}
		}
	}
	if ret == nil {
		ret = &VectorType{
			InnerType: inner,
		}
		c.AllTypes = append(c.AllTypes, ret)
	}
	return
}

func (c *Corntext) generateReflectionTypes() {
	c.TypeIDEnum = &EnumType{
		Name:   "ReflectTypeID",
		Values: []EnumValue{},
	}
	c.TypeKindEnum = &EnumType{
		Name: "ReflectTypeKind",
		Values: []EnumValue{
			{Name: "Primitive", Value: "0"},
			{Name: "Enum", Value: "1"},
			{Name: "Class", Value: "2"},
			{Name: "Vector", Value: "3"},
		},
	}
	c.ReflectField = &ClassType{
		Name: "ReflectField",
		Fields: []ClassField{
			{"typeID", c.TypeIDEnum},
			{"name", c.PrimitiveTypes["std::string"]},
			{"offset", c.PrimitiveTypes["size_t"]},
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
	c.ReflectEnumValue = &ClassType{
		Name: "ReflectEnumValue",
		Fields: []ClassField{

			{"name", c.PrimitiveTypes["std::string"]},
			{"value", c.PrimitiveTypes["int"]},
		},
		AdditionalCode: `
			ReflectEnumValue(){};
			ReflectEnumValue( std::string name, int value) {
				this->name = name;
				this->value = value;
			}
		`,
	}
	c.vectorOfReflectFields = &VectorType{
		InnerType: c.ReflectField,
	}
	c.vectorOfReflectEnumValues = &VectorType{
		InnerType: c.ReflectEnumValue,
	}
	c.ReflectType = &ClassType{
		Name: "ReflectType",
		Fields: []ClassField{
			{"typeID", c.TypeIDEnum},
			{"name", c.PrimitiveTypes["std::string"]},
			{"kind", c.TypeKindEnum},
			{"size", c.PrimitiveTypes["size_t"]},
			{"innerType", c.TypeIDEnum},
			{"fields", c.vectorOfReflectFields},
			{"enumValues", c.vectorOfReflectEnumValues},
		},
		AdditionalCode: `
		void (*_Construct)(void *mem);
		void (*_Destruct)(void *obj);
		VectorOperations vectorOps;
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
			VectorOperations vectorOps,
			void (*_Construct)(void *mem), void (*_Destruct)(void *obj)) {
			ReflectType t;
			t.kind = ReflectTypeKind::Vector;
			t.typeID = id;
			t.innerType = innerType;
			t.size = size;
			t._Construct = _Construct;
			t._Destruct = _Destruct;
			t.vectorOps = vectorOps;
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
}
