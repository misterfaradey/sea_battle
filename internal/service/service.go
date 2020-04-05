package service

type SeaBattleService interface {
	SeaBattleSquare
}

func NewSeaBattleService() SeaBattleService {
	return newSeaBattleSquareService()
}
