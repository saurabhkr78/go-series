package dto

// Used when creating a user (POST /users)
type CreateUserInput struct {
	Name  string
	Email string
}

// Used when updating a user (PUT /users/{id})
type UpdateUserInput struct {
	Name  string
	Email string
}

// Used when partially updating a user (PATCH /users/{id})
type PatchUserInput struct {
	Name  *string
	Email *string
}
