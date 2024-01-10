package image

import (
	"fmt"

	"github.com/jykuo-love-shiritori/twp/pkg/constants"
)

func GetUrl(filename string) string {
	if filename == "" {
		// ! Very bad
		return "/placeholder/parrot.gif"
	}

	return fmt.Sprintf("%s/%s", constants.IMAGE_BASE_PATH, filename)
}
