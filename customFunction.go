package govaluate

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type dateType struct {
	year   int
	month  int
	day    int
	hour   int
	minute int
	second int
}

func castDateTypeToString(date dateType) string {
	result := fmt.Sprintf("%d-", date.year)
	if date.month < 10 {
		result += "0"
	}
	result = fmt.Sprintf("%s%d-", result, date.month)
	if date.day < 10 {
		result += "0"
	}
	result = fmt.Sprintf("%s%d", result, date.day)
	return result
}

func castTimeStampToDateType(timestamp int64) (dateType, error) {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	timeInstance := time.Unix(timestamp, 0).In(bangkok)
	if err != nil {
		return dateType{}, err
	}
	return dateType{
		year:   timeInstance.Year(),
		month:  int(timeInstance.Month()),
		day:    timeInstance.Day(),
		hour:   timeInstance.Hour(),
		minute: timeInstance.Minute(),
		second: timeInstance.Second(),
	}, nil
}

func castToDateType(dateString string) (dateType, error) {
	datetimeSplitted := strings.Split(dateString, " ")
	splittedDate := strings.Split(datetimeSplitted[0], "-")
	hour := 0
	minute := 0
	second := 0

	if len(splittedDate) != 3 {
		return dateType{}, fmt.Errorf("datepart: %s is incorrect format (yyyy-mm-dd hh:mm:ss or yyyy-mm-dd)", dateString)
	}
	year, err := strconv.Atoi(strings.ReplaceAll(splittedDate[0], "\"", ""))
	if err != nil {
		return dateType{}, fmt.Errorf("data type’s input field is not date/datetime at year")
	}
	month, err := strconv.Atoi(strings.ReplaceAll(splittedDate[1], "\"", ""))
	if err != nil {
		return dateType{}, fmt.Errorf("data type’s input field is not date/datetime at month")
	}
	day, err := strconv.Atoi(strings.ReplaceAll(splittedDate[2], "\"", ""))
	if err != nil && len(splittedDate[2]) > 2 {
		return dateType{}, fmt.Errorf("data type’s input field is not date/datetime at day")
	}

	if len(datetimeSplitted) > 1 {
		splittedTime := strings.Split(datetimeSplitted[1], ":")

		if len(splittedTime) != 3 {
			return dateType{}, fmt.Errorf("datepart: %s is incorrect format (yyyy-mm-dd hh:mm:ss or yyyy-mm-dd)", dateString)
		}
		hour, err = strconv.Atoi(splittedTime[0])
		if err != nil {
			return dateType{}, fmt.Errorf("data type’s input field is not date/datetime at hour")
		}
		minute, err = strconv.Atoi(splittedTime[1])
		if err != nil {
			return dateType{}, fmt.Errorf("data type’s input field is not date/datetime at minute")
		}
		second, err = strconv.Atoi(splittedTime[2])
		if err != nil {
			return dateType{}, fmt.Errorf("data type’s input field is not date/datetime at second")
		}
	}

	return dateType{
		year:   year,
		month:  month,
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	}, nil
}

// ===================================================================

var minCustomFunction = func(args ...interface{}) (result interface{}, err error) {
	if len(args) == 0 {
		return "", fmt.Errorf("cannot find minimum of empty list")
	}
	result = math.Inf(0)
	for _, arg := range args {
		var value float64
		switch tmp := arg.(type) {
		case string:
			value, err = strconv.ParseFloat(tmp, 64)
			if err != nil {
				return "", err
			}
		case float64:
			value = tmp
		case float32:
			value = float64(tmp)
		case int:
			value = float64(tmp)
		case int32:
			value = float64(tmp)
		case int64:
			value = float64(tmp)
		default:
			return "", fmt.Errorf("wrong argument type to minimum function")
		}

		result = math.Min(result.(float64), value)
	}
	return fmt.Sprintf("%v", result), nil
}

var maxCustomFunction = func(args ...interface{}) (result interface{}, err error) {
	if len(args) == 0 {
		return "", fmt.Errorf("cannot find minimum of empty list")
	}
	result = math.Inf(-1)
	for _, arg := range args {
		var value float64
		switch tmp := arg.(type) {
		case string:
			value, err = strconv.ParseFloat(tmp, 64)
			if err != nil {
				return "", err
			}
		case float64:
			value = tmp
		case float32:
			value = float64(tmp)
		case int:
			value = float64(tmp)
		case int32:
			value = float64(tmp)
		case int64:
			value = float64(tmp)
		default:
			return "", fmt.Errorf("wrong argument type to minimum function")
		}

		result = math.Max(result.(float64), value)
	}
	return fmt.Sprintf("%v", result), nil
}

