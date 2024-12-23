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
  Floor *ebiten.Image
  WallLeft *ebiten.Image
  WallRight *ebiten.Image
  WallTop *ebiten.Image
  WallBottom *ebiten.Image
  Width int
  Height int
  
  CliffNarrowWidth int
  CliffWidth int
  CliffHeight int

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
  loadTopBottomCliff()

  return floorSheet, nil
}

func loadTopBottomCliff(spriteSheet *SpriteSheetFloor) {
  img, _, err := ebitenutil.NewImageFromFile("cliff1.png")
  if err != nil {
    log.Fatal(err)
  }
  sheet := ebiten.NewImageFromImage(img)
  sheet.CliffNarrowWidth = 80 // narrow width is the actual cliff
  sheet.CornerHeight = 448
  distanceFromTop := 512 - sheet.CornerHeight
  // load sheet an array
}

