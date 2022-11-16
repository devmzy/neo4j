package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
)

type QueryReq struct {
	Label string `json:"label" form:"label"`
	Key   string `json:"key" form:"key"`
	Limit string `json:"limit" form:"limit"`
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func accuracyQuery(c *gin.Context) {

}

func main() {
	r := gin.Default()
	r.Use(Cors())
	r.POST("/query", func(c *gin.Context) {
		var req QueryReq
		c.Bind(&req)
		fmt.Println(req.Label, req.Key, req.Limit)
		dbUri := "neo4j://localhost:7687"
		driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "123456", ""))
		if err != nil {
			panic(err)
		}
		c.ContentType()
		ctx := context.Background()
		defer driver.Close(ctx)

		cql := "match p=({name:'" + req.Key + "'})-[r*0..4]-() return p limit " + req.Limit

		if req.Label != "" {
			cql = "match p=(n:" + req.Label + "{name:'" + req.Key + "'})-[r*0..4]-() return p limit " + req.Limit
		}
		if req.Key == "" {
			cql = "match p=(n:" + req.Label + ")-[r*0..4]-() return p limit " + req.Limit
		}
		fmt.Println(cql)
		session := driver.NewSession(ctx, neo4j.SessionConfig{})
		defer session.Close(ctx)
		result, err := session.Run(ctx, cql, nil)
		if err != nil {
			log.Println("Query Run failed: ", err)
		}
		var nodeList []neo4j.Node
		var edgeList []neo4j.Relationship
		recordList, err := result.Collect(ctx)
		if err != nil {
			fmt.Println("record failed", err)
		}
		// fmt.Println(recordList)
		for _, record := range recordList {
			for _, value := range record.Values {
				path := value.(neo4j.Path)
				node := path.Nodes
				edge := path.Relationships
				nodeList = append(nodeList, node...)
				edgeList = append(edgeList, edge...)
			}
		}

		c.JSON(200, gin.H{
			"nodeList": RemoveNode(nodeList),
			"edgeList": RemoveEdge(edgeList),
		})
	})

	r.Run(":89") // listen and serve on 0.0.0.0:8080

}

func RemoveNode(nodeList []neo4j.Node) []neo4j.Node {
	var result []neo4j.Node
	tempMap := make(map[string]int)
	for _, element := range nodeList {
		l := len(tempMap)
		tempMap[element.ElementId] = 0
		if len(tempMap) != l {
			result = append(result, element)
		}
	}
	return result
}

func RemoveEdge(edgeList []neo4j.Relationship) []neo4j.Relationship {
	var result []neo4j.Relationship
	tempMap := make(map[string]int)
	for _, element := range edgeList {
		l := len(tempMap)
		tempMap[element.ElementId] = 0
		if len(tempMap) != l {
			result = append(result, element)
		}
	}
	return result
}
