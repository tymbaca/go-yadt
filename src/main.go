package main

import (
	"fmt"
	"templater/generator"
	"time"
)

var body string = `[
	{
		"filename": "1. Monday",
		"pages": [
			{
				"organisation": "Naturovo1",
				"address": "Pcholkovo 481"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo3",
				"address": "Pcholkovo 483"
			}
		]
	},
	{
		"filename": "2. Вторник",
		"pages": [
			{
				"organisation": "Naturovo1",
				"address": "Pcholkovo 481"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo2",
				"address": "Pcholkovo 482"
			},
			{
				"organisation": "Naturovo3",
				"address": "Pcholkovo 483"
			}
		]
	}
]`

func main() {
	start := time.Now()
	// filenames := []string{"go.mod", "go.sum"}
	// utils.CompressFiles(filenames, "output.zip")

	// filegen.NewFileGenerator(templateFileName: "template.docx", json_data: )

	fmt.Println("hello")
	fileGenerator, err := generator.New("template.docx", body)
	if err != nil {
		panic(err)
	}

	err = fileGenerator.GenerateZip("Output.zip")
	if err != nil {
		panic(err)
	}

	since := time.Since(start)
	fmt.Println(since)
}
