package students

import (
	"errors"
	"fmt"
	"net/http"

	"golang/utils"

	"github.com/gin-gonic/gin"
)

var emailChannel = make(chan string)

func Init(server *gin.Engine) {

	/*
		Listen for messages from `emailChannel`
	*/
	go func() {
		for {
			fmt.Println("Msg from channel: ", <-emailChannel)
		}
	}()

	/*
		Student CRUD routes
	*/
	studentRouter := server.Group("/students")
	{
		studentRouter.POST("/", addStudent)
		studentRouter.GET("/", getStudents)
		studentRouter.DELETE("/:id", deleteStudents)
	}
}

type Address struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	City    string `json:"city"`
	PinCode string `json:"pinCode"`
}

type Student struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Email    string  `json:"email"`
	Address  Address `json:"address"`
	Degree   string  `json:"degree"`
	IsAlumni bool    `json:"isAlumni"`
}

func (st *Student) GetEmail() string {
	return st.Email
}

var studentsList = []Student{
	{Id: "1", Name: "Test 1", Age: 20, Email: "test1@biofourmis.com", Address: Address{Line1: "line 1", Line2: "line 2", City: "City", PinCode: "100000"}, Degree: "B.E", IsAlumni: false},
	{Id: "2", Name: "Test 2", Age: 30, Email: "test2@biofourmis.com", Address: Address{Line1: "line 1", Line2: "line 2", City: "City", PinCode: "100000"}, Degree: "B.E", IsAlumni: true},
}

func addStudent(con *gin.Context) {
	var st Student
	if err := con.ShouldBindJSON(&st); err != nil {
		utils.SendErrorResponse(con, http.StatusBadRequest, err)
		return
	}

	fmt.Println("St: Email:", st.GetEmail())

	emailChannel <- st.GetEmail()

	studentsList = append(studentsList, st)

	utils.SendSuccessResponse(con, http.StatusCreated, st)
}

func getStudents(con *gin.Context) {
	var search string
	if name, ok := con.GetQuery("name"); ok {
		search = name
	} else {
		utils.SendSuccessResponse(con, http.StatusOK, studentsList)
		return
	}

	var result Student
	var found bool
	for _, v := range studentsList {
		if search == v.Name {
			found = true
			result = v
			break
		}
	}

	if !found {
		utils.SendErrorResponse(con, http.StatusNotFound, errors.New("Not found"))
		return
	}

	utils.SendSuccessResponse(con, http.StatusOK, result)
}

func deleteStudents(con *gin.Context) {
	id, ok := con.Params.Get("id")
	if !ok {
		utils.SendErrorResponse(con, http.StatusBadRequest, errors.New("`id` is required"))
		return
	}

	var found bool
	for k, v := range studentsList {
		if id == v.Id {
			found = true
			studentsList = append(studentsList[:k], studentsList[k+1:]...)
			break
		}
	}

	if !found {
		utils.SendErrorResponse(con, http.StatusBadRequest, errors.New("`id` not found"))
		return
	}

	utils.SendSuccessResponse(con, http.StatusNoContent, gin.H{})
}
