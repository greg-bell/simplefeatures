package geom_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	. "github.com/peterstace/simplefeatures/geom"
)

func TestIsEmptyDimension(t *testing.T) {
	for _, tt := range []struct {
		wkt       string
		wantEmpty bool
		wantDim   int
	}{
		{"POINT EMPTY", true, 0},
		{"POINT(1 1)", false, 0},
		{"LINESTRING EMPTY", true, 0},
		{"LINESTRING(0 0,1 1)", false, 1},
		{"LINESTRING(0 0,1 1,2 2)", false, 1},
		{"LINESTRING(0 0,1 1,1 0,0 0)", false, 1},
		{"LINEARRING(0 0,1 0,1 1,0 0)", false, 1},
		{"POLYGON EMPTY", true, 0},
		{"POLYGON((0 0,1 1,1 0,0 0))", false, 2},
		{"MULTIPOINT EMPTY", true, 0},
		{"MULTIPOINT((0 0))", false, 0},
		{"MULTIPOINT((0 0),(1 1))", false, 0},
		{"MULTILINESTRING EMPTY", true, 0},
		{"MULTILINESTRING((0 0,1 1,2 2))", false, 1},
		{"MULTILINESTRING(EMPTY)", true, 0},
		{"MULTIPOLYGON EMPTY", true, 0},
		{"MULTIPOLYGON(((0 0,1 0,1 1,0 0)))", false, 2},
		{"MULTIPOLYGON(((0 0,1 0,1 1,0 0)))", false, 2},
		{"MULTIPOLYGON(EMPTY)", true, 0},
		{"GEOMETRYCOLLECTION EMPTY", true, 0},
		{"GEOMETRYCOLLECTION(POINT EMPTY)", true, 0},
		{"GEOMETRYCOLLECTION(POLYGON EMPTY)", true, 0},
		{"GEOMETRYCOLLECTION(POINT(1 1))", false, 0},
		{"GEOMETRYCOLLECTION(POINT(1 1),LINESTRING(0 0,1 1))", false, 1},
		{"GEOMETRYCOLLECTION(POLYGON((0 0,1 1,1 0,0 0)),POINT(1 1),LINESTRING(0 0,1 1))", false, 2},
	} {
		t.Run(tt.wkt, func(t *testing.T) {
			geom, err := UnmarshalWKT(strings.NewReader(tt.wkt))
			if err != nil {
				t.Fatal(err)
			}
			t.Run("IsEmpty_"+tt.wkt, func(t *testing.T) {
				gotEmpty := geom.IsEmpty()
				if gotEmpty != tt.wantEmpty {
					t.Errorf("want=%v got=%v", tt.wantEmpty, gotEmpty)
				}
			})
			t.Run("Dimension_"+tt.wkt, func(t *testing.T) {
				gotDim := geom.Dimension()
				if gotDim != tt.wantDim {
					t.Errorf("want=%v got=%v", tt.wantDim, gotDim)
				}
			})
		})
	}
}

