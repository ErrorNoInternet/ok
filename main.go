package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/guptarohit/asciigraph"
	"github.com/peterbourgon/diskv"
)

var okDatabase *diskv.Diskv

func reverseArray(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func main() {
	databasePath := "."
	userObject, errorObject := user.Current()
	if errorObject != nil {
		panic(errorObject)
	}
	if runtime.GOOS == "linux" {
		databasePath = "/home/" + userObject.Username + "/.OkDatabase"
	} else if runtime.GOOS == "windows" {
		databasePath = "C:\\Users\\" + userObject.Username + "\\Documents\\OkDatabase"
	}

	rand.Seed(time.Now().UnixNano())
	currentTime := time.Now()
	timeParts := strings.Split(currentTime.Format("01-02-2016"), "-")
	currentDayInt64, _ := strconv.ParseInt(timeParts[1], 10, 0)
	currentDay := int(currentDayInt64)

	flatTransform := func(input string) []string {
		return []string{}
	}
	okDatabase := diskv.New(diskv.Options{
		BasePath:     databasePath,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	arguments := os.Args[1:]
	showStatistics := false
	resetValues := false
	showHelpPage := false
	for _, argument := range arguments {
		if argument == "stats" || argument == "statistics" {
			showStatistics = true
		} else if argument == "reset" {
			resetValues = true
		} else if argument == "help" {
			showHelpPage = true
		}
	}

	if showHelpPage {
		helpText := "<fg=white;op=bold;>ok</> - ok\n<fg=white;op=bold;>ok stats</> - shows your statistics\n<fg=white;op=bold;>ok reset</> - resets your statistics\n"
		color.Printf(helpText)
	} else if showStatistics {
		keyCount := 0
		for _ = range okDatabase.Keys(make(chan struct{})) {
			keyCount++
		}
		if keyCount == 0 {
			fmt.Println("No statistics...")
			return
		}

		currentCount := 1
		currentCountBytes, errorObject := okDatabase.Read("counter")
		if errorObject == nil {
			currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
			currentCount = int(currentCountInt64)
		}
		numberArray := []float64{}
		captionArray := []string{}
		heatmapArray := []string{}
		highestCount := 1
		keys := okDatabase.Keys(make(chan struct{}))
		for key := range keys {
			if strings.HasPrefix(string(key), "DAY.") {
				dayString := strings.Split(string(key), "DAY.")[1]
				repeatedTimes := 1
				repeatedTimesBytes, errorObject := okDatabase.Read(key)
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
					if strings.HasSuffix(dayString, "1") && dayString != "11" {
						dayUnit = "st"
					} else if strings.HasSuffix(dayString, "2") && dayString != "12" {
						dayUnit = "nd"
					} else if strings.HasSuffix(dayString, "3") && dayString != "13" {
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
			if strings.HasSuffix(caption, "1") && caption != "11" {
				dayUnit = "st"
			} else if strings.HasSuffix(caption, "2") && caption != "12" {
				dayUnit = "nd"
			} else if strings.HasSuffix(caption, "3") && caption != "13" {
				dayUnit = "rd"
			}
			captionText += caption + dayUnit + "    "
		}
		todayCounter := 1
		todayCounterBytes, errorObject := okDatabase.Read("DAY." + strconv.Itoa(currentDay))
		if errorObject == nil {
			todayCounterInt64, _ := strconv.ParseInt(string(todayCounterBytes), 10, 0)
			todayCounter = int(todayCounterInt64)
		}
		graph := "Not enough data..."
		heatmapOutput = heatmapOutput[:len(heatmapOutput)-2]
		if len(numberArray) > 0 {
			graph = asciigraph.Plot(numberArray, asciigraph.Width(20), asciigraph.Height(10), asciigraph.Caption(captionText))
		}
		color.Printf("<fg=white;op=bold;>OK Counter:</> %v\n<fg=white;op=bold;>Records:</> %v\n<fg=white;op=bold;>Graph:</>\n%v\n\nYou've said OK <fg=white;op=bold;>%v times</> today\n", currentCount, heatmapOutput, graph, todayCounter)
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
		currentCountBytes, errorObject := okDatabase.Read("DAY." + strconv.Itoa(currentDay))
		if errorObject == nil {
			currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
			currentCount = int(currentCountInt64)
		}
		okDatabase.Write("DAY."+strconv.Itoa(currentDay), []byte(strconv.Itoa(currentCount+1)))

		currentCount = 1
		currentCountBytes, errorObject = okDatabase.Read("counter")
		if errorObject == nil {
			currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
			currentCount = int(currentCountInt64)
		}
		okDatabase.Write("counter", []byte(strconv.Itoa(currentCount+1)))

		responses := []string{"ok", "ooka booka", "ok", "o k", "oooookaaa booookaaaa", "you said ok", "ok + 1", "ok = ok", "ok ok ok", "ok ok", "ok", "ooka", "booka"}
		randomIndex := rand.Intn(len(responses))
		outputResponse := responses[randomIndex]
		for _, letter := range outputResponse {
			red := uint8(rand.Intn(214) + 42)
			green := uint8(rand.Intn(214) + 42)
			blue := uint8(rand.Intn(214) + 42)
			color.RGB(red, green, blue).Print(string(letter))
		}
		fmt.Println("")
	}
}
