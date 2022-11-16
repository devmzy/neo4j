package logic

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"neo4j/util"
)

func Query(label string, key string, limit string) ([]neo4j.Node, []neo4j.Relationship) {
	cql := "match p=({name:'" + key + "'})-[r*0..4]-() return p limit " + limit

	if label != "" {
		cql = "match p=(n:" + label + "{name:'" + key + "'})-[r*0..4]-() return p limit " + limit
	}
	if key == "" {
		cql = "match p=(n:" + label + ")-[r*0..4]-() return p limit " + limit
	}
	nodeList, edgeList := util.Query(cql)
	fmt.Println(cql)
	// session := driver.NewSession(context.TODO(), neo4j.SessionConfig{})
	// defer session.Close(context.TODO())
	// result, err := session.Run(context.TODO(), cql, nil)
	// if err != nil {
	// 	log.Println("Query Run failed: ", err)
	// }
	// var nodeList []neo4j.Node
	// var edgeList []neo4j.Relationship
	// recordList, err := result.Collect(context.TODO())
	// if err != nil {
	// 	fmt.Println("record failed", err)
	// }
	//
	// for _, record := range recordList {
	// 	for _, value := range record.Values {
	// 		path := value.(neo4j.Path)
	// 		node := path.Nodes
	// 		edge := path.Relationships
	// 		nodeList = append(nodeList, node...)
	// 		edgeList = append(edgeList, edge...)
	// 	}
	// }
	return RemoveNode(nodeList), RemoveEdge(edgeList)
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
