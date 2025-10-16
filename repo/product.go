package repo

type Product struct{
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	BrandName string `json:"brand_name" db:"brand_name"`
	Category string `json:"category" db:"category"`
	SubCatergory string `json:"sub_category" db:"sub_category"`
	ImageURL string `json:"image_url" db:"image_url"`
	Price float32 `json:"price" db:"price"`
	//Score int `json:"score" db:"score"`
	PackagingMaterial string `json:"packaging_material" db:"packaging_material"`
	ManufacturingLocation string `json:"manufacturing_location" db:"manufacturing_location"`
	DisposalMethod string `json:"disposal_method" db:"disposal_method"`
}
