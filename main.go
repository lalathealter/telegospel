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

	b.Handle("/translation", ChooseTranslation)

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	var (
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		btnHelp     = menu.Text("ℹ Help")
		btnSettings = menu.Text("⚙ Settings")
	)

	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!", menu)
	})

	settingsMenu := &tele.ReplyMarkup{}
	localBtn := settingsMenu.Text("Language")
	translationLangBtn := settingsMenu.Text("Translation")
	translationVersionBtn := settingsMenu.Text("Version")

	settingsMenu.Reply(
		settingsMenu.Row(localBtn, translationLangBtn, translationVersionBtn),
	)

	localizationMenu := &tele.ReplyMarkup{}

	langBtnColl := controllers.ProduceLocalizationButtons(localizationMenu, b)
	localizationMenu.Reply(
		localizationMenu.Split(4, langBtnColl)...,
	)

	b.Handle(&translationLangBtn, func(ctx tele.Context) error {
		return ctx.Reply("translation", localizationMenu)
	})
	b.Handle(&btnSettings, func(c tele.Context) error {
		return c.Send("settings", settingsMenu)
	})

	b.Handle(&btnHelp, func(c tele.Context) error {
		return c.Send("Here is some help: ...")
	})

	fmt.Println("the bot has launched;")
	b.Start()
}

func ChooseTranslation(c tele.Context) error {
	c.Send("Choose a language of translation")
	c.Send("Choose a translation")
	m := c.Message()
	fmt.Println(m.Text)
	val := m.Text

	c.Set(keys.TRANSLATION, val)
	return nil
}
