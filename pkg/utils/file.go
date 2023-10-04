package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

/**
  @author: yhy
  @since: 2022/10/9
  @desc: //TODO
**/

func WriteFile(fileName string, fileData string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//Turn off the FILE
	defer file.Close()
	//When writing to the file, use the cache *writeer
	write := bufio.NewWriter(file)
	write.WriteString(fileData)
	//Flush really writes the cache file into the file
	write.Flush()
}

// LoadFile content to slice
func LoadFile(filename string) (lines []string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println("LoadFile err, ", err)
		return
	}
	defer f.Close() //nolint
	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() != "" {
			lines = append(lines, s.Text())
		}
	}
	return
}

func SaveFile(path string, data []byte) (err error) {
	// Remove file if exist
	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			log.Println("旧文件删除失败", err.Error())
		}
	}

	// save file
	return ioutil.WriteFile(path, data, 0644)
}

// DeCompress 解压
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

// RemoveDir delete ./github All projects
func RemoveDir() {
	Pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, _ := ioutil.ReadDir(Pwd + "/github")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"github", d.Name()}...))
	}
}

// Exists Determine the path of the path given/Whether the folder exists
func Exists(path string) bool {
	_, err := os.Stat(path) //os.STAT obtain file information
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
