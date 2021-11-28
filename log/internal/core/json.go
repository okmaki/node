package core

import "encoding/json"

func GetJson(model interface{}) string {
	data, err := json.Marshal(model)
	if err != nil {
		return "{}"
	}

	return string(data)
}
