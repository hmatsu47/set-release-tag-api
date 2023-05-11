package testdouble

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// モックパラメーター
type ECRParams struct {
	RepositoryName  string
	RegistryId      string
	ImageDetails    []types.ImageDetail
	MaxResults      int32
	AttachTagName   string
	SelectedTagName string
	Images          []types.Image
}

// モック生成用
type MockECRParams struct {
	ECRParams ECRParams
}

// モック化
type MockECRAPI struct {
	DescribeImagesAPI MockECRDescribeImagesAPI
	BatchGetImageAPI  MockECRBatchGetImageAPI
	PutImageAPI       MockECRPutImageAPI
}

type MockECRDescribeImagesAPI func(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
type MockECRBatchGetImageAPI func(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error)
type MockECRPutImageAPI func(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error)

func (m MockECRAPI) DescribeImages(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error) {
	return m.DescribeImagesAPI(ctx, params, optFns...)
}

func (m MockECRAPI) BatchGetImage(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error) {
	return m.BatchGetImageAPI(ctx, params, optFns...)
}

func (m MockECRAPI) PutImage(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error) {
	return m.PutImageAPI(ctx, params, optFns...)
}
