package main

import (
    "fmt"
   // "io/ioutil" // To read the YAML file
   // "log"

    //"gopkg.in/yaml.v2"

    "github.com/luyomo/maskdata/pkg/utils"
    "github.com/luyomo/maskdata/pkg/masking"
)

func main() {

    opts := utils.ReadOptions()
    fmt.Printf("Opts: %#v \n", opts)

    masking.MaskCSVData(opts)
}

