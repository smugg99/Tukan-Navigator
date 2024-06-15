package validator

// import (
// 	"net"
// 	"regexp"
// 	"strings"
// 	"unicode"
// 	"unicode/utf8"

// 	"smuggr.xyz/tukan-navigator/common/logger"
// )

// const (
// 	maxUsernameLength = 32
// 	minUsernameLength = 8
// 	minLoginLength    = 8
// 	maxLoginLength    = 16
// 	maxPasswordLength = 32
// 	minPasswordLength = 8

// 	minClientIDLength = 1
// 	maxClientIDLength = 23

// 	minUppercase = 1
// 	minLowercase = 1
// 	minDigit     = 1
// 	minSpecial   = 1
// )

// var (
// 	reservedClientIDs = []string{"inline"}
// 	reservedLogins    = []string{"administrator"}
// )

// func lengthInRange(s string, minLength, maxLength int) bool {
// 	length := utf8.RuneCountInString(s)
// 	return length >= minLength && length <= maxLength
// }

// func containsCategory(s string, category func(rune) bool, n int) bool {
// 	count := 0
// 	for _, char := range s {
// 		if category(char) {
// 			count++
// 			if count >= n {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
