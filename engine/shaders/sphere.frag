#version 330 core

out vec4 frag_color;

uniform float uTime;
uniform vec3 ourColor;

float sphereSDF(vec2 p, vec2 c, float r){
    return length(p - c) - r;
}


vec3 raymarch(in vec2 pos) {
    float sphere = sphereSDF(pos, vec2(0.0, 0.0), 0.25);

    if (sphere <= 0.0) {
        return ourColor;
    }

    if (sphere > 0.0) {
        return vec3(0.0, 0.0, 0.0);
    }
}

const vec2 resolution = vec2(800.0, 600.0);

void main() {
    vec2 p = (gl_FragCoord.xy - resolution) / resolution.y;
    vec2 center = vec2(0.0, sin(uTime));
    float r = 0.25;
    vec2 pos = vec2(0.0, (-1.0 + r));

    float wind = abs(p.y - pos.y) * sin(uTime);
    pos.x = pos.x + wind;
    // vec3 col = raymarch(p);
    float sdf = sphereSDF(p, pos, r);
    // float col = step(0.0, sphereSDF(p, vec2(0.0), 0.25));
    float col = smoothstep(-0.01, 0.01, sdf);
    // frag_color = vec4(val, val, val, 1.0);
    frag_color = vec4(vec3(col), 1.0);
}
