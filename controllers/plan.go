package controllers

import (
	"fmt"

	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)

type ReadingPlans map[string]ReadingPlan

func (rps ReadingPlans) getContentsOf(plan string) ReadingPlanContent {
	chPlan, ok := rps[plan]
	if !ok {
		return nil
	}
	return chPlan.Days
}

type ReadingPlan struct {
	Name string             `json:"name"`
	Days ReadingPlanContent `json:"days"`
}

type ReadingPlanContent [][][]string

func (rpc ReadingPlanContent) getPlanLength() int {
	return len(rpc)
}

var plansColl = parseCollFromFile[ReadingPlans]("./plans.json")

func ChooseReadingPlan(c tele.Context) error {
	plan, err := getArg(0, c)
	if err != nil {
		return sendDocsForReadingPlan(c)
	}

	err = setPlan(plan, c)
	if err != nil {
		return sendDocsForReadingPlan(c)
	}

	return setReadingDay(0, c)
}

func getCurrPlanSchedule(c tele.Context) ReadingPlanContent {
	planCode, ok := c.Get(keys.PLAN).(string)
	if !ok {
		setDefaultPlan(c)
		return getCurrPlanSchedule(c)
	}

	return plansColl.getContentsOf(planCode)
}

var ErrUnknownReadingPlan = fmt.Errorf("Ошибка: неизвестный план чтения")

func setPlan(planCode string, c tele.Context) error {
	_, ok := plansColl[planCode]
	if !ok {
		setDefaultPlan(c)
		return ErrUnknownReadingPlan
	}

	c.Set(keys.PLAN, planCode)
	msg := fmt.Sprintf("Выбран план %v", planCode)
	return c.Send(msg)
}

const DEFAULT_READING_PLAN = "RCL"

func setDefaultPlan(c tele.Context) {
	c.Set(keys.PLAN, DEFAULT_READING_PLAN)
}

var sendDocsForReadingPlan = func() tele.HandlerFunc {
	msg := fmt.Sprintf(
		"%v *код_плана*\nДля выбора доступны следующие планы чтения:\n*код_плана — название_плана*",
		keys.API_PLAN_PATH,
	)

	for code, planObj := range plansColl {
		msg += fmt.Sprintf("\n%v — %v", code, planObj.Name)
	}

	return bindMessageSender(msg)
}()
