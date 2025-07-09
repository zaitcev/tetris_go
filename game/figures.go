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

func (f *barFigure) Down() {

    // This should be unnecessary, because the caller only ever invokes
    // fig.Down() after checking for boundary. But we want to be safe,
    // so we add this additional checking pass before committing.
    // XXX maybe we can eliminate this with a clever re-factoring later.
    for i := 0; i < 4; i++ {
        if f.points[i].row == 0 {
            return
        }
    }

    for i := 0; i < 4; i++ {
        f.points[i].row--
    }
}

type Figure interface {
    Land() []Point     // make an imprint in the can
    Down()
}

// This only generates the figure, without checking for a conflict
// with the field in the can. We want to see the generated figure
// overlaid on top of the field even if we quit immediately.
func NewFigure(cols int, rows int) Figure {
    var f barFigure
    f.Init(cols, rows)
    return &f
}
