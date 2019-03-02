package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func corsPermitions(e *gin.Engine) *gin.Engine {
	e.OPTIONS("/api/startLecture", letRequestPermission)
	e.OPTIONS("/api/getStudents", letRequestPermission)
	e.OPTIONS("/api/putMark", letRequestPermission)
	e.OPTIONS("/api/getSubjectTables", letRequestPermission)
	e.OPTIONS("/api/admin/addStudent", letRequestPermission)
	e.OPTIONS("/api/admin/addLector", letRequestPermission)
	e.OPTIONS("/api/admin/deleteStudent", letRequestPermission)
	e.OPTIONS("/api/admin/deleteLector", letRequestPermission)
	e.OPTIONS("/api/admin/addDevice", letRequestPermission)
	e.OPTIONS("/api/admin/deleteDevice", letRequestPermission)
	e.OPTIONS("/api/admin/addSubject", letRequestPermission)
	e.OPTIONS("/api/admin/deleteSubject", letRequestPermission)
	e.OPTIONS("/api/admin/info/students", letRequestPermission)
	e.OPTIONS("/api/admin/info/lectors", letRequestPermission)
	e.OPTIONS("/api/admin/info/devices", letRequestPermission)
	e.OPTIONS("/api/admin/info/subjects", letRequestPermission)
	e.OPTIONS("/api/admin/teachStudent", letRequestPermission)
	return e
}

func letRequestPermission(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "AuthID,Origin, X-Requested-With, Content-Type, Accept, JWT")
	c.Header("Access-Control-Allow-Methods", "PUT, DELETE")
	c.Status(http.StatusOK)
}
