package docstopic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-project/rafter/internal/handler/docstopic"
	"github.com/kyma-project/rafter/internal/handler/docstopic/automock"
	"github.com/kyma-project/rafter/internal/webhookconfig"
	amcfg "github.com/kyma-project/rafter/internal/webhookconfig/automock"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("docstopic-test")

func Test_findSource(t *testing.T) {
	findSource := docstopic.FindSource()
	g := gomega.NewGomegaWithT(t)
	testSource1 := v1beta1.Source{
		Name:   "test",
		Type:   "me",
		URL:    "test.it",
		Mode:   "mode1",
		Filter: "*",
	}

	tests := []struct {
		name     string
		srcSlice []v1beta1.Source
		srcName  v1beta1.DocsTopicSourceName
		srcType  v1beta1.DocsTopicSourceType
		matcher  types.GomegaMatcher
	}{
		{
			name:     "empty slice",
			srcSlice: []v1beta1.Source{},
			srcName:  "name",
			srcType:  "type",
			matcher:  gomega.BeNil(),
		},
		{
			name:     "nil slice",
			srcSlice: nil,
			srcName:  "name",
			srcType:  "type",
			matcher:  gomega.BeNil(),
		},
		{
			name:     "found",
			srcSlice: []v1beta1.Source{testSource1},
			srcName:  "test",
			srcType:  "me",
			matcher:  gomega.Equal(&testSource1),
		},
		{
			name:     "not found",
			srcSlice: []v1beta1.Source{},
			srcName:  "test",
			srcType:  "me",
			matcher:  gomega.BeNil(),
		},
		{
			name: "found2",
			srcSlice: []v1beta1.Source{
				v1beta1.Source{
					Name:   "test",
					Type:   "me2",
					URL:    "test.it",
					Mode:   "mode1",
					Filter: "*",
				},
				testSource1,
			},
			srcName: "test",
			srcType: "me",
			matcher: gomega.Equal(&testSource1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := findSource(test.srcSlice, test.srcName, test.srcType)
			g.Expect(actual).To(test.matcher)
		})
	}
}

func TestDocstopicHandler_Handle_AddOrUpdate(t *testing.T) {
	sourceName := v1beta1.DocsTopicSourceName("t1")
	assetType := v1beta1.DocsTopicSourceType("swag")

	t.Run("Create", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{"test-bucket"}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(nil, nil).Once()
		assetSvc.On("Create", ctx, testData, mock.Anything).Return(nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicPending))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicWaitingForAssets))
	})

	t.Run("CreateWithMetadata", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		metadata := &runtime.RawExtension{Raw: []byte(`{"json":"true"}`)}
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, metadata)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{"test-bucket"}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(nil, nil).Once()
		assetSvc.On("Create", ctx, testData, mock.Anything).Return(nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicPending))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicWaitingForAssets))
	})

	t.Run("CreateError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{"test-bucket"}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(nil, nil).Once()
		assetSvc.On("Create", ctx, testData, mock.Anything).Return(errors.New("test-data")).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicAssetsCreationFailed))
	})

	t.Run("Update", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{
			testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil),
			testSource(sourceName, "markdown", "https://dummy.url", v1beta1.DocsTopicSingle, nil),
		}
		testData := testData("halo", sources)
		source, _ := getSourceByType(sources, sourceName)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		existingAsset.Spec.Source.Filter = "xyz"
		existingAssets := []docstopic.CommonAsset{existingAsset}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		assetSvc.On("Update", ctx, mock.Anything).Return(nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicPending))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicWaitingForAssets))
	})

	t.Run("UpdateErr2", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{
			testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil),
			testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil),
		}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicAssetsSpecValidationFailed))
	})

	t.Run("UpdateError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)
		source, _ := getSourceByType(sources, sourceName)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		existingAsset.Spec.Source.Filter = "xyz"
		existingAssets := []docstopic.CommonAsset{existingAsset}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		assetSvc.On("Update", ctx, mock.Anything).Return(errors.New("test-error")).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicAssetsUpdateFailed))
	})

	t.Run("Delete", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)
		source, ok := getSourceByType(sources, sourceName)
		g.Expect(ok, true)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		toRemove := commonAsset("papa", assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		existingAssets := []docstopic.CommonAsset{existingAsset, toRemove}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		assetSvc.On("Delete", ctx, toRemove).Return(nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicPending))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicWaitingForAssets))
	})

	t.Run("DeleteError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)
		source, ok := getSourceByType(sources, sourceName)
		g.Expect(ok, true)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		toRemove := commonAsset("papa", assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		existingAssets := []docstopic.CommonAsset{existingAsset, toRemove}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		assetSvc.On("Delete", ctx, toRemove).Return(errors.New("test-error")).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicAssetsDeletionFailed))
	})

	t.Run("CreateWithBucket", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return(nil, nil).Once()
		bucketSvc.On("Create", ctx, mock.Anything, false, map[string]string{"cms.kyma-project.io/access": "public"}).Return(nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(nil, nil).Once()
		assetSvc.On("Create", ctx, testData, mock.Anything).Return(nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicPending))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicWaitingForAssets))
	})

	t.Run("BucketCreationError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return(nil, nil).Once()
		bucketSvc.On("Create", ctx, mock.Anything, false, map[string]string{"cms.kyma-project.io/access": "public"}).Return(errors.New("test-error")).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicBucketError))
	})

	t.Run("BucketListingError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return(nil, errors.New("test-error")).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicBucketError))
	})

	t.Run("AssetsListingError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{"test-bucket"}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(nil, errors.New("test-error")).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicFailed))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicAssetsListingFailed))
	})
}

