package assets

// assets should contain all sorts of text variables and the like, there should not be any methods

const (
	GamesSavePath = "assets/game_save.json"
)

var ( // this variable should not be global, it is unique for each game, so it should be the fields of the Game structure
	BombCounter   int
	DeveloperMode bool // true = admin, false = all users
)
