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
//
// Note that rows start at the bottom of the can and go up, in order
// for the rotations to match the common or school layout of 2D trigonometry.
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
    var ret Figure = fig
    for {
        next := c.Down1(ret)
        if next == nil {
            break
        }
        ret = next
    }
    return ret
}

func (c *Can) Down1(fig Figure) Figure {
    next := fig.Down()
    if next == nil {
        return nil
    }
    if c.CheckConflict(next) {
        return nil
    }
    return next
}

func (c *Can) Left1(fig Figure) Figure {
    next := fig.Left()
    if next == nil {
        return nil
    }
    if c.CheckConflict(next) {
        return nil
    }
    return next
}

func (c *Can) Right1(fig Figure) Figure {
    next := fig.Right()
    if next == nil {
        return nil
    }
    if c.CheckConflict(next) {
        return nil
    }
    return next
}

func (c *Can) Rotate(fig Figure) Figure {
    next := fig.Rotate()
    if next == nil {
        return nil
    }
    if c.CheckConflict(next) {
        return nil
    }
    return next
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
        // if row < 0 {
        //     return true
        // }
        col := land[i].Column()
        // if col >= c.cols {
        //     return true
        // }
        if c.Matrix[row*c.cols + col] {
            return true
        }
    }
    return false
}
