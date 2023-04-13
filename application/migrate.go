package main

import (
	"context"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
)

func StartDbMigrate(conn *dgo.Dgraph) {

	err := conn.Alter(context.Background(), &api.Operation{
		DropAll: isDropFirst(),
	})

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateAccountType,
	})
	if err != nil {
		log.Fatal("Alter CreateAccountType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCaptchaType,
	})
	if err != nil {
		log.Fatal("Alter CreateCaptchaType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCountryType,
	})
	if err != nil {
		log.Fatal("Alter CreateCountryType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreatePostType,
	})
	if err != nil {
		log.Fatal("Alter CreatePostType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCityType,
	})
	if err != nil {
		log.Fatal("Alter CreateCityType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateCommentType,
	})
	if err != nil {
		log.Fatal("Alter CreateCommentType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateTagType,
	})
	if err != nil {
		log.Fatal("Alter CreateTagType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateFriendshipType,
	})
	if err != nil {
		log.Fatal("Alter CreateFriendshipType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateLikeType,
	})
	if err != nil {
		log.Fatal("Alter CreateLikeType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateDialogType,
	})
	if err != nil {
		log.Fatal("Alter CreateDialogType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateMessageType,
	})
	if err != nil {
		log.Fatal("Alter CreateMessageType schemas has been closed with error")
	}

	err = conn.Alter(context.Background(), &api.Operation{
		Schema: CreateSettingsType,
	})
	if err != nil {
		log.Fatal("Alter CreateSettingsType schemas has been closed with error")
	}

	if isDropFirst() {
		txn := conn.NewTxn()
		marshalR := []byte(InsertCountryRu)
		marshalRB := []byte(InsertRB)
		_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: marshalR})
		_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: marshalRB, CommitNow: true})
		if err != nil {
			log.Fatal("Import new data has been closed with error")
		}
	}
}

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
	friends
	statusCode
	messagePermission
	createdOn
	updatedOn
	birthDate
	lastOnlineTime
    settings
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
posts: [uid] @reverse .
friends: [uid] @reverse .
settings: uid .
`

const CreatePostType = `type Post {
	authorId
	postText
	title
	likes
	comments
    isDeleted
    publishDate
    myLike
    commentsCount
	likeAmount
	time
	isBlocked
	tags
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
likes: [uid] @reverse .
tags: [uid] @reverse .
comments: [uid] .
`

const CreateCommentType = `type Comment {
	commentText
	authorId
	parentId
	postId
	likes
	commentType
	commentsCount
	myLike
	likeAmount
	timeChanged
	time
	comments
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

const CreateFriendshipType = `type Friendship {
    friend
	status
	previousStatus
}
friend: uid .
status: string @index(hash) .
previousStatus: string .
`

const CreateTagType = `type Tag {
    name
}
name: string @index(hash) .
`

const CreateLikeType = `type Like {
    authorId
}
authorId: string @index(hash) .
`

const CreateCaptchaType = `type Captcha {
    captchaId
	captchaCode
	expiredTime
}

captchaId: string @index(term) @lang .
captchaCode: string @index(term) @lang .
expiredTime: datetime .`

const CreateMessageType = `type Message {
	messageText    
	authorId
	recipientId
	timeSend
    isRead
	isDeleted
}
messageText: string @index(term) @lang .
authorId: string @index(hash) .
recipientId: string @index(hash).
timeSend: int .
isRead: bool .
isDeleted: bool .`

const CreateDialogType = `type Dialog {
	participantOne
	participantTwo
    unreadCount
    messages
	isDeleted
}
participantOne: uid @reverse .
participantTwo: uid @reverse .
unreadCount: int .
messages: [uid] .
isDeleted: bool .`

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

const CreateSettingsType = `type Settings {
	enablePost
    enablePostComment
	enableCommentComment
	enableMessage
	enableFriendRequest
	enableFriendBirthday
	enableSendEmailMessage
	enableIsDeleted
}

enablePost: bool .
enablePostComment: bool .
enableCommentComment: bool .
enableMessage: bool .
enableFriendRequest: bool .
enableFriendBirthday: bool .
enableSendEmailMessage: bool .
enableIsDeleted: bool .`

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
