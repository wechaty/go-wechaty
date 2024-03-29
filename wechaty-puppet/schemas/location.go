package schemas

type LocationPayload struct {
	Accuracy  float32 `json:"accuracy"`  // Estimated horizontal accuracy of this location, radial, in meters. (same as Android & iOS API)
	Address   string  `json:"address"`   // 北京市北京市海淀区45 Chengfu Rd
	Latitude  float64 `json:"latitude"`  // 39.995120999999997
	Longitude float64 `json:"longitude"` // 116.3341
	Name      string  `json:"name"`      // 东升乡人民政府(海淀区成府路45号)
}
