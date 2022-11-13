package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
)

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

// func EdgeQuery(ctx context.Context, driver neo4j.DriverWithContext, Cypher string) ([]any, error) {
//
// 	// var list []neo4j.Relationship
// 	session := driver.NewSession(ctx, neo4j.SessionConfig{})
// 	defer session.Close(ctx)
// 	result, err := session.Run(ctx, Cypher, nil)
// 	if err != nil {
// 		log.Println("Query Run failed: ", err)
// 		return nil, err
// 	}
// 	var list []any
//
// 	for result.Next(ctx) {
// 		record := result.Record()
// 		list = append(list, record.Values)
// 		// if value, ok := record.Get("r"); ok {
//
// 		// relationship := value.(neo4j.Relationship)
// 		// list = append(list, relationship)
// 		//				log.Println("Edgeid:", relationship.Id, ">>>Node:", relationship.StartId, "---", relationship.Type, "--->","Node:",relationship.EndId)
// 		// }
// 	}
// 	if err = result.Err(); err != nil {
// 		return nil, err
// 	}
// 	return list, result.Err()
//
// }

func main() {
	r := gin.Default()
	r.Use(Cors())
	r.GET("/ping", func(c *gin.Context) {
		// nodeData := getNodeData()
		// edgeData := getEdgeData()
		dbUri := "neo4j://localhost:7687"
		driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "123456", ""))
		if err != nil {
			panic(err)
		}
		ctx := context.Background()
		defer driver.Close(ctx)
		// cql := "MATCH (m)-[r*3..]->(n) return m,n,r limit 100"
		cql := "match p=({name:'大渡河'})-[*]->(t) return p limit 60"

		// cql := "MATCH (n:`河流`) RETURN n LIMIT 25"
		session := driver.NewSession(ctx, neo4j.SessionConfig{})
		defer session.Close(ctx)
		// getNodeData()
		nodeList, relationList, err := Query(ctx, driver, cql)
		if err != nil {
			return
		}
		c.JSON(200, gin.H{
			"nodeList": nodeList,
			"edgeList": relationList,
			"data":     relationList,
		})
	})
	r.GET("/query/:key", func(c *gin.Context) {
		key := c.Param("key")
		// nodeData := getNodeData()
		// edgeData := getEdgeData()
		dbUri := "neo4j://localhost:7687"
		driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "123456", ""))
		if err != nil {
			panic(err)
		}
		ctx := context.Background()
		defer driver.Close(ctx)
		// cql := "MATCH (m)-[r*3..]->(n) return m,n,r limit 100"
		// cql := "match p=(n:`" + key + "`)-[*2..3]->(t) return p limit 25"
		cql := "match p=(  {name:'" + key + "'})-[*]->() return p limit 600"

		// cql := "MATCH (n:`河流`) RETURN n LIMIT 25"
		session := driver.NewSession(ctx, neo4j.SessionConfig{})
		defer session.Close(ctx)
		// getNodeData()
		nodeList, relationList, err := Query(ctx, driver, cql)
		if err != nil {
			return
		}
		print(nodeList)
		c.JSON(200, gin.H{
			// "nodeList": nodeList,
			// "edgeList": relationList,
			"data": relationList,
		})
	})
	r.Run(":89") // listen and serve on 0.0.0.0:8080

}

// func getNodeData() []neo4j.Node {
// 	dbUri := "neo4j://localhost:7687"
// 	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "123456", ""))
// 	if err != nil {
// 		panic(err)
// 	}
// 	ctx := context.Background()
// 	defer driver.Close(ctx)
// 	cql := "match p=(s)-[*2..3]->(t) return p limit 25"
// 	// cql := "MATCH (n:`河流`) RETURN n LIMIT 25"
// 	session := driver.NewSession(ctx, neo4j.SessionConfig{})
// 	defer session.Close(ctx)
// 	data, err := NodeQuery(ctx, driver, cql)
// 	for i := 0; i < len(data); i++ {
//
// 		fmt.Println(data[i].Props["name"].(string)) // / OID 是我自己创建 neo4j db entry 的时候，添加的私有属性
// 		fmt.Println(data[i].ElementId)
//
// 	}
// 	return data
// }
//
// func getEdgeData() []any {
// 	dbUri := "neo4j://localhost:7687"
// 	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "123456", ""))
// 	if err != nil {
// 		panic(err)
// 	}
// 	ctx := context.Background()
// 	defer driver.Close(ctx)
// 	// cql := "MATCH ()-[r]->() RETURN r"
//
// 	// cql := "MATCH ()-[r*3..]->(result) return result limit 100"
// 	cql := "MATCH (n:`河流`) RETURN n LIMIT 25"
//
// 	session := driver.NewSession(ctx, neo4j.SessionConfig{})
// 	defer session.Close(ctx)
// 	data, err := EdgeQuery(ctx, driver, cql)
// 	// for i := 0; i < len(data); i++ {
// 	//
// 	// 	fmt.Println(data[i].Props["name"].(string)) // / OID 是我自己创建 neo4j db entry 的时候，添加的私有属性
// 	// 	fmt.Println(data[i].ElementId)
// 	//
// 	// }
// 	return data
// }

// func NodeQuery(ctx context.Context, driver neo4j.DriverWithContext, Cypher string) ([]neo4j.Node, error) {
//
// 	var list []neo4j.Node
// 	session := driver.NewSession(ctx, neo4j.SessionConfig{})
// 	defer session.Close(ctx)
//
// 	result, err := session.Run(ctx, Cypher, nil)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for result.Next(ctx) {
// 		record := result.Record()
// 		if value, ok := record.Get("n"); ok {
// 			node := value.(neo4j.Node)
// 			list = append(list, node)
// 		}
// 	}
// 	if err = result.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return list, result.Err()
//
// }

func Query(ctx context.Context, driver neo4j.DriverWithContext, Cypher string) ([]neo4j.Node, []any, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	result, err := session.Run(ctx, Cypher, nil)
	if err != nil {
		log.Println("Query Run failed: ", err)
		return nil, nil, err
	}
	var list []any
	// var nodeList []neo4j.Node
	// var relationList []neo4j.Relationship

	for result.Next(ctx) {
		record := result.Record()
		// valueList := record.Values.([]interface{})
		// values := record.Values
		if value, ok := record.Get("p"); ok {
			// node := value.(neo4j.Node)
			// nodeList = append(nodeList, node)
			list = append(list, value)

		}
	}
	return nil, list, nil
	// if value, ok := record.Get("m"); ok {
	// 	node := value.(neo4j.Node)
	// 	nodeList = append(nodeList, node)
	// }
	// if value, ok := record.Get("n"); ok {
	// 	node := value.(neo4j.Node)
	// 	nodeList = append(nodeList, node)
	// }
	// if value, ok := record.Get("r"); ok {
	// 	valueList := value.([]interface{})
	// 	relationList = append(relationList, valueList...)
	// }
	// }
	// if err = result.Err(); err != nil {
	// 	return nil, nil, err
	// }
	// return nodeList, relationList, result.Err()

}
