package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/new", func(c tele.Context) error {
		admin, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
		if err != nil {
			log.Fatal(err)
			return err
		}

		if c.Sender() == nil || c.Sender().ID != admin {
			return c.Send("Sorry, you can't create new facts.")
		}
		//msg is c.Data()
		//TODO: manage creation errors
		factID, err := createFact(c.Data())
		if err != nil {
			c.Send("Whoops, something went wrong, try again later: " + err.Error())
		}
		return c.Send("Fact " + factID + " created!")
	})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		text := `I can help you create and manage new facts.
You can control facts by sending these commands:

 /new <fact> - create a new fact`

		return c.Send(text)
	})

	bot.Start()
}

type FactResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func createFact(msg string) (string, error) {
	url := os.Getenv("FACTS_URL")
	msg = strings.ReplaceAll(msg, `"`, `\"`)

	var jsonStr = []byte(`{"content":"` + msg + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FACTS_SECRET"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var fact FactResponse
	if err := json.Unmarshal(body, &fact); err != nil {
		fmt.Println("Can't unmarshal JSON")
		return "", err
	}
	return fact.ID, nil
}
