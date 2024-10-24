package utils

import (
    "fmt"
    "io/ioutil" // To read the YAML file
    "log"

    "gopkg.in/yaml.v2"
)

type MaskRule struct {
    MaskType     string            `yaml:"mask_type"`

    Mapping      map[string]string `yaml:"mapping,omitempty"`       // Optional mapping array
    MappingFile  string            `yaml:"mapping_file,omitempty"`   // Optional mapping file
    Value        string            `yaml:"value,omitempty"`

    // Used for regexp
    MatchStr     string            `yaml:"match_str,omitempty"`
    ReplaceWith  string            `yaml:"replace_with,omitempty"`
}

type Config struct {
    IndexMethod string              `yaml:"index_method"`
    MaskRules   map[string]MaskRule `yaml:"mask_rules"`
}

type ColumnMap struct {
    Config Config `yaml:"config"`
}

func ParseYaml(yamlFileName string) ColumnMap {
    // Read the YAML file
    yamlFile, err := ioutil.ReadFile(yamlFileName)
    if err != nil {
        log.Fatalf("Error reading YAML file: %v", err)
    }

    // Initialize the root config structure
    var rootConfig ColumnMap

    // Unmarshal the YAML file content into the rootConfig structure
    err = yaml.Unmarshal(yamlFile, &rootConfig)
    if err != nil {
        log.Fatalf("Error parsing YAML file: %v", err)
    }

    // Access and print the parsed values
    fmt.Printf("Index Method: %s\n", rootConfig.Config.IndexMethod)

    for column, rule := range rootConfig.Config.MaskRules {
        fmt.Printf("Column: %s, Mask Type: %s\n", column, rule.MaskType)
        if len(rule.Mapping) > 0 {
            fmt.Println("Mappings:")
            for key,  value:= range rule.Mapping {
                fmt.Printf("  Original: %s, Mask Value: %s\n", key, value)
            }
        }
        if rule.MappingFile != "" {
            fmt.Printf("Mapping File: %s\n", rule.MappingFile)
        }
    }

    return rootConfig
}

func ParseMapping(yamlFileName string) (map[string]string, error) {
    // Read the YAML file
    yamlFile, err := ioutil.ReadFile(yamlFileName)
    if err != nil {
        log.Fatalf("Error reading YAML file: %v", err)
        return nil, err
    }

    var columnMap map[string]string

    // Unmarshal the YAML file content into the rootConfig structure
    err = yaml.Unmarshal(yamlFile, &columnMap)
    if err != nil {
        log.Fatalf("Error parsing YAML file: %v", err)
        return nil, err
    }

    return columnMap, nil
}
