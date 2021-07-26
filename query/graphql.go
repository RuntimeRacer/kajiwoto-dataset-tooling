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
