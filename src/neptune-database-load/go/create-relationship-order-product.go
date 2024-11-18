package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

const (
	defaultDir    = "../../../data/kaggle-brazilian-ecommerce"
	defaultSource = defaultDir + "/olist_order_items_dataset.csv"
	outputDir     = "../../../data/bulk-loader-example-opencypher-format"
	defaultOutput = outputDir + "/relationship-order-to-product.csv"
)

func main() {
	source := flag.String("source", defaultSource, "The full path to the source file")
	output := flag.String("output", defaultOutput, "The full path to the output file")
	flag.Parse()

	file, err := os.Open(*source)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	df := dataframe.ReadCSV(file)

	// Rename the columns to conform to the openCypher syntax for data loading of Relationship types
	renameMap := map[string]string{
		":START_ID(order)":     "order_id",
		":END_ID(product)":     "product_id",
		"price:Double":         "price",
		"freight_value:Double": "freight_value",
	}
	for oldName, newName := range renameMap {
		df = df.Rename(oldName, newName)
	}

	// Create an ID column by manually using a sequence
	numRows := df.Nrow()
	idValues := make([]int, numRows)
	for i := 0; i < numRows; i++ {
		idValues[i] = i + 1
	}
	idSeries := series.New(idValues, series.Int, ":ID")
	df = df.Mutate(idSeries)

	// Create a type column called 'has_item' to identify the relationship type in openCypher
	typeValues := make([]string, numRows)
	for i := 0; i < numRows; i++ {
		typeValues[i] = "has_item"
	}
	typeSeries := series.New(typeValues, series.String, ":TYPE")
	df = df.Mutate(typeSeries)

	// Convert the original shipping_limit_date column which is a string to
	// an ISO 8601 date format
	// Parse the date using a specifc example format "2006-01-02 15:04:05"
	// Note a generic "yyyy-mm-dd hh:mm:ss" cannot be used
	// Format to ISO 8601 "2006-01-02T15:04:05Z"
	dateValues := make([]string, numRows)
	for i, val := range df.Col("shipping_limit_date").Records() {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", val)
		if err != nil {
			log.Fatalf("Failed to parse date: %v", err)
		}
		dateValues[i] = parsedTime.Format(time.RFC3339)
	}
	dateSeries := series.New(dateValues, series.String, "shipping_limit_date:DateTime")
	df = df.Mutate(dateSeries)

	df = df.Select([]string{
		":ID",
		":START_ID(order)",
		":END_ID(product)",
		":TYPE",
		"price:Double",
		"freight_value:Double",
		"shipping_limit_date:DateTime",
	})

	outputFile, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	if err := df.WriteCSV(outputFile); err != nil {
		log.Fatalf("failed to write DataFrame to CSV: %v", err)
	}
	fmt.Println("\nNew DataFrame successfully written to file.")
}
