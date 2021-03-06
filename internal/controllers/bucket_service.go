package controllers

import (
	"context"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type bucketService struct {
	client client.Client
	scheme *runtime.Scheme
	region string
}

func newBucketService(client client.Client, scheme *runtime.Scheme, region string) *bucketService {
	return &bucketService{
		client: client,
		scheme: scheme,
		region: region,
	}
}

func (s *bucketService) List(ctx context.Context, namespace string, labels map[string]string) ([]string, error) {
	instances := &v1beta1.BucketList{}
	err := s.client.List(ctx, instances, client.MatchingLabels(labels))
	if err != nil {
		return nil, errors.Wrapf(err, "while listing Buckets in namespace %s", namespace)
	}

	names := make([]string, 0, len(instances.Items))
	for _, instance := range instances.Items {
		if instance.Namespace != namespace {
			continue
		}
		namespacedName := instance.GetName()
		names = append(names, namespacedName)
	}

	return names, nil
}

func (s *bucketService) Create(ctx context.Context, namespacedName types.NamespacedName, private bool, labels map[string]string) error {
	policy := v1beta1.BucketPolicyReadOnly
	if private {
		policy = v1beta1.BucketPolicyNone
	}

	instance := &v1beta1.Bucket{
		ObjectMeta: v1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
			Labels:    labels,
		},
		Spec: v1beta1.BucketSpec{
			CommonBucketSpec: v1beta1.CommonBucketSpec{
				Policy: policy,
				Region: v1beta1.BucketRegion(s.region),
			},
		},
	}

	if err := s.client.Create(ctx, instance); err != nil {
		return errors.Wrapf(err, "while creating Bucket %s in namespace %s", namespacedName.Name, namespacedName.Namespace)
	}

	return nil
}
