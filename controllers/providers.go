package controllers

import (
	"fmt"

	bus "github.com/lalathealter/telegospel/business"
	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

func ChooseProvider(c tele.Context) error {
	provName, err := getArg(0, c)
	if err != nil {
		return sendDocsForProvider(c)
	}
  return setProvider(provName, c)
}

func GetCurrProvider(c tele.Context) bus.BSI {
  prov, ok := c.Get(keys.PROVIDER).(bus.BSI)
  if !ok || prov == nil {
    setProvider("", c)
    prov = GetCurrProvider(c)
  }
  return prov
}

func setProvider(provName string, c tele.Context) error {
	prov := bus.GetProviderInterface(provName)
	c.Set(keys.PROVIDER, prov)
  msg := fmt.Sprintf("Выбран провайдер %v", prov)
  return c.Send(msg)
}

var sendDocsForProvider = func() tele.HandlerFunc {
	msg := fmt.Sprintf(
		"%v *код_провайдера*\nДля выбора доступны следующие провайдеры:\n*код_провайдера — название_провайдера*",
		keys.API_PROVIDER_PATH,
	)
	for code, link := range bus.BibleSourcesColl {
		msg += fmt.Sprintf("\n*%v* — %v", code, link)
	}
	return bindMessageSender(msg)
}()
