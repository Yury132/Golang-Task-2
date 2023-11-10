package models

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

type ImageMeta struct {
	Name   string
	Type   string
	Height int
	Width  int
}

// Отображение информации о загруженных картинках и созданных миниатюрах
type AllImages struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	NameMini   string `json:"name_miniature"`
	WidthMini  int    `json:"width_miniature"`
	HeightMini int    `json:"height_miniature"`
}

// Получаем данные о картинке
func CollectImageMeta(data []byte, name string) (*ImageMeta, error) {
	// Из байтов декодируем изображение
	r := bytes.NewReader(data)
	//r.Seek(0, 0)
	imageData, imageType, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	b := imageData.Bounds()
	return &ImageMeta{
		Name:   name,
		Type:   imageType,
		Height: b.Max.Y,
		Width:  b.Max.X,
	}, nil
}

type InfoForThumbnail struct {
	Path string `json:"path"`
	Size int    `json:"size"`
}
