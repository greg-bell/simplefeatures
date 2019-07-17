package simplefeatures_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	. "github.com/peterstace/simplefeatures"
)

func hexStringToBytes(t *testing.T, s string) []byte {
	t.Helper()
	if len(s)%2 != 0 {
		t.Fatal("hex string must have even length")
	}
	var buf []byte
	for i := 0; i < len(s); i += 2 {
		x, err := strconv.ParseUint(s[i:i+2], 16, 8)
		if err != nil {
			t.Fatal(err)
		}
		buf = append(buf, byte(x))
	}
	return buf
}

func TestWKBParser(t *testing.T) {
	// Test cases generated from:
	/*
		SELECT
		        wkt,
		        ST_AsText(ST_Force2D(ST_GeomFromText(wkt))) AS flat,
		        ST_AsBinary(ST_GeomFromText(wkt)) AS wkb
		FROM (
		        VALUES
		        ('POINT EMPTY'),
		        ('POINTZ EMPTY'),
		        ('POINTM EMPTY'),
		        ('POINTZM EMPTY'),
		        ('POINT(1 2)'),
		        ('POINTZ(1 2 3)'),
		        ('POINTM(1 2 3)'),
		        ('POINTZM(1 2 3 4)'),
		        ('LINESTRING EMPTY'),
		        ('LINESTRINGZ EMPTY'),
		        ('LINESTRINGM EMPTY'),
		        ('LINESTRINGZM EMPTY'),
		        ('LINESTRING(1 2,3 4)'),
		        ('LINESTRINGZ(1 2 3,4 5 6)'),
		        ('LINESTRINGM(1 2 3,4 5 6)'),
		        ('LINESTRINGZM(1 2 3 4,5 6 7 8)'),
		        ('LINESTRING(1 2,3 4,5 6)'),
		        ('LINESTRINGZ(1 2 3,3 4 5,5 6 7)'),
		        ('LINESTRINGM(1 2 3,3 4 5,5 6 7)'),
		        ('LINESTRINGZM(1 2 3 4,3 4 5 6,5 6 7 8)'),
		        ('POLYGON EMPTY'),
		        ('POLYGONZ EMPTY'),
		        ('POLYGONM EMPTY'),
		        ('POLYGONZM EMPTY'),
		        ('POLYGON((0 0,4 0,0 4,0 0),(1 1,2 1,1 2,1 1))'),
		        ('POLYGONZ((0 0 9,4 0 9,0 4 9,0 0 9),(1 1 9,2 1 9,1 2 9,1 1 9))'),
		        ('POLYGONM((0 0 9,4 0 9,0 4 9,0 0 9),(1 1 9,2 1 9,1 2 9,1 1 9))'),
		        ('POLYGONZM((0 0 9 9,4 0 9 9,0 4 9 9,0 0 9 9),(1 1 9 9,2 1 9 9,1 2 9 9,1 1 9 9))'),
		        ('MULTIPOINT EMPTY'),
		        ('MULTIPOINTZ EMPTY'),
		        ('MULTIPOINTM EMPTY'),
		        ('MULTIPOINTZM EMPTY'),
		        ('MULTIPOINT(1 2)'),
		        ('MULTIPOINTZ(1 2 3)'),
		        ('MULTIPOINTM(1 2 3)'),
		        ('MULTIPOINTZM(1 2 3 4)'),
		        ('MULTIPOINT(1 2,3 4)'),
		        ('MULTIPOINTZ(1 2 3,3 4 5)'),
		        ('MULTIPOINTM(1 2 3,3 4 5)'),
		        ('MULTIPOINTZM(1 2 3 4,3 4 5 6)')
		) AS q (wkt);
	*/
	for i, tt := range []struct {
		wkb string
		wkt string
	}{
		{
			// POINT EMPTY
			wkb: "0101000000000000000000f87f000000000000f87f",
			wkt: "POINT EMPTY",
		},
		{
			// POINTZ EMPTY
			wkb: "01e9030000000000000000f87f000000000000f87f000000000000f87f",
			wkt: "POINT EMPTY",
		},
		{
			// POINTM EMPTY
			wkb: "01d1070000000000000000f87f000000000000f87f000000000000f87f",
			wkt: "POINT EMPTY",
		},
		{
			// POINTZM EMPTY
			wkb: "01b90b0000000000000000f87f000000000000f87f000000000000f87f000000000000f87f",
			wkt: "POINT EMPTY",
		},
		{
			// POINT(1 2)
			wkb: "0101000000000000000000f03f0000000000000040",
			wkt: "POINT(1 2)",
		},
		{
			// POINTZ(1 2 3)
			wkb: "01e9030000000000000000f03f00000000000000400000000000000840",
			wkt: "POINT(1 2)",
		},
		{
			// POINTM(1 2 3)
			wkb: "01d1070000000000000000f03f00000000000000400000000000000840",
			wkt: "POINT(1 2)",
		},
		{
			// POINTZM(1 2 3 4)
			wkb: "01b90b0000000000000000f03f000000000000004000000000000008400000000000001040",
			wkt: "POINT(1 2)",
		},
		{
			// LINESTRING EMPTY
			wkb: "010200000000000000",
			wkt: "LINESTRING EMPTY",
		},
		{
			// LINESTRINGZ EMPTY
			wkb: "01ea03000000000000",
			wkt: "LINESTRING EMPTY",
		},
		{
			// LINESTRINGM EMPTY
			wkb: "01d207000000000000",
			wkt: "LINESTRING EMPTY",
		},
		{
			// LINESTRINGZM EMPTY
			wkb: "01ba0b000000000000",
			wkt: "LINESTRING EMPTY",
		},
		{
			// LINESTRING(1 2,3 4)
			wkb: "010200000002000000000000000000f03f000000000000004000000000000008400000000000001040",
			wkt: "LINESTRING(1 2,3 4)",
		},

		{
			// LINESTRINGZ(1 2 3,4 5 6)
			wkb: "01ea03000002000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001840",
			wkt: "LINESTRING(1 2,4 5)",
		},
		{
			// LINESTRINGM(1 2 3,4 5 6)
			wkb: "01d207000002000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001840",
			wkt: "LINESTRING(1 2,4 5)",
		},
		{
			// LINESTRINGZM(1 2 3 4,5 6 7 8)
			wkb: "01ba0b000002000000000000000000f03f000000000000004000000000000008400000000000001040000000000000144000000000000018400000000000001c400000000000002040",
			wkt: "LINESTRING(1 2,5 6)",
		},
		{
			// LINESTRING(1 2,3 4,5 6)
			wkb: "010200000003000000000000000000f03f00000000000000400000000000000840000000000000104000000000000014400000000000001840",
			wkt: "LINESTRING(1 2,3 4,5 6)",
		},
		{
			// LINESTRINGZ(1 2 3,3 4 5,5 6 7)
			wkb: "01ea03000003000000000000000000f03f00000000000000400000000000000840000000000000084000000000000010400000000000001440000000000000144000000000000018400000000000001c40",
			wkt: "LINESTRING(1 2,3 4,5 6)",
		},
		{
			// LINESTRINGM(1 2 3,3 4 5,5 6 7)
			wkb: "01d207000003000000000000000000f03f00000000000000400000000000000840000000000000084000000000000010400000000000001440000000000000144000000000000018400000000000001c40",
			wkt: "LINESTRING(1 2,3 4,5 6)",
		},
		{
			// LINESTRINGZM(1 2 3 4,3 4 5 6,5 6 7 8)
			wkb: "01ba0b000003000000000000000000f03f0000000000000040000000000000084000000000000010400000000000000840000000000000104000000000000014400000000000001840000000000000144000000000000018400000000000001c400000000000002040",
			wkt: "LINESTRING(1 2,3 4,5 6)",
		},
		{
			// POLYGON EMPTY
			wkb: "010300000000000000",
			wkt: "POLYGON EMPTY",
		},
		{
			// POLYGONZ EMPTY
			wkb: "01eb03000000000000",
			wkt: "POLYGON EMPTY",
		},
		{
			// POLYGONM EMPTY
			wkb: "01d307000000000000",
			wkt: "POLYGON EMPTY",
		},
		{
			// POLYGONZM EMPTY
			wkb: "01bb0b000000000000",
			wkt: "POLYGON EMPTY",
		},
		{
			// POLYGON((0 0,4 0,0 4,0 0),(1 1,2 1,1 2,1 1))
			wkb: "010300000002000000040000000000000000000000000000000000000000000000000010400000000000000000000000000000000000000000000010400000000000000000000000000000000004000000000000000000f03f000000000000f03f0000000000000040000000000000f03f000000000000f03f0000000000000040000000000000f03f000000000000f03f",
			wkt: "POLYGON((0 0,4 0,0 4,0 0),(1 1,2 1,1 2,1 1))",
		},
		{
			// POLYGONZ((0 0 9,4 0 9,0 4 9,0 0 9),(1 1 9,2 1 9,1 2 9,1 1 9))
			wkb: "01eb030000020000000400000000000000000000000000000000000000000000000000224000000000000010400000000000000000000000000000224000000000000000000000000000001040000000000000224000000000000000000000000000000000000000000000224004000000000000000000f03f000000000000f03f00000000000022400000000000000040000000000000f03f0000000000002240000000000000f03f00000000000000400000000000002240000000000000f03f000000000000f03f0000000000002240",
			wkt: "POLYGON((0 0,4 0,0 4,0 0),(1 1,2 1,1 2,1 1))",
		},
		{
			// POLYGONM((0 0 9,4 0 9,0 4 9,0 0 9),(1 1 9,2 1 9,1 2 9,1 1 9))
			wkb: "01d3070000020000000400000000000000000000000000000000000000000000000000224000000000000010400000000000000000000000000000224000000000000000000000000000001040000000000000224000000000000000000000000000000000000000000000224004000000000000000000f03f000000000000f03f00000000000022400000000000000040000000000000f03f0000000000002240000000000000f03f00000000000000400000000000002240000000000000f03f000000000000f03f0000000000002240",
			wkt: "POLYGON((0 0,4 0,0 4,0 0),(1 1,2 1,1 2,1 1))",
		},
		{
			// POLYGONZM((0 0 9 9,4 0 9 9,0 4 9 9,0 0 9 9),(1 1 9 9,2 1 9 9,1 2 9 9,1 1 9 9))
			wkb: "01bb0b00000200000004000000000000000000000000000000000000000000000000002240000000000000224000000000000010400000000000000000000000000000224000000000000022400000000000000000000000000000104000000000000022400000000000002240000000000000000000000000000000000000000000002240000000000000224004000000000000000000f03f000000000000f03f000000000000224000000000000022400000000000000040000000000000f03f00000000000022400000000000002240000000000000f03f000000000000004000000000000022400000000000002240000000000000f03f000000000000f03f00000000000022400000000000002240",
			wkt: "POLYGON((0 0,4 0,0 4,0 0),(1 1,2 1,1 2,1 1))",
		},
		{
			// MULTIPOINT EMPTY
			wkb: "010400000000000000",
			wkt: "MULTIPOINT EMPTY",
		},
		{
			// MULTIPOINTZ EMPTY
			wkb: "01ec03000000000000",
			wkt: "MULTIPOINT EMPTY",
		},
		{
			// MULTIPOINTM EMPTY
			wkb: "01d407000000000000",
			wkt: "MULTIPOINT EMPTY",
		},
		{
			// MULTIPOINTZM EMPTY
			wkb: "01bc0b000000000000",
			wkt: "MULTIPOINT EMPTY",
		},
		{
			// MULTIPOINT(1 2)
			wkb: "0104000000010000000101000000000000000000f03f0000000000000040",
			wkt: "MULTIPOINT(1 2)",
		},
		{
			// MULTIPOINTZ(1 2 3)
			wkb: "01ec0300000100000001e9030000000000000000f03f00000000000000400000000000000840",
			wkt: "MULTIPOINT(1 2)",
		},
		{
			// MULTIPOINTM(1 2 3)
			wkb: "01d40700000100000001d1070000000000000000f03f00000000000000400000000000000840",
			wkt: "MULTIPOINT(1 2)",
		},
		{
			// MULTIPOINTZM(1 2 3 4)
			wkb: "01bc0b00000100000001b90b0000000000000000f03f000000000000004000000000000008400000000000001040",
			wkt: "MULTIPOINT(1 2)",
		},
		{
			// MULTIPOINT(1 2,3 4)
			wkb: "0104000000020000000101000000000000000000f03f0000000000000040010100000000000000000008400000000000001040",
			wkt: "MULTIPOINT(1 2,3 4)",
		},
		{
			// MULTIPOINTZ(1 2 3,3 4 5)
			wkb: "01ec0300000200000001e9030000000000000000f03f0000000000000040000000000000084001e9030000000000000000084000000000000010400000000000001440",
			wkt: "MULTIPOINT(1 2,3 4)",
		},
		{
			// MULTIPOINTM(1 2 3,3 4 5)
			wkb: "01d40700000200000001d1070000000000000000f03f0000000000000040000000000000084001d1070000000000000000084000000000000010400000000000001440",
			wkt: "MULTIPOINT(1 2,3 4)",
		},
		{
			// MULTIPOINTZM(1 2 3 4,3 4 5 6)
			wkb: "01bc0b00000200000001b90b0000000000000000f03f00000000000000400000000000000840000000000000104001b90b00000000000000000840000000000000104000000000000014400000000000001840",
			wkt: "MULTIPOINT(1 2,3 4)",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			geom, err := UnmarshalWKB(bytes.NewReader(hexStringToBytes(t, tt.wkb)))
			expectNoErr(t, err)
			expectDeepEqual(t, geom, geomFromWKT(t, tt.wkt))
		})
	}
}

func TestWKBParserInvalidGeometryType(t *testing.T) {
	// Same as POINT(1 2), but with the geometry type byte set to 0xff.
	const wkb = "01ff000000000000000000f03f0000000000000040"
	_, err := UnmarshalWKB(bytes.NewReader(hexStringToBytes(t, wkb)))
	if err == nil {
		t.Errorf("expected an error but got nil")
	}
	if !strings.Contains(err.Error(), "unknown geometry type") {
		t.Errorf("expected to be an error about unknown geometry type, but got: %v", err)
	}
}
