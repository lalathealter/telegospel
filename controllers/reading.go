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

const DEFAULT_READING_PLAN = "RCL"
func (rps ReadingPlans) getDefaultContents() ReadingPlanContent {
  return rps.getContentsOf(DEFAULT_READING_PLAN)
}

type ReadingPlan struct {
	Name string             `json:"name"`
	Days ReadingPlanContent `json:"days"`
}

type ReadingPlanContent [][]string

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
	return nil
}

var ErrCouldNotGetPassages = fmt.Errorf("Ошибка: провалилась попытка получить отрывки для запрашиваемого дня")

func GetTodayPassages(c tele.Context) error {
	msg := "Day 1:"
	prov := GetCurrProvider(c)
	passes := getPassagesFor(0, c)
	if passes == nil {
		return ErrCouldNotGetPassages
	}

	for _, pass := range passes {
		link := prov.GetPassageLink(pass, GetTranslation(c))
		msg += fmt.Sprintf("\n[%v](%v)", pass, link)
	}

	return bindMessageSender(msg)(c)
}

func getPassagesFor(day int, c tele.Context) []string {
	planDays := getCurrPlanSchedule(c)
	if day >= len(planDays) {
		return nil
	}
	pass := planDays[day]
	return pass
}

func getCurrPlanSchedule(c tele.Context) ReadingPlanContent {
	v, ok := c.Get(keys.PLAN).(ReadingPlanContent)
	if !ok {
    setDefaultPlan(c)
		v = getCurrPlanSchedule(c)
	}
	return v
}


var ErrUnknownReadingPlan = fmt.Errorf("Ошибка: неизвестный план чтения")

func setPlan(planCode string, c tele.Context) error {
  _, ok := plansColl[planCode]
  if !ok {
    setDefaultPlan(c)
    return ErrUnknownReadingPlan
  }

  c.Set(keys.PLAN, plansColl.getContentsOf(planCode))
  return nil
}

func setDefaultPlan(c tele.Context) {
  c.Set(keys.PLAN, plansColl.getDefaultContents())
}


var sendDocsForReadingPlan = func() func(tele.Context) error {
	msg := fmt.Sprintf(
		"%v *код_плана*\nДля выбора доступны следующие планы чтения:\n*код_плана — название_плана*",
		keys.API_PLAN_PATH,
	)

	for code, planObj := range plansColl {
		msg += fmt.Sprintf("\n%v — %v", code, planObj.Name)
	}

	return bindMessageSender(msg)
}()
