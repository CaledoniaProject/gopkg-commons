package commons

import "os"

func NormalFileExists(path string) (bool, error) {
	if info, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	} else {
		return info.Mode().IsRegular(), nil
	}
}

func FileSize(filename string) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func FileSize2(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
