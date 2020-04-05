package dto

/*
{
	"ship_count":int, // всего кораблей
	"destroyed" :int, // потоплено
	"knocked"   :int, // подбито
	"shot_count":int  // сделано выстрелов
}
*/
type State struct {
	ShipCount int `json:"ship_count" form:"ship_count"`
	Destroyed int `json:"destroyed" form:"destroyed"`
	Knocked   int `json:"knocked" form:"knocked"`
	ShotCount int `json:"shot_count" form:"shot_count"`
}
