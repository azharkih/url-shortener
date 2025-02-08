package hash

import (
	"math/rand"
	"time"
)

// Генерация случайной строки
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandString(countOfChars int) string {
	const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, countOfChars)
	for i := range bytes {
		bytes[i] = charSet[rng.Intn(len(charSet))]
	}
	return string(bytes)
}
