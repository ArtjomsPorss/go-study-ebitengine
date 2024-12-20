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
  Width int
  Height int
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
  // width = 160
  // height = 80
  floorSheet.Floor = sheet.SubImage(image.Rect(0,0,floorSheet.Width,floorSheet.Height)).(*ebiten.Image)

  return floorSheet, nil
}