var dateDiffFunction = func(args ...interface{}) (result interface{}, err error) {
	if len(args) != 3 {
		return "", fmt.Errorf("wrong number of argument for datediff")
	}
	var date1, date2 dateType
	var mode string

	switch tmp := args[0].(type) {
	case string:
		date1, err = castToDateType(tmp)
		if err != nil {
			return nil, err
		}
	case float64:
		date1, err = castTimeStampToDateType(int64(tmp))
		if err != nil {
			return nil, err
		}
	case int64:
		date1, err = castTimeStampToDateType(tmp)
		if err != nil {
			return nil, err
		}
	default:
		return "", fmt.Errorf("wrong first argument type to datediff")
	}

	switch tmp := args[1].(type) {
	case string:
		date2, err = castToDateType(tmp)
		if err != nil {
			return nil, err
		}
	case float64:
		date2, err = castTimeStampToDateType(int64(tmp))
		if err != nil {
			return nil, err
		}
	case int64:
		date2, err = castTimeStampToDateType(tmp)
		if err != nil {
			return nil, err
		}
	default:
		return "", fmt.Errorf("wrong second argument type to datediff")
	}

	switch tmp := args[2].(type) {
	case string:
		mode = tmp
	default:
		return "", fmt.Errorf("wrong third argument type to datediff")
	}

	// calculate

	// case 1: return abs(date2 - date1)
	// maxDate, minDate := maxMinDate(date1, date2)

	// case 2: return date2 - date1
	maxDate := date2
	minDate := date1

	valueInt := 0
	if mode == "\"y\"" {
		valueInt = maxDate.year - minDate.year
		if maxDate.month < minDate.month || (maxDate.month == minDate.month && maxDate.day < minDate.day) {
			valueInt--
		}
	} else if mode == "\"m\"" {
		valueInt = ((maxDate.year - minDate.year) * 12) + (maxDate.month - minDate.month)
		if maxDate.day < minDate.day {
			valueInt--
		}
	} else {
		maxEpoch := time.Date(maxDate.year, time.Month(maxDate.month), maxDate.day, 0, 0, 0, 0, time.UTC)
		minEpoch2 := time.Date(minDate.year, time.Month(minDate.month), minDate.day, 0, 0, 0, 0, time.UTC)
		valueInt = int(maxEpoch.Sub(minEpoch2).Hours() / 24)
	}
	result = fmt.Sprintf("(%s)", strconv.FormatInt(int64(valueInt), 10))
	return result, nil
}

var dateAddFunction = func(args ...interface{}) (result interface{}, err error) {
	if len(args) != 3 {
		return "", fmt.Errorf("wrong number of argument for dateadd")
	}
	var dateInput dateType
	var offsetResult float64
	var mode string

	switch tmp := args[0].(type) {
	case string:
		dateInput, err = castToDateType(tmp)
		if err != nil {
			return nil, err
		}
	case float64:
		dateInput, err = castTimeStampToDateType(int64(tmp))
		if err != nil {
			return nil, err
		}
	case int64:
		dateInput, err = castTimeStampToDateType(tmp)
		if err != nil {
			return nil, err
		}
	default:
		return "", fmt.Errorf("wrong first argument type to dateadd")
	}

	switch tmp := args[1].(type) {
	case float64:
		offsetResult = tmp
	default:
		return "", fmt.Errorf("wrong second argument type to dateadd")
	}

	switch tmp := args[2].(type) {
	case string:
		mode = tmp
	default:
		return "", fmt.Errorf("wrong third argument type to dateadd")
	}

	offsetInt := int(offsetResult)
	yearOffset := 0
	monthOffset := 0
	dayOffSet := 0
	// calculate
	if mode == "\"y\"" {
		yearOffset = offsetInt
	} else if mode == "\"m\"" {
		monthOffset = offsetInt
	} else {
		dayOffSet = offsetInt
	}
	// must check, if day exceed capability of month
	epoch := time.Date(
		dateInput.year+yearOffset,
		time.Month(dateInput.month+monthOffset),
		dateInput.day+dayOffSet,
		0, 0, 0, 0, time.UTC,
	)
	dateResult := dateType{
		year:  epoch.Year(),
		month: int(epoch.Month()),
		day:   epoch.Day(),
	}
	return castDateTypeToString(dateResult), nil
}

var datePartFunction = func(args ...interface{}) (result interface{}, err error) {
	if len(args) != 2 {
		return "", fmt.Errorf("wrong number of argument for datepart")
	}
	var dateInput dateType
	var mode string

	switch tmp := args[0].(type) {
	case string:
		dateInput, err = castToDateType(tmp)
		if err != nil {
			return nil, err
		}
	case float64:
		dateInput, err = castTimeStampToDateType(int64(tmp))
		if err != nil {
			return nil, err
		}
	case int64:
		dateInput, err = castTimeStampToDateType(tmp)
		if err != nil {
			return nil, err
		}
	default:
		return "", fmt.Errorf("wrong first argument type to datepart")
	}

	switch tmp := args[1].(type) {
	case string:
		mode = tmp
	default:
		return "", fmt.Errorf("wrong second argument type to datepart")
	}

	switch mode {
	case "date":
		result = fmt.Sprintf("%d", dateInput.day)
	case "day": // day of week
		instance := time.Date(dateInput.year, time.Month(dateInput.month), dateInput.day, 0, 0, 0, 0, time.Local)
		result = fmt.Sprintf("%d", instance.Weekday()+1)
	case "month":
		result = fmt.Sprintf("%d", dateInput.month)
	case "year":
		result = fmt.Sprintf("%d", dateInput.year)
	case "hour":
		result = fmt.Sprintf("%d", dateInput.hour)
	case "minute":
		result = fmt.Sprintf("%d", dateInput.minute)
	case "second":
		result = fmt.Sprintf("%d", dateInput.second) // panic(createSyntaxErrorMessage("Unimplemented mode for datepart, since there is no time in value", formula))
	default:
		return "", fmt.Errorf("unrecognized mode to datepart")
	}

	return result, nil
}
