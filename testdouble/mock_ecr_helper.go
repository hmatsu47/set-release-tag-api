package testdouble

import (
	"context"
	"errors"

	// "fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func GenerateMockECRAPI(mockParams MockECRParams) MockECRAPI {
	return MockECRAPI{
		DescribeImagesAPI: GenerateMockECRDescribeImagesAPI(mockParams),
		BatchGetImageAPI:  GenerateMockECRBatchGetImageAPI(mockParams),
		PutImageAPI:       GenerateMockECRPutImageAPI(mockParams),
	}
}

func GenerateMockECRDescribeImagesAPI(mockParams MockECRParams) MockECRDescribeImagesAPI {
	return MockECRDescribeImagesAPI(func(ctx context.Context, params *ecr.DescribeImagesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error) {
		// fmt.Printf("MockECRDescribeImagesAPI(Expect) : %d / %s / %s\n", mockParams.ECRParams.MaxResults, mockParams.ECRParams.RegistryId, mockParams.ECRParams.RepositoryName)
		// fmt.Printf("MockECRDescribeImagesAPI(Real) :   %d / %s / %s\n", aws.ToInt32(params.MaxResults), aws.ToString(params.RegistryId), aws.ToString(params.RepositoryName))

		if params.MaxResults == nil || aws.ToInt32(params.MaxResults) != mockParams.ECRParams.MaxResults {
			return nil, errors.New("DescribeImagesを呼び出すときのMaxResultsの指定が間違っています")
		}
		if params.RegistryId == nil || aws.ToString(params.RegistryId) != mockParams.ECRParams.RegistryId {
			return nil, errors.New("DescribeImagesを呼び出すときのRegistryIdの指定が間違っています")
		}
		if params.RepositoryName == nil || aws.ToString(params.RepositoryName) != mockParams.ECRParams.RepositoryName {
			return nil, errors.New("DescribeImagesを呼び出すときのRepositoryNameの指定が間違っています")
		}
		detailOutput := &ecr.DescribeImagesOutput{
			ImageDetails: mockParams.ECRParams.ImageDetails,
		}
		return detailOutput, nil
	})
}

func GenerateMockECRBatchGetImageAPI(mockParams MockECRParams) MockECRBatchGetImageAPI {
	return MockECRBatchGetImageAPI(func(ctx context.Context, params *ecr.BatchGetImageInput, optFns ...func(*ecr.Options)) (*ecr.BatchGetImageOutput, error) {
		// fmt.Printf("MockECRBatchGetImageAPI(Expect) : %d / %s / %s\n", 1, mockParams.ECRParams.RegistryId, mockParams.ECRParams.RepositoryName)
		// fmt.Printf("MockECRBatchGetImageAPI(Real) :   %d / %s / %s\n", len(params.ImageIds), aws.ToString(params.RegistryId), aws.ToString(params.RepositoryName))

		if params.ImageIds == nil || len(params.ImageIds) != 1 || params.ImageIds[0].ImageTag == nil || aws.ToString(params.ImageIds[0].ImageTag) != mockParams.ECRParams.SelectedTagName {
			return nil, errors.New("BatchGetImageを呼び出すときのImageIdsの指定が間違っています")
		}
		if params.RegistryId == nil || aws.ToString(params.RegistryId) != mockParams.ECRParams.RegistryId {
			return nil, errors.New("BatchGetImageを呼び出すときのRegistryIdの指定が間違っています")
		}
		if params.RepositoryName == nil || aws.ToString(params.RepositoryName) != mockParams.ECRParams.RepositoryName {
			return nil, errors.New("BatchGetImageを呼び出すときのRepositoryNameの指定が間違っています")
		}

		batchOutput := &ecr.BatchGetImageOutput{
			Images: mockParams.ECRParams.Images,
		}
		return batchOutput, nil
	})
}

func GenerateMockECRPutImageAPI(mockParams MockECRParams) MockECRPutImageAPI {
	return MockECRPutImageAPI(func(ctx context.Context, params *ecr.PutImageInput, optFns ...func(*ecr.Options)) (*ecr.PutImageOutput, error) {
		// fmt.Printf("MockECRPutImageAPI(Expect) : %s / %s / %s / %s\n", aws.ToString(mockParams.ECRParams.Images[0].ImageManifest), mockParams.ECRParams.RegistryId, mockParams.ECRParams.RepositoryName, mockParams.ECRParams.AttachTagName)
		// fmt.Printf("MockECRPutImageAPI(Real) :   %s / %s / %s / %s\n", aws.ToString(params.ImageManifest), aws.ToString(params.RegistryId), aws.ToString(params.RepositoryName), aws.ToString(params.ImageTag))

		if params.ImageManifest == nil || aws.ToString(params.ImageManifest) != aws.ToString(mockParams.ECRParams.Images[0].ImageManifest) {
			return nil, errors.New("PutImageを呼び出すときのImageManifestの指定が間違っています")
		}
		if params.RegistryId == nil || aws.ToString(params.RegistryId) != mockParams.ECRParams.RegistryId {
			return nil, errors.New("PutImageを呼び出すときのRegistryIdの指定が間違っています")
		}
		if params.RepositoryName == nil || aws.ToString(params.RepositoryName) != mockParams.ECRParams.RepositoryName {
			return nil, errors.New("PutImageを呼び出すときのRepositoryNameの指定が間違っています")
		}
		if params.ImageTag == nil || aws.ToString(params.ImageTag) != mockParams.ECRParams.AttachTagName {
			return nil, errors.New("PutImageを呼び出すときのImageTagの指定が間違っています")
		}

		PutImageOutput := &ecr.PutImageOutput{
			Image: &mockParams.ECRParams.Images[0],
		}
		return PutImageOutput, nil
	})
}
