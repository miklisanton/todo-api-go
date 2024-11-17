package requests

type PutTaskRequest struct {
	Title       *string `json:"title" validate:"required"`
	Description *string `json:"description" validate:"required"`
	DueDate     *string `json:"due_date" validate:"required"`
}

type PostTaskRequest struct {
	Title       *string `json:"title" validate:"required"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
}

type PatchTaskRequest struct {
	Completed bool `json:"completed" validate:"required"`
}
