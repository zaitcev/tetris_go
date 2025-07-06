//
package main

import (
    "fmt"
    // "math/rand"
    "os"
)

const ROWS int = 25;  // we'll do tcgetattr later
const COFF int = 5;   // just because

func _main() error {
    erase := []byte("\033[2J")
    os.Stdout.Write(erase)
    return nil
}

func main() {
    if err := _main(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}
