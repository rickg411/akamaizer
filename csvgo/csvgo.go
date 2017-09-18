package csvgo

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type csvGo struct {
	fileName string
	count    int
	quit     chan bool
	row      chan []string
	getCSV
}

type getCSV struct {
	path string
	rows map[int]map[string]string
}

//*** Log Functions ***//
func (c *csvGo) startLog() {

	appName := strings.Replace(os.Args[0], "./", "", -1)

	fpArr := []string{os.Getenv("HOME"), "Desktop", appName}

	fp := strings.Join(fpArr, "/")

	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		os.Mkdir(fp, 0775)
	}

	if c.fileName != "" {
		fpArr = append(fpArr, c.fileName+"-")
	}

	fpArr = append(fpArr, time.Now().Format(time.RFC850)+".log.csv")

	fp = strings.Join(fpArr, "/")

	f, err := os.Create(fp)
	if err != nil {
		fmt.Println("Could not Create logfile", fp)
		os.Exit(1)
	}

	c.addLog(f)
}

func (c *csvGo) addLog(f *os.File) {
	w := csv.NewWriter(f)
	w.UseCRLF = false
	defer f.Close()
	for {
		c.count++
		w.Write(<-c.row)
		w.Flush()
	}
}

func (c *csvGo) close() int {
	close(c.row)
	return c.count
}

// Converting CSV file to HashMap
func (c *getCSV) importFile(fp string) map[int]map[string]string {

	csvRows := make(map[int]map[string]string)
	var titles []string

	// Validating if option is CSV file and exists
	validFile, _ := regexp.MatchString(`\.csv$`, fp)
	if validFile {
		_, err := os.Stat(fp)
		if err != nil {
			fmt.Println("\nImport file does not exists:", fp)
			fmt.Println(err.Error() + "\n")
			os.Exit(1)
		}
	} else {
		fmt.Println("\nImport file must be in CSV:", fp+"\n")
		os.Exit(1)
	}

	f, _ := os.Open(fp)
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.FieldsPerRecord = -1

	d, err := r.ReadAll()
	if err != nil {
		panic("Can't read csv file\n")
	}

	for i := 0; i < len(d[0]); i++ {
		titles = append(titles, strings.ToLower(d[0][i]))
	}

	// Converting CSV file to map
	for i, val := range d {
		// Skipping first row/ title row of csv
		if i > 0 {
			csvRow := make(map[string]string)
			for t, title := range titles {
				csvRow[title] = val[t]
			}
			csvRows[i] = csvRow
		}
	}
	c.rows = csvRows
	return csvRows
}
