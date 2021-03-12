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

type Ship struct {
	Name int
	From [2]int
	To   [2]int

	//	количество живых клеток
	HealthMax int
	Health    int
}

func (s ShipReq) Parse(matrixRange int) (map[int]Ship, error) {

	if strings.ContainsAny(s.Coordinates, ":()-abcdefghijklmnopqrstuvwxyz/\\") {
		return nil, errors.New("недопустимые символы в запросе")
	}

	coords, err := splitStr(s.Coordinates)
	if err != nil {
		return nil, err
	}

	return fillShips(coords, matrixRange)
}

func fillShips(coords [][2]int, matrixRange int) (map[int]Ship, error) {
	res := make(map[int]Ship)

	for i := 0; i < len(coords); i = i + 2 {
		x := coords[0]
		y := coords[1]

		if x[0] >= matrixRange || x[1] >= matrixRange || y[0] >= matrixRange || y[1] >= matrixRange ||
			x[0] < 0 || x[1] < 0 || y[0] < 0 || y[1] < 0 {
			return nil, errors.New("координата находится вне поля")
		}

		var hor, vert int

		hor = x[0] - y[0]
		vert = x[1] - y[1]

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

		res[len(res)+1] = Ship{
			Name:      len(res) + 1,
			From:      x,
			To:        y,
			HealthMax: health,
			Health:    health,
		}
	}

	return res, nil
}

// 1B 2C,3D 5D
func splitStr(in string) ([][2]int, error) {

	shipsCoordsStr := strings.Split(in, ",")

	res := make([][2]int, 0, len(shipsCoordsStr)*2)

	for _, pairs := range shipsCoordsStr {

		p := strings.Split(pairs, " ")
		if len(p) != 2 {
			return nil, errors.New("error pair parses")
		}

		head, err := getCoordinates(p[0])
		if err != nil {
			return nil, err
		}
		tail, err := getCoordinates(p[1])
		if err != nil {
			return nil, err
		}

		res = append(res, head, tail)
	}

	return res, nil
}

//26Z
func getCoordinates(in string) (res [2]int, err error) {
	if len(in) < 2 {
		return res, errors.New("wrong len")
	}

	str := in[:len(in)-1]

	res[0], err = strconv.Atoi(str)
	if err != nil {
		return res, err
	}

	res[1] = int(in[len(in)-1] - 'A' + 1)

	res[0]--
	res[1]--

	return res, nil
}
