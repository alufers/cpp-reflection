package main

import "fmt"

type VectorHelperTypes struct {
}

func (vht *VectorHelperTypes) GenerateVectorOperationsStruct(cg *CppGenerator) {
	fmt.Fprintln(cg.Body, `
		class AnyRef;
		struct VectorOperations {
			void (*push_back)(AnyRef &vec, AnyRef &val);
			AnyRef (*at)(AnyRef &vec, size_t index);
			size_t (*size)(AnyRef &vec);
		};
		
	`)
}

func (vht *VectorHelperTypes) GenerateVectorManipulator(cg *CppGenerator) {
	fmt.Fprintln(cg.Body, `
	template<class T>
	class __VectorManipulator {
		public:
			static void push_back(AnyRef &vec, AnyRef &val) {
				auto theVector = reinterpret_cast<std::vector<T>*>(vec.value.voidptr);
				auto theValue = *reinterpret_cast<T*>(val.value.voidptr);
				theVector->push_back(theValue);
			};
			static AnyRef at(AnyRef &vec, size_t index) {
				auto theVector = reinterpret_cast<std::vector<T>*>(vec.value.voidptr);
				return AnyRef::of<T>(&(*theVector)[index]);
			};
			static size_t size(AnyRef &vec) {
				auto theVector = reinterpret_cast<std::vector<T>*>(vec.value.voidptr);
				return theVector->size();
			};
	};
	`)
}
