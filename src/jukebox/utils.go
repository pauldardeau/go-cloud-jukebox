package jukebox

import (
   "crypto/md5"
   "errors"
   "fmt"
   "io"
   "os"
   "path/filepath"
)


func FileExists(pathToFile string) bool {
   if _, err := os.Stat(pathToFile); errors.Is(err, os.ErrNotExist) {
      return false
   } else {
      return true
   }
}

func DeleteFile(pathToFile string) bool {
   err := os.Remove(pathToFile)
   return err == nil
}

func DirectoryExists(dirPath string) bool {
   file, err := os.Stat(dirPath)
   if err != nil {
      return false
   }
   return file.IsDir()
}

func CreateDirectory(dirPath string) bool {
   err := os.Mkdir(dirPath, 0755)
   return err == nil
}

func ListDirsInDirectory(dirPath string) ([]string, error) {
   fileList := make([]string, 0)
   files, err := os.ReadDir(dirPath)
   if err != nil {
      return nil, err
   }

   for _, file := range files {
      if file.IsDir() {
         fileList = append(fileList, file.Name())
      }
   }
   return fileList, nil
}

func DirectoryDeleteDirectory(dirPath string) bool {
   err := os.Remove(dirPath)
   return err == nil
}

func ListFilesInDirectory(dirPath string) ([]string, error) {
   fileList := make([]string, 0)
   files, err := os.ReadDir(dirPath)
   if err != nil {
      return nil, err
   }

   for _, file := range files {
      if ! file.IsDir() {
         if file.Name() != "." && file.Name() != ".." {
            fileList = append(fileList, file.Name())
         }
      }
   }
   return fileList, nil
}

func PathJoin(dirPath string, fileName string) string {
   return filepath.Join(dirPath, fileName)
}

func GetFileSize(filePath string) int64 {
   fi, err := os.Stat(filePath)
   if err != nil {
      return -1
   }
   return fi.Size()
}

func FileReadAllText(filePath string) (string, error) {
   //TODO: implement FileReadAllText
   return "", errors.New("function not implemented")
}

func FileWriteAllText(filePath string, fileContents string) bool {
   //TODO: implement FileWriteAllText
   return false
}

func FileWriteAllBytes(filePath string, fileContents []byte) bool {
   //TODO: implement FileWriteAllBytes
   return false
}

func FileReadAllBytes(filePath string) ([]byte, error) {
   //TODO: implement FileReadAllBytes
   return nil, errors.New("function not implemented")
}

func Md5ForFile(pathToFile string) (string, error) {
   f, err := os.Open(pathToFile)
   if err != nil {
      return "", err
   }
   defer f.Close()

   h := md5.New()
   if _, err := io.Copy(h, f); err != nil {
      return "", err
   }
   return fmt.Sprintf("%x", h.Sum(nil)), nil
}

