package main

import (
  "math"
  "image"
  _ "image/png"
  "log"
  "github.com/hajimehoshi/ebiten/v2"
  // "github.com/hajimehoshi/ebiten/v2/ebitenutil"
  "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
  screenWidth = 649
  screenHeight = 480
  centerX = screenWidth / 2
  centerY = screenHeight / 2

  frame0X = 0
  frame0Y = 0
  frameWidth = 71
  frameHeight = 84
  frameCount = 8


  frameWidthStand = 43
  frameHeightStand = 78
  frameCountStand = 6
  
  frameWidthAttack = 96
  frameHeightAttack = 85
  frameCountAttack = 16
)

var (
  // character image width and height
  tempHeight int
  tempWidth int

  tempFrameCount int

  floorSheet *SpriteSheet
  gameLevel *GameLevel

  characterState int // 0 - stand, 1 - run, 2 - attack
  angle float64
)

type Game struct {
  count int
  mouseX int
  mouseY int
  zone int
}

type Enemy struct {
  Pos *Point
  State int // 0 - stand, 1 - run, 2 - attack, 3 - dead
  // TODO to remove died sprites after some time - diedAt int
  // direction where it looks
  // TODO can be used also for player?
  Direction int
  Angle int
  Radius float64

  RunImage *ebiten.Image
  StandImage [][]*ebiten.Image
  AttackImage *ebiten.Image
  DeadImage *ebiten.Image
  ImageToRender *ebiten.Image
}

// game engine function that updates state
func (g *Game) Update() error {
  g.count++
  g.mouseX, g.mouseY = ebiten.CursorPosition()
 // log.Printf("X[%v] Y[%v] Zone[%v] Px[%v] Py[%v] angle[%v] walls[%v]", g.mouseX, g.mouseY, g.zone, gameLevel.PlayerXY.X, gameLevel.PlayerXY.Y, angle, len(floorSheet.WallTopBottom))
  playerFeetY := centerY + floorSheet.PlayerYPos
  g.zone = calculateZone(centerX, playerFeetY, g.mouseX, g.mouseY, 16)
  // log.Printf("X[%v] Y[%v] Zone[%v] Px[%v] Py[%v] angle[%v]", g.mouseX, g.mouseY, g.zone, gameLevel.PlayerXY.X, gameLevel.PlayerXY.Y + float64(playerFeetY), angle)

 
  updateCharacterState()
  updateCharacterPosition(gameLevel)
  updateCowState()
  return nil
}

// game engine function that draws everything
func (g *Game) Draw(screen *ebiten.Image) {
  // TODO prevent drawing images outside current view
  // draw floor
  drawGameLevel(screen)
  drawWall(screen)
 
  // draw character
  selectImageToDraw(g)
  // draw player
  drawPlayer(g, screen)
  // draw cows
  drawCows(g, screen)
}

