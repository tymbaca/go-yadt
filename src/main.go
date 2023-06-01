package main

import (
	"fmt"
	"templater/generator"
)

var body string = `[
	{
		"filename": "1. Monday",
		"pages": [
			{
				"organisation": "Naturovo",
				"address": "Pcholkovo 48"
			},
			{
				"organisation": "Naturovo",
				"address": "Pcholkovo 48"
			},
			{
				"organisation": "Naturovo",
				"address": "Pcholkovo 48"
			}
		]
	},
	{
		"filename": "2. Tuesday",
		"pages": [
			{
				"organisation": "Naturovo",
				"address": "Pcholkovo 48"
			},
			{
				"organisation": "Naturovo",
				"address": "Pcholkovo 48"
			},
			{
				"organisation": "Naturovo",
				"address": "Pcholkovo 48"
			}
		]
	}
]`

func main() {
	// filegen.NewFileGenerator(templateFileName: "template.docx", json_data: )
	fmt.Println("hello")
	fileGenerator := generator.NewFileGenerator("template.docx", body)
	fileGenerator.GenerateZip("Output.zip")
}
