package repositories

import (
	"context"
	"errors"

	domain "task_management/Domain"
	"task_management/db"
	"task_management/usecases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	CountDocuments(cts context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// implementation of userrepository and uses monogb

type UserRepository struct {
	Collection IMongoCollection
	Context    context.Context
}

// constructor to initialize userRepositoryImpl
func NewUserRepository() usecases.IUserRepository {

	col := db.GetUsersCollection()
	ctx := context.Background()

	return &UserRepository{
		Collection: col,
		Context:    ctx,
	}
}

// inserts a new user to mongodb collection
func (r *UserRepository) CreateUser(user *domain.User) error {
	user.ID = primitive.NewObjectID()

	_, err := r.Collection.InsertOne(r.Context, user)
	return err
}

// retrieves a user based on the given username
func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {

	var user domain.User
	err := r.Collection.FindOne(r.Context, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// counts the number of users that matches the username
func (r *UserRepository) CountByUsername(username string) (int64, error) {
	return r.Collection.CountDocuments(r.Context, bson.M{"username": username})
}

//counts the total number of user documents in the collection

func (r *UserRepository) CountAll() (int64, error) {
	return r.Collection.CountDocuments(r.Context, bson.M{})
}

// updates the user role to admin based the id provided
func (r *UserRepository) PromoteUser(userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user id")
	}
	//filters the id
	filter := bson.M{"_id": objID}
	// the thing to be updated
	update := bson.M{"$set": bson.M{"role": "Admin"}}

	result, err := r.Collection.UpdateOne(r.Context, filter, update)
	if err != nil {
		return err
	}
	//if no result or the count is 0 it returns error
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
