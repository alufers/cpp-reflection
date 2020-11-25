package main

import (
	"text/template"
)

func GenerateAnyTypes(gen *CppGenerator, primitiveTypes []GeneratableType, allTypes []GeneratableType) {
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
				

				ReflectType *reflectType() {
					return &reflectTypeInfo[static_cast<int>(this->typeID)];
				}
				AnyRef getField(int i) {
					auto info = this->reflectType();
					if(info->kind != ReflectTypeKind::Class) {
						throw "not a class";
					}
					return AnyRef(info->fields[i].typeID, this->value.voidptr + info->fields[i].offset);
				}
				template <typename T>
				static AnyRef of(T *obj)
				{
					ReflectTypeID typeID = T::_TYPE_ID;
					AnyRef a;
					a.typeID = typeID;
					a.value.voidptr = (void*) obj;
					return a;
				}
			
				union ReflectedTypes {
					void *voidptr;
					{{range .AllTypes}}{{.CppType}} *u_{{.IdentifierName}};
					{{end}}
				} value;
				private:
		
	};

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
	
	`)).Execute(gen.Body, map[string]interface{}{
		"PrimitiveTypes": primitiveTypes,
		"AllTypes":       allTypes,
	})
}
