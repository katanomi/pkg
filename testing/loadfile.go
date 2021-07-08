package testing

import (
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

// LoadYAML loads yaml
func LoadYAML(file string, obj interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}
	err = yaml.Unmarshal(data, obj)
	return
}
