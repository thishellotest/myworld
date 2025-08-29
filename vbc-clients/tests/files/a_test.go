package files

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func Test_a(t *testing.T) {

	filePath := "/path/to/your/file.ss.ccc.txt"
	fileNameWithSuffix := filepath.Base(filePath)
	fileName := strings.TrimSuffix(fileNameWithSuffix, filepath.Ext(fileNameWithSuffix))
	fmt.Println("File name without suffix:", fileName)

}
