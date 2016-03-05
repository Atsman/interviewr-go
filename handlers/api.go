package handlers

import (
	"github.com/atsman/interviewr-go/handlers/companies"
	"github.com/atsman/interviewr-go/handlers/users"
	"github.com/atsman/interviewr-go/middlewares"
	"github.com/gin-gonic/gin"
)

func NewEngine() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Connect())
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.Cors())

	r1 := r.Group("/api/v1")
	{
		//r1.POST("/login", login)
		r1.POST("/users", users.Create)
		//r1.GET("/images/:id", getImage)
	}

	authR := r1.Group("/")
	authR.Use(middlewares.Auth("secret"))
	{
		//authR.POST("/images", uploadImage)

		authR.GET("/users", users.List)
		authR.PUT("/users/:id", users.Update)
		authR.DELETE("/users/:id", users.Delete)
		//authR.GET("/users/:id/companies", getUserCompanies)

		//authR.GET("/companies", getCompaniesList)
		authR.POST("/companies", companies.Create)
		//authR.GET("/companies/:id", getCompany)
		authR.PUT("/companies/:id", companies.Update)
		authR.DELETE("/companies/:id", companies.Delete)

		/*authR.GET("/companies/:id/comments", getCompanyComments)
		authR.POST("/companies/:id/comments", createCompanyComment)

		authR.GET("/vacancies", getVacanciesList)
		authR.POST("/vacancies", createVacancy)
		authR.GET("/vacancies/:id", getVacancy)
		authR.PUT("/vacancies/:id", updateVacancy)
		authR.DELETE("/vacancies/:id", deleteVacancy)
		authR.GET("/vacancies/:id/subscription", getVacancySubscriptions)*/
	}
	return r
}
