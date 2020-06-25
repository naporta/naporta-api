package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Vendedor struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Condominio  string             `bson:"condominio" json:"condominio"`
	Empresa     string             `bson:"empresa" json:"empresa"`
	Responsavel string             `bson:"responsavel" json:"responsavel"`
	Produtos    string             `bson:"produtos" json:"produtos"`
	Whatsapp    int64              `bson:"whatsapp" json:"whatsapp"`
	Bloco       int64              `bson:"bloco" json:"bloco"`
	Apt         int64              `bson:"apt" json:"apt"`
	Pagamento   string             `bson:"pagamento" json:"pagamento"`
	Categoria   string             `bson:"categoria" json:"categoria"`
}

type Connection struct {
	User     string
	Password string
	Server   string
	Database string
	client   *mongo.Client
}

func (c *Connection) Connect() error {
	ctx := context.TODO()

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s/%s",
		c.User, c.Password, c.Server, c.Database)

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *Connection) FindAll() ([]Vendedor, error) {
	ctx := context.TODO()

	collection := c.client.Database(c.Database).Collection("vendedores")
	var vendedores []Vendedor

	findOptions := options.Find()
	query, err := collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	defer query.Close(ctx)

	for query.Next(ctx) {
		var elem Vendedor
		err := query.Decode(&elem)
		if err != nil {
			return nil, err
		}

		vendedores = append(vendedores, elem)
	}
	if err := query.Err(); err != nil {
		return nil, err
	}

	return vendedores, nil
}

func (c *Connection) Insert(v Vendedor) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	bson, err := bson.Marshal(v)
	res, err := c.client.Database(c.Database).Collection("vendedores").InsertOne(ctx, bson)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Connection) FindByID(id string) (Vendedor, error) {
	ctx := context.TODO()
	var result Vendedor
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Vendedor{}, err
	}

	filter := bson.D{{"_id", objID}}

	err = c.client.Database(c.Database).Collection("vendedores").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Vendedor{}, err
	}

	return result, nil
}

func (c *Connection) Delete(v Vendedor) (*mongo.DeleteResult, error) {
	ctx := context.TODO()
	filter := bson.D{{"_id", v.ID}}
	res, err := c.client.Database(c.Database).Collection("vendedores").DeleteOne(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Connection) Update(v Vendedor) (*mongo.UpdateResult, error) {
	ctx := context.TODO()

	filter := bson.D{{"_id", v.ID}}

	newData, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	update := bson.D{{"$set", newData}}

	res, err := c.client.Database(c.Database).Collection("vendedores").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}
