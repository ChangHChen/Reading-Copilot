package apis

type ChatRequest struct {
	Model       string `json:"model"`
	APIKey      string `json:"apikey"`
	BookID      int    `json:"bookID"`
	Progress    int    `json:"progress"`
	UserMessage string `json:"userMessage"`
}

type ChatResponse struct {
	ResponseMessage string `json:"responseMessage"`
	Error           string `json:"error"`
}

type BuildBookRequest struct {
	BookID int `json:"bookID"`
}

type BuildBookResponse struct {
	Error string `json:"error"`
}
