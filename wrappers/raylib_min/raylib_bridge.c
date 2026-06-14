#include "koda.h"
#include "raylib.h"
#include <string.h>

static int knum(KodaValue v) { return IS_NUMBER(v) ? (int)AS_NUMBER(v) : 0; }
static const char* kstr(KodaValue v) { return (IS_OBJ(v) && AS_OBJ(v)->type == OBJ_STRING) ? ((KodaString*)AS_OBJ(v))->chars : ""; }
static Color kcolor(KodaValue v) {
    int c = knum(v);
    Color out = { (unsigned char)((c >> 16) & 255), (unsigned char)((c >> 8) & 255), (unsigned char)(c & 255), 255 };
    return out;
}

KodaValue rl_init_window(int argCount, KodaValue* args) {
    if (argCount < 3) return BOOL_VAL(false);
    InitWindow(knum(args[0]), knum(args[1]), kstr(args[2]));
    return BOOL_VAL(IsWindowReady());
}

KodaValue rl_window_should_close(int argCount, KodaValue* args) {
    return BOOL_VAL(WindowShouldClose());
}

KodaValue rl_close_window(int argCount, KodaValue* args) {
    CloseWindow();
    return NULL_VAL;
}

KodaValue rl_begin_drawing(int argCount, KodaValue* args) {
    BeginDrawing();
    return NULL_VAL;
}

KodaValue rl_end_drawing(int argCount, KodaValue* args) {
    EndDrawing();
    return NULL_VAL;
}

KodaValue rl_clear_background(int argCount, KodaValue* args) {
    if (argCount >= 1) ClearBackground(kcolor(args[0]));
    return NULL_VAL;
}

KodaValue rl_draw_text(int argCount, KodaValue* args) {
    if (argCount >= 5) DrawText(kstr(args[0]), knum(args[1]), knum(args[2]), knum(args[3]), kcolor(args[4]));
    return NULL_VAL;
}

KodaValue rl_draw_rectangle(int argCount, KodaValue* args) {
    if (argCount >= 5) DrawRectangle(knum(args[0]), knum(args[1]), knum(args[2]), knum(args[3]), kcolor(args[4]));
    return NULL_VAL;
}

KodaValue rl_draw_circle(int argCount, KodaValue* args) {
    if (argCount >= 4) DrawCircle(knum(args[0]), knum(args[1]), (float)knum(args[2]), kcolor(args[3]));
    return NULL_VAL;
}

KodaValue rl_set_target_fps(int argCount, KodaValue* args) {
    if (argCount >= 1) SetTargetFPS(knum(args[0]));
    return NULL_VAL;
}

KodaValue rl_is_key_down(int argCount, KodaValue* args) {
    if (argCount < 1) return BOOL_VAL(false);
    return BOOL_VAL(IsKeyDown(knum(args[0])));
}

KodaValue rl_is_key_pressed(int argCount, KodaValue* args) {
    if (argCount < 1) return BOOL_VAL(false);
    return BOOL_VAL(IsKeyPressed(knum(args[0])));
}
