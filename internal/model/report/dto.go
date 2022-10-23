package report

type CreateReportDTO struct {
	Header string `json:"header" binding:"required"`
	Department  string `json:"department" binding:"required"`
	Body   string `json:"body"`
}

type UpdateReportDTO struct {
	ID     int
	Header string `json:"header"`
	Body   string `json:"body"`
	Department  string `json:"department"`
}

type GetAllReportsDTO struct {
	Reports []Report `json:"reports"`
}
