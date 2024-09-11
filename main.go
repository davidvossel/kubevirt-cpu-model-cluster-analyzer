package main

import (
	"fmt"
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

type metadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type node struct {
	Kind    string   `json:"kind"`
	ObjMeta metadata `json:"metadata"`
}

type nodeList struct {
	Items []node `json:"items"`
}

func main() {

	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("Failed to read stdin: %v", err))
	}

	nodeList := nodeList{}
	err = yaml.Unmarshal(bytes, &nodeList)
	if err != nil {
		panic(fmt.Errorf("Error unmarshalling YAML nodelist: %v", err))
	}

	fmt.Printf("node list length: %d\n", len(nodeList.Items))

	for _, node := range nodeList.Items {
		fmt.Printf("node: %v\n", node)
	}
}
