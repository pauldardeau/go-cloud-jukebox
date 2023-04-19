package jukebox

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type S3StorageSystem struct {
	debugMode      bool
	awsAccessKey   string
	awsSecretKey   string
	listContainers []string
	s3Client       *s3.S3
	s3Session      *session.Session
}

func NewS3StorageSystem(accessKey string,
	secretKey string,
	debugMode bool) *S3StorageSystem {

	ss := S3StorageSystem{
		debugMode:      debugMode,
		awsAccessKey:   accessKey,
		awsSecretKey:   secretKey,
		listContainers: nil,
		s3Client:       nil,
		s3Session:      nil,
	}

	if debugMode {
		fmt.Printf("Using accessKey='%s', secretKey='%s'\n", accessKey, secretKey)
	}

	return &ss
}

func (ss *S3StorageSystem) Enter() bool {
	connected := false
	if ss.debugMode {
		fmt.Println("attempting to connect to S3")
	}

	s3Config := aws.Config{
		Credentials:      credentials.NewStaticCredentials(ss.awsAccessKey, ss.awsSecretKey, ""),
		Endpoint:         aws.String("https://s3.us-central-1.wasabisys.com"),
		Region:           aws.String("us-central-1"),
		S3ForcePathStyle: aws.Bool(true),
	}

	s3Session, err := session.NewSessionWithOptions(session.Options{
		Config:  s3Config,
		Profile: "wasabi",
	})
	if err != nil {
		fmt.Println("unable to establish S3 session")
	} else {
		s3Client := s3.New(s3Session)

		fmt.Println("established S3 session")
		ss.s3Client = s3Client
		ss.s3Session = s3Session
		connected = true
	}

	if connected {
		listBuckets, err := ss.ListAccountContainers()
		if err == nil {
			ss.listContainers = *listBuckets
		} else {
			fmt.Printf("ss.ListAccountContainers returned error: %v\n", err)
		}
	}

	return connected
}

func (ss *S3StorageSystem) Exit() {
	if ss.s3Session != nil || ss.s3Client != nil {
		if ss.debugMode {
			fmt.Println("closing S3 connection object")
		}
		ss.s3Session = nil
		ss.s3Client = nil
	}
	ss.listContainers = nil
}

func (ss *S3StorageSystem) ListAccountContainers() (*[]string, error) {
	if ss.debugMode {
		fmt.Println("ListAccountContainers")
	}

	if ss.s3Client != nil {
		result, err := ss.s3Client.ListBuckets(&s3.ListBucketsInput{})
		if err != nil {
			fmt.Printf("ListBuckets failed\n")
			return nil, err
		} else {
			var listContainers []string
			for _, bucket := range result.Buckets {
				listContainers = append(listContainers, *bucket.Name)
			}
			return &listContainers, nil
		}
	} else {
		return nil, errors.New("no existing S3 session")
	}
}

func (ss *S3StorageSystem) CreateContainer(containerName string) bool {
	if ss.debugMode {
		fmt.Printf("CreateContainer: '%s'\n", containerName)
	}

	containerCreated := false

	if ss.s3Client != nil {
		bucket := aws.String(containerName)
		_, err := ss.s3Client.CreateBucket(&s3.CreateBucketInput{
			Bucket: bucket,
		})
		if err == nil {
			containerCreated = true
			ss.AddContainer(containerName)
		}
	}

	return containerCreated
}

func (ss *S3StorageSystem) DeleteContainer(containerName string) bool {
	if ss.debugMode {
		fmt.Printf("DeleteContainer: '%s'\n", containerName)
	}

	containerDeleted := false

	if ss.s3Client != nil {
		bucket := aws.String(containerName)
		_, err := ss.s3Client.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: bucket,
		})
		if err == nil {
			containerDeleted = true
			ss.DeleteContainer(containerName)
		}
	}

	return containerDeleted
}

func (ss *S3StorageSystem) ListContainerContents(containerName string) ([]string, error) {
	if ss.debugMode {
		fmt.Printf("ListContainerContents: '%s'\n", containerName)
	}

	if ss.s3Client != nil {
		bucket := aws.String(containerName)
		resp, err := ss.s3Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket})
		if err == nil {
			var listObjects []string
			for _, item := range resp.Contents {
				if len(*item.Key) > 0 {
					listObjects = append(listObjects, *item.Key)
				}
			}
			return listObjects, nil
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("no existing S3 session")
	}
}

