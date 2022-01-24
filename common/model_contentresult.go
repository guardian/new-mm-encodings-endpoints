package common

type ContentResult struct {
	Encoding

	RealMimeName string `json:"real_name"` //optional string containing the MIME type equivalent
	PosterURL    string `json:"posterurl"` //optional string containing a URL for the poster image if it exists
}
