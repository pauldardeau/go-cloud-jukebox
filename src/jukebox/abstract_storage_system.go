package jukebox

import (
   "fmt"
   "strings"
)


type AbstractStorageSystem struct {
   Debug_mode bool
   Authenticated bool
   Compress_files bool
   Encrypt_files bool
   List_container_names []string
   Container_prefix string
   Metadata_prefix string
   Storage_system_type string
   Css StorageSystem
}

func NewAbstractStorageSystem(storage_system_type string,
                              debug_mode bool,
                              css StorageSystem) *AbstractStorageSystem {
   var ss AbstractStorageSystem
   ss.Debug_mode = debug_mode
   ss.Authenticated = false
   ss.Compress_files = false
   ss.Encrypt_files = false
   ss.List_container_names = []string{}
   ss.Container_prefix = ""
   ss.Metadata_prefix = ""
   ss.Storage_system_type = storage_system_type
   ss.Css = css
   return &ss
}

func (ss *AbstractStorageSystem) Un_prefixed_container(container_name string) string {
   if len(ss.Container_prefix) > 0 && len(container_name) > 0 {
      if strings.HasPrefix(container_name, ss.Container_prefix) {
         return container_name[len(ss.Container_prefix):]
      }
   }
   return container_name
}

func (ss *AbstractStorageSystem) Prefixed_container(container_name string) string {
   return ss.Container_prefix + container_name
}

func (ss *AbstractStorageSystem) Has_container(container_name string) bool {
   container_name_in_list_containers := false
   if len(ss.List_container_names) > 0 {
      for _, cnr_name := range ss.List_container_names {
         if cnr_name == container_name {
            container_name_in_list_containers = true
            break
         }
      }
   }
   return container_name_in_list_containers
}

func (ss *AbstractStorageSystem) Get_container_names() []string {
   return ss.List_container_names
}

func (ss *AbstractStorageSystem) Add_container(container_name string) {
   ss.List_container_names = append(ss.List_container_names, container_name)
}

func (ss *AbstractStorageSystem) Remove_container(container_name string) {
   //TODO: implement remove_container
   //ss.List_containers.Remove(container_name)
}

func (ss *AbstractStorageSystem) Retrieve_file(fm *FileMetadata,
                                               local_directory string) int {
   if fm != nil && len(local_directory) > 0 {
      file_path := Path_join(local_directory, fm.File_uid)
      fmt.Printf("retrieving container=%s\n", fm.Container_name)
      fmt.Printf("retrieving object=%s\n", fm.Object_name)
      return ss.Css.Get_object(fm.Container_name,
                               fm.Object_name,
                               file_path)
   }
   return 0
}

func (ss *AbstractStorageSystem) store_file(fm *FileMetadata,
                                            file_contents []byte) bool {
   if fm != nil && file_contents != nil {
      return ss.Css.Put_object(fm.Container_name,
                               fm.Object_name,
                               file_contents,
                               fm.To_Dictionary_With_Prefix(ss.Metadata_prefix))
   }
   return false
}

func (ss *AbstractStorageSystem) Add_file_from_path(container_name string,
                                                    object_name string,
                                                    file_path string) bool {
   return false
}

