package controllers

import (
	tele "gopkg.in/telebot.v3"

	"encoding/json"
	"os"
)

func getArg(n int, c tele.Context) (string, error) {
	args := c.Args()
	if n >= len(args) {
		return "", tele.ErrEmptyText
	}
	return args[n], nil
}

func parseCollFromFile[T any](p string) T {
	coll := new(T)
	f, err := os.Open(p)
	if err != nil {
		panic(err)
	}

	if err := json.NewDecoder(f).Decode(coll); err != nil {
		panic(err)
	}

	return *coll
}

func bindMessageSender(msg string) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		return ctx.Send(msg, tele.ModeMarkdown)
	}
}
