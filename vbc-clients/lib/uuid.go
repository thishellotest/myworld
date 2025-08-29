package lib

import (
	"fmt"
	guuid "github.com/google/uuid"
	"time"
)

// UuidNumeric length 10
func UuidNumeric() string {
	uuid := guuid.New()
	id := InterfaceToString(uuid.ID())
	idLen := len(id)
	length := 10
	if idLen < length {
		id = fmt.Sprintf("%0*s", length, id)
	}
	id = time.Now().Format("06") + id[2:10]
	return id
}

//
//func UuidNumericTime() string {
//	str := UuidNumeric()
//	now := time.Now()
//	// 2006-01-02 15:04:05
//	n, _ := rand.Int(rand.Reader, big.NewInt(100))
//	if n == nil {
//		n = big.NewInt(0)
//	}
//	randStr := InterfaceToString(n.Int64())
//	randStr = fmt.Sprintf("%0*s", 3, randStr)
//	//time.RFC3339
//	str = now.Format("05") + str[5:] + randStr
//	return str
//}
