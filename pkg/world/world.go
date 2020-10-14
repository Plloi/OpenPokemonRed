package world

import (
	"pokered/pkg/data/worldmap/header"
	"pokered/pkg/store"
	"pokered/pkg/util"

	"github.com/hajimehoshi/ebiten"
)

// World data
type World struct {
	MapID  uint
	Image  *ebiten.Image
	Header *header.Header
}

var curWorld *World

// LoadWorldData load world data
func LoadWorldData(id uint) {
	h := header.Get(id)
	img, _ := ebiten.NewImage(int(h.Width*32), int(h.Height*32), ebiten.FilterDefault)
	loadBlockset(h.Tileset)

	for y := 0; y < int(h.Height); y++ {
		for x := 0; x < int(h.Width); x++ {
			blockID := h.Blk(y*int(h.Width) + x)
			block := curBlockset.Data[blockID]
			util.DrawImageBlock(img, block, x, y)
		}
	}
	curWorld = &World{
		MapID:  id,
		Image:  img,
		Header: h,
	}
}

// FrontTileID get tile ID in front of player
func FrontTileID() (uint, uint) {
	p := store.SpriteData[0]
	deltaX, deltaY := 0, 0
	px, py := p.MapXCoord, p.MapYCoord
	switch p.Direction {
	case util.Up:
		py--
		deltaY = -16
	case util.Down:
		py++
		deltaY = 16
	case util.Left:
		px--
		deltaX = -16
	case util.Right:
		px++
		deltaX = 16
	}

	blockX, blockY := (store.SCX+64+deltaX)/32, (store.SCY+64+deltaY)/32
	blockID := curWorld.Header.Blk(blockY*int(curWorld.Header.Width) + blockX)

	switch {
	case px%2 == 0 && py%2 == 0:
		return curBlockset.TilesetID, uint(curBlockset.Bytes[uint(blockID)*16+0])
	case px%2 == 1 && py%2 == 0:
		return curBlockset.TilesetID, uint(curBlockset.Bytes[uint(blockID)*16+2])
	case px%2 == 0 && py%2 == 1:
		return curBlockset.TilesetID, uint(curBlockset.Bytes[uint(blockID)*16+8])
	case px%2 == 1 && py%2 == 1:
		return curBlockset.TilesetID, uint(curBlockset.Bytes[uint(blockID)*16+10])
	}

	return curBlockset.TilesetID, 0
}

// VBlank script executed in VBlank
func VBlank() {
	if curWorld == nil {
		return
	}

	util.DrawImagePixel(store.TileMap, curWorld.Image, -store.SCX, -store.SCY)
}