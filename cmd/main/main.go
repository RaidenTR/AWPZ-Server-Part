package main

import (
	"AWPZ/internal/database"
	"AWPZ/internal/recognizer"
	"AWPZ/internal/registration"
	"AWPZ/internal/routes"
	"fmt"
	"reflect"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	testType := registration.StudentData{}
	testType.Name = "test"
	type1 := reflect.TypeOf(testType)
	field1, _ := type1.FieldByName("Name")
	fmt.Println(field1.Tag)
	database.InitializeDB("admin:admin@/AWPZDB")
	recognizer.InitializeRecognizor()
	defer database.CloseDB()
	r := routes.CreateRoutes()
	//recognizer.Teach("", 5)
	r.Run(":8030")
	//r.RunTLS(":8030", "certificate.pem", "key.pem")
}
