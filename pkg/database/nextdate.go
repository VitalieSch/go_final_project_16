package database

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Вычисляем следующую дату выполнения задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {

	//Анализируем исходную дату
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("ошибка анализа исходной даты")
	}
	//Обработка повторений
	switch {
	case strings.HasPrefix(repeat, "d"):
		part := strings.Split(repeat, " ")
		if len(part) != 2 {
			return "", errors.New("неверный формат")
		}

		days, err := strconv.Atoi(part[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неверный формат интервала")
		}

		nextDate := startDate

		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.After(now) {
				break
			}
		}
		return nextDate.Format("20060102"), nil

	case repeat == "y":
		startDate = startDate.AddDate(1, 0, 0)
		for !startDate.After(now) {
			startDate = startDate.AddDate(1, 0, 0)
		}
	default:
		return "", errors.New("неверный формат повторений")
	}
	return startDate.Format("20060102"), nil
}

// Обработчик запросов для "/api/nextdate"
func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	now, err := parseTime(nowStr)
	if err != nil {
		http.Error(w, `{"error": "нневерный формат даты"}`, http.StatusBadRequest)
		return
	}

	NextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, `{"error": "нневерный формат даты"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, NextDate)
}

func parseTime(nowStr string) (time.Time, error) {

	timeParse, err := time.Parse("20060102", nowStr)

	return timeParse, err
}
