package repository

import (
	"context"

	"github.com/ihezebin/promisys/component/storage"
	"github.com/ihezebin/promisys/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type exampleMysqlRepository struct {
	db *gorm.DB
}

func (e *exampleMysqlRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	example.Id = primitive.NewObjectID().Hex()
	err := e.db.Create(example).Error
	if err != nil {
		return err
	}

	return nil
}

func (e *exampleMysqlRepository) FindByUsername(ctx context.Context, username string) (example *entity.Example, err error) {
	example = &entity.Example{}
	err = e.db.Where("username = ?", username).First(example).Error
	if err != nil {
		return nil, err
	}

	return example, nil
}

func (e *exampleMysqlRepository) FindByEmail(ctx context.Context, email string) (example *entity.Example, err error) {
	example = &entity.Example{}
	err = e.db.Where("email = ?", email).First(example).Error
	if err != nil {
		return nil, err
	}

	return example, nil
}

func NewExampleMysqlRepository() ExampleRepository {
	return &exampleMysqlRepository{
		db: storage.MySQLStorageDatabase(),
	}
}

var _ ExampleRepository = (*exampleMysqlRepository)(nil)
