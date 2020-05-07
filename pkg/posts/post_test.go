package posts

import (
	"fmt"
	"net/http/httptest"
	"testing"

	mdbabstractlayer "unit-tests-mongo-golang/pkg/mongodbabstractlayer"

	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestNewPostsRepo(t *testing.T) {

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)

	NewPostsRepo(testCollection)

}

func TestGetAll(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	ctx := gomock.Any()

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}

	testCursor := mdbabstractlayer.NewMockIMongoCursor(ctrl)

	filter := bson.M{}
	// Сортировка по рейтингу
	options := options.Find()
	options.SetSort(map[string]int{"score": -1})

	// Good
	testCollection.EXPECT().Find(ctx, filter, options).Return(testCursor, nil)

	testCursor.EXPECT().Next(ctx)
	testCursor.EXPECT().Close(ctx).Return(nil)

	_, err := service.GetAll(req)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// Error Find
	errorText := "FIND_DB_ERROR"

	testCollection.EXPECT().Find(ctx, filter, options).Return(testCursor, fmt.Errorf(errorText))

	_, err = service.GetAll(req)
	if err != nil && err.Error() != errorText {
		t.Errorf("unexpected err: %s", err)
		return
	}

}
