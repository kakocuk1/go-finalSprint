package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102" // формат даты для хранения в БД и передачи в API

// afterNow проверяет, что дата date находится после даты now
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid date dstart: " + err.Error())
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {

	case "d":
		if len(parts) != 2 {
			return "", errors.New("dayily interval is not specified")
		}
		days, err := strconv.Atoi(parts[1])

		if err != nil {
			return "", errors.New("interval is not a number")
		}
		if days <= 0 || days > 400 {
			return "", errors.New("interval is out of range (1-400)")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("weekly interval is not specified")
		}
		weekdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		for {
			date = date.AddDate(0, 0, 1)
			// time.Weekday() возвращает день недели в виде числа от 0 (воскресенье) до 6 (суббота),
			// в repeat указывается день недели в виде числа от 1 (понедельник) до 7 (воскресенье)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			// Останавливаемся, как только одновременно выполнены два условия:
			// дата уже больше now, и день недели этой даты входит в разрешённый набор weekdays.
			if afterNow(date, now) && weekdays[wd] {
				break
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("monthly interval is not specified")
		}
		days, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}
		var months map[int]bool //Парсим список дней месяца
		if len(parts) >= 3 {
			months, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		}
		for {
			date = date.AddDate(0, 1, 0)
			if !afterNow(date, now) {
				continue
			}
			if months != nil && !months[int(date.Month())] {
				continue
			}
			lastDay := lastDayOfMonth(date)
			day := date.Day()
			if days[day] || (day == lastDay && days[-1]) || (day == lastDay-1 && days[-2]) {
				break
			}
		}

	default:
		return "", errors.New("invalid repeat format")
	}

	return date.Format(dateFormat), nil
}

// разбор списка дней недели
func parseWeekdays(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, p := range strings.Split(s, ",") { // разбиваем строку по запятым
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 7 {
			return nil, errors.New("incorrect weekday: " + p)
		}
		result[n] = true // сохраняем в результат, что день недели n разрешён для повторения
	}
	return result, nil
}

// разбор списка дней месяца
func parseMonthDays(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(p)
		if err != nil || n == 0 || n < -2 || n > 31 { // разрешаем указывать дни от 1 до 31, а также -1 (последний день месяца) и -2 (предпоследний день месяца)
			return nil, errors.New("incorrect month day: " + p)
		}
		result[n] = true
	}
	return result, nil
}

// разбор списка месяцев
func parseMonths(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 12 {
			return nil, errors.New("incorrect month: " + p)
		}
		result[n] = true
	}
	return result, nil
}

// вычисление последнего дня месяца
// Чтобы обработать m-1 (последний день месяца) и m-2 (предпоследний)
func lastDayOfMonth(date time.Time) int {
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastDay := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

// HTTP-обработчик для получения следующей даты по заданным параметрам
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	now := time.Now()
	if nowParam != "" {
		parsedNow, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "incorrect now parameter", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	next, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
