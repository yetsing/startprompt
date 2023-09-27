package startprompt

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) equal(other *Coordinate) bool {
	return c.X == other.X && c.Y == other.Y
}

// gt 大于
func (c *Coordinate) gt(other *Coordinate) bool {
	if c.Y > other.Y {
		return true
	} else if c.Y == other.Y {
		return c.X > other.X
	} else {
		return false
	}
}

func (c *Coordinate) add(x int, y int) {
	c.X += x
	c.Y += y
}

func (c *Coordinate) addX(n int) {
	c.X += n
}

func (c *Coordinate) addY(n int) {
	c.Y += n
}

type Location struct {
	Row int
	Col int
}

type section struct {
	start Location
	end   Location
}

type area struct {
	start Coordinate
	end   Coordinate
}

func (a *area) Contains(coordinate Coordinate) bool {
	start := a.getStart()
	end := a.getEnd()
	if start.Y == end.Y {
		return start.Y == coordinate.Y && start.X <= coordinate.X && coordinate.X < end.X
	}
	if coordinate.Y == start.Y {
		return start.X <= coordinate.X
	} else if coordinate.Y == end.Y {
		return coordinate.X < end.X
	}
	return coordinate.Y > start.Y && coordinate.Y < end.Y
}

// RectContains 判断点是否在开始和结束组成的矩形中
func (a *area) RectContains(coordinate Coordinate) bool {
	start := a.getStart()
	end := a.getEnd()
	return start.Y <= coordinate.Y && coordinate.Y < end.Y && start.X <= coordinate.X && coordinate.X < end.X
}

func (a *area) isEmpty() bool {
	start := a.getStart()
	end := a.getEnd()
	return start.Y > end.Y || (start.Y == end.Y && start.X >= end.X)
}

func (a *area) getStart() Coordinate {
	if a.start.gt(&a.end) {
		return a.end
	}
	return a.start
}

func (a *area) getEnd() Coordinate {
	if a.start.gt(&a.end) {
		return a.start
	}
	return a.end
}

// limitTo 坐标超出 coordinate 的设为 coordinate
func (a *area) limitTo(coordinate Coordinate) {
	if a.start.gt(&coordinate) {
		a.start = coordinate
	}
	if a.end.gt(&coordinate) {
		a.end = coordinate
	}
}
