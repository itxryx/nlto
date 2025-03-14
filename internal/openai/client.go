package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const OpenAiApiEndpoint = "https://api.openai.com/v1/chat/completions"

var OpenAiApiKey = os.Getenv("OPENAI_API_KEY")

func GenerateCommand(query string) (string, string, string, string, error) {
	if OpenAiApiKey == "" {
		return "", "", "", "", errors.New("OPENAI_API_KEY is not set, please run `export OPENAI_API_KEY=YOUR_API_KEY`")
	}

	contentMessage := `
	あなたはUNIX・Linux・Macのコマンドについて世界最高の技術と知識を持つエンジニアです。
	必ず以下のJSON形式でレスポンスを返してください。

	{
		"command": "<command>",
		"explanation": "<explanation>"
		"danger_level": "<danger_level>"
		"suggest": [
			"<suggest_command>": "<suggest_explanation>"
		]
	}

	レスポンス形式のルール
	- 出力言語: リクエストの言語を推測し、それに従ってください。基本的には日本語で返してください。
	- <command>（必須）: 無駄のない最適なコマンドを記述してください。
	- <explanation>（必須）: <command>の説明を100文字以内の正確で簡潔な文章で記述してください。オプションに関する説明も含めてください。
	- <danger_level>（必須）: <command>の実行時の影響を1（低）~10（高）の10段階で評価して記述してください。
	- suggest（必須）: 1つ以上4つ未満の、<suggest_command>をキー、<suggest_explanation>を値に持つ配列のセットを返してください。
	- <suggest_command>（必須）: <command>に近い動作を行うコマンドを記述してください。
	- <suggest_explanation>（必須）: <suggest_command>の説明を35文字以内の正確で簡潔な文章で記述してください。
	`

	messages := []ChatMessage{
		{Role: "system", Content: contentMessage},
		{Role: "user", Content: query},
	}

	requestData := OpenAIRequest{
		Model:       "gpt-4o",
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.2,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", "", "", "", err
	}

	req, err := http.NewRequest("POST", OpenAiApiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OpenAiApiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", "", "", err
	}

	var openAiRes OpenAIResponse
	if err := json.Unmarshal(body, &openAiRes); err != nil {
		return "", "", "", "", err
	}

	if len(openAiRes.Choices) == 0 {
		return "", "", "", "", errors.New("no response from OpenAI API")
	}

	resText := openAiRes.Choices[0].Message.Content
	resText = strings.TrimPrefix(resText, "```json")
	resText = strings.TrimSuffix(resText, "```")
	resText = strings.TrimSpace(resText)
	var commandRes OpenAICommandResponse
	if err := json.Unmarshal([]byte(resText), &commandRes); err != nil {
		return "", "", "", "", err
	}

	suggestList := make([]map[string]string, len(commandRes.Suggest))
	for i, s := range commandRes.Suggest {
		suggestList[i] = map[string]string{
			"suggest_command":     s.SuggestCommand,
			"suggest_explanation": s.SuggestExplanation,
		}
	}

	suggestJSON, err := json.Marshal(suggestList)
	if err != nil {
		return "", "", "", "", err
	}

	// 空の場合は空文字を返す（が、ありえないはず）
	suggestStr := string(suggestJSON)
	if len(suggestList) == 0 {
		suggestStr = ""
	}

	return commandRes.Command, commandRes.Explanation, suggestStr, commandRes.DangerLevel, nil
}
