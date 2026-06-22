package controllers

import (
	"SungClip/internal/services"
	"SungClip/internal/utils"
)

type controllers struct {
	utils *utils.Utils
	services services.IServices
}

func NewControllers(utils *utils.Utils, services services.IServices) IControllers {
	return &controllers{
		utils: utils,
		services: services,
	}
}