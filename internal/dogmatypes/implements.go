package dogmatypes

import "go/types"

// Implements returns true if a value of type t (or a pointer to such a value)
// implements i.
func Implements(t types.Type, i *types.Interface) (types.Type, bool) {
	if types.Implements(t, i) {
		return t, true
	}

	t = types.NewPointer(t)

	if types.Implements(t, i) {
		return t, true
	}

	return nil, false
}
