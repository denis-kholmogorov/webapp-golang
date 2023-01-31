package main

const CreateAccountType = `type Account {
    email
    firstName
    lastName
    password
    age
}

email: string @index(term) @lang .
firstName: string @index(term) @lang .
lastName: string @index(term) @lang .
password: string @index(term) @lang .
age: int @index(int) .`

const CreateCaptchaType = `type getCaptcha {
    captchaId
	captchaCode
	expiredTime
}

captchaId: string @index(term) @lang .
captchaCode: string @index(term) @lang .
expiredTime: datetime .`

const FindByEmail = `{
	q (func: eq(email, "%s")) {
		email
	}
}`
