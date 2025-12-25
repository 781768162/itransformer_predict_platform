package schemas

type LoginRequest struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	ExpireTime int    `json:"expire_time"`
}

type RegisterRequest struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type UserInfoRequest struct {
}

type UserInfoResponse struct {
	Username string `json:"user_name"`
}
