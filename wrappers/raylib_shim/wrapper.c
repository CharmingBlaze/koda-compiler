/*
 * Minimal Raylib ↔ Koda shim (Value native ABI).
 * Compiled with koda build via KODA_NATIVE_SOURCES.
 */
#include "koda_runtime.h"
#include <raylib.h>
#include <string.h>

static const char* koda_arg_cstr(Value* args, int i) {
    if (!args || !IS_OBJ(args[i]) || AS_OBJ(args[i])->type != OBJ_STRING) {
        return "";
    }
    return ((ObjString*)AS_OBJ(args[i]))->chars;
}

/* Koda rgb()/rgba() pack as 0xRRGGBBAA; Raylib Color is {r,g,b,a} bytes — do not cast the uint32. */
static Color koda_color_from_packed(unsigned int packed) {
    Color c;
    c.r = (unsigned char)((packed >> 24) & 255);
    c.g = (unsigned char)((packed >> 16) & 255);
    c.b = (unsigned char)((packed >> 8) & 255);
    c.a = (unsigned char)(packed & 255);
    return c;
}

Value koda_shim_InitWindow(int argCount, Value* args) {
    if (argCount < 3) return NIL_VAL;
    int w = (int)(IS_NUMBER(args[0]) ? AS_NUMBER(args[0]) : 0);
    int h = (int)(IS_NUMBER(args[1]) ? AS_NUMBER(args[1]) : 0);
    const char* title = koda_arg_cstr(args, 2);
    InitWindow(w, h, title);
    return NIL_VAL;
}

Value koda_shim_CloseWindow(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    CloseWindow();
    return NIL_VAL;
}

Value koda_shim_WindowShouldClose(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return BOOL_VAL(WindowShouldClose());
}

Value koda_shim_SetTargetFPS(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    int fps = (int)(IS_NUMBER(args[0]) ? AS_NUMBER(args[0]) : 0);
    SetTargetFPS(fps);
    return NIL_VAL;
}

Value koda_shim_BeginDrawing(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    BeginDrawing();
    return NIL_VAL;
}

Value koda_shim_EndDrawing(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    EndDrawing();
    return NIL_VAL;
}

