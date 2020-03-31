package asset_test

import (
	"context"
	"strings"
	"testing"
	"time"

	engine "github.com/kyma-project/rafter/internal/assethook"
	engineMock "github.com/kyma-project/rafter/internal/assethook/automock"
	"github.com/kyma-project/rafter/internal/handler/asset"
	loaderMock "github.com/kyma-project/rafter/internal/loader/automock"
	storeMock "github.com/kyma-project/rafter/internal/store/automock"
	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("asset-test")

const (
	remoteBucketName = "bucket-name"
)

func TestAssetHandler_Handle_OnAddOrUpdate(t *testing.T) {
	t.Run("OnAdd", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetPending))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetScheduled))
	})

	t.Run("OnAddWithDisplayName", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testDataWithDisplayName("test-asset", "test-bucket", "https://localhost/test.md", "source displayName")

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetPending))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetScheduled))
	})

	t.Run("OnUpdate", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.ObjectMeta.Generation = int64(2)
		asset.Status.ObservedGeneration = int64(1)

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetPending))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetScheduled))
	})
}

func TestAssetHandler_Handle_Default(t *testing.T) {
	// Given
	g := NewGomegaWithT(t)
	ctx := context.TODO()
	relistInterval := time.Minute
	now := time.Now()
	asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
	asset.ObjectMeta.Generation = int64(1)
	asset.Status.ObservedGeneration = int64(1)

	handler, mocks := newHandler(relistInterval)
	defer mocks.AssertExpectations(t)

	// When
	status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

	// Then
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(status).To(BeZero())
}

func TestAssetHandler_Handle_OnReady(t *testing.T) {
	t.Run("NotTaken", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.NewTime(now)
		asset.Status.CommonAssetStatus.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("NotChanged", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.NewTime(now.Add(-2 * relistInterval))
		asset.Status.CommonAssetStatus.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ContainsAllObjects", ctx, remoteBucketName, asset.Name, mock.AnythingOfType("[]string")).Return(true, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetReady))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetUploaded))
	})

	t.Run("BucketNotReady", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "notReady", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.NewTime(now.Add(-2 * relistInterval))
		asset.Status.CommonAssetStatus.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetPending))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetBucketNotReady))
	})

	t.Run("BucketError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "error", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.NewTime(now.Add(-2 * relistInterval))
		asset.Status.CommonAssetStatus.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetBucketError))
	})

	t.Run("MissingFiles", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.NewTime(now.Add(-2 * relistInterval))
		asset.Status.CommonAssetStatus.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ContainsAllObjects", ctx, remoteBucketName, asset.Name, mock.AnythingOfType("[]string")).Return(false, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetMissingContent))
	})

	t.Run("ContainsError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.NewTime(now.Add(-2 * relistInterval))
		asset.Status.CommonAssetStatus.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ContainsAllObjects", ctx, remoteBucketName, asset.Name, mock.AnythingOfType("[]string")).Return(false, errors.New("nope")).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetRemoteContentVerificationError))
	})
}

