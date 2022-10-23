package label

type Label struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name" binding:"required"`
	Department string `json:"department" db:"department" binding:"required"`
}
