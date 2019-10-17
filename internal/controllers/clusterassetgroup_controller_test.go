package controllers

import (
	"context"
	"time"

	"github.com/kyma-project/rafter/internal/source"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("Asset", func() {
	var (
		assetgroup *v1beta1.ClusterAssetGroup
		reconciler *ClusterAssetGroupReconciler
		request    ctrl.Request
	)

	BeforeEach(func() {
		assetgroup = newFixClusterAssetGroup()
		Expect(k8sClient.Create(context.TODO(), assetgroup)).To(Succeed())

		request = ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      assetgroup.Name,
				Namespace: assetgroup.Namespace,
			},
		}

		assetService := newClusterAssetService(k8sClient, scheme.Scheme)
		bucketService := newClusterBucketService(k8sClient, scheme.Scheme, "us-east-1")

		reconciler = &ClusterAssetGroupReconciler{
			Client:           k8sClient,
			Log:              ctrl.Log,
			recorder:         record.NewFakeRecorder(100),
			relistInterval:   60 * time.Hour,
			assetSvc:         assetService,
			bucketSvc:        bucketService,
			webhookConfigSvc: webhookConfigSvc,
		}
	})

	It("should successfully create, update and delete AssetGroup", func() {
		By("creating the ClusterAssetGroup")
		result, err := reconciler.Reconcile(request)
		validateReconcilation(err, result)
		assetgroup = &v1beta1.ClusterAssetGroup{}
		Expect(k8sClient.Get(context.TODO(), request.NamespacedName, assetgroup)).To(Succeed())
		validateAssetGroup(assetgroup.Status.CommonAssetGroupStatus, assetgroup.ObjectMeta, v1beta1.AssetGroupPending, v1beta1.AssetGroupWaitingForAssets)

		By("ClusterAssets changes states to ready")
		assets := &v1beta1.ClusterAssetList{}
		Expect(k8sClient.List(context.TODO(), assets)).To(Succeed())
		Expect(assets.Items).To(HaveLen(len(assetgroup.Spec.Sources)))

		for _, asset := range assets.Items {
			asset.Status.Phase = v1beta1.AssetReady
			asset.Status.LastHeartbeatTime = v1.Now()
			Expect(k8sClient.Status().Update(context.TODO(), &asset)).To(Succeed())

			if asset.Annotations["cms.kyma-project.io/asset-short-name"] == "source-one" {
				Expect(asset.Spec.Parameters).ToNot(BeNil())
				Expect(asset.Spec.Parameters).To(Equal(&fixParameters))
			} else {
				Expect(asset.Spec.Parameters).To(BeNil())
			}
		}

		result, err = reconciler.Reconcile(request)
		validateReconcilation(err, result)
		assetgroup = &v1beta1.ClusterAssetGroup{}
		Expect(k8sClient.Get(context.TODO(), request.NamespacedName, assetgroup)).To(Succeed())
		validateAssetGroup(assetgroup.Status.CommonAssetGroupStatus, assetgroup.ObjectMeta, v1beta1.AssetGroupReady, v1beta1.AssetGroupAssetsReady)

		By("updating the ClusterAssetGroup")
		assetgroup.Spec.Sources = source.FilterByType(assetgroup.Spec.Sources, "dita")
		markdownIndex := source.IndexByType(assetgroup.Spec.Sources, "markdown")
		Expect(markdownIndex).NotTo(Equal(-1))
		assetgroup.Spec.Sources[markdownIndex].Filter = "zyx"
		Expect(k8sClient.Update(context.TODO(), assetgroup)).To(Succeed())

		result, err = reconciler.Reconcile(request)
		validateReconcilation(err, result)
		assetgroup = &v1beta1.ClusterAssetGroup{}
		Expect(k8sClient.Get(context.TODO(), request.NamespacedName, assetgroup)).To(Succeed())
		validateAssetGroup(assetgroup.Status.CommonAssetGroupStatus, assetgroup.ObjectMeta, v1beta1.AssetGroupPending, v1beta1.AssetGroupWaitingForAssets)

		assets = &v1beta1.ClusterAssetList{}
		Expect(k8sClient.List(context.TODO(), assets)).To(Succeed())
		Expect(assets.Items).To(HaveLen(len(assetgroup.Spec.Sources)))
		for _, a := range assets.Items {
			if a.Annotations["cms.kyma-project.io/asset-short-name"] != "source-two" {
				continue
			}
			Expect(a.Spec.Source.Filter).To(Equal("zyx"))
		}

		By("deleting the ClusterAssetGroup")
		Expect(k8sClient.Delete(context.TODO(), assetgroup)).To(Succeed())

		_, err = reconciler.Reconcile(request)
		Expect(err).To(Succeed())

		assetgroup = &v1beta1.ClusterAssetGroup{}
		err = k8sClient.Get(context.TODO(), request.NamespacedName, assetgroup)
		Expect(err).To(HaveOccurred())
		Expect(apiErrors.IsNotFound(err)).To(BeTrue())

	})
})

func newFixClusterAssetGroup() *v1beta1.ClusterAssetGroup {
	return &v1beta1.ClusterAssetGroup{
		ObjectMeta: ctrl.ObjectMeta{
			Name: string(uuid.NewUUID()),
		},
		Spec: v1beta1.ClusterAssetGroupSpec{
			CommonAssetGroupSpec: v1beta1.CommonAssetGroupSpec{
				Description: "Test topic, have fun",
				DisplayName: "Test Topic",
				Sources: []v1beta1.Source{
					{
						Name:       "source-one",
						Type:       "openapi",
						Mode:       v1beta1.AssetGroupSingle,
						URL:        "https://dummy.url/single",
						Parameters: &fixParameters,
					},
					{
						Name:   "source-two",
						Type:   "markdown",
						Filter: "xyz",
						Mode:   v1beta1.AssetGroupPackage,
						URL:    "https://dummy.url/package",
					},
					{
						Name:   "source-three",
						Type:   "dita",
						Filter: "xyz",
						Mode:   v1beta1.AssetGroupIndex,
						URL:    "https://dummy.url/index",
					},
					{
						Name: "source-four",
						Type: "openapi",
						Mode: v1beta1.AssetGroupPackage,
						URL:  "https://dummy.url/single",
					},
				},
			},
		},
		Status: v1beta1.ClusterAssetGroupStatus{CommonAssetGroupStatus: v1beta1.CommonAssetGroupStatus{
			LastHeartbeatTime: v1.Now(),
		}},
	}
}
