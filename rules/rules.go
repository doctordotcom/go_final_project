package rules

import (
	"errors"
	"m/global"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate - вычисляет следующую дату для задачи согласно заданным правилам повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {

	// Проверка пустой строки в repeat
	if repeat == "" {
		return "", errors.New("The repetition rule is not specified")
	}

	// Проверка и парсинг исходной даты
	startDate, err := time.Parse(global.DateFormat, date)
	if err != nil {
		return "", errors.New("Incorrect date format: " + date)
	}

	// Переменная для хранения следующей даты.
	var nextDate time.Time

	// Определение следующей даты на основе правила повторения
	// Разделение правила повторения на части
	repeatSlice := strings.Split(repeat, " ")
	// Проверка корректности формата правила повторения по длине массива
	if len(repeatSlice) < 1 || len(repeatSlice) > 3 {
		return "", errors.New("Incorrect format of the repetition rule: " + repeat)
	}
	// Выбор обработки в зависимости от первого символа правила
	switch repeatSlice[0] {
	// В случае d - дни
	case "d":
		// Обработка правила повторения на определенные дни
		if len(repeatSlice) != 1 {
			// Чтение числа дней из правила
			days, err := strconv.Atoi(repeatSlice[1])
			if err != nil {
				return "", errors.New("Incorrect format")
			}
			// Проверка корректности числа дней
			if days < 1 || days > 400 {
				return "", errors.New("The number of days should be from 1 to 400")
			}
			nextDate = startDate.AddDate(0, 0, days)
			// Передвижение даты до тех пор, пока она не станет позже текущей
			for nextDate.Before(now) {
				nextDate = nextDate.AddDate(0, 0, days)
			}
		}
	// В случае y - раз в год
	case "y":
		nextDate = startDate.AddDate(1, 0, 0)
		// Передвижение даты до тех пор, пока она не станет позже текущей
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

	// В случае w - в определённые дни недели
	case "w":
		if len(repeatSlice) == 1 {
			return "", errors.New("Error: it is necessary to specify the days of the week in the format 'w <days>', for example, 'w 2'")
		}
		// Массив для хранения дней недели, в которые нужно повторить событие
		weekDays := make([]bool, 7) // понедельник-воскресенье
		days := strings.Split(repeatSlice[1], ",")

		// Парсинг дней недели
		for _, day := range days {
			dayNumber, err := strconv.Atoi(day)
			if err != nil || dayNumber < 1 || dayNumber > 7 {
				return "", errors.New("The number of days must be from 1 to 7 inclusive, here they are:" + day)
			}
			// Установка true для соответствующего дня
			weekDays[dayNumber-1] = true
		}

		// Поиск ближайшей даты в указанные дни недели
		nextDate = startDate
		for {
			// 0=воскресенье, 6=суббота
			dayOfWeek := int(nextDate.Weekday())
			if dayOfWeek == 0 {
				// Преобразование воскресенья к 7
				dayOfWeek = 7
			}
			if weekDays[dayOfWeek-1] && nextDate.After(now) {
				break
			}
			// Переход к следующему дню
			nextDate = nextDate.AddDate(0, 0, 1)
		}

	// В случае m - в определённые дни месяца
	case "m":
		// Разделение строки с правилами повторения на части по пробелам
		parts := strings.Split(strings.TrimPrefix(repeat, "m "), " ")
		// Парсинг дней месяца из первой части
		daysOfMonth := parseDays(parts[0])
		// Проверка на корректность парсинга дней
		if daysOfMonth == nil {
			return "", errors.New("The days of the month are specified incorrectly")
		}
		// Инициализация массива месяцев. По умолчанию все месяцы
		var months []int
		// Если указаны месяцы, то они парсятся из второй части
		if len(parts) > 1 {
			months = parseMonths(parts[1])
			// Проверка на корректность парсинга месяцев
			if months == nil {
				return "", errors.New("The months are specified incorrectly")
			}
		} else {
			// Если месяцы не указаны, то добавляются все месяцы от 1 до 12
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}
		// Сортировка дней и месяцев для удобства работы
		sort.Ints(daysOfMonth)
		sort.Ints(months)
		// Цикл для поиска ближайшей даты, соответствующей правилам
		for {
			// Текущий месяц и год
			month := int(startDate.Month())
			year := startDate.Year()
			// Список дат, подходящих под критерии
			var sortDates []time.Time
			// Цикл по каждому месяцу из списка
			for _, m := range months {
				// Последний день текущего месяца
				lastDayOfTheMonth := getLastDayOfMonth(year, time.Month(m))
				// Если текущий месяц больше или равен рассматриваемому месяцу
				if m >= month {
					// Цикл по каждому дню из списка дней месяца
					for i, d := range daysOfMonth {
						// Если день не превышает последний день текущего месяца
						if d <= lastDayOfTheMonth {
							// Расчёт новой даты
							newDate := calculateNewDate(year, m, d)
							// Если новая дата больше текущей даты
							if newDate.After(now) {
								// Добавление даты в список подходящих дат
								sortDates = append(sortDates, newDate)
								// Если это последняя дата в списке дней - возвращается самая ранняя
								if i == len(daysOfMonth)-1 {
									earliestDate := getEarliestDate(sortDates)
									return earliestDate, nil
								}
							}
						}
					}
				}
			}
			// Переход к следующему месяцу
			startDate = startDate.AddDate(0, 1, 0)
		}

	default:
		return "", errors.New("Unsupported repetition rule format: " + repeat)
	}

	// Проверка, больше ли следующая дата указанного времени
	if nextDate.Before(now) {
		return "", errors.New("The next date must be later than the current date.")
	}
	// Возвращается вычисленная следующая дата
	return nextDate.Format(global.DateFormat), nil
}

// parseDays - вспомогательная функция для парсинга дней месяца
// Принимает строку дней (через запятую), например "1,15,31"
// Возвращает срез целых чисел, представляющих дни
// Возвращает nil, если в строке есть некорректные данные
func parseDays(days string) []int {
	var result []int
	for _, day := range strings.Split(days, ",") {
		d, err := strconv.Atoi(day)
		if err != nil || d < -2 || d == 0 || d > 31 {
			return nil
		}
		result = append(result, d)
	}
	return result
}

// parseMonths - вспомогательная функция для парсинга месяцев
// Принимает строку месяцев (через запятую), например "1,3,12"
// Возвращает срез целых чисел, представляющих номера месяцев (1-12)
// Возвращает nil, если в строке есть некорректные данные
func parseMonths(months string) []int {
	var result []int
	for _, month := range strings.Split(months, ",") {
		m, err := strconv.Atoi(month)
		if err != nil || m < 1 || m > 12 {
			return nil
		}
		result = append(result, m)
	}
	return result
}

// getLastDayOfMonth - вспомогательная функция получения последнего дня месяца
// Принимает год и месяц
// Возвращает целое число - последний день месяца
func getLastDayOfMonth(year int, month time.Month) int {
	// Создание временной метки на первый день следующего месяца
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	// Вычисление последнего дня текущего месяца путём вычитания 24 часов из первого дня следующего месяца
	lastDayOfMonth := firstDayNextMonth.Add(-24 * time.Hour)

	return lastDayOfMonth.Day()
}

// calculateNewDate - вспомогательная функция расчёта новой даты
// Принимает год, месяц и день
// Возвращает time.Time значение новой даты
// Обрабатывает особые случаи, когда день больше максимального для месяца и когда день равен -1, -2
// Возвращает time.Time{} если входные данные некорректны
func calculateNewDate(year, month, day int) time.Time {
	lastDay := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	if day > 0 {
		if day > lastDay {
			day = lastDay
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	} else if day == -1 {
		return time.Date(year, time.Month(month), lastDay, 0, 0, 0, 0, time.UTC)
	} else if day == -2 {
		return time.Date(year, time.Month(month), lastDay-1, 0, 0, 0, 0, time.UTC)
	}
	return time.Time{}
}

// getEarliestDate - функция получения ближайшей даты
// Принимает срез дат
// Возвращает строку с самой ранней датой в формате YYYYMMDD
// Возвращает пустую строку, если входной массив пустой
func getEarliestDate(dates []time.Time) string {
	if len(dates) == 0 {
		return ""
	}
	// Предположение, что первая дата самая ранняя
	earliest := dates[0]

	// Итерация по остальным датам
	for _, date := range dates {
		if date.Before(earliest) { // Сравнение с текущей минимальной
			earliest = date
		}
	}

	// Возвращение самой ранней даты в формате YYYYMMDD
	return earliest.Format(global.DateFormat)
}
