// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// コンテナイメージ一覧の取得
	// (GET /images)
	GetImages(c *gin.Context)
	// リリースタグセット
	// (POST /images)
	PostImages(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetImages operation middleware
func (siw *ServerInterfaceWrapper) GetImages(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetImages(c)
}

// PostImages operation middleware
func (siw *ServerInterfaceWrapper) PostImages(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostImages(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {

	errorHandler := options.ErrorHandler

	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/images", wrapper.GetImages)

	router.POST(options.BaseURL+"/images", wrapper.PostImages)

	return router
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xWS4/cRBD+K6jg6Fk7uywPn3hoBSMRFIXcVquox1O2O3G7e7vbsLOjkcCWAHHhguAC",
	"Bx6C5SFAyiWICH5ME8HPQNX2vD0wQRwizaGn3FX1fV897CkkUihZYmkNxFPQeF6hsS/JMUdv4IJlaG62",
	"ZjIksrRY+iNTquAJs1yW4R0jS7KZJEfB6PSUxhRieDJcZgjbpyYcUtRbLIPZbBbAGE2iuaI4EINrvvO/",
	"B67+xdW/u/pnV//qmsY175O9/snVV/SI/n7qmvdc/SVQFI1GydK0sFFrqW92lv8N9glF7cVcX7nmW8Lc",
	"fE5QCfADV993zQ8e6meuuecPC8DBQtn/gJFbFGYvjSmRnSiEGJjWbNIP/h7Ba951zQeu/spTIPB/3H/7",
	"r6+/2U1hFnSQPJRWm3i6tzJfUJzmewhAaalQ267jBBpDyOPpHLmxmpdZV+TzimscQ3y6uHgWgOW2oJst",
	"iAVlObqDiYUALgbGSlXwLPfq8jHEUDx3KQwf5YejMU988FaxHgr9+uymMOZZNy0bDAJQlclxfJv5p6nU",
	"gk4wZhYHlgtcYl+6aFTScCv15HbJBPaGNfxy9UFZiRFqX3yWmbWW2XLdaI91jb17F34bSTBnusprpRyt",
	"oH3l4KVFXbIC4pQVBvsrlOkqw0iXF2mZqmWFaG/E0w3JbWtcZ9cf9u7h8Vt5epkeX5scddO8QXmLA6Xc",
	"r6suxPGxuDyX51qbIwpOk16mcj7fLPF3UTBeQAy5YNZUTz/7QkaGg0QKCKCtMrzKtZwwUz1xne7k3DAI",
	"oNLezVpl4jDMuM2rEbmF80jwKAv1z4+uXrwxhAAKnmC3hbrs14e39kkXGiwwsYNlawyYUuGokKNQMGNR",
	"h68NXz55/Y0T32udqAbJo0BmcGC9tm+iNi3cawcRXZUKS6Y4xHB0EB1E1GLM5r7WYbs66Zih3X9gu4X2",
	"zo8PP/z44W+fgE+i/ZIdUu1eQTtsI2+8TQ6jaNe2XdwLN/a5L0PKqsL+u+v628ov1koIpiePRKYd9dP2",
	"xQJntGyk6ZNndz9sSXJDmlVN5h8Hk92cVr4fwvWPh9njpOo/arCppPdFTT0K8el0ZSriMCxkwopcGhsf",
	"RVEEs7OF/759uRz5NiGF6N0tF5U6fCa6+zzP0hxms78DAAD//4XE0RG9CQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
