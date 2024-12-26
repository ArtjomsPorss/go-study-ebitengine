- animation
  - [x] read image
  - [x] remove background from druid image
  - [x] draw running image
  - [x] standing image
  - [x] switch from standing to running based on mouse left click
  - [x] attacking image
  - [x] switch on attacking image if right mouse clicked
  - [ ] (later) complete running and attacking animation when canceled
- hue / saturation / value - where spells are cast
  - in front - apply light and saturation over all character when spell in front of them
  - behind - at the edges of the character
- shader looks interesting - maybe can be used as a replacement of hue changes
- drag and drop - inventory
- isometric - the game itself - background, walls
  - [x] draw background
    - [x] select romboid as image
    - [x] draw a level of floor (outside screen also)
  - [x] move around level when walking
    - [x] level shold be drawn not in zig zag but rather each lower tile is shifted right
  - [x] draw walls
  - [x] walls prevent from moving forward - collision detection
  - [ ] [shader/hue] walls disappear when wall in front of a character 
- masking - draw a hole in the wall that is drawn in front of the character
- battle
  - [ ] draw an enemy
  - [ ] draw several enemies
  - [ ] enemy moves towards character
  - [ ] update closest enemy first
    - [ ] (later) enemy shouldn't move out of bounds
  - [ ] when enemy is close - attack
  - [ ] attacking animation
  - [ ] draw life
