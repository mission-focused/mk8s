package common

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

func CreateConfigFile(content string, fileName string) error {

	var configFile map[string]interface{}

	err := yaml.Unmarshal([]byte(content), &configFile)
	if err != nil {
		return err
	}

	var b bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&b)
	// yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(configFile)

	err = os.WriteFile(fileName, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
