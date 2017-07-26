package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

func main() {
	subreddits := flag.String("subreddits", "accidentalwesanderson", "Comma separated list of subreddits to pull photos from")
	t := flag.String("t", "day", "Time frame to search. Can either be hour, day, week, month, year, or all")
	sort := flag.String("sort", "top", "Default sorting method of posts. Can be hot, new, rising, top, controversial, gilded, or promoted")
	width := flag.Float64("w", 0, "Min width of image to use")
	height := flag.Float64("h", 0, "Min height of image to use")
	limit := flag.Float64("limit", 100, "Number of posts to search")
	verbose := flag.Bool("v", false, "Verbose mode. Causes debugging information to be output")
	dir := flag.String("d", os.TempDir(), "Directory to save images")
	r := flag.Bool("r", false, "Randomly select an image from response")
	flag.Parse()

	subs := split(*subreddits)

	url := fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?t=%s&limit=%v&raw_json=1", subs, *sort, *t, *limit)

	if *verbose {
		log.Printf("Making request to %s for images", url)
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	request.Header.Set("User-Agent", "reddit-bg/0.1")

	ctx := context.Background()
	request = request.WithContext(ctx)

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var content response
	err = json.Unmarshal(body, &content)
	if err != nil {
		log.Fatalln(err.Error())
	}

	post, err := content.selectPost(*width, *height, *r)
	if err != nil {
		log.Fatalln(err.Error())
	}

	u, err := post.getImageURL()
	if err != nil {
		log.Fatalln(err.Error())
	}

	imageRequest, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	imageResponse, err := client.Do(imageRequest)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer imageResponse.Body.Close()
	imageBody, err := ioutil.ReadAll(imageResponse.Body)

	filePath := path.Join(*dir, normalizeTitle(post.Title))

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer file.Close()

	if *verbose {
		log.Printf("Saving image to: %s \n", filePath)
	}

	_, err = file.Write(imageBody)
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", fmt.Sprintf("file://%s", filePath))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func normalizeTitle(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}
	s = reg.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "-", -1)
	return s
}

func split(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, ",", "+", -1)
	return s
}
