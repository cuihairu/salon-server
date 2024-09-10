package utils

import (
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
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

func DeleteFile(dirname string, filename string) error {
	if !filepath.IsAbs(dirname) {
		return fmt.Errorf("invalid path: %s", dirname)
	}
	filename = filepath.Join(dirname, filename)
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// 文件不存在，返回 nil
		return nil
	} else if err != nil {
		// 其他错误
		return err
	}
	// 检查是否有删除权限（即检查文件是否可写）
	if info.Mode().Perm()&0200 == 0 {
		// 文件不可写，返回权限错误
		return fmt.Errorf("no write permission for file: %s", filename)
	}
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(uploadFile multipart.File, dirname string, filename string) error {
	if !filepath.IsAbs(dirname) {
		return fmt.Errorf("invalid path: %s", dirname)
	}
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		// 目录不存在，尝试创建目录
		err = os.MkdirAll(dirname, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fullFilename := filepath.Join(dirname, filename)
	out, err := os.Create(fullFilename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, uploadFile)
	if err != nil {
		return err
	}
	return nil
}

type FileWithPermissions struct {
	Filename string `json:"filename" yaml:"filename"`
	Perm     string `json:"perm" yaml:"perm"`
}

func ListFilesWithPermissions(dir string) ([]FileWithPermissions, error) {
	var files []FileWithPermissions
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			fileMode := info.Mode()
			perm := fileMode.Perm()
			files = append(files, FileWithPermissions{Filename: relativePath, Perm: fmt.Sprintf("%o", perm)})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ListFilesWithPermissionsAndPaging(dir string, current int, pageSize int) ([]FileWithPermissions, int64, error) {
	var files []FileWithPermissions
	if current <= 0 {
		current = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var start = int64((current - 1) * pageSize)
	var end = start + int64(pageSize)
	var index int64
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if index >= start && index < end {
				relativePath, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}
				fileMode := info.Mode()
				perm := fileMode.Perm()
				files = append(files, FileWithPermissions{Filename: relativePath, Perm: fmt.Sprintf("%o", perm)})
			}
			index++
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return files, index, nil
}
