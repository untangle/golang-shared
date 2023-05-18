package booleval

import (
	"fmt"
	"strings"
	"time"
)

func tryToParseTimeString(timeOfDay string) (time.Time, error) {
	formats := []string{
		time.Kitchen,
		"3:04:00PM",
		time.UnixDate,
		time.Stamp,
		time.RFC822,
		time.RFC822Z,
	}
	timeOfDay = strings.ToUpper(timeOfDay)
	for _, format := range formats {
		if possibleTime, err := time.Parse(format, timeOfDay); err == nil {
			return possibleTime, nil
		}

	}
	return time.Unix(0, 0), fmt.Errorf("booleval TimeOfDayComparable: unable to parse time: %s", timeOfDay)
}

// TimeComparable is a comparable for time objects.
type TimeComparable struct {
	time time.Time
}

func (t TimeComparable) timeCompare(comp func(lhs time.Time, rhs time.Time) bool, other any) (bool, error) {
	switch val := other.(type) {
	case int, uint32, uint64, int64, uint, int32:
		int64Val, _ := getInt(val)
		return comp(t.time, time.Unix(int64Val, 0)), nil
	case string:
		if parsedTime, err := tryToParseTimeString(val); err == nil {
			return comp(t.time, parsedTime), nil
		} else {
			return false, fmt.Errorf(
				"booleval TimeComparable.timeCompare: unable to parse date string: %s",
				val)
		}
	case time.Time:
		return comp(t.time, val), nil
	}
	return false, fmt.Errorf("booleval TimeComparable.timeCompare: can't convert %v(%T) to a time",
		other, other)
}

// Equal returns true if other is an integer and, if interpreted as a
// unix timestamp, it is equal to t.time, or if other is a time object
// and is equal to t.time
func (t TimeComparable) Equal(other any) (bool, error) {
	return t.timeCompare(
		func(lhs time.Time, rhs time.Time) bool {
			return lhs.Equal(rhs)
		},
		other,
	)
}

// Greater is true if other is an integer and t.time is greater than its value interpreted as
// a unix timestamp, or if other is a time.Time
// and t.time is greater than it.
func (t TimeComparable) Greater(other any) (bool, error) {
	return t.timeCompare(
		func(lhs time.Time, rhs time.Time) bool {
			return lhs.After(rhs)
		},
		other,
	)
}

// TimeComparable is a comparable for time objects.
type DayOfWeekComparable struct {
	dayOfWeek time.Weekday
}

var dayMap = map[string]time.Weekday{
	"SUNDAY":    time.Sunday,
	"MONDAY":    time.Monday,
	"TUESDAY":   time.Tuesday,
	"WEDNESDAY": time.Wednesday,
	"THURSDAY":  time.Thursday,
	"FRIDAY":    time.Friday,
	"SATURDAY":  time.Saturday,
}

func NewDayOfWeekFromString(val string) (DayOfWeekComparable, error) {
	if day, dayFound := dayMap[strings.ToUpper(val)]; dayFound {
		return DayOfWeekComparable{dayOfWeek: day}, nil
	}
	return DayOfWeekComparable{}, fmt.Errorf("booleval.NewDayOfWeekFromString: %s is a bad day of week",
		val)
}
func (t DayOfWeekComparable) dowCompare(comp func(lhs time.Weekday, rhs time.Weekday) bool, other any) (bool, error) {
	switch val := other.(type) {
	case int, uint32, uint64, int64, uint, int32:
		intval, _ := getInt(val)
		return comp(t.dayOfWeek, time.Unix(intval, 0).Weekday()), nil
	case time.Weekday:
		return comp(t.dayOfWeek, val), nil
	case string:
		if day, ok := dayMap[strings.ToUpper(val)]; ok {
			return comp(t.dayOfWeek, day), nil
		}
		return false, fmt.Errorf(
			"booleval DayOfWeekComparable.dowCompare: not a valid weekday: %s",
			val)
	case time.Time:
		return comp(val.Weekday(), t.dayOfWeek), nil
	}
	return false, fmt.Errorf("booleval DayOfWeekComparable.dowCompare: can't convert %v(%T) to a day of the week",
		other, other)
}

// Equal returns true if other is an integer and, if interpreted as a
// unix timestamp, it is equal to t.time, or if other is a time object
// and is equal to t.time
func (t DayOfWeekComparable) Equal(other any) (bool, error) {
	return t.dowCompare(
		func(lhs time.Weekday, rhs time.Weekday) bool {
			return lhs == rhs
		},
		other)
}

// Greater is true if other is an integer and t.time is greater than its value interpreted as
// a unix timestamp, or if other is a time.Time
// and t.time is greater than it.
func (t DayOfWeekComparable) Greater(other any) (bool, error) {
	return t.dowCompare(
		func(lhs time.Weekday, rhs time.Weekday) bool {
			return lhs > rhs
		},
		other)
}

type TimeOfDayComparable struct {
	timeSinceDayStart time.Duration
}

func NewTimeOfDayFromTimeString(str string) (TimeOfDayComparable, error) {
	if parsedTime, err := tryToParseTimeString(str); err == nil {
		return TimeOfDayComparable{timeToTimeOfDay(parsedTime)}, nil
	} else {
		return TimeOfDayComparable{}, err
	}

}

func timeToTimeOfDay(t time.Time) time.Duration {
	hours, minutes, _ := t.Clock()
	return (time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute)
}

func (t TimeOfDayComparable) todCompare(comp func(time.Duration, time.Duration) bool,
	other any) (bool, error) {
	switch val := other.(type) {
	case int64, uint, uint32, int, int32, uint64:
		intval, _ := getInt(val)
		timevalue := time.Unix(intval, 0)
		return t.todCompare(comp, timevalue)
	case time.Duration:
		return comp(t.timeSinceDayStart, val), nil
	case string:
		if time, err := tryToParseTimeString(val); err == nil {
			return comp(t.timeSinceDayStart, timeToTimeOfDay(time)), nil
		} else {
			return false, err
		}

	case time.Time:
		return comp(timeToTimeOfDay(val), t.timeSinceDayStart), nil
	}
	return false, fmt.Errorf("booleval TimeComparable.todCompare: can't convert %v to a time", other)
}

func (t TimeOfDayComparable) Equal(other any) (bool, error) {
	return t.todCompare(
		func(lhs time.Duration, rhs time.Duration) bool {
			return lhs == rhs
		},
		other)
}

func (t TimeOfDayComparable) Greater(other any) (bool, error) {
	return t.todCompare(
		func(lhs time.Duration, rhs time.Duration) bool {
			return lhs > rhs
		},
		other)
}
