//
package game

import (
    "math/rand"
)

type barFigure struct {
    points [4]Point
    cols, rows int
}

type genericFigure struct {
    pos Point

    // Relative positions for 3 blocks that are not the center.
    points [3]Point

    // Ouch. A figure must carry the limits, because we made it
    // provide self-checking motion methods. It's not like we use
    // much of any memory, because only one figure ever exists
    // at any given time. But it's just annoying to have truth
    // replicated in such a silly way.
    cols, rows int
}

type normalElFigure genericFigure
type mirrorElFigure genericFigure
type zigFigure genericFigure
type zagFigure genericFigure
type squareFigure genericFigure

func (f *barFigure) Land() []Point {
    return f.points[0:4]
}

func (f *genericFigure) Land() []Point {
    ret := make([]Point, 4)
    ret[0] = f.pos
    for i := 0; i < 3; i++ {
        ret[i+1].col = f.points[i].col + f.pos.col
        ret[i+1].row = f.points[i].row + f.pos.row
    }
    return ret
}

func (f *normalElFigure) Land() []Point {  return (*genericFigure)(f).Land() }
func (f *mirrorElFigure) Land() []Point {  return (*genericFigure)(f).Land() }
func (f *zigFigure) Land() []Point {  return (*genericFigure)(f).Land() }
func (f *zagFigure) Land() []Point {  return (*genericFigure)(f).Land() }
func (f *squareFigure) Land() []Point {  return (*genericFigure)(f).Land() }

func (f *barFigure) Init(cols int, rows int) {

    row := rows - 1
    col := cols/2 - 2
    for i := 0; i < 4; i++ {
        f.points[i].col = col + i
        f.points[i].row = row
    }

    f.rows = rows
    f.cols = cols
}

func (f *normalElFigure) Init(cols int, rows int) {

    f.pos.row = rows - 3
    f.pos.col = cols/2

    //   [0]
    //   [1]
    //   [x][2]
    f.points[0].row = 2
    f.points[0].col = 0
    f.points[1].row = 1
    f.points[1].col = 0
    f.points[2].row = 0
    f.points[2].col = 1

    f.rows = rows
    f.cols = cols
}

func (f *mirrorElFigure) Init(cols int, rows int) {

    f.pos.row = rows - 3
    f.pos.col = cols/2

    //    [0]
    //    [1]
    // [2][x]
    f.points[0].row = 2
    f.points[0].col = 0
    f.points[1].row = 1
    f.points[1].col = 0
    f.points[2].row = 0
    f.points[2].col = -1

    f.rows = rows
    f.cols = cols
}

func (f *zigFigure) Init(cols int, rows int) {

    f.pos.row = rows - 2
    f.pos.col = cols/2

    //    [0]
    // [x][1]
    // [2]
    f.points[0].row = 1
    f.points[0].col = 1
    f.points[1].row = 0
    f.points[1].col = 1
    f.points[2].row = -1
    f.points[2].col = 0

    f.rows = rows
    f.cols = cols
}

func (f *zagFigure) Init(cols int, rows int) {

    f.pos.row = rows - 2
    f.pos.col = cols/2

    //    [0]
    //    [x][1]
    //       [2]
    f.points[0].row = 1
    f.points[0].col = 0
    f.points[1].row = 0
    f.points[1].col = 1
    f.points[2].row = -1
    f.points[2].col = 1

    f.rows = rows
    f.cols = cols
}

func (f *squareFigure) Init(cols int, rows int) {

    f.pos.row = rows - 2
    f.pos.col = cols/2

    //    [0][1]
    //    [x][2]
    f.points[0].row = 1
    f.points[0].col = 0
    f.points[1].row = 1
    f.points[1].col = 1
    f.points[2].row = 0
    f.points[2].col = 1

    f.rows = rows
    f.cols = cols
}

