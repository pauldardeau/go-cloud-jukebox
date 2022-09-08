package jukebox

import (
   //"context"
   "errors"
   "fmt"
   "github.com/aws/aws-sdk-go/aws"
   //"github.com/aws/aws-sdk-go-v2/config"
   "github.com/aws/aws-sdk-go/aws/session"
   //"github.com/aws/aws-sdk-go/service/s3"
   //"github.com/aws/aws-sdk-go-v2/service/s3"
   //"github.com/aws/aws-sdk-go-v2/credentials"
   "github.com/aws/aws-sdk-go/service/s3"
   "github.com/aws/aws-sdk-go/credentials"
)

// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#pkg-examples
// https://github.com/mathisve/golang-AWS-S3-SDK-Example/blob/master/main.go

type S3StorageSystem struct {
   AbstractStorageSystem
   aws_access_key string
   aws_secret_key string
   container_prefix string
   s3_session *s3.Client
}


func NewS3StorageSystem(access_key string,
                        secret_key string,
                        cnr_prefix string,
                        debug_mode bool) StorageSystem {

   s3_ss := S3StorageSystem {
      AbstractStorageSystem{
         Debug_mode: debug_mode,
         Authenticated: false,
         Compress_files: false,
         Encrypt_files: false,
         List_container_names: []string{},
         Container_prefix: "",
         Metadata_prefix: "",
         Storage_system_type: "S3",
      },
      access_key,
      secret_key,
      cnr_prefix,
      nil,
    }

    if debug_mode {
        fmt.Printf("Using access_key='%s', secret_key='%s'\n", access_key, secret_key)
    }

    if len(cnr_prefix) > 0 {
        if debug_mode {
            fmt.Printf("using container_prefix='%s'\n", cnr_prefix)
        }
    }

    return s3_ss
}

func (s3_ss S3StorageSystem) Enter() bool {
    connected := false
    if s3_ss.Debug_mode {
        fmt.Println("attempting to connect to S3")
    }

    s3Config := aws.Config{
        Credentials:      credentials.NewStaticCredentials(s3_ss.aws_access_key, s3_ss.aws_secret_key, ""),
        Endpoint:         aws.String("https://s3.wasabisys.com"),
        Region:           aws.String("us-central-1"),
        S3ForcePathStyle: aws.Bool(true),
    }

    session, err := session.NewSessionWithOptions(session.Options{
        Config: s3Config,
    })
    if err != nil {
        fmt.Println("unable to establish S3 session")
    } else {
        s3Client := s3.New(session)

        fmt.Println("established S3 session")
        s3_ss.s3_session = s3Client
        connected = true
    }

/*
    creds := credentials.NewStaticCredentialsProvider(s3_ss.aws_access_key,
                                                      s3_ss.aws_secret_key,
						      "")
    endpoint := "https://s3.us-central-1.wasabisys.com"

    //config := aws.Config{
    //    Credentials: creds,
    //    Endpoint: &endpoint,
    //}

    cfg, err := config.LoadDefaultConfig(context.TODO(),
                                         config.WithCredentialsProvider(creds),
                                         config.WithEndpoint(endpoint))
    if err != nil {
        fmt.Printf("error: %v\n", err)
        return
    }

    awsS3Client := s3.NewFromConfig(cfg)
    */
/*
    sess, err := session.NewSession(&config)
    if err != nil {
        fmt.Println("unable to establish S3 session")
    } else {
        fmt.Println("established S3 session")
        s3_ss.s3_session = sess
	connected = true
    }
    */

    if connected {
        s3_ss.Authenticated = true
        list_buckets, err := s3_ss.List_account_containers()
        if err == nil {
            s3_ss.List_container_names = *list_buckets
        }
    }

    return connected
}

func (s3_ss S3StorageSystem) Exit() {
    if s3_ss.s3_session != nil {
        if s3_ss.Debug_mode {
            fmt.Println("closing S3 connection object")
        }

        s3_ss.Authenticated = false
        s3_ss.List_container_names = nil
        s3_ss.s3_session = nil
    }
}

