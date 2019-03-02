package routes

import (
	"AWPZ/internal/authorization"
	"AWPZ/internal/authorizationdata"
	"AWPZ/internal/database"
	"AWPZ/internal/device"
	"AWPZ/internal/recognizer"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateRoutes() *gin.Engine {
	result := gin.Default()
	result.POST("/adminLogin", getAdminToken)
	result.POST("/login", getToken)
	result.POST("/loginDevice", getDeviceToken)
	result = corsPermitions(result)

	authorized := result.Group("/api", authMiddleware)
	authorized.PUT("/putMark", putMarkHandler)
	authorized.POST("/startLecture", startLectureHandler)
	authorized.GET("/getStudents", getStudentsHandler)
	authorized.GET("/getSubjectTables", getSubjectTables)
	authorized.POST("/recognize", recognizingHandler)
	authorized.POST("/recognizeAuditory", recognizingAuditoryHandler)

	setAdminRoutes(authorized)
	return result
}

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("JWT")
	if tokenString == "" || !authorization.ValidateToken(tokenString) {
		c.AbortWithStatus(http.StatusNonAuthoritativeInfo)
		return
	}
	c.Next()
}

//GetToken enpoint for /getToken
func getToken(c *gin.Context) {
	login, password := getInput(c)
	c.Header("Access-Control-Allow-Origin", "*")
	loginStruct := authorizationdata.Set{
		Login:     login,
		Password:  password,
		AccessLvl: authorization.Lector,
	}
	token, err := authorization.GetLectorToken(loginStruct)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id := authorization.GetIDFromToken(token)
	subjects := database.GetLectorSubjects(id)
	groups := database.GetGroups()
	rooms := database.GetDevicesList()
	c.JSON(http.StatusOK, gin.H{"JWT": token,
		"subjects": subjects,
		"groups":   groups,
		"rooms":    rooms})
}

func getDeviceToken(c *gin.Context) {
	macAdress, ok := c.GetPostForm("macAdress")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := authorization.GetDeviceToken(macAdress)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	room := authorization.GetDeviceRoom(token)
	c.JSON(http.StatusOK, gin.H{"JWT": token, "room": room})
}

func getInput(c *gin.Context) (string, string) {
	login, loginOK := c.GetPostForm("login")
	password, passwordOK := c.GetPostForm("password")
	if !(loginOK && passwordOK) {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return login, password
}

func putMarkHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id, err := strconv.Atoi(c.PostForm("id"))
	mark, err2 := strconv.Atoi(c.PostForm("mark"))
	lecture, err3 := strconv.Atoi(c.PostForm("lecture"))
	database.SetPresent(id, lecture)
	if err != nil || err2 != nil || err3 != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	database.PutMark(id, mark, lecture)
	c.Status(http.StatusOK)
}

func adminAddDBAction(c *gin.Context, form *interface{}, action func(interface{})) {
	if c.Bind(form) == nil {
		action(*form)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func startLectureHandler(c *gin.Context) {
	tokenString := c.GetHeader("JWT")
	lectorID := authorization.GetIDFromToken(tokenString)
	c.Header("Access-Control-Allow-Origin", "*")
	var lecture Lecture
	if c.Bind(&lecture) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	lectureId := database.StartLecture(lecture.Groups, lecture.SubjectID, lectorID)
	activateDevice(lecture.Room, int(lectureId))
	c.JSON(http.StatusOK, gin.H{"lectureId": lectureId})
}

func activateDevice(room string, lectureID int) {
	device.GetInstance().StartWatching(room, strconv.Itoa(lectureID))
	go func() {
		time.Sleep(120 * time.Minute)
		device.GetInstance().FinishWatching(room, strconv.Itoa(lectureID))
	}()
}

func getStudentsHandler(c *gin.Context) {
	idString := c.Query("lecture")
	c.Header("Access-Control-Allow-Origin", "*")
	id, _ := strconv.Atoi(idString)
	result := database.GetStudentsList(id)
	c.JSON(http.StatusOK, result)
}

func recognizingHandler(c *gin.Context) {
	room := authorization.GetDeviceRoom(c.GetHeader("JWT"))
	if !device.GetInstance().IsStreaming(room) {
		//return
	}
	photoBase64, _ := c.GetPostForm("photo")
	lectureIdString := device.GetInstance().GetLectureId(room)
	lectureId, _ := strconv.Atoi(lectureIdString)
	studentsID := recognizer.RecognizeBase64(photoBase64)
	fmt.Println(studentsID)
	for _, id := range studentsID {
		database.SetPresent(id, lectureId)
	}
	isSuccess := len(studentsID) > 0 && studentsID[0] != 0
	fmt.Println(isSuccess)
	file, _ := os.Create("photoLog/" + strconv.Itoa(rand.Int()) + " " + strconv.FormatBool(isSuccess))
	notBasePhoto, _ := base64.StdEncoding.DecodeString(photoBase64)
	file.Write(notBasePhoto)
	c.JSON(http.StatusOK, gin.H{"succes": isSuccess})
}

func recognizingAuditoryHandler(c *gin.Context) {
	lectureId, _ := strconv.Atoi(c.PostForm("lectureId"))
	photo := c.PostForm("photo")
	studentsID := recognizer.RecognizeBase64(photo)
	for _, id := range studentsID {
		database.SetPresent(id, lectureId)
	}
	c.Status(http.StatusOK)
}

///TODO delete this
func getSubjectTables(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	subjectIdStr := c.Query("subject")
	subjectId, _ := strconv.Atoi(subjectIdStr)
	result := database.GenerateJSONForLecuteCourse(subjectId)
	c.JSON(http.StatusOK, result)
}
