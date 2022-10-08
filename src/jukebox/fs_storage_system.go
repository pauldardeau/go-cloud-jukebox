package jukebox

import (
   "errors"
   "fmt"
)

type FSStorageSystem struct {
   rootDir string
   debugMode bool
}

func NewFSStorageSystem(rootDir string, debugMode bool) *FSStorageSystem {
    return &FSStorageSystem{
        rootDir: rootDir,
        debugMode: debugMode,
    }
}

func (fs *FSStorageSystem) Enter() bool {
   if !DirectoryExists(fs.rootDir) {
      CreateDirectory(fs.rootDir)
   }
   return DirectoryExists(fs.rootDir)
}

func (fs *FSStorageSystem) Exit() {
}

func (fs *FSStorageSystem) ListAccountContainers() ([]string, error) {
   return ListDirsInDirectory(fs.rootDir)
}

func (fs *FSStorageSystem) GetContainerNames() ([]string, error) {
   return fs.ListAccountContainers()
}

func (fs *FSStorageSystem) HasContainer(containerName string) bool {
   listContainers, err := fs.ListAccountContainers()
   if err != nil {
      return false
   } else {
      for _, container := range listContainers {
         if containerName == container {
            return true
         }
      }
      return false
   }
}

func (fs *FSStorageSystem) CreateContainer(containerName string) bool {
   containerDir := PathJoin(fs.rootDir, containerName)
   containerCreated := CreateDirectory(containerDir)
   if containerCreated {
      if fs.debugMode {
         fmt.Printf("container created: '%s'\n", containerName)
      }
   }
   return containerCreated
}

func (fs *FSStorageSystem) DeleteContainer(containerName string) bool {
   containerDir := PathJoin(fs.rootDir, containerName)
   containerDeleted := DirectoryDeleteDirectory(containerDir)
   if containerDeleted {
      if fs.debugMode {
         fmt.Printf("container deleted: '%s'\n", containerName)
      }
   }
   return containerDeleted
}

func (fs *FSStorageSystem) ListContainerContents(containerName string) ([]string, error) {
   containerDir := PathJoin(fs.rootDir, containerName)
   if DirectoryExists(containerDir) {
      return ListFilesInDirectory(containerDir)
   } else {
      return nil, errors.New("container does not exist")
   }
}

func (fs *FSStorageSystem) GetObjectMetadata(containerName string,
                                             objectName string,
                                             dictProps *PropertySet) bool {
   if len(containerName) > 0 && len(objectName) > 0 {
      containerDir := PathJoin(fs.rootDir, containerName)
      if DirectoryExists(containerDir) {
         objectPath := PathJoin(containerDir, objectName)
         metaPath := objectPath + ".meta"
         if FileExists(metaPath) {
            return dictProps.ReadFromFile(metaPath)
         }
      }
   }
   return false
}

func (fs *FSStorageSystem) PutObject(containerName string,
                                     objectName string,
                                     fileContents []byte,
                                     headers *PropertySet) bool {
   objectAdded := false
   if len(containerName) > 0 && len(objectName) > 0 && len(fileContents) > 0 {
      containerDir := PathJoin(fs.rootDir, containerName)
      if DirectoryExists(containerDir) {
         objectPath := PathJoin(containerDir, objectName)
         objectAdded = FileWriteAllBytes(objectPath, fileContents)
         if objectAdded {
            if fs.debugMode {
               fmt.Printf("object added: %s/%s\n", containerName, objectName)
            }
            if headers != nil {
               if headers.Count() > 0 {
                  metaPath := objectPath + ".meta"
                  headers.WriteToFile(metaPath)
               }
            }
         } else {
            fmt.Println("FileWriteAllBytes failed to write object contents, put failed")
         }
      } else {
         if fs.debugMode {
            fmt.Println("container doesn't exist, can't put object")
         }
      }
   } else {
      if fs.debugMode {
         if len(containerName) == 0 {
            fmt.Println("container name is missing, can't put object")
         } else {
            if len(objectName) == 0 {
               fmt.Println("object name is missing, can't put object")
            } else {
               if len(fileContents) == 0 {
                  fmt.Println("object content is empty, can't put object")
               }
            }
         }
      }
   }
   return objectAdded
}

func (fs *FSStorageSystem) DeleteObject(containerName string,
                                        objectName string) bool {
   objectDeleted := false
   if len(containerName) > 0 && len(objectName) > 0 {
      containerDir := PathJoin(fs.rootDir, containerName)
      objectPath := PathJoin(containerDir, objectName)
      if FileExists(objectPath) {
         objectDeleted = DeleteFile(objectPath)
         if objectDeleted {
            if fs.debugMode {
               fmt.Printf("object deleted: %s/%s\n", containerName, objectName)
            }
            metaPath := objectPath + ".meta"
            if FileExists(metaPath) {
               DeleteFile(metaPath)
            }
         } else {
            if fs.debugMode {
               fmt.Println("delete of object file failed")
            }
         }
      } else {
         if fs.debugMode {
            fmt.Println("cannot delete object, path doesn't exist")
         }
      }
   } else {
      if fs.debugMode {
         fmt.Println("cannot delete object, container name or object name is missing")
      }
   }
   return objectDeleted
}

func (fs *FSStorageSystem) GetObject(containerName string,
                                     objectName string,
                                     localFilePath string) int64 {
   var bytesRetrieved int64
   if len(containerName) > 0 &&
      len(objectName) > 0 &&
      len(localFilePath) > 0 {

      containerDir := PathJoin(fs.rootDir, containerName)
      objectPath := PathJoin(containerDir, objectName)
      if FileExists(objectPath) {
         objFileContents, err := FileReadAllBytes(objectPath)
         if err == nil {
            if fs.debugMode {
               fmt.Printf("attempting to write object to '%s'\n", localFilePath)
            }
            if FileWriteAllBytes(localFilePath, objFileContents) {
               bytesRetrieved = int64(len(objFileContents))
            }
         } else {
            fmt.Printf("error: unable to read object file '%s'\n", objectPath)
            fmt.Printf("error: %v\n", err)
         }
      }
   }
   return bytesRetrieved
}

