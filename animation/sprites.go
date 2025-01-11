package main

import (
  "image"
  _ "image/png"
  "log"

  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageContainerInterface interface {
  Images() [][]*ebiten.Image
  Width() int
  Height() int
  TopToFeet() int
}

type ImageContainer struct {
  sprites [][]*ebiten.Image
  height int
  width int
  topToFeet int
}

func (i ImageContainer) Images() [][]*ebiten.Image {
  return i.sprites
}
func (i ImageContainer) Height() int {
  return i.height
}
func (i ImageContainer) Width() int {
  return i.width
}
func (i ImageContainer) TopToFeet() int {
  return i.topToFeet
}

// spritesheet represents a collection of sprite images.
type SpriteSheet struct {
  // floor 
  Floor *ebiten.Image
  Width int
  Height int

  // wall - cliff - top/bottom and left/right
  WallTopBottom [13]*ebiten.Image
  WallLeftRight [13]*ebiten.Image
  CliffNarrowWidth int
  CliffWidth int
  CliffFullHeight int
  CliffHeight int
  CliffYDrawStartingPoint float64
  CliffYRightDrawStartingPoint float64 // when drawing right walls - substract
  // wall - cliff - corner
  WallCorner *ebiten.Image
  CornerWidth int
  CornerHeight int
  
  CowSpriteWidth int
  CowSpriteHeight int
  CowYPos int // center bottom of cow's sprite height

  CowWalkSpriteWidth int
  CowWalkSpriteHeight int
  CowAttackWidth int
  CowAttackHeight int

  CowWalk ImageContainer
  CowAttack ImageContainer
  CowStand ImageContainer
  CowToRender ImageContainer

  // player sheets
  RunnerImage *ebiten.Image
  StandingImage *ebiten.Image
  PlayerYPos int // used to identify center bottom of player's sprite when drawing
  AttackingImage *ebiten.Image
  ImageToRender *ebiten.Image
}

func LoadSprites() (*SpriteSheet) {
  sheet := &SpriteSheet{}
  
  // load floor
  sheet.loadFloorSpriteSheet()
  // load cliff
  sheet.loadTopBottomCliff()
  sheet.loadLeftRightCliff()
  // load player
  sheet.loadPlayer()
  // load cow - enemy
  sheet.loadCowSheet()
  return sheet
}

func (floorSheet *SpriteSheet) loadPlayer() {
  floorSheet.PlayerYPos = 39 // to adjust other sprites for the position of the player
  img, _, err := ebitenutil.NewImageFromFile("resources/druid-run.png")
  if err != nil {
    log.Fatal(err)
  }
  floorSheet.RunnerImage = ebiten.NewImageFromImage(img)

  img2, _, err := ebitenutil.NewImageFromFile("resources/druid-stand.png")
  if err != nil {
    log.Fatal(err)
  }
  floorSheet.StandingImage = ebiten.NewImageFromImage(img2)

  img3, _, err := ebitenutil.NewImageFromFile("resources/druid-attack.png")
  if err != nil {
    log.Fatal(err)
  }
  floorSheet.AttackingImage = ebiten.NewImageFromImage(img3)
}


// Load Spritesheet loads the embedded spritesheet
func (floorSheet *SpriteSheet) loadFloorSpriteSheet() {
  sheet := loadSpriteSheet("resources/floor-swamp.png")

  floorSheet.Width = 160
  floorSheet.Height = 80
  // load floor image
  floorSheet.Floor = sheet.SubImage(image.Rect(0,0,floorSheet.Width,floorSheet.Height)).(*ebiten.Image)
}

func (spriteSheet *SpriteSheet) loadTopBottomCliff() {
  sheet := loadSpriteSheet("resources/cliff1.png")
  spriteSheet.CliffWidth = 160
  spriteSheet.CliffHeight = 448
  spriteSheet.CliffNarrowWidth = 80 // narrow width is the actual cliff
  spriteSheet.CornerHeight = 448
  spriteSheet.CliffFullHeight = 512
  spriteSheet.CliffYDrawStartingPoint = 432.0
  // load sheet an array
  for i:=0; i < len(spriteSheet.WallTopBottom); i++ {
    spriteSheet.WallTopBottom[i] = sheet.SubImage(image.Rect(i * spriteSheet.CliffWidth,0,(i+1) * spriteSheet.CliffWidth,spriteSheet.CliffFullHeight)).(*ebiten.Image)  
  }
}

func loadCorner(spriteSheet *SpriteSheet) {

}

func (spriteSheet *SpriteSheet) loadLeftRightCliff() {
  spriteSheet.CliffYRightDrawStartingPoint = 360 
  sheet := loadSpriteSheet("resources/cliff2.png")
  // load sheet an array
  for i:=0; i < len(spriteSheet.WallLeftRight); i++ {
    spriteSheet.WallLeftRight[i] = sheet.SubImage(image.Rect(i * spriteSheet.CliffWidth,0,(i+1) * spriteSheet.CliffWidth,spriteSheet.CliffFullHeight)).(*ebiten.Image)  
    // log.Printf("cliff image bounds[%v]", spriteSheet.WallLeftRight[i].Bounds())
  }
}

func (ss *SpriteSheet) loadCowSheet() {
  cowStandTopToFeet := 150 // to adjust position to cow's feet where it stands
  cowStandWidth := 164
  cowStandHeight := 156

  // COW STAND
  sheet := loadSpriteSheet("resources/cow-stand.png")
  // log.Printf("cowsheet bounds [%v]", sheet.Bounds())

	// Define dimensions for the 2D slice
	width := 8
	height := 10

	// Create a 2D slice of *ebiten.Image
	cowStand := make([][]*ebiten.Image, height) // Height rows
	for i := range cowStand {
		cowStand[i] = make([]*ebiten.Image, width) // Width columns
	}

  // load sheet into array
  for y:=0; y < len(cowStand); y++ {
    for x:=0; x < len(cowStand[y]); x++ {
      rect := image.Rect(x * cowStandWidth,y * cowStandHeight,(x + 1) * cowStandWidth,(y + 1) * cowStandHeight)
      cowStand[y][x] = sheet.SubImage(rect).(*ebiten.Image)  
    }
  }

  ss.CowStand = ImageContainer {
    sprites: cowStand,
    height: cowStandHeight,
    width: cowStandWidth,
    topToFeet: cowStandTopToFeet,
  }

  // COW WALK
  sheet = loadSpriteSheet("resources/cow-walk.png")
  cowWalkWidth := 157
  cowWalkHeight := 151
  cowWalkTopToFeet := 140

	// Define dimensions for the 2D slice
	width = 8
	height = 8 

	// Create a 2D slice of *ebiten.Image
	cowWalk := make([][]*ebiten.Image, height) // Height rows
	for i := range cowWalk {
		cowWalk[i] = make([]*ebiten.Image, width) // Width columns
	}

  // load sheet into array
  for y:=0; y < len(cowWalk); y++ {
    for x:=0; x < len(cowWalk[y]); x++ {
      rect := image.Rect(x * cowWalkWidth,y * cowWalkHeight,(x + 1) * cowWalkWidth,(y + 1) * cowWalkHeight)
      cowWalk[y][x] = sheet.SubImage(rect).(*ebiten.Image)  
    }
  }
  ss.CowWalk = ImageContainer{
    sprites: cowWalk,
    height: cowWalkHeight,
    width: cowWalkWidth,
    topToFeet: cowWalkTopToFeet,
  }

  // COW ATTACK
  sheet = loadSpriteSheet("resources/cow-attack.png")
  cowAttackHeight := 218
  cowAttackWidth := 262
  cowAttackTopToFeet := 174

	// Define dimensions for the 2D slice
	width = 19
	height = 8 

	// Create a 2D slice of *ebiten.Image
	cowAttack := make([][]*ebiten.Image, height) // Height rows
	for i := range cowAttack {
		cowAttack[i] = make([]*ebiten.Image, width) // Width columns
	}

  // load sheet into array
  for y:=0; y < len(cowAttack); y++ {
    for x:=0; x < len(cowAttack[y]); x++ {
      rect := image.Rect(x * cowAttackWidth,y * cowAttackHeight,(x + 1) * cowAttackWidth,(y + 1) * cowAttackHeight)
      cowAttack[y][x] = sheet.SubImage(rect).(*ebiten.Image)
    }
  }
  ss.CowAttack = ImageContainer{
    sprites: cowAttack,
    height: cowAttackHeight,
    width: cowAttackWidth,
    topToFeet: cowAttackTopToFeet,
  }
}

func loadSpriteSheet(path string) *ebiten.Image {
  img, _, err := ebitenutil.NewImageFromFile(path)
  if err != nil {
    log.Fatal(err)
  }
  return ebiten.NewImageFromImage(img)
}
