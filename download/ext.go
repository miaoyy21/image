package download

import (
	"fmt"
	"strings"
)

func parseExt(s string) (string, error) {
	if strings.Contains(s, ".png") {
		return "png", nil
	} else if strings.Contains(s, ".jpeg") {
		return "jpeg", nil
	} else if strings.Contains(s, ".jpg") {
		return "jpg", nil
	}

	return "", fmt.Errorf("unknown file type %s", s)
}
