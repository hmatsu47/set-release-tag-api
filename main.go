package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/hmatsu47/set-release-tag-api/api"
)

func NewGinSetReleaseTagServer(setReleaseTag *api.SetReleaseTag, port int) *http.Server {
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Swagger specの読み取りに失敗しました\n: %s", err)
		os.Exit(1)
	}

	// Swagger Document 非公開
	swagger.Servers = nil

	// Gin Router 設定
	r := gin.Default()

	// HTTP Request の Validation 設定
	r.Use(middleware.OapiRequestValidator(swagger))

	// Handler 実装
	r = api.RegisterHandlers(r, setReleaseTag)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}
	return s
}

func main() {
	port := flag.Int("port", 18080, "Port for API server")
	flag.Parse()
	// リポジトリ URI・付与するタグはコマンドラインパラメータで取得
	repositoryUri := flag.Arg(0)
	if repositoryUri == "" {
		panic("リポジトリの指定がありません")
	}
	tagName := flag.Arg(1)
	if tagName == "" {
		tagName = "release"
	}
	// Server Instance 生成
	setReleaseTag := api.NewSetReleaseTag(repositoryUri, tagName)
	s := NewGinSetReleaseTagServer(setReleaseTag, *port)
	// 停止まで HTTP Request を処理
	log.Fatal(s.ListenAndServe())
}
