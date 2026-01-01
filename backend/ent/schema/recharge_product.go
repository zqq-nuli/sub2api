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

// RechargeProduct holds the schema definition for the RechargeProduct entity.
type RechargeProduct struct {
	ent.Schema
}

func (RechargeProduct) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "recharge_products"},
	}
}

func (RechargeProduct) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(255).
			NotEmpty(),
		field.Float("amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,2)"}).
			Comment("Price in CNY"),
		field.Float("balance").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Comment("Base balance in USD"),
		field.Float("bonus_balance").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0).
			Comment("Bonus balance in USD"),
		field.String("description").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Default(""),
		field.Int("sort_order").
			Default(0),
		field.Bool("is_hot").
			Default(false),
		field.String("discount_label").
			MaxLen(50).
			Default(""),
		field.String("status").
			MaxLen(20).
			Default("active").
			Comment("active/inactive"),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

func (RechargeProduct) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("orders", Order.Type),
	}
}

func (RechargeProduct) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status", "sort_order"),
	}
}