// This approach is dumb. It is both laborous and error-prone.
// We only leave it here as a monument to our shame. The barFigure
// implements the same Figure as the genericFigure does, so it plugs
// right in despite being like this.
func (f *barFigure) Rotate() Figure {
    var ret barFigure
    ret.cols = f.cols
    ret.rows = f.rows

    if f.points[0].col < f.points[2].col {          // pointing left
        ret.points[0].row = f.points[0].row + 2
        ret.points[0].col = f.points[0].col + 2
        ret.points[1].row = f.points[1].row + 1
        ret.points[1].col = f.points[1].col + 1
        ret.points[2].row = f.points[2].row
        ret.points[2].col = f.points[2].col
        ret.points[3].row = f.points[3].row - 1
        ret.points[3].col = f.points[3].col - 1
    } else if f.points[0].row > f.points[2].row {   // pointing up
        ret.points[0].row = f.points[0].row - 2
        ret.points[0].col = f.points[0].col + 2
        ret.points[1].row = f.points[1].row - 1
        ret.points[1].col = f.points[1].col + 1
        ret.points[2].row = f.points[2].row
        ret.points[2].col = f.points[2].col
        ret.points[3].row = f.points[3].row + 1
        ret.points[3].col = f.points[3].col - 1
    } else if f.points[0].col > f.points[2].col {   // pointing right
        ret.points[0].row = f.points[0].row - 2
        ret.points[0].col = f.points[0].col - 2
        ret.points[1].row = f.points[1].row - 1
        ret.points[1].col = f.points[1].col - 1
        ret.points[2].row = f.points[2].row
        ret.points[2].col = f.points[2].col
        ret.points[3].row = f.points[3].row + 1
        ret.points[3].col = f.points[3].col + 1
    } else {                                        // pointing down presumably
        ret.points[0].row = f.points[0].row + 2
        ret.points[0].col = f.points[0].col - 2
        ret.points[1].row = f.points[1].row + 1
        ret.points[1].col = f.points[1].col - 1
        ret.points[2].row = f.points[2].row
        ret.points[2].col = f.points[2].col
        ret.points[3].row = f.points[3].row - 1
        ret.points[3].col = f.points[3].col + 1
    }

    for i := 0; i < 4; i++ {
        if ret.points[i].col < 0 || ret.points[i].col >= ret.cols ||
          ret.points[i].row < 0 || ret.points[i].row >= ret.rows {
            return nil
        }
    }
    return &ret
}

// XXX should we push the figure away from the side if it cannot rotate?
func (f *genericFigure) Rotate() Figure {
    ret := *f

    // The rotation matrix for -0.5*pi is:
    //  [col] [  0   1 ]
    //  [row] [ -1   0 ]
    for i := 0; i < 3; i++ {
        ret.points[i].col = f.points[i].row
        ret.points[i].row = f.points[i].col * -1
    }

    // unnecessary for a rotation
    // if ret.pos.col < 0 || ret.pos.col >= ret.cols ||
    //    ret.pos.row < 0 || ret.pos.row >= ret.rows {
    //     return nil
    // }
    for i := 0; i < 3; i++ {
        newcol := ret.pos.col + ret.points[i].col
        newrow := ret.pos.row + ret.points[i].row
        if newcol < 0 || newcol >= ret.cols ||
          newrow < 0 || newrow >= ret.rows {
            return nil
        }
    }
    return &ret
}

func (f *normalElFigure) Rotate() Figure { return (*genericFigure)(f).Rotate() }
func (f *mirrorElFigure) Rotate() Figure { return (*genericFigure)(f).Rotate() }
func (f *zigFigure) Rotate() Figure { return (*genericFigure)(f).Rotate() }
func (f *zagFigure) Rotate() Figure { return (*genericFigure)(f).Rotate() }
func (f *squareFigure) Rotate() Figure { return (*genericFigure)(f).Rotate() }

func (f *barFigure) Down() Figure {
    var ret barFigure
    ret.cols = f.cols
    ret.rows = f.rows

    for i := 0; i < 4; i++ {
        if f.points[i].row == 0 {
            return nil
        }
        ret.points[i].row = f.points[i].row - 1
        ret.points[i].col = f.points[i].col
    }

    return &ret
}

