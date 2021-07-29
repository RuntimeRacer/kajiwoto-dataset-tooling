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
	ID          graphql.String
	UserMessage graphql.String
	Message     graphql.String
	/*	ASM Cheat sheet

		... No idea what "ASM" stands for in this context. But its holding the emotional values for the Kaji dialogues.

		HAPPY => Happy or Excited
		SAD   => Sad
		HUNGRY => Hungry
		FULL => Full
		EXCITED => Excited
		ANGRY => Angry
		SCARED => Scared
		BULLIED => Bullied
		ATTACKED => Attacked
		POWERED => Powered
		DRUNK => Drunk
		SICK => Sick
		SLEEPY => Sleepy
	*/
	ASM graphql.String
	/*  Condition cheat sheet

	Seems to be oriented from linux permissions. five digits; last two seem to be never used.

	//// Attachment Keys
	XX 1 XX Disliked
	XX 2 XX Any-Emotion
	XX 3 XX Liked
	XX 4 XX -- NOT USED --
	XX 5 XX Disliked/Neutral

	//// Daytime Keys
	// Default (single) conditions
	1 XXXX Early Morning AM
	2 XXXX Morning
	3 XXXX Afternoon
	4 XXXX Evening
	5 XXXX Middle of Sleep AM
	6 XXXX -- NOT USED --
	// Combined Conditions
	7 XXXX Early Morning AM - Morning
	8 XXXX Evening AM - Middle of Sleep AM
	9 XXXX Morning - Afternoon

	//// Last seen keys
	X 1 XXX Seen 2 hrs ago
	X 2 XXX Seen 12 hrs ago
	X 3 XXX Seen 5 days ago
	X 4 XXX Seen 5 days+ ago
	*/
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
