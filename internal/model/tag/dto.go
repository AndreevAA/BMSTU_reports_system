package tag

type CreateTagDTO struct {
	Name  string `json:"name" binding:"required"`
	Department string `json:"department" binding:"required"`
}

type UpdateTagDTO struct {
	Name  string `json:"name" db:"name"`
	Department string `json:"department" db:"department"`
}

type GetAllTagsDTO struct {
	Tags []Tag `json:"tags"`
}
