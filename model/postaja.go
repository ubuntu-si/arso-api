package model

// Postaja holds info about weather
type Postaja struct {
	ID            string  `xml:"metData>domain_meteosiId"`
	Title         string  `xml:"metData>domain_longTitle"`
	Lat           float64 `xml:"metData>domain_lat"`
	Lon           float64 `xml:"metData>domain_lon"`
	Altitude      float64 `xml:"metData>domain_altitude"`
	Issued        string  `xml:"metData>tsUpdated_RFC822"`
	Temp          float64 `xml:"metData>t"`
	Wind          float64 `xml:"metData>ff_val" json:",omitempty"`
	WindDirection string  `xml:"metData>dd_icon" json:",omitempty"`
	RH            float64 `xml:"metData>rh" json:",omitempty"`
	Pressure      float64 `xml:"metData>p" json:",omitempty"`
	Sky           string  `xml:"metData>nn_shortText" json:",omitempty"`
	Valid         string  `xml:"metData>tsValid_issued_UTC"`
	Auto          bool
}
