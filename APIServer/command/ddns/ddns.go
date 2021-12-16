package ddns

import "github.com/gin-gonic/gin"

func Get(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get"})
}

func List(c *gin.Context) {
	c.JSON(200, gin.H{"message": "List"})
}
func Update(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create"})
}

func Delete(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete"})
}
