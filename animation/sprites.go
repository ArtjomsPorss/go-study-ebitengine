package main

import (
  "image"
  _ "image/png"
  "log"

  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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
  CowStand [8][10]*ebiten.Image

  // player sheets
  RunnerImage *ebiten.Image
  StandingImage *ebiten.Image
  PlayerYPos int // used to identify center bottom of sprite when drawing
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
  ss.CowYPos = 150 // to adjust position to cow's feet where it stands
  ss.CowSpriteWidth = 164
  ss.CowSpriteHeight = 156

  // cow stand
  sheet := loadSpriteSheet("resources/cow-stand.png")
  // log.Printf("cowsheet bounds [%v]", sheet.Bounds())

  // load sheet into array
  for y:=0; y < len(ss.CowStand); y++ {
    for x:=0; x < len(ss.CowStand[y]); x++ {
      rect := image.Rect(x * ss.CowSpriteWidth,y * ss.CowSpriteHeight,(x + 1) * ss.CowSpriteWidth,(y + 1) * ss.CowSpriteHeight)
      ss.CowStand[y][x] = sheet.SubImage(rect).(*ebiten.Image)  
      // log.Printf("cow image bounds[%v] width[%v] height[%v]", ss.CowStand[y][x].Bounds(), ss.CowSpriteWidth, ss.CowSpriteHeight)
    }
  }
}

func loadSpriteSheet(path string) *ebiten.Image {
  img, _, err := ebitenutil.NewImageFromFile(path)
  if err != nil {
    log.Fatal(err)
  }
  return ebiten.NewImageFromImage(img)
}
