## gocodewalker

```
fileListQueue := make(chan *gocodewalker.File, 100)

fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

// restrict to only process files that have the .go extension
fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")

// handle the errors by printing them out and then ignore
errorHandler := func(e error) bool {
    fmt.Println("ERR", e.Error())
    return true
}
fileWalker.SetErrorHandler(errorHandler)

go fileWalker.Start()

for f := range fileListQueue {
    fmt.Println(f.Location)
}
```

## tiktoken-go

```
package main

import (
    "fmt"
    "github.com/tiktoken-go/tokenizer"
)

func main() {
    enc, err := tokenizer.Get(tokenizer.Cl100kBase)
    if err != nil {
        panic("oh oh")
    }

    // this should print a list of token ids
    ids, _, _ := enc.Encode("supercalifragilistic")
    fmt.Println(ids)

    // this should print the original string back
    text, _ := enc.Decode(ids)
    fmt.Println(text)
}
```