package posts

import (
	"testing"

	mdbabstractlayer "unit-tests-mongo-golang/pkg/mongodbabstractlayer"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

func TestСalcScoreData__01(t *testing.T) {

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}

	plugID, _ := primitive.ObjectIDFromHex("5e7600f01631b82037c9bf76")

	tempVotesStruct :=
		&PostVote{
			User: "monorepo",
			Vote: 1,
		}

	votes := []*PostVote{}
	votes = append(votes, tempVotesStruct)

	post := &Post{
		ID:    plugID,
		Votes: votes,
	}

	expectedResult := &ScoredData{
		upvotePercentage: 100,
		score:            1,
	}

	actualResult := service.СalcScoreData(post)

	assert.Equal(t, expectedResult, actualResult)
	assert.NotNil(t, actualResult)

}

func TestСalcScoreData__02(t *testing.T) {

	// Передаём t сюда, для того чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	testCollection := mdbabstractlayer.NewMockIMongoCollection(ctrl)
	service := &PostsRepo{
		Collection: testCollection,
	}

	plugID, _ := primitive.ObjectIDFromHex("5e7600f01631b82037c9bf76")

	tempVotesStruct :=
		&PostVote{
			User: "monorepo",
			Vote: -1,
		}

	votes := []*PostVote{}
	votes = append(votes, tempVotesStruct)

	post := &Post{
		ID:    plugID,
		Votes: votes,
	}

	expectedResult := &ScoredData{
		upvotePercentage: 0,
		score:            -1,
	}

	actualResult := service.СalcScoreData(post)

	assert.Equal(t, expectedResult, actualResult)
	assert.NotNil(t, actualResult)

}
