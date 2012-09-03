## pkgreflect - A Go preprocessor for package scoped reflection

Go reflection does not support enumerating types, variables and functions of packages.

This preprocessor generates a file called pkgreflect.go in every parsed package directory.
This file contains the follwing maps of exported names to reflection types/values:

	var Types = map[string]reflect.Type{ ... }
	var Functions = map[string]reflect.Value{ ... }
	var Variables = map[string]reflect.Value{ ... }

Command line usage:

	pkgreflect --help
	pkgreflect [-notypes][-nofuncs][-novars][-unexported][-gofile=filename.go] [DIR_NAME]