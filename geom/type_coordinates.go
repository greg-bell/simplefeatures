package geom

import (
	"strconv"
	"strings"
)

// Coordinates represents a point location. Coordinates values may be
// constructed manually using the type definition directly. Alternatively, one
// of the New(XYZM)Coordinates constructor functions can be used.
type Coordinates struct {
	// XY represents the XY position of the point location.
	XY

	// Z represents the height of the location. Its value is zero
	// for non-3D coordinate types.
	Z float64

	// M represents a user defined measure associated with the
	// location. Its value is zero for non-measure coordinate
	// types.
	M float64

	// Type indicates the coordinates type, and therefore whether
	// or not Z and M are populated.
	Type CoordinatesType
}

// String gives a string representation of the coordinates.
func (c Coordinates) String() string {
	var sb strings.Builder
	sb.WriteString("Coordinates[")
	sb.WriteString(c.Type.String())
	sb.WriteString("] ")
	sb.WriteString(strconv.FormatFloat(c.X, 'f', -1, 64))
	sb.WriteRune(' ')
	sb.WriteString(strconv.FormatFloat(c.Y, 'f', -1, 64))
	if c.Type.Is3D() {
		sb.WriteRune(' ')
		sb.WriteString(strconv.FormatFloat(c.Z, 'f', -1, 64))
	}
	if c.Type.IsMeasured() {
		sb.WriteRune(' ')
		sb.WriteString(strconv.FormatFloat(c.M, 'f', -1, 64))
	}
	return sb.String()
}

// appendFloat64s appends the coordinates to dst, taking into
// consideration the coordinate type.
func (c Coordinates) appendFloat64s(dst []float64) []float64 {
	switch c.Type {
	case DimXY:
		return append(dst, c.X, c.Y)
	case DimXYZ:
		return append(dst, c.X, c.Y, c.Z)
	case DimXYM:
		return append(dst, c.X, c.Y, c.M)
	case DimXYZM:
		return append(dst, c.X, c.Y, c.Z, c.M)
	default:
		panic(c.Type.String())
	}
}

// NewXYCoordinates constructs a new set of coordinates of type XY.
func NewXYCoordinates(x, y float64) Coordinates {
	return Coordinates{
		Type: DimXY,
		XY:   XY{x, y},
	}
}

// NewXYZCoordinates constructs a new set of coordinates of type XYZ.
func NewXYZCoordinates(x, y, z float64) Coordinates {
	return Coordinates{
		Type: DimXYZ,
		XY:   XY{x, y},
		Z:    z,
	}
}

// NewXYMCoordinates constructs a new set of coordinates of type XYM.
func NewXYMCoordinates(x, y, m float64) Coordinates {
	return Coordinates{
		Type: DimXYM,
		XY:   XY{x, y},
		M:    m,
	}
}

// NewXYZMCoordinates constructs a new set of coordinates of type XYZM.
func NewXYZMCoordinates(x, y, z, m float64) Coordinates {
	return Coordinates{
		Type: DimXYZM,
		XY:   XY{x, y},
		Z:    z,
		M:    m,
	}
}
