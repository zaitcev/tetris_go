//
package game

import (
    "time"
)

type Point struct {
    col, row int
}
func (p Point) Column() int { return p.col }
func (p Point) Row() int { return p.row }

// Initial representation is hand-rolled matrix, made out of a slice.
// We use a slice in case we want to change the size of the can dynamically.
// Note that rows start at the bottom of the can and go up, just because.
//
// This has turned into a general game state class.
type Can struct {
    cols, rows int
    Matrix []bool

    MissionStart time.Time
    MissionTime time.Duration    // last displayed
}

func NewCan(cols int, rows int) *Can {
    var ret *Can

    ret = new(Can)
    ret.cols = cols
    ret.rows = rows
    ret.Matrix = make([]bool, cols*rows)

    ret.MissionStart = time.Now()
    ret.MissionTime = 0

    return ret
}
