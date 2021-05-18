package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/guptarohit/asciigraph"
	"github.com/howeyc/gopass"
	"github.com/peterbourgon/diskv"
)

type player struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type playerList struct {
	Count   int      `json:"count"`
	Players []player `json:"people"`
}

type githubAssetData struct {
	Name          string `json:"name"`
	URL           string `json:"browser_download_url"`
	DownloadCount int    `json:"download_count"`
}

type githubRelease struct {
	HtmlURL     string            `json:"html_url"`
	TagName     string            `json:"tag_name"`
	ReleaseName string            `json:"name"`
	Prerelease  bool              `json:"prerelease"`
	AssetData   []githubAssetData `json:"assets"`
}

func reverseArray(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func reverseIntArray(arr []int) []int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

var okDatabase *diskv.Diskv
var currentVersion string = "1.4.4"

func main() {
	databasePath := "OkDatabase"
	userObject, errorObject := user.Current()
	if errorObject != nil {
		panic(errorObject)
	}
	if runtime.GOOS == "linux" {
		databasePath = "/home/" + userObject.Username + "/.OkDatabase"
	} else if runtime.GOOS == "windows" {
		databasePath = userObject.HomeDir + "\\OkDatabase"
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
	showPlayerList := false
	submitPlayer := false
	postMessage := false
	receiveMessage := false
	updateProgram := false
	showVersion := false
	leaveLeaderboard := false
	extraText := ""
	for _, argument := range arguments {
		if argument == "stats" || argument == "statistics" {
			showStatistics = true
		} else if argument == "reset" {
			resetValues = true
		} else if argument == "help" {
			showHelpPage = true
		} else if argument == "list" || argument == "leaderboard" || argument == "lb" {
			showPlayerList = true
		} else if argument == "submit" || argument == "join" {
			submitPlayer = true
		} else if argument == "post" || argument == "send" {
			postMessage = true
		} else if argument == "receive" {
			receiveMessage = true
		} else if argument == "update" {
			updateProgram = true
		} else if argument == "version" {
			showVersion = true
		} else if argument == "leave" || argument == "remove" {
			leaveLeaderboard = true
		} else {
			extraText += " " + argument
		}
	}

	if showHelpPage {
		helpText := "<fg=white;op=bold;>ok</> - ok\n<fg=white;op=bold;>ok stats</> - shows your statistics\n<fg=white;op=bold;>ok reset</> - resets your statistics\n<fg=white;op=bold;>ok list</> - shows the OK leaderboard\n<fg=white;op=bold;>ok join</> - join the OK leaderboard\n<fg=white;op=bold;>ok leave</> - leave the OK leaderboard\n<fg=white;op=bold;>ok post</> - post a public message\n<fg=white;op=bold;>ok receive</> - receive a random message\n<fg=white;op=bold;>ok version</> - shows the OK version\n<fg=white;op=bold;>ok update</> - checks for OK updates\n"
		color.Printf(helpText)
	} else if leaveLeaderboard {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Username: ")
		scanner.Scan()
		userInput := scanner.Text()
		fmt.Print("Password: ")
		userPassword := ""
		userPasswordBytes, errorObject := gopass.GetPasswd()
		if errorObject == nil {
			userPassword = string(userPasswordBytes)
		}
		if userInput == "" {
			fmt.Println("Please enter a name!")
			return
		} else if userPassword == "" {
			fmt.Println("Please enter a password!")
			return
		} else {
			fmt.Printf("Leaving leaderboard...")
			httpResponse, errorObject := http.Get(fmt.Sprintf("http://ok-server.herokuapp.com/remove/%v/%v", userInput, userPassword))
			if errorObject != nil {
				fmt.Println("\rFailed to leave leaderboard...")
				return
			}
			response := ""
			responseBytes, errorObject := ioutil.ReadAll(httpResponse.Body)
			if errorObject != nil {
				fmt.Println("\rFailed to leave leaderboard...")
				return
			} else {
				response = string(responseBytes)
			}
			if strings.HasPrefix(response, "ERROR.") {
				errorName := strings.Split(response, "ERROR.")[1]
				fmt.Println("\rError: " + errorName)
			} else {
				fmt.Println("\rSuccessfully removed player from leaderboard")
			}
			return
		}
	} else if showVersion {
		color.Printf("OK Version: <fg=white;op=bold;>%v</>\n", currentVersion)
		return
	} else if updateProgram {
		fmt.Printf("Checking for updates...")
		httpResponse, errorObject := http.Get("https://api.github.com/repos/ErrorNoInternet/ok/releases/latest")
		if errorObject != nil {
			fmt.Println("\rFailed to check for updates...")
			return
		}
		var response githubRelease
		responseBytes, errorObject := ioutil.ReadAll(httpResponse.Body)
		if errorObject != nil {
			fmt.Println("\rFailed to check for updates...")
			return
		}
		_ = json.Unmarshal(responseBytes, &response)
		if strings.Contains(response.TagName, "termux") {
			response.TagName = strings.Replace(response.TagName, "-termux", "", -1)
			response.ReleaseName = strings.Replace(response.ReleaseName, "-termux", "", -1)
		}
		if response.TagName != currentVersion {
			boldTag := "<fg=white;op=bold;>"
			color.Printf("\r%vNew update!</> Version %v%v</>: %v%v</>\nGitHub URL: %v\n", boldTag, boldTag, response.TagName, boldTag, response.ReleaseName, response.HtmlURL)
		} else {
			fmt.Println("\rThere are no new updates...")
		}
		return
	} else if showPlayerList {
		fmt.Print("Fetching leaderboard...")
		httpResponse, errorObject := http.Get("http://ok-server.herokuapp.com/list")
		if errorObject != nil {
			fmt.Println("\rFailed to fetch player list")
			return
		}
		var response playerList
		responseBytes, errorObject := ioutil.ReadAll(httpResponse.Body)
		if errorObject != nil {
			fmt.Println("\rFailed to fetch player list")
			return
		}
		_ = json.Unmarshal(responseBytes, &response)
		numberArray := []int{}
		playerList := make(map[int]string)
		if response.Count > 0 {
			color.Println("\r<fg=white;op=bold;>OK Leaderboard:</>          \n")

			for _, player := range response.Players {
				numberArray = append(numberArray, player.Score)
				playerList[player.Score] = player.Name
			}
			sort.Ints(numberArray)
			numberArray = reverseIntArray(numberArray)
			for index, number := range numberArray {
				playerName := playerList[number]
				color.Printf("<fg=white;op=bold;>%v.</> %v - <fg=white;op=bold;>%v OKs</>\n", index+1, playerName, number)
				if index == 9 {
					return
				}
			}
		} else {
			fmt.Println("\rThere are no players on the OK leaderboard...")
		}
		return
	} else if postMessage {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Message: ")
		scanner.Scan()
		message := scanner.Text()
		fmt.Print("Sending message...")
		_, errorObject := http.Get("http://ok-server.herokuapp.com/send/" + message)
		if errorObject != nil {
			fmt.Println("\rFailed to send message...")
		} else {
			fmt.Println("\rSuccessfully sent message!")
		}
		return
	} else if receiveMessage {
		fmt.Printf("Fetching random message...")
		httpResponse, errorObject := http.Get("http://ok-server.herokuapp.com/message")
		if errorObject != nil {
			fmt.Println("\rFailed to get random message")
			return
		}
		responseBytes, _ := ioutil.ReadAll(httpResponse.Body)
		response := string(responseBytes)
		if strings.HasPrefix(response, "ERROR.") {
			fmt.Println("\rNo one has sent any message...")
		} else {
			color.Println("\r<fg=white;op=bold;>Here's a random message sent by someone:</>\n" + response)
		}
		return
	} else if submitPlayer {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Username: ")
		scanner.Scan()
		userInput := scanner.Text()
		fmt.Print("Password: ")
		userPassword := ""
		userPasswordBytes, errorObject := gopass.GetPasswd()
		if errorObject == nil {
			userPassword = string(userPasswordBytes)
		}
		if userInput == "" {
			fmt.Println("Please enter a name!")
			return
		} else if userPassword == "" {
			fmt.Println("Please enter a password!")
			return
		} else {
			fmt.Printf("Submitting profile...")
			currentCount := 1
			currentCountBytes, errorObject := okDatabase.Read("counter")
			if errorObject == nil {
				currentCountInt64, _ := strconv.ParseInt(string(currentCountBytes), 10, 0)
				currentCount = int(currentCountInt64)
			}
			httpResponse, errorObject := http.Get(fmt.Sprintf("http://ok-server.herokuapp.com/submit/%v/%v/%v/%v", time.Now().Unix(), userInput, userPassword, currentCount))
			if errorObject != nil {
				fmt.Println("\rFailed to submit profile")
				return
			}
			responseBytes, _ := ioutil.ReadAll(httpResponse.Body)
			response := string(responseBytes)

			if strings.HasPrefix(response, "ERROR.") {
				errorName := strings.Split(response, ".")[1]
				fmt.Println("\rFailed to submit profile: " + errorName)
			} else {
				fmt.Println("\rSuccessfully submitted profile to leaderboard!")
			}
			return
		}
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
		if extraText != "" {
			for _, letter := range extraText {
				red := uint8(rand.Intn(214) + 42)
				green := uint8(rand.Intn(214) + 42)
				blue := uint8(rand.Intn(214) + 42)
				color.RGB(red, green, blue).Print(string(letter))
			}
		}
		fmt.Println("")
	}
}
