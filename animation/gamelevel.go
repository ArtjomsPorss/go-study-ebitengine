package main

import (
  "log"
  "github.com/hajimehoshi/ebiten/v2"
)

type GameLevel struct {
  // for starters square level
  Level [12][8]*ebiten.Image
  // starting position
  StartX int
  StartY int

  // current character position
  // TODO probably can be moved to character struct?
  // so that each npc can also contain position on level
  PlayerX int
  PlayerY int
  
  SpriteWidth int
  SpriteHeight int
}

func CreateGameLevel() (*GameLevel) {
  spriteSheet, err := LoadFloorSpriteSheet()
  if err != nil {
    log.Fatal(err)
  }
  lvl := &GameLevel{}
  lvl.SpriteWidth = spriteSheet.Width
  lvl.SpriteHeight = spriteSheet.Height
  // create a array of arrays width 8 x 12 down
  // Initialive each cell in the array
  for y := 0; y < len(lvl.Level); y++ {
    for x := 0; x < len(lvl.Level[y]); x++ {
      lvl.Level[y][x] = spriteSheet.Floor
    }
  }

  return lvl
}

