package model

import (
	"errors"
	"math"
	"math/rand"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`

	G bool
}

func Dist(p, q Point) float64 {
	return math.Sqrt(math.Pow((q.X-p.X), 2) + math.Pow((q.Y-p.Y), 2))
}

var iii = 0

func NewRandPoint(maxX, maxY float64) Point {
	p := Point{
		X: rand.Float64() * maxX,
		Y: rand.Float64() * maxY,
	}
	//	p := Point{
	//		X: 93.96288101783308,
	//		Y: 36.740258457930565,
	//	}
	return p
}

type Polygon struct {
	Points []Point
}

type Segment struct {
	A Point
	B Point
}

func (p Polygon) edges() []Segment {
	ret := make([]Segment, 0, len(p.Points))
	for i := 1; i < len(p.Points); i++ {
		ret = append(ret, Segment{A: p.Points[i-1], B: p.Points[i]})
	}

	ret = append(ret, Segment{A: p.Points[len(p.Points)-1], B: p.Points[0]})

	return ret
}

func onSegment(p, q, r Point) bool {
	if q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) &&
		q.Y <= math.Max(p.Y, r.Y) && q.Y >= math.Min(p.Y, r.Y) {
		return true
	}
	return false
}

func orientation(p, q, r Point) int {
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	switch {
	case val == 0:
		return 0
	case val > 0:
		return 1
	}
	return 2
}

func doIntersect(p1, q1, p2, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}

	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}

	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}

func edgesIntersect(e1 Segment, e2 Segment) bool {
	return doIntersect(e1.A, e1.B, e2.A, e2.B)
}

func (p Polygon) IntersectsWithSegment(e Segment) bool {
	edges := p.edges()
	for i := 0; i < len(edges); i++ {
		if edgesIntersect(edges[i], e) {
			return true
		}
	}

	return false
}

func (p Polygon) Inside(in Point) bool {
	c := 0
	edges := p.edges()
	for i := 0; i < len(edges); i++ {
		if (in.Y < edges[i].A.Y) != (in.Y < edges[i].B.Y) &&
			in.X < edges[i].A.X+((in.Y-edges[i].A.Y)/(edges[i].B.Y-edges[i].A.Y))*(edges[i].B.X-edges[i].A.X) {
			c += 1
		}
	}

	return c%2 == 1
}

type Object struct {
	Type   string   `json:"type"`
	X      *float64 `json:"x,omitempty"`
	Y      *float64 `json:"y,omitempty"`
	Points *[]Point `json:"points,omitempty"`
}

func (o Object) ToPoint() (Point, error) {
	if o.X == nil || o.Y == nil {
		return Point{}, errors.New("can't cast object to point")
	}

	return Point{
		X: *o.X,
		Y: *o.Y,
	}, nil
}

func (o Object) ToPolygon() (Polygon, error) {
	if o.Points == nil {
		return Polygon{}, errors.New("can't cast object to polygon")
	}

	return Polygon{
		Points: *o.Points,
	}, nil
}
