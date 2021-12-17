package engine

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Shader struct {
	id         uint32
	vertSource string
	fragSource string
}

func NewShader(shaderPath string) *Shader {
	shader := new(Shader)
	err := shader.load(shaderPath+".vert", shaderPath+".frag")
	if err != nil {
		fmt.Printf("Error opening: %q", err)
	}

	vertex, err := compileShader(shader.vertSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	frag, err := compileShader(shader.fragSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shader.id = gl.CreateProgram()
	gl.AttachShader(shader.id, vertex)
	gl.AttachShader(shader.id, frag)
	gl.LinkProgram(shader.id)

	// Should check if linking issues here

	// Shaders are now linked to program, so we can delete them.
	gl.DeleteShader(vertex)
	gl.DeleteShader(frag)

	return shader
}

func (s *Shader) use() {
	gl.UseProgram(s.id)
}

func (s *Shader) setFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.id, gl.Str(name+"\x00")), value)
}

func (s *Shader) load(vertPath string, fragPath string) error {
	vertFile, errV := os.Open(vertPath)
	defer vertFile.Close()
	if errV != nil {
		return errV
	}
	fragFile, errF := os.Open(fragPath)
	defer fragFile.Close()
	if errF != nil {
		return errF
	}

	vertBytes, errV := ioutil.ReadAll(vertFile)
	if errV != nil {
		return errV
	}
	fragBytes, errF := ioutil.ReadAll(fragFile)
	if errF != nil {
		return errF
	}

	s.vertSource = string(vertBytes) + "\x00"
	s.fragSource = string(fragBytes) + "\x00"

	return nil
}

func drawSimple(window *glfw.Window, prog uint32) {
	// This draw function relies on the Vertex Shader
	// to specify the placement of 4 vertices (The entire screen)
	quad := []float32{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, 1.0, 0.0,
		1.0, 1.0, 0.0,
	}

	var vbo uint32
	var vao uint32

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(quad), gl.Ptr(quad), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	timeVal := glfw.GetTime()
	greenVal := (math.Sin(timeVal) / 2.0) + 0.5
	vertexColLoc := gl.GetUniformLocation(prog, gl.Str("ourColor\x00"))
	timeUniform := gl.GetUniformLocation(prog, gl.Str("uTime\x00"))

	gl.UseProgram(prog)

	gl.Uniform3f(vertexColLoc, 0.0, float32(greenVal), 0.0)
	gl.Uniform1f(timeUniform, float32(timeVal))

	glfw.PollEvents()
	window.SwapBuffers()
}

func draw(window *glfw.Window, prog uint32, vertices *[]float32, vcolors *[]float32) {

	var vbo uint32
	var vcbo uint32
	var vao uint32

	gl.GenBuffers(1, &vbo) // Seems possible to provide array for multiple objects in C
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(*vertices), gl.Ptr(*vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &vcbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vcbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(*vcolors), gl.Ptr(*vcolors), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, vcbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	timeVal := glfw.GetTime()
	greenVal := (math.Sin(timeVal) / 2.0) + 0.5
	vertexColLoc := gl.GetUniformLocation(prog, gl.Str("ourColor\x00"))

	gl.UseProgram(prog)
	gl.Uniform3f(vertexColLoc, 0.0, float32(greenVal), 0.0)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(*vertices)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var LogLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &LogLength)

		log := strings.Repeat("\x00", int(LogLength+1))
		gl.GetShaderInfoLog(shader, LogLength, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// initialises a Vertex Array Object from the provided points
func makeVAO(points []float32) uint32 {
	var vbo uint32 // Vertex Buffer Object
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
