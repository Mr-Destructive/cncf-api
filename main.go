package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	cncf_web_app_url := "https://landscape.cncf.io/"

	// Get the HTML content
	resp, err := http.Get(cncf_web_app_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the HTML content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	htmlContent := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	// Find script tags with multiple lines of content
	scriptContents := make([]string, 0)
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		content := strings.TrimSpace(s.Text())
		if strings.Contains(content, "\n") {
			scriptContents = append(scriptContents, content)
		}
	})

	first := strings.ReplaceAll(scriptContents[0], "window.baseDS = ", "[")
	second := strings.ReplaceAll(first, ";\n      window.statsDS = ", ",")
	third := strings.ReplaceAll(second, ";", "]")
	final_content := strings.ReplaceAll(third, "\n", "")
	var json_data interface{}
	fmt.Println(final_content)
	err = json.Unmarshal([]byte(final_content), &json_data)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	fmt.Println(json_data)
}
