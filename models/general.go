package models

type ResponseSuccess struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApplicationError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Success    bool   `json:"success"`
}

type Message struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
