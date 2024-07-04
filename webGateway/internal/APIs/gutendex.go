package apis

type Author struct {
	Name string `json:"name"`
}

type BookListAPIResponse struct {
	Results []struct {
		GutenID int               `json:"id"`
		Title   string            `json:"title"`
		Authors []Author          `json:"authors"`
		Formats map[string]string `json:"formats"`
	} `json:"results"`
}

type BookAPIResponse struct {
	GutenID int               `json:"id"`
	Title   string            `json:"title"`
	Authors []Author          `json:"authors"`
	Formats map[string]string `json:"formats"`
}
