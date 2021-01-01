package deploy

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/simse/hermes/internal/bucket"
)

// Manifest represents a hermes stack manifest consisting of: hermes init version, hermes deploy version, associated resources and deployed files
type Manifest struct {
	InitVersion   string        `json:"init_version"`
	DeployVersion string        `json:"deploy_version"`
	Files         []File        `json:"files"`
	CloudFront    string        `json:"cloudfront"`
	Bucket        string        `json:"bucket"`
	Domain        string        `json:"domain"`
	DomainAliases []string      `json:"domain_aliases"`
	EdgeHandlers  []EdgeHandler `json:"edge_handlers"`
}

// File represents a file in a bucket
type File struct {
	Key      string `json:"key"`
	Checksum string `json:"checksum"` // Algorithm: SHA256
	Size     int64  `json:"size"`
}

// EdgeHandler represents a lambda@edge definition
type EdgeHandler struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Type   string `json:"type"` // [ORIGIN_REQUEST, ORIGIN_RESPONSE, VIEWER_REQUEST, VIEWER_RESPONSE]
}

// ScanDir scans a directory and returns a list of files with name, checksum and size
func ScanDir(path string) ([]File, error) {
	var files []File

	root := path

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// Calculate SHA256 checksum
		fileReader, fileOpenErr := os.Open(path)
		if fileOpenErr != nil {
			return nil
		}
		defer fileReader.Close()

		hash := sha256.New()
		if _, err := io.Copy(hash, fileReader); err != nil {
			return nil
		}

		file := File{
			Key:      strings.TrimPrefix(path, root),
			Size:     info.Size(),
			Checksum: hex.EncodeToString(hash.Sum(nil)),
		}

		files = append(files, file)

		return nil
	})

	return files, nil
}

// CreateDefault deploys the default site to a bucket returning a list of files
func CreateDefault() ([]File, error) {
	var files []File

	// TODO: Update dynamically
	files = append(files, File{
		Key:      "index.html",
		Size:     857,
		Checksum: "05e5f045efbe9bd5669e0193591d5a1e35ff184a254529abed242ea87bb01288",
	})

	return files, nil
}

// ReadManifest reads a manifest from a bucket
func ReadManifest(bucketName string) (Manifest, error) {
	manifest := Manifest{}

	getObjectInput := s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(".hermes.json"),
	}
	getObjectOutput, err := bucket.S3.GetObject(&getObjectInput)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(getObjectOutput.Body)
	manifestBytes := buf.Bytes()

	json.Unmarshal(manifestBytes, &manifest)

	return manifest, nil
}

// WriteManifest writes a manifest to a bucket
func WriteManifest(bucketName string, manifest Manifest) error {
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		panic(err)
	}

	putObjectInput := s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(".hermes.json"),
		Body:   aws.ReadSeekCloser(bytes.NewReader(manifestBytes)),
	}

	putObjectOutput, putError := bucket.S3.PutObject(&putObjectInput)
	if putError != nil {
		panic(err)
	}

	fmt.Println(putObjectOutput)

	return nil
}
