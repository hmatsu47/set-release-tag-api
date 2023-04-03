package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/go-playground/assert/v2"
	"github.com/hmatsu47/set-release-tag-api/api"
)

func doGet(t *testing.T, handler http.Handler, url string) *httptest.ResponseRecorder {
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, handler)
	return response.Recorder
}

// モック化
type ECRClientAPI interface {
	ListImages(ctx context.Context, params *ecr.ListImagesInput, optFns ...func(*ecr.Options)) (*ecr.ListImagesOutput, error)
	DescribeImages(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
	BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
	PutImage(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error)
}

type MockECRClient struct {
	ECRClientAPI
	mockListImages     func(ctx context.Context, params *ecr.ListImagesInput, optFns ...func(*ecr.Options)) (*ecr.ListImagesOutput, error)
	mockDescribeImages func(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
	mockBatchGetImage  func(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
	mockPutImage       func(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error)
}

func (m *MockECRClient) ListImages(ctx context.Context, params *ecr.ListImagesInput, optFns ...func(*ecr.Options)) (*ecr.ListImagesOutput, error) {
	return m.mockListImages(ctx, params, optFns...)
}

func (m *MockECRClient) DescribeImages(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error) {
	return m.mockDescribeImages(ctx, params, optFns...)
}

func (m *MockECRClient) BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error) {
	return m.mockBatchGetImage(ctx, params, optFns...)
}

// go test -v で実行する
func TestSetReleaseTag1(t *testing.T) {
	t.Run("イメージ取得（GetImageListのみ／2つ中1つがタグ付き）", func(t *testing.T) {
		// テスト用の ListImages の結果を生成
		digest1 := "sha256:4d2653f861f1c4cb187f1a61f97b9af7adec9ec1986d8e253052cfa60fd7372f"
		tag1 := "latest"
		imageId1 :=
			types.ImageIdentifier{
				ImageDigest: aws.String(digest1),
				ImageTag:    aws.String(tag1),
			}
		digest2 := "sha256:20b39162cb057eab7168652ab012ae3712f164bf2b4ef09e6541fca4ead3df62"
		imageId2 :=
			types.ImageIdentifier{
				ImageDigest: aws.String(digest2),
			}
		var imageIds []types.ImageIdentifier
		imageIds = append(imageIds, imageId1)
		imageIds = append(imageIds, imageId2)

		// テスト用の DescribeImages の結果を生成
		expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:27:02Z")
		expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:07:10Z")
		registryId := "000000000000"
		repositoryName := "repository1"
		size1 := float32(10017365)
		size1Int64 := int64(10017365)
		var tags1 []string
		tags1 = append(tags1, tag1)
		imageDetail1 :=
			types.ImageDetail{
				ImageDigest:      aws.String(digest1),
				ImagePushedAt:    aws.Time(expectedTime1),
				ImageSizeInBytes: aws.Int64(size1Int64),
				ImageTags:        tags1,
				RegistryId:       aws.String(registryId),
				RepositoryName:   aws.String(repositoryName),
			}
		size2Int64 := int64(10017367)
		imageDetail2 :=
			types.ImageDetail{
				ImageDigest:      aws.String(digest2),
				ImagePushedAt:    aws.Time(expectedTime2),
				ImageSizeInBytes: aws.Int64(size2Int64),
				RegistryId:       aws.String(registryId),
				RepositoryName:   aws.String(repositoryName),
			}
		var imageDetails []types.ImageDetail
		imageDetails = append(imageDetails, imageDetail1)
		imageDetails = append(imageDetails, imageDetail2)

		repositoryUri := "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1"
		imageList := api.GetImageList(imageIds, imageDetails, repositoryName, repositoryUri)
		assert.Equal(t, 1, len(imageList))
		assert.Equal(t, digest1, imageList[0].Digest)
		assert.Equal(t, expectedTime1, imageList[0].PushedAt)
		assert.Equal(t, repositoryName, imageList[0].RepositoryName)
		assert.Equal(t, size1, imageList[0].Size)
		assert.Equal(t, 1, len(imageList[0].Tags))
		assert.Equal(t, tag1, imageList[0].Tags[0])
	})

	t.Run("イメージ取得（GetImageListのみ／2つ中2つがタグ付き）", func(t *testing.T) {
		// テスト用の ListImages の結果を生成
		digest1 := "sha256:4d2653f861f1c4cb187f1a61f97b9af7adec9ec1986d8e253052cfa60fd7372f"
		tag1 := "latest"
		imageId1 :=
			types.ImageIdentifier{
				ImageDigest: aws.String(digest1),
				ImageTag:    aws.String(tag1),
			}
		digest2 := "sha256:20b39162cb057eab7168652ab012ae3712f164bf2b4ef09e6541fca4ead3df62"
		tag2 := "old"
		imageId2 :=
			types.ImageIdentifier{
				ImageDigest: aws.String(digest2),
				ImageTag:    aws.String(tag2),
			}
		var imageIds []types.ImageIdentifier
		imageIds = append(imageIds, imageId1)
		imageIds = append(imageIds, imageId2)

		// テスト用の DescribeImages の結果を生成
		expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:27:02Z")
		expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:07:10Z")
		registryId := "000000000000"
		repositoryName := "repository1"
		size1 := float32(10017365)
		size1Int64 := int64(10017365)
		var tags1 []string
		tags1 = append(tags1, tag1)
		imageDetail1 :=
			types.ImageDetail{
				ImageDigest:      aws.String(digest1),
				ImagePushedAt:    aws.Time(expectedTime1),
				ImageSizeInBytes: aws.Int64(size1Int64),
				ImageTags:        tags1,
				RegistryId:       aws.String(registryId),
				RepositoryName:   aws.String(repositoryName),
			}
		size2 := float32(10017367)
		size2Int64 := int64(10017367)
		var tags2 []string
		tags2 = append(tags2, tag2)
		imageDetail2 :=
			types.ImageDetail{
				ImageDigest:      aws.String(digest2),
				ImagePushedAt:    aws.Time(expectedTime2),
				ImageSizeInBytes: aws.Int64(size2Int64),
				ImageTags:        tags2,
				RegistryId:       aws.String(registryId),
				RepositoryName:   aws.String(repositoryName),
			}
		var imageDetails []types.ImageDetail
		imageDetails = append(imageDetails, imageDetail1)
		imageDetails = append(imageDetails, imageDetail2)

		repositoryUri := "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1"
		imageList := api.GetImageList(imageIds, imageDetails, repositoryName, repositoryUri)
		assert.Equal(t, 2, len(imageList))
		assert.Equal(t, digest1, imageList[0].Digest)
		assert.Equal(t, expectedTime1, imageList[0].PushedAt)
		assert.Equal(t, repositoryName, imageList[0].RepositoryName)
		assert.Equal(t, size1, imageList[0].Size)
		assert.Equal(t, 1, len(imageList[0].Tags))
		assert.Equal(t, tag1, imageList[0].Tags[0])
		assert.Equal(t, digest2, imageList[1].Digest)
		assert.Equal(t, expectedTime2, imageList[1].PushedAt)
		assert.Equal(t, size2, imageList[1].Size)
		assert.Equal(t, 1, len(imageList[1].Tags))
		assert.Equal(t, tag2, imageList[1].Tags[0])
	})
}

func TestSetReleaseTag2(t *testing.T) {
	// t.Run("イメージ取得（モック利用／2つ中1つがタグ付き）", func(t *testing.T) {
	// 	// テスト用の ListImages の結果を生成
	// 	digest1 := "sha256:4d2653f861f1c4cb187f1a61f97b9af7adec9ec1986d8e253052cfa60fd7372f"
	// 	tag1 := "latest"
	// 	imageId1 :=
	// 		types.ImageIdentifier{
	// 			ImageDigest: aws.String(digest1),
	// 			ImageTag:    aws.String(tag1),
	// 		}
	// 	digest2 := "sha256:20b39162cb057eab7168652ab012ae3712f164bf2b4ef09e6541fca4ead3df62"
	// 	imageId2 :=
	// 		types.ImageIdentifier{
	// 			ImageDigest: aws.String(digest2),
	// 		}
	// 	var imageIds []types.ImageIdentifier
	// 	imageIds = append(imageIds, imageId1)
	// 	imageIds = append(imageIds, imageId2)

	// 	// テスト用の DescribeImages の結果を生成
	// 	expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:27:02Z")
	// 	expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:07:10Z")
	// 	registryId := "000000000000"
	// 	repositoryName := "repository1"
	// 	// size1 := float32(10017365)
	// 	size1Int64 := int64(10017365)
	// 	var tags1 []string
	// 	tags1 = append(tags1, tag1)
	// 	imageDetail1 :=
	// 		types.ImageDetail{
	// 			ImageDigest:      aws.String(digest1),
	// 			ImagePushedAt:    aws.Time(expectedTime1),
	// 			ImageSizeInBytes: aws.Int64(size1Int64),
	// 			ImageTags:        tags1,
	// 			RegistryId:       aws.String(registryId),
	// 			RepositoryName:   aws.String(repositoryName),
	// 		}
	// 	size2Int64 := int64(10017367)
	// 	imageDetail2 :=
	// 		types.ImageDetail{
	// 			ImageDigest:      aws.String(digest2),
	// 			ImagePushedAt:    aws.Time(expectedTime2),
	// 			ImageSizeInBytes: aws.Int64(size2Int64),
	// 			RegistryId:       aws.String(registryId),
	// 			RepositoryName:   aws.String(repositoryName),
	// 		}
	// 	var imageDetails []types.ImageDetail
	// 	imageDetails = append(imageDetails, imageDetail1)
	// 	imageDetails = append(imageDetails, imageDetail2)

	// 	// モッククライアントの定義
	// 	mockECR := &MockECRClient{
	// 		mockListImages: func(ctx context.Context, params *ecr.ListImagesInput, optFns ...func(*ecr.Options)) (*ecr.ListImagesOutput, error) {
	// 			return &ecr.ListImagesOutput{
	// 				ImageIds: []types.ImageIdentifier{
	// 					{
	// 						ImageDigest: aws.String(digest1),
	// 						ImageTag:    aws.String(tag1),
	// 					},
	// 					{
	// 						ImageDigest: aws.String("digest2"),
	// 						ImageTag:    aws.String("tag2"),
	// 					},
	// 				},
	// 			}, nil
	// 		},
	// 		mockDescribeImages: func(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error) {
	// 			return &ecr.DescribeImagesOutput{
	// 				ImageDetails: []types.ImageDetail{
	// 					{
	// 						ImageDigest: aws.String("digest1"),
	// 						ImageTags:   []string{"tag1"},
	// 					},
	// 					{
	// 						ImageDigest: aws.String("digest2"),
	// 						ImageTags:   []string{"tag2"},
	// 					},
	// 				},
	// 			}, nil
	// 		},
	// 		mockBatchGetImage: func(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error) {
	// 			return &ecr.BatchGetImageOutput{
	// 				Images: []types.Image{
	// 					{
	// 						ImageId: &types.ImageIdentifier{
	// 							ImageDigest: aws.String("digest1"),
	// 							ImageTag:    aws.String("tag1"),
	// 						},
	// 					},
	// 					{
	// 						ImageId: &types.ImageIdentifier{
	// 							ImageDigest: aws.String("digest2"),
	// 							ImageTag:    aws.String("tag2"),
	// 						},
	// 					},
	// 				},
	// 			}, nil
	// 		},
	// 		mockPutImage: func(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error) {
	// 			return &ecr.PutImageOutput{
	// 				Image: &types.Image{
	// 					ImageId: &types.ImageIdentifier{
	// 						ImageDigest: aws.String("new-digest"),
	// 						ImageTag:    aws.String("new-tag"),
	// 					},
	// 				},
	// 			}, nil
	// 		},
	// 	}
	// })
}
