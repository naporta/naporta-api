package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (c *Connection) Insert(v Vendedor) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	novo := bson.M{
		"condominio": v.Condominio,
		"nome":       v.Nome,
		"empresa":    v.Empresa,
		"profissao":  v.Profissao,
		"produtos":   v.Produtos,
		"whatsapp":   v.Whatsapp,
		"facebook":   v.Facebook,
		"instagram":  v.Instagram,
		"bloco":      v.Bloco,
		"apt":        v.Apt,
		"pagamento":  v.Pagamento,
		"categoria":  v.Categoria,
		"tags":       v.Tags,
		"verificado": v.Verificado,
		"assinante":  v.Assinante,
		"assinatura": v.Assinatura,
	}
	res, err := c.client.Database(c.Database).Collection("vendedores").InsertOne(ctx, novo)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Connection) FindAll(condominio string, categoria string) ([]bson.M, error) {
	ctx := context.TODO()

	collection := c.client.Database(c.Database).Collection("vendedores")
	var vendedores []bson.M

	var query bson.D
	if condominio != "" {
		query = bson.D{
			primitive.E{Key: "verificado", Value: true},
			primitive.E{Key: "condominio", Value: condominio},
			primitive.E{Key: "tags", Value: categoria},
		}
	} else {
		query = bson.D{
			primitive.E{Key: "verificado", Value: true},
		}
	}

	filter := options.Find()
	filter.SetProjection(bson.M{
		"_id":        0,
		"verificado": 0,
		"assinante":  0,
		"assinatura": 0,
	})
	result, err := collection.Find(ctx, query, filter)
	if err != nil {
		return nil, err
	}
	if err = result.All(ctx, &vendedores); err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return vendedores, nil
}

func (c *Connection) FindByID(id string) (Vendedor, error) {
	ctx := context.TODO()
	var result Vendedor
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Vendedor{}, err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	err = c.client.Database(c.Database).Collection("vendedores").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return Vendedor{}, err
	}

	return result, nil
}

func (c *Connection) FindOneFalse(condominio string) (Vendedor, error) {
	ctx := context.TODO()
	var result Vendedor

	var query bson.D
	if condominio != "" {
		query = bson.D{
			primitive.E{Key: "verificado", Value: false},
			primitive.E{Key: "condominio", Value: condominio},
		}
	} else {
		query = bson.D{primitive.E{Key: "verificado", Value: false}}
	}
	err := c.client.Database(c.Database).Collection("vendedores").FindOne(ctx, query).Decode(&result)
	if err != nil {
		return Vendedor{}, err
	}

	return result, nil
}

func (c *Connection) GetTags(condominio string) (bson.M, error) {
	ctx := context.TODO()

	collection := c.client.Database(c.Database).Collection("vendedores")
	var tags []bson.M

	aggr := []bson.M{
		{"$match": bson.M{"condominio": condominio}},
		{"$unwind": "$tags"},
		{"$group": bson.M{"_id": 0, "tags": bson.M{"$addToSet": "$tags"}}},
		{"$project": bson.M{"_id": 0}},
	}

	result, err := collection.Aggregate(ctx, aggr)
	if err != nil {
		return nil, err
	}
	if err = result.All(ctx, &tags); err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return tags[0], nil
}

func (c *Connection) GetProdutos() (bson.M, error) {
	ctx := context.TODO()

	collection := c.client.Database(c.Database).Collection("vendedores")
	var tags []bson.M

	aggr := []bson.M{
		{"$unwind": "$produtos"},
		{"$group": bson.M{"_id": 0, "produtos": bson.M{"$addToSet": "$produtos"}}},
		{"$project": bson.M{"_id": 0}},
	}

	result, err := collection.Aggregate(ctx, aggr)
	if err != nil {
		return nil, err
	}
	if err = result.All(ctx, &tags); err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return tags[0], nil
}

func (c *Connection) GetCategorias() (bson.M, error) {
	ctx := context.TODO()

	collection := c.client.Database(c.Database).Collection("vendedores")
	var cat []bson.M

	aggr := []bson.M{
		{"$unwind": "$categoria"},
		{"$group": bson.M{"_id": 0, "categoria": bson.M{"$addToSet": "$categoria"}}},
		{"$project": bson.M{"_id": 0}},
	}

	result, err := collection.Aggregate(ctx, aggr)
	if err != nil {
		return nil, err
	}
	if err = result.All(ctx, &cat); err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	if err := result.Err(); err != nil {
		return nil, err
	}

	return cat[0], nil
}
func (c *Connection) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx := context.TODO()
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
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

	filter := bson.D{primitive.E{Key: "_id", Value: v.ID}}

	newData, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	update := bson.D{primitive.E{Key: "$set", Value: newData}}

	res, err := c.client.Database(c.Database).Collection("vendedores").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Connection) UpdateVerificado(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	ctx := context.TODO()

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	update := bson.M{"$set": bson.M{"verificado": true}}

	res, err := c.client.Database(c.Database).Collection("vendedores").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}
