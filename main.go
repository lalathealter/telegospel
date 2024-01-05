package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	bus "github.com/lalathealter/telegospel/business"
	"github.com/lalathealter/telegospel/controllers"
	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv(keys.TG_TOKEN),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	prov := bus.ChooseProvider(bus.BibleGatewayLink)
	fmt.Println(prov.GetPassageLink("Matt 1:1", "NRSVUE"))

	var (
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		btnHelp = menu.Text("ℹ Помощь")
	)

	menu.Reply(
		menu.Row(btnHelp),
	)

	b.Handle(keys.API_TRANSLATION_PATH, controllers.ChooseTranslation)

	b.Handle(&btnHelp, controllers.GetHelp)

	fmt.Println("the bot has been launched;")
	b.Start()
}
