//
package game

type barFigure struct {
    points [4]Point

    // Ouch. A figure must carry the limits, because we made it
    // provide self-checking motion methods. It's not like we use
    // much of any memory, because only one figure ever exists
    // at any given time. But it's just annoying to have truth
    // replicated in such a silly way.
    cols, rows int
}

func (f *barFigure) Land() []Point {
    return f.points[0:4]
}

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

// func (f *barFigure) Copy() *barFigure {
//     var ret barFigure
//     ret.points = f.points
//     return &ret
// }

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
    var f barFigure
    f.Init(cols, rows)
    return &f
}
