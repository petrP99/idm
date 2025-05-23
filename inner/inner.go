package inner

import "math/rand"

func RandomInt(a int) int {
	return a * rand.Int()
}
