package services

import "frugal-hero/outputs"

type IService interface {
	Inspect(output outputs.OutputInterface)
}