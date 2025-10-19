package logic

import (
	"ecoscan.com/repo"
)

func CalculateScore(product repo.Product) {
	packagingScore := calculatePackagingScore(product.PackagingMaterial)
	transportScore := calculateTransportScore(product.ManufacturingLocation)
	disposalScore := calculateDisposalScore(product.DisposalMethod)

	overallScore := (float64(packagingScore) * 0.35) + (float64(transportScore) * 0.30) + (float64(disposalScore) * 0.35)

	return overallScore
}

func calculatePackagingScore(material string) int {
	switch material {
	case "none", "compostable_paper":
		return 100
	case "glass", "paper", "cardboard":
		return 80
	case "aluminum", "recyclable_plastic":
		return 60
	case "plastic", "mixed_materials":
		return 20
	default:
		return 40 //for unknown material
	}
}

func calculateTransportScore(location string) int {
	switch location {
	case "local", "regional":
		return 95
	case "national":
		return 70
	case "international":
		return 30
	default:
		return 50
	}
}

func calculateDisposalScore(method string) int {
	switch method {
	case "compostable", "reusable": 
		return 100 
	case "recyclable": 
		return 85 
	case "minimal_impact":
		return 70 
	case "landfill":
		return 10 
	default:
		return 40
	}
}