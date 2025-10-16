package repo

type Score struct{
	PackagingScore int `json:"packaging_score" db:"packaging_score"`
	LocationScore int `json:"" db:"location_score"`
	DisposalScore int 
}