package main

import "github.com/ChangHChen/Reading-Copilot/webGateway/internal/validator"

type userSignupForm struct {
	UserName            string `form:"username"`
	Email               string `form:"email"`
	PWD                 string `form:"pwd"`
	PWDConfirm          string `form:"pwdconfirm"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	PWD                 string `form:"pwd"`
	validator.Validator `form:"-"`
}

type updatePWDForm struct {
	CurPWD              string `form:"curpwd"`
	NewPWD              string `form:"newpwd"`
	NewPWDConfirm       string `form:"newpwdconfirm"`
	validator.Validator `form:"-"`
}
