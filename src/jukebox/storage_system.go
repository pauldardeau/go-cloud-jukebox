package jukebox


type StorageSystem interface {
    Enter() bool
    Exit()

    List_containers() []string
    Has_container(container_name string) bool
    Add_container(container_name string)
    Remove_container(container_name string)
    Create_container(container_name string) bool
    Delete_container(container_name string) bool
    List_container_contents(container_name string) []string
    Get_container_names() []string

    Retrieve_file(fm *FileMetadata, local_directory string) int
    Store_file(fm *FileMetadata, file_contents []byte) bool
    Add_file_from_path(container_name string,
                       object_name string,
		       file_path string) bool

    Get_object_metadata(container_name string,
                        object_name string) *map[string]interface{}

    Put_object(container_name string,
               object_name string,
               file_contents []byte,
               headers map[string]string) bool

    Delete_object(container_name string,
                  object_name string) bool

    Get_object(container_name string,
               object_name string,
               local_file_path string) int
}

