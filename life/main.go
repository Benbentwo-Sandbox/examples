package life

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/fyne-io/examples/img/icon"
)

type board struct {
	cells  [][]bool
	width  int
	height int
}

func (b *board) ifAlive(x, y int) int {
	if x < 0 || x >= b.width {
		return 0
	}

	if y < 0 || y >= b.height {
		return 0
	}

	if b.cells[y][x] {
		return 1
	}
	return 0
}

func (b *board) countNeighbours(x, y int) int {
	sum := 0

	sum += b.ifAlive(x-1, y-1)
	sum += b.ifAlive(x, y-1)
	sum += b.ifAlive(x+1, y-1)

	sum += b.ifAlive(x-1, y)
	sum += b.ifAlive(x+1, y)

	sum += b.ifAlive(x-1, y+1)
	sum += b.ifAlive(x, y+1)
	sum += b.ifAlive(x+1, y+1)

	return sum
}

func (b *board) nextGen() [][]bool {
	state := make([][]bool, b.height)

	for y := 0; y < b.height; y++ {
		state[y] = make([]bool, b.width)

		for x := 0; x < b.width; x++ {
			n := b.countNeighbours(x, y)

			if b.cells[y][x] {
				state[y][x] = n == 2 || n == 3
			} else {
				state[y][x] = n == 3
			}
		}
	}

	return state
}

func (b *board) renderState(state [][]bool) {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x] = state[y][x]
		}
	}
}

func (b *board) load() {
	// gun
	b.cells[5][1] = true
	b.cells[5][2] = true
	b.cells[6][1] = true
	b.cells[6][2] = true

	b.cells[3][13] = true
	b.cells[3][14] = true
	b.cells[4][12] = true
	b.cells[4][16] = true
	b.cells[5][11] = true
	b.cells[5][17] = true
	b.cells[6][11] = true
	b.cells[6][15] = true
	b.cells[6][17] = true
	b.cells[6][18] = true
	b.cells[7][11] = true
	b.cells[7][17] = true
	b.cells[8][12] = true
	b.cells[8][16] = true
	b.cells[9][13] = true
	b.cells[9][14] = true

	b.cells[1][25] = true
	b.cells[2][23] = true
	b.cells[2][25] = true
	b.cells[3][21] = true
	b.cells[3][22] = true
	b.cells[4][21] = true
	b.cells[4][22] = true
	b.cells[5][21] = true
	b.cells[5][22] = true
	b.cells[6][23] = true
	b.cells[6][25] = true
	b.cells[7][25] = true

	b.cells[3][35] = true
	b.cells[3][36] = true
	b.cells[4][35] = true
	b.cells[4][36] = true

	// spaceship
	b.cells[34][2] = true
	b.cells[34][3] = true
	b.cells[34][4] = true
	b.cells[34][5] = true
	b.cells[35][1] = true
	b.cells[35][5] = true
	b.cells[36][5] = true
	b.cells[37][1] = true
	b.cells[37][4] = true
}

func newBoard() *board {
	b := &board{nil, 60, 50}
	b.cells = make([][]bool, b.height)

	for y := 0; y < b.height; y++ {
		b.cells[y] = make([]bool, b.width)
	}

	return b
}

type game struct {
	board  *board
	paused bool

	size     fyne.Size
	position fyne.Position
	hidden   bool
}

func (g *game) Size() fyne.Size {
	return g.size
}

func (g *game) Resize(size fyne.Size) {
	g.size = size
	widget.Renderer(g).Layout(size)
}

func (g *game) Position() fyne.Position {
	return g.position
}

func (g *game) Move(pos fyne.Position) {
	g.position = pos
	widget.Renderer(g).Layout(g.size)
}

func (g *game) MinSize() fyne.Size {
	return widget.Renderer(g).MinSize()
}

func (g *game) Visible() bool {
	return !g.hidden
}

func (g *game) Show() {
	g.hidden = false
}

func (g *game) Hide() {
	g.hidden = true
}

type gameRenderer struct {
	render   *canvas.Raster
	objects  []fyne.CanvasObject
	imgCache *image.RGBA

	aliveColor color.Color
	deadColor  color.Color

	game *game
}

func (g *gameRenderer) MinSize() fyne.Size {
	return fyne.NewSize(g.game.board.width*10, g.game.board.height*10)
}

func (g *gameRenderer) Layout(size fyne.Size) {
	g.render.Resize(size)
}

func (g *gameRenderer) ApplyTheme() {
	g.aliveColor = theme.TextColor()
	g.deadColor = theme.BackgroundColor()
}

func (g *gameRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (g *gameRenderer) Refresh() {
	canvas.Refresh(g.render)
}

func (g *gameRenderer) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *gameRenderer) Destroy() {
}

func (g *gameRenderer) draw(w, h int) image.Image {
	img := g.imgCache
	if img == nil || img.Bounds().Size().X != w || img.Bounds().Size().Y != h {
		img = image.NewRGBA(image.Rect(0, 0, w, h))
		g.imgCache = img
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			xpos, ypos := g.game.cellForCoord(x, y, w, h)

			if xpos < g.game.board.width && ypos < g.game.board.height && g.game.board.cells[ypos][xpos] {
				img.Set(x, y, g.aliveColor)
			} else {
				img.Set(x, y, g.deadColor)
			}
		}
	}

	return img
}

func (g *game) CreateRenderer() fyne.WidgetRenderer {
	renderer := &gameRenderer{game: g}

	render := canvas.NewRaster(renderer.draw)
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	renderer.ApplyTheme()

	return renderer
}

func (g *game) cellForCoord(x, y, w, h int) (int, int) {
	xpos := int(float64(g.board.width) * (float64(x) / float64(w)))
	ypos := int(float64(g.board.height) * (float64(y) / float64(h)))

	return xpos, ypos
}

func (g *game) run() {
	g.paused = false
}

func (g *game) stop() {
	g.paused = true
}

func (g *game) toggleRun() {
	g.paused = !g.paused
}

func (g *game) animate() {
	go func() {
		tick := time.NewTicker(time.Second / 6)

		for {
			select {
			case <-tick.C:
				if g.paused {
					continue
				}

				state := g.board.nextGen()
				g.board.renderState(state)
				widget.Refresh(g)
			}
		}
	}()
}

func (g *game) typedRune(r rune) {
	if r == ' ' {
		g.toggleRun()
	}
}

func (g *game) Tapped(ev *fyne.PointEvent) {
	xpos, ypos := g.cellForCoord(ev.Position.X, ev.Position.Y, g.size.Width, g.size.Height)

	if ev.Position.X < 0 || ev.Position.Y < 0 || xpos >= g.board.width || ypos >= g.board.height {
		return
	}

	g.board.cells[ypos][xpos] = !g.board.cells[ypos][xpos]

	widget.Refresh(g)
}

func (g *game) TappedSecondary(ev *fyne.PointEvent) {
}

func newGame(b *board) *game {
	g := &game{board: b}

	return g
}

// Show starts a new game of life
func Show(app fyne.App) {
	board := newBoard()
	board.load()

	game := newGame(board)

	window := app.NewWindow("Life")
	window.SetIcon(icon.LifeBitmap)
	window.SetContent(game)
	window.Canvas().SetOnTypedRune(game.typedRune)

	// start the board animation before we show the window - it will block
	game.animate()

	window.Show()
}
