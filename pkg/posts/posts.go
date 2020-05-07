package posts

import (
	"math"
	"net/http"

	mdbabstractlayer "unit-tests-mongo-golang/pkg/mongodbabstractlayer"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostsRepo репозиторий для работы с пользователями системы
type PostsRepo struct {
	Collection mdbabstractlayer.IMongoCollection
}

// NewPostsRepo инициализация репозитория Posts
func NewPostsRepo(collection mdbabstractlayer.IMongoCollection) *PostsRepo {
	return &PostsRepo{
		Collection: collection,
	}
}

// GetAll получить данные всех постов из БД
func (repo *PostsRepo) GetAll(r *http.Request) ([]*Post, error) {
	// Сортировка по рейтингу
	options := options.Find()
	options.SetSort(map[string]int{"score": -1})

	items := []*Post{}

	cursor, err := repo.Collection.Find(r.Context(), bson.M{}, options)
	if err != nil {
		return nil, err
	}

	for cursor.Next(r.Context()) {
		elem := &Post{}
		err := cursor.Decode(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, elem)
	}

	// IDEA: Курсор требуется закрывать, согласно документации https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
	cursor.Close(r.Context())

	return items, nil
}

// CalcPostScoreData пересчитать процент положительных голосов и общий рейтинг поста
func (repo *PostsRepo) CalcPostScoreData(objID primitive.ObjectID, r *http.Request) (*Post, error) {

	updatedDocument := &Post{}

	postFilter :=
		primitive.E{
			Key:   "_id",
			Value: objID,
		}
	filter := bson.D{
		postFilter,
	}

	// Выполним поиск документа
	err := repo.Collection.FindOne(r.Context(), filter).Decode(updatedDocument)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("INSIDE PKG updatedDocument= %#v \n", updatedDocument.ID.Hex())

	// Пересчитать процент положительных голосов и общий рейтинг поста
	var upVotesLength int32 = 0
	var downVotesLength int32 = 0
	for i := 0; i < len(updatedDocument.Votes); i++ {
		if updatedDocument.Votes[i].Vote > 0 {
			upVotesLength++
		} else {
			downVotesLength++
		}
	}
	var allVotesLength int = 1
	if len(updatedDocument.Votes) != 0 {
		allVotesLength = len(updatedDocument.Votes)
	}
	var upvotePercentage float64 = math.Floor(
		float64(
			(float64(upVotesLength) / float64(allVotesLength)) * 100,
		),
	)
	var upvotePercentageINT32 int32 = int32(upvotePercentage)
	upvotePercentageData := primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{
				Key:   "upvotePercentage",
				Value: upvotePercentageINT32,
			},
		},
	}
	score := upVotesLength - downVotesLength
	scoreData := primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{
				Key:   "score",
				Value: score,
			},
		},
	}

	update := bson.D{
		// SCORE
		scoreData,
		// UPVOTE_PERCENTAGE
		upvotePercentageData,
	}

	// Обновить процент положительных голосов и общий рейтинг поста
	err = repo.Collection.FindOneAndUpdate(r.Context(), filter, update).Decode(updatedDocument)
	if err != nil {
		return nil, err
	}

	// Т.к. FindOneAndUpdate, не возвращает обновления, внесенные в документ
	// Выполним повторный поиск обновленного документа
	err = repo.Collection.FindOne(r.Context(), filter).Decode(updatedDocument)
	if err != nil {
		return nil, err
	}

	return updatedDocument, nil
}

