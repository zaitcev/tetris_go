//
package game

type barFigure struct {
    points [4]Point
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
}

func (f *barFigure) Down() Figure {
    var ret barFigure

    for i := 0; i < 4; i++ {
        if f.points[i].row == 0 {
            return nil
        }
        ret.points[i].row = f.points[i].row - 1
        ret.points[i].col = f.points[i].col
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
}

// This only generates the figure, without checking for a conflict
// with the field in the can. We want to see the generated figure
// overlaid on top of the field even if we quit immediately.
func NewFigure(cols int, rows int) Figure {
    var f barFigure
    f.Init(cols, rows)
    return &f
}
