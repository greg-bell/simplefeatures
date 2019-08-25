package geom_test

import (
	"testing"

	"github.com/peterstace/simplefeatures/geom"
)

func TestEqualsExact(t *testing.T) {
	wkts := map[string]string{
		"pt_a": "POINT(2 3)",
		"pt_b": "POINT(3 -1)",
		"pt_c": "POINT(2.09 2.91)",
		"pt_d": "POINT(2.08 2.92)",
		"pt_e": "POINT EMPTY",
		"pt_f": "POINT(3.125 -1)",

		"ln_a": "LINESTRING(1 2,3 4)",
		"ln_b": "LINESTRING(1 2,3 3.9)",
		"ln_c": "LINESTRING(1.1 2,3 4)",
		"ln_d": "LINESTRING(3 4,1 2)",

		"ls_a": "LINESTRING(1 2,3 4,5 6)",
		"ls_b": "LINESTRING(1 2,3 4,5 6.1)",
		"ls_c": "LINESTRING(5 6,3 4,1 2)",

		// ccw rings
		"ls_m": "LINESTRING(0 0,1 0,0 1,0 0)",
		"ls_n": "LINESTRING(1 0,0 1,0 0,1 0)",
		"ls_o": "LINESTRING(0 1,0 0,1 0,0 1)",
		// cw rings
		"ls_p": "LINESTRING(0 0,0 1,1 0,0 0)",
		"ls_q": "LINESTRING(1 0,0 0,0 1,1 0)",
		"ls_r": "LINESTRING(0 1,1 0,0 0,0 1)",

		"lr_a": "LINEARRING(0 0,0 1,1 1,1 0,0 0)",
		"lr_b": "LINEARRING(0 0,1 0,1 1,0 1,0 0)",
		"lr_c": "LINESTRING(0 0,1 0,1 1,0 1,0 0)",

		"lr_d": "LINEARRING(1 1,0 1,1 0,1 1)",
		"lr_e": "LINEARRING(1 1,0 1,1 0.1,1 1)",

		"ls_empty": "LINESTRING EMPTY",

		// TODO: Polygon
		"p_a": "POLYGON((0 0,0 1,1 0,0 0))",

		//"p_b": "POLYGON((1 1,1 9,9 1,1 1),(2 2,2 3,3 2,2 2))",
		//"p_c": "POLYGON((1 1,1 9,9 1,1 1),(4 4,4 5,5 4,4 4))",

		"p_empty": "POLYGON EMPTY",

		// TODO: MultiPoint
		// TODO: MultiLineString
		// TODO: MultiPolygon
	}

	type pair struct{ keyA, keyB string }
	eqWithTolerance := []pair{
		{"pt_a", "pt_d"},
		{"pt_c", "pt_d"},
		{"pt_f", "pt_b"},

		{"ln_a", "ln_b"},
		{"ln_b", "ln_c"},
		{"ln_a", "ln_c"},

		{"ls_a", "ls_b"},

		{"lr_b", "lr_c"},
		{"lr_d", "lr_e"},
	}

	eqWithoutOrder := []pair{
		{"ln_a", "ln_d"},
		{"ls_a", "ls_c"},

		//{"ls_m", "ls_n"},
		//{"ls_n", "ls_o"},
		//{"ls_o", "ls_m"},

		{"ls_m", "ls_n"},
		{"ls_m", "ls_o"},
		{"ls_m", "ls_p"},
		{"ls_m", "ls_q"},
		{"ls_m", "ls_r"},
		{"ls_n", "ls_o"},
		{"ls_n", "ls_p"},
		{"ls_n", "ls_q"},
		{"ls_n", "ls_r"},
		{"ls_o", "ls_p"},
		{"ls_o", "ls_q"},
		{"ls_o", "ls_r"},
		{"ls_p", "ls_q"},
		{"ls_p", "ls_r"},
		{"ls_q", "ls_r"},

		{"lr_a", "lr_b"},
		{"lr_b", "lr_c"},
		{"lr_c", "lr_a"},
	}

	for _, p := range append(eqWithTolerance, eqWithoutOrder...) {
		if _, ok := wkts[p.keyA]; !ok {
			t.Fatalf("bad test setup: %v doesn't exist", p.keyA)
		}
		if _, ok := wkts[p.keyB]; !ok {
			t.Fatalf("bad test setup: %v doesn't exist", p.keyB)
		}
	}

	t.Run("reflexive", func(t *testing.T) {
		for key, wkt := range wkts {
			t.Run(key, func(t *testing.T) {
				t.Logf("WKT: %v", wkt)
				g := geomFromWKT(t, wkt)
				t.Run("no options", func(t *testing.T) {
					if !g.EqualsExact(g) {
						t.Errorf("should be equal to itself")
					}
				})
			})
		}
	})
	t.Run("equal with tolerance", func(t *testing.T) {
		for keyA := range wkts {
			for keyB := range wkts {
				t.Run(keyA+" and "+keyB, func(t *testing.T) {
					var want bool
					if keyA == keyB {
						want = true
					}
					for _, p := range eqWithTolerance {
						if (keyA == p.keyA && keyB == p.keyB) || (keyA == p.keyB && keyB == p.keyA) {
							want = true
							break
						}
					}
					gA := geomFromWKT(t, wkts[keyA])
					gB := geomFromWKT(t, wkts[keyB])
					got := gA.EqualsExact(gB, geom.Tolerance(0.125))
					if got != want {
						t.Errorf("got=%v want=%v", got, want)
					}
				})
			}
		}
	})
	t.Run("equal ignoring order", func(t *testing.T) {
		for keyA := range wkts {
			for keyB := range wkts {
				t.Run(keyA+" and "+keyB, func(t *testing.T) {
					var want bool
					if keyA == keyB {
						want = true
					}
					for _, p := range eqWithoutOrder {
						if (keyA == p.keyA && keyB == p.keyB) || (keyA == p.keyB && keyB == p.keyA) {
							want = true
							break
						}
					}
					gA := geomFromWKT(t, wkts[keyA])
					gB := geomFromWKT(t, wkts[keyB])
					got := gA.EqualsExact(gB, geom.IgnoreOrder)
					if got != want {
						t.Errorf("got=%v want=%v", got, want)
					}
				})
			}
		}
	})
}

func TestEqualsExactOrthogonal(t *testing.T) {
	// TODO: check that the two options don't interact with each other badly.
}