func (ss *S3StorageSystem) GetObjectMetadata(containerName string,
	objectName string, dictProps *PropertySet) bool {
	if ss.debugMode {
		fmt.Printf("GetObjectMetadata: container='%s', object='%s'\n",
			containerName,
			objectName)
	}

	gotMetadata := false

	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.HeadObject
	if ss.s3Client != nil && len(containerName) > 0 && len(objectName) > 0 {
		//TODO: implement GetObjectMetadata
		/*
			        request := s3.HeadObjectInput{
			           Bucket: &containerName,
			           Key: &objectName,
			        }
			        response, err := ss.s3Session.HeadObject(&request)
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
	return gotMetadata
}

func (ss *S3StorageSystem) PutObject(containerName string,
	objectName string,
	fileContents []byte,
	headers *PropertySet) bool {

	objectAdded := false

	if ss.s3Client != nil && len(containerName) > 0 &&
		len(objectName) > 0 && len(fileContents) > 0 {
		//TODO: implement PutObject
	}

	return objectAdded
}

func (ss *S3StorageSystem) DeleteObject(containerName string, objectName string) bool {
	if ss.debugMode {
		fmt.Printf("DeleteObject: container='%s', object='%s'\n",
			containerName,
			objectName)
	}

	objectDeleted := false

	if ss.s3Client != nil && len(containerName) > 0 && len(objectName) > 0 {
		bucket := aws.String(containerName)
		item := aws.String(objectName)
		_, err := ss.s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: bucket,
			Key:    item,
		})
		if err == nil {
			objectDeleted = true
		}
	}

	return objectDeleted
}

func (ss *S3StorageSystem) GetObject(containerName string,
	objectName string,
	localFilePath string) int64 {

	if ss.debugMode {
		fmt.Printf("GetObject: container='%s', object='%s', localFilePath='%s'\n",
			containerName,
			objectName,
			localFilePath)
	}

	bytesRetrieved := int64(0)

	if ss.s3Client != nil && len(containerName) > 0 &&
		len(objectName) > 0 && len(localFilePath) > 0 {

		file, err := os.Create(localFilePath)
		if err != nil {
			return 0
		}
		defer file.Close()

		downloader := s3manager.NewDownloader(ss.s3Session)

		bucket := aws.String(containerName)
		filename := aws.String(objectName)

		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: bucket,
				Key:    filename,
			})
		if err != nil {
			fmt.Printf("error: unable to download %s/%s to %s\n", containerName, objectName, localFilePath)
		} else {
			bytesRetrieved = GetFileSize(localFilePath)
		}
	}

	return bytesRetrieved
}

func (ss *S3StorageSystem) AddContainer(containerName string) {
	if ss.listContainers != nil {
		ss.listContainers = append(ss.listContainers, containerName)
	}
}

func (ss *S3StorageSystem) AddFileFromPath(containerName string,
	objectName string,
	filePath string) bool {

	fileAdded := false

	if ss.s3Session != nil {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Unable to open file " + filePath)
			return false
		}
		defer file.Close()

		uploader := s3manager.NewUploader(ss.s3Session)

		bucket := aws.String(containerName)
		filename := aws.String(objectName)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: bucket,
			Key:    filename,
			Body:   file,
		})

		if err != nil {
			fmt.Printf("error: unable to add file from path: %s to %s/%s\n", filePath, containerName, objectName)
		} else {
			fileAdded = true
		}
	}

	return fileAdded
}

func (ss *S3StorageSystem) GetContainerNames() ([]string, error) {
	if ss.listContainers != nil {
		return ss.listContainers, nil
	} else {
		return nil, errors.New("no existing S3 session")
	}
}

func (ss *S3StorageSystem) HasContainer(containerName string) bool {
	haveContainer := false

	if ss.listContainers != nil {
		for _, existingContainerName := range ss.listContainers {
			if containerName == existingContainerName {
				haveContainer = true
				break
			}
		}
	} else {
		fmt.Printf("ss.listContainers is nil\n")
	}
	return haveContainer
}

func (ss *S3StorageSystem) RemoveContainer(containerName string) {
	if ss.listContainers != nil {
		foundIndex := -1
		for index, existingContainerName := range ss.listContainers {
			if existingContainerName == containerName {
				foundIndex = index
				break
			}
		}

		if foundIndex > -1 {
			// copy the last item in the slice to the index of the item we're removing
			// and then make the slice a sub-slice of 1 element less
			ss.listContainers[foundIndex] = ss.listContainers[len(ss.listContainers)-1]
			ss.listContainers = ss.listContainers[:len(ss.listContainers)-1]
		}
	}
}

func (ss *S3StorageSystem) RetrieveFile(fm *FileMetadata,
	localDirectory string) int64 {

	if len(localDirectory) > 0 {
		filePath := PathJoin(localDirectory, fm.FileUid)
		if ss.debugMode {
			fmt.Printf("retrieving container=%s\n", fm.ContainerName)
			fmt.Printf("retrieving object=%s\n", fm.ObjectName)
		}
		return ss.GetObject(fm.ContainerName, fm.ObjectName, filePath)
	} else {
		return 0
	}
}

func (ss *S3StorageSystem) StoreFile(fm *FileMetadata,
	fileContents []byte) bool {
	return ss.PutObject(fm.ContainerName, fm.ObjectName, fileContents, nil)
}
