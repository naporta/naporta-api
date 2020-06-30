package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vendedor struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Condominio string             `bson:"condominio" json:"condominio"`
	Nome       string             `bson:"nome" json:"nome"`
	Empresa    string             `bson:"empresa" json:"empresa"`
	Profissao  string             `bson:"profissao" json:"profissao"`
	Produtos   []string           `bson:"produtos" json:"produtos"`
	Whatsapp   int64              `bson:"whatsapp" json:"whatsapp"`
	Facebook   string             `bson:"facebook" json:"facebook"`
	Instagram  string             `bson:"instagram" json:"instagram"`
	Bloco      int64              `bson:"bloco" json:"bloco"`
	Apt        int64              `bson:"apt" json:"apt"`
	Pagamento  []string           `bson:"pagamento" json:"pagamento"`
	Tags       []string           `bson:"tags" json:"tags"`
	Verificado bool               `bson:"verificado" json:"verificado"`
	Assinante  bool               `bson:"assinante" json:"assinante"`
	Assinatura *time.Time         `bson:"assinatura" json:"assinatura"`
}
