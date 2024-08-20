package toolkit

import "crypto/rand"

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Tools struct {
}

func (t *Tools) RandomString(length int) string {
	s, r := make([]rune, length), []rune(randomStringSource)
	randomStringSourceLen := len(r)
	for i := range s {
		prime, err := rand.Prime(rand.Reader, randomStringSourceLen)
		if err != nil {
			return ""
		}
		x, y := prime.Uint64(), uint64(randomStringSourceLen)
		s[i] = r[x%y]
	}
	return string(s)
}
