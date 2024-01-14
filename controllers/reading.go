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

func bindMoveReadingDay(mov int) tele.HandlerFunc {
	return func(c tele.Context) error {
		d := getReadingDay(c)
		sum := d + mov
		res := clampDay(sum, c)
		if mov > 0 {
			if sum != res {
				sendDocsExceededPlanRange(c)
			}
		}

		return setReadingDay(res, c)
	}
}

func bindMoveReadingDayBy(signMult int) tele.HandlerFunc {
	return func(c tele.Context) error {
		arg, err := getArg(0, c)
		if err != nil {
			return sendDocsForMovingDayBy(c)
		}

		mov, err := strconv.Atoi(arg)
		if err != nil {
			return sendDocsForMovingDayBy(c)
		}

		if mov < 0 {
			mov = -mov
		}

		return bindMoveReadingDay(signMult * mov)(c)
	}
}

var MoveNextReadingDay = bindMoveReadingDay(1)
var MovePrevReadingDay = bindMoveReadingDay(-1)

var MoveReadingDayForwardBy = bindMoveReadingDayBy(1)
var MoveReadingDayBackwardBy = bindMoveReadingDayBy(-1)

var sendDocsExceededPlanRange = func() tele.HandlerFunc {
	msg := "Был достигнут последний день плана"
	return bindMessageSender(msg)
}()

var sendDocsForMovingDayBy = func() tele.HandlerFunc {
	msg := fmt.Sprintf(
		"%v *количество_дней*\nГде *количество_дней* — количество дней (больше 0), на которое вы хотите переместиться по плану",
		keys.API_READING_DAY_PATH,
	)

	return bindMessageSender(msg)
}()

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

func getPlanLengthFrom(c tele.Context) int {
	return getCurrPlanSchedule(c).getPlanLength()
}

func clampDay(v int, c tele.Context) int {
	planLen := getPlanLengthFrom(c)
	if v < 0 {
		v = 0
	}

	if v >= planLen {
		v = planLen - 1
	}

	return v
}

var sendDocsForReadingDay = func() tele.HandlerFunc {
	msg := fmt.Sprintf(
		"%v *день*\nГде *день* - целое число больше 0",
		keys.API_READING_DAY_PATH,
	)

	return bindMessageSender(msg)
}()

var ErrCouldNotGetPassages = fmt.Errorf("Ошибка: провалилась попытка получить отрывки для запрашиваемого дня")

func GetTodayPassages(c tele.Context) error {
	day := getReadingDay(c)
	prov := GetCurrProvider(c)
	sections := getPassagesFor(day, c)
	if sections == nil {
		return ErrCouldNotGetPassages
	}

	tran := GetTranslation(c)
	msg := fmt.Sprintf("День %v", day+1)
	msg += prov.GetPassageLink(sections[0], tran)
	if len(sections) > 1 {
		msg += fmt.Sprintf("\nДополнение:")
		msg += prov.GetPassageLink(sections[1], tran)
	}

	return bindMessageSender(msg)(c)
}

func getPassagesFor(dayIndex int, c tele.Context) [][]string {
	planDays := getCurrPlanSchedule(c)
	if dayIndex >= planDays.getPlanLength() {
		return nil
	}
	pass := planDays[dayIndex]
	return pass
}