func TestAssetHandler_Handle_OnPending(t *testing.T) {
	t.Run("WithWebhooks", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.store.On("PutObjects", ctx, remoteBucketName, asset.Name, "/tmp", mock.AnythingOfType("[]string")).Return(nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.validator.On("Validate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.ValidationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.metadataExtractor.On("Extract", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MetadataWebhookService).Return(nil, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetReady))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetUploaded))
	})

	t.Run("WithoutWebhooks", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation
		asset.Spec.Source.ValidationWebhookService = nil
		asset.Spec.Source.MutationWebhookService = nil
		asset.Spec.Source.MetadataWebhookService = nil

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.store.On("PutObjects", ctx, remoteBucketName, asset.Name, "/tmp", mock.AnythingOfType("[]string")).Return(nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetReady))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetUploaded))
	})

	t.Run("LoadError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, errors.New("nope")).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetPullingFailed))
	})

	t.Run("MutationFailed", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: false}, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetMutationFailed))
	})

	t.Run("MutationError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: false}, errors.New("nope")).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetMutationError))
	})

	t.Run("MetadataExtractionFailed", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.validator.On("Validate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.ValidationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.metadataExtractor.On("Extract", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MetadataWebhookService).Return(nil, errors.New("nope")).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetMetadataExtractionFailed))
	})

	t.Run("ValidationError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.validator.On("Validate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.ValidationWebhookService).Return(engine.Result{Success: false}, errors.New("nope")).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetValidationError))
	})

	t.Run("ValidationFailed", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.validator.On("Validate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.ValidationWebhookService).Return(engine.Result{Success: false}, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetValidationFailed))
	})

	t.Run("UploadError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.store.On("PutObjects", ctx, remoteBucketName, asset.Name, "/tmp", mock.AnythingOfType("[]string")).Return(errors.New("nope")).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.validator.On("Validate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.ValidationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.metadataExtractor.On("Extract", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MetadataWebhookService).Return(nil, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetUploadFailed))
	})

	t.Run("BucketNotReady", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "notReady", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetPending))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetBucketNotReady))
	})

	t.Run("BucketStatusError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "error", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetFailed))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetBucketError))
	})

	t.Run("OnBucketNotReadyBeforeTime", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "error", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetPending
		asset.Status.CommonAssetStatus.Reason = v1beta1.AssetBucketNotReady
		asset.Status.CommonAssetStatus.LastHeartbeatTime = v1.Now()
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})
}

func TestAssetHandler_Handle_OnFailed(t *testing.T) {
	t.Run("ShouldHandle", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetFailed
		asset.Status.CommonAssetStatus.Reason = v1beta1.AssetBucketError
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()
		mocks.store.On("PutObjects", ctx, remoteBucketName, asset.Name, "/tmp", mock.AnythingOfType("[]string")).Return(nil).Once()
		mocks.loader.On("Load", asset.Spec.Source.URL, asset.Name, asset.Spec.Source.Mode, asset.Spec.Source.Filter).Return("/tmp", nil, nil).Once()
		mocks.loader.On("Clean", "/tmp").Return(nil).Once()
		mocks.mutator.On("Mutate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MutationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.validator.On("Validate", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.ValidationWebhookService).Return(engine.Result{Success: true}, nil).Once()
		mocks.metadataExtractor.On("Extract", ctx, "/tmp", mock.AnythingOfType("[]string"), asset.Spec.Source.MetadataWebhookService).Return(nil, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).ToNot(BeZero())
		g.Expect(status.Phase).To(Equal(v1beta1.AssetReady))
		g.Expect(status.Reason).To(Equal(v1beta1.AssetUploaded))
	})

	t.Run("ValidationFailed", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetFailed
		asset.Status.CommonAssetStatus.Reason = v1beta1.AssetValidationFailed
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("MutationFailed", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		asset.Status.CommonAssetStatus.Phase = v1beta1.AssetFailed
		asset.Status.CommonAssetStatus.Reason = v1beta1.AssetMutationFailed
		asset.Status.ObservedGeneration = asset.Generation

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})
}

func TestAssetHandler_Handle_OnDelete(t *testing.T) {
	t.Run("NoFiles", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		nowMeta := v1.Now()
		asset.ObjectMeta.DeletionTimestamp = &nowMeta

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("MultipleFiles", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		nowMeta := v1.Now()
		asset.ObjectMeta.DeletionTimestamp = &nowMeta
		files := []string{"test/a.txt", "test/b.txt"}

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(files, nil).Once()
		mocks.store.On("DeleteObjects", ctx, remoteBucketName, asset.Name).Return(nil).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("ListError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		nowMeta := v1.Now()
		asset.ObjectMeta.DeletionTimestamp = &nowMeta

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(nil, errors.New("nope")).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("Error", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "test-bucket", "https://localhost/test.md")
		nowMeta := v1.Now()
		asset.ObjectMeta.DeletionTimestamp = &nowMeta
		files := []string{"test/a.txt", "test/b.txt"}

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)
		mocks.store.On("ListObjects", ctx, remoteBucketName, asset.Name).Return(files, nil).Once()
		mocks.store.On("DeleteObjects", ctx, remoteBucketName, asset.Name).Return(errors.New("nope")).Once()

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("BucketNotReady", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "notReady", "https://localhost/test.md")
		nowMeta := v1.Now()
		asset.ObjectMeta.DeletionTimestamp = &nowMeta

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(status).To(BeZero())
	})

	t.Run("BucketStatusError", func(t *testing.T) {
		// Given
		g := NewGomegaWithT(t)
		ctx := context.TODO()
		relistInterval := time.Minute
		now := time.Now()
		asset := testData("test-asset", "error", "https://localhost/test.md")
		nowMeta := v1.Now()
		asset.ObjectMeta.DeletionTimestamp = &nowMeta

		handler, mocks := newHandler(relistInterval)
		defer mocks.AssertExpectations(t)

		// When
		status, err := handler.Do(ctx, now, asset, asset.Spec.CommonAssetSpec, asset.Status.CommonAssetStatus)

		// Then
		g.Expect(err).To(HaveOccurred())
		g.Expect(status).To(BeZero())
	})
}