func (s3_ss *S3StorageSystem) List_account_containers() (*[]string, error) {
    if s3_ss.Debug_mode {
        fmt.Println("list_account_containers")
    }

    if s3_ss.s3_session != nil {
	    /*
        input := &s3.ListBucketsInput{}

        result, err := s3_ss.s3_session.ListBuckets(input)
        if err != nil {
            if aerr, ok := err.(awserr.Error); ok {
                switch aerr.Code() {
                    default:
                        fmt.Println(aerr.Error())
                }
            } else {
                // Print the error, cast err to awserr.Error to get the Code and
                // Message from an error.
                fmt.Println(err.Error())
            }
	    err := errors.New("ListBuckets not implemented")
	    return nil, err
        } else {
            list_container_names := []string{}
            for _, bucket := range result.Buckets {
                container_name := *bucket.Name
                list_container_names = append(list_container_names,
                                              s3_ss.Un_prefixed_container(container_name))
            }

            return &list_container_names, nil
        }
	*/
	return nil, errors.New("ListBuckets not implemented")
    }

    err := errors.New("No existing S3 session")

    return nil, err
}

func (s3_ss S3StorageSystem) Create_container(container_name string) bool {
    if s3_ss.Debug_mode {
        fmt.Printf("create_container: '%s'\n", container_name)
    }

    container_created := false

    if s3_ss.s3_session != nil {
       /*
       resp, err := s3_ss.s3_session.CreateBucket(&s3.CreateBucketInput{
          Bucket: aws.String(BUCKET_NAME),
          CreateBucketConfiguration: &s3.CreateBucketConfiguration{
             LocationConstraint: aws.String(REGION),
          },
       })

       if err != nil {
          if aerr, ok := err.(awserr.Error); ok {
             switch aerr.Code() {
            case s3.ErrCodeBucketAlreadyExists:
        fmt.Println("Bucket name already in use!")
        panic(err)
      case s3.ErrCodeBucketAlreadyOwnedByYou:
        fmt.Println("Bucket exists and is owned by you!")
      default:
        panic(err)
           }
       }
       */
    }

    if container_created {
       s3_ss.Add_container(container_name)
    }

    return container_created
}

func (s3_ss S3StorageSystem) Delete_container(container_name string) bool {
    if s3_ss.Debug_mode {
        fmt.Printf("delete_container: '%s'\n", container_name)
    }

    container_deleted := false

    if s3_ss.s3_session != nil {

	    /*
input := &s3.DeleteBucketInput{
    Bucket: aws.String(container_name),
}

result, err := svc.DeleteBucket(input)
if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        // Print the error, cast err to awserr.Error to get the Code and
        // Message from an error.
        fmt.Println(err.Error())
    }
}
*/
            container_deleted = true
    }

    if container_deleted {
       s3_ss.Remove_container(container_name)
    }

    return container_deleted
}

func (s3_ss S3StorageSystem) List_container_contents(container_name string) []string {
    if s3_ss.Debug_mode {
        fmt.Printf("list_container_contents: '%s'\n", container_name)
    }

    if s3_ss.s3_session != nil {
        //try:
        /*
            response = s3.conn.list_objects_v2(container_name)
            meta = response['ResponseMetadata']
            status_code = meta['HTTPStatusCode']
            if status_code == 200:
                list_contents = []
                contents = response['Contents']

                for objDict in contents:
                    list_contents.append(objDict['Key'])

                return list_contents
        */
        //except Exception as exception:
        //    print("exception caught: %s" % type(exception).__name__)
        // except S3.Client.exceptions.NoSuchBucket:
        //    print("bucket does not exist")
        //    pass
    }

    return nil
}

