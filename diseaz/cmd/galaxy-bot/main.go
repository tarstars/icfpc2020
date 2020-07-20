package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	gx "github.com/tarstars/icfpc2020/diseaz/interpreter"
)

type ICons = gx.ICons
type Int = gx.Int

type Result struct {
	Picture *gx.Picture `json:",inline"`
	Results []string    `json:""`
}

func (r *Result) AddResults(anss ...gx.Token) {
	for _, ans := range anss {
		r.Results = append(r.Results, ans.Galaxy())
	}
}

func command(c gx.Context, name string, program string) *GameResponse {
	rs := gx.ParseString(c, program)
	if len(rs) == 0 {
		return nil
	}

	logr := &Result{}
	logr.AddResults(rs...)
	logr.Picture = c.Picture()
	json.NewEncoder(os.Stdout).Encode(logr)

	gr := ParseGameResponse(rs[0].(ICons))
	grJSON, err := json.Marshal(gr)
	if err != nil {
		log.Panicf("GameResponse marshaling to JSON failed: %s", err)
	}
	log.Printf("GameResponse: %s", string(grJSON))
	if gr.Error {
		log.Panic(fmt.Errorf("%s failed", name))
	}

	return &gr
}

func checkListEnd(name string, endItem gx.Token) {
	endItemC, isCons := endItem.(ICons)
	if !isCons || !endItemC.IsNil() {
		log.Printf("Unexpected extra data for %s: %s", name, endItem.Galaxy())
	}
}

func ParseGameResponse(itemX0 ICons) (gr GameResponse) {
	status := itemX0.Car().(Int).V
	if status == 0 {
		return GameResponse{Error: true}
	}

	itemX1 := itemX0.Cdr().(ICons)
	gr.Stage = GameStage(itemX0.Car().(Int).V)
	itemX2 := itemX1.Cdr().(ICons)
	gr.StaticInfo = ParseGameStaticInfo(itemX2.Car().(ICons))
	itemX3 := itemX2.Cdr().(ICons)
	gr.State = ParseGameState(itemX3.Car().(ICons))

	checkListEnd("GameResponse", itemX3.Cdr())

	return gr
}

type GameResponse struct {
	Error      bool
	Stage      GameStage
	StaticInfo *GameStaticInfo
	State      *GameState
}

type GameState struct {
	Tick  int64
	X1    string
	Ships []*ShipAndCommands
}

func ParseGameState(itemX0 ICons) (gs *GameState) {
	if itemX0.IsNil() {
		return nil
	}

	gs = &GameState{}

	gs.Tick = itemX0.Car().(Int).V
	itemX1 := itemX0.Cdr().(ICons)
	gs.X1 = itemX1.Car().Galaxy()
	itemX2 := itemX1.Cdr().(ICons)
	gs.Ships = ParseShipsAndCommands(itemX2.Car().(ICons))

	checkListEnd("ParseGameState", itemX2.Cdr())

	return gs
}

func listNext(item ICons) (gx.Token, ICons) {
	return item.Car(), item.Cdr().(ICons)
}

func ParseShipsAndCommands(itemX0 ICons) []*ShipAndCommands {
	var ships []*ShipAndCommands
	for item := itemX0; !item.IsNil(); item = item.Cdr().(ICons) {
		ships = append(ships, ParseShip(item.Car().(ICons)))
	}
	return ships
}

type ShipAndCommands struct {
	Ship     ShipState
	Commands []string
}

func ParseShip(itemX0 ICons) (ship *ShipAndCommands) {
	if itemX0.IsNil() {
		return nil
	}

	ship.Ship = ParseShipState(itemX0.Car().(ICons))
	itemX1 := itemX0.Cdr().(ICons)
	ship.Commands = ParseCommands(itemX1.Car().(ICons))

	checkListEnd("ParseShip", itemX1.Cdr())

	return ship
}

func ParseCommands(itemX0 ICons) (cmds []string) {
	for item := itemX0; !item.IsNil(); item = item.Cdr().(ICons) {
		cmds = append(cmds, item.Car().Galaxy())
	}
	return cmds
}

type ShipState struct {
	Role     Role
	ID       int64
	Position gx.Point
	Velocity gx.Point
	X4       string
	X5       string
	X6       string
	X7       string
}

func ParseShipState(itemX0 ICons) (ss ShipState) {

	ss.Role = Role(itemX0.Car().(Int).V)
	itemX1 := itemX0.Cdr().(ICons)
	ss.ID = itemX1.Car().(Int).V
	itemX2 := itemX1.Cdr().(ICons)
	ss.Position = ParsePoint(itemX2.Car().(ICons))
	itemX3 := itemX2.Cdr().(ICons)
	ss.Velocity = ParsePoint(itemX3.Car().(ICons))
	itemX4 := itemX3.Cdr().(ICons)
	ss.X4 = itemX4.Car().Galaxy()
	itemX5 := itemX4.Cdr().(ICons)
	ss.X5 = itemX5.Car().Galaxy()
	itemX6 := itemX4.Cdr().(ICons)
	ss.X6 = itemX6.Car().Galaxy()
	itemX7 := itemX6.Cdr().(ICons)
	ss.X7 = itemX7.Car().Galaxy()

	checkListEnd("ParseShipState", itemX7.Cdr())

	return ss
}

func ParsePoint(itemX0 ICons) (p gx.Point) {

	p.X = int(itemX0.Car().(Int).V)
	itemX1 := itemX0.Cdr().(ICons)
	p.Y = int(itemX1.Car().(Int).V)

	checkListEnd("ParsePoint", itemX1.Cdr())

	return p
}

type GameStaticInfo struct {
	X0   int64
	Role Role
	X2   string
	X3   string
	X4   string
}

func ParseGameStaticInfo(infoX0Item ICons) *GameStaticInfo {
	if infoX0Item.IsNil() {
		return nil
	}
	gi := GameStaticInfo{}

	gi.X0 = infoX0Item.Car().(Int).V
	infoX1Item := infoX0Item.Cdr().(ICons)
	gi.Role = Role(infoX1Item.Car().(Int).V)
	infoX2Item := infoX1Item.Cdr().(ICons)
	gi.X2 = infoX2Item.Car().Galaxy()
	infoX3Item := infoX2Item.Cdr().(ICons)
	gi.X3 = infoX3Item.Car().Galaxy()
	infoX4Item := infoX3Item.Cdr().(ICons)
	gi.X4 = infoX4Item.Car().Galaxy()

	checkListEnd("GameStaticInfo", infoX4Item.Cdr())

	return &gi
}

type Role int

const (
	RoleAttacker Role = 0
	RoleDefender Role = 1
)

type GameStage int

const (
	GamePending  GameStage = 0
	GameStarted  GameStage = 1
	GameFinished GameStage = 2
)

func main() {
	flag.Parse()

	serverURL, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Panic(err)
	}
	serverURL.Path = "/aliens/send"
	playerKey, err := strconv.ParseInt(flag.Arg(1), 10, 64)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("ServerUrl: %s; PlayerKey: %d", serverURL, playerKey)

	c := gx.NewContext(serverURL)
	gs := command(c, "JOIN", fmt.Sprintf("ap send (2, %d, nil)", playerKey))
	if gs.Stage == GameFinished {
		return
	}
	command(c, "START", fmt.Sprintf("ap send (3, %d, (1, 1, 1, 1))", playerKey))
}
