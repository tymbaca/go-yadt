package generator

import "encoding/json"

func parseJson(json_data string) ParseData {
	parseData := ParseData{}
	json.Unmarshal([]byte(json_data), &parseData)
	return parseData
}

func (s *FileGenerator) Generate() {}
