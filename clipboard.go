package clip_img

import (
	"io"
	"os"
)

// Write write image to clipboard
func Write(r io.Reader) error {
	f, err := writeTemp(r)
	if err != nil {
		return err
	}
	defer os.Remove(f)
	return write(f)
}

// Read read image from clipboard
func Read() (io.Reader, error) {
	return read()
}

func writeTemp(r io.Reader) (string, error) {
	f, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	return f.Name(), nil
}
