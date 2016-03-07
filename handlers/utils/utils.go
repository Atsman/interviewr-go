package utils

import (
	"net/url"

	"github.com/atsman/interviewr-go/commons/strutils"
	"github.com/atsman/interviewr-go/middlewares"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetDb(c *gin.Context) *mgo.Database {
	db := c.MustGet("db").(*mgo.Database)
	return db
}

func GetUserId(c *gin.Context) string {
	return c.MustGet(middlewares.USER_ID).(string)
}

func AddIfExist(propName string, values *url.Values, query map[string]interface{}) {
	value := values.Get(propName)
	if !strutils.IsEmpty(value) {
		query[propName] = value
	}
}

type valueProcessor func(interface{}) interface{}

func ProcessAndAddIfExist(propName string, values *url.Values, query map[string]interface{}, pr valueProcessor) {
	value := values.Get(propName)
	if !strutils.IsEmpty(value) {
		query[propName] = pr(value)
	}
}

func ConvertToObjectId(value interface{}) interface{} {
	strVal := value.(string)
	hexVal := bson.ObjectIdHex(strVal)
	return hexVal
}
