package main

import (
  // "log"
  "github.com/hajimehoshi/ebiten/v2"
)

type Point struct {
  X, Y float64
}

type GameLevel struct {
  // for starters square level
  Level [8][8]*ebiten.Image
  LevelCoords [4]*Point
  // starting position
  StartX int
  StartY int

  // current character position
  // TODO probably can be moved to character struct?
  // so that each npc can also contain position on level
  PlayerXY *Point
  
  SpriteWidth int
  SpriteHeight int
  
  Enemies [1]*Enemy
}

func CreateGameLevel(spriteSheet *SpriteSheet) (*GameLevel) {
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
/*
  lvl.LevelCoords[0] = &Point{-300,-200}
  lvl.LevelCoords[1] = &Point{320,-510}
  lvl.LevelCoords[2] = &Point{910,-230}
  lvl.LevelCoords[3] = &Point{300,80}
*/ 
// testing coordinates
  lvl.LevelCoords[0] = &Point{380,0}
  lvl.LevelCoords[1] = &Point{460,-80}
  lvl.LevelCoords[2] = &Point{540,0}
  lvl.LevelCoords[3] = &Point{460,80}
  lvl.PlayerXY = &Point{400,0}

  lvl.Enemies[0] = &Enemy{}
  lvl.Enemies[0].Pos = &Point{400,0}

  return lvl
}

// IsPointInPolygon determines if a point is inside a polygon
func (lvl GameLevel) IsPointInPolygon(point Point) bool {
	inside := false
	n := len(lvl.LevelCoords)
	for i := 0; i < n; i++ {
		// Current vertex and the next vertex in the polygon
		current := lvl.LevelCoords[i]
		next := lvl.LevelCoords[(i+1)%n] // Wrap around to the first point after the last point

		// Check if the test point is inside the polygon edge
		if ((current.Y > point.Y) != (next.Y > point.Y)) &&
			(point.X < (next.X-current.X)*(point.Y-current.Y)/(next.Y-current.Y)+current.X) {
			inside = !inside
		}
	}

	return inside
}


