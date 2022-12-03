package logic

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
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
	log.Println(cql)
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
