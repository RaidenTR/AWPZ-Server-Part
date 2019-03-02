package routes

import (
	"AWPZ/internal/authorization"
	"AWPZ/internal/authorizationdata"
	"AWPZ/internal/database"
	"AWPZ/internal/recognizer"
	"AWPZ/internal/registration"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setAdminRoutes(group *gin.RouterGroup) {
	newGroup := group.Group("/admin", adminMiddleware)
	newGroup.POST("/addStudent", addStudentHandler)
	newGroup.DELETE("/deleteStudent", deleteStudentHandler)
	newGroup.POST("/addLector", addLectorHandler)
	newGroup.DELETE("/deleteLector", deleteLectorHandler)
	newGroup.POST("/addSubject", addSubjectHandler)
	newGroup.DELETE("/deleteSubject", deleteSubjectHandler)
	newGroup.POST("/addDevice", addDeviceHandler)
	newGroup.DELETE("/deleteDevice", deleteDeviceHandler)
	newGroup.GET("/info/students", getStudentsToDelete)
	newGroup.GET("/info/lectors", getLectorsToDelete)
	newGroup.GET("/info/devices", getDevicesToDelete)
	newGroup.GET("/info/subjects", getSubjectsToDelete)
	newGroup.POST("/teachStudent", teachStudentFaceHandler)
}

func adminMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("JWT")
	if !authorization.IsAdmin(tokenString) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}

func getAdminToken(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	login, password := getInput(c)
	loginStruct := authorizationdata.Set{
		Login:     login,
		Password:  password,
		AccessLvl: authorization.Admin,
	}
	token, err := authorization.GetAdminToken(loginStruct)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"JWT": token})
}

func addStudentHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var student registration.StudentData
	if c.Bind(&student) == nil {
		studentId := database.AddStudent(student)
		photo := c.PostForm("photo")
		recognizer.Teach(photo, studentId)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func deleteStudentHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	idString, ok := c.GetPostForm("id")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err == nil {
		database.DeleteStudent(id)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func addLectorHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	///TODO fix binding
	var lector registration.LectorData
	if c.Bind(&lector) == nil {
		database.AddLector(lector)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func deleteLectorHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	idString, ok := c.GetPostForm("id")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err == nil {
		database.DeleteLector(id)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func addSubjectHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var subject registration.SubjectData
	if c.Bind(&subject) == nil {
		database.AddSubject(subject)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func deleteSubjectHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	idString, ok := c.GetPostForm("id")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err == nil {
		database.DeleteSubject(id)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func addDeviceHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	var device registration.DeviceData
	if c.Bind(&device) == nil {
		database.AddDevice(device)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func deleteDeviceHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	idString, ok := c.GetPostForm("id")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err == nil {
		database.DeleteDevice(id)
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func getStudentsToDelete(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := database.GetStudentsDataList()
	c.JSON(http.StatusOK, data)
}

func getLectorsToDelete(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := database.GetLectorsList()
	c.JSON(http.StatusOK, data)
}

func getSubjectsToDelete(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := database.GetSubjectsList()
	c.JSON(http.StatusOK, data)
}

func getDevicesToDelete(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := database.GetDevicesList()
	c.JSON(http.StatusOK, data)
}

func teachStudentFaceHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	id := c.PostForm("id")
	photo := c.PostForm("photo")
	idInt, _ := strconv.Atoi(id)
	recognizer.Teach(photo, int64(idInt))
	c.Status(http.StatusOK)
}
