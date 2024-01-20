package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lalathealter/telegospel/controllers"
	"github.com/lalathealter/telegospel/db"
	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return
	}

  if err := db.Get().Ping(); err != nil {
    log.Fatal(err)
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

	var (
		menu    = &tele.ReplyMarkup{ResizeKeyboard: true}
		btnHelp = menu.Text("ℹ Помощь")
		btnPrev = menu.Text("⬅️")
		btnCurr = menu.Text("📖")
		btnNext = menu.Text("➡️")
	)

	menu.Reply(
		menu.Row(btnPrev, btnCurr, btnNext),
		menu.Row(btnHelp),
	)

	b.Use(controllers.ManageUserContext)
	b.Handle(keys.API_PROVIDER_PATH, controllers.ChooseProvider)
	b.Handle(keys.API_TRANSLATION_PATH, controllers.ChooseTranslation)
  b.Handle(keys.API_PLAN_PATH, controllers.ChooseReadingPlan)
  b.Handle(keys.API_READING_DAY_PATH, controllers.ChooseReadingDay)
  b.Handle(keys.API_READING_DAY_MOVE_BACKWARD_PATH, controllers.MoveReadingDayBackwardBy)
  b.Handle(keys.API_READING_DAY_MOVE_FORWARD_PATH, controllers.MoveReadingDayForwardBy)

  b.Handle(&btnPrev, controllers.MovePrevReadingDay)
  b.Handle(&btnNext, controllers.MoveNextReadingDay)
	b.Handle(&btnCurr, controllers.GetTodayPassages)
	b.Handle(&btnHelp, controllers.GetHelp)
  b.Handle("/start", func(ctx tele.Context) error {
    return ctx.Send("Добро пожаловать в TeleGospel!", menu)
	})

	fmt.Println("the bot has been launched;")
	b.Start()
}
