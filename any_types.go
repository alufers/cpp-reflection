package main

import (
	"fmt"
	"text/template"
)

func GenerateAnyTypes(gen *CppGenerator, primitiveTypes []GeneratableType, allTypes []GeneratableType) {

	exceptionalTypes := []GeneratableType{}

	for _, t := range allTypes {
		_, isPrimitve := t.(*PrimitiveType)
		_, isEnum := t.(*EnumType)
		if isPrimitve || isEnum {
			exceptionalTypes = append(exceptionalTypes, t)
		}
	}

	template.Must(template.New("any.cpp").Parse(`

	
	class AnyRef {
		public:
				ReflectTypeID typeID;
				AnyRef() {};
				AnyRef(ReflectTypeID typeID, void *obj) {
					this->typeID = typeID;
					this->value.voidptr = obj;
				}
				template<typename T>
				T *as() {
					// if(T::_TYPE_ID != this->typeID) {
					// 	throw "invalid as call";
					// }
					return (T*) this->value.voidptr;
				}

				template<typename T>
				bool is() {
					{{range .PrimitiveTypes}}if constexpr(std::is_same<T, {{.CppType}}>::value) {
						return ReflectTypeID::{{.IdentifierName}} == this->typeID;
					} else 
					{{end}} {
						return T::_TYPE_ID == this->typeID;
					}
				}
				

				ReflectType *reflectType();
				AnyRef getField(int i);
				template <typename T>
				static AnyRef of(T *obj)
				{
					ReflectTypeID typeID;
					{{range .PrimitiveTypes}}if constexpr(std::is_same<T, {{.CppType}}>::value) {
						typeID = ReflectTypeID::{{.IdentifierName}};
					} else 
					{{end}} {
						typeID = T::_TYPE_ID;
					}
					AnyRef a;
					a.typeID = typeID;
					a.value.voidptr = (void*) obj;
					return a;
				}
			
				union ReflectedTypes {
					void *voidptr;
					{{range .allTypes}}{{.CppType}} *u_{{.IdentifierName}};
					{{end}}
				} value;
				private:
		
	};
	
	
	`)).Execute(gen.Body, map[string]interface{}{
		"PrimitiveTypes": exceptionalTypes,
		"allTypes":       allTypes,
	})
}

func GenerateAnyTypesImplementation(gen *CppGenerator) {
	fmt.Fprintf(gen.Body, `
	ReflectType *AnyRef::reflectType() {
		return &reflectTypeInfo[static_cast<int>(this->typeID)];
	}
	AnyRef AnyRef::getField(int i) {
		auto info = this->reflectType();
		if(info->kind != ReflectTypeKind::Class) {
			throw "not a class";
		}
		return AnyRef(info->fields[i].typeID, this->value.voidptr + info->fields[i].offset);
	}

	class UniqueAny: public AnyRef {
		public:
			UniqueAny() {
				this->value.voidptr = nullptr;
			};
			UniqueAny(ReflectTypeID typeID) {
				this->typeID = typeID;
				auto typeInfo = &reflectTypeInfo[static_cast<int>(typeID)];
				AnyRef a;
				this->value.voidptr = new unsigned char[typeInfo->size];
				typeInfo->_Construct(this->value.voidptr);
			};
			~UniqueAny() {
				auto typeInfo = &reflectTypeInfo[static_cast<int>(typeID)];
				typeInfo->_Destruct(this->value.voidptr);
				delete this->value.voidptr;
			};
	};

	class AnyVectorRef {
		public:
			AnyRef ref;
			AnyVectorRef(AnyRef r): ref(r) {}
			void push_back(AnyRef &v) {
				auto typeInfo = &reflectTypeInfo[static_cast<int>(this->ref.typeID)];
				typeInfo->vectorOps.push_back(ref, v);
			}
			size_t size() {
				auto typeInfo = &reflectTypeInfo[static_cast<int>(this->ref.typeID)];
				return typeInfo->vectorOps.size(ref);
			}

			AnyRef at(size_t index) {
				auto typeInfo = &reflectTypeInfo[static_cast<int>(this->ref.typeID)];
				return typeInfo->vectorOps.at(ref, index);
			}
	};

	`)
}