func (s3_ss S3StorageSystem) Get_object_metadata(container_name string,
                                                 object_name string) *map[string]interface{} {
    if s3_ss.Debug_mode {
        fmt.Printf("get_object_metadata: container='%s', object='%s'\n",
                   container_name,
                   object_name)
    }

    // https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.HeadObject
    if s3_ss.s3_session != nil && len(container_name) > 0 && len(object_name) > 0 {
	    /*
        request := s3.HeadObjectInput{
	   Bucket: &container_name,
	   Key: &object_name,
	}
        response, err := s3_ss.s3_session.HeadObject(&request)
	if err != nil {
           return nil
        } else {
           dict_headers := make(map[string]interface{})
	   //if response.ArchiveStatus != nil {
	   //   dict_headers["ArchiveStatus"] = *response.ArchiveStatus
           //}
	   //if response.BucketKeyEnabled != nil {
           //   enabled := *response.BucketKeyEnabled
	   //   var enabled_value string
	   //   if enabled {
           //      enabled_value = "true"
           //   } else {
           //      enabled_value = "false"
           //   }
	   //   dict_headers["BucketKeyEnabled"] = enabled_value
           //}
	   if response.CacheControl != nil {
              dict_headers["CacheControl"] = *response.CacheControl
           }
	   if response.ChecksumCRC32 != nil {
              dict_headers["ChecksumCRC32"] = *response.ChecksumCRC32
           }
	   if response.ChecksumCRC32C != nil {
              dict_headers["ChecksumCRC32C"] = *response.ChecksumCRC32C
           }
	   if response.ChecksumSHA1 != nil {
              dict_headers["ChecksumSHA1"] = *response.ChecksumSHA1
           }
	   if response.ChecksumSHA256 != nil {
              dict_headers["ChecksumSHA256"] = *response.ChecksumSHA256
           }
	   if response.ContentDisposition != nil {
              dict_headers["ContentDisposition"] = *response.ContentDisposition
           }
	   if response.ContentEncoding != nil {
              dict_headers["ContentEncoding"] = *response.ContentEncoding
           }
	   if response.ContentLanguage != nil {
              dict_headers["ContentLanguage"] = *response.ContentLanguage
           }
	   //if response.ContentLength != nil {
           //   //*int64
           //}
	   if response.ContentType != nil {
              dict_headers["ContentType"] = *response.ContentType
           }
	   //if response.DeleteMarker != nil {
           //   delete_marker := *response.DeleteMarker
           //   var delete_marker_value string
           //   if delete_marker {
           //      delete_marker_value = "true"
           //   } else {
           //      delete_marker_value = "false"
           //   }
	   //   dict_headers["DeleteMarker"] = delete_marker_value
           //}
	   if response.ETag != nil {
              dict_headers["ETag"] = *response.ETag
           }
	   if response.Expiration != nil {
              dict_headers["Expiration"] = *response.Expiration
           }
	   if response.Expires != nil {
              dict_headers["Expires"] = *response.Expires
           }
	   if response.LastModified != nil {
              //*time.Time
           }
	   if response.Metadata != nil {
              //map[string]*string
           }
	   //if response.MissingMeta != nil {
           //   //*int64
           //}
	   //if response.ObjectLockLegalHoldStatus != nil {
           //   dict_headers["ObjectLockLegalHoldStatus"] = *response.ObjectLockLegalHoldStatus
           //}
	   //if response.ObjectLockMode != nil {
           //   dict_headers["ObjectLockMode"] = *response.ObjectLockMode
           //}
	   if response.ObjectLockRetainUntilDate != nil {
              //*time.Time
           }
	   //if response.PartsCount != nil {
           //   dict_headers["PartsCount"] = fmt.Sprintf("%d", *response.PartsCount)
           //}
	   //if response.ReplicationStatus != nil {
           //   dict_headers["ReplicationStatus"] = *response.ReplicationStatus
           //}
	   //if response.RequestCharged != nil {
           //   dict_headers["RequestCharged"] = *response.RequestCharged
           //}
	   if response.Restore != nil {
              dict_headers["Restore"] = *response.Restore
           }
	   if response.SSECustomerAlgorithm != nil {
              dict_headers["SSECustomerAlgorithm"] = *response.SSECustomerAlgorithm
           }
	   if response.SSECustomerKeyMD5 != nil {
              dict_headers["SSECustomerKeyMD5"] = *response.SSECustomerKeyMD5
           }
	   if response.SSEKMSKeyId != nil {
              dict_headers["SSEKMSKeyId"] = *response.SSEKMSKeyId
           }
	   //if response.ServerSideEncryption != nil {
           //   dict_headers["ServerSideEncryption"] = *response.ServerSideEncryption
           //}
	   //if response.StorageClass != nil {
           //   dict_headers["StorageClass"] = *response.StorageClass
           //}
	   if response.VersionId != nil {
              dict_headers["VersionId"] = *response.VersionId
           }
	   if response.WebsiteRedirectLocation != nil {
              dict_headers["WebsiteRedirectLocation"] = *response.WebsiteRedirectLocation
           }
	   return &dict_headers
        }
	*/
    }
    return nil
}

