package utils

type OutputResult struct {
	Target      string
	Version     string
	Players     string
	Description string
}

type FullResponse struct {
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	}

	Version struct {
		Name string `json:"name"`
	}

	Description string
}

type Response struct {
	Players struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`

	Version struct {
		Name string `json:"name"`
	} `json:"version"`
}

type ReponseMOTD struct {
	Description struct {
		Text string `json:"text"`
	}
}
