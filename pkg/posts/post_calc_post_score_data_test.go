package posts

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	mdbabstractlayer "unit-tests-mongo-golang/pkg/mongodbabstractlayer"

	monkey "bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestCalcPostScoreData(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	reqCtx := req.Context()
	ctx := gomock.Any()

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}
	testSingleResult := mdbabstractlayer.NewMockIMongoSingleResult(ctrl)

	plugID, _ := primitive.ObjectIDFromHex("5e7600f01631b82037c9bf76")

	// #1
	postFilter :=
		primitive.E{
			Key:   "_id",
			Value: plugID,
		}
	filter := bson.D{
		postFilter,
	}

	// #2
	var upvotePercentageINT32 int32 = 0
	upvotePercentageData := primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{
				Key:   "upvotePercentage",
				Value: upvotePercentageINT32,
			},
		},
	}
	var score int32 = 0
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

	updatedDocument := &Post{}

	// Good
	// #1
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	// #2
	testCollection.EXPECT().FindOneAndUpdate(ctx, gomock.Eq(filter), gomock.Eq(update)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	// #3
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	expectedResult := &Post{
		ID:               primitive.ObjectID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Score:            0,
		Views:            0,
		Type:             "",
		Title:            "",
		Text:             "",
		URL:              "",
		Category:         "",
		Created:          "",
		UpvotePercentage: 0,
	}

	// USE MONKEY PATCHING LUKE
	// Это вызовы внутри CalcScoreData
	var repo *PostsRepo
	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(repo), "СalcScoreData",
		func(repo *PostsRepo, post *Post) *ScoredData {
			guard.Unpatch()

			result := &ScoredData{
				upvotePercentage: 0,
				score:            0,
			}
			return result
		})

	actualResult, err := service.CalcPostScoreData(reqCtx, plugID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.Equal(t, expectedResult, actualResult)

}

func TestCalcPostScoreData__ERROR_01(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	reqCtx := req.Context()
	ctx := gomock.Any()

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}
	testSingleResult := mdbabstractlayer.NewMockIMongoSingleResult(ctrl)

	plugID, _ := primitive.ObjectIDFromHex("5e7600f01631b82037c9bf76")

	// #1
	postFilter :=
		primitive.E{
			Key:   "_id",
			Value: plugID,
		}
	filter := bson.D{
		postFilter,
	}

	updatedDocument := &Post{}

	// Error #1
	errorText := "DB_ERROR"
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(fmt.Errorf(errorText))

	var expectedResult Post

	actualResult, err := service.CalcPostScoreData(reqCtx, plugID)
	if err != nil && err.Error() != errorText {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.NotEqual(t, expectedResult, actualResult)

}

func TestCalcPostScoreData__ERROR_02(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	reqCtx := req.Context()
	ctx := gomock.Any()

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}
	testSingleResult := mdbabstractlayer.NewMockIMongoSingleResult(ctrl)

	plugID, _ := primitive.ObjectIDFromHex("5e7600f01631b82037c9bf76")

	// #1
	postFilter :=
		primitive.E{
			Key:   "_id",
			Value: plugID,
		}
	filter := bson.D{
		postFilter,
	}

	// #2
	var upvotePercentageINT32 int32 = 0
	upvotePercentageData := primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{
				Key:   "upvotePercentage",
				Value: upvotePercentageINT32,
			},
		},
	}
	var score int32 = 0
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

	updatedDocument := &Post{}

	var expectedResult Post

	errorText := "DB_ERROR"
	// Error #2
	// #1
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	// #2
	testCollection.EXPECT().FindOneAndUpdate(ctx, gomock.Eq(filter), gomock.Eq(update)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(fmt.Errorf(errorText))

	actualResult, err := service.CalcPostScoreData(reqCtx, plugID)
	if err != nil && err.Error() != errorText {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.NotEqual(t, expectedResult, actualResult)
	assert.Nil(t, actualResult)

}

func TestCalcPostScoreData__ERROR_03(t *testing.T) {

	req := httptest.NewRequest("GET", "/", nil)
	reqCtx := req.Context()
	ctx := gomock.Any()

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}
	testSingleResult := mdbabstractlayer.NewMockIMongoSingleResult(ctrl)

	plugID, _ := primitive.ObjectIDFromHex("5e7600f01631b82037c9bf76")

	// #1
	postFilter :=
		primitive.E{
			Key:   "_id",
			Value: plugID,
		}
	filter := bson.D{
		postFilter,
	}

	// #2
	var upvotePercentageINT32 int32 = 0
	upvotePercentageData := primitive.E{
		Key: "$set",
		Value: bson.D{
			primitive.E{
				Key:   "upvotePercentage",
				Value: upvotePercentageINT32,
			},
		},
	}
	var score int32 = 0
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

	updatedDocument := &Post{}

	var expectedResult Post

	errorText := "DB_ERROR"

	// Error #3
	// #1
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	// #2
	testCollection.EXPECT().FindOneAndUpdate(ctx, gomock.Eq(filter), gomock.Eq(update)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	// // #3
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(fmt.Errorf(errorText))

	actualResult, err := service.CalcPostScoreData(reqCtx, plugID)
	if err != nil && err.Error() != errorText {
		t.Errorf("unexpected err: %s", err)
		return
	}

	assert.NotEqual(t, expectedResult, actualResult)
	assert.Nil(t, actualResult)
}
