package masking

import (
    "fmt"
    "encoding/csv"
    "os"
    "github.com/luyomo/maskdata/pkg/utils"
    "crypto/md5"
    "encoding/hex"
    "path/filepath"
)

var mapColumnFunc = map[string]func(string)string{}

func MaskCSVData(opts *utils.CommandOptions) error {
    columnMapCfg := utils.ParseYaml(*opts.ConfigPath)

    fmt.Printf("The config index method : %#v \n", columnMapCfg.Config.IndexMethod)
    for colName, mapFile := range columnMapCfg.Config.MaskRules {
        fmt.Printf("mapping rule: %s vs %#v \n", colName, mapFile)
	if mapFile.MaskType == "fix_data" {
	    mapColumnFunc[colName] = returnFuncFixValue(mapFile.Value)
        }
	if mapFile.MaskType == "md5" {
	    mapColumnFunc[colName] = returnFuncMD5()
        }
	if mapFile.MaskType == "map" {
	    if mapFile.Mapping != nil {
	        mapColumnFunc[colName] = returnFuncMapConfig(mapFile.Mapping)
            } else if mapFile.MappingFile != "" {
                valueMap, err := utils.ParseMapping(mapFile.MappingFile)
		if err != nil {
                    return err
		}
	        mapColumnFunc[colName] = returnFuncMapConfig(valueMap)
	    }
        }
    }

    // Open the CSV file
    file, err := os.Open(*opts.DataPath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return err
    }
    defer file.Close()

    // Create a new CSV reader
    reader := csv.NewReader(file)

    // Read all the CSV records
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading file:", err)
        return err
    }

    // 002. Prepare the output file
    baseFileName := filepath.Base(*opts.DataPath)

    outputFile, err := os.Create(fmt.Sprintf("%s/%s", *opts.OutputPath, baseFileName))
    if err != nil {
        fmt.Println("Error creating file:", err)
        return err
    }
    defer outputFile.Close()

    // Create a new CSV writer
    fileWriter := csv.NewWriter(outputFile)

    mapCSV := map[int]func(string)string {}

    // header := records[0]
    // for idx, colName := range header {
    //     if theFunc, exists :=  mapColumnFunc[colName]; exists {
    //         mapCSV[idx] = theFunc
    //     }
    // }
    // fmt.Printf("The value: %#v \n", mapCSV)

    // Loop through the records and print them
    for i, record := range records {
        if i == 0 {
            for idx, colName := range record {
                if theFunc, exists :=  mapColumnFunc[colName]; exists {
                    mapCSV[idx] = theFunc
                }
            }
            fmt.Printf("The value: %#v \n", mapCSV)
            err := fileWriter.Write(record)
            if err != nil {
                fmt.Println("Error writing record to file:", err)
                return err
            }
	} else{
            fmt.Printf("Record %d: %v\n", i, record)

            for j, column := range record {
                fmt.Printf("Column: %d: %#v \n", j, column)

                if theFunc, exists := mapCSV[j]; exists {
	            record[j] = theFunc(column)
                    // mapCSV[idx] = theFunc
	        }
	    }

            err := fileWriter.Write(record)
            if err != nil {
                fmt.Println("Error writing record to file:", err)
                return err
            }
        }
    }

    // Flush any buffered data to the file
    fileWriter.Flush()

    // Check if there were any errors during the flush
    if err := fileWriter.Error(); err != nil {
        fmt.Println("Error flushing file:", err)
       return err
    }

    return nil
}

func returnFuncMapConfig(mapV map[string]string) func(string)string {
    return func(theV string) string {
        if newV, exists := mapV[theV]; exists {
            return newV
	}else{
            return "Default"
	}
    }
}

func returnFuncFixValue(inputV string) func(string)string {
    return func(string) string {
        return inputV
    }
}

func returnFuncMD5() func(string)string {
    return func(theV string) string {
        hash := md5.Sum([]byte(theV))

	return hex.EncodeToString(hash[:])
    }
}
