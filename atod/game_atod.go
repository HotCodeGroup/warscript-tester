package atod

import (
	"encoding/json"

	"github.com/HotCodeGroup/warscript-tester/games"
	"github.com/pkg/errors"
)

const (
	lEPS  = float64(.0000001)
	qEPS  = float64(.01)
	image = "atod"
)

func (a *Atod) Init() {
	a = stdField()
}

func (a *Atod) Images() (string, string) {
	return image, image
}

func (a *Atod) Snapshots() (shot1, shot2 []byte) {
	shot1, _ = json.Marshal(a.createShot1())
	shot2, _ = json.Marshal(a.createShot2())
	return
}

func (a *Atod) SaveSnapshots(shot1, shot2 []byte) error {
	var s1, s2 shot
	err1 := json.Unmarshal(shot1, &s1)
	if err1 != nil {
		a.isEnded = true
		a.occuredError = games.ErrPlayer1Fail
		return errors.Wrap(err1, games.ErrPlayer1Fail.Error())
	}

	err2 := json.Unmarshal(shot2, &s2)
	if err2 != nil {
		a.isEnded = true
		a.occuredError = games.ErrPlayer2Fail
		return errors.Wrap(err2, games.ErrPlayer2Fail.Error())
	}

	a.loadSnapShots(&s1, &s2)
	return nil
}

func (a *Atod) GetInfo() (info games.Info) {
	i := &Info{
		Player1Dropzone: dropzoneToResp(a.dropzone1, a.heihgt, a.width),
		Player2Dropzone: dropzoneToResp(a.dropzone2, a.heihgt, a.width),
		Ratio:           a.width / a.heihgt,
	}

	return i
}

func (a *Atod) GetState() (state games.State, fin bool) {
	return &State{
		Projectiles: projectilesToResp(a.projectiles, a.heihgt, a.width),
		Obstacles:   obstaclesToResp(a.obstacles, a.heihgt, a.width),
		P1Units:     unitsToResp(a.player1Units, a.heihgt, a.width),
		P2Units:     unitsToResp(a.player2Units, a.heihgt, a.width),
		P1Flags:     flagsToResp(a.flags1, a.heihgt, a.width),
		P2Flags:     flagsToResp(a.flags2, a.heihgt, a.width),
	}, a.isEnded
}

func (a *Atod) GetResult() (result games.Result) {
	if !a.isEnded {
		return nil
	}

	return &Result{
		Projectiles: projectilesToResp(a.projectiles, a.heihgt, a.width),
		Obstacles:   obstaclesToResp(a.obstacles, a.heihgt, a.width),
		P1Units:     unitsToResp(a.player1Units, a.heihgt, a.width),
		P2Units:     unitsToResp(a.player2Units, a.heihgt, a.width),
		P1Flags:     flagsToResp(a.flags1, a.heihgt, a.width),
		P2Flags:     flagsToResp(a.flags2, a.heihgt, a.width),

		Winner: a.winner,
		Error:  a.occuredError,
	}
}