func (f *genericFigure) Down() Figure {
    ret := *f
    ret.pos.row = f.pos.row - 1

    if ret.pos.row < 0 { return nil }
    for i := 0; i < 3; i++ {
        if ret.pos.row + ret.points[i].row < 0 { return nil }
    }
    return &ret
}

func (f *normalElFigure) Down() Figure {  return (*genericFigure)(f).Down() }
func (f *mirrorElFigure) Down() Figure {  return (*genericFigure)(f).Down() }
func (f *zigFigure) Down() Figure {  return (*genericFigure)(f).Down() }
func (f *zagFigure) Down() Figure {  return (*genericFigure)(f).Down() }
func (f *squareFigure) Down() Figure {  return (*genericFigure)(f).Down() }

func (f *barFigure) Left() Figure {
    var ret barFigure
    ret.cols = f.cols
    ret.rows = f.rows

    for i := 0; i < 4; i++ {
        ret.points[i].row = f.points[i].row
        ret.points[i].col = f.points[i].col - 1
        if ret.points[i].col < 0 {
            return nil
        }
    }

    return &ret
}

func (f *genericFigure) Left() Figure {
    ret := *f
    ret.pos.col = f.pos.col - 1

    if ret.pos.col < 0 { return nil }
    for i := 0; i < 3; i++ {
        if ret.pos.col + ret.points[i].col < 0 { return nil }
    }
    return &ret
}

func (f *normalElFigure) Left() Figure {  return (*genericFigure)(f).Left() }
func (f *mirrorElFigure) Left() Figure {  return (*genericFigure)(f).Left() }
func (f *zigFigure) Left() Figure {  return (*genericFigure)(f).Left() }
func (f *zagFigure) Left() Figure {  return (*genericFigure)(f).Left() }
func (f *squareFigure) Left() Figure {  return (*genericFigure)(f).Left() }

func (f *barFigure) Right() Figure {
    var ret barFigure
    ret.cols = f.cols
    ret.rows = f.rows

    for i := 0; i < 4; i++ {
        ret.points[i].row = f.points[i].row
        ret.points[i].col = f.points[i].col + 1
        if ret.points[i].col >= ret.cols {
            return nil
        }
    }

    return &ret
}

func (f *genericFigure) Right() Figure {
    ret := *f
    ret.pos.col = f.pos.col + 1

    if ret.pos.col >= ret.cols { return nil }
    for i := 0; i < 3; i++ {
        if ret.pos.col + ret.points[i].col >= ret.cols { return nil }
    }
    return &ret
}

func (f *normalElFigure) Right() Figure {  return (*genericFigure)(f).Right() }
func (f *mirrorElFigure) Right() Figure {  return (*genericFigure)(f).Right() }
func (f *zigFigure) Right() Figure {  return (*genericFigure)(f).Right() }
func (f *zagFigure) Right() Figure {  return (*genericFigure)(f).Right() }
func (f *squareFigure) Right() Figure {  return (*genericFigure)(f).Right() }

type Figure interface {
    Land() []Point     // make an imprint in the can
    Down() Figure
    Rotate() Figure
    Left() Figure
    Right() Figure
}

// This only generates the figure, without checking for a conflict
// with the field in the can. We want to see the generated figure
// overlaid on top of the field even if we quit immediately.
func NewFigure(cols int, rows int) Figure {
    var ret Figure
    switch rand.Intn(6) {
    case 0:
        var f0 barFigure
        f0.Init(cols, rows)
        ret = &f0
    case 1:
        var f1 normalElFigure
        f1.Init(cols, rows)
        ret = &f1
    case 2:
        var f2 mirrorElFigure
        f2.Init(cols, rows)
        ret = &f2
    case 3:
        var f3 zigFigure
        f3.Init(cols, rows)
        ret = &f3
    case 4:
        var f4 zagFigure
        f4.Init(cols, rows)
        ret = &f4
    case 5:
        var f5 squareFigure
        f5.Init(cols, rows)
        ret = &f5
    }
    return ret
}
