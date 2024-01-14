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
%v *номер_дня* — выбор дня внутри плана чтения
%v *количество_дней* — продвинуться в плане чтения на заданное количество дней вперёд
%v *количество_дней* — вернуться в плане чтения на заданное количество дней назад

Для ознакомления со списком доступных опций достаточно просто ввести интересующую вас команду`,
		keys.API_TRANSLATION_PATH, keys.API_PROVIDER_PATH,
		keys.API_PLAN_PATH, keys.API_READING_DAY_PATH,
		keys.API_READING_DAY_MOVE_FORWARD_PATH,
		keys.API_READING_DAY_MOVE_BACKWARD_PATH,
	)

	return bindMessageSender(msg)
}()
