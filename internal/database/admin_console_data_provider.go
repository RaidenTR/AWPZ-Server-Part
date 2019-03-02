package database

import (
	"database/sql"
)

type DataPair struct {
	ID       int64  `json:"id"`
	TextInfo string `json:"text"`
}

func GetStudentsDataList() []DataPair {
	resultRows, _ := dbInstance.Query("select ID, Surname from Student")
	return dbRowsToObjects(resultRows)
}

func dbRowsToObjects(rows *sql.Rows) []DataPair {
	result := make([]DataPair, 0)
	for rows.Next() {
		var data DataPair
		rows.Scan(&data.ID, &data.TextInfo)
		result = append(result, data)
	}
	return result
}

func GetLectorsList() []DataPair {
	resultRows, _ := dbInstance.Query("select ID, Surname from Lector")
	return dbRowsToObjects(resultRows)
}

func GetDevicesList() []DataPair {
	resultRows, _ := dbInstance.Query("select ID, Room from Device")
	return dbRowsToObjects(resultRows)
}

func GetSubjectsList() []DataPair {
	resultRows, _ := dbInstance.Query("select ID, Title from Subject")
	return dbRowsToObjects(resultRows)
}
