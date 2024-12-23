package main

import (
  "image"
  _ "image/png"
  "log"

  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// spritesheet represents a collection of sprite images.
type SpriteSheetFloor struct {
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
  // wall - cliff - corner
  WallCorner *ebiten.Image
  CornerWidth int
  CornerHeight int
}

// Load Spritesheet loads the embedded spritesheet
func LoadFloorSpriteSheet() (*SpriteSheetFloor, error) {
  img, _, err := ebitenutil.NewImageFromFile("floor-swamp.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)

  floorSheet := &SpriteSheetFloor{}
  floorSheet.Width = 160
  floorSheet.Height = 80
  // load floor image
  floorSheet.Floor = sheet.SubImage(image.Rect(0,0,floorSheet.Width,floorSheet.Height)).(*ebiten.Image)
  // load 
  loadTopBottomCliff(floorSheet)

  return floorSheet, nil
}

func loadTopBottomCliff(spriteSheet *SpriteSheetFloor) {
  img, _, err := ebitenutil.NewImageFromFile("cliff1.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)
  spriteSheet.CliffWidth = 160
  spriteSheet.CliffHeight = 448
  spriteSheet.CliffNarrowWidth = 80 // narrow width is the actual cliff
  spriteSheet.CornerHeight = 448
  spriteSheet.CliffFullHeight = 512
  // load sheet an array
  for i:=0; i < len(spriteSheet.WallTopBottom); i++ {
    spriteSheet.WallTopBottom[i] = sheet.SubImage(image.Rect(i * spriteSheet.CliffWidth,0,(i+1) * spriteSheet.CliffWidth,spriteSheet.CliffFullHeight)).(*ebiten.Image)  
  }

}

