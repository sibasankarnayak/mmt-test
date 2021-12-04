package controllers

import (
	"encoding/json"
	ds "mmt/data_structure"
	"mmt/forms"
	"mmt/model"
	"mmt/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var column = []utils.CustomStrType{
	utils.CustomStrType{
		"column_name": "flight_no",
		"column_type": "string",
	},
	utils.CustomStrType{
		"column_name": "from_code",
		"column_type": "string",
	},
	utils.CustomStrType{
		"column_name": "to_code",
		"column_type": "string",
	},
	utils.CustomStrType{
		"column_name": "start_time",
		"column_type": "time",
	},
	utils.CustomStrType{
		"column_name": "end_time",
		"column_type": "time",
	},
}

//FilghtController ...
type FilghtController struct{}

var flightForm = new(forms.FlightForm)

//Create ...
func (ctrl FilghtController) GetFastestFlight(c *gin.Context) {
	// userID := getUserID(c)

	var form forms.GetFlightForm
	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := flightForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	var flightData []model.FilghtData
	err, vertexMap := utils.ReadCsv("ivtest-sched.csv", column, &flightData)
	if err != nil {
		message := "unable to get flight details. " + err.Error()
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}
	graph := ds.NewGraph()

	for vertex := range vertexMap {
		graph.AddVertex(vertex, nil)
	}
	for _, flight := range flightData {
		timeTaken := utils.TimeDiffMin(flight.StartTime, flight.EndTime, "minute")
		graph.AddEdge(flight.FromCode, flight.ToCode, flight.FlightNo, float64(timeTaken), nil)
	}
	fileData, _ := json.MarshalIndent(flightData, "", "")
	_ = os.WriteFile("flightData.json", fileData, 0644)
	timeTaken, path, code, _ := graph.Yen(form.FromCode, form.ToCode, 5)
	var result []map[string]interface{}
	for index, val := range path {
		result = append(result, map[string]interface{}{
			ds.GetJoinID(val): map[string]interface{}{
				ds.GetJoinID(code[index]): timeTaken[index],
			},
		})
	}
	c.JSON(http.StatusOK, gin.H{"message": "flight path fetched successfully", "data": result})
}
