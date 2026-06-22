package services

import "SungClip/internal/utils"

type services struct {
	utils *utils.Utils
}

func NewServices(utils *utils.Utils) IServices {
	return &services{
		utils: utils,
	}
}