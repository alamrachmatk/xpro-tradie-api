package lib

import (
	"net/http"

	"github.com/labstack/echo"
)

type Status struct {
	Code          uint   `json:"code"`
	MessageServer string `json:"message_server"`
	MessageClient string `json:"message_client"`
}

type Response struct {
	Status Status      `json:"status"`
	Meta   *Pagination `json:"meta"`
	Data   interface{} `json:"data"`
}

type ResponseToken struct {
	Status 		Status      `json:"status"`
	Email		string 		`json:"email"`
	Token		string 		`json:"token"`
}

// CustomError Return custom error. Param: (code, message_server, message_client).
//
// Parameters:
// code         Error code
// message[0]   Message server (return as parameter message_server)
// message[1]   Message client (return as parameter message_client)
//
// Output:
// {
//   "$schema": "http://json-schema.org/draft-06/schema#",
//   "properties": {
//     "status": {
//       "properties": {
//         "code": {
//           "type": "integer"
//         },
//         "message_client": {
//           "type": "string"
//         },
//         "message_server": {
//           "type": "string"
//         }
//       },
//       "required": [
//         "code",
//         "message_client",
//         "message_server"
//       ],
//       "type": "object"
//     },
//     "meta": {
//       "type": "null"
//     },
//     "data": {
//       "type": "null"
//     }
//   },
//   "required": [
//       "status",
//       "meta",
//       "data"
//     ],
//     "type": "object"
//   }
// }
func CustomError(code int, messages ...string) *echo.HTTPError {
	var response Response
	response.Status.Code = uint(code)
	response.Status.MessageServer = http.StatusText(code)
	response.Status.MessageClient = http.StatusText(code)
	for index, value := range messages {
		if index == 0 {
			response.Status.MessageServer = value
		}
		if index == 1 {
			response.Status.MessageClient = value
		}
	}
	return echo.NewHTTPError(code, response)
}

func CustomErrorToken(code int, messages ...string) *echo.HTTPError {
	var response ResponseToken
	response.Status.Code = uint(code)
	response.Status.MessageServer = http.StatusText(code)
	response.Status.MessageClient = http.StatusText(code)
	for index, value := range messages {
		if index == 0 {
			response.Status.MessageServer = value
		}
		if index == 1 {
			response.Status.MessageClient = value
		}
	}
	return echo.NewHTTPError(code, response)
}
