package controllers

import (
	"context"

	"github.com/kyma-project/kyma/components/cms-controller-manager/internal/handler/docstopic"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type clusterAssetService struct {
	client client.Client
	scheme *runtime.Scheme
}

func newClusterAssetService(client client.Client, scheme *runtime.Scheme) *clusterAssetService {
	return &clusterAssetService{
		client: client,
		scheme: scheme,
	}
}

func (s *clusterAssetService) List(ctx context.Context, namespace string, labels map[string]string) ([]docstopic.CommonAsset, error) {
	instances := &v1beta1.ClusterAssetList{}
	err := s.client.List(ctx, instances, client.MatchingLabels(labels))
	if err != nil {
		return nil, errors.Wrap(err, "while listing ClusterAssets")
	}

	commons := make([]docstopic.CommonAsset, 0, len(instances.Items))
	for _, instance := range instances.Items {
		common := s.assetToCommon(instance)
		commons = append(commons, common)
	}

	return commons, nil
}

func (s *clusterAssetService) Create(ctx context.Context, docsTopic v1.Object, commonAsset docstopic.CommonAsset) error {
	instance := &v1beta1.ClusterAsset{
		ObjectMeta: commonAsset.ObjectMeta,
		Spec: v1beta1.ClusterAssetSpec{
			CommonAssetSpec: commonAsset.Spec,
		},
	}

	if err := controllerutil.SetControllerReference(docsTopic, instance, s.scheme); err != nil {
		return errors.Wrapf(err, "while creating ClusterAsset %s", commonAsset.Name)
	}

	return s.client.Create(ctx, instance)
}

func (s *clusterAssetService) Update(ctx context.Context, commonAsset docstopic.CommonAsset) error {
	instance := &v1beta1.ClusterAsset{}
	err := s.client.Get(ctx, types.NamespacedName{Name: commonAsset.Name, Namespace: commonAsset.Namespace}, instance)
	if err != nil {
		return errors.Wrapf(err, "while updating ClusterAsset %s", commonAsset.Name)
	}

	updated := instance.DeepCopy()
	updated.Spec.CommonAssetSpec = commonAsset.Spec

	return s.client.Update(ctx, updated)
}

func (s *clusterAssetService) Delete(ctx context.Context, commonAsset docstopic.CommonAsset) error {
	instance := &v1beta1.ClusterAsset{}
	err := s.client.Get(ctx, types.NamespacedName{Name: commonAsset.Name, Namespace: commonAsset.Namespace}, instance)
	if err != nil {
		return errors.Wrapf(err, "while deleting ClusterAsset %s", commonAsset.Name)
	}

	return s.client.Delete(ctx, instance)
}

func (s *clusterAssetService) assetToCommon(instance v1beta1.ClusterAsset) docstopic.CommonAsset {
	return docstopic.CommonAsset{
		ObjectMeta: instance.ObjectMeta,
		Spec:       instance.Spec.CommonAssetSpec,
		Status:     instance.Status.CommonAssetStatus,
	}
}
