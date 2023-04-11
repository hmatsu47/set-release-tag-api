# set-release-tag-api
Go で内部管理用 API を作るテスト (2)

## `.yaml`ファイルから API コードの枠組みを生成

- こちらを踏襲
  - https://github.com/hmatsu47/select-repository-api

```sh:install
go mod init github.com/hmatsu47/select-repository-api
mkdir internal
cd internal
（作成した`.yaml`ファイルを`internal`内にコピー）
cd ..
mkdir api
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
oapi-codegen -output-config -old-config-style -package=api -generate=types -alias-types internal/set-release-tag.yaml > api/config-types.yaml
oapi-codegen -output-config -old-config-style -package=api -generate=gin,spec -alias-types internal/set-release-tag.yaml > api/config-server.yaml
oapi-codegen -config api/config-types.yaml internal/set-release-tag.yaml > api/types.gen.go
oapi-codegen -config api/config-server.yaml internal/set-release-tag.yaml > api/server.gen.go
go mod tidy
```

## 起動方法

`go run main.go [-port=待機ポート番号（TCP）] 対象ECRリポジトリURI [付与するタグ]`