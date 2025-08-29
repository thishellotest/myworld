package lib

import (
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExt 输入：aaapdf.doc.pdf 输出：aaapdf.doc pdf
func FileExt(filename string, returnLowercase bool) (name string, suffix string) {
	//filename := "example.txt"
	fileNameWithSuffix := filepath.Base(filename)
	extension := filepath.Ext(fileNameWithSuffix)

	name = strings.TrimSuffix(fileNameWithSuffix, extension)
	extension = strings.Replace(extension, ".", "", -1)
	if returnLowercase {
		return name, strings.ToLower(extension)
	}
	return name, extension
}

// TrimHiddenCharacter 隐藏字符的ASCII码为 0-31，127 , 逐个循环字符进行过滤
func TrimHiddenCharacter(originStr string) string {
	srcRunes := []rune(originStr)
	dstRunes := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRunes = append(dstRunes, c)
	}
	return string(dstRunes)
}

func TrimCharacterForFileName(str string) string {
	str = TrimHiddenCharacter(str)
	str = strings.ReplaceAll(str, "\\", "")
	str = strings.ReplaceAll(str, "/", "")
	return str
}
