package main

const CreateAccountType = `type Account {
    email
    firstName
    lastName
    password
    age
	isDeleted
	isBlocked
	isOnline
	phone
	photo
	photoId
	photoName
	about
	city
	country
	posts
	statusCode
	messagePermission
	createdOn
	updatedOn
	birthDate
	lastOnlineTime
}

email: string @index(term) @lang .
firstName: string @index(term) @lang .
lastName: string @index(term) @lang .
password: string @index(term) @lang .
age: int @index(int) .
isDeleted: bool .
isBlocked: bool .
isOnline: bool .
phone: string .
photo: string .
photoId: string .
photoName: string .
about: string .
city: string .
country: string .
statusCode: string .
messagePermission: string .
createdOn: datetime .
updatedOn: datetime .
birthDate: datetime .
lastOnlineTime: datetime .
posts: [uid] .`

const CreatePostType = `type Post {
    isDeleted
    publishDate
    myLike
    commentsCount
	likeAmount
	time
	isBlocked
	type
}

isDeleted: bool .
publishDate: datetime .
myLike: bool .
commentsCount: int .
likeAmount: int .
time: datetime .
isBlocked: bool .
type: string .
likes: [uid] .
tags: [uid] .
`

const CreateCaptchaType = `type Captcha {
    captchaId
	captchaCode
	expiredTime
}

captchaId: string @index(term) @lang .
captchaCode: string @index(term) @lang .
expiredTime: datetime .`

const CreateCityType = `type City {
    title
}

title: string @lang .`

const CreateCountryType = `type Country {
    title
	cities
}

title: string @lang .
cities: [uid] .`

const InsertCountryRu = `{
"title":"Россия",
"cities":[
{"title":"Москва","dgraph.type":["City"]},
{"title":"Ярославль","dgraph.type":["City"]},
{"title":"Саратов","dgraph.type":["City"]},
{"title":"Санкт-Петербург","dgraph.type":["City"]},
{"title":"Казань","dgraph.type":["City"]}
],
"dgraph.type":["Country"]
}`

const InsertRB = `{
"title":"Беларуссия",
"cities":[
{"title":"Минск","dgraph.type":["City"]},
{"title":"Гродно","dgraph.type":["City"]},
{"title":"Сморгонь","dgraph.type":["City"]},
{"title":"Витебск","dgraph.type":["City"]},
{"title":"Орша","dgraph.type":["City"]}
],
"dgraph.type":["Country"]
}`

const FindByEmail = `{
	q (func: eq(email, "%s")) {
		email
	}
}`
