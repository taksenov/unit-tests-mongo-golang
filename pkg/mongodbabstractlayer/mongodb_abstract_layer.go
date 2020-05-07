package mongodbabstractlayer

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockgen command:
// mockgen -source=mongodb_abstract_layer.go -destination=mongodb_abstract_layer_mock.go -package=mongodbabstractlayer IMongoDatabase

/*
================
Интерфейсы
================
*/

// IMongoDatabase интерфейс БД MongoDB
type IMongoDatabase interface {
	Collection(name string) IMongoCollection
}

// IMongoCollection интерфейс коллекции
type IMongoCollection interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) IMongoSingleResult
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) IMongoSingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
}

// IMongoSingleResult интерфейс для одиночного результата
type IMongoSingleResult interface {
	Decode(v interface{}) error
}

// IMongoCursor интерфейс курсора
type IMongoCursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
}

/*
================
Структуры
================
*/

// MongoCollection коллекция
type MongoCollection struct {
	Сoll *mongo.Collection
}

// MongoSingleResult одиночный результат
type MongoSingleResult struct {
	sr *mongo.SingleResult
}

// MongoCursor курсор
type MongoCursor struct {
	cur *mongo.Cursor
}

/*
================
Методы структур
================
*/

// Decode раскодировать одиночный результат в заданную структуру
func (msr *MongoSingleResult) Decode(v interface{}) error {
	return msr.sr.Decode(v)
}

// Close закрыть курсор
func (mc *MongoCursor) Close(ctx context.Context) error {
	return mc.cur.Close(ctx)
}

// Next следующий элемент в курсоре
func (mc *MongoCursor) Next(ctx context.Context) bool {
	return mc.cur.Next(ctx)
}

// Decode раскодировать курсор в заданную структуру
func (mc *MongoCursor) Decode(val interface{}) error {
	return mc.cur.Decode(val)
}

// Find найти элементы в коллекции
func (mc *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (IMongoCursor, error) {
	cursorResult, err := mc.Сoll.Find(ctx, filter, opts...)
	return &MongoCursor{cur: cursorResult}, err
}

// FindOne найти элемент в коллекции
func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) IMongoSingleResult {
	singleResult := mc.Сoll.FindOne(ctx, filter, opts...)
	return &MongoSingleResult{sr: singleResult}
}

// FindOneAndUpdate найти элемент в коллекции обновить
func (mc *MongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) IMongoSingleResult {
	singleResult := mc.Сoll.FindOneAndUpdate(ctx, filter, update, opts...)
	return &MongoSingleResult{sr: singleResult}
}

// FindOneAndDelete найти элемент в коллекции и удалить
func (mc *MongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) IMongoSingleResult {
	singleResult := mc.Сoll.FindOneAndDelete(ctx, filter, opts...)
	return &MongoSingleResult{sr: singleResult}
}

// InsertOne добавить документ в коллекцию
func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.Сoll.InsertOne(ctx, document)
	return id.InsertedID, err
}
