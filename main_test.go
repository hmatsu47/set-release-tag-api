package main

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/hmatsu47/set-release-tag-api/api"
	"github.com/hmatsu47/set-release-tag-api/testdouble"
	"github.com/stretchr/testify/assert"
)

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
		expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:07:10Z")
		expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:27:02Z")
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
		assert.Equal(t, digest2, imageList[0].Digest)
		assert.Equal(t, expectedTime2, imageList[0].PushedAt)
		assert.Equal(t, repositoryName, imageList[0].RepositoryName)
		assert.Equal(t, size2, imageList[0].Size)
		assert.Equal(t, 1, len(imageList[0].Tags))
		assert.Equal(t, tag2, imageList[0].Tags[0])
		assert.Equal(t, digest1, imageList[1].Digest)
		assert.Equal(t, expectedTime1, imageList[1].PushedAt)
		assert.Equal(t, size1, imageList[1].Size)
		assert.Equal(t, 1, len(imageList[1].Tags))
		assert.Equal(t, tag1, imageList[1].Tags[0])
	})
}

func TestSetReleaseTag2(t *testing.T) {
	// テスト用のパラメーターを生成
	repositoryUri := "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1"
	repositoryName := "repository1"
	registryId := "000000000000"
	maxResults := int32(1000)
	attachTagName := "release"
	selectedTagName := "latest"

	// テスト用の ImageIds を生成
	digest1 := "sha256:4d2653f861f1c4cb187f1a61f97b9af7adec9ec1986d8e253052cfa60fd7372f"
	imageId1 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest1),
		}
	digest2 := "sha256:20b39162cb057eab7168652ab012ae3712f164bf2b4ef09e6541fca4ead3df62"
	tag2 := selectedTagName
	imageId2 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest2),
			ImageTag:    aws.String(tag2),
		}
	var imageIds []types.ImageIdentifier
	imageIds = append(imageIds, imageId1)
	imageIds = append(imageIds, imageId2)

	// テスト用の ImageDetails を生成
	expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:27:02Z")
	expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:07:10Z")
	size1Int64 := int64(10017365)
	imageDetail1 :=
		types.ImageDetail{
			ImageDigest:      aws.String(digest1),
			ImagePushedAt:    aws.Time(expectedTime1),
			ImageSizeInBytes: aws.Int64(size1Int64),
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

	// テスト用の Images（ECR）を生成
	manifest := "{\"test\":\"testtext\"}"
	image1 := types.Image{
		ImageId:        &imageIds[1],
		ImageManifest:  aws.String(manifest),
		RegistryId:     aws.String(registryId),
		RepositoryName: aws.String(repositoryName),
	}
	var images []types.Image
	images = append(images, image1)

	// テストケース
	testParams := testdouble.ParamsECRMock{
		RepositoryName:     repositoryName,
		RegistryId:         registryId,
		ImageIds:           imageIds,
		ImageDetails:       imageDetails,
		MaxResults:         maxResults,
		AttachTagName:      attachTagName,
		SelectedTagName:    selectedTagName,
		Images:             images,
	}
	mockParams := testdouble.MockECRParams{
		ParamsECRMock: testParams,
	}

	t.Run("イメージ取得（モック利用／2つ中1つがタグ付き）", func(t *testing.T) {
		ecrClient := func(t *testing.T) testdouble.MockECRAPI {
			return testdouble.GenerateMockECRAPI(mockParams)
		}
		ctx := context.TODO()
		// ImageList のテスト
		imageList, err := api.ImageList(ctx, ecrClient(t), repositoryUri)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(imageList))
		assert.Equal(t, digest2, imageList[0].Digest)
		assert.Equal(t, expectedTime2, imageList[0].PushedAt)
		assert.Equal(t, repositoryName, imageList[0].RepositoryName)
		assert.Equal(t, size2, imageList[0].Size)
		assert.Equal(t, 1, len(imageList[0].Tags))
		assert.Equal(t, tag2, imageList[0].Tags[0])
		// SetTag のテスト
		err = api.SetTag(ctx, ecrClient(t), repositoryUri, attachTagName, selectedTagName)
		assert.NoError(t, err)
	})
}

func TestSetReleaseTag3(t *testing.T) {
	// テスト用のパラメーターを生成
	repositoryUri := "000000000000.dkr.ecr.ap-northeast-1.amazonaws.com/repository1"
	repositoryName := "repository1"
	registryId := "000000000000"
	maxResults := int32(1000)
	attachTagName := "release"
	selectedTagName := "latest"

	// テスト用の ImageIds を生成
	digest1 := "sha256:4d2653f861f1c4cb187f1a61f97b9af7adec9ec1986d8e253052cfa60fd7372f"
	tag1 := selectedTagName
	imageId1 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest1),
			ImageTag:    aws.String(tag1),
		}
	digest2 := "sha256:20b39162cb057eab7168652ab012ae3712f164bf2b4ef09e6541fca4ead3df62"
	tag2 := "test2"
	imageId2 :=
		types.ImageIdentifier{
			ImageDigest: aws.String(digest2),
		}
	var imageIds []types.ImageIdentifier
	imageIds = append(imageIds, imageId1)
	imageIds = append(imageIds, imageId2)

	// テスト用の ImageDetails を生成
	expectedTime1, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:27:02Z")
	expectedTime2, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2022-09-02T05:07:10Z")
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

	// テスト用の Images（ECR）を生成
	manifest := "{\"test\":\"testtext\"}"
	image1 := types.Image{
		ImageId:        &imageIds[0],
		ImageManifest:  aws.String(manifest),
		RegistryId:     aws.String(registryId),
		RepositoryName: aws.String(repositoryName),
	}
	var images []types.Image
	images = append(images, image1)

	// テストケース
	testParams := testdouble.ParamsECRMock{
		RepositoryName:     repositoryName,
		RegistryId:         registryId,
		ImageIds:           imageIds,
		ImageDetails:       imageDetails,
		MaxResults:         maxResults,
		AttachTagName:      attachTagName,
		SelectedTagName:    selectedTagName,
		Images:             images,
	}
	mockParams := testdouble.MockECRParams{
		ParamsECRMock: testParams,
	}

	t.Run("イメージ取得（モック利用／2つ中2つがタグ付き）", func(t *testing.T) {
		ecrClient := func(t *testing.T) testdouble.MockECRAPI {
			return testdouble.GenerateMockECRAPI(mockParams)
		}
		ctx := context.TODO()
		// ImageList のテスト
		imageList, err := api.ImageList(ctx, ecrClient(t), repositoryUri)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(imageList))
		assert.Equal(t, digest1, imageList[0].Digest)
		assert.Equal(t, expectedTime1, imageList[0].PushedAt)
		assert.Equal(t, repositoryName, imageList[0].RepositoryName)
		assert.Equal(t, size1, imageList[0].Size)
		assert.Equal(t, 1, len(imageList[0].Tags))
		assert.Equal(t, tag1, imageList[0].Tags[0])
		assert.Equal(t, digest2, imageList[1].Digest)
		assert.Equal(t, expectedTime2, imageList[1].PushedAt)
		assert.Equal(t, repositoryName, imageList[1].RepositoryName)
		assert.Equal(t, size2, imageList[1].Size)
		assert.Equal(t, 1, len(imageList[1].Tags))
		assert.Equal(t, tag2, imageList[1].Tags[0])
		// SetTag のテスト
		err = api.SetTag(ctx, ecrClient(t), repositoryUri, attachTagName, selectedTagName)
		assert.NoError(t, err)
	})
}