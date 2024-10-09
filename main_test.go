package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func TestAbs(t *testing.T) {
	testCases := []string{"test_case_1"}

	for _, testCase := range testCases {
		input, err := readFile("test-assets/" + testCase + "_input.yaml")
		if err != nil {
			t.Errorf("%v", err)
		}
		expectedOutput, err := readFile("test-assets/" + testCase + "_expected_output.yaml")
		if err != nil {
			t.Errorf("%v", err)
		}

		nodeList := nodeList{}
		err = yaml.Unmarshal(input, &nodeList)
		if err != nil {
			t.Errorf("%v", err)
		}

		res := findCommonSupportedModels(nodeList.Items)

		resBytes, err := yaml.Marshal(&res)
		if err != nil {
			t.Errorf("%v", err)
		}

		if strings.Compare(string(expectedOutput), string(resBytes)) != 0 {

			fmt.Printf("\n---EXPECTED---\n%s\n---ACTUAL---\n%s\n---\n", expectedOutput, resBytes)

			t.Errorf("test case [%s] does not match expected results", testCase)
		}

	}
}
