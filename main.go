package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/exec"
	"strings"
)

var commands = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/jobs"),
		tgbotapi.NewKeyboardButton("/run"),
	),
)

var jobs = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/run one-click"),
	),
)

func main() {
	token := "TOKEN"
	log.SetPrefix("[X-RELEASE] ")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Error creating bot: ", err)
	}

	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		if update.Message.IsCommand() { // ignore any non-command Messages
			switch update.Message.Command() {
			case "start":
				msg.ReplyMarkup = commands
			case "jobs":
				msg.ReplyMarkup = jobs
			case "run":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				runJob(msg.Text)
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func runJob(run string) {
	fmt.Printf("Starting job: %s \n", run)

	currentRepository := "git@github.com:sab1tm/one-click.git"
	currentPath := "/Users/sab1tm/xrunner/"
	repositoryName := getRepositoryName(currentRepository)

	fmt.Println("Current repository:", currentRepository)
	fmt.Println("Current path:", currentPath)

	// clone repository
	goToPath(currentPath)
	execCmd("git", "clone", currentRepository)
	fmt.Println("Cloning successfully")

	// build maven
	goToPath(currentPath + "/" + repositoryName)
	execCmd("mvn", "clean", "package")

	// clean up
	goToPath(currentPath)
	execCmd("rm", "-rf", repositoryName)
	fmt.Println("Remove files")

	fmt.Println("Job finished")
}

func goToPath(path string) bool {
	if err := os.Chdir(path); err != nil {
		log.Fatalf("Error %s: %v", path, err)
		return false
	}
	return false
}

func execCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func getRepositoryName(url string) string {
	parts := strings.Split(url, "/")
	repoWithExt := parts[len(parts)-1]
	return strings.TrimSuffix(repoWithExt, ".git")
}
