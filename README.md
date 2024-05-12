# CSGen
CSGen is a Golang utility package for simplifying development of code generation tools.  It contains abstractions around the standard library AST packages to make using them more intuitive for developers.

## To Install
    go get github.com/cscoding21/csgen

## Core Uses
The primary use-case for the library is to get a list of all the structs in a file including details about each field.

    //---get all of the structs within a file.  The file argument
    structs, err := csgen.GetStructs("test_struct.go")