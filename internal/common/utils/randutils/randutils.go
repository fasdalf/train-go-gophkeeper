package randutils

import "math/rand"

const (
	aInt  = int('a')
	azInt = int(byte('z') - byte('a'))
)

// RandomString generates a random string of a given length.
func RandomString(length int) string {
	ran_str := make([]byte, length)
	for i := 0; i < length; i++ {
		ran_str[i] = byte(aInt + rand.Intn(azInt))
	}
	return string(ran_str)
}
