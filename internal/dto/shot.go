package dto

import (
	"errors"
	"strconv"
	"strings"
)

//{"сoord": string}
type ShotReq struct {
	Coordinates string `json:"сoord" form:"сoord"`
}

func (s ShotReq) Parse(matrixRange int) (x, y int, err error) {

	if len(s.Coordinates) == 0 {
		return 0, 0, errors.New("пустой запрос")
	}

	if strings.ContainsAny(s.Coordinates, " :()-abcdefghijklmnopqrstuvwxyz/\\") {
		return 0, 0, errors.New("недопустимые символы в запросе")
	}

	for key, value := range dict {
		s.Coordinates = strings.ReplaceAll(s.Coordinates, key, value)
	}

	coordsStr := strings.Split(s.Coordinates, ":")

	if len(coordsStr) != 2 {
		return 0, 0, errors.New("недопустимые символы в запросе")
	}

	x, err = strconv.Atoi(coordsStr[0])
	if err != nil {
		return 0, 0, err
	}

	y, err = strconv.Atoi(coordsStr[1])
	if err != nil {
		return 0, 0, err
	}

	if x <= 0 || y <= 0 || x > matrixRange || y > matrixRange {
		return 0, 0, errors.New("координата находится вне поля")
	}

	return x - 1, y - 1, nil
}

//{"destroy":false,"knock":true,"end":false},
type ShotRes struct {
	Destroy bool `json:"destroy" form:"destroy"`
	Knock   bool `json:"knock" form:"knock"`
	End     bool `json:"end" form:"end"`
}
