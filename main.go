//
package main

import (
    "errors"
    "fmt"
    // "math/rand"
    "os"
    "time"

    "golang.org/x/term"

    // main.go:13:5: package game is not in std (/usr/lib/golang/src/game)
    // found packages game (game.go) and main (main.go)
    // XXX move to subdirectory?
    "zaitcev.us/tetris/game"
)

const ROWS int = 20;  // we'll do tcgetattr later, maybe
const COLS int = 10;
const TROFF int = 1;   // Header line
const TCOFF int = 5;   // just because

type Event int
const (
    EV_EXIT Event = iota
    EV_ERROR
    EV_UP
    EV_LEFT
    EV_RIGHT
    EV_SPACE
    EV_TIME
)

var lastLetter error	// Letter from beyond the grave

type Display struct {
    // This looks exactly like a Can but represents what is displayed
    // at present.
    Matrix []bool
    DP *os.File
}

type ReaderState int
const (
    R_IDLE ReaderState = iota
    R_ESC
    R_CSI
)
type Reader struct {
    State ReaderState
}

func NewDisplay() *Display {
    var ret *Display

    ret = new(Display)
    ret.Matrix = make([]bool, COLS*ROWS)
    ret.DP = os.Stdout
    return ret
}

func (d *Display) Erase() {
    d.DP.Write([]byte("\033[2J"))
}

// Update only updates the representation of the game field within the can.
// It does not update the mission status banner. Maybe change it?
func (d *Display) Update(newcan *game.Can, curfig game.Figure) {

    field := make([]bool, COLS*ROWS)
    copy(field, newcan.Matrix)

    // XXX Check for a conflict of the new figure with the can, it's game over.
    land := curfig.Land()
    for i := range land {
        field[land[i].Row()*COLS + land[i].Column()] = true
    }

    // Full refresh
    // We can at least economize on syscalls.
    for i := 0; i < ROWS; i++ {
        line := fmt.Sprintf("\033[%d;%dH", TROFF+i+1, TCOFF+1)
        for j := 0; j < COLS; j++ {
            v := field[(ROWS-1-i)*COLS + j]
            if v {
                line += "[=]"
            } else {
                line += " . "
            }
        }
        d.DP.Write([]byte(line))
    }
    d.DP.Write([]byte(fmt.Sprintf("\033[1;1H")))
    d.Matrix = field
}

func (d *Display) Time(mt time.Duration) {
    line := fmt.Sprintf("\033[%d;%dH", 1, TCOFF+20+1)
    line += fmt.Sprintf("%ds", int(mt.Seconds()))
    line += fmt.Sprintf("\033[1;1H")
    d.DP.Write([]byte(line))
}

func reader(mainChan chan Event) {

    var state Reader

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

        // XXX This needs a timer to catch a standalone ESC.
        switch state.State {
        case R_IDLE:
            if buf[0] == 0x1b {
                state.State = R_ESC
            } else if buf[0] == ' ' {
                mainChan <- EV_SPACE
            }
        case R_ESC:
            if buf[0] == 0x5b {
                state.State = R_CSI
            } else {
                state.State = R_IDLE
            }
        case R_CSI:
            if buf[0] == 'A' {
                mainChan <- EV_UP
            } else if buf[0] == 'D' {
                mainChan <- EV_LEFT
            } else if buf[0] == 'C' {
                mainChan <- EV_RIGHT
            }
            state.State = R_IDLE
        }
    }
}

func timer(mainChan chan Event) {
    for {
        // XXX de-skew
        time.Sleep(1000 * time.Millisecond)
        mainChan <- EV_TIME
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

    can := game.NewCan(COLS, ROWS)
    dp := NewDisplay()
    dp.Erase()

    curfig := game.NewFigure(COLS, ROWS)
    dp.Update(can, curfig)

    go reader(mainChan)
    go timer(mainChan)
    dp.Time(can.MissionTime)

    // This is a STEM model, but only because the game is so simple.
    // If it had enemies or phenomena, it would have more threads.
    /// XXX do something for a situation of the updater not keeping up
    var ev Event
    for {
        ev = <- mainChan

        if ev == EV_EXIT || ev == EV_ERROR {
            break
        }

        if ev == EV_SPACE {          // land
        } else if ev == EV_UP {      // rotate
        } else if ev == EV_LEFT {    // left
        } else if ev == EV_RIGHT {   // right
        } else {  // EV_TIME
            d := time.Since(can.MissionStart)
            if int(d.Seconds()) != int(can.MissionTime.Seconds()) {
                can.MissionTime = d
                dp.Time(can.MissionTime)
            }
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