func drawPlayer(g *Game, screen *ebiten.Image) {

  op := &ebiten.DrawImageOptions{}
  // adjust for half character width and height
  op.GeoM.Translate(-float64(tempWidth / 2), -float64(tempHeight / 2))
  // adjust to draw at the center of the screen
  op.GeoM.Translate(screenWidth/2, screenHeight/2)
  // select image based on ingame timer and frame
  i := (g.count / 5) % tempFrameCount
  // and direction of where the cursor is pointing from player
  sx, sy := frame0X+i*tempWidth, frame0Y + (g.zone * tempHeight)

  screen.DrawImage(floorSheet.ImageToRender.SubImage(image.Rect(sx, sy, sx+tempWidth, sy+tempHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
  return screenWidth, screenHeight
}

func main() {
  loadImages()
  gameLevel = CreateGameLevel(floorSheet)
  ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
  ebiten.SetWindowTitle("Animation Demo")
  if err := ebiten.RunGame(&Game{}); err != nil {
    log.Fatal(err)
  }
}

func calculateZone(fx, fy, tx, ty, zoneCount int) int {
  // calculate deltas from center
  dx := float64(tx - fx)
  dy := float64(fy - ty) // invert y axis to make upward positive
  // get angle in degrees
  angle = math.Atan2(dy, dx) * (180 / math.Pi)
  // to make South area as value 0 - currently 0 is east-north
  modAngle := angle + 75
  // modAngle := angle 
  // normalize angle to 0-360
  if modAngle < 0 {
    modAngle += 360 
  } else if modAngle > 360 {
    modAngle -= 360
  }
  modAngle = 360 - modAngle
  // divide circle into 16 regions (22.5 degrees each)
  regions := float64(360 / zoneCount)
  region := int(math.Floor(modAngle/regions))
  return region % zoneCount
}

func selectImageToDraw(g *Game) {
  // pick image and cutout size for drawing depending on whether character is standing or running
  if floorSheet.ImageToRender == nil || characterState == 0 {
    floorSheet.ImageToRender = floorSheet.StandingImage
    tempWidth = frameWidthStand
    tempHeight = frameHeightStand
    tempFrameCount = frameCountStand
  } else if characterState == 1 {
    floorSheet.ImageToRender = floorSheet.RunnerImage
    tempWidth = frameWidth
    tempHeight = frameHeight
    tempFrameCount = frameCount
  } else if characterState == 2 {
    floorSheet.ImageToRender = floorSheet.AttackingImage
    tempWidth = frameWidthAttack
    tempHeight = frameHeightAttack
    tempFrameCount = frameCountAttack
  }
}

func loadImages() {
  floorSheet = LoadSprites()
}

func drawGameLevel(screen *ebiten.Image) {
  // TODO repeated calculations can be done only once outside loop
  for y := 0; y < len(gameLevel.Level); y++ {
    for x := 0; x < len(gameLevel.Level[y]); x++ {
      gameLevelOptions := &ebiten.DrawImageOptions{}

      // adjust for isometric positioning column/row
      transX := float64(gameLevel.SpriteWidth / 2 * x + gameLevel.SpriteWidth / 2 * y) 
      transY := float64(-gameLevel.SpriteHeight / 2 * x + gameLevel.SpriteHeight / 2 * y) 
      gameLevelOptions.GeoM.Translate(transX, transY)
      // adjust for half of floor tile height
      gameLevelOptions.GeoM.Translate(0, -float64(gameLevel.SpriteHeight / 2))
      // adjust for player position
      gameLevelOptions.GeoM.Translate(float64(-gameLevel.PlayerXY.X),float64(-gameLevel.PlayerXY.Y))
      // adjust for player screen positioning
      gameLevelOptions.GeoM.Translate(screenWidth/2, screenHeight/2)
      // adjust for player feet position (39p lower)
      gameLevelOptions.GeoM.Translate(float64(0), float64(floorSheet.PlayerYPos))
      // ?? adjust for character width ??

      screen.DrawImage(gameLevel.Level[y][x], gameLevelOptions)
    }
  }
}

func drawWall(screen *ebiten.Image) {
  // TODO repeated calculations can be done only once outside loop
  maxLen := len(floorSheet.WallTopBottom)
  gameLevelTilesX := len(gameLevel.Level[0])
  gameLevelTilesY := len(gameLevel.Level)
  drawOpts := &ebiten.DrawImageOptions{}

  // draw TOP walls
  for x := 0; x < maxLen; x++ {
    if x == gameLevelTilesX {
      break
    }
    drawOpts.GeoM.Reset()

    transX := float64(x * floorSheet.CliffWidth / 2) 
    transY := float64(x * -floorSheet.Height / 2)
    drawOpts.GeoM.Translate(transX, transY - floorSheet.CliffYDrawStartingPoint)
    // adjust for half of floor tile height
    drawOpts.GeoM.Translate(0, -float64(gameLevel.SpriteHeight / 2))
    // adjust for player position
    drawOpts.GeoM.Translate(float64(-gameLevel.PlayerXY.X),float64(-gameLevel.PlayerXY.Y))
    // adjust for player screen positioning
    drawOpts.GeoM.Translate(screenWidth/2, screenHeight/2)
    // adjust for player feet position (39p lower)
    drawOpts.GeoM.Translate(float64(0), float64(floorSheet.PlayerYPos))
 
    screen.DrawImage(floorSheet.WallTopBottom[x], drawOpts)
  }


  // draw BOTTOM walls
  for y := 0; y < maxLen; y++ {
    if y == gameLevelTilesX {
      break
    }
    drawOpts.GeoM.Reset()
  
    // adjust for half of floor tile height
    drawOpts.GeoM.Translate(0, -float64(gameLevel.SpriteHeight / 2))
    // adjust for player position
    drawOpts.GeoM.Translate(float64(-gameLevel.PlayerXY.X),float64(-gameLevel.PlayerXY.Y))
    // adjust for player screen positioning
    drawOpts.GeoM.Translate(screenWidth/2, screenHeight/2)
    // adjust for player feet position (39p lower)
    drawOpts.GeoM.Translate(float64(0), float64(floorSheet.PlayerYPos))
 
    // adjust for tile position
    transX := float64(y * floorSheet.CliffWidth / 2 + gameLevelTilesY * floorSheet.Width / 2) 
    transY := float64(y * -floorSheet.Height / 2 + gameLevelTilesY * floorSheet.Width / 2 )

    drawOpts.GeoM.Translate(transX, transY - floorSheet.CliffYDrawStartingPoint)

    screen.DrawImage(floorSheet.WallTopBottom[y], drawOpts)
  }

  // draw LEFT wall
  for x := 0; x < maxLen; x++ {
     if x == gameLevelTilesY {
      break;
    }
    drawOpts.GeoM.Reset()
 
     // adjust for half of floor tile height
    drawOpts.GeoM.Translate(0, -float64(gameLevel.SpriteHeight / 2))
    // adjust for player position
    drawOpts.GeoM.Translate(float64(-gameLevel.PlayerXY.X),float64(-gameLevel.PlayerXY.Y))
    // adjust for player screen positioning
    drawOpts.GeoM.Translate(screenWidth/2, screenHeight/2)
    // adjust for player feet position (39p lower)
    drawOpts.GeoM.Translate(float64(0), float64(floorSheet.PlayerYPos))

    // adjust for tile size and positioning
    transX := float64(x * floorSheet.CliffWidth / 2 - floorSheet.Height) 
    transY := float64(x * floorSheet.Height / 2 + floorSheet.CliffHeight - floorSheet.Height)
    drawOpts.GeoM.Translate(transX, transY - floorSheet.CliffYDrawStartingPoint)

    screen.DrawImage(floorSheet.WallLeftRight[x], drawOpts)
  }

  // draw RIGHT wall
  for x := 0; x < maxLen; x++ {
     if x == gameLevelTilesY {
      break;
    }
    drawOpts.GeoM.Reset()

     // adjust for half of floor tile height
    drawOpts.GeoM.Translate(0, -float64(gameLevel.SpriteHeight / 2))
    // adjust for player position
    drawOpts.GeoM.Translate(float64(-gameLevel.PlayerXY.X),float64(-gameLevel.PlayerXY.Y))
    // adjust for player screen positioning
    drawOpts.GeoM.Translate(screenWidth/2, screenHeight/2)
    // adjust for player feet position (39p lower)
    drawOpts.GeoM.Translate(float64(0), float64(floorSheet.PlayerYPos))

    // adjust for tile size and positioning
    transX := float64(x * floorSheet.CliffWidth / 2 - floorSheet.Height + floorSheet.Width / 2 * 8) 
    transY := float64(x * floorSheet.Height / 2 - floorSheet.Height / 2 * 8) - floorSheet.CliffYRightDrawStartingPoint 
    drawOpts.GeoM.Translate(transX, transY)

    screen.DrawImage(floorSheet.WallLeftRight[x], drawOpts)
  }

  // TODO draw CORNERS
    // last wall pieces on all sides replace by walls
  
  // TODO skip drawing sprites outside user window - or it is done automatically?
}

func drawCows(g *Game, screen *ebiten.Image) {
  zone := selectCowImageToDraw(g)

  // log.Printf("zone[%v]", zone)

  divider := len(floorSheet.CowToRender[zone])
  // divide by 5 makes it slow enough
  i := g.count / 5 % divider

  op := &ebiten.DrawImageOptions{}
  // adjust for player XY
  op.GeoM.Translate(gameLevel.Enemies[0].Pos.X - gameLevel.PlayerXY.X, gameLevel.Enemies[0].Pos.Y - gameLevel.PlayerXY.Y)
  // adjust for where player is rendered on center of screen
  op.GeoM.Translate(screenWidth/2, screenHeight/2)
  // adjust for where cow's standing center mass at bottom of legs are on a sprite
  xy := floorSheet.CowToRender[0][i].Bounds().Size()
  op.GeoM.Translate(-float64(xy.X/2), -float64(floorSheet.CowYPos))
  // adjust for player feet position (39p lower)
  op.GeoM.Translate(float64(0), float64(floorSheet.PlayerYPos))

  screen.DrawImage(floorSheet.CowToRender[zone][i], op)
}

func selectCowImageToDraw(g *Game) int {
  ex := gameLevel.Enemies[0].Pos.X
  ey := gameLevel.Enemies[0].Pos.Y
  // playerFeetY := centerY + floorSheet.PlayerYPos
  // px := centerX
  // py := playerFeetY
  px := gameLevel.PlayerXY.X
  py := gameLevel.PlayerXY.Y

  // find where the player is
  // find angle towards player and which zone it is
  zone := calculateZone(int(ex), int(ey), int(px), int(py), 8)

	log.Printf("cow xy[%v,%v] pl xy[%v,%v] zone[%v]", ex,ey, px, py, zone)
  // get the area group
  // TODO set the proper image for rendering
	return zone
}

func updateCowState() {
  // cow radius is larger than player's
  // is cow point + movement distance distance to player point is less than cow radius + player radius
  // then move cow towards player
  cPos := gameLevel.Enemies[0].Pos
  pPos := gameLevel.PlayerXY
  // calculate distance between 2 positions
  dx := cPos.X - pPos.X
  dy := cPos.Y - pPos.Y
  d := math.Sqrt(dx * dx + dy * dy)
  cowMovementDistance := 2.0
  if d + cowMovementDistance > gameLevel.Enemies[0].Radius + gameLevel.PlayerRadius {
    // set cow as moving
    gameLevel.Enemies[0].State = 1
    floorSheet.CowToRender = floorSheet.CowWalk
    // move towards player

    // calculate angle from cow to player
    dx = pPos.X - cPos.X
    dy = cPos.Y - pPos.Y
 
    modAngle := math.Atan2(dy, dx) * (180 / math.Pi)
    if modAngle < 0 {
      modAngle += 360 
    } else if modAngle > 360 {
      modAngle -= 360
    }
    modAngle = 360 - modAngle

    angleRadians := modAngle * math.Pi / 180
    deltaX := math.Cos(angleRadians) * cowMovementDistance
    deltaY := math.Sin(angleRadians) * cowMovementDistance
    // new position
    p := Point{cPos.X + deltaX, cPos.Y + deltaY}
    gameLevel.Enemies[0].Pos = &p
  } else {
    // cow is attacking
    gameLevel.Enemies[0].State = 2
    floorSheet.CowToRender = floorSheet.CowAttack
  }
}

func updateCharacterState() {
  if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
    characterState = 0 // player standing
  } else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
    characterState = 1 // player running
  } else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
    characterState = 2 // player attacking 
  }
}

func updateCharacterPosition(gl *GameLevel) {
  if (characterState == 1) { // running - change the position
    distance := 3.0

    modAngle := angle
    if modAngle < 0 {
      modAngle += 360 
    } else if modAngle > 360 {
      modAngle -= 360
    }
    modAngle = 360 - modAngle

    angleRadians := modAngle * math.Pi / 180
    deltaX := math.Cos(angleRadians) * distance
    deltaY := math.Sin(angleRadians) * distance
    p := Point{gameLevel.PlayerXY.X + deltaX, gameLevel.PlayerXY.Y + deltaY}

    // if new position will be within the level
    // update current position
    // TODO uncomment border
    if (gl.IsPointInPolygon(p)) {
      gl.PlayerXY = &p
    }
  }
}
