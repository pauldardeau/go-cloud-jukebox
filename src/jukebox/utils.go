package jukebox

import (
   "crypto/md5"
   "fmt"
   "io"
   "os"
   "path/filepath"
   "strings"
   "time"
)


func FileExists(pathToFile string) bool {
   file, err := os.Stat(pathToFile)
   if err != nil {
      return false
   }
   return !file.IsDir()
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

func PathSplitExt(path string) (string, string) {
   // python: os.path.splitext

   // splitext("bar") -> ("bar", "")
   // splitext("foo.bar.exe") -> ("foo.bar", ".exe")
   // splitext("/foo/bar.exe") -> ("/foo/bar", ".exe")
   // splitext(".cshrc") -> (".cshrc", "")
   // splitext("/foo/....jpg") -> ("/foo/....jpg", "")

   root := ""
   ext := ""

   if len(path) > 0 {
      pos_last_dot := strings.LastIndex(path, ".")
      if pos_last_dot == -1 {
         // no '.' exists in path (i.e., "bar")
         root = path
      } else {
         // is the last '.' the first character? (i.e., ".cshrc")
         if pos_last_dot == 0 {
            root = path
         } else {
            preceding := path[pos_last_dot-1]
            // is the preceding char also a '.'? (i.e., "/foo/....jpg")
            if preceding == '.' {
               root = path
            } else {
               // splitext("foo.bar.exe") -> ("foo.bar", ".exe")
               // splitext("/foo/bar.exe") -> ("/foo/bar", ".exe")
	       root = path[0:pos_last_dot]
	       ext = path[pos_last_dot:]
            }
         }
      }
   }

   return root, ext
}

func PathGetMtime(path string) (time.Time, error) {
   fi, err := os.Stat(path)
   if err != nil {
      return time.Now(), err
   } else {
      return fi.ModTime(), nil
   }
}

func GetFileSize(filePath string) int64 {
   fi, err := os.Stat(filePath)
   if err != nil {
      return -1
   }
   return fi.Size()
}

func FileReadAllText(filePath string) (string, error) {
   content, err := os.ReadFile(filePath)
   if err != nil {
      fmt.Printf("error: unable to read file '%s': %v\n", filePath, err)
      return "", err
   } else {
      return string(content), nil
   }
}

func FileWriteAllText(filePath string, fileContents string) bool {
   f, err := os.Create(filePath)
   if err != nil {
      fmt.Printf("error: unable to create file '%s': %v\n", filePath, err)
      return false
   }
   defer f.Close()
 
   f.Write([]byte(fileContents))
   return true
}

func FileWriteAllBytes(filePath string, fileContents []byte) bool {
   f, err := os.Create(filePath)
   if err != nil {
      fmt.Printf("error: unable to create file '%s': %v\n", filePath, err)
      return false
   }
   defer f.Close()

   f.Write(fileContents)
   return true
}

func FileReadAllBytes(filePath string) ([]byte, error) {
   content, err := os.ReadFile(filePath)
   if err != nil {
      fmt.Printf("error: unable to read file '%s': %v\n", filePath, err)
      return nil, err
   } else {
      return content, nil
   }
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

func TimeSleepSeconds(seconds int) {
   time.Sleep(time.Duration(seconds) * time.Second)
}

