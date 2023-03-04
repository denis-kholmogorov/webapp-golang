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
firstName: string @index(trigram) @lang .
lastName: string @index(trigram) @lang .
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
birthDate: datetime @index(day).
lastOnlineTime: datetime .
posts: [uid] @reverse .`

const CreatePostType = `type Post {
	authorId
	postText
	title
    isDeleted
    publishDate
    myLike
    commentsCount
	likeAmount
	time
	isBlocked
	type
}
authorId: string @index(hash) .
postText: string @index(fulltext) .
title: string @index(fulltext) .
isDeleted: bool .
publishDate: datetime .
myLike: bool .
commentsCount: int .
likeAmount: int .
time: datetime .
isBlocked: bool .
type: string .
likes: [uid] .
tagged: [uid] @reverse .
comments: [uid] .
`
const CreateCommentType = `type Comment {
	commentText
	authorId
	parentId
	postId
	commentType
	commentsCount
	myLike
	likeAmount
	timeChanged
	time
	imagePath
	isBlocked
	isDeleted
}
commentText: string @index(fulltext) .
authorId: string @index(hash) .
parentId: string @index(hash) .
postId: string @index(hash) .
commentType: string @index(hash) .
commentsCount: int .
myLike: bool .
likeAmount: int .
likes: [uid] .
timeChanged: datetime .
time: datetime .
imagePath: string .
isDeleted: bool .
isBlocked: bool .
comments: [uid] .
`

const CreateTagType = `type Tag {
    name
}
name: string @index(hash) .
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
    cityTitle
}

cityTitle: string @lang .`

const CreateCountryType = `type Country {
    countryTitle
	cities
}

countryTitle: string @lang .
cities: [uid] .`

const InsertCountryRu = `{
"countryTitle":"Россия",
"cities":[
{"cityTitle":"Москва","dgraph.type":["City"]},
{"cityTitle":"Ярославль","dgraph.type":["City"]},
{"cityTitle":"Саратов","dgraph.type":["City"]},
{"cityTitle":"Санкт-Петербург","dgraph.type":["City"]},
{"cityTitle":"Казань","dgraph.type":["City"]}
],
"dgraph.type":["Country"]
}`

const InsertRB = `{
"countryTitle":"Беларуссия",
"cities":[
{"cityTitle":"Минск","dgraph.type":["City"]},
{"cityTitle":"Гродно","dgraph.type":["City"]},
{"cityTitle":"Сморгонь","dgraph.type":["City"]},
{"cityTitle":"Витебск","dgraph.type":["City"]},
{"cityTitle":"Орша","dgraph.type":["City"]}
],
"dgraph.type":["Country"]
}`

const FindByEmail = `{
	q (func: eq(email, "%s")) {
		email
	}
}`
