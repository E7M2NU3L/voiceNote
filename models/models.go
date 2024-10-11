package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	Profile   Profile            `bson:"profile"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Profile struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
	Bio       string `bson:"bio"`
	// ... other profile details
}

// ... other model definitions
type Freelancer struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"userId"`
	Skills []string           `bson:"skills"`
	// ... other freelancer-specific details
}

type Client struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"userId"`
	// ... other client-specific details
}

type Project struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ClientID     primitive.ObjectID `bson:"clientId"`
	FreelancerID primitive.ObjectID `bson:"freelancerId"`
	Title        string             `bson:"title"`
	Description  string             `bson:"description"`
	// ... other project details
}

type Invoice struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID    primitive.ObjectID `bson:"projectId"`
	FreelancerID primitive.ObjectID `bson:"freelancerId"`
	Amount       float64            `bson:"amount"`
	Status       string             `bson:"status"`
	// ... other invoice details
}
