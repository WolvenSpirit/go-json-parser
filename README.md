[![Go-JSON-Parser](https://github.com/WolvenSpirit/go-json-parser/actions/workflows/go.yml/badge.svg)](https://github.com/WolvenSpirit/go-json-parser/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/WolvenSpirit/go-json-parser/branch/master/graph/badge.svg?token=6YPUD64XSC)](https://codecov.io/gh/WolvenSpirit/go-json-parser)
# Go-JSON-Parser

Go-JSON-Parser is a small library that parses json strings into maps without having to use a struct with declared fields and tags.

```go
package main

import (
    jsonparser "github.com/WolvenSpirit/go-json-parser"
)

func main() {
    str := `{
        "email":"me@example.com", "token":"xyz", "child":{
            "value": 3
        }
    }
    `
    // Parse 2 dimensional Json
    m, err := jsonparser.Parse2Dimensional(str)
    if err != nil {
        log.Println(err.Error())
    }
    fmt.Println(m["child"].Map["value"])
    // output: "3"

    // Parse one level at a time
    m, err := jsonparser.Parse(str)
    fmt.Println(m["child"])
    // output: "child":{
    //        "value": 3
    //      }

    // Parse the child json object later
    m, err := jsonparser.Parse(m["child"])
    fmt.Println(m["value"])
    // output: "3"
    
    // If you have more depth to your json just do subsequent parses if and when needed
}

```