Value koda_shim_ClearBackground(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    unsigned int c = (unsigned int)(IS_NUMBER(args[0]) ? AS_NUMBER(args[0]) : 0);
    ClearBackground(koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_BeginMode3D(int argCount, Value* args) {
    if (argCount < 10) return NIL_VAL;
    Camera3D cam = {0};
    cam.position.x = (float)AS_NUMBER(args[0]);
    cam.position.y = (float)AS_NUMBER(args[1]);
    cam.position.z = (float)AS_NUMBER(args[2]);
    cam.target.x = (float)AS_NUMBER(args[3]);
    cam.target.y = (float)AS_NUMBER(args[4]);
    cam.target.z = (float)AS_NUMBER(args[5]);
    cam.up.x = (float)AS_NUMBER(args[6]);
    cam.up.y = (float)AS_NUMBER(args[7]);
    cam.up.z = (float)AS_NUMBER(args[8]);
    cam.fovy = (float)AS_NUMBER(args[9]);
    cam.projection = CAMERA_PERSPECTIVE;
    BeginMode3D(cam);
    return NIL_VAL;
}

Value koda_shim_EndMode3D(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    EndMode3D();
    return NIL_VAL;
}

Value koda_shim_DrawGrid(int argCount, Value* args) {
    if (argCount < 2) return NIL_VAL;
    int slices = (int)AS_NUMBER(args[0]);
    float spacing = (float)AS_NUMBER(args[1]);
    DrawGrid(slices, spacing);
    return NIL_VAL;
}

Value koda_shim_DrawCube(int argCount, Value* args) {
    if (argCount < 7) return NIL_VAL;
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    Vector3 size = {(float)AS_NUMBER(args[3]), (float)AS_NUMBER(args[4]), (float)AS_NUMBER(args[5])};
    unsigned int c = (unsigned int)AS_NUMBER(args[6]);
    DrawCubeV(pos, size, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_DrawCubeWires(int argCount, Value* args) {
    if (argCount < 7) return NIL_VAL;
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    Vector3 size = {(float)AS_NUMBER(args[3]), (float)AS_NUMBER(args[4]), (float)AS_NUMBER(args[5])};
    unsigned int c = (unsigned int)AS_NUMBER(args[6]);
    DrawCubeWiresV(pos, size, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_DrawText(int argCount, Value* args) {
    if (argCount < 5) return NIL_VAL;
    const char* text = koda_arg_cstr(args, 0);
    int x = (int)AS_NUMBER(args[1]);
    int y = (int)AS_NUMBER(args[2]);
    int fs = (int)AS_NUMBER(args[3]);
    unsigned int c = (unsigned int)AS_NUMBER(args[4]);
    DrawText(text, x, y, fs, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_GetFrameTime(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetFrameTime());
}

Value koda_shim_IsKeyDown(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    int key = (int)AS_NUMBER(args[0]);
    return BOOL_VAL(IsKeyDown(key));
}

Value koda_shim_IsKeyPressed(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    int key = (int)AS_NUMBER(args[0]);
    return BOOL_VAL(IsKeyPressed(key));
}

Value koda_shim_DrawRectangle(int argCount, Value* args) {
    if (argCount < 5) return NIL_VAL;
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    int w = (int)AS_NUMBER(args[2]);
    int h = (int)AS_NUMBER(args[3]);
    unsigned int c = (unsigned int)AS_NUMBER(args[4]);
    DrawRectangle(x, y, w, h, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_DrawCircle(int argCount, Value* args) {
    if (argCount < 4) return NIL_VAL;
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    float r = (float)AS_NUMBER(args[2]);
    unsigned int c = (unsigned int)AS_NUMBER(args[3]);
    DrawCircle(x, y, r, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_DrawLine3D(int argCount, Value* args) {
    if (argCount < 7) return NIL_VAL;
    Vector3 a = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    Vector3 b = {(float)AS_NUMBER(args[3]), (float)AS_NUMBER(args[4]), (float)AS_NUMBER(args[5])};
    unsigned int c = (unsigned int)AS_NUMBER(args[6]);
    DrawLine3D(a, b, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_GetMouseX(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetMouseX());
}

Value koda_shim_GetMouseY(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetMouseY());
}

Value koda_shim_GetMouseDeltaX(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetMouseDelta().x);
}

Value koda_shim_GetMouseDeltaY(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetMouseDelta().y);
}

Value koda_shim_DisableCursor(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    DisableCursor();
    return NIL_VAL;
}

Value koda_shim_EnableCursor(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    EnableCursor();
    return NIL_VAL;
}

Value koda_shim_IsMouseButtonDown(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    int btn = (int)AS_NUMBER(args[0]);
    return BOOL_VAL(IsMouseButtonDown(btn));
}

Value koda_shim_IsMouseButtonPressed(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    int btn = (int)AS_NUMBER(args[0]);
    return BOOL_VAL(IsMouseButtonPressed(btn));
}

Value koda_shim_GetMouseWheelMove(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetMouseWheelMove());
}

Value koda_shim_GetScreenWidth(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetScreenWidth());
}

Value koda_shim_GetScreenHeight(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetScreenHeight());
}

Value koda_shim_SetWindowTitle(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    SetWindowTitle(koda_arg_cstr(args, 0));
    return NIL_VAL;
}

Value koda_shim_GetFPS(int argCount, Value* args) {
    (void)argCount;
    (void)args;
    return NUMBER_VAL((double)GetFPS());
}

Value koda_shim_DrawLine(int argCount, Value* args) {
    if (argCount < 5) return NIL_VAL;
    int x1 = (int)AS_NUMBER(args[0]);
    int y1 = (int)AS_NUMBER(args[1]);
    int x2 = (int)AS_NUMBER(args[2]);
    int y2 = (int)AS_NUMBER(args[3]);
    unsigned int c = (unsigned int)AS_NUMBER(args[4]);
    DrawLine(x1, y1, x2, y2, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_DrawCircleLines(int argCount, Value* args) {
    if (argCount < 4) return NIL_VAL;
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    float r = (float)AS_NUMBER(args[2]);
    unsigned int c = (unsigned int)AS_NUMBER(args[3]);
    DrawCircleLines(x, y, r, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_DrawRectangleLines(int argCount, Value* args) {
    if (argCount < 5) return NIL_VAL;
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    int w = (int)AS_NUMBER(args[2]);
    int h = (int)AS_NUMBER(args[3]);
    unsigned int c = (unsigned int)AS_NUMBER(args[4]);
    DrawRectangleLines(x, y, w, h, koda_color_from_packed(c));
    return NIL_VAL;
}

#define KODA_MAX_TEXTURES 256
static Texture2D koda_textures[KODA_MAX_TEXTURES];
static int koda_texture_count = 0;

Value koda_shim_LoadTexture(int argCount, Value* args) {
    if (argCount < 1) return NUMBER_VAL(-1);
    if (koda_texture_count >= KODA_MAX_TEXTURES) return NUMBER_VAL(-1);
    const char* path = koda_arg_cstr(args, 0);
    koda_textures[koda_texture_count] = LoadTexture(path);
    int id = koda_texture_count++;
    return NUMBER_VAL((double)id);
}

Value koda_shim_DrawTexture(int argCount, Value* args) {
    if (argCount < 4) return NIL_VAL;
    int id = (int)AS_NUMBER(args[0]);
    int x = (int)AS_NUMBER(args[1]);
    int y = (int)AS_NUMBER(args[2]);
    unsigned int c = (unsigned int)AS_NUMBER(args[3]);
    if (id < 0 || id >= koda_texture_count) return NIL_VAL;
    DrawTexture(koda_textures[id], x, y, koda_color_from_packed(c));
    return NIL_VAL;
}

Value koda_shim_UnloadTexture(int argCount, Value* args) {
    if (argCount < 1) return NIL_VAL;
    int id = (int)AS_NUMBER(args[0]);
    if (id < 0 || id >= koda_texture_count) return NIL_VAL;
    UnloadTexture(koda_textures[id]);
    return NIL_VAL;
}

