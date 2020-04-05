package dto

import "errors"

type CreateMatrix struct {
	MatrixRange int `json:"range" form:"range"`
}

func (c CreateMatrix) Validate() error {
	if c.MatrixRange < 1 || c.MatrixRange > 26 {
		return errors.New("размер поля должен быть больше 1 и меньше 26")
	}

	return nil
}