func TestDocstopicHandler_Handle_Status(t *testing.T) {
	sourceName := v1beta1.DocsTopicSourceName("t1")
	assetType := v1beta1.DocsTopicSourceType("swag")

	t.Run("NotChanged", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)
		testData.Status.Phase = v1beta1.DocsTopicPending
		source, ok := getSourceByType(sources, sourceName)
		g.Expect(ok, true)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetPending)
		existingAssets := []docstopic.CommonAsset{existingAsset}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).To(gomega.BeNil())
	})

	t.Run("Changed", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)
		testData.Status.Phase = v1beta1.DocsTopicPending
		source, ok := getSourceByType(sources, sourceName)
		g.Expect(ok, true)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetReady)
		existingAssets := []docstopic.CommonAsset{existingAsset}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicReady))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicAssetsReady))
	})

	t.Run("AssetError", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)
		ctx := context.TODO()
		bucketName := "test-bucket"
		sources := []v1beta1.Source{testSource(sourceName, assetType, "https://dummy.url", v1beta1.DocsTopicSingle, nil)}
		testData := testData("halo", sources)
		testData.Status.Phase = v1beta1.DocsTopicReady
		source, ok := getSourceByType(sources, sourceName)
		g.Expect(ok, true)
		existingAsset := commonAsset(sourceName, assetType, testData.Name, bucketName, *source, v1beta1.AssetFailed)
		existingAssets := []docstopic.CommonAsset{existingAsset}

		assetSvc := new(automock.AssetService)
		defer assetSvc.AssertExpectations(t)
		bucketSvc := new(automock.BucketService)
		defer bucketSvc.AssertExpectations(t)
		webhookConfSvc := new(amcfg.AssetWebhookConfigService)
		defer webhookConfSvc.AssertExpectations(t)

		bucketSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/access": "public"}).Return([]string{bucketName}, nil).Once()
		assetSvc.On("List", ctx, testData.Namespace, map[string]string{"cms.kyma-project.io/docs-topic": testData.Name}).Return(existingAssets, nil).Once()
		webhookConfSvc.On("Get", ctx).Return(webhookconfig.AssetWebhookConfigMap{}, nil).Once()

		handler := docstopic.New(log, fakeRecorder(), assetSvc, bucketSvc, webhookConfSvc)

		// When
		status, err := handler.Handle(ctx, testData, testData.Spec.CommonDocsTopicSpec, testData.Status.CommonDocsTopicStatus)

		// Then
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status.Phase).To(gomega.Equal(v1beta1.DocsTopicPending))
		g.Expect(status.Reason).To(gomega.Equal(v1beta1.DocsTopicWaitingForAssets))
	})
}

func fakeRecorder() record.EventRecorder {
	return record.NewFakeRecorder(20)
}

func testSource(sourceName v1beta1.DocsTopicSourceName, sourceType v1beta1.DocsTopicSourceType, url string, mode v1beta1.DocsTopicSourceMode, parameters *runtime.RawExtension) v1beta1.Source {
	return v1beta1.Source{
		Name:       sourceName,
		Type:       sourceType,
		URL:        url,
		Mode:       mode,
		Parameters: parameters,
	}
}

func testData(name string, sources []v1beta1.Source) *v1beta1.DocsTopic {
	return &v1beta1.DocsTopic{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "test",
		},
		Spec: v1beta1.DocsTopicSpec{
			CommonDocsTopicSpec: v1beta1.CommonDocsTopicSpec{
				DisplayName: fmt.Sprintf("%s Display", name),
				Description: fmt.Sprintf("%s Description", name),
				Sources:     sources,
			},
		},
	}
}

func commonAsset(name v1beta1.DocsTopicSourceName, assetType v1beta1.DocsTopicSourceType, docsName, bucketName string, source v1beta1.Source, phase v1beta1.AssetPhase) docstopic.CommonAsset {
	return docstopic.CommonAsset{
		ObjectMeta: v1.ObjectMeta{
			Name:      string(name),
			Namespace: "test",
			Labels: map[string]string{
				"cms.kyma-project.io/docs-topic": docsName,
				"cms.kyma-project.io/type":       string(assetType),
			},
			Annotations: map[string]string{
				"cms.kyma-project.io/asset-short-name": string(name),
			},
		},
		Spec: v1beta1.CommonAssetSpec{
			Source: v1beta1.AssetSource{
				URL:    source.URL,
				Mode:   v1beta1.AssetMode(source.Mode),
				Filter: source.Filter,
			},
			BucketRef: v1beta1.AssetBucketRef{
				Name: bucketName,
			},
		},
		Status: v1beta1.CommonAssetStatus{
			Phase: phase,
		},
	}
}

func getSourceByType(slice []v1beta1.Source, sourceName v1beta1.DocsTopicSourceName) (*v1beta1.Source, bool) {
	for _, source := range slice {
		if source.Name != sourceName {
			continue
		}
		return &source, true
	}
	return nil, false
}
