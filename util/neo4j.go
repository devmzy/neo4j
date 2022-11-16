package util

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/viper"
	"log"
)

func CreateDriver() neo4j.DriverWithContext {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/db.yaml")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置文件失败")
	}
	uri := viper.Get("neo4j.uri").(string)
	username := viper.Get("neo4j.username").(string)
	password := viper.Get("neo4j.password").(string)
	realm := viper.Get("neo4j.realm").(string)

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, realm))
	if err != nil {
		fmt.Println("neo4j连接失败")
	}
	return driver
}

func CloseDriver(driver neo4j.DriverWithContext) error {
	return driver.Close(context.TODO())
}

func Query(cql string) ([]neo4j.Node, []neo4j.Relationship) {
	driver := CreateDriver()
	session := driver.NewSession(context.TODO(), neo4j.SessionConfig{})
	defer session.Close(context.TODO())
	result, err := session.Run(context.TODO(), cql, nil)
	if err != nil {
		log.Println("Query Run failed: ", err)
	}
	var nodeList []neo4j.Node
	var edgeList []neo4j.Relationship
	recordList, err := result.Collect(context.TODO())
	if err != nil {
		fmt.Println("record failed", err)
	}

	for _, record := range recordList {
		for _, value := range record.Values {
			path := value.(neo4j.Path)
			node := path.Nodes
			edge := path.Relationships
			nodeList = append(nodeList, node...)
			edgeList = append(edgeList, edge...)
		}
	}
	return nodeList, edgeList

}
