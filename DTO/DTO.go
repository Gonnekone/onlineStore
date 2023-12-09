package DTO

type BrandDTO struct {
	Name string `json:"name"`
}

type DeviceInfoDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TypeDTO struct {
	Name string `json:"name"`
}

type DeviceDTOImage struct {
	Name  string          `json:"name"`
	Price float64         `json:"price"`
	Img   string          `json:"img"`
	Type  TypeDTO         `json:"type"`
	Brand BrandDTO        `json:"brand"`
	Info  []DeviceInfoDTO `json:"info"`
}

type DeviceDTO struct {
	Name  string          `json:"name"`
	Price float64         `json:"price"`
	Type  TypeDTO         `json:"type"`
	Brand BrandDTO        `json:"brand"`
	Info  []DeviceInfoDTO `json:"info"`
}
