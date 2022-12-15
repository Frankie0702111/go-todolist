package request

// Create login request struct when user login from /login URL
type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email,max=50"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
}

// Create register request struct when user register from /register URL
type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}
