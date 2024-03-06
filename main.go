package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mr-destructive/cncf-api/data"
)

func crawler(linkUrl string) string {
	resp, err := http.Get(linkUrl)
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
	return htmlContent
}

func scrapeScriptTag(htmlContent string) []string {
	// Find script tags with multiple lines of content
	scriptContents := make([]string, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		content := strings.TrimSpace(s.Text())
		if strings.Contains(content, "\n") {
			scriptContents = append(scriptContents, content)
		}
	})
	return scriptContents
}

func parseJsonData(data string) []interface{} {
	first := strings.ReplaceAll(data, "window.baseDS = ", "[")
	second := strings.ReplaceAll(first, ";\n      window.statsDS = ", ",")
	third := strings.ReplaceAll(second, ";", "]")
	final_content := strings.ReplaceAll(third, "\n", "")
	var json_data interface{}
	err := json.Unmarshal([]byte(final_content), &json_data)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	list, ok := json_data.([]interface{})
	if !ok {
		log.Fatal("Error converting data")
	}
	return list
}

func getRegistryData(list []interface{}) *[]data.RegistryItem {
	registryData := make([]data.RegistryItem, 0)
	for _, item := range list[:1] {
		items, ok := item.(map[string]interface{})
		if !ok {
			log.Fatal("Error converting item")
		}
		for _, subitem := range items["items"].([]interface{}) {
			items, ok := subitem.(map[string]interface{})
			registryData = append(registryData, data.RegistryItem{
				ID:          items["id"].(string),
				Name:        items["name"].(string),
				Logo:        items["logo"].(string),
				Category:    items["category"].(string),
				Subcategory: items["subcategory"].(string),
			})
			if !ok {
				log.Fatal("Error converting subitem")
				return nil
			}
			db := data.GetDb()
			err := data.CreateTable(db)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			err = data.InsertData(db, registryData[len(registryData)-1])
			if err != nil {
				log.Fatal(err)
				return nil
			}
			registryData, err = data.GetRegistry(db)
			if err != nil {
				log.Fatal(err)
				return nil
			}
		}
	}
	return &registryData
}

func updateRegistry() {
	cncf_web_app_url := "https://landscape.cncf.io/"
	htmlContent := crawler(cncf_web_app_url)
	scriptContents := scrapeScriptTag(htmlContent)
	parseJsonData(scriptContents[0])
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/update", scrapeUpdateRegistry)
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db := data.GetDb()
	registryData, err := data.GetRegistry(db)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryData)
}

func scrapeUpdateRegistry(w http.ResponseWriter, r *http.Request) {
	updateRegistry()
}
