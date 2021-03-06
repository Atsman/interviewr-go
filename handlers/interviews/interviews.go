package interviews

import (
	"net/http"

	"github.com/atsman/interviewr-go/db/interviewdb"
	"github.com/atsman/interviewr-go/handlers/utils"
	"github.com/atsman/interviewr-go/models"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("handlers.interviews")

func bindInterview(c *gin.Context) (error, *models.Interview) {
	interview := models.Interview{}
	err := c.Bind(&interview)
	return err, &interview
}

func notValidModel(err error) *utils.ApiError {
	return &utils.ApiError{
		Status:      http.StatusBadRequest,
		Title:       "Interview model not valid",
		Description: err.Error(),
	}
}

func Create(c *gin.Context) {
	err, interview := bindInterview(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, notValidModel(err))
		return
	}

	db := utils.GetDb(c)
	userId := utils.GetUserId(c)
	err = interviewdb.Create(db, userId, interview)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, interview)
}

func Update(c *gin.Context) {
	updateModel := models.InterviewUpdateModel{}
	err := c.BindJSON(&updateModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, notValidModel(err))
		return
	}

	db := utils.GetDb(c)
	id := c.Params.ByName("id")
	userID := utils.GetUserId(c)
	err, updatedInterview := interviewdb.Update(db, userID, id, &updateModel)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedInterview)
}

func Delete(c *gin.Context) {
	db := utils.GetDb(c)
	userID := utils.GetUserId(c)
	id := c.Params.ByName("id")
	err, interview := interviewdb.Delete(db, userID, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, interview)
}

func GetOne(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")
	err, interview := interviewdb.GetOne(db, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, interview)
}

func GetList(c *gin.Context) {
	db := utils.GetDb(c)
	err, query := BuildQuery(c)
	if err != nil {
		c.Error(err)
		return
	}
	log.Debug("interviews - GetList, query: ", query)
	err, interviews := interviewdb.GetList(db, query)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, interviews)
}

func CreateFeedback(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	feedback := models.Feedback{}
	err := c.BindJSON(&feedback)
	if err != nil {
		c.Error(err)
		return
	}

	err = interviewdb.CreateFeedback(db, id, &feedback)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, feedback)
}

func GetFeedback(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	err, feedback := interviewdb.GetFeedback(db, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, feedback)
}

func Start(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	err := interviewdb.Start(db, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func End(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	err := interviewdb.End(db, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
