package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	res := readFromSTDIN()
	err := writeIntoCSVFile(res)
	if err != nil {
		fmt.Printf("Error While Writing the file!\n%v", err)
	}
	fmt.Print("res.csv has been made sucessfully")

}

func readFromSTDIN() [][]string {
	var res [][]string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord\n")
	for scanner.Scan() {
		if scanner.Text() == "" {
			fmt.Print("EXIT")
			break
		}
		s := checkDomain(scanner.Text())
		res = append(res, s)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Couldnt not read from input: %v\n", err)
	}
	return res
}

func checkDomain(domain string) []string {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string = "NULL", "NULL"
	mxRecord, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecord) > 0 {
		hasMX = true
	}
	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)

	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}

	}
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("ErrorL%v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	var res []string
	res = append(res, domain)
	res = append(res, strconv.FormatBool(hasMX))
	res = append(res, strconv.FormatBool(hasSPF))
	res = append(res, spfRecord)
	res = append(res, strconv.FormatBool(hasDMARC))
	res = append(res, dmarcRecord)

	return res
}

func writeIntoCSVFile(records [][]string) error {
	file, err := os.Create("./res.csv")
	if err != nil {
		return err
	}
	var format [][]string = [][]string{{"domain", "hasMX", "hasSPF", "spfRecord", "hasDMARC", "dmarcRecord"}}
	// var data = [][]string{
	// 	{"Name", "Age", "Occupation"},
	// 	{"Sally", "22", "Nurse"},
	// 	{"Joe", "43", "Sportsman"},
	// 	{"Louis", "39", "Author"},
	// }
	format = append(format, records...)

	w := csv.NewWriter(file)
	err = w.WriteAll(format)
	if err != nil {
		return err
	}
	return nil
}
