package commons

import (
	"io"
	"os"
)

func GetFilesInDirectory(dirPath string, maxRead int) ([]string, error) {
	var (
		filenames []string
	)

	dirp, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dirp.Close()

	if files, err := dirp.Readdir(maxRead); err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, err
	} else {
		for _, file := range files {
			if !file.IsDir() {
				filenames = append(filenames, file.Name())
			}
		}
	}

	return filenames, nil
}
