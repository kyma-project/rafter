package controllers

import (
	"context"

	"github.com/kyma-project/rafter/internal/handler/assetgroup"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type assetService struct {
	client client.Client
	scheme *runtime.Scheme
}

func newAssetService(client client.Client, scheme *runtime.Scheme) *assetService {
	return &assetService{
		client: client,
		scheme: scheme,
	}
}

func (s *assetService) List(ctx context.Context, namespace string, labels map[string]string) ([]assetgroup.CommonAsset, error) {
	instances := &v1beta1.AssetList{}
	err := s.client.List(ctx, instances, client.MatchingLabels(labels))
	if err != nil {
		return nil, errors.Wrapf(err, "while listing Assets in namespace %s", namespace)
	}

	commons := make([]assetgroup.CommonAsset, 0, len(instances.Items))
	for _, instance := range instances.Items {
		if instance.Namespace != namespace {
			continue
		}
		common := s.assetToCommon(instance)
		commons = append(commons, common)
	}

	return commons, nil
}

func (s *assetService) Create(ctx context.Context, assetGroup v1.Object, commonAsset assetgroup.CommonAsset) error {
	instance := &v1beta1.Asset{
		ObjectMeta: commonAsset.ObjectMeta,
		Spec: v1beta1.AssetSpec{
			CommonAssetSpec: commonAsset.Spec,
		},
	}

	if err := controllerutil.SetControllerReference(assetGroup, instance, s.scheme); err != nil {
		return errors.Wrapf(err, "while creating Asset %s in namespace %s", commonAsset.Name, commonAsset.Namespace)
	}

	return s.client.Create(ctx, instance)
}

func (s *assetService) Update(ctx context.Context, commonAsset assetgroup.CommonAsset) error {
	instance := &v1beta1.Asset{}
	err := s.client.Get(ctx, types.NamespacedName{Name: commonAsset.Name, Namespace: commonAsset.Namespace}, instance)
	if err != nil {
		return errors.Wrapf(err, "while updating Asset %s in namespace %s", commonAsset.Name, commonAsset.Namespace)
	}

	updated := instance.DeepCopy()
	updated.Spec.CommonAssetSpec = commonAsset.Spec

	return s.client.Update(ctx, updated)
}

func (s *assetService) Delete(ctx context.Context, commonAsset assetgroup.CommonAsset) error {
	instance := &v1beta1.Asset{}
	err := s.client.Get(ctx, types.NamespacedName{Name: commonAsset.Name, Namespace: commonAsset.Namespace}, instance)
	if err != nil {
		return errors.Wrapf(err, "while deleting Asset %s in namespace %s", commonAsset.Name, commonAsset.Namespace)
	}

	return s.client.Delete(ctx, instance)
}

func (s *assetService) assetToCommon(instance v1beta1.Asset) assetgroup.CommonAsset {
	return assetgroup.CommonAsset{
		ObjectMeta: instance.ObjectMeta,
		Spec:       instance.Spec.CommonAssetSpec,
		Status:     instance.Status.CommonAssetStatus,
	}
}
