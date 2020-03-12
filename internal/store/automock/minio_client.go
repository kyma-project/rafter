// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	minio "github.com/minio/minio-go"
	mock "github.com/stretchr/testify/mock"
)

// MinioClient is an autogenerated mock type for the MinioClient type
type MinioClient struct {
	mock.Mock
}

// BucketExists provides a mock function with given fields: bucketName
func (_m *MinioClient) BucketExists(bucketName string) (bool, error) {
	ret := _m.Called(bucketName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(bucketName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bucketName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FPutObjectWithContext provides a mock function with given fields: ctx, bucketName, objectName, filePath, opts
func (_m *MinioClient) FPutObjectWithContext(ctx context.Context, bucketName string, objectName string, filePath string, opts minio.PutObjectOptions) (int64, error) {
	ret := _m.Called(ctx, bucketName, objectName, filePath, opts)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, minio.PutObjectOptions) int64); ok {
		r0 = rf(ctx, bucketName, objectName, filePath, opts)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, minio.PutObjectOptions) error); ok {
		r1 = rf(ctx, bucketName, objectName, filePath, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBucketPolicy provides a mock function with given fields: bucketName
func (_m *MinioClient) GetBucketPolicy(bucketName string) (string, error) {
	ret := _m.Called(bucketName)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(bucketName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(bucketName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListObjects provides a mock function with given fields: bucketName, objectPrefix, recursive, doneCh
func (_m *MinioClient) ListObjects(bucketName string, objectPrefix string, recursive bool, doneCh <-chan struct{}) <-chan minio.ObjectInfo {
	ret := _m.Called(bucketName, objectPrefix, recursive, doneCh)

	var r0 <-chan minio.ObjectInfo
	if rf, ok := ret.Get(0).(func(string, string, bool, <-chan struct{}) <-chan minio.ObjectInfo); ok {
		r0 = rf(bucketName, objectPrefix, recursive, doneCh)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan minio.ObjectInfo)
		}
	}

	return r0
}

// MakeBucket provides a mock function with given fields: bucketName, location
func (_m *MinioClient) MakeBucket(bucketName string, location string) error {
	ret := _m.Called(bucketName, location)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(bucketName, location)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveBucket provides a mock function with given fields: bucketName
func (_m *MinioClient) RemoveBucket(bucketName string) error {
	ret := _m.Called(bucketName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(bucketName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveObjectsWithContext provides a mock function with given fields: ctx, bucketName, objectsCh
func (_m *MinioClient) RemoveObjectsWithContext(ctx context.Context, bucketName string, objectsCh <-chan string) <-chan minio.RemoveObjectError {
	ret := _m.Called(ctx, bucketName, objectsCh)

	var r0 <-chan minio.RemoveObjectError
	if rf, ok := ret.Get(0).(func(context.Context, string, <-chan string) <-chan minio.RemoveObjectError); ok {
		r0 = rf(ctx, bucketName, objectsCh)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan minio.RemoveObjectError)
		}
	}

	return r0
}

// SetBucketPolicy provides a mock function with given fields: bucketName, policy
func (_m *MinioClient) SetBucketPolicy(bucketName string, policy string) error {
	ret := _m.Called(bucketName, policy)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(bucketName, policy)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
