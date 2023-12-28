package controllers

import (
	"errors"
	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

var russianResponsesColl = LocalizedResponseMap{
	UNKNOWN:         "Возникла ошибка",
	LANGUAGE_CHANGE: "Язык приложения был изменён на русский",
}
var englishResponsesColl = LocalizedResponseMap{
	UNKNOWN:         "An error has occured",
	LANGUAGE_CHANGE: "App language was changed to english",
}

var ErrUnknownResponseCode = errors.New("ERROR: cannot find an appropriate response in localization map")
var ErrUnknownLanguage = errors.New("ERROR: tried to access a non-existent localization")

var GetResponseMessage = func() func(tele.Context, ResponseCode) string {
	responseMap := LocalizationColl{
		RU_LANG: russianResponsesColl,
		EN_LANG: englishResponsesColl,
	}
	return func(c tele.Context, respCode ResponseCode) string {
		langVal := c.Get(keys.LOCALIZATION)
		lang := langVal.(LanguageCode)
		localResponses, ok := responseMap[lang]
		if !ok {
			return ErrUnknownLanguage.Error()
		}

		msg, ok := localResponses[ResponseCode(respCode)]
		if !ok {
			return ErrUnknownResponseCode.Error()
		}
		return msg
	}
}()

func ProduceLocalizationButtons(menu *tele.ReplyMarkup, b *tele.Bot) []tele.Btn {
	langsColl := LanguageCodes
	langBtnColl := make([]tele.Btn, len(langsColl))

	for i, v := range langsColl {
		option := menu.Text(string(v))
		langBtnColl[i] = option
		b.Handle(&option, bindChooseLocalization(v))
	}

	return langBtnColl
}

func bindChooseLocalization(v LanguageCode) func(tele.Context) error {
	return func(ctx tele.Context) error {
		ctx.Set(keys.LOCALIZATION, v)
		msg := GetResponseMessage(ctx, LANGUAGE_CHANGE)
		if msg == ErrUnknownLanguage.Error() {
			ctx.Set(keys.LOCALIZATION, DEFAULT_LANG)
			msg = GetResponseMessage(ctx, LANGUAGE_CHANGE)
		}
		return ctx.Send(msg)
	}
}

type LanguageCode string

const (
	RU_LANG      LanguageCode = "ru"
	EN_LANG      LanguageCode = "en"
	DEFAULT_LANG LanguageCode = RU_LANG
)

var LanguageCodes = [...]LanguageCode{RU_LANG, EN_LANG}

type ResponseCode string

const (
	UNKNOWN         ResponseCode = ""
	LANGUAGE_CHANGE ResponseCode = "lang_change"
)

type LocalizationColl map[LanguageCode]LocalizedResponseMap
type LocalizedResponseMap map[ResponseCode]string
