//
package main

import (
    "errors"
    "fmt"
    "os"
    "strings"
    "time"

    "golang.org/x/term"

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

    for i := 0; i < ROWS; i++ {
        var line []string
        line = append(line, fmt.Sprintf("\033[%d;%dH", TROFF+i+1, TCOFF+1))
        for j := 0; j < COLS; j++ { line = append(line, " . ") }
        d.DP.Write([]byte(strings.Join(line, "")))
    }
    d.DP.Write([]byte(fmt.Sprintf("\033[1;1H")))
}

// Update only updates the representation of the game field within the can.
// It does not update the mission status banner. Maybe change it?
func (d *Display) Update(newcan *game.Can, curfig game.Figure) {

    field := make([]bool, COLS*ROWS)
    copy(field, newcan.Matrix)

    // Note that this blithedly lands regardless of the conflicts.
    // We do this because that is what we display at the last moment
    // when the can is full and a new figure appears.
    if curfig != nil {
        land := curfig.Land()
        for i := range land {
            field[land[i].Row()*COLS + land[i].Column()] = true
        }
    }

    // For now we only do a very basic row-by-row optimization.
    // This way we don't have to inject additional cursor positions.
    for i := 0; i < ROWS; i++ {
        if !d.rowsEqual(field, ROWS-1-i) {
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
    }
    d.DP.Write([]byte(fmt.Sprintf("\033[1;1H")))
    d.Matrix = field
}

// This is a technical subordinate of Display.Update.
func (d *Display) rowsEqual(field []bool, row int) bool {
    old_row := d.Matrix[row*COLS : (row+1)*COLS]
    new_row := field[row*COLS : (row+1)*COLS]
    for i, v := range old_row { if v != new_row[i] { return false } }
    return true
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

func landAndCollapse(dp *Display, can *game.Can, curfig *game.Figure) {

    // Step 1: Land
    can.Land(*curfig)
    *curfig = nil

    // Step 2: Collapse
    for row := 0; row < ROWS; row++ {
        for can.RowIsFull(row) {
            can.Collapse(row)
            dp.Update(can, nil)
            time.Sleep(50 * time.Millisecond)  // make it visible
        }
    }
}

func _main() error {
    var curfig game.Figure

    mainChan := make(chan Event)

    // This might need to be associated with the DP, use int(os.Stdin.Fd()).
    // But our current Display does not have open and close.
    // XXX Save or extract the keyboard interrupt character.
    termFd := int(os.Stdin.Fd())   // Fd() on Windows is uintptr, make it int
    termState, err := term.MakeRaw(termFd)
    if err != nil {
        return err
    }
    // We want to print something after we restore the terminal.
    // defer term.Restore(termFd, termState)

    can := game.NewCan(COLS, ROWS)
    dp := NewDisplay()
    dp.Erase()

    curfig = game.NewFigure(COLS, ROWS)
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

        if ev == EV_SPACE {          // drop
            curfig = can.Drop(curfig)
            dp.Update(can, curfig)
            landAndCollapse(dp, can, &curfig)
            curfig = game.NewFigure(COLS, ROWS)
            dp.Update(can, curfig)
            if can.CheckConflict(curfig) {
                break
            }
        } else if ev == EV_UP {      // rotate
            next := can.Rotate(curfig)
            if next != nil {
                curfig = next
                dp.Update(can, curfig)
            }
        } else if ev == EV_LEFT {    // left
            next := can.Left1(curfig)
            if next != nil {
                curfig = next
                dp.Update(can, curfig)
            }
        } else if ev == EV_RIGHT {   // right
            next := can.Right1(curfig)
            if next != nil {
                curfig = next
                dp.Update(can, curfig)
            }
        } else {  // EV_TIME
            d := time.Since(can.MissionStart)
            if int(d.Seconds()) != int(can.MissionTime.Seconds()) {
                can.MissionTime = d
                dp.Time(can.MissionTime)
            }

            // XXX multiples of timer for better resolution
            // XXX reset timer when new figure appears at drop
            next := can.Down1(curfig)
            if next == nil {
                landAndCollapse(dp, can, &curfig)
                curfig = game.NewFigure(COLS, ROWS)
                dp.Update(can, curfig)
                if can.CheckConflict(curfig) {
                    break
                }
            } else {
                curfig = next
                dp.Update(can, curfig)
            }
        }
    }

    dp.DP.Write([]byte(fmt.Sprintf("\033[%d;1H", TROFF+ROWS+1)))
    term.Restore(termFd, termState)

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
