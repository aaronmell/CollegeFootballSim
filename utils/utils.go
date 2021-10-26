package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ResponseFormat(c *gin.Context, respStatus *Code, data interface{}) {
	if respStatus == nil {
		log.Error().Msg("response status param not found!")
		respStatus = RequestParamError
	}
	c.JSON(respStatus.Status, gin.H{
		"code": respStatus.Code,
		"msg":  respStatus.Message,
		"data": data,
	})
}
