package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

const (
	cpuModelLabel  = "cpu-model.node.kubevirt.io/"
	hostModelLabel = "host-model-cpu.node.kubevirt.io/"
)

type modelStats struct {
	HostModelNodeCount  int `json:"isHostModel"`
	CompatibleNodeCount int `json:"compatibleNodeCount"`
	TotalNodeCount      int `json:"totalNodeCount"`
}
type results struct {
	CPUModels map[string]modelStats
}

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

func parseModelLabel(key string) (string, bool) {
	isHostModel := false
	if strings.Contains(key, cpuModelLabel) {
		//fall through
	} else if strings.Contains(key, hostModelLabel) {
		isHostModel = true
	} else {
		// not a cpu model related label
		return "", false
	}

	str := strings.Split(key, "/")
	if len(str) != 2 {
		panic(fmt.Errorf("invalid cpu model label encountered [%s]", key))
	}
	return str[1], isHostModel
}

func findCommonSupportedModels(nodes []node) results {
	res := results{
		CPUModels: make(map[string]modelStats),
	}

	for _, node := range nodes {
		// TODO ignore nodes that are not virt capable
		for key, val := range node.ObjMeta.Labels {
			if val != "true" {
				continue
			}
			model, isHostModel := parseModelLabel(key)
			if model == "" {
				continue
			}
			stats, exists := res.CPUModels[model]
			if !exists {
				stats = modelStats{
					TotalNodeCount: len(nodes),
				}
			}

			stats.CompatibleNodeCount++
			if isHostModel {
				stats.HostModelNodeCount++
			}
			res.CPUModels[model] = stats

		}
	}

	return res
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

	res := findCommonSupportedModels(nodeList.Items)

	// TODO source CPUModel results by rank

	for model, stats := range res.CPUModels {
		fmt.Printf("CPU Model [%s] compatible with [%d] out of [%d] nodes\n", model, stats.CompatibleNodeCount, stats.TotalNodeCount)
		fmt.Printf("CPU Model [%s] hostModel with [%d] out of [%d] nodes\n", model, stats.HostModelNodeCount, stats.TotalNodeCount)
	}
}