func TestEnvelope(t *testing.T) {
	xy := func(x, y float64) XY {
		return XY{x, y}
	}
	for i, tt := range []struct {
		wkt string
		min XY
		max XY
	}{
		{"POINT(1 1)", xy(1, 1), xy(1, 1)},
		{"LINESTRING(1 2,3 4)", xy(1, 2), xy(3, 4)},
		{"LINESTRING(4 1,2 3)", xy(2, 1), xy(4, 3)},
		{"LINESTRING(1 1,3 1,2 2,2 4)", xy(1, 1), xy(3, 4)},
		{"LINEARRING(1 1,3 1,2 2,2 4,1 1)", xy(1, 1), xy(3, 4)},
		{"POLYGON((1 1,3 1,2 2,2 4,1 1))", xy(1, 1), xy(3, 4)},
		{"MULTIPOINT(1 1,3 1,2 2,2 4,1 1)", xy(1, 1), xy(3, 4)},
		{"MULTILINESTRING((1 1,3 1,2 2,2 4,1 1),(4 1,4 2))", xy(1, 1), xy(4, 4)},
		{"MULTILINESTRING((4 1,4 2),(1 1,3 1,2 2,2 4,1 1))", xy(1, 1), xy(4, 4)},
		{"MULTIPOLYGON(((4 1,4 2,3 2,4 1)),((1 1,3 1,2 2,2 4,1 1)))", xy(1, 1), xy(4, 4)},
		{"GEOMETRYCOLLECTION(POINT(4 1),POINT(2 3))", xy(2, 1), xy(4, 3)},
		{"GEOMETRYCOLLECTION(GEOMETRYCOLLECTION(POINT(4 1),POINT(2 3)))", xy(2, 1), xy(4, 3)},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log("wkt:", tt.wkt)
			g := geomFromWKT(t, tt.wkt)
			env, have := g.Envelope()
			if !have {
				t.Fatalf("expected to have envelope but didn't")
			}
			if !env.Min().Equals(tt.min) {
				t.Errorf("min: got=%v want=%v", env.Min(), tt.min)
			}
			if !env.Max().Equals(tt.max) {
				t.Errorf("max: got=%v want=%v", env.Max(), tt.max)
			}
		})
	}
}

func TestNoEnvelope(t *testing.T) {
	for _, wkt := range []string{
		"POINT EMPTY",
		"MULTIPOINT EMPTY",
		"MULTILINESTRING EMPTY",
		"MULTIPOLYGON EMPTY",
		"GEOMETRYCOLLECTION EMPTY",
		"GEOMETRYCOLLECTION(POINT EMPTY)",
	} {
		t.Run(wkt, func(t *testing.T) {
			g := geomFromWKT(t, wkt)
			if _, have := g.Envelope(); have {
				t.Errorf("have envelope but expected not to")
			}
		})
	}
}

func TestIsSimple(t *testing.T) {
	for i, tt := range []struct {
		wkt        string
		wantSimple bool
	}{
		{"POINT EMPTY", true},
		{"POINT(1 2)", true},

		{"LINESTRING EMPTY", true},
		{"LINESTRING(0 0,1 2)", true},
		{"LINESTRING(0 0,1 1,1 1)", true},
		{"LINESTRING(0 0,0 0,1 1)", true},
		{"LINESTRING(0 0,1 1,0 0)", false},
		{"LINESTRING(0 0,1 1,0 1)", true},
		{"LINESTRING(0 0,1 1,0 1,0 0)", true},
		{"LINESTRING(0 0,1 1,0 1,1 0)", false},
		{"LINESTRING(0 0,1 1,0 1,1 0,0 0)", false},
		{"LINESTRING(0 0,1 1,0 1,1 0,2 0)", false},
		{"LINESTRING(0 0,1 1,0 1,0 0,1 1)", false},
		{"LINESTRING(0 0,1 1,0 1,0 0,2 2)", false},
		{"LINESTRING(1 1,2 2,0 0)", false},
		{"LINESTRING(1 1,2 2,3 2,3 3,0 0)", false},
		{"LINESTRING(0 0,1 1,2 2)", true},

		{"LINEARRING(0 0,0 1,1 0,0 0)", true},

		{"POLYGON((0 0,0 1,1 0,0 0))", true},

		{"MULTIPOINT((1 2),(3 4),(5 6))", true},
		{"MULTIPOINT((1 2),(3 4),(1 2))", false},
		{"MULTIPOINT EMPTY", true},

		{"POLYGON EMPTY", true},

		{"MULTILINESTRING EMPTY", true},
		{"MULTILINESTRING((0 0,1 0))", true},
		{"MULTILINESTRING((0 0,1 0,0 1,0 0))", true},
		{"MULTILINESTRING((0 0,1 1,2 2),(0 2,1 1,2 0))", false},
		{"MULTILINESTRING((0 0,2 1,4 2),(4 2,2 3,0 4))", true},
		{"MULTILINESTRING((0 0,2 0,4 0),(2 0,2 1))", false},

		{"MULTIPOLYGON(((0 0,1 0,0 1,0 0)))", true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := geomFromWKT(t, tt.wkt).(HeterogenousGeometry)
			got := g.IsSimple()
			if got != tt.wantSimple {
				t.Logf("wkt: %s", tt.wkt)
				t.Errorf("got=%v want=%v", got, tt.wantSimple)
			}
		})
	}
}

