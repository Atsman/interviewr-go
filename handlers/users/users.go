package users

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/atsman/interviewr-go/db/companydb"
	"github.com/atsman/interviewr-go/db/userdb"
	"github.com/atsman/interviewr-go/handlers/utils"
	"github.com/atsman/interviewr-go/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"gopkg.in/mgo.v2/bson"
)

var log = logging.MustGetLogger("handlers.users")

var userNotFoundError = utils.ApiError{
	Status: http.StatusNotFound,
	Title:  "User not found",
}

func getUser(c *gin.Context) (error, *models.User) {
	user := models.User{}
	err := c.Bind(&user)
	return err, &user
}

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateJwtToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["_id"] = userId
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	return tokenString, err
}

func Login(c *gin.Context) {
	loginData := LoginData{}
	err := c.Bind(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	db := utils.GetDb(c)
	err, user := userdb.GetOne(db, bson.M{
		"username": loginData.Username,
	})

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "The username or password don't match",
		})
		return
	}

	token, err := CreateJwtToken(user.ID.Hex())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"id_token": token,
	})
}

func Create(c *gin.Context) {
	db := utils.GetDb(c)

	err, user := getUser(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = userdb.Create(db, user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func Update(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	log.Debugf("Update, Id=%v", id)

	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.Error(err)
		return
	}

	user.ID = bson.ObjectId("")

	err, updatedUser := userdb.Update(db, id, &user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func Delete(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	err, user := userdb.Delete(db, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetList(c *gin.Context) {
	db := utils.GetDb(c)

	err, users := userdb.GetList(db, &bson.M{})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetOne(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")
	if id == "me" {
		id = utils.GetUserId(c)
	}

	log.Debug("Userid :%v", id)
	err, user := userdb.GetOneById(db, id)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, userNotFoundError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserCompanies(c *gin.Context) {
	db := utils.GetDb(c)
	id := c.Params.ByName("id")

	err, companies := companydb.GetList(db, &bson.M{
		"owner": bson.ObjectIdHex(id),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, userNotFoundError)
		return
	}

	c.JSON(http.StatusOK, companies)
}
