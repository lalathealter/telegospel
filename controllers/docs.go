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

%v *код_варианта* — выбор варианта перевода
%v *код_провайдера* — выбор сайта-провайдера текста
%v *код_плана* — выбор плана чтения (лекционария)

Для ознакомления со списком доступных опций достаточно просто ввести интересующую вас команду`,
		keys.API_TRANSLATION_PATH, keys.API_PROVIDER_PATH, keys.API_PLAN_PATH,
	)

	return bindMessageSender(msg)
}()
