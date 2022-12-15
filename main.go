package main

import (
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/samandar2605/docker_exam/storage"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CreateUser struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

func main() {

	psqlUrl := "host=postgres port=5432 user=postgres password=1 dbname=users sslmode=disable"
	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	strg := storage.NewDBManager(psqlConn)

	router := gin.Default()

	router.POST("/users/", func(ctx *gin.Context) {
		var (
			req CreateUser
		)
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		resp, err := strg.Create(&storage.User{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			PhoneNumber: req.PhoneNumber,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	})

	router.GET("users/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		resp, err := strg.Get(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, storage.User{
			Id:          resp.Id,
			FirstName:   resp.FirstName,
			LastName:    resp.LastName,
			PhoneNumber: resp.PhoneNumber,
		})
	})
	err = router.Run(":8000")
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