func TestBoundary(t *testing.T) {
	for i, tt := range []struct {
		wkt, boundary string
	}{
		{"POINT EMPTY", "POINT EMPTY"},
		{"LINESTRING EMPTY", "LINESTRING EMPTY"},
		{"POLYGON EMPTY", "POLYGON EMPTY"},
		{"MULTIPOINT EMPTY", "MULTIPOINT EMPTY"},
		{"MULTILINESTRING EMPTY", "MULTILINESTRING EMPTY"},
		{"MULTIPOLYGON EMPTY", "MULTIPOLYGON EMPTY"},

		{"POINT(1 2)", "GEOMETRYCOLLECTION EMPTY"},
		{"LINESTRING(1 2,3 4)", "MULTIPOINT(1 2,3 4)"},
		{"LINESTRING(1 2,3 4,5 6)", "MULTIPOINT(1 2,5 6)"},
		{"LINESTRING(1 2,3 4,5 6,7 8)", "MULTIPOINT(1 2,7 8)"},
		{"LINESTRING(0 0,1 0,0 1,0 0)", "MULTIPOINT EMPTY"},
		{"LINEARRING(0 0,1 0,0 1,0 0)", "MULTIPOINT EMPTY"},

		{"POLYGON((0 0,1 0,1 1,0 1,0 0))", "LINESTRING(0 0,1 0,1 1,0 1,0 0)"},
		{"POLYGON((0 0,3 0,3 3,0 3,0 0),(1 1,2 1,2 2,1 2,1 1))", "MULTILINESTRING((0 0,3 0,3 3,0 3,0 0),(1 1,2 1,2 2,1 2,1 1))"},

		{"MULTIPOINT((1 2))", "GEOMETRYCOLLECTION EMPTY"},
		{"MULTIPOINT((1 2),(3 4))", "GEOMETRYCOLLECTION EMPTY"},

		{
			"MULTILINESTRING((0 0,1 1))",
			"MULTIPOINT(0 0,1 1)",
		},
		{
			"MULTILINESTRING((0 0,1 0),(0 1,1 1))",
			"MULTIPOINT(0 0,1 0,0 1,1 1)",
		},
		{
			"MULTILINESTRING((0 0,1 1),(1 1,1 0))",
			"MULTIPOINT(0 0,1 0)",
		},
		{
			"MULTILINESTRING((0 0,1 0,1 1),(0 0,0 1,1 1))",
			"MULTIPOINT EMPTY",
		},
		{
			"MULTILINESTRING((0 0,1 1),(0 1,1 1),(1 0,1 1))",
			"MULTIPOINT(0 0,1 1,0 1,1 0)",
		},
		{
			"MULTILINESTRING((0 0,0 1,1 1),(0 1,0 0,1 0))",
			"MULTIPOINT(0 0,1 1,0 1,1 0)",
		},
		{
			"MULTILINESTRING((0 1,1 1),(1 1,1 0),(1 1,2 1),(1 2,1 1))",
			"MULTIPOINT(0 1,1 0,2 1,1 2)",
		},
		{
			"MULTILINESTRING((1 1,2 2),(1 1,2 2))",
			"MULTIPOINT EMPTY",
		},

		{
			"MULTIPOLYGON(((0 0,3 0,3 3,0 3,0 0),(1 1,2 1,2 2,1 2,1 1)),((4 0,5 0,5 1,4 1,4 0)))",
			"MULTILINESTRING((0 0,3 0,3 3,0 3,0 0),(1 1,2 1,2 2,1 2,1 1),(4 0,5 0,5 1,4 1,4 0))",
		},
		{
			"MULTIPOLYGON(((0 0,3 0,3 3,0 3,0 0)))",
			"MULTILINESTRING((0 0,3 0,3 3,0 3,0 0))",
		},

		{
			"GEOMETRYCOLLECTION EMPTY",
			"GEOMETRYCOLLECTION EMPTY",
		},
		{
			"GEOMETRYCOLLECTION(GEOMETRYCOLLECTION EMPTY)",
			"GEOMETRYCOLLECTION(GEOMETRYCOLLECTION EMPTY)",
		},
		{
			"GEOMETRYCOLLECTION(POINT EMPTY, GEOMETRYCOLLECTION EMPTY)",
			"GEOMETRYCOLLECTION(POINT EMPTY, GEOMETRYCOLLECTION EMPTY)",
		},
		{
			"GEOMETRYCOLLECTION(POINT(1 1))",
			"GEOMETRYCOLLECTION EMPTY",
		},
		{
			`GEOMETRYCOLLECTION(
				LINESTRING(1 0,0 5,5 2),
				POINT(2 3),
				POLYGON((0 0,1 0,0 1,0 0))
			)`,
			`GEOMETRYCOLLECTION(
				MULTIPOINT(1 0,5 2),
				LINESTRING(0 0,1 0,0 1,0 0)
			)`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			want := geomFromWKT(t, tt.boundary)
			got := geomFromWKT(t, tt.wkt).Boundary()
			if !reflect.DeepEqual(got, want) {
				t.Logf("want: %s", string(want.AsText()))
				t.Logf("got:  %s", string(got.AsText()))
				t.Errorf("mismatch")
			}
		})
	}
}

