package api

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// ECR クライアント生成
func EcrClient(region string) (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("AWS（API）の認証に失敗しました : %s", err)
	}
	return ecr.NewFromConfig(cfg), nil
}

// ECR ListImages
type EcrListImagesAPI interface {
	ListImages(ctx context.Context, params *ecr.ListImagesInput, optFns ...func(*ecr.Options)) (*ecr.ListImagesOutput, error)
}

func EcrListImages(ctx context.Context, api EcrListImagesAPI, repositoryName string, registryId string) ([]types.ImageIdentifier, error) {
	// ページネーションさせないために最大件数を 1,000 に（実際には数十個程度の想定）
	maxResults := int32(1000)

	ecrImageIds, err := api.ListImages(ctx, &ecr.ListImagesInput{
		RepositoryName: &repositoryName,
		RegistryId:     &registryId,
		MaxResults:     &maxResults,
	})
	if err != nil {
		return nil, fmt.Errorf("リポジトリ（%s）のイメージ一覧の取得に失敗しました : %s", repositoryName, err)
	}
	return ecrImageIds.ImageIds, nil
}

// ECR DescribeImages
type EcrDescribeImagesAPI interface {
	DescribeImages(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
}

func EcrDescribeImages(ctx context.Context, api EcrDescribeImagesAPI, repositoryName string, registryId string) ([]types.ImageDetail, error) {
	// ページネーションさせないために最大件数を 1,000 に（実際には数十個程度の想定）
	maxResults := int32(1000)

	ecrImages, err := api.DescribeImages(ctx, &ecr.DescribeImagesInput{
		RepositoryName: &repositoryName,
		RegistryId:     &registryId,
		MaxResults:     &maxResults,
	})
	if err != nil {
		return nil, fmt.Errorf("リポジトリ（%s）のイメージ詳細一覧の取得に失敗しました : %s", repositoryName, err)
	}
	return ecrImages.ImageDetails, nil
}

// ECR BatchGetImage
type EcrBatchGetImageAPI interface {
	BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
}

func EcrBatchGetImage(ctx context.Context, api EcrBatchGetImageAPI, repositoryName string, registryId string, selectedTagName string) ([]types.Image, error) {
	var imageIds []types.ImageIdentifier
	imageIds = append(imageIds, types.ImageIdentifier{
		ImageTag: &selectedTagName,
	})
	ecrImage, err := api.BatchGetImage(ctx, &ecr.BatchGetImageInput{
		ImageIds:       imageIds,
		RepositoryName: &repositoryName,
		RegistryId:     &registryId,
	})
	if err != nil {
		return nil, fmt.Errorf("リポジトリ（%s）のイメージ情報の取得に失敗しました : %s", repositoryName, err)
	}
	if ecrImage == nil {
		return nil, fmt.Errorf("リポジトリ（%s）のイメージ情報の取得に失敗しました : 対象のイメージ（%s）が存在しません", repositoryName, selectedTagName)
	}

	var images []types.Image
	images = ecrImage.Images
	if len(images) == 0 {
		return nil, fmt.Errorf("リポジトリ（%s）のイメージ情報の取得に失敗しました : 対象のイメージ（%s）が存在しません", repositoryName, selectedTagName)
	}
	return images, nil
}

// ECR PutImage
type EcrPutImageAPI interface {
	PutImage(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error)
}

func EcrPutImage(ctx context.Context, api EcrPutImageAPI, imageManifest string, repositoryName string, registryId string, attachTagName string) error {
	_, err := api.PutImage(ctx, &ecr.PutImageInput{
		ImageManifest:  &imageManifest,
		RepositoryName: &repositoryName,
		ImageTag:       &attachTagName,
		RegistryId:     &registryId,
	})
	return err
}

// ImageList を取得
func GetImageList(imageIds []types.ImageIdentifier, imageDetails []types.ImageDetail, repositoryName string, repositoryUri string) []Image {
	var imageList []Image
	for _, v := range imageDetails {
		tags := v.ImageTags

		if len(tags) > 0 {
			// タグがあるイメージのみ一覧に追加
			digest := v.ImageDigest
			pushedAt := v.ImagePushedAt
			size := v.ImageSizeInBytes
			image := Image{
				Digest:         *digest,
				PushedAt:       *pushedAt,
				RepositoryName: repositoryName,
				Size:           float32(*size),
				Tags:           tags,
			}
			imageList = append(imageList, image)
		}
	}
	// 結果をプッシュ時間の降順でソート
	sort.Slice(imageList, func(i, j int) bool {
		return imageList[i].PushedAt.After(imageList[j].PushedAt)
	})
	return imageList
}

// ECR リポジトリ内イメージ一覧取得
func ImageList(repositoryUri string) ([]Image, error) {
	var err error
	region := strings.Split(repositoryUri, ".")[3]
	repositoryName := strings.Split(repositoryUri, "/")[1]
	registryId := strings.Split(repositoryUri, ".")[0]
	ecrClient, err := EcrClient(region)
	if err != nil {
		return nil, err
	}

	imageIds, err := EcrListImages(context.TODO(), ecrClient, repositoryName, registryId)
	if err != nil {
		return nil, err
	}
	imageDetails, err := EcrDescribeImages(context.TODO(), ecrClient, repositoryName, registryId)
	if err != nil {
		return nil, err
	}

	imageList := GetImageList(imageIds, imageDetails, repositoryName, repositoryUri)
	return imageList, nil
}

// 対象タグを持つイメージにリリースタグを付加
func SetTag(repositoryUri string, attachTagName string, selectedTagName string) error {
	var err error
	region := strings.Split(repositoryUri, ".")[3]
	repositoryName := strings.Split(repositoryUri, "/")[1]
	registryId := strings.Split(repositoryUri, ".")[0]
	ecrClient, err := EcrClient(region)
	if err != nil {
		return err
	}

	var images []types.Image
	images, err = EcrBatchGetImage(context.TODO(), ecrClient, repositoryName, registryId, selectedTagName)
	if err != nil {
		return err
	}

	imageManifest := *images[0].ImageManifest
	err = EcrPutImage(context.TODO(), ecrClient, imageManifest, repositoryName, registryId, attachTagName)
	return err
}