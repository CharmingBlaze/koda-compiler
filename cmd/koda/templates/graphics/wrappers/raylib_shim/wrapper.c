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
    ClearBackground(*(Color*)&c);
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
    DrawCubeV(pos, size, *(Color*)&c);
    return NIL_VAL;
}

Value koda_shim_DrawCubeWires(int argCount, Value* args) {
    if (argCount < 7) return NIL_VAL;
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    Vector3 size = {(float)AS_NUMBER(args[3]), (float)AS_NUMBER(args[4]), (float)AS_NUMBER(args[5])};
    unsigned int c = (unsigned int)AS_NUMBER(args[6]);
    DrawCubeWiresV(pos, size, *(Color*)&c);
    return NIL_VAL;
}

Value koda_shim_DrawText(int argCount, Value* args) {
    if (argCount < 5) return NIL_VAL;
    const char* text = koda_arg_cstr(args, 0);
    int x = (int)AS_NUMBER(args[1]);
    int y = (int)AS_NUMBER(args[2]);
    int fs = (int)AS_NUMBER(args[3]);
    unsigned int c = (unsigned int)AS_NUMBER(args[4]);
    DrawText(text, x, y, fs, *(Color*)&c);
    return NIL_VAL;
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
    DrawRectangle(x, y, w, h, *(Color*)&c);
    return NIL_VAL;
}

Value koda_shim_DrawCircle(int argCount, Value* args) {
    if (argCount < 4) return NIL_VAL;
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    float r = (float)AS_NUMBER(args[2]);
    unsigned int c = (unsigned int)AS_NUMBER(args[3]);
    DrawCircle(x, y, r, *(Color*)&c);
    return NIL_VAL;
}

Value koda_shim_DrawLine3D(int argCount, Value* args) {
    if (argCount < 7) return NIL_VAL;
    Vector3 a = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    Vector3 b = {(float)AS_NUMBER(args[3]), (float)AS_NUMBER(args[4]), (float)AS_NUMBER(args[5])};
    unsigned int c = (unsigned int)AS_NUMBER(args[6]);
    DrawLine3D(a, b, *(Color*)&c);
    return NIL_VAL;
}

