// Code generated by mockery v1.0.0. DO NOT EDIT.
package automock

import fileheader "github.com/kyma-project/rafter/pkg/fileheader"
import mock "github.com/stretchr/testify/mock"

// FileHeader is an autogenerated mock type for the FileHeader type
type FileHeader struct {
	mock.Mock
}

// Filename provides a mock function with given fields:
func (_m *FileHeader) Filename() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Open provides a mock function with given fields:
func (_m *FileHeader) Open() (fileheader.File, error) {
	ret := _m.Called()

	var r0 fileheader.File
	if rf, ok := ret.Get(0).(func() fileheader.File); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fileheader.File)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields:
func (_m *FileHeader) Size() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}
