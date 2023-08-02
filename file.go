package tools

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 作用:判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
func files_exist(path *[]string) (string, bool) {
	for _, f := range *path {
		if FileExist(f) {
			return f, true
		}
	}

	return "", false
}
func FileCopy(dstFileName string, srcFileName string) (written int64, err error) {

	srcFile, err := os.Open(srcFileName)

	if err != nil {
		return 0, err
	}

	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return 0, err
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func FileRemove(name string) error {
	if FileIsDir(name) {
		return os.RemoveAll(name)
	}
	return os.Remove(name)
}
func FileRemove_ext(path, suffix string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// 如果找到的是目录，直接返回
		if info.IsDir() {
			return nil
		}

		// 如果找到的是目标文件
		if strings.HasSuffix(info.Name(), suffix) {
			err = os.Remove(path) // 删除该文件
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func FileIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 下载文件
// 变量说明: 根据url保存为filename
func FileDownload(url, filename string) error {

	res, err := SendGetRequest(url)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err2 := f.Write(res); !(err2 == nil && err2 != io.EOF) {
		return err
	}
	return nil
}

func FileMd5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	defer f.Close()

	f_body, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	h := md5.New()
	h.Write(f_body)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 写入文件
func FileWrite(file string, content []byte) error {

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}
	f.Write(content)
	f.Close()
	return nil
}

// 将以GBK的编码写入文件
func FileWrite_gbk(file string, content []byte) error {

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}
	f.Write(StringSetGBK(content))
	f.Close()
	return nil
}

// 将逐行存储并可选以GBK的编码写入文件
func FileWrite_gbk_list(file string, content []string, gbk bool) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}
	for i, line := range content {
		if i != len(content)-1 {
			line = line + "\r\n"
		}
		if !gbk {
			f.WriteString(line)
		} else {
			f.Write(StringSetGBK([]byte(line)))
		}
	}
	f.Close()
	return nil
}

func FileCreate(name string) (*os.File, error) {
	dir := path.Dir(name)
	if dir != "" {
		_, err := os.Lstat(dir)
		if err != nil {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Create(name)
}

type dirInfo struct {
	Name    string
	ModTime time.Time
}

func FileDecompress_gz(filePath string, targetPath string) error {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	if targetPath != "" {
		folder_mkdir(targetPath)
	}
	currentDir := dirInfo{}
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				if currentDir.Name != "" {
					RemodifyTime(currentDir.Name, currentDir.ModTime)
				}
				break
			} else {
				return err
			}
		}
		fi := header.FileInfo()
		fileName := filepath.Join(targetPath, header.Name)
		if !strings.HasPrefix(fileName, currentDir.Name) {
			RemodifyTime(currentDir.Name, currentDir.ModTime)
		}
		if fi.IsDir() {
			folder_mkdir(fileName)

			currentDir = dirInfo{
				fileName,
				fi.ModTime(),
			}
			continue
		}
		file, err := FileCreate(fileName)
		if err != nil {
			return fmt.Errorf("can not create file %v: %v", fileName, err)
		}
		io.Copy(file, tr)
		file.Close()
		RemodifyTime(fileName, header.ModTime)
	}
	return nil
}
func FileRead(file string) ([]byte, error) {
	if !FileExist(file) {
		return nil, fmt.Errorf("文件不存在:%s", file)
	}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func FileReadGBK(file string) ([]byte, error) {
	if !FileExist(file) {
		return nil, fmt.Errorf("文件不存在:%s", file)
	}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return StringSetGBK(f), nil
}

// 按列读取文件
func FileRead_list(file string) ([]string, int, error) {
	if !FileExist(file) {
		return nil, 0, fmt.Errorf("The file does not exist")
	}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, 0, err
	}
	out := string(StringSetGBK(f))

	list := strings.Split(out, "\r\n")
	return list, len(list), nil
}

// 生成随机的文件名,以tmp结尾
func FileRandom_tmp() string {
	return Rangdom(6) + ".tmp"
}

func folder_mkdir(dir string) {
	os.Mkdir(dir, 0755)
}
