package data

type Id string

type Poi struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SearchArea struct {
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	RadiusInMeter uint64  `json:"radius"`
}

type Pois []Poi
