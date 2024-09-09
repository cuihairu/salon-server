package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CreateDirIfNotExist(paths ...string) error {
	for _, path := range paths {
		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateFileIfNotExistInCurPath(dirname string, createFile string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return "", err
	}

	// 构建 logs 目录路径
	logsDir := filepath.Join(currentDir, "logs")

	// 检查 logs 目录是否存在
	_, err = os.Stat(logsDir)
	if os.IsNotExist(err) {
		// 目录不存在，尝试创建 logs 目录
		err = os.Mkdir(logsDir, 0755)
		if err != nil {
			fmt.Println("Error creating logs directory:", err)
			return "", err
		}
		fmt.Println("Logs directory created at:", logsDir)
	} else if err != nil {
		// 发生其他错误（例如权限问题）
		fmt.Println("Error checking logs directory:", err)
		return "", err
	} else {
		fmt.Println("Logs directory already exists:", logsDir)
	}
	if len(createFile) > 0 {
		// 检查目录写入权限
		testFile := filepath.Join(logsDir, createFile)
		file, err := os.OpenFile(testFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("No write permission to logs directory or file exists:", err)
			return "", err
		}
		file.Close()
		return testFile, nil
	}
	return currentDir, nil
}

func ListFiles(dir string) ([]string, error) {
	var files []string
	// type WalkFunc func(path string, info fs.FileInfo, err error) error
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files = append(files, relativePath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
