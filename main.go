package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/eskriett/confusables"
	"github.com/spf13/viper"
	zerowidth "github.com/trubitsyn/go-zero-width"
)

func isReady(session *discordgo.Session, ready *discordgo.Ready) {
	log.Println("Successfully connected to discord!")
}

func clean(unclean string) string {
	return zerowidth.RemoveZeroWidthSpace(confusables.ToASCII(strings.ToLower(unclean)))
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	wordsWrong := make([]string, 0)

	searchString := clean(message.Content)

	for canadianWord, americanWord := range WordMap {
		wordIndex := strings.Index(searchString, canadianWord)
		if wordIndex != -1 {
			log.Printf("Replaced %s with %s", canadianWord, americanWord)

			wordsWrong = append(wordsWrong, "*"+americanWord)
			searchString = strings.Replace(searchString, canadianWord+" ", americanWord+" ", 1)
		}
	}

	if len(wordsWrong) != 0 {
		err := session.ChannelTyping(message.ChannelID)
		if err != nil {
			log.Println("Error 'typing': ", err)
			return
		}
		time.Sleep(1 * time.Second)

		finalMessage := strings.Join(wordsWrong, ", ")
		_, err = session.ChannelMessageSend(message.ChannelID, finalMessage)
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
