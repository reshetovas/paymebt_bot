package utils

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func ParseCustomReportDates(input string) (time.Time, time.Time, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 6 {
		log.Error().Msg("Некорректный формат дат, ожидается ГГГГ-ММ-ДД - ГГГГ-ММ-ДД")
		log.Debug().Msgf("INPUT: type: %T, value: %v\n Len parts %d", input, input, len(parts))
		return time.Time{}, time.Time{}, fmt.Errorf("Некорректный формат дат, ожидается ГГГГ-ММ-ДД - ГГГГ-ММ-ДД")
	}

	startStr := strings.TrimSpace(parts[0] + "-" + parts[1] + "-" + parts[2])
	endStr := strings.TrimSpace(parts[3] + "-" + parts[4] + "-" + parts[5])

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("не удалось распарсить дату начала")
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("не удалось распарсить дату конца")
	}

	return start, end, nil
}
