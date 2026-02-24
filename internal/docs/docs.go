package docs

import (
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
	ginredoc "github.com/mvrilo/go-redoc/gin"
)

func InitDoc(g *gin.Engine) {
	doc := redoc.Redoc{
		Title:       "Payment service API",
		Description: "Service de paiment en ligne",
		DocsPath:    "/docs",
		SpecPath:    "/Payment.yaml",
		SpecFile:    "../../internal/docs/OpenApi/Payment.yaml",
	}
	g.Use(ginredoc.New(doc))
}
