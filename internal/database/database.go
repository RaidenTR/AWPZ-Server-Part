package database

import (
	"AWPZ/internal/authorizationdata"
	"AWPZ/internal/registration"
	"database/sql"
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"

	_ "github.com/go-sql-driver/mysql"
)

var dbInstance *sql.DB

func InitializeDB(DSN string) {
	var err error
	dbInstance, err = sql.Open("mysql", DSN)
	if err != nil {
		panic("DB wasn't opened")
	}
}

func CloseDB() {
	dbInstance.Close()
}

func StartLecture(groups string, idSub, idLector int) int64 {
	query := fmt.Sprintf("insert into Lecture (SubjectID,Date) "+
		"values (%d , NOW())", idSub)
	transaction, _ := dbInstance.Begin()
	res, _ := transaction.Exec(query)
	transaction.Commit()
	idLecture, _ := res.LastInsertId()
	groupsAsArr := strings.Split(groups, ",")
	for a, b := range groupsAsArr {
		groupsAsArr[a] = `'` + b + `'`
	}
	addStudentsInLecture(idLecture, strings.Join(groupsAsArr, ", "))
	return idLecture
}

func addStudentsInLecture(lectureID int64, groups string) {
	query := fmt.Sprintf("insert into Mark (StudentID,LectureID,Value, IsPresent) "+
		"select ID, %d ,0 ,0 "+
		"from Student "+
		"where Student.GroupName in (%s) "+
		"order by Student.GroupName", lectureID, groups)
	makeTransaction(query)
}

func PutMark(id, mark, lecture int) {
	query := fmt.Sprintf("update Mark set Value=%d where Mark.StudentID=%d and Mark.LectureID=%d", mark, id, lecture)
	makeTransaction(query)
}

func SetPresent(studentID, lectureID int) {
	query := fmt.Sprintf("update Mark set IsPresent=1 where Mark.StudentID=%d and Mark.LectureID=%d", studentID, lectureID)
	makeTransaction(query)
}

func makeTransaction(query string) {
	transaction, _ := dbInstance.Begin()
	transaction.Exec(query)
	transaction.Commit()
}

func IsAuthenticatedLector(loginData authorizationdata.Set) bool {
	var result bool
	err := dbInstance.QueryRow("select count(*) from Lector where Login=(?) and Password=(?)", loginData.Login, loginData.Password).Scan(&result)
	if err != nil {
		return false
	}
	return result
}

func IsAuthenticatedAdmin(loginData authorizationdata.Set) bool {
	var result bool
	err := dbInstance.QueryRow("select count(*) from Admin where Login=(?) and Password=(?)", loginData.Login, loginData.Password).Scan(&result)
	if err != nil {
		return false
	}
	return result
}

func IsAuthenticatedDevice(loginData authorizationdata.Set) bool {
	var result bool
	err := dbInstance.QueryRow("select count(*) from Device where MACAdress=(?) ", loginData.Login).Scan(&result)
	if err != nil {
		return false
	}
	return result
}

func GetLectorID(loginData authorizationdata.Set) int {
	var result int
	dbInstance.QueryRow("select ID from Lector where Login=(?) and Password=(?)", loginData.Login, loginData.Password).Scan(&result)
	return result
}

func GetAdminID(loginData authorizationdata.Set) int {
	var result int
	dbInstance.QueryRow("select ID from Admin where Login=(?) and Password=(?)", loginData.Login, loginData.Password).Scan(&result)
	return result
}

func GetDeviceID(loginData authorizationdata.Set) int {
	var result int
	dbInstance.QueryRow("select ID from Device where MACAdress=(?)", loginData.Login).Scan(&result)
	return result
}

func GetDeviceRoom(id int) string {
	var result string
	dbInstance.QueryRow("select Room from Device where ID=(?)", id).Scan(&result)
	return result
}

func GetStudentsList(lectureID int) []Student {
	rows := getStudentsOnLectureRows(lectureID)
	result := formStudentsList(rows)
	return result
}

func getStudentsOnLectureRows(lectureID int) *sql.Rows {
	rows, _ := dbInstance.Query(
		"select Student.Surname, Student.ID, Student.GroupName, Mark.Value, Mark.IsPresent "+
			"from Mark "+
			"inner join Student on Student.ID=Mark.StudentID "+
			"where Mark.LectureID=(?) "+
			"order by Student.GroupName", lectureID)
	return rows
}

func formStudentsList(rows *sql.Rows) []Student {
	resultSlice := make([]Student, 0)
	for rows.Next() {
		var s Student
		rows.Scan(&s.Surname, &s.ID, &s.Group, &s.Value, &s.IsPresent)
		resultSlice = append(resultSlice, s)
	}
	return resultSlice
}

