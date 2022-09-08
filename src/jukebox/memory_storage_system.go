package jukebox


type MemoryStorageSystem struct {
   AbstractStorageSystem
   container_objects map[string]map[string][]byte
   container_headers map[string]map[string]map[string]string
}


func NewMemoryStorageSystem(debug_mode bool) *MemoryStorageSystem {
   mss := MemoryStorageSystem {
      AbstractStorageSystem{
         Debug_mode: debug_mode,
	 Authenticated: false,
	 Compress_files: false,
	 Encrypt_files: false,
	 List_container_names: []string{},
	 Container_prefix: "",
	 Metadata_prefix: "",
	 Storage_system_type: "Memory",
      },
      make(map[string]map[string][]byte),
      make(map[string]map[string]map[string]string),
   }
   return &mss
}

func (mss *MemoryStorageSystem) Enter() bool {
   return true
}

func (mss *MemoryStorageSystem) Exit() {
}

func (mss *MemoryStorageSystem) List_account_containers() []string {
   return mss.List_container_names
}

func (mss *MemoryStorageSystem) Create_container(container_name string) bool {
   container_created := false
   mss.container_objects[container_name] = make(map[string][]byte)
   mss.container_headers[container_name] = make(map[string]map[string]string)
   if ! mss.Has_container(container_name) {
      mss.List_container_names = append(mss.List_container_names, container_name)
      container_created = true
   }
   return container_created
}

func (mss *MemoryStorageSystem) Delete_container(container_name string) bool {
    container_deleted := false
    if mss.Has_container(container_name) {
        //TODO:
        //del mss.container_objects[container_name]
        //del mss.container_headers[container_name]
        //mss.list_containers.remove(container_name)
        container_deleted = true
    }
    return container_deleted
}

func (mss *MemoryStorageSystem) List_container_contents(container_name string) []string {
    list_contents := []string{}
    object_container, exists := mss.container_objects[container_name]
    if exists {
        for key, _ := range object_container {
            list_contents = append(list_contents, key)
        }
    }
    return list_contents
}

func (mss *MemoryStorageSystem) Get_object_metadata(container_name string,
                                                    object_name string) *map[string]string {
    if len(container_name) > 0 && len(object_name) > 0 {
        header_container, exists := mss.container_headers[container_name]
        if exists {
            headers := header_container[object_name]
            return &headers
        }
    }
    return nil
}

func (mss *MemoryStorageSystem) Put_object(container_name string,
                                           object_name string,
                                           file_contents []byte,
                                           headers map[string]string) bool {
   object_added := false
   if len(container_name) > 0 &&
      len(object_name) > 0 &&
      len(file_contents) > 0 {

      object_container := mss.container_objects[container_name]
      object_container[object_name] = file_contents
      header_container := mss.container_headers[container_name]
      header_container[object_name] = headers
      object_added = true
   }
   return object_added
}

func (mss *MemoryStorageSystem) Delete_object(container_name string,
                                              object_name string) bool {
    object_deleted := false
    if len(container_name) > 0 && len(object_name) > 0 {
        if mss.Has_container(container_name) {
            object_container := mss.container_objects[container_name]
	    _, object_exists_in_container := object_container[object_name]
            if object_exists_in_container {
                //TODO: delete object
                //del object_container[object_name]
                object_deleted = true
            }
	    header_container := mss.container_headers[container_name]
	    _, object_exists_in_headers := header_container[object_name]
            if object_exists_in_headers {
                //TODO: delete from headers
                //del header_container[object_name]
                object_deleted = true
            }
        }
    }
    return object_deleted
}

func (mss *MemoryStorageSystem) Get_object(container_name string,
                                           object_name string,
                                           local_file_path string) int {
   bytes_retrieved := 0
   if len(container_name) > 0 && len(object_name) > 0 && len(local_file_path) > 0 {
      object_container, container_exists := mss.container_objects[container_name]
      if container_exists {
	 _, object_exists_in_container := object_container[object_name]
         if object_exists_in_container {
            //TODO: write object to file
            //with open(local_file_path, 'w') as f:
            //   f.write(object_container[object_name])
            bytes_retrieved = len(object_container[object_name])
         }
      }
   }
   return bytes_retrieved
}