func (s3_ss S3StorageSystem) Put_object(container_name string,
                                        object_name string,
                                        file_contents []byte,
                                        headers map[string]string) bool {

    object_added := false

    if s3_ss.s3_session != nil && len(container_name) > 0 &&
       len(object_name) > 0 && len(file_contents) > 0 {

        /*
        bucket := container_name
        result = s3.conn.put_object(file_contents, bucket, object_name)
        if "HTTPStatusCode" in result {
            status_code = result["HTTPStatusCode"]
            if status_code == 200 {
                object_added = true
            }
        } else {
            if "ResponseMetadata" in result {
                resp_meta = result["ResponseMetadata"]
                if "HTTPStatusCode" in resp_meta {
                    status_code = resp_meta["HTTPStatusCode"]
                    if status_code == 200 {
                        object_added = true
                    }
                }
            } else {
                //print(repr(result))
            }
        }
        */
    }

    return object_added
}

func (s3_ss S3StorageSystem) Delete_object(container_name string, object_name string) bool {
    if s3_ss.Debug_mode {
        fmt.Printf("delete_object: container='%s', object='%s'\n",
                   container_name,
                   object_name)
    }

    object_deleted := false

    if s3_ss.s3_session != nil && len(container_name) > 0 && len(object_name) > 0 {
	    /*
       request := &s3.DeleteObjectInput{
		Bucket: &container_name,
		Key:    &object_name,
       }

       _, err := s3_ss.s3_session.DeleteObject(request)
       if err == nil {
          object_deleted = true
       }
       */
    }

    return object_deleted
}

func (s3_ss S3StorageSystem) Get_object(container_name string,
                                        object_name string,
                                        local_file_path string) int {
    if s3_ss.Debug_mode {
        fmt.Printf("get_object: container='%s', object='%s', local_file_path='%s'\n",
                   container_name,
                   object_name,
                   local_file_path)
    }

    bytes_retrieved := 0

    if s3_ss.s3_session != nil && len(container_name) > 0 &&
       len(object_name) > 0 && len(local_file_path) > 0 {

        //s3_ss.s3_session.download_file(container_name, object_name, local_file_path)
        //if File_exists(local_file_path) {
        //    bytes_retrieved = int(Get_file_size(local_file_path))
        //}
    }

    return bytes_retrieved
}

func (s3 S3StorageSystem) Add_container(container_name string) {
   //ss.List_container_names = append(ss.List_container_names, container_name)
}

func (ss S3StorageSystem) Add_file_from_path(container_name string,
                                             object_name string,
                                             file_path string) bool {
   return false
}

func (ss S3StorageSystem) Get_container_names() []string {
   names := []string{}
   return names
}

func (ss S3StorageSystem) Has_container(container_name string) bool {
   return false
}

func (ss S3StorageSystem) List_containers() []string {
   containers := []string{}
   return containers
}

func (ss S3StorageSystem) Remove_container(constiner_name string) {
}

func (ss S3StorageSystem) Retrieve_file(fm *FileMetadata,
                                        local_directory string) int {
   return 0
}

func (ss S3StorageSystem) Store_file(fm *FileMetadata,
                                     file_contents []byte) bool {
   return false
}

