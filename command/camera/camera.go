package camera

import "github.com/gin-gonic/gin"

func TriggerPushNotification(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Done"})
}
