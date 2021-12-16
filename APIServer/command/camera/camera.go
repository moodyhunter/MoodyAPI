package camera

import "github.com/gin-gonic/gin"

func TriggerPushNotification(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Done"})
}

func GetCameraState(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "State OK"})
}

func StartCamera(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Started"})
}

func StopCamera(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Stopped"})
}

func GetVideoList(c *gin.Context) {
	c.JSON(200, gin.H{"Message": "Stopped"})
}