func TestCoordinates(t *testing.T) {
	cmp0d := func(t *testing.T, got Coordinates, want [2]float64) {
		if got.XY.X != want[0] {
			t.Errorf("coordinate mismatch: got=%v want=%v", got, want)
		}
		if got.XY.Y != want[1] {
			t.Errorf("coordinate mismatch: got=%v want=%v", got, want)
		}
	}
	cmp1d := func(t *testing.T, got []Coordinates, want [][2]float64) {
		if len(got) != len(want) {
			t.Errorf("length mismatch: got=%v want=%v", len(got), len(want))
		}
		for i := range got {
			cmp0d(t, got[i], want[i])
		}
	}
	cmp2d := func(t *testing.T, got [][]Coordinates, want [][][2]float64) {
		if len(got) != len(want) {
			t.Errorf("length mismatch: got=%v want=%v", len(got), len(want))
		}
		for i := range got {
			cmp1d(t, got[i], want[i])
		}
	}
	cmp3d := func(t *testing.T, got [][][]Coordinates, want [][][][2]float64) {
		if len(got) != len(want) {
			t.Errorf("length mismatch: got=%v want=%v", len(got), len(want))
		}
		for i := range got {
			cmp2d(t, got[i], want[i])
		}
	}
	t.Run("Point", func(t *testing.T) {
		cmp0d(t,
			geomFromWKT(t, "POINT(1 2)").(Point).Coordinates(),
			[2]float64{1, 2},
		)
	})
	t.Run("Line-LineString-LinearRing-MultiPoint", func(t *testing.T) {
		for _, tt := range []struct {
			wkt  string
			want [][2]float64
		}{
			{"LINESTRING(0 1,2 3)", [][2]float64{{0, 1}, {2, 3}}},
			{"LINESTRING(0 1,2 3,4 5)", [][2]float64{{0, 1}, {2, 3}, {4, 5}}},
			{"LINEARRING(0 0,1 0,0 1,0 0)", [][2]float64{{0, 0}, {1, 0}, {0, 1}, {0, 0}}},
			{"MULTIPOINT(0 1,2 3,4 5)", [][2]float64{{0, 1}, {2, 3}, {4, 5}}},
		} {
			cmp1d(t,
				geomFromWKT(t, tt.wkt).(interface{ Coordinates() []Coordinates }).Coordinates(),
				tt.want,
			)
		}
	})
	t.Run("Polygon-MultiLineString", func(t *testing.T) {
		for _, tt := range []struct {
			wkt  string
			want [][][2]float64
		}{
			{
				"POLYGON((0 0,0 10,10 0,0 0),(2 2,2 7,7 2,2 2))",
				[][][2]float64{
					{{0, 0}, {0, 10}, {10, 0}, {0, 0}},
					{{2, 2}, {2, 7}, {7, 2}, {2, 2}},
				},
			},
			{
				"MULTILINESTRING((0 0,0 10,10 0,0 0),(2 2,2 8,8 2,2 2))",
				[][][2]float64{
					{{0, 0}, {0, 10}, {10, 0}, {0, 0}},
					{{2, 2}, {2, 8}, {8, 2}, {2, 2}},
				},
			},
		} {
			cmp2d(t,
				geomFromWKT(t, tt.wkt).(interface{ Coordinates() [][]Coordinates }).Coordinates(),
				tt.want,
			)
		}

	})
	t.Run("MultiPolygon", func(t *testing.T) {
		const wkt = `MULTIPOLYGON(
			((0 0,0 10,10 0,0 0),(2 2,2 7,7 2,2 2)),
			((100 100,100 110,110 100,100 100),(102 102,102 107,107 102,102 102))
		)`
		cmp3d(t,
			geomFromWKT(t, wkt).(interface{ Coordinates() [][][]Coordinates }).Coordinates(),
			[][][][2]float64{
				{
					{{0, 0}, {0, 10}, {10, 0}, {0, 0}},
					{{2, 2}, {2, 7}, {7, 2}, {2, 2}},
				},
				{
					{{100, 100}, {100, 110}, {110, 100}, {100, 100}},
					{{102, 102}, {102, 107}, {107, 102}, {102, 102}},
				},
			},
		)

	})
}

