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
	defaultDir    = "../../../data/bulk-loader-example-opencypher-format"
	defaultSource = defaultDir + "/node-olist-orders.csv"
	defaultOutput = defaultDir + "/relationship-customer-to-order.csv"
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
	fmt.Println("DataFrame:")
	fmt.Println(df)

	// Select the columns from the child table (orders) to build the relationship to customers
	newDf := df.Select([]string{
		"order_id:ID(order)",
		"customer_id:String",
		"order_purchase_timestamp:String",
	})
	// Rename the columns to conform to the openCypher syntax for data loading of Relationship types
	renameMap := map[string]string{
		":START_ID(customer)": "customer_id:String",
		":END_ID(order)":      "order_id:ID(order)",
	}
	for oldName, newName := range renameMap {
		newDf = newDf.Rename(oldName, newName)
	}
	// Create an ID column by manually using a sequence
	numRows := newDf.Nrow()
	idValues := make([]int, numRows)
	for i := 0; i < numRows; i++ {
		idValues[i] = i + 1
	}
	idSeries := series.New(idValues, series.Int, ":ID")
	newDf = newDf.Mutate(idSeries)

	// Create a type column called 'ordered' to identify the relationship type in openCypher
	typeValues := make([]string, numRows)
	for i := 0; i < numRows; i++ {
		typeValues[i] = "ordered"
	}
	typeSeries := series.New(typeValues, series.String, ":TYPE")
	newDf = newDf.Mutate(typeSeries)

	// Convert the original order_purchase_timestamp column which is a string to
	// an ISO 8601 date format
	// Parse the date using a specifc example format "2006-01-02 15:04:05"
	// Note a generic "yyyy-mm-dd hh:mm:ss" cannot be used
	// Format to ISO 8601 "2006-01-02T15:04:05Z"
	dateValues := make([]string, numRows)
	for i, val := range newDf.Col("order_purchase_timestamp:String").Records() {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", val)
		if err != nil {
			log.Fatalf("Failed to parse date: %v", err)
		}
		//dateValues = append(dateValues, parsedTime.Format(time.RFC3339))
		dateValues[i] = parsedTime.Format(time.RFC3339)
	}
	dateSeries := series.New(dateValues, series.String, "order_purchase_timestamp:Date")
	newDf = newDf.Mutate(dateSeries)

	newDf = newDf.Select([]string{
		":ID",
		":START_ID(customer)",
		":END_ID(order)",
		":TYPE",
		"order_purchase_timestamp:Date",
	})

	fmt.Println(newDf)

	outputFile, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	if err := newDf.WriteCSV(outputFile); err != nil {
		log.Fatalf("failed to write DataFrame to CSV: %v", err)
	}
	fmt.Println("\nNew DataFrame successfully written to file.")
}
