package geom

import "fmt"

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rank(g Geometry) int {
	switch g.(type) {
	case EmptySet:
		return 1
	case Point:
		return 2
	case Line:
		return 3
	case LineString:
		return 4
	case LinearRing:
		return 5
	case Polygon:
		return 6
	case MultiPoint:
		return 7
	case MultiLineString:
		return 8
	case MultiPolygon:
		return 9
	case GeometryCollection:
		return 10
	default:
		panic(fmt.Sprintf("unknown geometry type: %T", g))
	}
}

func must(x Geometry, err error) Geometry {
	if err != nil {
		panic(err)
	}
	return x
}