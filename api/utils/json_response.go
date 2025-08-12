package utils

import (
	"fmt"

	"clevergo.tech/jsend"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/gofiber/contrib/websocket"
	fiber "github.com/gofiber/fiber/v2"
)

// JSONSuccess is a generic success output writer
func JSONSuccess(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(jsend.New(data))
}

// JSONFail is a generic fail output writer
// JSONFail can used for 4xx status code response
func JSONFail(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(jsend.NewFail(data))
}

// JSONError is a generic error output writer
// JSONError can used for 5xx status code response
func JSONError(c *fiber.Ctx, statusCode int, err string) error {
	return c.Status(statusCode).JSON(jsend.NewError(err, statusCode, nil))
}

// WsJSONSuccess is a generic success output writer
func JSONSuccessWs(c *websocket.Conn, eventName string, data interface{}) error {
	if c == nil {
		return nil
	} else {
		return c.WriteJSON(jsend.New(structs.SocketResponseFormat{EventName: eventName, Data: data}))
	}
}

// WsJSONFail is a generic fail output writer
// WsJSONFail can used for 4xx status code response
func JSONFailWs(c *websocket.Conn, eventName string, data interface{}) error {
	if c == nil {
		return nil
	} else {
		return c.WriteJSON(jsend.NewFail(structs.SocketResponseFormat{EventName: eventName, Data: data}))
	}
}

// JsonErrorWs is a generic error output writer
// JsonErrorWs can used for 5xx status code response
func JSONErrorWs(c *websocket.Conn, eventName string, data interface{}) error {
	if c == nil {
		return nil
	} else {
		return c.WriteJSON(jsend.NewError("Error", -1, structs.SocketResponseFormat{EventName: eventName, Data: data}))
	}

}

func JSONSuccessPdf(c *fiber.Ctx, statusCode int, pdfBytes []byte) error {
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, constants.PdfName))
	return c.Status(statusCode).Send(pdfBytes)
}
