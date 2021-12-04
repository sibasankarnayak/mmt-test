package utils

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type CustomStrType map[string]string

func ReadCsv(fileName string, column []CustomStrType, v interface{}) (error, CustomStrType) {
	vertexMap := make(CustomStrType)
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return err, vertexMap
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
		return err, vertexMap
	}
	var objArr []map[string]interface{}
	for _, line := range csvLines {
		obj := make(map[string]interface{})
		for index, value := range line {
			switch column[index]["column_type"] {
			case "int64":
				if intVal, err := strconv.Atoi(value); err == nil {
					obj[column[index]["column_name"]] = int64(intVal)
				}
			case "string":
				obj[column[index]["column_name"]] = value
				if _, err := strconv.Atoi(value); err != nil {
					vertexMap[value] = ""
				}

			case "time":
				switch len(value) {
				case 0:
					return errors.New("invalid time format found"), vertexMap
				case 1:
					value = "0" + value + "00"
				case 2:
					value = "0" + value + "0"
				case 3:
					value = "0" + value
				}
				layout := "15:04"
				str := value[:2] + ":" + value[2:]
				t, err := time.Parse(layout, str)
				if err != nil {
					return err, vertexMap
				}
				obj[column[index]["column_name"]] = t
			}
		}
		objArr = append(objArr, obj)
	}
	byteArr, _ := json.Marshal(objArr)
	if err != nil {
		return err, vertexMap
	}
	return json.Unmarshal(byteArr, v), vertexMap
}

func TimeDiffMin(t1, t2 time.Time, typeStr string) int {
	diff := t2.Sub(t1)
	if diff.Hours() < 0 {
		t2 = t2.AddDate(0, 0, 1)
		diff = t2.Sub(t1)
	}

	switch typeStr {
	case "hour":
		return int(diff.Hours())
	case "minute":
		return int(diff.Minutes())
	case "second":
		return int(diff.Seconds())
	default:
		return int(diff.Milliseconds())
	}
}

func makePositive(val int) int {
	return val
}
