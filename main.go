package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"sigs.k8s.io/yaml"
)

const (
	cpuModelLabel  = "cpu-model.node.kubevirt.io/"
	hostModelLabel = "host-model-cpu.node.kubevirt.io/"
)

type modelStats struct {
	CPUModel            string `json:"cpuModel"`
	HostModelNodeCount  int    `json:"isHostModel"`
	CompatibleNodeCount int    `json:"compatibleNodeCount"`
}
type results struct {
	CPUModelNodeInfo []modelStats `json:"cpuModelNodeInfo"`
	TotalNodeCount   int          `json:"totalNodeCount"`
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

type BestMatch []modelStats

func (b BestMatch) Len() int      { return len(b) }
func (b BestMatch) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b BestMatch) Less(i, j int) bool {

	if b[i].CompatibleNodeCount > b[j].CompatibleNodeCount {
		return true
	} else if b[i].CompatibleNodeCount == b[j].CompatibleNodeCount &&
		b[i].HostModelNodeCount > b[j].HostModelNodeCount {
		return true
	}
	return false
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
		CPUModelNodeInfo: []modelStats{},
		TotalNodeCount:   len(nodes),
	}

	modelsMap := map[string]modelStats{}

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
			stats, exists := modelsMap[model]
			if !exists {
				stats = modelStats{
					CPUModel: model,
				}
			}

			stats.CompatibleNodeCount++
			if isHostModel {
				stats.HostModelNodeCount++
			}
			modelsMap[model] = stats
		}
	}

	for _, stat := range modelsMap {
		res.CPUModelNodeInfo = append(res.CPUModelNodeInfo, stat)
	}

	sort.Sort(BestMatch(res.CPUModelNodeInfo))

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

	for _, stats := range res.CPUModelNodeInfo {
		fmt.Printf("CPU Model [%s] compatible with [%d] out of [%d] nodes\n", stats.CPUModel, stats.CompatibleNodeCount, res.TotalNodeCount)
		fmt.Printf("CPU Model [%s] hostModel with [%d] out of [%d] nodes\n", stats.CPUModel, stats.HostModelNodeCount, res.TotalNodeCount)
	}
}
