package controllers

import (
	"fmt"

	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

type TranslationsByLanguage map[string]TranslationVersions
type TranslationVersions map[string]string

var DEFAULT_TRANSLATION = "SYNOD"
var Translations = TranslationsByLanguage{
	"русский": {
		"SYNOD": "Синодальный",
	},
	"english": {
		"NKJV": "New King James' Version",
		"KJV":  "King James' Version",
	},
}

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

func sendDocsForTranslation(c tele.Context) error {
	return c.Send(getDocsForTranslation(), tele.ModeMarkdown)
}

var getDocsForTranslation = func()func() string {
	msg := fmt.Sprintf(
		"%v *код_варианта*\nДля выбора доступны следующие варианты\n{код_варианта — название}:",
		keys.API_TRANSLATION_PATH,
	)

	for lang, vers := range Translations {
		versList := ""
		for code, name := range vers {
			versList += fmt.Sprintf("\n\t%v — %v", code, name)
		}
		msg += fmt.Sprintf("\n*%v*:%v", lang, versList)
	}
	return func() string {
		return msg
	}
}()

func getArg(n int, c tele.Context) (string, error) {
	args := c.Args()
	if n >= len(args) {
		return "", tele.ErrEmptyText
	}
	return args[n], nil
}
