package controllers

import (
	"fmt"

	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

var GetHelp = func() func(tele.Context) error {
	msg := fmt.Sprintf(
    `Добро пожаловать в TeleGospel — приложение для проведения ежедневных библейских чтений
Вам доступны следующие команды:

%v *код_варианта* — выбор варианта перевода; для ознакомления со списком доступных переводов достаточно просто ввести эту команду`,
		keys.API_TRANSLATION_PATH,
	)
	return func(c tele.Context) error {
		return c.Send(msg, tele.ModeMarkdown)
	}
}()
