package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/telebot.v3"
)

func main() {

	token := os.Getenv("K8S_KUBECTL_BOT_TOKEN")
	if token == "" {
		log.Fatal("K8S_KUBECTL_BOT_TOKEN is not set")
	}

	b, err := telebot.NewBot(telebot.Settings{
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	b.Handle(telebot.OnText, func(c telebot.Context) error {

		// If you want you can check if the user is authorized
		// by adding a check using c.Sender().Username here if necessary

		userCommand := c.Text()
		log.Printf("Received command: \"%s\" from user %s", userCommand, c.Sender().Username)

		// Determine the command prefix
		if strings.HasPrefix(userCommand, "kubectl") {
			// userCommand = "kubectl get pods -n monitoring"
			// strings.TrimPrefix(userCommand, "kubectl") removes kubectl --> " get pods -n monitoring"
			// strings.TrimSpace(..) removes any leading spaces --> "get pods -n monitoring"
			userCommand = strings.TrimSpace(strings.TrimPrefix(userCommand, "kubectl"))
		} else if strings.HasPrefix(userCommand, "k") {
			// same thing as before but for command that start with k
			userCommand = strings.TrimSpace(strings.TrimPrefix(userCommand, "k"))
		} else {
			return c.Send("Invalid command. Please start your command with `kubectl` or `k`.")
		}

		// The command is split for exec.Command
		commandParts := strings.Fields(userCommand)
		if len(commandParts) == 0 {
			return c.Send("No command provided after prefix.")
		}

		// Even if the user used k the command will start with kubectl
		cmd := exec.Command("kubectl", commandParts...)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return c.Send(fmt.Sprintf("Failed to execute kubectl: %s\n%s", err.Error(), string(output)))
		}

		response := string(output)
		if len(response) > 4000 { // Telegram message length limit is 4096
			response = response[:4000] + "\n...output truncated..."
		}

		// Format the response string with MarkdownV2 syntax, wrapping the output in triple backticks to preserve
		// the formatting when displayed in Telegram. In this way it is better if the output is very large
		// like "kubectl get po -A -o wide"
		response = fmt.Sprintf("```\n%s\n```", response)

		return c.Send(response, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	})

	log.Println("k8s-kubectl-bot is running...")
	b.Start()
}
