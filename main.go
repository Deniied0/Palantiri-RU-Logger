package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/apatters/go-conlog"
)

var robloxRUApi = "https://search.roblox.com/catalog/json?CatalogContext=2&Category=6&SortType=3&ResultsPerPage=1"

var lastLoggedId int64 // dont wanna log the same update twice lol

type Errors struct {
	code    int
	message string
}

type ModelData struct {
	AssetId      int64
	Name         string
	Creator      string
	CreatorID    int64
	ThumbnailUrl string
	errors       []Errors
}

func handleError(err error) {
	if err != nil {
		log.Error(err)
	}
}

func getAllHashes(AssetId int64) map[int]string {
	var idString string = strconv.FormatInt(AssetId, 10)
	var version string = "0"
	var hashes = make(map[int]string)
	for {
		var data map[string]interface{}
		i, _ := strconv.Atoi(version)
		version = strconv.Itoa(i + 1)
		resp, err := http.Get("https://assetdelivery.roblox.com/v1/assetId/" + idString + "/version/" + version)
		handleError(err)
		defer resp.Body.Close()

		read, err2 := io.ReadAll(resp.Body)
		handleError(err2)
		err3 := json.Unmarshal(read, &data)
		handleError(err3)
		if data["errors"] != nil {
			break
		}
		log.Info("Obtained hash location " + fmt.Sprint(data["location"]) + "for ID " + idString)
		var wa = hashes
		i2, _ := strconv.Atoi(version)
		wa[i2] = fmt.Sprint(data["location"])
		hashes = wa
	}
	return hashes
}

func GetModel() {
	var data []ModelData
	resp, err0 := http.Get(robloxRUApi)
	handleError(err0)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		read, err := io.ReadAll(resp.Body)
		handleError(err)
		err2 := json.Unmarshal(read, &data)
		handleError(err2)

		if lastLoggedId != data[0].AssetId {
			lastLoggedId = data[0].AssetId

			var idString string = strconv.FormatInt(data[0].AssetId, 10)
			log.Info("Logging model " + idString)

			var hashes = getAllHashes(data[0].AssetId)

			f, err := os.Open("logged.json")
			handleError(err)
			fileRead, err3 := ioutil.ReadAll(f)
			handleError(err3)

			var loggedData map[string]interface{}
			err4 := json.Unmarshal(fileRead, &loggedData)
			handleError(err4)

			var model = make(map[string]interface{})
			model["Name"] = data[0].Name
			model["Creator"] = data[0].Creator
			model["CreatorId"] = data[0].CreatorID
			model["ThumbnailUrl"] = data[0].ThumbnailUrl
			model["VersionHashes"] = hashes

			var assetidstr = strconv.FormatInt(data[0].AssetId, 10)
			loggedData[assetidstr] = model

			var datanew, err5 = json.Marshal(loggedData)
			handleError(err5)

			err6 := os.WriteFile("logged.json", []byte(datanew), 0644)
			handleError(err6)
		}
	}
}

func main() {
	var formatter = log.NewStdFormatter()
	formatter.Options.TimestampType = log.TimestampTypeWall
	log.SetFormatter(formatter)
	log.Info("Palantíri model logger by Deniied")
	log.Info("Starting...")
	for {
		GetModel()
	}
}
