package simplefeatures

type Coordinates struct {
	X, Y float64
	// TODO: Put optional Z and M here.
}

type OptionalCoordinates struct {
	Empty bool
	Value Coordinates
}