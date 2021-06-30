package bucket_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kyma-project/rafter/internal/bucket"
	"github.com/kyma-project/rafter/internal/bucket/automock"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateSystemBuckets(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		privatePrefix := "private"
		publicPrefix := "public"
		region := "region"
		cfg := bucket.Config{
			PrivatePrefix: privatePrefix,
			PublicPrefix:  publicPrefix,
			Region:        region,
		}

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(publicPrefix))).Return(false, nil).Once()
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(publicPrefix)), region).Return(nil).Once()
		minioCli.On("SetBucketPolicy", mock.MatchedBy(testBucketNameFn(publicPrefix)), mock.MatchedBy(func(policy string) bool { return true })).Return(nil).Once()
		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(privatePrefix))).Return(false, nil).Once()
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(privatePrefix)), region).Return(nil).Once()
		defer minioCli.AssertExpectations(t)

		// When
		buckets, err := handler.CreateSystemBuckets()

		// Then
		g.Expect(buckets.Private).To(gomega.HavePrefix(privatePrefix))
		g.Expect(buckets.Public).To(gomega.HavePrefix(publicPrefix))
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Exists", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		privatePrefix := "private"
		publicPrefix := "public"
		region := "region"
		cfg := bucket.Config{
			PrivatePrefix: privatePrefix,
			PublicPrefix:  publicPrefix,
			Region:        region,
		}

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(publicPrefix))).Return(true, nil).Once()
		minioCli.On("SetBucketPolicy", mock.MatchedBy(testBucketNameFn(publicPrefix)), mock.MatchedBy(func(policy string) bool { return true })).Return(nil).Once()
		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(privatePrefix))).Return(true, nil).Once()
		defer minioCli.AssertExpectations(t)

		// When
		buckets, err := handler.CreateSystemBuckets()

		// Then
		g.Expect(buckets.Private).To(gomega.HavePrefix(privatePrefix))
		g.Expect(buckets.Public).To(gomega.HavePrefix(publicPrefix))
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Temporary Errors", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		privatePrefix := "private"
		publicPrefix := "public"
		region := "region"
		cfg := bucket.Config{
			PrivatePrefix: privatePrefix,
			PublicPrefix:  publicPrefix,
			Region:        region,
		}
		testErr := errors.New("Test err")

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(publicPrefix))).Return(false, testErr).Once()
		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(publicPrefix))).Return(false, nil).Once()
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(publicPrefix)), region).Return(nil).Once()
		minioCli.On("SetBucketPolicy", mock.MatchedBy(testBucketNameFn(publicPrefix)), mock.MatchedBy(func(policy string) bool { return true })).Return(nil).Once()

		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(privatePrefix))).Return(false, nil).Twice()
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(privatePrefix)), region).Return(testErr).Once()
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(privatePrefix)), region).Return(nil).Once()

		defer minioCli.AssertExpectations(t)

		// When
		buckets, err := handler.CreateSystemBuckets()

		// Then
		g.Expect(buckets.Private).To(gomega.HavePrefix(privatePrefix))
		g.Expect(buckets.Public).To(gomega.HavePrefix(publicPrefix))
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Fatal Errors", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		privatePrefix := "private"
		publicPrefix := "public"
		region := "region"
		cfg := bucket.Config{
			PrivatePrefix: privatePrefix,
			PublicPrefix:  publicPrefix,
			Region:        region,
		}
		testErr := errors.New("Test err")

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		times := 5

		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(publicPrefix))).Return(false, nil).Maybe()
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(publicPrefix)), region).Return(testErr).Maybe()

		minioCli.On("BucketExists", mock.MatchedBy(testBucketNameFn(privatePrefix))).Return(false, nil).Times(times)
		minioCli.On("MakeBucket", mock.MatchedBy(testBucketNameFn(privatePrefix)), region).Return(testErr).Times(times)

		defer minioCli.AssertExpectations(t)

		// When
		_, err := handler.CreateSystemBuckets()

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(err.Error()).To(gomega.ContainSubstring(testErr.Error()))
	})
}

