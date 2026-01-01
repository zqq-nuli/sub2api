package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

func (Order) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "orders"},
	}
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.String("order_no").
			MaxLen(32).
			NotEmpty().
			Unique(),
		field.Int64("user_id").
			Positive(),
		field.Int64("product_id").
			Optional().
			Nillable(),
		field.String("product_name").
			MaxLen(255).
			NotEmpty(),
		field.Float("amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,2)"}).
			Comment("Payment amount in CNY"),
		field.Float("bonus_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Bonus amount in USD"),
		field.Float("actual_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Comment("Actual credited amount in USD"),
		field.String("payment_method").
			MaxLen(50).
			NotEmpty().
			Comment("alipay/wxpay/epusdt"),
		field.String("payment_gateway").
			MaxLen(50).
			Default("epay"),
		field.String("trade_no").
			MaxLen(255).
			Optional().
			Nillable().
			Comment("Third-party order number"),
		field.String("status").
			MaxLen(20).
			Default("pending").
			Comment("pending/paid/failed/expired"),
		field.Time("created_at"),
		field.Time("paid_at").
			Optional().
			Nillable(),
		field.Time("expired_at").
			Comment("Order expiration time (15 minutes)"),
		field.String("notes").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default(""),
		field.String("callback_data").
			SchemaType(map[string]string{dialect.Postgres: "jsonb"}).
			Optional().
			Nillable(),
	}
}

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("orders").
			Field("user_id").
			Required().
			Unique(),
		edge.From("product", RechargeProduct.Type).
			Ref("orders").
			Field("product_id").
			Unique(),
	}
}

func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}
