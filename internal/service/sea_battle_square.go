package service

import (
	"errors"
	"local/sea_fight/internal/constants"
	"local/sea_fight/internal/dto"
	"sync"
)

type SeaBattleSquare interface {
	CreateMatrix(matrixRange int) error
	GetRange() int
	Ship(ships map[int]dto.Ship) error
	Shot(x, y int) (dto.ShotRes, error)
	Clear()
	State() dto.State
}

type seaBattleSquareService struct {
	isPlaced bool
	isEnded  bool

	mut sync.RWMutex

	matrix matrixType
	ships  map[int]dto.Ship

	shipCount int
	destroyed int
	shotCount int
}

type matrixType []matrixRow
type matrixRow []int

func newSeaBattleSquareService() SeaBattleSquare {
	return &seaBattleSquareService{
		isPlaced: false,
		mut:      sync.RWMutex{},
	}
}

func (s *seaBattleSquareService) CreateMatrix(matrixRange int) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.isPlaced && !s.isEnded {
		return errors.New("игра уже начата. сначала закончите или очистите предыдущую игру")
	}

	s.clear(matrixRange)

	return nil
}

func (s *seaBattleSquareService) GetRange() int {
	s.mut.RLock()
	defer s.mut.RUnlock()

	return len(s.matrix)
}

func (s *seaBattleSquareService) Ship(ships map[int]dto.Ship) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.isEnded {
		s.clear(len(s.matrix))
	}

	if s.isPlaced {
		return errors.New("игра уже начата. сначала закончите или очистите предыдущую игру")
	}

	if len(s.matrix) == 0 {
		return errors.New("размер поля не задан")
	}

	for _, ship := range ships {

		var fromX, toX, fromY, toY int

		if ship.From[0] > ship.To[0] {
			fromX = ship.To[0]
			toX = ship.From[0]
		} else {
			fromX = ship.From[0]
			toX = ship.To[0]
		}

		if ship.From[1] > ship.To[1] {
			fromY = ship.To[1]
			toY = ship.From[1]
		} else {
			fromY = ship.From[1]
			toY = ship.To[1]
		}

		for i := fromX; i <= toX; i++ {
			for j := fromY; j <= toY; j++ {

				if s.matrix[i][j] != 0 {

					//сброс матрицы
					s.matrix = make(matrixType, len(s.matrix))
					for i := range s.matrix {
						s.matrix[i] = make(matrixRow, len(s.matrix))
					}

					return errors.New("корабли наслаиваются друг на друга")
				}

				s.matrix[i][j] = ship.Name
			}
		}

	}

	s.ships = ships
	s.shipCount = len(ships)
	s.isPlaced = true

	return nil
}

func (s *seaBattleSquareService) Shot(x, y int) (dto.ShotRes, error) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.isEnded {
		return dto.ShotRes{}, errors.New("игра закончена")
	}

	t := s.matrix[x][y]

	switch t {
	case constants.IsSpent:
		return dto.ShotRes{}, errors.New("сюда уже стреляли")
	case 0:
		s.matrix[x][y] = constants.IsSpent

		s.shotCount++

		return dto.ShotRes{
			Destroy: false,
			Knock:   false,
			End:     false,
		}, nil
	}

	ship, ok := s.ships[t]
	if !ok {
		return dto.ShotRes{}, errors.New("неизвестная ошибка, найден неизвестный корабль")
	}

	s.shotCount++

	ship.Health--
	s.ships[t] = ship
	//если нужно помнить где стояли корабли, то можно другой флаг ставить
	s.matrix[x][y] = constants.IsSpent

	// ранили корабль
	if ship.Health > 0 {
		return dto.ShotRes{
			Destroy: false,
			Knock:   true,
			End:     false,
		}, nil
	}

	s.destroyed++

	//потопили корабль
	if s.destroyed != s.shipCount {
		return dto.ShotRes{
			Destroy: true,
			Knock:   true,
			End:     false,
		}, nil
	}

	//закончили игру

	s.isEnded = true

	return dto.ShotRes{
		Destroy: true,
		Knock:   true,
		End:     true,
	}, nil
}

func (s *seaBattleSquareService) Clear() {
	s.mut.Lock()
	defer s.mut.Unlock()

	if len(s.matrix) == 0 {
		return
	}

	s.clear(len(s.matrix))
}

func (s *seaBattleSquareService) clear(matrixRange int) {

	s.ships = map[int]dto.Ship{}
	s.isPlaced = false
	s.isEnded = false
	s.shipCount = 0
	s.destroyed = 0
	s.shotCount = 0

	//сброс матрицы
	s.matrix = make(matrixType, matrixRange)
	for i := range s.matrix {
		s.matrix[i] = make(matrixRow, matrixRange)
	}
}

func (s *seaBattleSquareService) State() dto.State {
	s.mut.RLock()
	defer s.mut.RUnlock()

	var knocked int

	for _, ship := range s.ships {
		if ship.Health > 0 && ship.HealthMax != ship.Health {
			knocked++
		}
	}

	return dto.State{
		ShipCount: s.shipCount,
		Destroyed: s.destroyed,
		Knocked:   knocked,
		ShotCount: s.shotCount,
	}
}