func TestBucketHandler_CheckBuckets(t *testing.T) {
	t.Run("Buckets exist", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		publicBucket := "public-bucket"
		buckets := bucket.SystemBucketNames{
			Private: "private-bucket",
			Public:  publicBucket,
		}
		cfg := bucket.Config{
			Region: "region",
		}
		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", buckets.Private).Return(true, nil).Once()
		minioCli.On("BucketExists", buckets.Public).Return(true, nil).Once()
		minioCli.On("SetBucketPolicy", publicBucket, mock.MatchedBy(func(policy string) bool { return true })).Return(nil).Once()

		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CheckBuckets(buckets)

		// Then
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Create buckets", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		publicBucket := "public-bucket"
		buckets := bucket.SystemBucketNames{
			Private: "private-bucket",
			Public:  publicBucket,
		}
		region := "region"
		cfg := bucket.Config{
			Region: region,
		}
		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", buckets.Private).Return(false, nil).Once()
		minioCli.On("MakeBucket", buckets.Private, region).Return(nil).Once()
		minioCli.On("BucketExists", buckets.Public).Return(false, nil).Once()
		minioCli.On("MakeBucket", buckets.Public, region).Return(nil).Once()
		minioCli.On("SetBucketPolicy", publicBucket, mock.MatchedBy(func(policy string) bool { return true })).Return(nil).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CheckBuckets(buckets)

		// Then
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Checking private bucket error", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		
		buckets := bucket.SystemBucketNames{
			Private: "private-bucket",
			Public:  "public-bucket",
		}
		cfg := bucket.Config{
			Region: "region",
		}
		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", buckets.Private).Return(false, errors.New("test error")).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CheckBuckets(buckets)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
	})

	t.Run("Checking public bucket error", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		buckets := bucket.SystemBucketNames{
			Private: "public-bucket",
			Public:  "private-bucket",
		}
		cfg := bucket.Config{
			Region: "region",
		}
		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", buckets.Private).Return(true, nil).Once()
		minioCli.On("BucketExists", buckets.Public).Return(false, errors.New("test error")).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CheckBuckets(buckets)

		// Then
		g.Expect(err).To(gomega.HaveOccurred())
	})
}

func TestBucketHandler_CreateIfDoesntExist(t *testing.T) {
	t.Run("SuccessRegion", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		bucketName := "bucket"
		region := "region"
		cfg := bucket.Config{
			Region: region,
		}

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", bucketName).Return(false, nil).Once()
		minioCli.On("MakeBucket", bucketName, region).Return(nil).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CreateIfDoesntExist(bucketName, region)

		// Then
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Exists", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		bucketName := "bucket"
		region := "region"
		cfg := bucket.Config{
			Region: region,
		}

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", bucketName).Return(true, nil).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CreateIfDoesntExist(bucketName, region)

		// Then
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})

	t.Run("Error", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		bucketName := "bucket"
		region := "region"
		cfg := bucket.Config{
			Region: region,
		}
		testErr := errors.New("test error")

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", bucketName).Return(false, nil).Once()
		minioCli.On("MakeBucket", bucketName, region).Return(testErr).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CreateIfDoesntExist(bucketName, region)

		// Then
		g.Expect(err.Error()).To(gomega.ContainSubstring(testErr.Error()))
	})

	t.Run("ErrorCheckingIfExists", func(t *testing.T) {
		// Given
		g := gomega.NewGomegaWithT(t)

		bucketName := "bucket"
		region := "region"
		cfg := bucket.Config{
			Region: region,
		}
		testErr := errors.New("test error")

		minioCli := &automock.BucketClient{}
		handler := bucket.NewHandler(minioCli, cfg)

		minioCli.On("BucketExists", bucketName).Return(false, testErr).Once()
		defer minioCli.AssertExpectations(t)

		// When
		err := handler.CreateIfDoesntExist(bucketName, region)

		// Then
		g.Expect(err.Error()).To(gomega.ContainSubstring(testErr.Error()))
	})
}

func testBucketNameFn(prefix string) func(string) bool {
	return func(bucketName string) bool {
		fmt.Println("bucket name", bucketName)
		return strings.HasPrefix(bucketName, prefix)
	}
}
