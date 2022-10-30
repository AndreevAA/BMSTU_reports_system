package label

type CreateLabelDTO struct {
	Name       string `json:"name" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type UpdateLabelDTO struct {
	Name       string `json:"name" db:"name"`
	Department string `json:"department" db:"department"`
}

type GetAllLabelsDTO struct {
	Labels []Label `json:"labels"`
}
