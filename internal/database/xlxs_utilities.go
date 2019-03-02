package database

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const startFillXLXSFrom = 2

func formHeadOfFile(sheet string, file *excelize.File) {
	file.NewSheet(sheet)
	file.SetCellValue(sheet, "A1", "Surname")
	file.SetCellValue(sheet, "B1", "Group")
	file.SetCellValue(sheet, "C1", "IsPresent")
	file.SetCellValue(sheet, "D1", "Mark")
}

func fillRow(sheet string, row int, info Student, file *excelize.File) {
	file.SetCellValue(sheet, "A"+strconv.Itoa(row), info.Surname)
	file.SetCellValue(sheet, "B"+strconv.Itoa(row), info.Group)
	file.SetCellValue(sheet, "C"+strconv.Itoa(row), info.IsPresent)
	file.SetCellValue(sheet, "D"+strconv.Itoa(row), info.Value)
}
