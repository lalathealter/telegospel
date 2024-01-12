package controllers

import (
	"fmt"
	"strconv"

	"github.com/lalathealter/telegospel/keys"
	tele "gopkg.in/telebot.v3"
)


func getReadingDay(c tele.Context) int {
  dayData := c.Get(keys.READING_DAY)
  
  var day int
  switch dayData.(type) {
  case float64:
    day = int(dayData.(float64))
  case string:
    dayc, err := strconv.Atoi(dayData.(string))
    day = dayc
    if err != nil {
      setReadingDay(0, c)
      return 0
    }
  default:
    dayc, ok := dayData.(int)
    if !ok {
      setReadingDay(0, c)
      return 0
    }
    day = dayc
  }

  ans := clampDay(day, c)
  return ans
}

func ChooseReadingDay(c tele.Context) error {
  arg, err := getArg(0, c)
  if err != nil {
    return sendDocsForReadingDay(c)
  }

  i, err := strconv.Atoi(arg)
  if err != nil {
    return sendDocsForReadingDay(c)
  }

  return setReadingDay(i-1, c)
}

func setReadingDay(dayIndex int, c tele.Context) error {
  dayIndex = clampDay(dayIndex, c)

  c.Set(keys.READING_DAY, dayIndex)
  msg := fmt.Sprintf("Выбран день %v", dayIndex+1)

  return c.Send(msg)
}

func clampDay(v int, c tele.Context) int {
  planLen := getCurrPlanSchedule(c).getPlanLength()
  if v < 0 {
    v = -v
  }

  if v >= planLen {
    v = planLen - 1
  }

  return v
}

var sendDocsForReadingDay = func()tele.HandlerFunc{
  msg :=  fmt.Sprintf(
		"%v *день*\nГде *день* - целое число больше 0",
		keys.API_READING_DAY_PATH,
	)

  return bindMessageSender(msg)
}()


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


var ErrCouldNotGetPassages = fmt.Errorf("Ошибка: провалилась попытка получить отрывки для запрашиваемого дня")

func GetTodayPassages(c tele.Context) error {
  day := getReadingDay(c)
  prov := GetCurrProvider(c)
	passes := getPassagesFor(day, c)
	if passes == nil {
		return ErrCouldNotGetPassages
	}

  msg := fmt.Sprintf("День %v", day+1)
  msg += prov.GetPassageLink(passes, GetTranslation(c))

	return bindMessageSender(msg)(c)
}

func getPassagesFor(dayIndex int, c tele.Context) []string {
	planDays := getCurrPlanSchedule(c)
	if dayIndex >= planDays.getPlanLength() {
		return nil
	}
	pass := planDays[dayIndex]
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
