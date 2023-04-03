package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SetReleaseTag struct {
	RepositoryUri string
	TagName       string
}

func NewSetReleaseTag(repositoryUri string, tagName string) *SetReleaseTag {
	return &SetReleaseTag{
		RepositoryUri: repositoryUri,
		TagName:       tagName,
	}
}

// エラーメッセージ返却用
func sendError(c *gin.Context, code int, message string) {
	selectErr := Error{
		Message: message,
	}
	c.JSON(code, selectErr)
}

// コンテナイメージ一覧の取得
func (s *SetReleaseTag) GetImages(c *gin.Context) {
	var result []Image
	region := strings.Split(s.RepositoryUri, ".")[3]
	ecrClient, err := EcrClient(region)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	result, err = ImageList(context.TODO(), ecrClient, s.RepositoryUri)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	c.JSON(http.StatusOK, result)
}

// リリース対象のタグ設定後コンテナイメージ一覧取得
func (s *SetReleaseTag) PostImages(c *gin.Context) {
	var imageTag ImageTag
	err := c.Bind(&imageTag)
	if err != nil {
		sendError(c, http.StatusBadRequest, fmt.Sprintf("パラメーターの形式が誤っています : %s", err))
		return
	}

	// リリースタグ設定
	region := strings.Split(s.RepositoryUri, ".")[3]
	ecrClient, err := EcrClient(region)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	err = SetTag(context.TODO(), ecrClient, s.RepositoryUri, s.TagName, imageTag.Tag)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("タグの設定が失敗しました : %s", err))
		return
	}

	// タグ設定後のコンテナイメージ一覧取得
	var result []Image
	result, err = ImageList(context.TODO(), ecrClient, s.RepositoryUri)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	c.JSON(http.StatusOK, result)
}
