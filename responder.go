package gallop

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"sync"
)

const RspBodyKey = "gallop_rsp_key"

var responderList []Responder
var once_resp_list sync.Once

func get_responder_list() []Responder {
	once_resp_list.Do(func() {
		responderList = []Responder{(StringResponder)(nil),
			(JsonResponder)(nil),
			(XMLResponder)(nil),
			(FileResponder)(nil),
		}
	})
	return responderList
}
func Convert(handler interface{}) gin.HandlerFunc {
	h_ref := reflect.ValueOf(handler)
	for _, r := range get_responder_list() {
		r_ref := reflect.TypeOf(r)
		if h_ref.Type().ConvertibleTo(r_ref) {
			return h_ref.Convert(r_ref).Interface().(Responder).RespondTo()
		}
	}
	return nil
}

type Responder interface {
	RespondTo() gin.HandlerFunc
}

type StringResponder func(*Context) string

func (s StringResponder) RespondTo() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := s(&Context{c})
		c.Set(RspBodyKey, body)
		c.String(200, body)
	}
}

type Json interface{}
type JsonResponder func(*Context) Json

func (j JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		body := j(&Context{context})
		context.Set(RspBodyKey, body)
		context.JSON(200, body)
	}
}

type XML interface{}

type XMLResponder func(*Context) XML

func (s XMLResponder) RespondTo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header()["content-type"] = []string{"application/xml; charset=utf-8"}
		body := s(&Context{c}).(string)
		c.Set(RspBodyKey, body)
		c.String(200, body)
	}
}

type File struct {
	Data        []byte
	ContentType string
	FileName    string
}

type FileResponder func(*Context) File

func (f FileResponder) RespondTo() gin.HandlerFunc {
	return func(c *gin.Context) {
		file := f(&Context{c})
		c.Header("Content-Disposition", "attachment; filename="+file.FileName)
		c.Header("Content-Transfer-Encoding", "binary")
		c.Data(200, file.ContentType, file.Data)
	}
}
