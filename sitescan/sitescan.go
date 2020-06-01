package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const nTests = 3
const sleepInSeconds = 5
const FILE_SITES = "sites/sites.txt"
const FILE_LOG_INFO = "logs/log.txt"
const FILE_LOG_ERROR = "logs/error.txt"

func main() {

	printIntroduction()

	for {
		option := getOption()

		switch option {
		case 1:
			startScan()

		case 2:
			fmt.Println("Printing SCAN logs...")
			printFile(FILE_LOG_INFO)

		case 3:
			fmt.Println("Printing ERROR logs...")
			printFile(FILE_LOG_ERROR)

		case 0:
			fmt.Println("Program finished.")
			os.Exit(0)

		default:
			fmt.Println("Option not available.")
			os.Exit(-1)
		}
	}
}

func printIntroduction() {
	name := "Joca"
	age := 36
	version := "1.0.0"

	fmt.Println("Hello world Mr.", name, "your age is", age)
	fmt.Println("This program is in version", version)
}

func getOption() int {
	fmt.Println("\nPlease, select one of the options bellow:")
	fmt.Println("1- Start scan")
	fmt.Println("2- Print scan logs")
	fmt.Println("3- Print error logs")
	fmt.Println("0- Exit")

	option := 0
	fmt.Scan(&option)

	return option
}

func startScan() {
	fmt.Println("Scanning...")

	sites := readSitesFromFile()

	for i := 0; i < nTests; i++ {

		for _, site := range sites {

			resp, err := http.Get(site)
			if err != nil {
				logError(site, err)
			}

			if resp == nil {
				logSiteAccess(site, 0)
			} else {
				logSiteAccess(site, resp.StatusCode)
			}
		}

		time.Sleep(sleepInSeconds * time.Second)
	}

}

func readSitesFromFile() []string {
	sites := []string{}

	fileSites, err := os.Open(FILE_SITES)
	if err != nil {
		fmt.Println(err)
		return sites
	}

	reader := bufio.NewReader(fileSites)
	for {
		line, err := reader.ReadString('\n')

		line = strings.TrimSpace(line)

		sites = append(sites, line)

		if err == io.EOF {
			break
		}
	}

	fileSites.Close()

	return sites
}

func logSiteAccess(site string, statusCode int) {
	logFile, err := os.OpenFile(FILE_LOG_INFO, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now()
	logLine := "[" + now.Format("01/02/2006 - 15:04:05") + "] The access on the site " + site

	if statusCode != 200 {
		logLine = logLine + " was unsuccessfull. Status code: " + strconv.Itoa(statusCode) + "\n"
		logFile.WriteString(logLine)
	} else {
		logLine = logLine + " was done with success. Status code 200\n"
		logFile.WriteString(logLine)
	}

	logFile.Close()
}

func logError(site string, errLog error) {
	if errLog == nil {
		return
	}

	logFile, err := os.OpenFile(FILE_LOG_ERROR, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	now := time.Now()
	logLine := "[" + now.Format("01/02/2006 - 15:04:05") + "] Error when site " + site + " was accessed.\n"

	logFile.WriteString(logLine)
	logFile.WriteString(errLog.Error() + "\n")
	logFile.Close()
}

func printFile(fileName string) {
	infoLog, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(infoLog))
}
