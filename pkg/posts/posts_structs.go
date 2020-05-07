package posts

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Author автор поста или комментария
type Author struct {
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
}

// PostVote голоса отданные посту
type PostVote struct {
	User string `json:"user" bson:"user"`
	Vote int32  `json:"vote" bson:"vote"`
}

// PostComment комментарий
type PostComment struct {
	ID      string  `json:"id" bson:"id"`
	Created string  `json:"created" bson:"created"`
	Author  *Author `json:"author" bson:"author"`
	Body    string  `json:"body" bson:"body"`
}

// Post пост
type Post struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Score            int32              `json:"score" bson:"score"`
	Views            int32              `json:"views" bson:"views"`
	Type             string             `json:"type" bson:"type"`
	Title            string             `json:"title" bson:"title"`
	Text             string             `json:"text,omitempty" bson:"text,omitempty"`
	URL              string             `json:"url" bson:"url"`
	Category         string             `json:"category" bson:"category"`
	Created          string             `json:"created" bson:"created"`
	UpvotePercentage int32              `json:"upvotePercentage" bson:"upvotePercentage"`
	Author           *Author            `json:"author" bson:"author"`
	Votes            []*PostVote        `json:"votes" bson:"votes"`
	Comments         []*PostComment     `json:"comments" bson:"comments"`
}

// User структура для данных пользователя
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Admin    bool               `json:"admin" bson:"admin"`
}

// ResPost результаты, когда для поста не нужно возвращать сам пост
type ResPost struct {
	Message string
}

// PostFormPayload форма создания поста
type PostFormPayload struct {
	Category string
	Type     string
	Title    string
	Text     string
	URL      string
}
