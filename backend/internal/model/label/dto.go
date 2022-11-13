package label

type CreateLabelDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateLabelDTO struct {
	Name string `json:"name" db:"name"`
}

type GetAllLabelsDTO struct {
	Labels []Label `json:"labels"`
}
