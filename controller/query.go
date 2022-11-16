package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neo4j/logic"
	"neo4j/model"
)

func QueryHandler(c *gin.Context) {
	var req model.QueryReq
	err := c.Bind(&req)
	if err != nil {
		return
	}
	fmt.Println(req.Label, req.Key, req.Limit)
	nodeList, edgeList := logic.Query(req.Label, req.Key, req.Limit)

	c.JSON(200, gin.H{
		"nodeList": nodeList,
		"edgeList": edgeList,
	})
}
