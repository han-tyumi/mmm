package download

import (
	"io"
	"net/http"
	"os"

	"github.com/han-tyumi/mmm/utils"
)

// FromURL downloads a file from a URL to the current directory under a name.
func FromURL(name, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		utils.Error(res.Status)
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}
