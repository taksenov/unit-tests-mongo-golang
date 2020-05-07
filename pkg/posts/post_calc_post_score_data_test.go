package posts

import (
	"fmt"
	"net/http/httptest"
	"testing"

	mdbabstractlayer "unit-tests-mongo-golang/pkg/mongodbabstractlayer"

	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html


// Идея, в том чтобы получить замоканный результат, для того чтобы проверить условия в указанном методе
func TestCalcPostScoreDataMockedResults(t *testing.T) {

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

	updatedDocument := &Post{
		ID: plugID,
	}
	updatedDocument2 := &Post{}

	// Good
	// #1
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument)).Return(nil)

	// #2
	testCollection.EXPECT().FindOneAndUpdate(ctx, gomock.Eq(filter), gomock.Eq(update)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument2)).Return(nil)

	// #3
	testCollection.EXPECT().FindOne(ctx, gomock.Eq(filter)).Return(testSingleResult)
	testSingleResult.EXPECT().Decode(gomock.AssignableToTypeOf(updatedDocument2)).Return(nil)

	_, err := service.CalcPostScoreData(plugID, req)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	fmt.Printf("TEST updatedDocument= %#v \n", updatedDocument.ID.Hex())
	fmt.Printf("TEST testSingleResult= %#v \n", testSingleResult)
	fmt.Printf("TEST testSingleResult.EXPECT()= %#v \n", testSingleResult.EXPECT())
	// fmt.Printf("TEST testSingleResult= %#v \n", testSingleResult)

}
