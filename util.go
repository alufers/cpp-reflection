package main

func AddIncludeForType(t GeneratableType, gen *CppGenerator) {
	switch v := t.(type) {
	case *ClassType:
		gen.AddLocalInclude(v.Name + ".h")
	case *EnumType:
		gen.AddLocalInclude(v.Name + ".h")
	case *VectorType:
		gen.AddLibraryInclude("vector")
		AddIncludeForType(v.InnerType, gen)
	case *OptionalType:
		gen.AddLibraryInclude("optional")
		AddIncludeForType(v.InnerType, gen)
	}
}
