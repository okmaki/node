package core

import (
	"errors"
	"fmt"
)

type LogLevel = uint8

const (
	LevelDebug uint8 = 1 << iota
	LevelInfo
	LevelWarning
	LevelError
)

func IsValidLogLevel(level LogLevel) bool {
	switch level {
	case LevelDebug:
	case LevelInfo:
	case LevelWarning:
	case LevelError:
	default:
		return false
	}

	return true
}

func ContainsLevel(levels LogLevel, level LogLevel) bool {
	return levels&level == level
}

type IdValidatorFunc = func(field string, value string) error

type Log struct {
	SourceId             string   `json:"sourceId"`
	SourceType           string   `json:"sourceType"`
	TransactionId        string   `json:"transactionId"`
	TransactionTimestamp int64    `json:"transactionTimestamp"`
	Level                LogLevel `json:"level"`
	Timestamp            int64    `json:"timestamp"`
	Location             string   `json:"location"`
	Data                 string   `json:"data"`
}

func (l Log) Validate(validateId IdValidatorFunc) error {
	if err := validateId("source id", l.SourceId); err != nil {
		return err
	}

	if l.SourceType == "" {
		return errors.New("source type required")
	}

	if l.TransactionId != "" {
		if err := validateId("transaction id", l.TransactionId); err != nil {
			return err
		}
	}

	if !IsValidLogLevel(l.Level) {
		return errors.New(fmt.Sprint("invalid tranaction timestamp", l.TransactionTimestamp))
	}

	return nil
}

type LogFilter struct {
	SourceType    string   `json:"sourceType"`
	SourceId      string   `json:"sourceId"`
	TransactionId string   `json:"transactionId"`
	Levels        LogLevel `json:"levels"`
	After         int64    `json:"after"`
	Before        int64    `json:"before"`
}

func NewLogFilter() LogFilter {
	return LogFilter{
		SourceType:    "",
		SourceId:      "",
		TransactionId: "",
		Levels:        0,
		After:         0,
		Before:        0,
	}
}

func (f LogFilter) Validate(validateId IdValidatorFunc) error {
	if f.SourceType == "" && f.TransactionId == "" {
		return errors.New("source type or transaction id is required")
	}

	if f.TransactionId != "" {
		if err := validateId("transaction id", f.TransactionId); err != nil {
			return err
		}
	}

	if f.SourceId != "" {
		if err := validateId("source id", f.SourceId); err != nil {
			return err
		}
	}

	return nil
}

func (f LogFilter) HasLevel(level LogLevel) bool {
	return ContainsLevel(f.Levels, level)
}

func (f LogFilter) ShouldFilter(log Log) bool {
	if f.SourceType != "" && log.SourceType != f.SourceType {
		return true
	}

	if f.TransactionId != "" && log.TransactionId != f.TransactionId {
		return true
	}

	if f.SourceId != "" && log.SourceId != f.SourceId {
		return true
	}

	if f.Levels > 0 && !f.HasLevel(log.Level) {
		return true
	}

	if f.After > 0 && !(log.Timestamp > f.After) {
		return true
	}

	if f.Before > 0 && !(log.Timestamp < f.Before) {
		return true
	}

	return false
}
