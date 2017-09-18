package akamaizer

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	limit         = 8000
	regexProto    = `https?:\/\/`
	domainExt     = `com|it|de|fr|uk|br|jp|ca|au|co|sg`
	regexEndpoint = `^(` + regexProto + `)?([^\/]+\.(` + domainExt + `))`
	regexDomain   = `([^\/]+\.(` + domainExt + `))`
	regexURLext   = `\.[^\/^\s]{2,5}$`
	regexParams   = `\?[^$^\s]+$`
	driverParams  = regexProto + regexDomain + `?([^?]+)(\?.*)?$`
)

func getArrCharCount(arr []string) int {
	var count int
	for i := 0; i < len(arr); i++ {
		count = count + len(arr[i])
	}
	return count
}

func main() {
	importFile := flag.String("csv", "", "csv import file")

	flag.Parse()

	bucket := make(map[string][]string)
	domain := regexp.MustCompile(regexEndpoint)

	ak := csvGo{row: make(chan []string)}
	ak.importFile(*importFile)

	go ak.startLog()

	ak.row <- []string{"redirect", "paths"}

	for i := 0; i < len(ak.getCSV.rows); i++ {
		row := ak.getCSV.rows[i]
		path := strings.TrimSuffix(domain.ReplaceAllString(row["origin"], ""), "/") + "*"
		bucket[row["redirect"]] = append(bucket[row["redirect"]], path)
	}

	for k, v := range bucket {
		if getArrCharCount(v) > limit {
			//	test := make(map[int][]string)
			mapArrbyChar(v, ak.row, k)
		} else if k != "" {
			ak.row <- []string{k, strings.Join(v, "\n")}
		}
	}

	ak.close()
	fmt.Println("\nURLs have been Akamaized")
	fmt.Println()
}

func mapArrbyChar(arr []string, ch chan []string, key string) {
	var count int
	var newArr, newArr2 []string

	for i := 0; i < len(arr); i++ {
		count = count + len(arr[i])
		if getArrCharCount(newArr)+len(arr[i]) < limit {
			newArr = append(newArr, arr[i])
		} else {
			newArr2 = append(newArr2, arr[i])
		}

	}

	ch <- []string{key, strings.Join(newArr, "\n")}
	// temp[len(temp)+1] = newArr

	fmt.Println(key, strconv.Itoa(getArrCharCount(newArr)), newArr)

	if len(newArr2) > 0 {
		mapArrbyChar(newArr2, ch, key)
	}

	//	return temp
}
