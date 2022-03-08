package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const downloadFolder = "patreon"

type Config struct {
	SubFolder       string `json:"subFolder"`
	CampaignId      string `json:"campaignId"`
	MonthStart      string `json:"monthStart"`
	MonthEnd        string `json:"monthEnd"`
	Accept          string `json:"accept"`
	ContentType     string `json:"content-type"`
	SecChUaMobile   string `json:"sec-ch-ua-mobile"`
	SecFetchDest    string `json:"sec-fetch-dest"`
	SecFetchMode    string `json:"sec-fetch-mode"`
	SecFetchSite    string `json:"sec-fetch-site"`
	UserAgent       string `json:"user-agent"`
	SecChUaPlatform string `json:"sec-ch-ua-platform"`
	SecChUa         string `json:"sec-ch-ua"`
	Cookie          string `json:"cookie"`
	Interval        int    `json:"interval"`
}

func ReadConfig() Config {
	var config Config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("Error %v", err)
		return config
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	return config
}

func GetPatreonPostsFromMonth(config Config, month string) {
	fmt.Printf("Getting data for month: %v\n", month)
	data, err := GetPatreonPosts(month, config)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	newpath := filepath.Join(".", downloadFolder, config.SubFolder)
	os.MkdirAll(newpath, os.ModePerm)
	for _, entry := range data.Data {
		if entry.Type == "post" {
			folderName := strings.TrimPrefix(entry.Attributes.URL, "https://www.patreon.com/posts/")
			postMedia := map[string]*string{}
			for _, media := range entry.Relationships.Media.Data {
				postMedia[media.ID] = &media.ID
			}
			filesMap, err := GetFilesFromFolder(config.SubFolder, folderName)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			newpath2 := filepath.Join(".", downloadFolder, config.SubFolder, folderName)
			os.MkdirAll(newpath2, os.ModePerm)
			for _, included := range data.Included {
				if postMedia[included.ID] != nil {
					fileName := included.Attributes.FileName
					fullFile := filepath.Join(".", downloadFolder, config.SubFolder, folderName, fileName)
					if filesMap[fullFile] == nil {
						fmt.Printf("Downloading image: %v\n", fullFile)
						err := DownloadPatreonImage(included.Attributes.ImageUrls.Original, fullFile)
						if err != nil {
							fmt.Printf("%v\n", err)
						}
						time.Sleep(time.Duration(config.Interval) * time.Second)
					} else {
						fmt.Printf("SKIPPING: %v, already exists\n", fullFile)
					}
				}
			}
		}
	}
}

func main() {
	config := ReadConfig()
	monthStart, err := time.Parse("2006-01", config.MonthStart)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	monthEnd, err := time.Parse("2006-01", config.MonthEnd)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for d := monthStart; !d.After(monthEnd); d = d.AddDate(0, 1, 0) {
		GetPatreonPostsFromMonth(config, d.Format("2006-01"))
	}
}

func GetFilesFromFolder(subFolder string, folder string) (map[string]*string, error) {
	var files []string
	filesMap := map[string]*string{}
	folder2 := filepath.Join(".", downloadFolder, subFolder, folder)
	err := filepath.Walk(folder2, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return filesMap, err
	}
	for _, file := range files {
		filesMap[file] = &file
	}
	return filesMap, nil
}
