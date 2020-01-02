package geom_test

import (
	"io/ioutil"
	"testing"

	. "github.com/peterstace/simplefeatures/geom"
)

func TestZeroGeometry(t *testing.T) {
	var z Geometry
	expectDeepEqual(t, z.IsGeometryCollection(), true)
	z.AsGeometryCollection() // Doesn't crash.
	expectDeepEqual(t, z.AsText(), "GEOMETRYCOLLECTION EMPTY")
	z.AsBinary(ioutil.Discard) // Doesn't crash
	expectDeepEqual(t, z.Dimension(), 0)

	// TODO: continue further tests
}
