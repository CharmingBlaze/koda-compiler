/*
 * Fast numeric paths for hot Raylib draw/camera calls.
 * Linked alongside wrapper.c; invoked by codegen fusion (see native_fusion.go).
 */
#include "koda_wrapgen_abi.h"
#include <raylib.h>

static Color koda_color_from_packed(unsigned int packed) {
    Color c;
    c.r = (unsigned char)((packed >> 24) & 255);
    c.g = (unsigned char)((packed >> 16) & 255);
    c.b = (unsigned char)((packed >> 8) & 255);
    c.a = (unsigned char)(packed & 255);
    return c;
}

static Color koda_color_from_value(KodaValue v) {
    if (IS_NUMBER(v)) {
        return koda_color_from_packed((unsigned int)AS_NUMBER(v));
    }
    Color c = {0, 0, 0, 255};
    KodaValue r = koda_get_index(v, koda_copy_string("r", 1));
    KodaValue g = koda_get_index(v, koda_copy_string("g", 1));
    KodaValue b = koda_get_index(v, koda_copy_string("b", 1));
    KodaValue a = koda_get_index(v, koda_copy_string("a", 1));
    if (IS_NUMBER(r)) {
        c.r = (unsigned char)AS_NUMBER(r);
    }
    if (IS_NUMBER(g)) {
        c.g = (unsigned char)AS_NUMBER(g);
    }
    if (IS_NUMBER(b)) {
        c.b = (unsigned char)AS_NUMBER(b);
    }
    if (IS_NUMBER(a)) {
        c.a = (unsigned char)AS_NUMBER(a);
    }
    return c;
}

/* px py pz width height length r g b a — all numbers (color packed 0xRRGGBBAA). */
KodaValue koda_fast_DrawCube(int argCount, KodaValue* args) {
    if (argCount < 10) {
        return koda_err_str("koda_fast_DrawCube requires 10 numeric arguments");
    }
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    float w = (float)AS_NUMBER(args[3]);
    float h = (float)AS_NUMBER(args[4]);
    float l = (float)AS_NUMBER(args[5]);
    Color col = {
        (unsigned char)AS_NUMBER(args[6]),
        (unsigned char)AS_NUMBER(args[7]),
        (unsigned char)AS_NUMBER(args[8]),
        (unsigned char)AS_NUMBER(args[9]),
    };
    DrawCube(pos, w, h, l, col);
    return NULL_VAL;
}

/* px py pz width height length color — position/size numeric; color object or packed number. */
KodaValue koda_fast_DrawCube8(int argCount, KodaValue* args) {
    if (argCount < 7) {
        return koda_err_str("koda_fast_DrawCube8 requires 7 arguments");
    }
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    float w = (float)AS_NUMBER(args[3]);
    float h = (float)AS_NUMBER(args[4]);
    float l = (float)AS_NUMBER(args[5]);
    Color col = koda_color_from_value(args[6]);
    DrawCube(pos, w, h, l, col);
    return NULL_VAL;
}

KodaValue koda_fast_DrawCubeWires(int argCount, KodaValue* args) {
    if (argCount < 10) {
        return koda_err_str("koda_fast_DrawCubeWires requires 10 numeric arguments");
    }
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    float w = (float)AS_NUMBER(args[3]);
    float h = (float)AS_NUMBER(args[4]);
    float l = (float)AS_NUMBER(args[5]);
    Color col = {
        (unsigned char)AS_NUMBER(args[6]),
        (unsigned char)AS_NUMBER(args[7]),
        (unsigned char)AS_NUMBER(args[8]),
        (unsigned char)AS_NUMBER(args[9]),
    };
    DrawCubeWires(pos, w, h, l, col);
    return NULL_VAL;
}

KodaValue koda_fast_DrawCubeWires8(int argCount, KodaValue* args) {
    if (argCount < 7) {
        return koda_err_str("koda_fast_DrawCubeWires8 requires 7 arguments");
    }
    Vector3 pos = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    float w = (float)AS_NUMBER(args[3]);
    float h = (float)AS_NUMBER(args[4]);
    float l = (float)AS_NUMBER(args[5]);
    Color col = koda_color_from_value(args[6]);
    DrawCubeWires(pos, w, h, l, col);
    return NULL_VAL;
}

static const char* koda_fast_cstr(KodaValue v) {
    if (!IS_OBJ(v) || AS_OBJ(v)->type != OBJ_STRING) {
        return "";
    }
    return ((ObjString*)AS_OBJ(v))->chars;
}

