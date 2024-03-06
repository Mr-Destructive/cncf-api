package data

type RegistryItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}
