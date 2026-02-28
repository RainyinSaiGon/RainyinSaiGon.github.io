package main

import (
"log"

"portfolio/internal/builder"
)

func main() {
err := builder.Build(builder.Config{
ContentDir: "content",
OutputDir:  "docs",
})
if err != nil {
log.Fatal(err)
}
}
