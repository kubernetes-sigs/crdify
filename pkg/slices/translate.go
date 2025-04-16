package slices

// Translate is a generic function for translating an input slice
// of type S into a new slice of type E
func Translate[S any, E any](translation func(S) E, in ...S) []E {
	e := []E{}
	for _, s := range in {
		e = append(e, translation(s))
	}
	return e
}
