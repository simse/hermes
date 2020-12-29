package bucket

import "errors"

// ErrBucketExists is an error indicating a bucket is already owned by the user
var ErrBucketExists error = errors.New("Bucket already exists and is owned by user")

// ErrBucketExistsForeign is an error indicating a bucket is not owned by user
var ErrBucketExistsForeign error = errors.New("Bucket exists but is not owned by user")
