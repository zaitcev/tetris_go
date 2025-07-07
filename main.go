//
package main

import (
    "fmt"
    // "math/rand"
    "os"

    "golang.org/x/term"
)

const ROWS int = 20;  // we'll do tcgetattr later, maybe
const COLS int = 10;
const TROFF int = 1;   // Header line
const TCOFF int = 5;   // just because

// Initial representation is hand-rolled matrix, made out of a slice.
// We use a slice in case we want to change the size of the can dynamically.
// Note that rows start at the bottom of the can and go up, just because.
type Can struct {
    Matrix []bool
}

type Display struct {
    // This looks exactly like a Can but represents what is displayed
    // at present.
    Matrix []bool
    DP *os.File
}

func NewCan() *Can {
    var ret *Can

    ret = new(Can)
    ret.Matrix = make([]bool, COLS*ROWS)
    return ret
}

func NewDisplay() *Display {
    var ret *Display

    ret = new(Display)
    ret.Matrix = make([]bool, COLS*ROWS)
    ret.DP = os.Stdout
    return ret
}

func (d *Display) Update(newcan *Can) {
    // Full refresh
    // We can at least economize on syscalls.
    d.DP.Write([]byte(fmt.Sprintf("\033[%d;1H", TROFF+1)))
    for i := 0; i < ROWS; i++ {
        line := ""
        for j := 0; j < COLS; j++ {
            v := d.Matrix[(ROWS-1-i)*COLS + j]
            if v {
                line += "[=]"
            } else {
                line += " . "
            }
        }
        line += "\r\n"
        d.DP.Write([]byte(line))
    }
}

func reader(exitChan chan int) {
    var buf []byte
    buf = make([]byte, 1)

    _, err := os.Stdin.Read(buf)

    if err != nil {
        exitChan <- 1
    } else {
        exitChan <- 0
    }
}

func _main() error {

    exitChan := make(chan int)

    // This might need to be associated with the DP, use int(os.Stdin.Fd()).
    // But our current Display does not have open and close.
    // XXX Save or extract the keyboard interrupt character.
    termState, err := term.MakeRaw(1)
    if err != nil {
        return err
    }
    // We want to print something after we restore the terminal.
    // defer term.Restore(1, termState)

    erase := []byte("\033[2J")
    os.Stdout.Write(erase)

    can := NewCan()
    dp := NewDisplay()

    dp.Update(can)

    go reader(exitChan)

    var n int
    n = <- exitChan

    dp.DP.Write([]byte(fmt.Sprintf("\033[%d;1H", TROFF+ROWS+1)))
    term.Restore(1, termState)

    os.Stderr.Write([]byte(fmt.Sprintf("Goodbye %d\n", n)))

    return nil
}

func main() {
    if err := _main(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}
