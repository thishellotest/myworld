package lib

import (
	"math/rand"
	"time"
)

func GeneratePassword(length int) string {

	const (
		lowerCharset   = "abcdefghijklmnopqrstuvwxyz"
		upperCharset   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numberCharset  = "0123456789"
		specialCharset = "!@#$%^&*()"
	)

	if length < 4 {
		panic("密码长度必须至少为4，以保证包含所有类型的字符")
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 确保每种字符类型至少包含一个
	password := []byte{
		lowerCharset[seededRand.Intn(len(lowerCharset))],
		upperCharset[seededRand.Intn(len(upperCharset))],
		numberCharset[seededRand.Intn(len(numberCharset))],
		specialCharset[seededRand.Intn(len(specialCharset))],
	}

	// 填充剩余的字符
	allCharset := lowerCharset + upperCharset + numberCharset + specialCharset
	for i := 4; i < length; i++ {
		password = append(password, allCharset[seededRand.Intn(len(allCharset))])
	}

	// 打乱字符顺序
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
	//
	//const charset = "abcdefghijklmnopqrstuvwxyz" +
	//	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	//	"!@#$%^&*()"
	//
	//seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	//b := make([]byte, length)
	//for i := range b {
	//	b[i] = charset[seededRand.Intn(len(charset))]
	//}
	//return string(b)
}
