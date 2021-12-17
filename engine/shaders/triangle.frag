#version 330 core

in vec3 col;
out vec4 frag_color;

uniform vec3 ourColor;

void main() {
    // frag_color = vec4(ourColor, 1.0);
    frag_color = vec4(1.0, 1.0, 1.0, 1.0);
}