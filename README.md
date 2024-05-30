
<p align="center"><img src="https://github.com/cscoding21/cscoding/blob/main/assets/csc-banner.png?raw=true" width=728></p>

<p align="center">
    <a href="https://github.com/cscoding21/csgen"><img src="https://img.shields.io/badge/built_with-Go-29BEB0.svg?style=flat-square"></a>&nbsp;
    <a href="https://goreportcard.com/report/github.com/cscoding21/csgen"><img src="https://goreportcard.com/badge/github.com/cscoding21/csgen?style=flat-square"></a>&nbsp;
 <a href="https://pkg.go.dev/mod/github.com/cscoding21/csgen"><img src="https://pkg.go.dev/badge/mod/github.com/cscoding21/csgen"></a>&nbsp;
    <a href="https://github.com/cscoding21/csgen/" alt="Stars">
        <img src="https://img.shields.io/github/stars/cscoding21/csgen?color=0052FF&labelColor=090422" /></a>&nbsp;
    <a href="https://github.com/cscoding21/csgen/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/cscoding21/csgen?color=0052FF&labelColor=090422" /></a>
    <br />
    <a href="https://discord.gg/BjV88Bys" alt="Discord">
        <img src="https://img.shields.io/discord/1196192809120710779" /></a>&nbsp;
    <a href="https://www.youtube.com/@CommonSenseCoding-ge5dn" alt="YouTube">
        <img src="https://img.shields.io/badge/youtube-watch_videos-red.svg?color=0052FF&labelColor=090422&logo=youtube" /></a>&nbsp;
    <a href="https://twitter.com/cscoding21" alt="YouTube">
        <img src="https://img.shields.io/twitter/follow/cscoding21" /></a>&nbsp;
</p>


# CSGen
CSGen is a Golang utility package for simplifying development of code generation tools.  It contains abstractions around the standard library AST packages to make using them more intuitive for developers.

## Get the Package
    go get github.com/cscoding21/csgen

## Core Uses
The primary use-case for the library is to get a list of all the structs in a file including details about each field.

    
    import (
        "fmt"
        "github.com/cscoding21/csgen"
    )
    
    //---get all of the structs within a file.
    structs, err := csgen.GetStructs("test_struct.go")

    if err != nil {
        panic("error loading file")
    }

    for _, s := range structs {
        fmt.Println(s.Name)
        for _, f := range s.Fields {
            fmt.Println("  -- ", f.Name, f.Type)
        }
    } 

The __struct__ object contains information about the struct's definition as well as all fields contained within.

## Struct Properties
| Property | Type | Description |
| --- | --- | --- |
|Name|string|The name of the struct|
|Type|string|The type of the struct|
|FilePath|string|The path to the file containing the struct|
|Package|string|The package the struct is contained within|
|Fields|[]Field|The list of fields contained within the struct|

Additionally, a struct has a __GetField(name string)__ method that will return an individual field object by its name.

## Field Properties
| Property | Type | Description |
| --- | --- | --- |
|Name|string|The name of the field|
|Type|string|The type of the field|
|TagString|string|The tag string for the field|
|IsPrimitive|bool|True if the field is a primitive type|
|IsPointer|bool|True if the field is defined as a pointer to value|
|IsSlice|bool|True if the field is a slice|
|IsPublic|bool|True if the field is available outside the package|

Additionally, a field has a __GetTag(name string)__ method.  This method can be used to extract the value of an individual tag from the tag string.  For example, consider the tag string __`json:"email" csval:"req,email"`__.

    s.GetTag("json")    //---returns "email"
    s.GetTag("caval")   //---returns "req,email"






