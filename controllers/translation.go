package controllers

import (
	"fmt"

	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

type TranslationsByLanguage map[string]map[string]string

var DEFAULT_TRANSLATION = "SYNOD"
var Translations = parseCollFromFile[TranslationsByLanguage]("./translations.json")

var ErrNoTranslationEntries = fmt.Errorf("Ошибка: нет информации о доступных вариантах перевода")


var ErrUnknownTranslation = fmt.Errorf("Ошибка: неизвестный вариант перевода")

func GetTranslation(c tele.Context) string {
	code, ok := c.Get(keys.TRANSLATION).(string)
	if !ok {
		c.Set(keys.TRANSLATION, DEFAULT_TRANSLATION)
		code = DEFAULT_TRANSLATION
	}
	return code
}

func setTranslationVersion(version string, c tele.Context) error {
	for _, versions := range Translations {
		for code := range versions {
			if code == version {
				c.Set(keys.TRANSLATION, version)
				return nil
			}
		}
	}

	return ErrUnknownTranslation
}

func ChooseTranslation(c tele.Context) error {
	version, err := getArg(0, c)
	if err != nil {
		return sendDocsForTranslation(c)
	}

	err = setTranslationVersion(version, c)
	if err != nil {
		return sendDocsForTranslation(c)
	}

	return nil
}

var sendDocsForTranslation = func() func(tele.Context) error {
	msg := fmt.Sprintf(
		"%v *код_варианта*\nДля выбора доступны следующие варианты\n*код_варианта — название*:",
		keys.API_TRANSLATION_PATH,
	)

	for lang, vers := range Translations {
		versList := ""
		for code, name := range vers {
			versList += fmt.Sprintf("\n\t%v — %v", code, name)
		}
		msg += fmt.Sprintf("\n*%v*:%v", lang, versList)
	}

  return bindMessageSender(msg)
}()

