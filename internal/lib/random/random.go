package random

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewRandomString(stringLength int) string {
	randomString := make([]rune, stringLength)

	for index := range randomString {
		randomString[index] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(randomString)
}
