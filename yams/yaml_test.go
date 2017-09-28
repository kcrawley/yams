package yams

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

var parseValueTests = []struct {
	argument        string
	expectedResult  interface{}
	testDescription string
}{
	{"true", true, "boolean"},
	{"\"true\"", "true", "boolean as string"},
	{"3.4", 3.4, "number"},
	{"\"3.4\"", "3.4", "number as string"},
	{"", "", "empty string"},
}

func TestParseValue(t *testing.T) {
	for _, tt := range parseValueTests {
		assertResultWithContext(t, tt.expectedResult, parseValue(tt.argument), tt.testDescription)
	}
}

func TestRead(t *testing.T) {
	result, _ := read([]string{"examples/sample.yaml", "b.c"})
	assertResult(t, 2, result)
}

func TestReadArray(t *testing.T) {
	result, _ := read([]string{"examples/sample_array.yaml", "[1]"})
	assertResult(t, 2, result)
}

func TestReadString(t *testing.T) {
	result, _ := read([]string{"examples/sample_text.yaml"})
	assertResult(t, "hi", result)
}

func TestOrder(t *testing.T) {
	result, _ := read([]string{"examples/order.yaml"})
	formattedResult, _ := yamlToString(result)
	assertResult(t,
		`version: 3
application: MyApp`,
		formattedResult)
}

func TestNewYaml(t *testing.T) {
	result, _ := newYaml([]string{"b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{b [{c 3}]}]",
		formattedResult)
}

func TestNewYamlArray(t *testing.T) {
	result, _ := newYaml([]string{"[0].cat", "meow"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[[{cat meow}]]",
		formattedResult)
}

func TestUpdateYaml(t *testing.T) {
	result, _ := updateYaml([]string{"examples/sample.yaml", "b.c", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[{a Easy! as one two three} {b [{c 3} {d [3 4]} {e [[{name fred} {value 3}] [{name sam} {value 4}]]}]}]",
		formattedResult)
}

func TestUpdateYamlArray(t *testing.T) {
	result, _ := updateYaml([]string{"examples/sample_array.yaml", "[0]", "3"})
	formattedResult := fmt.Sprintf("%v", result)
	assertResult(t,
		"[3 2 3]",
		formattedResult)
}

func TestUpdateYaml_WithScript(t *testing.T) {
	writeScript = "examples/instruction_sample.yaml"
	_, _ = updateYaml([]string{"examples/sample.yaml"})
}

func TestUpdateYaml_WithUnknownScript(t *testing.T) {
	writeScript = "fake-unknown"
	_, err := updateYaml([]string{"examples/sample.yaml"})
	if err == nil {
		t.Error("Expected error due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, err.Error())
}

func TestNewYaml_WithScript(t *testing.T) {
	writeScript = "examples/instruction_sample.yaml"
	expectedResult := `b:
  c: cat
  e:
  - name: Mike Farah`
	result, _ := newYaml([]string{""})
	actualResult, _ := yamlToString(result)
	assertResult(t, expectedResult, actualResult)
}

func TestNewYaml_WithUnknownScript(t *testing.T) {
	writeScript = "fake-unknown"
	_, err := newYaml([]string{""})
	if err == nil {
		t.Error("Expected error due to unknown file")
	}
	expectedOutput := `open fake-unknown: no such file or directory`
	assertResult(t, expectedOutput, err.Error())
}

func TestWriteProperty(t *testing.T) {
	testYaml := `foo:
  bar: something`
	outputPath := filepath.Join(os.TempDir(), "yaml", "test")
	outputFile := filepath.Join(outputPath, "foo.yml")
	os.MkdirAll(outputPath, os.ModePerm)
	ioutil.WriteFile(outputFile, []byte(testYaml), 0644)
	WriteProperty(outputFile, "foo.bar", "thang")
	result, err := ioutil.ReadFile(outputFile)
	defer os.RemoveAll(outputPath)
	expected := `foo:
  bar: thang`
	assertResult(t, nil, err)
	assertResult(t, expected, string(result))
}

func TestReadProperty(t *testing.T) {
	content, _ := ioutil.ReadFile(filepath.Join("examples", "compose.yml"))
	outputPath := filepath.Join(os.TempDir(), "yaml", "test")
	file := filepath.Join(outputPath, "foo.yml")
	os.MkdirAll(outputPath, os.ModePerm)
	ioutil.WriteFile(file, content, 0644)
	result, err := ReadProperty(file, "services")
	assertResult(t, nil, err)
	results := result.(yaml.MapSlice)
	defer os.RemoveAll(outputPath)
	assertResult(t, "service-one", results[0].Key)
	assertResult(t, "service-two", results[1].Key)
}
