package model

// Image - структура представляющая исходное изображение и требуемые размеры.
type Image struct {
	Source    string // URL источника изображения.
	GetParams string // Get-параметры запроса.
	Width     int    // Требуемая ширина изображения.
	Height    int    // Требуемая длина изображения.
	Data      []byte // Изображение в виде массива байтов.
}

type Headers map[string][]string

type Request struct {
	URL     string
	Params  string
	Headers Headers
}
