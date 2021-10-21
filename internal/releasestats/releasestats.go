/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package releasestats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type release struct {
	Name          string `json:"name"`
	Publised_Date string `json:"published_at"`
	Assets        []assets
}

type assets struct {
	Name  string `json:"name"`
	Count int    `json:"download_count"`
}

// GetStats is the main function
func GetStats(repoLink string) {
	link, err := url.ParseRequestURI(repoLink)
	if err != nil {
		panic(err)
	}

	response, err := comms(splitter(fmt.Sprint(link)))
	if err != nil {
		panic(err)
	}

	rel := []release{}
	err = json.Unmarshal([]byte(response), &rel)
	if err != nil {
		panic(err)
	}

	totalCount := 0

	for _, v1 := range rel {
		for _, a1 := range v1.Assets {
			totalCount += a1.Count
		}
	}
	fmt.Println("============================================================")
	fmt.Println("Repository:", repoLink)
	fmt.Println("Total Downloads:", totalCount)
	fmt.Println("============================================================")

	for _, v := range rel {
		count := 0
		for _, a := range v.Assets {
			count += a.Count
		}
		fmt.Println("Name:      ", v.Name)
		fmt.Println("Published: ", v.Publised_Date)
		fmt.Println("Downloads: ", count)
		for _, a1 := range v.Assets {
			fmt.Println("--------------------------------------")
			fmt.Println("Release Name:", a1.Name)
			fmt.Println("Downloads:", a1.Count)
		}
		fmt.Println("============================================================")
	}

}

func splitter(repoLink string) string {
	s := strings.Split(strings.Split(fmt.Sprint(repoLink), "//")[1], "/")
	return "https://api." + s[0] + "/repos/" + s[1] + "/" + s[2] + "/releases"
}

func comms(apiLink string) (string, error) {
	req, err := http.Get(apiLink)
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"Accept": []string{"application/vnd.github.v3+json"},
	}

	defer req.Body.Close()

	if req.StatusCode != 200 {
		return "", fmt.Errorf("unable to fetch repo information: %v", req.StatusCode)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