/* color — packed number or {r,g,b,a} object. */
KodaValue koda_fast_ClearBackground(int argCount, KodaValue* args) {
    if (argCount < 1) {
        return koda_err_str("koda_fast_ClearBackground requires 1 argument");
    }
    ClearBackground(koda_color_from_value(args[0]));
    return NULL_VAL;
}

/* posX posY width height color — color packed or object. */
KodaValue koda_fast_DrawRectangle(int argCount, KodaValue* args) {
    if (argCount < 5) {
        return koda_err_str("koda_fast_DrawRectangle requires 5 arguments");
    }
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    int w = (int)AS_NUMBER(args[2]);
    int h = (int)AS_NUMBER(args[3]);
    DrawRectangle(x, y, w, h, koda_color_from_value(args[4]));
    return NULL_VAL;
}

/* centerX centerY radius color */
KodaValue koda_fast_DrawCircle(int argCount, KodaValue* args) {
    if (argCount < 4) {
        return koda_err_str("koda_fast_DrawCircle requires 4 arguments");
    }
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    float r = (float)AS_NUMBER(args[2]);
    DrawCircle(x, y, r, koda_color_from_value(args[3]));
    return NULL_VAL;
}

/* centerX centerY radius color — outline circle. */
KodaValue koda_fast_DrawCircleLines(int argCount, KodaValue* args) {
    if (argCount < 4) {
        return koda_err_str("koda_fast_DrawCircleLines requires 4 arguments");
    }
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    float r = (float)AS_NUMBER(args[2]);
    DrawCircleLines(x, y, r, koda_color_from_value(args[3]));
    return NULL_VAL;
}

/* text posX posY fontSize color */
KodaValue koda_fast_DrawText(int argCount, KodaValue* args) {
    if (argCount < 5) {
        return koda_err_str("koda_fast_DrawText requires 5 arguments");
    }
    const char* text = koda_fast_cstr(args[0]);
    int x = (int)AS_NUMBER(args[1]);
    int y = (int)AS_NUMBER(args[2]);
    int fs = (int)AS_NUMBER(args[3]);
    DrawText(text, x, y, fs, koda_color_from_value(args[4]));
    return NULL_VAL;
}

/* x1 y1 x2 y2 color */
KodaValue koda_fast_DrawLine(int argCount, KodaValue* args) {
    if (argCount < 5) {
        return koda_err_str("koda_fast_DrawLine requires 5 arguments");
    }
    int x1 = (int)AS_NUMBER(args[0]);
    int y1 = (int)AS_NUMBER(args[1]);
    int x2 = (int)AS_NUMBER(args[2]);
    int y2 = (int)AS_NUMBER(args[3]);
    DrawLine(x1, y1, x2, y2, koda_color_from_value(args[4]));
    return NULL_VAL;
}

/* x y w h color — outline rect. */
KodaValue koda_fast_DrawRectangleLines(int argCount, KodaValue* args) {
    if (argCount < 5) {
        return koda_err_str("koda_fast_DrawRectangleLines requires 5 arguments");
    }
    int x = (int)AS_NUMBER(args[0]);
    int y = (int)AS_NUMBER(args[1]);
    int w = (int)AS_NUMBER(args[2]);
    int h = (int)AS_NUMBER(args[3]);
    DrawRectangleLines(x, y, w, h, koda_color_from_value(args[4]));
    return NULL_VAL;
}

/* x1 y1 z1 x2 y2 z2 color — packed color or object. */
KodaValue koda_fast_DrawLine3D(int argCount, KodaValue* args) {
    if (argCount < 7) {
        return koda_err_str("koda_fast_DrawLine3D requires 7 arguments");
    }
    Vector3 a = {(float)AS_NUMBER(args[0]), (float)AS_NUMBER(args[1]), (float)AS_NUMBER(args[2])};
    Vector3 b = {(float)AS_NUMBER(args[3]), (float)AS_NUMBER(args[4]), (float)AS_NUMBER(args[5])};
    DrawLine3D(a, b, koda_color_from_value(args[6]));
    return NULL_VAL;
}

/* ex ey ez tx ty tz upx upy upz fovy — all numbers. */
KodaValue koda_fast_BeginMode3D(int argCount, KodaValue* args) {
    if (argCount < 10) {
        return koda_err_str("koda_fast_BeginMode3D requires 10 numeric arguments");
    }
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
    return NULL_VAL;
}
