package date

import "time"

type DateModifier func(t *time.Time)

func AddSecond(t *time.Time) {
	t = t.Add(time.Second)
}

func AddHour(t *time.Time) {
	t = t.Add(t.Hour())
}

func AddDay(t *time.Time) {
	t = t.AddDate(0, 0, 1)
}

func AddWeek(t *time.Time) {
	t = t.AddDate(0, 0, 7)
}

func AddMonth(t *time.Time) {
	t = t.AddDate(0, 1, 0)
}

func AddYear(t *time.Time) {
	t = t.AddDate(1, 0, 0)
}
