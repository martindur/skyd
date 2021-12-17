package engine

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/google/uuid"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

type Pos struct {
	X, Y int // Works for software rendering, but with OpenGL, probably float32
}

type Player struct {
	Pos
	Radius int
	ID     uuid.UUID
}

type Game struct {
	Local   uuid.UUID
	Players map[uuid.UUID]*Player
	Move    chan Pos
	Host    bool
}

func (player *Player) draw(pixels []byte) {
	for y := -player.Radius; y < player.Radius; y++ {
		for x := -player.Radius; x < player.Radius; x++ {
			if x*x+y*y < player.Radius*player.Radius {
				setPixelWhite(int(player.X)+x, int(player.Y)+y, pixels)
			}
		}
	}
}

func (player *Player) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_A] != 0 {
		player.X--
	}
	if keyState[sdl.SCANCODE_D] != 0 {
		player.X++
	}
}

func (game *Game) Init() (*sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	defer renderer.Destroy()

	return renderer, nil
}

func (game *Game) SpawnPlayer() *Player {
	return &Player{
		Pos:    Pos{X: 100, Y: 100},
		Radius: 25,
		ID:     uuid.New(),
	}
}

// Only need to update local player, as other players
// are only drawn (for now)
func (game *Game) Update(keyState []uint8) {
	game.Players[game.Local].update(keyState)
}

func (game *Game) Render(pixels []byte) {
	for _, player := range game.Players {
		player.draw(pixels)
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixelWhite(x, y int, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = 255
		pixels[index+1] = 255
		pixels[index+2] = 255
	}
}

func makeTexture(r *sdl.Renderer) (*sdl.Texture, error) {
	tex, err := r.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		return nil, err
	}
	defer tex.Destroy()

	return tex, nil
}

func GameLoop(game *Game) {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	// program := initOpenGL()
	initOpenGL()

	shader := NewShader("../engine/shaders/sphere")

	// triangles := []float32{
	// 	0, 0.5, 0,
	// 	-0.5, -0.5, 0,
	// 	0.5, -0.5, 0,
	// 	0.5, 0.5, 0,
	// 	0.25, 0.25, 0,
	// 	0.75, 0.25, 0,
	// }

	// colors := []float32{
	// 	1.0, 0.0, 0.0,
	// 	0.0, 1.0, 0.0,
	// 	0.0, 0.0, 1.0,
	// 	0.0, 0.0, 1.0,
	// 	1.0, 0.0, 0.0,
	// 	0.0, 1.0, 0.0,
	// }

	for !window.ShouldClose() {
		// draw(window, shader.id, &triangles, &colors)
		drawSimple(window, shader.id)
	}
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL Version", version)

	// var defaultShader Shader
	// // Relative to client/server
	// err := defaultShader.load("../engine/shaders/triangle.vert", "../engine/shaders/triangle.frag")
	// if err != nil {
	// 	fmt.Println(os.Getwd())
	// 	fmt.Printf("Error opening: %q", err)
	// }

	// vertexShader, err := compileShader(defaultShader.vertSource, gl.VERTEX_SHADER)
	// if err != nil {
	// 	panic(err)
	// }
	// fragShader, err := compileShader(defaultShader.fragSource, gl.FRAGMENT_SHADER)
	// if err != nil {
	// 	panic(err)
	// }

	// prog := gl.CreateProgram()
	// gl.AttachShader(prog, vertexShader)
	// gl.AttachShader(prog, fragShader)
	// gl.LinkProgram(prog)
	// return prog
	var tmp uint32
	tmp = 5
	return tmp
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(winWidth, winHeight, "Skyd", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// func GameLoop(game *Game) {
// 	pixels := make([]byte, winWidth*winHeight*4)
// 	keyState := sdl.GetKeyboardState()

// 	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
// 		panic(err)
// 	}
// 	defer sdl.Quit()

// 	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
// 		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer window.Destroy()

// 	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer renderer.Destroy()

// 	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer tex.Destroy()

// 	for {
// 		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 			switch event.(type) {
// 			case *sdl.QuitEvent:
// 				return
// 			}
// 		}
// 		clear(pixels)

// 		// Update the local player's position, based on key input
// 		game.Update(keyState)
// 		game.Render(pixels)

// 		tex.Update(nil, pixels, winWidth*4)
// 		renderer.Copy(tex, nil, nil)
// 		renderer.Present()

// 		sdl.Delay(16)
// 		if !game.Host {
// 			game.Move <- game.Players[game.Local].Pos
// 		}
// 	}
// }
