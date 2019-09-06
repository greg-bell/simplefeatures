package geom

import (
	"database/sql/driver"
	"errors"
	"io"
	"math"
)

// EmptySet is a 0-dimensional geometry that represents the empty pointset.
type EmptySet struct {
	wkt       string
	wkbType   uint32
	jsonType  string
	dimension int
}

func NewEmptyPoint(opts ...ConstructorOption) EmptySet {
	return EmptySet{"POINT EMPTY", wkbGeomTypePoint, "Point", 0}
}

func NewEmptyLineString(opts ...ConstructorOption) EmptySet {
	return EmptySet{"LINESTRING EMPTY", wkbGeomTypeLineString, "LineString", 1}
}

func NewEmptyPolygon(opts ...ConstructorOption) EmptySet {
	return EmptySet{"POLYGON EMPTY", wkbGeomTypePolygon, "Polygon", 2}
}

func (e EmptySet) AsText() string {
	return e.wkt
}

func (e EmptySet) AppendWKT(dst []byte) []byte {
	return append(dst, e.wkt...)
}

func (e EmptySet) IsSimple() bool {
	return true
}

func (e EmptySet) Intersection(g Geometry) Geometry {
	return intersection(e, g)
}

func (e EmptySet) Intersects(g Geometry) bool {
	has, _ := hasIntersection(e, g)
	return has
}

func (e EmptySet) IsEmpty() bool {
	return true
}

func (e EmptySet) Dimension() int {
	return e.dimension
}

func (e EmptySet) Equals(other Geometry) bool {
	return equals(e, other)
}

func (e EmptySet) Envelope() (Envelope, bool) {
	return Envelope{}, false
}

func (e EmptySet) Boundary() Geometry {
	return e
}

func (e EmptySet) Value() (driver.Value, error) {
	return wkbAsBytes(e)
}

func (e EmptySet) AsBinary(w io.Writer) error {
	marsh := newWKBMarshaller(w)
	marsh.writeByteOrder()
	marsh.writeGeomType(e.wkbType)
	switch e.wkbType {
	case wkbGeomTypePoint:
		marsh.writeFloat64(math.NaN())
		marsh.writeFloat64(math.NaN())
	case wkbGeomTypeLineString, wkbGeomTypePolygon:
		marsh.writeCount(0)
	default:
		marsh.setErr(errors.New("unknown empty geometry type (this shouldn't ever happen)"))
	}
	return marsh.err
}

// ConvexHull returns the convex hull of this geometry. The convex hull of an
// empty set is always an empty set.
func (e EmptySet) ConvexHull() Geometry {
	return convexHull(e)
}

func (e EmptySet) convexHullPointSet() []XY {
	return nil
}

func (e EmptySet) MarshalJSON() ([]byte, error) {
	return marshalGeoJSON(e.jsonType, []int{})
}

// TransformXY transforms this EmptySet into another EmptySet according to
// fn. It does this by ignoring fn and returning itself.
func (e EmptySet) TransformXY(fn func(XY) XY, opts ...ConstructorOption) (Geometry, error) {
	return e, nil
}

// EqualsExact checks if this EmptySet is exactly equal to another geometry
// by checking if the other geometry is an empty set of the same type.
func (e EmptySet) EqualsExact(other Geometry, opts ...EqualsExactOption) bool {
	o, ok := other.(EmptySet)
	return ok && e.wkbType == o.wkbType
}

// IsValid checks if this EmptySet is valid. However, this is no constraints on
// EmptySet, so this function always returns true
func (e EmptySet) IsValid() bool {
	return true
}
