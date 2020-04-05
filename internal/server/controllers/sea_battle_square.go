package controllers

import (
	"github.com/gin-gonic/gin"
	"local/sea_fight/internal/dto"
	"local/sea_fight/internal/server"
	service "local/sea_fight/internal/service"
	"net/http"
)

type seaBattleController struct {
	battleService service.SeaBattleSquare
	actions       []server.Action
}

func NewSeaBattleSquareController(battleService service.SeaBattleSquare) *seaBattleController {

	c := seaBattleController{
		battleService: battleService,
	}

	c.actions = append(c.actions,
		server.Action{
			HttpMethod:   http.MethodPost,
			RelativePath: "/api/sea_battle_square/create-matrix",
			ActionExec:   c.createMatrix,
		},
		server.Action{
			HttpMethod:   http.MethodPost,
			RelativePath: "/api/sea_battle_square/ship",
			ActionExec:   c.ship,
		},
		server.Action{
			HttpMethod:   http.MethodPost,
			RelativePath: "/api/sea_battle_square/shot",
			ActionExec:   c.shot,
		},
		server.Action{
			HttpMethod:   http.MethodPost,
			RelativePath: "/api/sea_battle_square/autoshooting",
			ActionExec:   c.autoShooting,
		},
		server.Action{
			HttpMethod:   http.MethodPost,
			RelativePath: "/api/sea_battle_square/clear",
			ActionExec:   c.clear,
		},
		server.Action{
			HttpMethod:   http.MethodGet,
			RelativePath: "/api/sea_battle_square/state",
			ActionExec:   c.state,
		},
	)

	return &c
}

func (s *seaBattleController) Actions() []server.Action {
	return s.actions
}

func (s *seaBattleController) createMatrix(ctx *gin.Context) {

	var createMatrix dto.CreateMatrix

	err := ctx.ShouldBind(&createMatrix)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	err = createMatrix.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	err = s.battleService.CreateMatrix(createMatrix.MatrixRange)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Ok",
	})
}

func (s *seaBattleController) ship(ctx *gin.Context) {
	var req dto.ShipReq

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	ships, err := req.Parse(s.battleService.GetRange())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	err = s.battleService.Ship(ships)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Ok"})
}

func (s *seaBattleController) shot(ctx *gin.Context) {

	var req dto.ShotReq

	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	x, y, err := req.Parse(s.battleService.GetRange())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	res, err := s.battleService.Shot(x, y)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"msg": err.Error()}})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *seaBattleController) autoShooting(ctx *gin.Context) {

	matrixRange := s.battleService.GetRange()

	results := make([]interface{}, 0, matrixRange*matrixRange)

	for i := 0; i < matrixRange; i++ {
		for j := 0; j < matrixRange; j++ {
			res, err := s.battleService.Shot(i, j)
			if err != nil {
				results = append(results, gin.H{"msg": err.Error()})
			} else {
				results = append(results, res)
			}
		}
	}

	ctx.JSON(http.StatusOK, results)
}

func (s *seaBattleController) clear(ctx *gin.Context) {

	s.battleService.Clear()

	ctx.JSON(http.StatusOK, gin.H{"status": "Ok"})
}

func (s *seaBattleController) state(ctx *gin.Context) {

	res := s.battleService.State()

	ctx.JSON(http.StatusOK, res)
}