func bucketStatusFinder(ctx context.Context, namespace, name string) (*v1beta1.CommonBucketStatus, bool, error) {
	switch {
	case strings.Contains(name, "notReady"):
		return nil, false, nil
	case strings.Contains(name, "error"):
		return nil, false, errors.New("test-error")
	default:
		return &v1beta1.CommonBucketStatus{
			Phase:      v1beta1.BucketReady,
			URL:        "http://test-url.com/bucket-name",
			RemoteName: remoteBucketName,
		}, true, nil
	}
}

type mocks struct {
	store             *storeMock.Store
	loader            *loaderMock.Loader
	validator         *engineMock.Validator
	mutator           *engineMock.Mutator
	metadataExtractor *engineMock.MetadataExtractor
}

func (m *mocks) AssertExpectations(t *testing.T) {
	m.store.AssertExpectations(t)
	m.loader.AssertExpectations(t)
	m.validator.AssertExpectations(t)
	m.mutator.AssertExpectations(t)
	m.metadataExtractor.AssertExpectations(t)
}

func newHandler(relistInterval time.Duration) (asset.Handler, mocks) {
	mocks := mocks{
		store:             new(storeMock.Store),
		loader:            new(loaderMock.Loader),
		validator:         new(engineMock.Validator),
		mutator:           new(engineMock.Mutator),
		metadataExtractor: new(engineMock.MetadataExtractor),
	}

	handler := asset.New(log, fakeRecorder(), mocks.store, mocks.loader, bucketStatusFinder, mocks.validator, mocks.mutator, mocks.metadataExtractor, relistInterval)

	return handler, mocks
}

func fakeRecorder() record.EventRecorder {
	return record.NewFakeRecorder(20)
}

func testData(assetName, bucketName, url string) *v1beta1.Asset {
	return &v1beta1.Asset{
		ObjectMeta: v1.ObjectMeta{
			Name:       assetName,
			Generation: int64(1),
		},
		Spec: v1beta1.AssetSpec{
			CommonAssetSpec: v1beta1.CommonAssetSpec{
				BucketRef: v1beta1.AssetBucketRef{Name: bucketName},
				Source: v1beta1.AssetSource{
					URL:                      url,
					Mode:                     v1beta1.AssetSingle,
					ValidationWebhookService: make([]v1beta1.AssetWebhookService, 3),
					MutationWebhookService:   make([]v1beta1.AssetWebhookService, 3),
					MetadataWebhookService:   make([]v1beta1.WebhookService, 3),
				},
			},
		},
	}
}

func testDataWithDisplayName(assetName, bucketName, url string, displayName string) *v1beta1.Asset {
	return &v1beta1.Asset{
		ObjectMeta: v1.ObjectMeta{
			Name:       assetName,
			Generation: int64(1),
		},
		Spec: v1beta1.AssetSpec{
			CommonAssetSpec: v1beta1.CommonAssetSpec{
				BucketRef: v1beta1.AssetBucketRef{Name: bucketName},
				Source: v1beta1.AssetSource{
					URL:                      url,
					Mode:                     v1beta1.AssetSingle,
					ValidationWebhookService: make([]v1beta1.AssetWebhookService, 3),
					MutationWebhookService:   make([]v1beta1.AssetWebhookService, 3),
					MetadataWebhookService:   make([]v1beta1.WebhookService, 3),
				},
				DisplayName: displayName,
			},
		},
	}
}
