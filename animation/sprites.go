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
  CowStand [8][10]*ebiten.Image
}


// Load Spritesheet loads the embedded spritesheet
func LoadFloorSpriteSheet() (*SpriteSheet, error) {
  img, _, err := ebitenutil.NewImageFromFile("resources/floor-swamp.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)

  floorSheet := &SpriteSheet{}
  floorSheet.Width = 160
  floorSheet.Height = 80
  // load floor image
  floorSheet.Floor = sheet.SubImage(image.Rect(0,0,floorSheet.Width,floorSheet.Height)).(*ebiten.Image)
  // load 
  loadTopBottomCliff(floorSheet)
  loadLeftRightCliff(floorSheet)

  floorSheet.loadCowSheet()
  return floorSheet, nil
}

func loadTopBottomCliff(spriteSheet *SpriteSheet) {
  img, _, err := ebitenutil.NewImageFromFile("resources/cliff1.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)
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

func loadLeftRightCliff(spriteSheet *SpriteSheet) {
  spriteSheet.CliffYRightDrawStartingPoint = 360 
  img, _, err := ebitenutil.NewImageFromFile("resources/cliff2.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)
  // load sheet an array
  for i:=0; i < len(spriteSheet.WallLeftRight); i++ {
    spriteSheet.WallLeftRight[i] = sheet.SubImage(image.Rect(i * spriteSheet.CliffWidth,0,(i+1) * spriteSheet.CliffWidth,spriteSheet.CliffFullHeight)).(*ebiten.Image)  
    // log.Printf("cliff image bounds[%v]", spriteSheet.WallLeftRight[i].Bounds())
  }
}

func (ss *SpriteSheet) loadCowSheet() {
  ss.CowSpriteWidth = 164
  ss.CowSpriteHeight = 156

  // cow stand
  img, _, err := ebitenutil.NewImageFromFile("resources/cow-stand.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)
  // log.Printf("cowsheet bounds [%v]", sheet.Bounds())

  // load sheet into array
  for y:=0; y < len(ss.CowStand); y++ {
    for x:=0; x < len(ss.CowStand[y]); x++ {
      rect := image.Rect(x * ss.CowSpriteWidth,y * ss.CowSpriteHeight,(x + 1) * ss.CowSpriteWidth,(y + 1) * ss.CowSpriteHeight)
      ss.CowStand[y][x] = sheet.SubImage(rect).(*ebiten.Image)  
      log.Printf("cow image bounds[%v] width[%v] height[%v]", ss.CowStand[y][x].Bounds(), ss.CowSpriteWidth, ss.CowSpriteHeight)
    }
  }
}
