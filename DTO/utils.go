package DTO

import (
	"github.com/Gonnekone/onlineStore/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BrandRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DeviceRequest struct {
	ID    uint            `json:"id"`
	Name  string          `json:"name"`
	Price float64         `json:"price"`
	Type  TypeDTO         `json:"type"`
	Brand BrandDTO        `json:"brand"`
	Info  []DeviceInfoDTO `json:"info"`
}

type TypeRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserReg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (r UserReg) Validate() interface{} {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&r.Password, validation.Required),
	)
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l UserLogin) Validate() interface{} {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&l.Password, validation.Required),
	)
}

func BrandDTOToBrand(brandDTO BrandDTO) models.Brand {
	var brand models.Brand
	brand.Name = brandDTO.Name
	return brand
}

func DeviceInfoDTOToDeviceInfo(deviceInfoDTO []DeviceInfoDTO) []models.DeviceInfo {
	var deviceInfo []models.DeviceInfo
	for _, i := range deviceInfoDTO {
		var d models.DeviceInfo
		d.Description = i.Description
		d.Title = i.Title
		deviceInfo = append(deviceInfo, d)
	}
	return deviceInfo
}

func TypeDTOToType(typeDTO TypeDTO) models.Type {
	var t models.Type
	t.Name = typeDTO.Name
	return t
}

func DeviceDTOToDevice(deviceDTO DeviceDTO) models.Device {
	var device models.Device
	device.Name = deviceDTO.Name
	device.Brand = BrandDTOToBrand(deviceDTO.Brand)
	device.Price = deviceDTO.Price
	device.Info = DeviceInfoDTOToDeviceInfo(deviceDTO.Info)
	device.Type = TypeDTOToType(deviceDTO.Type)
	return device
}

func BrandToBrandRequest(brand models.Brand) BrandRequest {
	var brandRequest BrandRequest
	brandRequest.ID = brand.ID
	brandRequest.Name = brand.Name

	return brandRequest
}

func TypeToTypeRequest(t models.Type) TypeRequest {
	var typeRequest TypeRequest
	typeRequest.ID = t.ID
	typeRequest.Name = t.Name

	return typeRequest
}

func DeviceInfoToDeviceInfoDTO(deviceInfo []models.DeviceInfo) []DeviceInfoDTO {
	var deviceInfoDTO []DeviceInfoDTO
	for _, i := range deviceInfo {
		var d DeviceInfoDTO
		d.Description = i.Description
		d.Title = i.Title
		deviceInfoDTO = append(deviceInfoDTO, d)
	}
	return deviceInfoDTO
}

func TypeToTypeDTO(t models.Type) TypeDTO {
	var typeDTO TypeDTO
	typeDTO.Name = t.Name
	return typeDTO
}

func BrandToBrandDTO(brand models.Brand) BrandDTO {
	var b BrandDTO
	b.Name = brand.Name
	return b
}

func DeviceToDeviceDTOImage(device models.Device) DeviceDTOImage {
	var deviceDTOImage DeviceDTOImage
	deviceDTOImage.Info = DeviceInfoToDeviceInfoDTO(device.Info)
	deviceDTOImage.Type = TypeToTypeDTO(device.Type)
	deviceDTOImage.Name = device.Name
	deviceDTOImage.Price = device.Price
	deviceDTOImage.Brand = BrandToBrandDTO(device.Brand)
	deviceDTOImage.Img = device.Img
	return deviceDTOImage
}

func DeviceRequestToDevice(deviceRequest DeviceRequest, device models.Device) models.Device {
	device.Name = deviceRequest.Name
	device.Brand = BrandDTOToBrand(deviceRequest.Brand)
	device.Price = deviceRequest.Price
	device.Info = DeviceInfoDTOToDeviceInfo(deviceRequest.Info)
	device.Type = TypeDTOToType(deviceRequest.Type)
	return device
}
