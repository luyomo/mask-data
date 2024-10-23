package utils

import (
    "flag"
    "fmt"
    "os"
//    "path/filepath"
)


type CommandOptions struct {
    ConfigPath  *string
    DataPath    *string
    OutputPath  *string
}
func ReadOptions() *CommandOptions {
// Define command-line flags
    var opts CommandOptions 
    opts.ConfigPath = flag.String("config", "", "Path to the config file (YAML)")
    opts.DataPath   = flag.String("data", "", "Path to the data file (CSV)")
    opts.OutputPath = flag.String("output", "", "Path to the output directory")

    // Parse the command-line flags
    flag.Parse()

    // Validate the required flags
    if *opts.ConfigPath == "" || *opts.DataPath == "" || *opts.OutputPath == "" {
        fmt.Println("Usage: go_app --config configFile.yaml --data data.csv --output new_folder")
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Check if config file exists
    if _, err := os.Stat(*opts.ConfigPath); os.IsNotExist(err) {
        fmt.Printf("Error: Config file does not exist: %s\n", *opts.ConfigPath)
        os.Exit(1)
    }

    // Check if data file exists
    if _, err := os.Stat(*opts.DataPath); os.IsNotExist(err) {
        fmt.Printf("Error: Data file does not exist: %s\n", *opts.DataPath)
        os.Exit(1)
    }

    // Create the output directory if it doesn't exist
    err := os.MkdirAll(*opts.OutputPath, os.ModePerm)
    if err != nil {
        fmt.Printf("Error: Could not create output directory: %s\n", *opts.OutputPath)
        os.Exit(1)
    }

    fmt.Printf("Config file: %s\n", *opts.ConfigPath)
    fmt.Printf("Data file: %s\n", *opts.DataPath)
    fmt.Printf("Output directory: %s\n", *opts.OutputPath)

    return &opts
}