func GenerateJSONForLecuteCourse(subjectID int) map[string][]Student {
	result := make(map[string][]Student)
	rowsLectures := getLecturesForOutput(subjectID)
	for rowsLectures.Next() {
		var date string
		var lectureID int
		rowsLectures.Scan(&lectureID, &date)
		result[date] = make([]Student, 0)
		rowsStudents := getStudentsForOutput(lectureID)
		for rowsStudents.Next() {
			var data Student
			rowsStudents.Scan(&data.Surname, &data.Group, &data.IsPresent, &data.Value)
			result[date] = append(result[date], data)
		}
	}
	return result
}

func GenerateXLSXForLectureCourse(subjectID int) *excelize.File {
	resultFile := formXLSXFile(subjectID)
	return resultFile
}

func formXLSXFile(subjectID int) *excelize.File {
	resultFile := excelize.NewFile()
	rowsLectures := getLecturesForOutput(subjectID)
	for rowsLectures.Next() {
		currentRowOfSheet := startFillXLXSFrom
		var date string
		var lectureID int
		rowsLectures.Scan(&lectureID, &date)
		formHeadOfFile(date, resultFile)
		rowsStudents := getStudentsForOutput(lectureID)
		for rowsStudents.Next() {
			var data Student
			rowsStudents.Scan(&data.Surname, &data.Group, &data.IsPresent, &data.Value)
			fillRow(date, currentRowOfSheet, data, resultFile)
			currentRowOfSheet++
		}
	}
	return resultFile
}

func getLecturesForOutput(subjectID int) *sql.Rows {
	result, _ := dbInstance.Query(
		"select ID, Date "+
			"from Lecture "+
			"where SubjectID=(?) "+
			"order by Date", subjectID)
	return result
}

func getStudentsForOutput(lectureID int) *sql.Rows {
	result, _ := dbInstance.Query(
		"select Student.Surname, Student.GroupName, Mark.IsPresent, Mark.Value "+
			"from Mark "+
			"inner join Student on Mark.StudentID = Student.ID "+
			"where Mark.LectureID = (?) "+
			"order by Mark.LectureID, Student.GroupName ", lectureID)
	return result
}

func AddStudent(form registration.StudentData) int64 {
	query := fmt.Sprintf("insert into Student (Name, Surname, GroupName) values ('%s', '%s', '%s')", form.Name, form.Surname, form.Group)
	fmt.Println(query)
	transaction, _ := dbInstance.Begin()
	res, _ := transaction.Exec(query)
	transaction.Commit()
	studentID, _ := res.LastInsertId()
	return studentID
}

func AddLector(form registration.LectorData) {
	query := fmt.Sprintf("insert into Lector (Name, Surname, Login, Password) values ('%s', '%s', '%s', '%s')",
		form.Name, form.Surname, form.Login, form.Password)
	makeTransaction(query)
}

func AddSubject(form registration.SubjectData) {
	query := fmt.Sprintf("insert into Subject (LectorID, Title) values ('%d', '%s')",
		form.LectorID, form.Title)
	makeTransaction(query)
}

func DeleteStudent(ID int) {
	query := fmt.Sprintf("delete from Student where ID='%d'", ID)
	makeTransaction(query)
}

func DeleteLector(ID int) {
	query := fmt.Sprintf("delete from Lector where ID='%d'", ID)
	makeTransaction(query)
}

func DeleteSubject(ID int) {
	query := fmt.Sprintf("delete from Subject where ID='%d'", ID)
	makeTransaction(query)
}

func AddDevice(form registration.DeviceData) {
	query := fmt.Sprintf("insert into Device (Room, MACAdress) values ('%s', '%s')",
		form.Room, form.MACAdress)
	fmt.Println(query)
	makeTransaction(query)
}

func DeleteDevice(ID int) {
	query := fmt.Sprintf("delete from Device where ID='%d'", ID)
	makeTransaction(query)
}

func GetLectorSubjects(LectorID int) []Subject {
	resultSubjectsStructs := make([]Subject, 0)
	resultRows, _ := dbInstance.Query("select ID, Title from Subject where LectorID=(?)", LectorID)
	for resultRows.Next() {
		var subject Subject
		resultRows.Scan(&subject.ID, &subject.Title)
		resultSubjectsStructs = append(resultSubjectsStructs, subject)
	}
	return resultSubjectsStructs
}

func GetGroups() []string {
	var result []string
	resultRows, _ := dbInstance.Query("select distinct GroupName from Student")
	for resultRows.Next() {
		var groupName string
		resultRows.Scan(&groupName)
		result = append(result, groupName)
	}
	return result
}
