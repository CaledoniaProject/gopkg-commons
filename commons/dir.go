package commons

import (
	"io"
	"os"
	"path/filepath"
)

func RemoveFolderContents(dirname string) error {
	var (
		batchSize = 4096
	)

	dirp, err := os.Open(dirname)
	if err != nil {
		return err
	}
	defer dirp.Close()

	for {
		files, err := dirp.Readdir(batchSize)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for _, file := range files {
			if err := os.RemoveAll(filepath.Join(dirname, file.Name())); err != nil {
				return err
			}
		}

		if len(files) < batchSize {
			break
		}
	}

	return nil
}

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
