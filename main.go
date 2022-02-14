package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/eskriett/confusables"
	"github.com/spf13/viper"
)

func isReady(session *discordgo.Session, ready *discordgo.Ready) {
	log.Println("Successfully connected to discord!")
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	wordsWrong := make([]string, 0)

	searchStringWithI := strings.ReplaceAll(message.Content, "I", "l")
	searchString := confusables.ToASCII(strings.ToLower(message.Content))

	for canadianWord, americanWord := range WordMap {
		wordIndex := strings.Index(searchString, canadianWord)
		wordWithIIndex := strings.Index(searchStringWithI, canadianWord)

		if wordIndex != -1 || wordWithIIndex != -1 {
			log.Printf("Replaced %s with %s", canadianWord, americanWord)

			wordsWrong = append(wordsWrong, "*"+americanWord)
			searchString = strings.Replace(searchString, canadianWord+" ", americanWord+" ", 1)
			searchString = strings.Replace(searchString, canadianWord+" ", americanWord+" ", 1)
		}
	}

	if len(wordsWrong) != 0 {
		finalMessage := strings.Join(wordsWrong, ", ")
		_, err := session.ChannelMessageSend(message.ChannelID, finalMessage)
		if err != nil {
			log.Println("Error editing message", message.Content, ". Error: ", err)
			return
		}
	} else {
		log.Println("Skipping: ", searchString)
	}
}

func main() {
	if os.Getenv("MODE") == "PRODUCTION" {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("Error getting config: ", err)
			panic(err)
		}

	}

	token := viper.GetString("TOKEN")
	log.Println(token)

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error while opening discord: ", err)
		panic(err)
	}

	session.AddHandler(messageHandler)
	session.AddHandler(isReady)

	err = session.Open()
	if err != nil {
		log.Println("Error connecting to discord: ", err)
		panic(err)
	}

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
