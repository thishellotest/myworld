package a

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_aac(t *testing.T) {

	regex := regexp.MustCompile(`^(\(\d{3}\) \d{3}-\d{4}|\+\d{1,3} \d{3}-\d{4}-\d{3}|\d{10}|\d{3}-\d{3}-\d{4})$`)

	testCases := []string{
		"(123) 123-1234",
		"+49 123-1234-123",
		"1238923741",
		"123-123-123",
		"1234-567-890", // 无效示例
	}

	for _, testCase := range testCases {
		if regex.MatchString(testCase) {
			fmt.Printf("Valid: %s\n", testCase)
		} else {
			fmt.Printf("Invalid: %s\n", testCase)
		}
	}

}
