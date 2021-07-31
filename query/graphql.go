package query

import "github.com/runtimeracer/go-graphql-client"

type Plus struct {
	ExpireAt  uint64
	Cancelled graphql.Boolean
	Icon      graphql.Int
	Coins     graphql.Int
	Type      graphql.String
}

type Creator struct {
	AllowSubscriptions graphql.Boolean
	DatasetTags        []graphql.String
}

type Profile struct {
	FirstName   graphql.String
	LastName    graphql.String
	Description graphql.String
	Gender      graphql.String
	Birthday    graphql.String
	PhotoUri    graphql.String
}

type Email struct {
	Address  graphql.String
	Verified graphql.Boolean
}

type Settings struct {
	PersonalRoomOrder []graphql.String
	FavoriteRoomIds   []graphql.String
	FavoriteEmojis    []graphql.String
}

type User struct {
	ID          graphql.String
	Activated   graphql.Boolean
	Moderator   graphql.Boolean
	Username    graphql.String
	DisplayName graphql.String
	Plus        Plus
	Creator     Creator
	Profile     Profile
	Email       Email
}

type Login struct {
	AuthToken string
	User      User
	Settings  Settings
}

type Announcement struct {
	Date      uint64
	Title     graphql.String
	Emojis    graphql.String
	Content   []graphql.String
	TextColor graphql.String
}

type Welcome struct {
	WebVersion   graphql.String
	Announcement Announcement
}

type LoginResult struct {
	Login   Login
	Welcome Welcome
}

type AITrainerGroup struct {
	ID              graphql.String
	Name            graphql.String
	Count           graphql.Int
	Deleted         graphql.Boolean
	Description     graphql.String
	Documents       []AIDocument
	DominantColors  []graphql.String
	Kudos           Kudos
	NSFW            graphql.Boolean
	Personalities   interface{}
	PetSpeciesIds   []graphql.String
	Price           graphql.Int
	ProfilePhotoUri graphql.String
	Purchased       graphql.Boolean
	Status          graphql.String
	Tags            []graphql.String
	UpdatedAt       uint64
	User            User
}

type Kudos struct {
	ID       graphql.String
	Upvoted  graphql.Boolean
	Upvotes  graphql.Int
	Comments graphql.Int
}

type AITrained struct {
	ID               graphql.String
	UserMessage      graphql.String
	Message          graphql.String
	ASM              graphql.String
	Condition        graphql.String
	Deleted          graphql.Boolean
	History          []graphql.String
	AITrainerGroupID graphql.String
}

type AIDocument struct {
	ID          graphql.String
	Order       graphql.String
	Title       graphql.String
	Content     graphql.String
	QueueStatus graphql.String
	QueuedAt    graphql.String
	BuiltAt     uint64
	CreatedAt   uint64
	UpdatedAt   uint64
}

type TrainDatasetResult struct {
	Calibrations []TrainCalibration
	Count        graphql.Int
}

type TrainCalibration struct {
	Chance      graphql.Float
	UserMessage graphql.String
	Message     graphql.String
}

type AITraining struct {
	Condition   graphql.String
	UserMessage graphql.String
	Message     graphql.String
}
