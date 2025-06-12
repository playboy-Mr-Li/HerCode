package readfile

import (
	"fmt"
	"os"
	"strings"
)

func ReadFile(fileName string) ([]byte, error) {

	if !strings.HasSuffix(fileName, ".hc") {
		return nil, fmt.Errorf("not a hercode file, must be .hc")

	}

	return os.ReadFile(fileName)
}
