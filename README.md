[![Go Report Card](https://goreportcard.com/badge/github.com/koepkeca/goSafeDataStruct)](https://goreportcard.com/report/github.com/koepkeca/goSafeDataStruct)

# Overview

goSafeDataStruct is a library which provides guarded data structures for use in concurrent applications.
You can read more about the design and methodology of this library here. 

# Installation

To install the library you just do a go get:

```
go get github.com/koepkeca/goSafeDataStruct
``` 

# Usage

Usage is easy, you just include the package you want to use, then call the 
New() method for the package. Make sure you call a corresponding Destroy for each
New.

There are examples of each data structure in the examples folder for each data type.

# Data Types Supported
Data Structure | Current Implementation Status | Unit Test Status | Benchmark Status
-----------|-------------------------------|------------------|-------------------
stack | yes | yes | no
queue | no | no | no
trie | no | no | no
graph | no | no | no
