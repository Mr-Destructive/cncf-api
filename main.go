package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	handler "github.com/mr-destructive/cncf-landscape-api/api"
	"github.com/mr-destructive/cncf-landscape-api/data"
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
	db := data.GetDb()
	err := data.CreateTable(db)
	if err != nil {
		log.Fatal(err)
		return nil
	}
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
			err = data.InsertData(db, registryData[len(registryData)-1])
			if err != nil {
				log.Fatal(err)
				return nil
			}
			registryData, err = data.GetRegistry(db, "")
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
	list := parseJsonData(scriptContents[0])
	getRegistryData(list)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/update", scrapeUpdateRegistry)
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db := data.GetDb()
	params := r.URL.Query()
	var filter string
	var args []interface{}

	if name := params.Get("name"); name != "" {
		filter += "name LIKE ?"
		args = append(args, "%"+name+"%")
	} else if category := params.Get("category"); category != "" {
		filter += "category LIKE ?"
		args = append(args, "%"+category+"%")
	} else if subcategory := params.Get("subcategory"); subcategory != "" {
		filter += "subcategory LIKE ?"
		args = append(args, "%"+subcategory+"%")
	}

	registryData, err := data.GetRegistry(db, filter, args...)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryData)
	defer db.Close()
}

func scrapeUpdateRegistry(w http.ResponseWriter, r *http.Request) {
	updateRegistry()
	w.WriteHeader(http.StatusOK)
}

func handleRequests(port int) {
	http.HandleFunc("/", handler.Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
