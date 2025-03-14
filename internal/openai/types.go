package openai

type OpenAIRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

type OpenAICommandResponse struct {
	Command     string `json:"command"`
	Explanation string `json:"explanation"`
	DangerLevel string `json:"danger_level"`
	Suggest     []struct {
		SuggestCommand     string `json:"suggest_command"`
		SuggestExplanation string `json:"suggest_explanation"`
	} `json:"suggest"`
}
