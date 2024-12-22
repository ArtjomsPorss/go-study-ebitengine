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
  screenWidth = 320
  screenHeight = 240
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
  log.Printf("X[%v] Y[%v] Zone[%v] Px[%v] Py[%v] angle[%v]", g.mouseX, g.mouseY, g.zone, gameLevel.PlayerX, gameLevel.PlayerY, angle)

  updateCharacterState()
  updateCharacterPosition()
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  // draw floor
  drawGameLevel(screen)
 
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
  /*
  floorOp := &ebiten.DrawImageOptions{}
  // TODO Height position could be better maybe - currently it is relative to character
  // position. Maybe character position is calculated incorrectly?
  floorOp.GeoM.Translate(-float64(floorSheet.Width)/2, -float64(floorSheet.Height)/8)
  floorOp.GeoM.Translate(screenWidth/2, screenHeight/2)
  screen.DrawImage(floorSheet.Floor, floorOp)
  */
  for y := 0; y < len(gameLevel.Level); y++ {
    for x := 0; x < len(gameLevel.Level[y]); x++ {
      gameLevelOptions := &ebiten.DrawImageOptions{}
      transX := float64(gameLevel.SpriteWidth * x)
      transY := float64(gameLevel.SpriteHeight / 2 * y)
      // since rombs must fit inbetween
      // the even layer should be shifted to sit inbetween
      if y % 2 == 1 {
        transX -= float64(gameLevel.SpriteWidth/2)
      }
      gameLevelOptions.GeoM.Translate(transX, transY)
      screen.DrawImage(gameLevel.Level[y][x], gameLevelOptions)
    }
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

func updateCharacterPosition() {
  if (characterState == 1) { // running - change the position
    distance := 2.0

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
    // update current position
    gameLevel.PlayerX += deltaX
    gameLevel.PlayerY += deltaY
  }
}
