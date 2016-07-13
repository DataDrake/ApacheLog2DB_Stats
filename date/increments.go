package date

import "time"

type DateModifier func(t time.Time) time.Time

func AddSecond(t time.Time) time.Time {
	return t.Add(time.Second)
}

func AddHour(t time.Time) time.Time {
	return t.Add(time.Hour)
}

func AddDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func AddWeek(t time.Time) time.Time {
	return t.AddDate(0, 0, 7)
}

func AddMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, 0)
}

func AddYear(t time.Time) time.Time {
	return t.AddDate(1, 0, 0)
}
