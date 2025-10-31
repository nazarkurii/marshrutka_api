package languages

import (
	"github.com/gin-gonic/gin"
)

func GinMiddlewear(c *gin.Context) {

	c.Set("lang", defineLanguage(c.GetHeader("Content-Language")))

	c.Next()

}

type Language interface {
	Code() string
	Name() string
}

type language struct {
	name string
	code string
}

func (l language) Code() string {
	return l.code
}
func (l language) Name() string {
	return l.name
}

var english = language{
	name: "English",
	code: "en",
}

var ukrainian = language{
	name: "Ukrainian",
	code: "uk",
}

func defineLanguage(value string) Language {
	switch value {
	case "uk":
		return ukrainian
	default:
		return english
	}
}
