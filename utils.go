package main

import "math/rand"

var letters = []rune("abcdefghijklkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genRandomUsername() string {
	str := make([]rune, 4)

	for index := range str {
		str[index] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
