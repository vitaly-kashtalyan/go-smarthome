package responses

type RelayStatus struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Relay `json:"data"`
}

type Relay struct {
	Id    int32 `json:"id"`
	State int32 `json:"state"`
}

type dht22 struct {
	Pin         int     `json:"pin"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`
	Status      string  `json:"status"`
}
type ds18b20 struct {
	Pin         int     `json:"pin"`
	Temperature float32 `json:"temperature"`
	Dec         string  `json:"dec"`
	Status      string  `json:"status"`
}

type ArduinoSensors struct {
	Dht22   []dht22   `json:"dht22"`
	Ds18b20 []ds18b20 `json:"ds18b20"`
}
