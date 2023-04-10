package jukebox

type StorageSystem interface {
	Enter() bool
	Exit()

	HasContainer(containerName string) bool
	CreateContainer(containerName string) bool
	DeleteContainer(containerName string) bool
	ListContainerContents(containerName string) ([]string, error)
	GetContainerNames() ([]string, error)

	RetrieveFile(fm *FileMetadata, localDirectory string) int64
	StoreFile(fm *FileMetadata, fileContents []byte) bool
	AddFileFromPath(containerName string,
		objectName string,
		filePath string) bool

	GetObjectMetadata(containerName string,
		objectName string,
		dictProps *PropertySet) bool

	PutObject(containerName string,
		objectName string,
		fileContents []byte,
		headers *PropertySet) bool

	DeleteObject(containerName string,
		objectName string) bool

	GetObject(containerName string,
		objectName string,
		localFilePath string) int64
}
