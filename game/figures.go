//
package game

type barFigure struct {
    points [4]Point
}

func (f barFigure) Land() []Point {
    return f.points[0:4]
}

func (f *barFigure) Init(cols int, rows int) {
    row := rows - 1
    col := cols/2 - 2
    for i := 0; i < 4; i++ {
        f.points[i].col = col + i
        f.points[i].row = row
    }
}

type Figure interface {
    Land() []Point     // make an imprint in the can
}

// This only generates the figure, without checking for a conflict
// with the field in the can.
func NewFigure(cols int, rows int) Figure {
    var f barFigure
    f.Init(cols, rows)
    return f
}
