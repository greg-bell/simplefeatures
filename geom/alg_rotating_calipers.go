package geom

import (
	"fmt"
	"math"
)

// RotatedMinimumAreaBoundingRectangle finds a rectangle with minimum area that
// fully encloses the geometry. If the geometry is empty, the empty geometry of
// the same type is returned. If the bounding rectangle would be degenerate
// (zero area), then point or line string (with a single line segment) will be
// returned.
func RotatedMinimumAreaBoundingRectangle(g Geometry) Geometry {
	hull := g.ConvexHull()
	if hull.IsEmpty() {
		return hull
	}
	var seq Sequence
	switch hull.Type() {
	case TypePoint, TypeLineString:
		return hull
	case TypePolygon:
		seq = hull.MustAsPolygon().ExteriorRing().Coordinates()
	default:
		panic(fmt.Sprintf("unexpected convex hull geometry type: %s", hull.Type()))
	}

	rect := findBestMBR(seq)
	return rect.asPoly().AsGeometry()
}

type rotatedRectangle struct {
	origin XY // one of the corners
	span1  XY // origin to first adjacent corner
	span2  XY // origin to second adjacent corner
}

// asPoly converts the rectangle to a polygon by traversing from the
// rectangle's origin to its other corners via its spans.
func (r rotatedRectangle) asPoly() Polygon {
	pts := [5]XY{
		r.origin,
		r.origin.Add(r.span1),
		r.origin.Add(r.span1).Add(r.span2),
		r.origin.Add(r.span2),
		r.origin,
	}
	coords := make([]float64, 2*len(pts))
	for i, pt := range pts {
		coords[2*i+0] = pt.X
		coords[2*i+1] = pt.Y
	}
	ring := NewLineString(NewSequence(coords, DimXY))
	poly := NewPolygon([]LineString{ring})
	return poly
}

// findBestMBR finds the minimum area bounding rectangle for a convex ring
// specified as a sequence. It does this by enumerating each candidate rotated
// bounding rectangle, and finding the one with the minimum area. There is a
// candidate rectangle corresponding to each edge in the convex ring.
func findBestMBR(seq Sequence) rotatedRectangle {
	rhs := caliper{orient: func(in XY) XY { return in }}
	far := caliper{orient: XY.rotateCCW90}
	lhs := caliper{orient: XY.rotate180}

	var rect rotatedRectangle
	bestArea := math.Inf(+1)

	for i := 0; i+1 < seq.Length(); i++ {
		rhs.update(seq, i)
		if i == 0 {
			far.idx = rhs.idx
		}
		far.update(seq, i)
		if i == 0 {
			lhs.idx = far.idx
		}
		lhs.update(seq, i)

		area := rhs.proj.Sub(lhs.proj).Cross(far.proj)
		if area < bestArea {
			bestArea = area
			rect = rotatedRectangle{
				origin: seq.GetXY(i).Add(lhs.proj),
				span1:  rhs.proj.Sub(lhs.proj),
				span2:  far.proj,
			}
		}
	}
	return rect
}

// caliper is a helper struct for finding the maximum perpendicular distance
// between a line segment on the (convex) ring and a point on the same ring. The
// perpendicular distance may be to the "right hand side" the ring, "across"
// the ring, or to the "left hand side" of the ring.
//
// It assumes that the line segment measured against rotates around the ring
// iteratively (which is one of the key properties that allows this algorithm
// to execute quickly).
//
// This is an example of a "rotating calipers" algorithm. See
// [rotating_calipers] for a general explanation of the technique.
//
// [rotating_calipers]: https://en.wikipedia.org/wiki/Rotating_calipers
type caliper struct {
	orient func(XY) XY
	idx    int
	proj   XY
}

func (c *caliper) update(seq Sequence, lnIdx int) {
	offset := seq.GetXY(lnIdx)
	dir := seq.GetXY(lnIdx + 1).Sub(offset)
	dir = c.orient(dir)

	n := seq.Length()
	pt := func() XY {
		return seq.GetXY(c.idx).Sub(offset)
	}

	d0 := pt().Dot(dir)
	for {
		c.idx = (c.idx + 1) % n
		d1 := pt().Dot(dir)
		if d1 < d0 {
			c.idx = (c.idx - 1 + n) % n
			c.proj = pt().proj(dir)
			break
		}
		d0 = d1
	}
}
