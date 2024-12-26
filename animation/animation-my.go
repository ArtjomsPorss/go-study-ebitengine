package main

import (
  "math"
  "image"
  _ "image/png"
  "log"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
  runnerImage *ebiten.Image
  standingImage *ebiten.Image
  attackingImage *ebiten.Image
  imageToRender *ebiten.Image

  // character image width and height
  tempHeight int
  tempWidth int

  tempFrameCount int

  floorSheet *SpriteSheetFloor
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

func (g *Game) Update() error {
  g.count++
  g.mouseX, g.mouseY = ebiten.CursorPosition()
  g.zone = calculateZone(g.mouseX, g.mouseY)
  log.Printf("X[%v] Y[%v] Zone[%v] Px[%v] Py[%v] angle[%v] walls[%v]", g.mouseX, g.mouseY, g.zone, gameLevel.PlayerXY.X, gameLevel.PlayerXY.Y, angle, len(floorSheet.WallTopBottom))

  updateCharacterState()
  updateCharacterPosition(gameLevel)
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  // TODO prevent drawing images outside current view
  // draw floor
  drawGameLevel(screen)
  drawWall(screen)
 
  // draw character
  selectImageToDraw(g)
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(tempWidth)/2, -float64(tempHeight)/2)
  op.GeoM.Translate(screenWidth/2, screenHeight/2)
  i := (g.count / 5) % tempFrameCount
  sx, sy := frame0X+i*tempWidth, frame0Y + (g.zone * tempHeight)

  screen.DrawImage(imageToRender.SubImage(image.Rect(sx, sy, sx+tempWidth, sy+tempHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
  return screenWidth, screenHeight
}

func main() {
  loadImages()
  gameLevel = CreateGameLevel()
  ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
  ebiten.SetWindowTitle("Animation Demo")
  if err := ebiten.RunGame(&Game{}); err != nil {
    log.Fatal(err)
  }
}

func calculateZone(x, y int) int {
  // calculate deltas from center
  dx := float64(x - centerX)
  dy := float64(centerY - y) // invert y axis to make upward positive
  // get angle in degrees
  angle = math.Atan2(dy, dx) * (180 / math.Pi)
  // to make South area as value 0 - currently 0 is east-north
  modAngle := angle + 101.25 
  // normalize angle to 0-360
  if modAngle < 0 {
    modAngle += 360 
  } else if modAngle > 360 {
    modAngle -= 360
  }
  modAngle = 360 - modAngle
  // divide circle into 16 regions (22.5 degrees each)
  region := int(math.Floor(modAngle/22.5))
  return region % 16
}

func selectImageToDraw(g *Game) {
  // pick image and cutout size for drawing depending on whether character is standing or running
  if imageToRender == nil || characterState == 0 {
    imageToRender = standingImage
    tempWidth = frameWidthStand
    tempHeight = frameHeightStand
    tempFrameCount = frameCountStand
  } else if characterState == 1 {
    imageToRender = runnerImage
    tempWidth = frameWidth
    tempHeight = frameHeight
    tempFrameCount = frameCount
  } else if characterState == 2 {
    imageToRender = attackingImage
    tempWidth = frameWidthAttack
    tempHeight = frameHeightAttack
    tempFrameCount = frameCountAttack
  }
}

func loadImages() {
  img, _, err := ebitenutil.NewImageFromFile("druid-run.png")
  if err != nil {
    log.Fatal(err)
  }
  runnerImage = ebiten.NewImageFromImage(img)

  img2, _, err := ebitenutil.NewImageFromFile("druid-stand.png")
  if err != nil {
    log.Fatal(err)
  }
  standingImage = ebiten.NewImageFromImage(img2)

  img3, _, err := ebitenutil.NewImageFromFile("druid-attack.png")
  if err != nil {
    log.Fatal(err)
  }
  attackingImage = ebiten.NewImageFromImage(img3)

  ss, err := LoadFloorSpriteSheet()
  if err != nil {
    log.Fatal(err)
  }
  floorSheet = ss
}

func drawGameLevel(screen *ebiten.Image) {
  // TODO repeated calculations can be done only once outside loop
  for y := 0; y < len(gameLevel.Level); y++ {
    for x := 0; x < len(gameLevel.Level[y]); x++ {
      gameLevelOptions := &ebiten.DrawImageOptions{}
      transX := float64(gameLevel.SpriteWidth / 2 * x + gameLevel.SpriteWidth / 2 * y) - gameLevel.PlayerXY.X
      transY := float64(-gameLevel.SpriteHeight / 2 * x + gameLevel.SpriteHeight / 2 * y) - gameLevel.PlayerXY.Y
      
      gameLevelOptions.GeoM.Translate(transX, transY)
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
    transX := float64(x * floorSheet.CliffWidth / 2) 
    transY := float64(x * -floorSheet.Height / 2)

    drawOpts.GeoM.Reset()
    drawOpts.GeoM.Translate(transX - gameLevel.PlayerXY.X, transY - gameLevel.PlayerXY.Y - floorSheet.CliffYDrawStartingPoint)
    screen.DrawImage(floorSheet.WallTopBottom[x], drawOpts)
  }


  // draw BOTTOM walls
  for y := 0; y < maxLen; y++ {
    if y == gameLevelTilesX {
      break
    }
    transX := float64(y * floorSheet.CliffWidth / 2 + gameLevelTilesY * floorSheet.Width / 2) 
    transY := float64(y * -floorSheet.Height / 2 + gameLevelTilesY * floorSheet.Width / 2 )

    drawOpts.GeoM.Reset()
    drawOpts.GeoM.Translate(transX - gameLevel.PlayerXY.X, transY - gameLevel.PlayerXY.Y - floorSheet.CliffYDrawStartingPoint)
    screen.DrawImage(floorSheet.WallTopBottom[y], drawOpts)
  }

  // draw LEFT wall
  for x := 0; x < maxLen; x++ {
     if x == gameLevelTilesY {
      break;
    }
    transX := float64(x * floorSheet.CliffWidth / 2 - floorSheet.Height) 
    transY := float64(x * floorSheet.Height / 2 + floorSheet.CliffHeight - floorSheet.Height)

    drawOpts.GeoM.Reset()
    drawOpts.GeoM.Translate(transX - gameLevel.PlayerXY.X, transY - gameLevel.PlayerXY.Y - floorSheet.CliffYDrawStartingPoint)
    screen.DrawImage(floorSheet.WallLeftRight[x], drawOpts)
  }

  // draw RIGHT wall
  for x := 0; x < maxLen; x++ {
     if x == gameLevelTilesY {
      break;
    }
    transX := float64(x * floorSheet.CliffWidth / 2 - floorSheet.Height + floorSheet.Width / 2 * 8) 
    transY := float64(x * floorSheet.Height / 2 - floorSheet.Height / 2 * 8) - floorSheet.CliffYRightDrawStartingPoint 

    drawOpts.GeoM.Reset()
    drawOpts.GeoM.Translate(transX - gameLevel.PlayerXY.X, transY - gameLevel.PlayerXY.Y)
    screen.DrawImage(floorSheet.WallLeftRight[x], drawOpts)
  }

  // draw CORNERS
    // last wall pieces on all sides replace by walls
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
    if (gl.IsPointInPolygon(p)) {
      gl.PlayerXY = &p
    }
  }
}
