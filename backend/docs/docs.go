package docs

import "github.com/swaggo/swag"

const docTemplate = `{"swagger":"2.0","info":{"title":"Fitness API","version":"1.0"},"basePath":"/","paths":{"/health":{"get":{"responses":{"200":{"description":"ok"}}}}}}`

type swaggerInfo struct{}

func (s *swaggerInfo) ReadDoc() string { return docTemplate }

func init() { swag.Register(swag.Name, &swaggerInfo{}) }
