package report

type CreateReportDTO struct {
	Header string `json:"header" binding:"required"`
	Body   string `json:"body"`
}

type UpdateReportDTO struct {
	ID     int
	Header string `json:"header"`
	Body   string `json:"body"`
}

type GetAllReportsDTO struct {
	Reports []Report `json:"reports"`
}
