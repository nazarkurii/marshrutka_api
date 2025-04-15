package documents

import (
	ginutil "maryan_api/pkg/ginutils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, s *gin.Engine, client *http.Client) {
	//-----------------------Trip Routes---------------------------------------
	var customer = s.Group("/customer")
	customer.GET("/legal-documents-uk.pdf", ginutil.ServeFileAttachment("../../static/pdf/legal-policy-uk.pdf", "legal-policy-uk.pdf"))
	customer.GET("/legal-documents-en.pdf", ginutil.ServeFileAttachment("../../static/pdf/legal-policy-en.pdf", "legal-policy-en.pdf"))
}
