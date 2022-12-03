package controller

import (
	"github.com/gin-gonic/gin"
	"neo4j/logic"
)

func AddVisit(c *gin.Context) {
	flag := logic.AddVisit()
	c.JSON(200, gin.H{
		"status": flag,
	})
}

func QueryVisit(c *gin.Context) {
	num := logic.QueryVisit()
	c.JSON(200, gin.H{
		"number": num,
	})
}
