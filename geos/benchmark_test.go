package geos_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/peterstace/simplefeatures/geom"
	"github.com/peterstace/simplefeatures/geos"
)

// regularPolygon computes a regular polygon circumscribed by a circle with the
// given center and radius. Sides must be at least 3 or it will panic.
func regularPolygon(center geom.XY, radius float64, sides int) geom.Polygon {
	if sides <= 2 {
		panic(sides)
	}
	coords := make([]float64, 2*(sides+1))
	for i := 0; i < sides; i++ {
		angle := math.Pi/2 + float64(i)/float64(sides)*2*math.Pi
		coords[2*i+0] = center.X + math.Cos(angle)*radius
		coords[2*i+1] = center.Y + math.Sin(angle)*radius
	}
	coords[2*sides+0] = coords[0]
	coords[2*sides+1] = coords[1]
	ring, err := geom.NewLineString(geom.NewSequence(coords, geom.DimXY))
	if err != nil {
		panic(err)
	}
	poly, err := geom.NewPolygonFromRings([]geom.LineString{ring})
	if err != nil {
		panic(err)
	}
	return poly
}

func BenchmarkIntersection(b *testing.B) {
	for _, sz := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("n=%d", sz), func(b *testing.B) {
			inputA := regularPolygon(geom.XY{X: 0, Y: 0}, 1.0, sz).AsGeometry()
			inputB := regularPolygon(geom.XY{X: 1, Y: 0}, 1.0, sz).AsGeometry()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := geos.Intersection(inputA, inputB)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
