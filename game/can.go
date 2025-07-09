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

// Drop the figure down as far as it can go in the can. Return the new position.
func (c *Can) Drop(fig Figure) Figure {
    for {
        land := fig.Land()
        for i := 0; i < len(land); i++ {
            row := land[i].Row()
            if row == 0 {
                // Already touching the bottom.
                return fig
            }
            col := land[i].Column()
            if c.Matrix[(row - 1) * c.cols + col] {
                // Touching a build-up.
                return fig
            }
        }
        fig.Down()
    }
    return fig
}

func (c *Can) Land(fig Figure) {
    land := fig.Land()
    for i := 0; i < len(land); i++ {
        row := land[i].Row()
        col := land[i].Column()
        c.Matrix[row*c.cols + col] = true
    }
}

func (c *Can) CheckConflict(fig Figure) bool {
    land := fig.Land()
    for i := 0; i < len(land); i++ {
        row := land[i].Row()
        if row == 0 {
            return true
        }
        col := land[i].Column()
        if c.Matrix[(row - 1) * c.cols + col] {
            return true
        }
    }
    return false
}