func TestTransformXY(t *testing.T) {
	transform := func(in XY) XY {
		return XY{in.X * 1.5, in.Y}
	}
	for i, tt := range []struct {
		wktIn, wktOut string
	}{
		{"POINT EMPTY", "POINT EMPTY"},
		{"LINESTRING EMPTY", "LINESTRING EMPTY"},
		{"POLYGON EMPTY", "POLYGON EMPTY"},

		{"POINT(1 3)", "POINT(1.5 3)"},
		{"LINESTRING(1 2,3 4)", "LINESTRING(1.5 2,4.5 4)"},
		{"LINESTRING(1 2,3 4,5 6)", "LINESTRING(1.5 2,4.5 4,7.5 6)"},
		{"LINEARRING(0 0,0 1,1 0,0 0)", "LINEARRING(0 0,0 1,1.5 0,0 0)"},
		{"POLYGON((0 0,0 1,1 0,0 0))", "POLYGON((0 0,0 1,1.5 0,0 0))"},
		{"MULTIPOINT(0 0,0 1,1 0,0 0)", "MULTIPOINT(0 0,0 1,1.5 0,0 0)"},
		{"MULTILINESTRING((1 2,3 4,5 6))", "MULTILINESTRING((1.5 2,4.5 4,7.5 6))"},
		{"MULTIPOLYGON(((0 0,0 1,1 0,0 0)))", "MULTIPOLYGON(((0 0,0 1,1.5 0,0 0)))"},

		{"GEOMETRYCOLLECTION EMPTY", "GEOMETRYCOLLECTION EMPTY"},
		{"GEOMETRYCOLLECTION(POINT EMPTY)", "GEOMETRYCOLLECTION(POINT EMPTY)"},
		{"GEOMETRYCOLLECTION(POINT(1 2))", "GEOMETRYCOLLECTION(POINT(1.5 2))"},
		{"GEOMETRYCOLLECTION(GEOMETRYCOLLECTION(POINT(1 2)))", "GEOMETRYCOLLECTION(GEOMETRYCOLLECTION(POINT(1.5 2)))"},
	} {

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := geomFromWKT(t, tt.wktIn)
			got, err := g.TransformXY(transform)
			expectNoErr(t, err)
			want := geomFromWKT(t, tt.wktOut)
			expectDeepEqual(t, got, want)
		})
	}
}
