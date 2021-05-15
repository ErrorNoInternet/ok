package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/guptarohit/asciigraph"
	"github.com/prologic/bitcask"
)

var okDatabase *bitcask.Bitcask

func reverseArray(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func main() {
	databasePath := "/home/ryan/Scripts/OkDatabase"

	rand.Seed(time.Now().UnixNano())
	currentTime := time.Now()
	timeParts := strings.Split(currentTime.Format("01-02-2016"), "-")
	currentDayInt64, _ := strconv.ParseInt(timeParts[1], 10, 0)
	currentDay := int(currentDayInt64)
	okDatabase, _ := bitcask.Open(databasePath)

	arguments := os.Args[1:]
	showStatistics := false
	resetValues := false
	for _, argument := range arguments {
		if argument == "stats" || argument == "statistics" {
			showStatistics = true
		} else if argument == "reset" {
			resetValues = true
		}
	}
	if showStatistics {
		if len(okDatabase.Keys()) == 0 {
			fmt.Println("No statistics...")
			return
		}

		currentCount := 1
		currentCountBytes, errorObject := okDatabase.Get([]byte("counter"))
		if errorObject == nil {
			currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
			currentCount = int(currentCountInt64)
		}
		numberArray := []float64{}
		captionArray := []string{}
		heatmapArray := []string{}
		highestCount := 1
		keys := okDatabase.Keys()
		for key := range keys {
			if strings.HasPrefix(string(key), "DAY.") {
				dayString := strings.Split(string(key), "DAY.")[1]
				repeatedTimes := 1
				repeatedTimesBytes, errorObject := okDatabase.Get(key)
				if errorObject == nil {
					repeatedTimesInt64, _ := strconv.ParseInt(string(repeatedTimesBytes), 10, 0)
					repeatedTimes = int(repeatedTimesInt64)
				}

				dayInt := 0
				dayInt64, errorObject := strconv.ParseInt(dayString, 10, 0)
				if errorObject == nil {
					dayInt = int(dayInt64)
				}
				if currentDay-dayInt < 3 {
					captionArray = append(captionArray, dayString)
					numberArray = append(numberArray, float64(repeatedTimes))
				}

				if repeatedTimes > highestCount {
					highestCount = repeatedTimes
					dayUnit := "th"
					if strings.HasSuffix(dayString, "1") {
						dayUnit = "st"
					} else if strings.HasSuffix(dayString, "2") {
						dayUnit = "nd"
					} else if strings.HasSuffix(dayString, "3") {
						dayUnit = "rd"
					}
					heatmapArray = append(heatmapArray, fmt.Sprintf("%v%v - %v times", dayString, dayUnit, repeatedTimes))
				}
			}
		}
		heatmapOutput := ""
		heatmapArray = reverseArray(heatmapArray)
		for index, entry := range heatmapArray {
			if index != 3 {
				heatmapOutput += entry + ", "
			} else {
				break
			}
		}
		captionText := ""
		for _, caption := range captionArray {
			dayUnit := "th"
			if strings.HasSuffix(caption, "1") {
				dayUnit = "st"
			} else if strings.HasSuffix(caption, "2") {
				dayUnit = "nd"
			} else if strings.HasSuffix(caption, "3") {
				dayUnit = "rd"
			}
			captionText += caption + dayUnit + "  "
		}
		heatmapOutput = heatmapOutput[:len(heatmapOutput)-2]
		graph := asciigraph.Plot(numberArray, asciigraph.Width(14), asciigraph.Height(10), asciigraph.Caption(captionText))
		color.Printf("<fg=white;op=bold;>OK Counter:</> %v\n<fg=white;op=bold;>Records:</> %v\n<fg=white;op=bold;>Graph:</>\n%v\n", currentCount, heatmapOutput, graph)
	} else if resetValues {
		scanner := bufio.NewScanner(os.Stdin)
		color.Danger.Println("Are you sure you want to reset all values?")
		fmt.Print("Please enter Y or N: ")
		scanner.Scan()
		userInput := strings.ToLower(scanner.Text())
		if userInput == "y" {
			color.Danger.Print("Please enter Y or N again: ")
			scanner.Scan()
			confirmation := strings.ToLower(scanner.Text())
			if confirmation == "y" {
				errorObject := os.RemoveAll(databasePath)
				if errorObject == nil {
					fmt.Println("\nSuccessfully deleted all values.")
				} else {
					fmt.Println("\nFailed to delete all values.\n" + errorObject.Error())
				}
				return
			} else {
				fmt.Println("\nOperation cancelled.")
				return
			}
		} else {
			fmt.Println("\nOperation cancelled.")
			return
		}
	} else {
		currentCount := 1
		currentCountBytes, errorObject := okDatabase.Get([]byte("DAY." + strconv.Itoa(currentDay)))
		if errorObject == nil {
			currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
			currentCount = int(currentCountInt64)
		}
		okDatabase.Put([]byte("DAY."+strconv.Itoa(currentDay)), []byte(strconv.Itoa(currentCount+1)))

		currentCount = 1
		currentCountBytes, errorObject = okDatabase.Get([]byte("counter"))
		if errorObject == nil {
			currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
			currentCount = int(currentCountInt64)
		}
		okDatabase.Put([]byte("counter"), []byte(strconv.Itoa(currentCount+1)))

		red := uint8(rand.Intn(214) + 42)
		green := uint8(rand.Intn(214) + 42)
		blue := uint8(rand.Intn(214) + 42)
		color.RGB(red, green, blue).Print("o")
		red = uint8(rand.Intn(214) + 42)
		green = uint8(rand.Intn(214) + 42)
		blue = uint8(rand.Intn(214) + 42)
		color.RGB(red, green, blue).Println("k")
	}
}
