package jukebox


type StorageSystem interface {
    Enter() bool
    Exit()

    ListContainers() []string
    HasContainer(containerName string) bool
    AddContainer(containerName string)
    RemoveContainer(containerName string)
    CreateContainer(containerName string) bool
    DeleteContainer(containerName string) bool
    ListContainerContents(containerName string) []string
    GetContainerNames() []string

    RetrieveFile(fm *FileMetadata, localDirectory string) int
    StoreFile(fm *FileMetadata, fileContents []byte) bool
    AddFileFromPath(containerName string,
                    objectName string,
                    filePath string) bool

    GetObjectMetadata(containerName string,
                      objectName string) *map[string]interface{}

    PutObject(containerName string,
              objectName string,
              fileContents []byte,
              headers map[string]string) bool

    DeleteObject(containerName string,
                 objectName string) bool

    GetObject(containerName string,
              objectName string,
              localFilePath string) int
}

