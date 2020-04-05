package dto

import (
	"errors"
	"strconv"
	"strings"
)

//{"Coordinates": string}
type ShipReq struct {
	Coordinates string `json:"coordinates" form:"coordinates"`
}

var dict = map[string]string{
	"A": ":1",
	"B": ":2",
	"C": ":3",
	"D": ":4",
	"E": ":5",
	"F": ":6",
	"G": ":7",
	"H": ":8",
	"I": ":9",
	"J": ":10",
	"K": ":11",
	"L": ":12",
	"M": ":13",
	"N": ":14",
	"O": ":15",
	"P": ":16",
	"Q": ":17",
	"R": ":18",
	"S": ":19",
	"T": ":20",
	"U": ":21",
	"V": ":22",
	"W": ":23",
	"X": ":24",
	"Y": ":25",
	"Z": ":26",
}

func (s ShipReq) Parse(matrixRange int) (map[int]Ship, error) {

	if strings.ContainsAny(s.Coordinates, ":()-abcdefghijklmnopqrstuvwxyz/\\") {
		return nil, errors.New("недопустимые символы в запросе")
	}

	for key, value := range dict {
		s.Coordinates = strings.ReplaceAll(s.Coordinates, key, value)
	}

	shipsCoordsStr := strings.Split(s.Coordinates, ",")

	out := make(map[int]Ship)

	for _, value := range shipsCoordsStr {
		fromToStr := strings.Split(value, " ")
		if len(fromToStr) != 2 {
			return nil, errors.New("неверное количество координат для корабля")
		}

		xStr := strings.Split(fromToStr[0], ":")
		if len(xStr) != 2 {
			return nil, errors.New("неверные данные")
		}

		yStr := strings.Split(fromToStr[1], ":")
		if len(yStr) != 2 {
			return nil, errors.New("неверные данные")
		}

		x1, err := strconv.Atoi(xStr[0])
		if err != nil {
			return nil, err
		}

		x2, err := strconv.Atoi(xStr[1])
		if err != nil {
			return nil, err
		}

		y1, err := strconv.Atoi(yStr[0])
		if err != nil {
			return nil, err
		}

		y2, err := strconv.Atoi(yStr[1])
		if err != nil {
			return nil, err
		}

		if x1 > matrixRange || x2 > matrixRange || y1 > matrixRange || y2 > matrixRange ||
			x1 <= 0 || x2 <= 0 || y1 <= 0 || y2 <= 0 {
			return nil, errors.New("координата находится вне поля")
		}

		var hor, vert int

		hor = x1 - y1
		vert = x2 - y2

		if hor >= 0 {
			hor++
		} else {
			hor = -hor + 1
		}

		if vert >= 0 {
			vert++
		} else {
			vert = -vert + 1
		}

		health := hor * vert

		out[len(out)+1] = Ship{
			Name:      len(out) + 1,
			From:      [2]int{x1 - 1, x2 - 1},
			To:        [2]int{y1 - 1, y2 - 1},
			HealthMax: health,
			Health:    health,
		}
	}

	return out, nil
}

type Ship struct {
	Name int
	From [2]int
	To   [2]int

	//	количество живых клеток
	HealthMax int
	Health    int
}
