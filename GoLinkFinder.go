/*
 https://0xsha.io
 by @0xsha 1/2020
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	// repos makes things easier
	"github.com/PuerkitoBio/goquery"
	"github.com/akamensky/argparse"
	"github.com/tomnomnom/gahttp"
)

/// Global variables ///

const version = "1.0.0-alpha"

// you can change it
const concurrency = 10

// regex foo from https://github.com/GerbenJavado/LinkFinder
const regexStr = `(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;| *()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\?|#][^"|']{0,}|)))(?:"|')`

// will add everything to this list
// you can change it to SQLite
var founds []string

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func downloadJSFile(urls []string, concurrency int) {

	pipeLine := gahttp.NewPipeline()
	pipeLine.SetConcurrency(concurrency)
	for _, u := range urls {
		pipeLine.Get(u, gahttp.Wrap(parseFile, gahttp.CloseBody))
	}
	pipeLine.Done()
	pipeLine.Wait()

}

func parseFile(req *http.Request, resp *http.Response, err error) {
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)

	// if you like it beautiful
	//options := jsbeautifier.DefaultOptions()
	//code := jsbeautifier.BeautifyFile(string(fileBytes), options)

	matchAndAdd(string(body))

}

func extractUrlFromJS(urls []string, baseUrl string) []string {

	urls = unique(urls)

	var cleaned []string

	for i := 1; i < len(urls); i++ {

		urls[i] = strings.ReplaceAll(urls[i], "'", "")
		urls[i] = strings.ReplaceAll(urls[i], "\"", "")

		if len(urls[i]) < 5 {
			continue
		}

		if !strings.Contains(urls[i], ".js") {
			continue
		}

		if urls[i][:4] == "http" || urls[i][:5] == "https" {
			cleaned = append(cleaned, urls[i])
			continue
		}

		if urls[i][:2] == "//" {

			cleaned = append(cleaned, "https:"+urls[i])
			continue
		}

		if urls[i][:1] == "/" {
			{
				cleaned = append(cleaned, baseUrl+urls[i])

			}

		}
	}
	return cleaned
}

func matchAndAdd(content string) []string {

	regExp, err := regexp.Compile(regexStr)
	if err != nil {
		log.Fatal(err)
	}

	links := regExp.FindAllString(content, -1)
	linksLength := len(links)
	if linksLength > 1 {
		for i := 0; i < linksLength; i++ {
			founds = append(founds, links[i])
		}

	}
	return founds

}

func appendBaseUrl(urls []string, baseUrl string) []string {
	urls = unique(urls)
	var n []string
	for i := 0; i < len(urls); i++ {
		n = append(n, baseUrl+strings.TrimSpace(urls[i]))
	}

	return n
}

func extractJSLinksFromHTML(baseUrl string) []string {

	var resp, err = http.Get(baseUrl)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)

	}

	goos, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var htmlJS = matchAndAdd(goos.Find("script").Text())
	var urls = extractUrlFromJS(htmlJS, baseUrl)

	goos.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		urls = append(urls, src)
	})

	urls = appendBaseUrl(urls, baseUrl)
	return urls
}

func main() {

	parser := argparse.NewParser("goLinkFinder", "GoLinkFinder")
	domain := parser.String("d", "domain", &argparse.Options{Required: true, Help: "Input a URL."})
	output := parser.String("o", "out", &argparse.Options{Required: false, Help: "File name :  (e.g : output.txt)"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))

	}

	var baseUrl = *domain

	if !strings.HasPrefix(baseUrl, "http://") &&
		!strings.HasPrefix(baseUrl, "https://") {
		baseUrl = "https://" + baseUrl

	}

	var htmlUrls = extractJSLinksFromHTML(baseUrl)
	downloadJSFile(htmlUrls, concurrency)
	founds = unique(founds)

	for _, found := range founds {
		found = strings.ReplaceAll(found, "\"", "")
		found = strings.ReplaceAll(found, "'", "")
		fmt.Println(found)
	}

	if *output != "" {

		f, err := os.OpenFile("./"+*output,
			os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()

		if err != nil {
			log.Println(err)
		}

		for _, found := range founds {
			found = strings.ReplaceAll(found, "\"", "")
			found = strings.ReplaceAll(found, "'", "")
			if _, err := f.WriteString(found + "\n"); err != nil {
				log.Fatal(err)
			}
		}

	}

}
