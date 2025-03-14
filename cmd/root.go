package cmd

import (
	"encoding/json"
	"fmt"
	"nlto/internal/openai"
	"nlto/internal/ui"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("exit")
		os.Exit(1)
	}

	args := os.Args[1:]
	query := strings.Join(args, " ")

	command, explanation, suggestJSON, dangerLevel, err := openai.GenerateCommand(query)
	if err != nil {
		fmt.Println(ui.Red + "Error: " + err.Error() + ui.Reset)
		os.Exit(1)
	}

	dangerLevelNum, err := strconv.Atoi(dangerLevel)
	if err != nil {
		dangerLevelNum = 0
	}

	// 画面のクリア
	fmt.Print("\033[H\033[2J")

	if dangerLevelNum > 7 {
		fmt.Println(ui.Red + "☠️   generated: " + command + ui.Reset)
	} else {
		fmt.Println(ui.Cyan + "🤖  generated: " + command + ui.Reset)
	}
	fmt.Println(ui.Blue + "📖  " + explanation + ui.Reset)

	var suggestList []map[string]string
	if suggestJSON != "" {
		if err := json.Unmarshal([]byte(suggestJSON), &suggestList); err != nil {
			fmt.Println(ui.Red + "Failed to parse suggestions: " + err.Error() + ui.Reset)
			os.Exit(1)
		}
	}

	if len(suggestList) > 0 {
		fmt.Println(ui.Cyan + "──────💡──────" + ui.Reset)
		for _, s := range suggestList {
			suggestCommand := s["suggest_command"]
			suggestExplanation := s["suggest_explanation"]
			fmt.Printf(ui.Yellow+"%s"+ui.Reset+" : %s\n", suggestCommand, suggestExplanation)
		}
	}

	fmt.Print(ui.Red + "🚀  Execute this command? (y/n): " + ui.Reset)
	var userInput string
	fmt.Scanln(&userInput)
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	userReinput := "y"
	if userInput == "y" && dangerLevelNum > 7 {
		fmt.Print(ui.Red + "☠️  This is DANGER command, execute? (y/n): " + ui.Reset)
		fmt.Scanln(&userReinput)
		userReinput = strings.ToLower(strings.TrimSpace(userReinput))
	}

	if userInput == "y" && userReinput == "y" {
		fmt.Println("Executing:", ui.Green+command+ui.Reset)

		cmdParts := strings.Fields(command)
		if len(cmdParts) == 0 {
			fmt.Println(ui.Red + "Error: Generated command is empty" + ui.Reset)
			return
		}

		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(ui.Red + "Execution failed: " + err.Error() + ui.Reset)
		}
	} else {
		fmt.Println(ui.Red + "Canceled" + ui.Reset)
	}
}
