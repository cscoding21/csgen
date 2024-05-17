
<p align="center"><img src="https://github.com/cscoding21/cscoding/blob/main/assets/csc-banner.png?raw=true" width=728></p>

<p align="center">
    <a href="https://github.com/cscoding21/csgen/" alt="Stars">
        <img src="https://img.shields.io/github/stars/cscoding21/csgen?color=0052FF&labelColor=090422" /></a>
    <a href="https://github.com/cscoding21/csgen/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/cscoding21/csgen?color=0052FF&labelColor=090422" /></a>
    <a href="https://discord.gg/BjV88Bys" alt="Discord">
        <img src="https://img.shields.io/discord/1196192809120710779" /></a>
    <a href="https://www.youtube.com/@CommonSenseCoding-ge5dn" alt="YouTube">
        <img src="https://img.shields.io/badge/youtube-watch_videos-red.svg?color=0052FF&labelColor=090422&logo=youtube" /></a>
</p>


# CSGen
CSGen is a Golang utility package for simplifying development of code generation tools.  It contains abstractions around the standard library AST packages to make using them more intuitive for developers.

## To Install
    go get github.com/cscoding21/csgen

## Core Uses
The primary use-case for the library is to get a list of all the structs in a file including details about each field.

    //---get all of the structs within a file.  The file argument
    structs, err := csgen.GetStructs("test_struct.go")