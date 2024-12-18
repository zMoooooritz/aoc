package spcl

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) Add(d Vector) {
	c.X += d.X
	c.Y += d.Y
}

func (c *Coordinate) Sub(d Vector) {
	c.X -= d.X
	c.Y -= d.Y
}

func (c *Coordinate) Mul(f int) {
	c.X *= f
	c.Y *= f
}

func (c Coordinate) CardinalNeighbours() []Coordinate {
	return c.neighours(CARDINAL_DIRS)
}

func (c Coordinate) IntercardinalNeighbours() []Coordinate {
	return c.neighours(INTERCARDINAL_DIRS)
}

func (c Coordinate) neighours(dirs []Vector) []Coordinate {
	neighbours := []Coordinate{}
	for _, dir := range dirs {
		neighbours = append(neighbours, Coordinate{c.X + dir.X, c.Y + dir.Y})
	}
	return neighbours

}

type Vector Coordinate

var CARDINAL_DIRS = []Vector{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
var INTERCARDINAL_DIRS = []Vector{{0, -1}, {1, 0}, {0, 1}, {-1, 0}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}

func (v *Vector) RotateCW() {
	v.X, v.Y = v.Y, -v.X
}

func (v *Vector) RotateCCW() {
	v.X, v.Y = -v.Y, v.X
}

func (v *Vector) Mirror() {
	v.X *= -1
	v.Y *= -1
}
