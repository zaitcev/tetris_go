//
package game

// import (
//     "math/rand"
// }

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

// XXX This implementation is just dumb. It is both laborous and error-prone.
// The figure must have a different implementation: center point and direction.
// Land() has to compute the landed footprint, but rotation becomes trivial.
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
    // var f barFigure
    var f normalElFigure
    f.Init(cols, rows)
    return &f
}
