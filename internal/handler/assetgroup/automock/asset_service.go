// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	assetgroup "github.com/kyma-project/rafter/internal/handler/assetgroup"

	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AssetService is an autogenerated mock type for the AssetService type
type AssetService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, assetGroup, commonAsset
func (_m *AssetService) Create(ctx context.Context, assetGroup v1.Object, commonAsset assetgroup.CommonAsset) error {
	ret := _m.Called(ctx, assetGroup, commonAsset)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.Object, assetgroup.CommonAsset) error); ok {
		r0 = rf(ctx, assetGroup, commonAsset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, commonAsset
func (_m *AssetService) Delete(ctx context.Context, commonAsset assetgroup.CommonAsset) error {
	ret := _m.Called(ctx, commonAsset)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, assetgroup.CommonAsset) error); ok {
		r0 = rf(ctx, commonAsset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: ctx, namespace, labels
func (_m *AssetService) List(ctx context.Context, namespace string, labels map[string]string) ([]assetgroup.CommonAsset, error) {
	ret := _m.Called(ctx, namespace, labels)

	var r0 []assetgroup.CommonAsset
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]string) []assetgroup.CommonAsset); ok {
		r0 = rf(ctx, namespace, labels)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]assetgroup.CommonAsset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]string) error); ok {
		r1 = rf(ctx, namespace, labels)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, commonAsset
func (_m *AssetService) Update(ctx context.Context, commonAsset assetgroup.CommonAsset) error {
	ret := _m.Called(ctx, commonAsset)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, assetgroup.CommonAsset) error); ok {
		r0 = rf(ctx, commonAsset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
