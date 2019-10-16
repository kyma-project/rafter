package testsuite

import (
	"time"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/kyma-project/rafter/tests/asset-store/pkg/resource"
	"github.com/kyma-project/rafter/tests/asset-store/pkg/waiter"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type bucket struct {
	resCli      *resource.Resource
	name        string
	namespace   string
	waitTimeout time.Duration
}

func newBucket(dynamicCli dynamic.Interface, name, namespace string, waitTimeout time.Duration, logFn func(format string, args ...interface{})) *bucket {
	return &bucket{
		resCli: resource.New(dynamicCli, schema.GroupVersionResource{
			Version:  v1beta1.GroupVersion.Version,
			Group:    v1beta1.GroupVersion.Group,
			Resource: "buckets",
		}, namespace, logFn),
		name:        name,
		namespace:   namespace,
		waitTimeout: waitTimeout,
	}
}

func (b *bucket) Create() error {
	bucket := &v1beta1.Bucket{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Bucket",
			APIVersion: v1beta1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
		},
		Spec: v1beta1.BucketSpec{
			CommonBucketSpec: v1beta1.CommonBucketSpec{
				Policy: v1beta1.BucketPolicyReadOnly,
			},
		},
	}

	err := b.resCli.Create(bucket)
	if err != nil {
		return errors.Wrapf(err, "while creating Bucket %s in namespace %s", b.name, b.namespace)
	}

	return err
}

func (b *bucket) WaitForStatusReady() error {
	err := waiter.WaitAtMost(func() (bool, error) {
		res, err := b.Get(b.name)
		if err != nil {
			return false, err
		}

		if res.Status.Phase != v1beta1.BucketReady {
			return false, nil
		}

		return true, nil
	}, b.waitTimeout)
	if err != nil {
		return errors.Wrapf(err, "while waiting for ready Bucket resources")
	}

	return nil
}

func (b *bucket) Get(name string) (*v1beta1.Bucket, error) {
	u, err := b.resCli.Get(name)
	if err != nil {
		return nil, err
	}

	var res v1beta1.Bucket
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &res)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting Bucket %s", name)
	}

	return &res, nil
}

func (b *bucket) Delete() error {
	err := b.resCli.Delete(b.name)
	if err != nil {
		return errors.Wrapf(err, "while deleting Bucket %s in namespace %s", b.name, b.namespace)
	}

	return nil
}
