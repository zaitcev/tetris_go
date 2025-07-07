//
package main

import (
    "errors"
    "fmt"
    // "math/rand"
    "os"

    "golang.org/x/term"
)

const ROWS int = 20;  // we'll do tcgetattr later, maybe
const COLS int = 10;
const TROFF int = 1;   // Header line
const TCOFF int = 5;   // just because

type Event int
const (
    EV_EXIT Event = iota
    EV_ERROR
    EV_LEFT
    EV_RIGHT
    EV_SPACE
    EV_TIME
)

var lastLetter error	// Letter from beyond the grave

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
    for i := 0; i < ROWS; i++ {
        line := fmt.Sprintf("\033[%d;%dH", TROFF+i+1, TCOFF+1)
        for j := 0; j < COLS; j++ {
            v := d.Matrix[(ROWS-1-i)*COLS + j]
            if v {
                line += "[=]"
            } else {
                line += " . "
            }
        }
        d.DP.Write([]byte(line))
    }
    d.DP.Write([]byte(fmt.Sprintf("\033[1;1H")))
}

func reader(mainChan chan Event) {
    var buf []byte
    buf = make([]byte, 1)

    for {
        n, err := os.Stdin.Read(buf)
        if err != nil {
            lastLetter = err
            mainChan <- EV_ERROR
            break
        }
        if n != 1 {
            lastLetter = errors.New(fmt.Sprintf("Read %d", n))
            mainChan <- EV_ERROR
            break
        }

        if buf[0] == ([]byte("q"))[0] || buf[0] == 0x03 {
            mainChan <- EV_EXIT
            break
        }
    }
}

func _main() error {

    mainChan := make(chan Event)

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

    go reader(mainChan)

    // This is a STEM model, but only because the game is so simple.
    // If it had enemies or phenomena, it would have more threads.
    var ev Event
    for {
        ev = <- mainChan

        if ev == EV_EXIT || ev == EV_ERROR {
            break
        }
    }

    dp.DP.Write([]byte(fmt.Sprintf("\033[%d;1H", TROFF+ROWS+1)))
    term.Restore(1, termState)

    if ev == EV_ERROR {
        return lastLetter
    }
    os.Stderr.Write([]byte("Goodbye\n"))
    return nil
}

func main() {
    if err := _main(); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}
