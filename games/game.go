package games

// State - state of the game
type State interface {
	String() string
	Bytes() []byte
}

// Result - result of the game
type Result interface {
	String() string
	Bytes() []byte
}

// Game interface for working with different games
type Game interface {
	Init()
	Images() (imgName1, imgName2 string)
	Snapshots() (shot1, shot2 []byte)
	SaveSnapshots(shot1, shot2 []byte) (gameErr error)
	GetState() (state State, fin bool)
	GetResult() (result Result)
}
