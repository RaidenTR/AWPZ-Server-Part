package routes

type Lecture struct {
	Room      string `form:"room" json:"room" binding:"required"`
	Groups    string `form:"groups[]" json:"groups[]" binding:"required"`
	SubjectID int    `form:"subjectId" json:"subjectId" binding:"required"`
}
