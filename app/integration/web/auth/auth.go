package auth

type AdminAuth struct {
	username string
	password string
}

func (a *AdminAuth) IsValid(username string, password string) bool {
	return a.username == username && a.password == password
}

func NewAdminAuth(username string, password string) AdminAuth {
	if username == "" || password == "" {
		panic("admin username or password cannot be empty")
	}

	return AdminAuth{
		username: username,
		password: password,
	}
}
