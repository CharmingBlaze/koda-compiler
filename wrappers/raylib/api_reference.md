# raylib - API Reference

## Contents

- [Functions](#functions)
- [Structs](#structs)
- [Enums](#enums)
- [Macros](#macros)

---

## Functions

### void

```c
typedef void()
```

**Koda usage**

```koda
let result = void();
```

---

### bool

```c
typedef bool()
```

**Koda usage**

```koda
let result = bool();
```

---

### bool

```c
typedef bool()
```

**Koda usage**

```koda
let result = bool();
```

---

### InitWindow

```c
void InitWindow(int width, int height, const char* title)
```

**Koda usage**

```koda
InitWindow(width, height, title);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `title` | `const char*` |

---

### CloseWindow

```c
void CloseWindow()
```

**Koda usage**

```koda
CloseWindow();
```

---

### WindowShouldClose

```c
bool WindowShouldClose()
```

**Koda usage**

```koda
let result = WindowShouldClose();
```

---

### IsWindowReady

```c
bool IsWindowReady()
```

**Koda usage**

```koda
let result = IsWindowReady();
```

---

### IsWindowFullscreen

```c
bool IsWindowFullscreen()
```

**Koda usage**

```koda
let result = IsWindowFullscreen();
```

---

### IsWindowHidden

```c
bool IsWindowHidden()
```

**Koda usage**

```koda
let result = IsWindowHidden();
```

---

### IsWindowMinimized

```c
bool IsWindowMinimized()
```

**Koda usage**

```koda
let result = IsWindowMinimized();
```

---

### IsWindowMaximized

```c
bool IsWindowMaximized()
```

**Koda usage**

```koda
let result = IsWindowMaximized();
```

---

### IsWindowFocused

```c
bool IsWindowFocused()
```

**Koda usage**

```koda
let result = IsWindowFocused();
```

---

### IsWindowResized

```c
bool IsWindowResized()
```

**Koda usage**

```koda
let result = IsWindowResized();
```

---

### IsWindowState

```c
bool IsWindowState(unsigned int flag)
```

**Koda usage**

```koda
let result = IsWindowState(flag);
```

| Parameter | Type |
|-----------|------|
| `flag` | `unsigned int` |

---

### SetWindowState

```c
void SetWindowState(unsigned int flags)
```

**Koda usage**

```koda
SetWindowState(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### ClearWindowState

```c
void ClearWindowState(unsigned int flags)
```

**Koda usage**

```koda
ClearWindowState(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### ToggleFullscreen

```c
void ToggleFullscreen()
```

**Koda usage**

```koda
ToggleFullscreen();
```

---

### ToggleBorderlessWindowed

```c
void ToggleBorderlessWindowed()
```

**Koda usage**

```koda
ToggleBorderlessWindowed();
```

---

### MaximizeWindow

```c
void MaximizeWindow()
```

**Koda usage**

```koda
MaximizeWindow();
```

---

### MinimizeWindow

```c
void MinimizeWindow()
```

**Koda usage**

```koda
MinimizeWindow();
```

---

### RestoreWindow

```c
void RestoreWindow()
```

**Koda usage**

```koda
RestoreWindow();
```

---

### SetWindowIcon

```c
void SetWindowIcon(Image image)
```

**Koda usage**

```koda
SetWindowIcon(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### SetWindowIcons

```c
void SetWindowIcons(Image* images, int count)
```

**Koda usage**

```koda
SetWindowIcons(images, count);
```

| Parameter | Type |
|-----------|------|
| `images` | `Image*` |
| `count` | `int` |

---

### SetWindowTitle

```c
void SetWindowTitle(const char* title)
```

**Koda usage**

```koda
SetWindowTitle(title);
```

| Parameter | Type |
|-----------|------|
| `title` | `const char*` |

---

### SetWindowPosition

```c
void SetWindowPosition(int x, int y)
```

**Koda usage**

```koda
SetWindowPosition(x, y);
```

| Parameter | Type |
|-----------|------|
| `x` | `int` |
| `y` | `int` |

---

### SetWindowMonitor

```c
void SetWindowMonitor(int monitor)
```

**Koda usage**

```koda
SetWindowMonitor(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### SetWindowMinSize

```c
void SetWindowMinSize(int width, int height)
```

**Koda usage**

```koda
SetWindowMinSize(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### SetWindowMaxSize

```c
void SetWindowMaxSize(int width, int height)
```

**Koda usage**

```koda
SetWindowMaxSize(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### SetWindowSize

```c
void SetWindowSize(int width, int height)
```

**Koda usage**

```koda
SetWindowSize(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### SetWindowOpacity

```c
void SetWindowOpacity(float opacity)
```

**Koda usage**

```koda
SetWindowOpacity(opacity);
```

| Parameter | Type |
|-----------|------|
| `opacity` | `float` |

---

### SetWindowFocused

```c
void SetWindowFocused()
```

**Koda usage**

```koda
SetWindowFocused();
```

---

### GetScreenWidth

```c
int GetScreenWidth()
```

**Koda usage**

```koda
let result = GetScreenWidth();
```

---

### GetScreenHeight

```c
int GetScreenHeight()
```

**Koda usage**

```koda
let result = GetScreenHeight();
```

---

### GetRenderWidth

```c
int GetRenderWidth()
```

**Koda usage**

```koda
let result = GetRenderWidth();
```

---

### GetRenderHeight

```c
int GetRenderHeight()
```

**Koda usage**

```koda
let result = GetRenderHeight();
```

---

### GetMonitorCount

```c
int GetMonitorCount()
```

**Koda usage**

```koda
let result = GetMonitorCount();
```

---

### GetCurrentMonitor

```c
int GetCurrentMonitor()
```

**Koda usage**

```koda
let result = GetCurrentMonitor();
```

---

### GetMonitorPosition

```c
Vector2 GetMonitorPosition(int monitor)
```

**Koda usage**

```koda
let result = GetMonitorPosition(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorWidth

```c
int GetMonitorWidth(int monitor)
```

**Koda usage**

```koda
let result = GetMonitorWidth(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorHeight

```c
int GetMonitorHeight(int monitor)
```

**Koda usage**

```koda
let result = GetMonitorHeight(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorPhysicalWidth

```c
int GetMonitorPhysicalWidth(int monitor)
```

**Koda usage**

```koda
let result = GetMonitorPhysicalWidth(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorPhysicalHeight

```c
int GetMonitorPhysicalHeight(int monitor)
```

**Koda usage**

```koda
let result = GetMonitorPhysicalHeight(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorRefreshRate

```c
int GetMonitorRefreshRate(int monitor)
```

**Koda usage**

```koda
let result = GetMonitorRefreshRate(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetWindowPosition

```c
Vector2 GetWindowPosition()
```

**Koda usage**

```koda
let result = GetWindowPosition();
```

---

### GetWindowScaleDPI

```c
Vector2 GetWindowScaleDPI()
```

**Koda usage**

```koda
let result = GetWindowScaleDPI();
```

---

### SetClipboardText

```c
void SetClipboardText(const char* text)
```

**Koda usage**

```koda
SetClipboardText(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |

---

### GetClipboardImage

```c
Image GetClipboardImage()
```

**Koda usage**

```koda
let result = GetClipboardImage();
```

---

### EnableEventWaiting

```c
void EnableEventWaiting()
```

**Koda usage**

```koda
EnableEventWaiting();
```

---

### DisableEventWaiting

```c
void DisableEventWaiting()
```

**Koda usage**

```koda
DisableEventWaiting();
```

---

### ShowCursor

```c
void ShowCursor()
```

**Koda usage**

```koda
ShowCursor();
```

---

### HideCursor

```c
void HideCursor()
```

**Koda usage**

```koda
HideCursor();
```

---

### IsCursorHidden

```c
bool IsCursorHidden()
```

**Koda usage**

```koda
let result = IsCursorHidden();
```

---

### EnableCursor

```c
void EnableCursor()
```

**Koda usage**

```koda
EnableCursor();
```

---

### DisableCursor

```c
void DisableCursor()
```

**Koda usage**

```koda
DisableCursor();
```

---

### IsCursorOnScreen

```c
bool IsCursorOnScreen()
```

**Koda usage**

```koda
let result = IsCursorOnScreen();
```

---

### ClearBackground

```c
void ClearBackground(Color color)
```

**Koda usage**

```koda
ClearBackground(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### BeginDrawing

```c
void BeginDrawing()
```

**Koda usage**

```koda
BeginDrawing();
```

---

### EndDrawing

```c
void EndDrawing()
```

**Koda usage**

```koda
EndDrawing();
```

---

### BeginMode2D

```c
void BeginMode2D(Camera2D camera)
```

**Koda usage**

```koda
BeginMode2D(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera2D` |

---

### EndMode2D

```c
void EndMode2D()
```

**Koda usage**

```koda
EndMode2D();
```

---

### BeginMode3D

```c
void BeginMode3D(Camera3D camera)
```

**Koda usage**

```koda
BeginMode3D(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera3D` |

---

### EndMode3D

```c
void EndMode3D()
```

**Koda usage**

```koda
EndMode3D();
```

---

### BeginTextureMode

```c
void BeginTextureMode(RenderTexture2D target)
```

**Koda usage**

```koda
BeginTextureMode(target);
```

| Parameter | Type |
|-----------|------|
| `target` | `RenderTexture2D` |

---

### EndTextureMode

```c
void EndTextureMode()
```

**Koda usage**

```koda
EndTextureMode();
```

---

### BeginShaderMode

```c
void BeginShaderMode(Shader shader)
```

**Koda usage**

```koda
BeginShaderMode(shader);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |

---

### EndShaderMode

```c
void EndShaderMode()
```

**Koda usage**

```koda
EndShaderMode();
```

---

### BeginBlendMode

```c
void BeginBlendMode(int mode)
```

**Koda usage**

```koda
BeginBlendMode(mode);
```

| Parameter | Type |
|-----------|------|
| `mode` | `int` |

---

### EndBlendMode

```c
void EndBlendMode()
```

**Koda usage**

```koda
EndBlendMode();
```

---

### BeginScissorMode

```c
void BeginScissorMode(int x, int y, int width, int height)
```

**Koda usage**

```koda
BeginScissorMode(x, y, width, height);
```

| Parameter | Type |
|-----------|------|
| `x` | `int` |
| `y` | `int` |
| `width` | `int` |
| `height` | `int` |

---

### EndScissorMode

```c
void EndScissorMode()
```

**Koda usage**

```koda
EndScissorMode();
```

---

### BeginVrStereoMode

```c
void BeginVrStereoMode(VrStereoConfig config)
```

**Koda usage**

```koda
BeginVrStereoMode(config);
```

| Parameter | Type |
|-----------|------|
| `config` | `VrStereoConfig` |

---

### EndVrStereoMode

```c
void EndVrStereoMode()
```

**Koda usage**

```koda
EndVrStereoMode();
```

---

### LoadVrStereoConfig

```c
VrStereoConfig LoadVrStereoConfig(VrDeviceInfo device)
```

**Koda usage**

```koda
let result = LoadVrStereoConfig(device);
```

| Parameter | Type |
|-----------|------|
| `device` | `VrDeviceInfo` |

---

### UnloadVrStereoConfig

```c
void UnloadVrStereoConfig(VrStereoConfig config)
```

**Koda usage**

```koda
UnloadVrStereoConfig(config);
```

| Parameter | Type |
|-----------|------|
| `config` | `VrStereoConfig` |

---

### LoadShader

```c
Shader LoadShader(const char* vsFileName, const char* fsFileName)
```

**Koda usage**

```koda
let result = LoadShader(vsFileName, fsFileName);
```

| Parameter | Type |
|-----------|------|
| `vsFileName` | `const char*` |
| `fsFileName` | `const char*` |

---

### LoadShaderFromMemory

```c
Shader LoadShaderFromMemory(const char* vsCode, const char* fsCode)
```

**Koda usage**

```koda
let result = LoadShaderFromMemory(vsCode, fsCode);
```

| Parameter | Type |
|-----------|------|
| `vsCode` | `const char*` |
| `fsCode` | `const char*` |

---

### IsShaderValid

```c
bool IsShaderValid(Shader shader)
```

**Koda usage**

```koda
let result = IsShaderValid(shader);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |

---

### GetShaderLocation

```c
int GetShaderLocation(Shader shader, const char* uniformName)
```

**Koda usage**

```koda
let result = GetShaderLocation(shader, uniformName);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `uniformName` | `const char*` |

---

### GetShaderLocationAttrib

```c
int GetShaderLocationAttrib(Shader shader, const char* attribName)
```

**Koda usage**

```koda
let result = GetShaderLocationAttrib(shader, attribName);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `attribName` | `const char*` |

---

### SetShaderValue

```c
void SetShaderValue(Shader shader, int locIndex, const void* value, int uniformType)
```

**Koda usage**

```koda
SetShaderValue(shader, locIndex, value, uniformType);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `value` | `const void*` |
| `uniformType` | `int` |

---

### SetShaderValueV

```c
void SetShaderValueV(Shader shader, int locIndex, const void* value, int uniformType, int count)
```

**Koda usage**

```koda
SetShaderValueV(shader, locIndex, value, uniformType, count);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `value` | `const void*` |
| `uniformType` | `int` |
| `count` | `int` |

---

### SetShaderValueMatrix

```c
void SetShaderValueMatrix(Shader shader, int locIndex, Matrix mat)
```

**Koda usage**

```koda
SetShaderValueMatrix(shader, locIndex, mat);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `mat` | `Matrix` |

---

### SetShaderValueTexture

```c
void SetShaderValueTexture(Shader shader, int locIndex, Texture2D texture)
```

**Koda usage**

```koda
SetShaderValueTexture(shader, locIndex, texture);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `texture` | `Texture2D` |

---

### UnloadShader

```c
void UnloadShader(Shader shader)
```

**Koda usage**

```koda
UnloadShader(shader);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |

---

### GetScreenToWorldRay

```c
Ray GetScreenToWorldRay(Vector2 position, Camera camera)
```

**Koda usage**

```koda
let result = GetScreenToWorldRay(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera` |

---

### GetScreenToWorldRayEx

```c
Ray GetScreenToWorldRayEx(Vector2 position, Camera camera, int width, int height)
```

**Koda usage**

```koda
let result = GetScreenToWorldRayEx(position, camera, width, height);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera` |
| `width` | `int` |
| `height` | `int` |

---

### GetWorldToScreen

```c
Vector2 GetWorldToScreen(Vector3 position, Camera camera)
```

**Koda usage**

```koda
let result = GetWorldToScreen(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `camera` | `Camera` |

---

### GetWorldToScreenEx

```c
Vector2 GetWorldToScreenEx(Vector3 position, Camera camera, int width, int height)
```

**Koda usage**

```koda
let result = GetWorldToScreenEx(position, camera, width, height);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `camera` | `Camera` |
| `width` | `int` |
| `height` | `int` |

---

### GetWorldToScreen2D

```c
Vector2 GetWorldToScreen2D(Vector2 position, Camera2D camera)
```

**Koda usage**

```koda
let result = GetWorldToScreen2D(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera2D` |

---

### GetScreenToWorld2D

```c
Vector2 GetScreenToWorld2D(Vector2 position, Camera2D camera)
```

**Koda usage**

```koda
let result = GetScreenToWorld2D(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera2D` |

---

### GetCameraMatrix

```c
Matrix GetCameraMatrix(Camera camera)
```

**Koda usage**

```koda
let result = GetCameraMatrix(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |

---

### GetCameraMatrix2D

```c
Matrix GetCameraMatrix2D(Camera2D camera)
```

**Koda usage**

```koda
let result = GetCameraMatrix2D(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera2D` |

---

### SetTargetFPS

```c
void SetTargetFPS(int fps)
```

**Koda usage**

```koda
SetTargetFPS(fps);
```

| Parameter | Type |
|-----------|------|
| `fps` | `int` |

---

### GetFrameTime

```c
float GetFrameTime()
```

**Koda usage**

```koda
let result = GetFrameTime();
```

---

### GetTime

```c
double GetTime()
```

**Koda usage**

```koda
let result = GetTime();
```

---

### GetFPS

```c
int GetFPS()
```

**Koda usage**

```koda
let result = GetFPS();
```

---

### SwapScreenBuffer

```c
void SwapScreenBuffer()
```

**Koda usage**

```koda
SwapScreenBuffer();
```

---

### PollInputEvents

```c
void PollInputEvents()
```

**Koda usage**

```koda
PollInputEvents();
```

---

### WaitTime

```c
void WaitTime(double seconds)
```

**Koda usage**

```koda
WaitTime(seconds);
```

| Parameter | Type |
|-----------|------|
| `seconds` | `double` |

---

### SetRandomSeed

```c
void SetRandomSeed(unsigned int seed)
```

**Koda usage**

```koda
SetRandomSeed(seed);
```

| Parameter | Type |
|-----------|------|
| `seed` | `unsigned int` |

---

### GetRandomValue

```c
int GetRandomValue(int min, int max)
```

**Koda usage**

```koda
let result = GetRandomValue(min, max);
```

| Parameter | Type |
|-----------|------|
| `min` | `int` |
| `max` | `int` |

---

### UnloadRandomSequence

```c
void UnloadRandomSequence(int* sequence)
```

**Koda usage**

```koda
UnloadRandomSequence(sequence);
```

| Parameter | Type |
|-----------|------|
| `sequence` | `int*` |

---

### TakeScreenshot

```c
void TakeScreenshot(const char* fileName)
```

**Koda usage**

```koda
TakeScreenshot(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### SetConfigFlags

```c
void SetConfigFlags(unsigned int flags)
```

**Koda usage**

```koda
SetConfigFlags(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### OpenURL

```c
void OpenURL(const char* url)
```

**Koda usage**

```koda
OpenURL(url);
```

| Parameter | Type |
|-----------|------|
| `url` | `const char*` |

---

### SetTraceLogLevel

```c
void SetTraceLogLevel(int logLevel)
```

**Koda usage**

```koda
SetTraceLogLevel(logLevel);
```

| Parameter | Type |
|-----------|------|
| `logLevel` | `int` |

---

### TraceLog

```c
void TraceLog(int logLevel, const char* text)
```

**Koda usage**

```koda
TraceLog(logLevel, text);
```

| Parameter | Type |
|-----------|------|
| `logLevel` | `int` |
| `text` | `const char*` |

---

### SetTraceLogCallback

```c
void SetTraceLogCallback(TraceLogCallback callback)
```

**Koda usage**

```koda
SetTraceLogCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `TraceLogCallback` |

---

### MemFree

```c
void MemFree(void* ptr)
```

**Koda usage**

```koda
MemFree(ptr);
```

| Parameter | Type |
|-----------|------|
| `ptr` | `void*` |

---

### UnloadFileData

```c
void UnloadFileData(unsigned char* data)
```

**Koda usage**

```koda
UnloadFileData(data);
```

| Parameter | Type |
|-----------|------|
| `data` | `unsigned char*` |

---

### SaveFileData

```c
bool SaveFileData(const char* fileName, const void* data, int dataSize)
```

**Koda usage**

```koda
let result = SaveFileData(fileName, data, dataSize);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `data` | `const void*` |
| `dataSize` | `int` |

---

### ExportDataAsCode

```c
bool ExportDataAsCode(const unsigned char* data, int dataSize, const char* fileName)
```

**Koda usage**

```koda
let result = ExportDataAsCode(data, dataSize, fileName);
```

| Parameter | Type |
|-----------|------|
| `data` | `const unsigned char*` |
| `dataSize` | `int` |
| `fileName` | `const char*` |

---

### UnloadFileText

```c
void UnloadFileText(char* text)
```

**Koda usage**

```koda
UnloadFileText(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `char*` |

---

### SaveFileText

```c
bool SaveFileText(const char* fileName, const char* text)
```

**Koda usage**

```koda
let result = SaveFileText(fileName, text);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `text` | `const char*` |

---

### SetLoadFileDataCallback

```c
void SetLoadFileDataCallback(LoadFileDataCallback callback)
```

**Koda usage**

```koda
SetLoadFileDataCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `LoadFileDataCallback` |

---

### SetSaveFileDataCallback

```c
void SetSaveFileDataCallback(SaveFileDataCallback callback)
```

**Koda usage**

```koda
SetSaveFileDataCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `SaveFileDataCallback` |

---

### SetLoadFileTextCallback

```c
void SetLoadFileTextCallback(LoadFileTextCallback callback)
```

**Koda usage**

```koda
SetLoadFileTextCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `LoadFileTextCallback` |

---

### SetSaveFileTextCallback

```c
void SetSaveFileTextCallback(SaveFileTextCallback callback)
```

**Koda usage**

```koda
SetSaveFileTextCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `SaveFileTextCallback` |

---

### FileRename

```c
int FileRename(const char* fileName, const char* fileRename)
```

**Koda usage**

```koda
let result = FileRename(fileName, fileRename);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `fileRename` | `const char*` |

---

### FileRemove

```c
int FileRemove(const char* fileName)
```

**Koda usage**

```koda
let result = FileRemove(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### FileCopy

```c
int FileCopy(const char* srcPath, const char* dstPath)
```

**Koda usage**

```koda
let result = FileCopy(srcPath, dstPath);
```

| Parameter | Type |
|-----------|------|
| `srcPath` | `const char*` |
| `dstPath` | `const char*` |

---

### FileMove

```c
int FileMove(const char* srcPath, const char* dstPath)
```

**Koda usage**

```koda
let result = FileMove(srcPath, dstPath);
```

| Parameter | Type |
|-----------|------|
| `srcPath` | `const char*` |
| `dstPath` | `const char*` |

---

### FileTextReplace

```c
int FileTextReplace(const char* fileName, const char* search, const char* replacement)
```

**Koda usage**

```koda
let result = FileTextReplace(fileName, search, replacement);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `search` | `const char*` |
| `replacement` | `const char*` |

---

### FileTextFindIndex

```c
int FileTextFindIndex(const char* fileName, const char* search)
```

**Koda usage**

```koda
let result = FileTextFindIndex(fileName, search);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `search` | `const char*` |

---

### FileExists

```c
bool FileExists(const char* fileName)
```

**Koda usage**

```koda
let result = FileExists(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### DirectoryExists

```c
bool DirectoryExists(const char* dirPath)
```

**Koda usage**

```koda
let result = DirectoryExists(dirPath);
```

| Parameter | Type |
|-----------|------|
| `dirPath` | `const char*` |

---

### IsFileExtension

```c
bool IsFileExtension(const char* fileName, const char* ext)
```

**Koda usage**

```koda
let result = IsFileExtension(fileName, ext);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `ext` | `const char*` |

---

### GetFileLength

```c
int GetFileLength(const char* fileName)
```

**Koda usage**

```koda
let result = GetFileLength(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### GetFileModTime

```c
long GetFileModTime(const char* fileName)
```

**Koda usage**

```koda
let result = GetFileModTime(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### MakeDirectory

```c
int MakeDirectory(const char* dirPath)
```

**Koda usage**

```koda
let result = MakeDirectory(dirPath);
```

| Parameter | Type |
|-----------|------|
| `dirPath` | `const char*` |

---

### ChangeDirectory

```c
bool ChangeDirectory(const char* dirPath)
```

**Koda usage**

```koda
let result = ChangeDirectory(dirPath);
```

| Parameter | Type |
|-----------|------|
| `dirPath` | `const char*` |

---

### IsPathFile

```c
bool IsPathFile(const char* path)
```

**Koda usage**

```koda
let result = IsPathFile(path);
```

| Parameter | Type |
|-----------|------|
| `path` | `const char*` |

---

### IsFileNameValid

```c
bool IsFileNameValid(const char* fileName)
```

**Koda usage**

```koda
let result = IsFileNameValid(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadDirectoryFiles

```c
FilePathList LoadDirectoryFiles(const char* dirPath)
```

**Koda usage**

```koda
let result = LoadDirectoryFiles(dirPath);
```

| Parameter | Type |
|-----------|------|
| `dirPath` | `const char*` |

---

### LoadDirectoryFilesEx

```c
FilePathList LoadDirectoryFilesEx(const char* basePath, const char* filter, bool scanSubdirs)
```

**Koda usage**

```koda
let result = LoadDirectoryFilesEx(basePath, filter, scanSubdirs);
```

| Parameter | Type |
|-----------|------|
| `basePath` | `const char*` |
| `filter` | `const char*` |
| `scanSubdirs` | `bool` |

---

### UnloadDirectoryFiles

```c
void UnloadDirectoryFiles(FilePathList files)
```

**Koda usage**

```koda
UnloadDirectoryFiles(files);
```

| Parameter | Type |
|-----------|------|
| `files` | `FilePathList` |

---

### IsFileDropped

```c
bool IsFileDropped()
```

**Koda usage**

```koda
let result = IsFileDropped();
```

---

### LoadDroppedFiles

```c
FilePathList LoadDroppedFiles()
```

**Koda usage**

```koda
let result = LoadDroppedFiles();
```

---

### UnloadDroppedFiles

```c
void UnloadDroppedFiles(FilePathList files)
```

**Koda usage**

```koda
UnloadDroppedFiles(files);
```

| Parameter | Type |
|-----------|------|
| `files` | `FilePathList` |

---

### GetDirectoryFileCount

```c
unsigned int GetDirectoryFileCount(const char* dirPath)
```

**Koda usage**

```koda
let result = GetDirectoryFileCount(dirPath);
```

| Parameter | Type |
|-----------|------|
| `dirPath` | `const char*` |

---

### GetDirectoryFileCountEx

```c
unsigned int GetDirectoryFileCountEx(const char* basePath, const char* filter, bool scanSubdirs)
```

**Koda usage**

```koda
let result = GetDirectoryFileCountEx(basePath, filter, scanSubdirs);
```

| Parameter | Type |
|-----------|------|
| `basePath` | `const char*` |
| `filter` | `const char*` |
| `scanSubdirs` | `bool` |

---

### ComputeCRC32

```c
unsigned int ComputeCRC32(const unsigned char* data, int dataSize)
```

**Koda usage**

```koda
let result = ComputeCRC32(data, dataSize);
```

| Parameter | Type |
|-----------|------|
| `data` | `const unsigned char*` |
| `dataSize` | `int` |

---

### LoadAutomationEventList

```c
AutomationEventList LoadAutomationEventList(const char* fileName)
```

**Koda usage**

```koda
let result = LoadAutomationEventList(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### UnloadAutomationEventList

```c
void UnloadAutomationEventList(AutomationEventList list)
```

**Koda usage**

```koda
UnloadAutomationEventList(list);
```

| Parameter | Type |
|-----------|------|
| `list` | `AutomationEventList` |

---

### ExportAutomationEventList

```c
bool ExportAutomationEventList(AutomationEventList list, const char* fileName)
```

**Koda usage**

```koda
let result = ExportAutomationEventList(list, fileName);
```

| Parameter | Type |
|-----------|------|
| `list` | `AutomationEventList` |
| `fileName` | `const char*` |

---

### SetAutomationEventList

```c
void SetAutomationEventList(AutomationEventList* list)
```

**Koda usage**

```koda
SetAutomationEventList(list);
```

| Parameter | Type |
|-----------|------|
| `list` | `AutomationEventList*` |

---

### SetAutomationEventBaseFrame

```c
void SetAutomationEventBaseFrame(int frame)
```

**Koda usage**

```koda
SetAutomationEventBaseFrame(frame);
```

| Parameter | Type |
|-----------|------|
| `frame` | `int` |

---

### StartAutomationEventRecording

```c
void StartAutomationEventRecording()
```

**Koda usage**

```koda
StartAutomationEventRecording();
```

---

### StopAutomationEventRecording

```c
void StopAutomationEventRecording()
```

**Koda usage**

```koda
StopAutomationEventRecording();
```

---

### PlayAutomationEvent

```c
void PlayAutomationEvent(AutomationEvent event)
```

**Koda usage**

```koda
PlayAutomationEvent(event);
```

| Parameter | Type |
|-----------|------|
| `event` | `AutomationEvent` |

---

### IsKeyPressed

```c
bool IsKeyPressed(int key)
```

**Koda usage**

```koda
let result = IsKeyPressed(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyPressedRepeat

```c
bool IsKeyPressedRepeat(int key)
```

**Koda usage**

```koda
let result = IsKeyPressedRepeat(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyDown

```c
bool IsKeyDown(int key)
```

**Koda usage**

```koda
let result = IsKeyDown(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyReleased

```c
bool IsKeyReleased(int key)
```

**Koda usage**

```koda
let result = IsKeyReleased(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyUp

```c
bool IsKeyUp(int key)
```

**Koda usage**

```koda
let result = IsKeyUp(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### GetKeyPressed

```c
int GetKeyPressed()
```

**Koda usage**

```koda
let result = GetKeyPressed();
```

---

### GetCharPressed

```c
int GetCharPressed()
```

**Koda usage**

```koda
let result = GetCharPressed();
```

---

### SetExitKey

```c
void SetExitKey(int key)
```

**Koda usage**

```koda
SetExitKey(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsGamepadAvailable

```c
bool IsGamepadAvailable(int gamepad)
```

**Koda usage**

```koda
let result = IsGamepadAvailable(gamepad);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |

---

### IsGamepadButtonPressed

```c
bool IsGamepadButtonPressed(int gamepad, int button)
```

**Koda usage**

```koda
let result = IsGamepadButtonPressed(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### IsGamepadButtonDown

```c
bool IsGamepadButtonDown(int gamepad, int button)
```

**Koda usage**

```koda
let result = IsGamepadButtonDown(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### IsGamepadButtonReleased

```c
bool IsGamepadButtonReleased(int gamepad, int button)
```

**Koda usage**

```koda
let result = IsGamepadButtonReleased(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### IsGamepadButtonUp

```c
bool IsGamepadButtonUp(int gamepad, int button)
```

**Koda usage**

```koda
let result = IsGamepadButtonUp(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### GetGamepadButtonPressed

```c
int GetGamepadButtonPressed()
```

**Koda usage**

```koda
let result = GetGamepadButtonPressed();
```

---

### GetGamepadAxisCount

```c
int GetGamepadAxisCount(int gamepad)
```

**Koda usage**

```koda
let result = GetGamepadAxisCount(gamepad);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |

---

### GetGamepadAxisMovement

```c
float GetGamepadAxisMovement(int gamepad, int axis)
```

**Koda usage**

```koda
let result = GetGamepadAxisMovement(gamepad, axis);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `axis` | `int` |

---

### SetGamepadMappings

```c
int SetGamepadMappings(const char* mappings)
```

**Koda usage**

```koda
let result = SetGamepadMappings(mappings);
```

| Parameter | Type |
|-----------|------|
| `mappings` | `const char*` |

---

### SetGamepadVibration

```c
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
```

**Koda usage**

```koda
SetGamepadVibration(gamepad, leftMotor, rightMotor, duration);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `leftMotor` | `float` |
| `rightMotor` | `float` |
| `duration` | `float` |

---

### IsMouseButtonPressed

```c
bool IsMouseButtonPressed(int button)
```

**Koda usage**

```koda
let result = IsMouseButtonPressed(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### IsMouseButtonDown

```c
bool IsMouseButtonDown(int button)
```

**Koda usage**

```koda
let result = IsMouseButtonDown(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### IsMouseButtonReleased

```c
bool IsMouseButtonReleased(int button)
```

**Koda usage**

```koda
let result = IsMouseButtonReleased(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### IsMouseButtonUp

```c
bool IsMouseButtonUp(int button)
```

**Koda usage**

```koda
let result = IsMouseButtonUp(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### GetMouseX

```c
int GetMouseX()
```

**Koda usage**

```koda
let result = GetMouseX();
```

---

### GetMouseY

```c
int GetMouseY()
```

**Koda usage**

```koda
let result = GetMouseY();
```

---

### GetMousePosition

```c
Vector2 GetMousePosition()
```

**Koda usage**

```koda
let result = GetMousePosition();
```

---

### GetMouseDelta

```c
Vector2 GetMouseDelta()
```

**Koda usage**

```koda
let result = GetMouseDelta();
```

---

### SetMousePosition

```c
void SetMousePosition(int x, int y)
```

**Koda usage**

```koda
SetMousePosition(x, y);
```

| Parameter | Type |
|-----------|------|
| `x` | `int` |
| `y` | `int` |

---

### SetMouseOffset

```c
void SetMouseOffset(int offsetX, int offsetY)
```

**Koda usage**

```koda
SetMouseOffset(offsetX, offsetY);
```

| Parameter | Type |
|-----------|------|
| `offsetX` | `int` |
| `offsetY` | `int` |

---

### SetMouseScale

```c
void SetMouseScale(float scaleX, float scaleY)
```

**Koda usage**

```koda
SetMouseScale(scaleX, scaleY);
```

| Parameter | Type |
|-----------|------|
| `scaleX` | `float` |
| `scaleY` | `float` |

---

### GetMouseWheelMove

```c
float GetMouseWheelMove()
```

**Koda usage**

```koda
let result = GetMouseWheelMove();
```

---

### GetMouseWheelMoveV

```c
Vector2 GetMouseWheelMoveV()
```

**Koda usage**

```koda
let result = GetMouseWheelMoveV();
```

---

### SetMouseCursor

```c
void SetMouseCursor(int cursor)
```

**Koda usage**

```koda
SetMouseCursor(cursor);
```

| Parameter | Type |
|-----------|------|
| `cursor` | `int` |

---

### GetTouchX

```c
int GetTouchX()
```

**Koda usage**

```koda
let result = GetTouchX();
```

---

### GetTouchY

```c
int GetTouchY()
```

**Koda usage**

```koda
let result = GetTouchY();
```

---

### GetTouchPosition

```c
Vector2 GetTouchPosition(int index)
```

**Koda usage**

```koda
let result = GetTouchPosition(index);
```

| Parameter | Type |
|-----------|------|
| `index` | `int` |

---

### GetTouchPointId

```c
int GetTouchPointId(int index)
```

**Koda usage**

```koda
let result = GetTouchPointId(index);
```

| Parameter | Type |
|-----------|------|
| `index` | `int` |

---

### GetTouchPointCount

```c
int GetTouchPointCount()
```

**Koda usage**

```koda
let result = GetTouchPointCount();
```

---

### SetGesturesEnabled

```c
void SetGesturesEnabled(unsigned int flags)
```

**Koda usage**

```koda
SetGesturesEnabled(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### IsGestureDetected

```c
bool IsGestureDetected(unsigned int gesture)
```

**Koda usage**

```koda
let result = IsGestureDetected(gesture);
```

| Parameter | Type |
|-----------|------|
| `gesture` | `unsigned int` |

---

### GetGestureDetected

```c
int GetGestureDetected()
```

**Koda usage**

```koda
let result = GetGestureDetected();
```

---

### GetGestureHoldDuration

```c
float GetGestureHoldDuration()
```

**Koda usage**

```koda
let result = GetGestureHoldDuration();
```

---

### GetGestureDragVector

```c
Vector2 GetGestureDragVector()
```

**Koda usage**

```koda
let result = GetGestureDragVector();
```

---

### GetGestureDragAngle

```c
float GetGestureDragAngle()
```

**Koda usage**

```koda
let result = GetGestureDragAngle();
```

---

### GetGesturePinchVector

```c
Vector2 GetGesturePinchVector()
```

**Koda usage**

```koda
let result = GetGesturePinchVector();
```

---

### GetGesturePinchAngle

```c
float GetGesturePinchAngle()
```

**Koda usage**

```koda
let result = GetGesturePinchAngle();
```

---

### UpdateCamera

```c
void UpdateCamera(Camera* camera, int mode)
```

**Koda usage**

```koda
UpdateCamera(camera, mode);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera*` |
| `mode` | `int` |

---

### UpdateCameraPro

```c
void UpdateCameraPro(Camera* camera, Vector3 movement, Vector3 rotation, float zoom)
```

**Koda usage**

```koda
UpdateCameraPro(camera, movement, rotation, zoom);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera*` |
| `movement` | `Vector3` |
| `rotation` | `Vector3` |
| `zoom` | `float` |

---

### SetShapesTexture

```c
void SetShapesTexture(Texture2D texture, Rectangle source)
```

**Koda usage**

```koda
SetShapesTexture(texture, source);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `source` | `Rectangle` |

---

### GetShapesTexture

```c
Texture2D GetShapesTexture()
```

**Koda usage**

```koda
let result = GetShapesTexture();
```

---

### GetShapesTextureRectangle

```c
Rectangle GetShapesTextureRectangle()
```

**Koda usage**

```koda
let result = GetShapesTextureRectangle();
```

---

### DrawPixel

```c
void DrawPixel(int posX, int posY, Color color)
```

**Koda usage**

```koda
DrawPixel(posX, posY, color);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `color` | `Color` |

---

### DrawPixelV

```c
void DrawPixelV(Vector2 position, Color color)
```

**Koda usage**

```koda
DrawPixelV(position, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `color` | `Color` |

---

### DrawLine

```c
void DrawLine(int startPosX, int startPosY, int endPosX, int endPosY, Color color)
```

**Koda usage**

```koda
DrawLine(startPosX, startPosY, endPosX, endPosY, color);
```

| Parameter | Type |
|-----------|------|
| `startPosX` | `int` |
| `startPosY` | `int` |
| `endPosX` | `int` |
| `endPosY` | `int` |
| `color` | `Color` |

---

### DrawLineV

```c
void DrawLineV(Vector2 startPos, Vector2 endPos, Color color)
```

**Koda usage**

```koda
DrawLineV(startPos, endPos, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `color` | `Color` |

---

### DrawLineEx

```c
void DrawLineEx(Vector2 startPos, Vector2 endPos, float thick, Color color)
```

**Koda usage**

```koda
DrawLineEx(startPos, endPos, thick, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawLineStrip

```c
void DrawLineStrip(const Vector2* points, int pointCount, Color color)
```

**Koda usage**

```koda
DrawLineStrip(points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawLineBezier

```c
void DrawLineBezier(Vector2 startPos, Vector2 endPos, float thick, Color color)
```

**Koda usage**

```koda
DrawLineBezier(startPos, endPos, thick, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawLineDashed

```c
void DrawLineDashed(Vector2 startPos, Vector2 endPos, int dashSize, int spaceSize, Color color)
```

**Koda usage**

```koda
DrawLineDashed(startPos, endPos, dashSize, spaceSize, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `dashSize` | `int` |
| `spaceSize` | `int` |
| `color` | `Color` |

---

### DrawCircle

```c
void DrawCircle(int centerX, int centerY, float radius, Color color)
```

**Koda usage**

```koda
DrawCircle(centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawCircleV

```c
void DrawCircleV(Vector2 center, float radius, Color color)
```

**Koda usage**

```koda
DrawCircleV(center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawCircleGradient

```c
void DrawCircleGradient(Vector2 center, float radius, Color inner, Color outer)
```

**Koda usage**

```koda
DrawCircleGradient(center, radius, inner, outer);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `inner` | `Color` |
| `outer` | `Color` |

---

### DrawCircleSector

```c
void DrawCircleSector(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color)
```

**Koda usage**

```koda
DrawCircleSector(center, radius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawCircleSectorLines

```c
void DrawCircleSectorLines(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color)
```

**Koda usage**

```koda
DrawCircleSectorLines(center, radius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawCircleLines

```c
void DrawCircleLines(int centerX, int centerY, float radius, Color color)
```

**Koda usage**

```koda
DrawCircleLines(centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawCircleLinesV

```c
void DrawCircleLinesV(Vector2 center, float radius, Color color)
```

**Koda usage**

```koda
DrawCircleLinesV(center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawEllipse

```c
void DrawEllipse(int centerX, int centerY, float radiusH, float radiusV, Color color)
```

**Koda usage**

```koda
DrawEllipse(centerX, centerY, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawEllipseV

```c
void DrawEllipseV(Vector2 center, float radiusH, float radiusV, Color color)
```

**Koda usage**

```koda
DrawEllipseV(center, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawEllipseLines

```c
void DrawEllipseLines(int centerX, int centerY, float radiusH, float radiusV, Color color)
```

**Koda usage**

```koda
DrawEllipseLines(centerX, centerY, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawEllipseLinesV

```c
void DrawEllipseLinesV(Vector2 center, float radiusH, float radiusV, Color color)
```

**Koda usage**

```koda
DrawEllipseLinesV(center, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawRing

```c
void DrawRing(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color)
```

**Koda usage**

```koda
DrawRing(center, innerRadius, outerRadius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `innerRadius` | `float` |
| `outerRadius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRingLines

```c
void DrawRingLines(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color)
```

**Koda usage**

```koda
DrawRingLines(center, innerRadius, outerRadius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `innerRadius` | `float` |
| `outerRadius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRectangle

```c
void DrawRectangle(int posX, int posY, int width, int height, Color color)
```

**Koda usage**

```koda
DrawRectangle(posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### DrawRectangleV

```c
void DrawRectangleV(Vector2 position, Vector2 size, Color color)
```

**Koda usage**

```koda
DrawRectangleV(position, size, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `size` | `Vector2` |
| `color` | `Color` |

---

### DrawRectangleRec

```c
void DrawRectangleRec(Rectangle rec, Color color)
```

**Koda usage**

```koda
DrawRectangleRec(rec, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `color` | `Color` |

---

### DrawRectanglePro

```c
void DrawRectanglePro(Rectangle rec, Vector2 origin, float rotation, Color color)
```

**Koda usage**

```koda
DrawRectanglePro(rec, origin, rotation, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `color` | `Color` |

---

### DrawRectangleGradientV

```c
void DrawRectangleGradientV(int posX, int posY, int width, int height, Color top, Color bottom)
```

**Koda usage**

```koda
DrawRectangleGradientV(posX, posY, width, height, top, bottom);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `top` | `Color` |
| `bottom` | `Color` |

---

### DrawRectangleGradientH

```c
void DrawRectangleGradientH(int posX, int posY, int width, int height, Color left, Color right)
```

**Koda usage**

```koda
DrawRectangleGradientH(posX, posY, width, height, left, right);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `left` | `Color` |
| `right` | `Color` |

---

### DrawRectangleGradientEx

```c
void DrawRectangleGradientEx(Rectangle rec, Color topLeft, Color bottomLeft, Color bottomRight, Color topRight)
```

**Koda usage**

```koda
DrawRectangleGradientEx(rec, topLeft, bottomLeft, bottomRight, topRight);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `topLeft` | `Color` |
| `bottomLeft` | `Color` |
| `bottomRight` | `Color` |
| `topRight` | `Color` |

---

### DrawRectangleLines

```c
void DrawRectangleLines(int posX, int posY, int width, int height, Color color)
```

**Koda usage**

```koda
DrawRectangleLines(posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### DrawRectangleLinesEx

```c
void DrawRectangleLinesEx(Rectangle rec, float lineThick, Color color)
```

**Koda usage**

```koda
DrawRectangleLinesEx(rec, lineThick, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `lineThick` | `float` |
| `color` | `Color` |

---

### DrawRectangleRounded

```c
void DrawRectangleRounded(Rectangle rec, float roundness, int segments, Color color)
```

**Koda usage**

```koda
DrawRectangleRounded(rec, roundness, segments, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `roundness` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRectangleRoundedLines

```c
void DrawRectangleRoundedLines(Rectangle rec, float roundness, int segments, Color color)
```

**Koda usage**

```koda
DrawRectangleRoundedLines(rec, roundness, segments, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `roundness` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRectangleRoundedLinesEx

```c
void DrawRectangleRoundedLinesEx(Rectangle rec, float roundness, int segments, float lineThick, Color color)
```

**Koda usage**

```koda
DrawRectangleRoundedLinesEx(rec, roundness, segments, lineThick, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `roundness` | `float` |
| `segments` | `int` |
| `lineThick` | `float` |
| `color` | `Color` |

---

### DrawTriangle

```c
void DrawTriangle(Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda usage**

```koda
DrawTriangle(v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### DrawTriangleGradient

```c
void DrawTriangleGradient(Vector2 v1, Vector2 v2, Vector2 v3, Color c1, Color c2, Color c3)
```

**Koda usage**

```koda
DrawTriangleGradient(v1, v2, v3, c1, c2, c3);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `c1` | `Color` |
| `c2` | `Color` |
| `c3` | `Color` |

---

### DrawTriangleLines

```c
void DrawTriangleLines(Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda usage**

```koda
DrawTriangleLines(v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### DrawTriangleFan

```c
void DrawTriangleFan(const Vector2* points, int pointCount, Color color)
```

**Koda usage**

```koda
DrawTriangleFan(points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawTriangleStrip

```c
void DrawTriangleStrip(const Vector2* points, int pointCount, Color color)
```

**Koda usage**

```koda
DrawTriangleStrip(points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawPoly

```c
void DrawPoly(Vector2 center, int sides, float radius, float rotation, Color color)
```

**Koda usage**

```koda
DrawPoly(center, sides, radius, rotation, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `sides` | `int` |
| `radius` | `float` |
| `rotation` | `float` |
| `color` | `Color` |

---

### DrawPolyLines

```c
void DrawPolyLines(Vector2 center, int sides, float radius, float rotation, Color color)
```

**Koda usage**

```koda
DrawPolyLines(center, sides, radius, rotation, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `sides` | `int` |
| `radius` | `float` |
| `rotation` | `float` |
| `color` | `Color` |

---

### DrawPolyLinesEx

```c
void DrawPolyLinesEx(Vector2 center, int sides, float radius, float rotation, float lineThick, Color color)
```

**Koda usage**

```koda
DrawPolyLinesEx(center, sides, radius, rotation, lineThick, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `sides` | `int` |
| `radius` | `float` |
| `rotation` | `float` |
| `lineThick` | `float` |
| `color` | `Color` |

---

### DrawSplineLinear

```c
void DrawSplineLinear(const Vector2* points, int pointCount, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineLinear(points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineBasis

```c
void DrawSplineBasis(const Vector2* points, int pointCount, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineBasis(points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineCatmullRom

```c
void DrawSplineCatmullRom(const Vector2* points, int pointCount, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineCatmullRom(points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineBezierQuadratic

```c
void DrawSplineBezierQuadratic(const Vector2* points, int pointCount, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineBezierQuadratic(points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineBezierCubic

```c
void DrawSplineBezierCubic(const Vector2* points, int pointCount, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineBezierCubic(points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentLinear

```c
void DrawSplineSegmentLinear(Vector2 p1, Vector2 p2, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineSegmentLinear(p1, p2, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentBasis

```c
void DrawSplineSegmentBasis(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineSegmentBasis(p1, p2, p3, p4, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentCatmullRom

```c
void DrawSplineSegmentCatmullRom(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineSegmentCatmullRom(p1, p2, p3, p4, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentBezierQuadratic

```c
void DrawSplineSegmentBezierQuadratic(Vector2 p1, Vector2 c2, Vector2 p3, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineSegmentBezierQuadratic(p1, c2, p3, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `p3` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentBezierCubic

```c
void DrawSplineSegmentBezierCubic(Vector2 p1, Vector2 c2, Vector2 c3, Vector2 p4, float thick, Color color)
```

**Koda usage**

```koda
DrawSplineSegmentBezierCubic(p1, c2, c3, p4, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `c3` | `Vector2` |
| `p4` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### GetSplinePointLinear

```c
Vector2 GetSplinePointLinear(Vector2 startPos, Vector2 endPos, float t)
```

**Koda usage**

```koda
let result = GetSplinePointLinear(startPos, endPos, t);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointBasis

```c
Vector2 GetSplinePointBasis(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float t)
```

**Koda usage**

```koda
let result = GetSplinePointBasis(p1, p2, p3, p4, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointCatmullRom

```c
Vector2 GetSplinePointCatmullRom(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float t)
```

**Koda usage**

```koda
let result = GetSplinePointCatmullRom(p1, p2, p3, p4, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointBezierQuadratic

```c
Vector2 GetSplinePointBezierQuadratic(Vector2 p1, Vector2 c2, Vector2 p3, float t)
```

**Koda usage**

```koda
let result = GetSplinePointBezierQuadratic(p1, c2, p3, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `p3` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointBezierCubic

```c
Vector2 GetSplinePointBezierCubic(Vector2 p1, Vector2 c2, Vector2 c3, Vector2 p4, float t)
```

**Koda usage**

```koda
let result = GetSplinePointBezierCubic(p1, c2, c3, p4, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `c3` | `Vector2` |
| `p4` | `Vector2` |
| `t` | `float` |

---

### CheckCollisionRecs

```c
bool CheckCollisionRecs(Rectangle rec1, Rectangle rec2)
```

**Koda usage**

```koda
let result = CheckCollisionRecs(rec1, rec2);
```

| Parameter | Type |
|-----------|------|
| `rec1` | `Rectangle` |
| `rec2` | `Rectangle` |

---

### CheckCollisionCircles

```c
bool CheckCollisionCircles(Vector2 center1, float radius1, Vector2 center2, float radius2)
```

**Koda usage**

```koda
let result = CheckCollisionCircles(center1, radius1, center2, radius2);
```

| Parameter | Type |
|-----------|------|
| `center1` | `Vector2` |
| `radius1` | `float` |
| `center2` | `Vector2` |
| `radius2` | `float` |

---

### CheckCollisionCircleRec

```c
bool CheckCollisionCircleRec(Vector2 center, float radius, Rectangle rec)
```

**Koda usage**

```koda
let result = CheckCollisionCircleRec(center, radius, rec);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `rec` | `Rectangle` |

---

### CheckCollisionCircleLine

```c
bool CheckCollisionCircleLine(Vector2 center, float radius, Vector2 p1, Vector2 p2)
```

**Koda usage**

```koda
let result = CheckCollisionCircleLine(center, radius, p1, p2);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `p1` | `Vector2` |
| `p2` | `Vector2` |

---

### CheckCollisionPointRec

```c
bool CheckCollisionPointRec(Vector2 point, Rectangle rec)
```

**Koda usage**

```koda
let result = CheckCollisionPointRec(point, rec);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `rec` | `Rectangle` |

---

### CheckCollisionPointCircle

```c
bool CheckCollisionPointCircle(Vector2 point, Vector2 center, float radius)
```

**Koda usage**

```koda
let result = CheckCollisionPointCircle(point, center, radius);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `center` | `Vector2` |
| `radius` | `float` |

---

### CheckCollisionPointTriangle

```c
bool CheckCollisionPointTriangle(Vector2 point, Vector2 p1, Vector2 p2, Vector2 p3)
```

**Koda usage**

```koda
let result = CheckCollisionPointTriangle(point, p1, p2, p3);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |

---

### CheckCollisionPointLine

```c
bool CheckCollisionPointLine(Vector2 point, Vector2 p1, Vector2 p2, int threshold)
```

**Koda usage**

```koda
let result = CheckCollisionPointLine(point, p1, p2, threshold);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `threshold` | `int` |

---

### CheckCollisionPointPoly

```c
bool CheckCollisionPointPoly(Vector2 point, const Vector2* points, int pointCount)
```

**Koda usage**

```koda
let result = CheckCollisionPointPoly(point, points, pointCount);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `points` | `const Vector2*` |
| `pointCount` | `int` |

---

### CheckCollisionLines

```c
bool CheckCollisionLines(Vector2 startPos1, Vector2 endPos1, Vector2 startPos2, Vector2 endPos2, Vector2* collisionPoint)
```

**Koda usage**

```koda
let result = CheckCollisionLines(startPos1, endPos1, startPos2, endPos2, collisionPoint);
```

| Parameter | Type |
|-----------|------|
| `startPos1` | `Vector2` |
| `endPos1` | `Vector2` |
| `startPos2` | `Vector2` |
| `endPos2` | `Vector2` |
| `collisionPoint` | `Vector2*` |

---

### GetCollisionRec

```c
Rectangle GetCollisionRec(Rectangle rec1, Rectangle rec2)
```

**Koda usage**

```koda
let result = GetCollisionRec(rec1, rec2);
```

| Parameter | Type |
|-----------|------|
| `rec1` | `Rectangle` |
| `rec2` | `Rectangle` |

---

### LoadImage

```c
Image LoadImage(const char* fileName)
```

**Koda usage**

```koda
let result = LoadImage(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadImageRaw

```c
Image LoadImageRaw(const char* fileName, int width, int height, int format, int headerSize)
```

**Koda usage**

```koda
let result = LoadImageRaw(fileName, width, height, format, headerSize);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `width` | `int` |
| `height` | `int` |
| `format` | `int` |
| `headerSize` | `int` |

---

### LoadImageAnim

```c
Image LoadImageAnim(const char* fileName, int* frames)
```

**Koda usage**

```koda
let result = LoadImageAnim(fileName, frames);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `frames` | `int*` |

---

### LoadImageAnimFromMemory

```c
Image LoadImageAnimFromMemory(const char* fileType, const unsigned char* fileData, int dataSize, int* frames)
```

**Koda usage**

```koda
let result = LoadImageAnimFromMemory(fileType, fileData, dataSize, frames);
```

| Parameter | Type |
|-----------|------|
| `fileType` | `const char*` |
| `fileData` | `const unsigned char*` |
| `dataSize` | `int` |
| `frames` | `int*` |

---

### LoadImageFromMemory

```c
Image LoadImageFromMemory(const char* fileType, const unsigned char* fileData, int dataSize)
```

**Koda usage**

```koda
let result = LoadImageFromMemory(fileType, fileData, dataSize);
```

| Parameter | Type |
|-----------|------|
| `fileType` | `const char*` |
| `fileData` | `const unsigned char*` |
| `dataSize` | `int` |

---

### LoadImageFromTexture

```c
Image LoadImageFromTexture(Texture2D texture)
```

**Koda usage**

```koda
let result = LoadImageFromTexture(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |

---

### LoadImageFromScreen

```c
Image LoadImageFromScreen()
```

**Koda usage**

```koda
let result = LoadImageFromScreen();
```

---

### IsImageValid

```c
bool IsImageValid(Image image)
```

**Koda usage**

```koda
let result = IsImageValid(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### UnloadImage

```c
void UnloadImage(Image image)
```

**Koda usage**

```koda
UnloadImage(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### ExportImage

```c
bool ExportImage(Image image, const char* fileName)
```

**Koda usage**

```koda
let result = ExportImage(image, fileName);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `fileName` | `const char*` |

---

### ExportImageAsCode

```c
bool ExportImageAsCode(Image image, const char* fileName)
```

**Koda usage**

```koda
let result = ExportImageAsCode(image, fileName);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `fileName` | `const char*` |

---

### GenImageColor

```c
Image GenImageColor(int width, int height, Color color)
```

**Koda usage**

```koda
let result = GenImageColor(width, height, color);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### GenImageGradientLinear

```c
Image GenImageGradientLinear(int width, int height, int direction, Color start, Color end)
```

**Koda usage**

```koda
let result = GenImageGradientLinear(width, height, direction, start, end);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `direction` | `int` |
| `start` | `Color` |
| `end` | `Color` |

---

### GenImageGradientRadial

```c
Image GenImageGradientRadial(int width, int height, float density, Color inner, Color outer)
```

**Koda usage**

```koda
let result = GenImageGradientRadial(width, height, density, inner, outer);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `density` | `float` |
| `inner` | `Color` |
| `outer` | `Color` |

---

### GenImageGradientSquare

```c
Image GenImageGradientSquare(int width, int height, float density, Color inner, Color outer)
```

**Koda usage**

```koda
let result = GenImageGradientSquare(width, height, density, inner, outer);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `density` | `float` |
| `inner` | `Color` |
| `outer` | `Color` |

---

### GenImageChecked

```c
Image GenImageChecked(int width, int height, int checksX, int checksY, Color col1, Color col2)
```

**Koda usage**

```koda
let result = GenImageChecked(width, height, checksX, checksY, col1, col2);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `checksX` | `int` |
| `checksY` | `int` |
| `col1` | `Color` |
| `col2` | `Color` |

---

### GenImageWhiteNoise

```c
Image GenImageWhiteNoise(int width, int height, float factor)
```

**Koda usage**

```koda
let result = GenImageWhiteNoise(width, height, factor);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `factor` | `float` |

---

### GenImagePerlinNoise

```c
Image GenImagePerlinNoise(int width, int height, int offsetX, int offsetY, float scale)
```

**Koda usage**

```koda
let result = GenImagePerlinNoise(width, height, offsetX, offsetY, scale);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `offsetX` | `int` |
| `offsetY` | `int` |
| `scale` | `float` |

---

### GenImageCellular

```c
Image GenImageCellular(int width, int height, int tileSize)
```

**Koda usage**

```koda
let result = GenImageCellular(width, height, tileSize);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `tileSize` | `int` |

---

### GenImageText

```c
Image GenImageText(int width, int height, const char* text)
```

**Koda usage**

```koda
let result = GenImageText(width, height, text);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `text` | `const char*` |

---

### ImageCopy

```c
Image ImageCopy(Image image)
```

**Koda usage**

```koda
let result = ImageCopy(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### ImageFromImage

```c
Image ImageFromImage(Image image, Rectangle rec)
```

**Koda usage**

```koda
let result = ImageFromImage(image, rec);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `rec` | `Rectangle` |

---

### ImageFromChannel

```c
Image ImageFromChannel(Image image, int selectedChannel)
```

**Koda usage**

```koda
let result = ImageFromChannel(image, selectedChannel);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `selectedChannel` | `int` |

---

### ImageText

```c
Image ImageText(const char* text, int fontSize, Color color)
```

**Koda usage**

```koda
let result = ImageText(text, fontSize, color);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `fontSize` | `int` |
| `color` | `Color` |

---

### ImageTextEx

```c
Image ImageTextEx(Font font, const char* text, float fontSize, float spacing, Color tint)
```

**Koda usage**

```koda
let result = ImageTextEx(font, text, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `text` | `const char*` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### ImageFormat

```c
void ImageFormat(Image* image, int newFormat)
```

**Koda usage**

```koda
ImageFormat(image, newFormat);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `newFormat` | `int` |

---

### ImageToPOT

```c
void ImageToPOT(Image* image, Color fill)
```

**Koda usage**

```koda
ImageToPOT(image, fill);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `fill` | `Color` |

---

### ImageCrop

```c
void ImageCrop(Image* image, Rectangle crop)
```

**Koda usage**

```koda
ImageCrop(image, crop);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `crop` | `Rectangle` |

---

### ImageAlphaCrop

```c
void ImageAlphaCrop(Image* image, float threshold)
```

**Koda usage**

```koda
ImageAlphaCrop(image, threshold);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `threshold` | `float` |

---

### ImageAlphaClear

```c
void ImageAlphaClear(Image* image, Color color, float threshold)
```

**Koda usage**

```koda
ImageAlphaClear(image, color, threshold);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `color` | `Color` |
| `threshold` | `float` |

---

### ImageAlphaMask

```c
void ImageAlphaMask(Image* image, Image alphaMask)
```

**Koda usage**

```koda
ImageAlphaMask(image, alphaMask);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `alphaMask` | `Image` |

---

### ImageAlphaPremultiply

```c
void ImageAlphaPremultiply(Image* image)
```

**Koda usage**

```koda
ImageAlphaPremultiply(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageBlurGaussian

```c
void ImageBlurGaussian(Image* image, int blurSize)
```

**Koda usage**

```koda
ImageBlurGaussian(image, blurSize);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `blurSize` | `int` |

---

### ImageKernelConvolution

```c
void ImageKernelConvolution(Image* image, const float* kernel, int kernelSize)
```

**Koda usage**

```koda
ImageKernelConvolution(image, kernel, kernelSize);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `kernel` | `const float*` |
| `kernelSize` | `int` |

---

### ImageResize

```c
void ImageResize(Image* image, int newWidth, int newHeight)
```

**Koda usage**

```koda
ImageResize(image, newWidth, newHeight);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `newWidth` | `int` |
| `newHeight` | `int` |

---

### ImageResizeNN

```c
void ImageResizeNN(Image* image, int newWidth, int newHeight)
```

**Koda usage**

```koda
ImageResizeNN(image, newWidth, newHeight);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `newWidth` | `int` |
| `newHeight` | `int` |

---

### ImageResizeCanvas

```c
void ImageResizeCanvas(Image* image, int newWidth, int newHeight, int offsetX, int offsetY, Color fill)
```

**Koda usage**

```koda
ImageResizeCanvas(image, newWidth, newHeight, offsetX, offsetY, fill);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `newWidth` | `int` |
| `newHeight` | `int` |
| `offsetX` | `int` |
| `offsetY` | `int` |
| `fill` | `Color` |

---

### ImageMipmaps

```c
void ImageMipmaps(Image* image)
```

**Koda usage**

```koda
ImageMipmaps(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageDither

```c
void ImageDither(Image* image, int rBpp, int gBpp, int bBpp, int aBpp)
```

**Koda usage**

```koda
ImageDither(image, rBpp, gBpp, bBpp, aBpp);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `rBpp` | `int` |
| `gBpp` | `int` |
| `bBpp` | `int` |
| `aBpp` | `int` |

---

### ImageFlipVertical

```c
void ImageFlipVertical(Image* image)
```

**Koda usage**

```koda
ImageFlipVertical(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageFlipHorizontal

```c
void ImageFlipHorizontal(Image* image)
```

**Koda usage**

```koda
ImageFlipHorizontal(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageRotate

```c
void ImageRotate(Image* image, int degrees)
```

**Koda usage**

```koda
ImageRotate(image, degrees);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `degrees` | `int` |

---

### ImageRotateCW

```c
void ImageRotateCW(Image* image)
```

**Koda usage**

```koda
ImageRotateCW(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageRotateCCW

```c
void ImageRotateCCW(Image* image)
```

**Koda usage**

```koda
ImageRotateCCW(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageColorTint

```c
void ImageColorTint(Image* image, Color color)
```

**Koda usage**

```koda
ImageColorTint(image, color);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `color` | `Color` |

---

### ImageColorInvert

```c
void ImageColorInvert(Image* image)
```

**Koda usage**

```koda
ImageColorInvert(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageColorGrayscale

```c
void ImageColorGrayscale(Image* image)
```

**Koda usage**

```koda
ImageColorGrayscale(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |

---

### ImageColorContrast

```c
void ImageColorContrast(Image* image, int contrast)
```

**Koda usage**

```koda
ImageColorContrast(image, contrast);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `contrast` | `int` |

---

### ImageColorBrightness

```c
void ImageColorBrightness(Image* image, int brightness)
```

**Koda usage**

```koda
ImageColorBrightness(image, brightness);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `brightness` | `int` |

---

### ImageColorReplace

```c
void ImageColorReplace(Image* image, Color color, Color replace)
```

**Koda usage**

```koda
ImageColorReplace(image, color, replace);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image*` |
| `color` | `Color` |
| `replace` | `Color` |

---

### UnloadImageColors

```c
void UnloadImageColors(Color* colors)
```

**Koda usage**

```koda
UnloadImageColors(colors);
```

| Parameter | Type |
|-----------|------|
| `colors` | `Color*` |

---

### UnloadImagePalette

```c
void UnloadImagePalette(Color* colors)
```

**Koda usage**

```koda
UnloadImagePalette(colors);
```

| Parameter | Type |
|-----------|------|
| `colors` | `Color*` |

---

### GetImageAlphaBorder

```c
Rectangle GetImageAlphaBorder(Image image, float threshold)
```

**Koda usage**

```koda
let result = GetImageAlphaBorder(image, threshold);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `threshold` | `float` |

---

### GetImageColor

```c
Color GetImageColor(Image image, int x, int y)
```

**Koda usage**

```koda
let result = GetImageColor(image, x, y);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `x` | `int` |
| `y` | `int` |

---

### ImageClearBackground

```c
void ImageClearBackground(Image* dst, Color color)
```

**Koda usage**

```koda
ImageClearBackground(dst, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `color` | `Color` |

---

### ImageDrawPixel

```c
void ImageDrawPixel(Image* dst, int posX, int posY, Color color)
```

**Koda usage**

```koda
ImageDrawPixel(dst, posX, posY, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `posX` | `int` |
| `posY` | `int` |
| `color` | `Color` |

---

### ImageDrawPixelV

```c
void ImageDrawPixelV(Image* dst, Vector2 position, Color color)
```

**Koda usage**

```koda
ImageDrawPixelV(dst, position, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `position` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawLine

```c
void ImageDrawLine(Image* dst, int startPosX, int startPosY, int endPosX, int endPosY, Color color)
```

**Koda usage**

```koda
ImageDrawLine(dst, startPosX, startPosY, endPosX, endPosY, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `startPosX` | `int` |
| `startPosY` | `int` |
| `endPosX` | `int` |
| `endPosY` | `int` |
| `color` | `Color` |

---

### ImageDrawLineV

```c
void ImageDrawLineV(Image* dst, Vector2 start, Vector2 end, Color color)
```

**Koda usage**

```koda
ImageDrawLineV(dst, start, end, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `start` | `Vector2` |
| `end` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawLineEx

```c
void ImageDrawLineEx(Image* dst, Vector2 start, Vector2 end, int thick, Color color)
```

**Koda usage**

```koda
ImageDrawLineEx(dst, start, end, thick, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `start` | `Vector2` |
| `end` | `Vector2` |
| `thick` | `int` |
| `color` | `Color` |

---

### ImageDrawCircle

```c
void ImageDrawCircle(Image* dst, int centerX, int centerY, int radius, Color color)
```

**Koda usage**

```koda
ImageDrawCircle(dst, centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawCircleV

```c
void ImageDrawCircleV(Image* dst, Vector2 center, int radius, Color color)
```

**Koda usage**

```koda
ImageDrawCircleV(dst, center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `center` | `Vector2` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawCircleLines

```c
void ImageDrawCircleLines(Image* dst, int centerX, int centerY, int radius, Color color)
```

**Koda usage**

```koda
ImageDrawCircleLines(dst, centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawCircleLinesV

```c
void ImageDrawCircleLinesV(Image* dst, Vector2 center, int radius, Color color)
```

**Koda usage**

```koda
ImageDrawCircleLinesV(dst, center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `center` | `Vector2` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawRectangle

```c
void ImageDrawRectangle(Image* dst, int posX, int posY, int width, int height, Color color)
```

**Koda usage**

```koda
ImageDrawRectangle(dst, posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### ImageDrawRectangleV

```c
void ImageDrawRectangleV(Image* dst, Vector2 position, Vector2 size, Color color)
```

**Koda usage**

```koda
ImageDrawRectangleV(dst, position, size, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `position` | `Vector2` |
| `size` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawRectangleRec

```c
void ImageDrawRectangleRec(Image* dst, Rectangle rec, Color color)
```

**Koda usage**

```koda
ImageDrawRectangleRec(dst, rec, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `rec` | `Rectangle` |
| `color` | `Color` |

---

### ImageDrawRectangleLines

```c
void ImageDrawRectangleLines(Image* dst, int posX, int posY, int width, int height, Color color)
```

**Koda usage**

```koda
ImageDrawRectangleLines(dst, posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### ImageDrawRectangleLinesEx

```c
void ImageDrawRectangleLinesEx(Image* dst, Rectangle rec, int thick, Color color)
```

**Koda usage**

```koda
ImageDrawRectangleLinesEx(dst, rec, thick, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `rec` | `Rectangle` |
| `thick` | `int` |
| `color` | `Color` |

---

### ImageDrawTriangle

```c
void ImageDrawTriangle(Image* dst, Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda usage**

```koda
ImageDrawTriangle(dst, v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawTriangleGradient

```c
void ImageDrawTriangleGradient(Image* dst, Vector2 v1, Vector2 v2, Vector2 v3, Color c1, Color c2, Color c3)
```

**Koda usage**

```koda
ImageDrawTriangleGradient(dst, v1, v2, v3, c1, c2, c3);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `c1` | `Color` |
| `c2` | `Color` |
| `c3` | `Color` |

---

### ImageDrawTriangleLines

```c
void ImageDrawTriangleLines(Image* dst, Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda usage**

```koda
ImageDrawTriangleLines(dst, v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawTriangleFan

```c
void ImageDrawTriangleFan(Image* dst, const Vector2* points, int pointCount, Color color)
```

**Koda usage**

```koda
ImageDrawTriangleFan(dst, points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### ImageDrawTriangleStrip

```c
void ImageDrawTriangleStrip(Image* dst, const Vector2* points, int pointCount, Color color)
```

**Koda usage**

```koda
ImageDrawTriangleStrip(dst, points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `points` | `const Vector2*` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### ImageDraw

```c
void ImageDraw(Image* dst, Image src, Rectangle srcRec, Rectangle dstRec, Color tint)
```

**Koda usage**

```koda
ImageDraw(dst, src, srcRec, dstRec, tint);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `src` | `Image` |
| `srcRec` | `Rectangle` |
| `dstRec` | `Rectangle` |
| `tint` | `Color` |

---

### ImageDrawText

```c
void ImageDrawText(Image* dst, const char* text, int posX, int posY, int fontSize, Color color)
```

**Koda usage**

```koda
ImageDrawText(dst, text, posX, posY, fontSize, color);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `text` | `const char*` |
| `posX` | `int` |
| `posY` | `int` |
| `fontSize` | `int` |
| `color` | `Color` |

---

### ImageDrawTextEx

```c
void ImageDrawTextEx(Image* dst, Font font, const char* text, Vector2 position, float fontSize, float spacing, Color tint)
```

**Koda usage**

```koda
ImageDrawTextEx(dst, font, text, position, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Image*` |
| `font` | `Font` |
| `text` | `const char*` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### LoadTexture

```c
Texture2D LoadTexture(const char* fileName)
```

**Koda usage**

```koda
let result = LoadTexture(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadTextureFromImage

```c
Texture2D LoadTextureFromImage(Image image)
```

**Koda usage**

```koda
let result = LoadTextureFromImage(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### LoadTextureCubemap

```c
TextureCubemap LoadTextureCubemap(Image image, int layout)
```

**Koda usage**

```koda
let result = LoadTextureCubemap(image, layout);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `layout` | `int` |

---

### LoadRenderTexture

```c
RenderTexture2D LoadRenderTexture(int width, int height)
```

**Koda usage**

```koda
let result = LoadRenderTexture(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### IsTextureValid

```c
bool IsTextureValid(Texture2D texture)
```

**Koda usage**

```koda
let result = IsTextureValid(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |

---

### UnloadTexture

```c
void UnloadTexture(Texture2D texture)
```

**Koda usage**

```koda
UnloadTexture(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |

---

### IsRenderTextureValid

```c
bool IsRenderTextureValid(RenderTexture2D target)
```

**Koda usage**

```koda
let result = IsRenderTextureValid(target);
```

| Parameter | Type |
|-----------|------|
| `target` | `RenderTexture2D` |

---

### UnloadRenderTexture

```c
void UnloadRenderTexture(RenderTexture2D target)
```

**Koda usage**

```koda
UnloadRenderTexture(target);
```

| Parameter | Type |
|-----------|------|
| `target` | `RenderTexture2D` |

---

### UpdateTexture

```c
void UpdateTexture(Texture2D texture, const void* pixels)
```

**Koda usage**

```koda
UpdateTexture(texture, pixels);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `pixels` | `const void*` |

---

### UpdateTextureRec

```c
void UpdateTextureRec(Texture2D texture, Rectangle rec, const void* pixels)
```

**Koda usage**

```koda
UpdateTextureRec(texture, rec, pixels);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `rec` | `Rectangle` |
| `pixels` | `const void*` |

---

### GenTextureMipmaps

```c
void GenTextureMipmaps(Texture2D* texture)
```

**Koda usage**

```koda
GenTextureMipmaps(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D*` |

---

### SetTextureFilter

```c
void SetTextureFilter(Texture2D texture, int filter)
```

**Koda usage**

```koda
SetTextureFilter(texture, filter);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `filter` | `int` |

---

### SetTextureWrap

```c
void SetTextureWrap(Texture2D texture, int wrap)
```

**Koda usage**

```koda
SetTextureWrap(texture, wrap);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `wrap` | `int` |

---

### DrawTexture

```c
void DrawTexture(Texture2D texture, int posX, int posY, Color tint)
```

**Koda usage**

```koda
DrawTexture(texture, posX, posY, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `posX` | `int` |
| `posY` | `int` |
| `tint` | `Color` |

---

### DrawTextureV

```c
void DrawTextureV(Texture2D texture, Vector2 position, Color tint)
```

**Koda usage**

```koda
DrawTextureV(texture, position, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `position` | `Vector2` |
| `tint` | `Color` |

---

### DrawTextureEx

```c
void DrawTextureEx(Texture2D texture, Vector2 position, float rotation, float scale, Color tint)
```

**Koda usage**

```koda
DrawTextureEx(texture, position, rotation, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `position` | `Vector2` |
| `rotation` | `float` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawTextureRec

```c
void DrawTextureRec(Texture2D texture, Rectangle source, Vector2 position, Color tint)
```

**Koda usage**

```koda
DrawTextureRec(texture, source, position, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `position` | `Vector2` |
| `tint` | `Color` |

---

### DrawTexturePro

```c
void DrawTexturePro(Texture2D texture, Rectangle source, Rectangle dest, Vector2 origin, float rotation, Color tint)
```

**Koda usage**

```koda
DrawTexturePro(texture, source, dest, origin, rotation, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `dest` | `Rectangle` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `tint` | `Color` |

---

### DrawTextureNPatch

```c
void DrawTextureNPatch(Texture2D texture, NPatchInfo nPatchInfo, Rectangle dest, Vector2 origin, float rotation, Color tint)
```

**Koda usage**

```koda
DrawTextureNPatch(texture, nPatchInfo, dest, origin, rotation, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `nPatchInfo` | `NPatchInfo` |
| `dest` | `Rectangle` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `tint` | `Color` |

---

### ColorIsEqual

```c
bool ColorIsEqual(Color col1, Color col2)
```

**Koda usage**

```koda
let result = ColorIsEqual(col1, col2);
```

| Parameter | Type |
|-----------|------|
| `col1` | `Color` |
| `col2` | `Color` |

---

### Fade

```c
Color Fade(Color color, float alpha)
```

**Koda usage**

```koda
let result = Fade(color, alpha);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `alpha` | `float` |

---

### ColorToInt

```c
int ColorToInt(Color color)
```

**Koda usage**

```koda
let result = ColorToInt(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### ColorNormalize

```c
Vector4 ColorNormalize(Color color)
```

**Koda usage**

```koda
let result = ColorNormalize(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### ColorFromNormalized

```c
Color ColorFromNormalized(Vector4 normalized)
```

**Koda usage**

```koda
let result = ColorFromNormalized(normalized);
```

| Parameter | Type |
|-----------|------|
| `normalized` | `Vector4` |

---

### ColorToHSV

```c
Vector3 ColorToHSV(Color color)
```

**Koda usage**

```koda
let result = ColorToHSV(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### ColorFromHSV

```c
Color ColorFromHSV(float hue, float saturation, float value)
```

**Koda usage**

```koda
let result = ColorFromHSV(hue, saturation, value);
```

| Parameter | Type |
|-----------|------|
| `hue` | `float` |
| `saturation` | `float` |
| `value` | `float` |

---

### ColorTint

```c
Color ColorTint(Color color, Color tint)
```

**Koda usage**

```koda
let result = ColorTint(color, tint);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `tint` | `Color` |

---

### ColorBrightness

```c
Color ColorBrightness(Color color, float factor)
```

**Koda usage**

```koda
let result = ColorBrightness(color, factor);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `factor` | `float` |

---

### ColorContrast

```c
Color ColorContrast(Color color, float contrast)
```

**Koda usage**

```koda
let result = ColorContrast(color, contrast);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `contrast` | `float` |

---

### ColorAlpha

```c
Color ColorAlpha(Color color, float alpha)
```

**Koda usage**

```koda
let result = ColorAlpha(color, alpha);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `alpha` | `float` |

---

### ColorAlphaBlend

```c
Color ColorAlphaBlend(Color dst, Color src, Color tint)
```

**Koda usage**

```koda
let result = ColorAlphaBlend(dst, src, tint);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Color` |
| `src` | `Color` |
| `tint` | `Color` |

---

### ColorLerp

```c
Color ColorLerp(Color color1, Color color2, float factor)
```

**Koda usage**

```koda
let result = ColorLerp(color1, color2, factor);
```

| Parameter | Type |
|-----------|------|
| `color1` | `Color` |
| `color2` | `Color` |
| `factor` | `float` |

---

### GetColor

```c
Color GetColor(unsigned int hexValue)
```

**Koda usage**

```koda
let result = GetColor(hexValue);
```

| Parameter | Type |
|-----------|------|
| `hexValue` | `unsigned int` |

---

### GetPixelColor

```c
Color GetPixelColor(void* srcPtr, int format)
```

**Koda usage**

```koda
let result = GetPixelColor(srcPtr, format);
```

| Parameter | Type |
|-----------|------|
| `srcPtr` | `void*` |
| `format` | `int` |

---

### SetPixelColor

```c
void SetPixelColor(void* dstPtr, Color color, int format)
```

**Koda usage**

```koda
SetPixelColor(dstPtr, color, format);
```

| Parameter | Type |
|-----------|------|
| `dstPtr` | `void*` |
| `color` | `Color` |
| `format` | `int` |

---

### GetPixelDataSize

```c
int GetPixelDataSize(int width, int height, int format)
```

**Koda usage**

```koda
let result = GetPixelDataSize(width, height, format);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `format` | `int` |

---

### GetFontDefault

```c
Font GetFontDefault()
```

**Koda usage**

```koda
let result = GetFontDefault();
```

---

### LoadFont

```c
Font LoadFont(const char* fileName)
```

**Koda usage**

```koda
let result = LoadFont(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadFontEx

```c
Font LoadFontEx(const char* fileName, int fontSize, const int* codepoints, int codepointCount)
```

**Koda usage**

```koda
let result = LoadFontEx(fileName, fontSize, codepoints, codepointCount);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |
| `fontSize` | `int` |
| `codepoints` | `const int*` |
| `codepointCount` | `int` |

---

### LoadFontFromImage

```c
Font LoadFontFromImage(Image image, Color key, int firstChar)
```

**Koda usage**

```koda
let result = LoadFontFromImage(image, key, firstChar);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `key` | `Color` |
| `firstChar` | `int` |

---

### LoadFontFromMemory

```c
Font LoadFontFromMemory(const char* fileType, const unsigned char* fileData, int dataSize, int fontSize, const int* codepoints, int codepointCount)
```

**Koda usage**

```koda
let result = LoadFontFromMemory(fileType, fileData, dataSize, fontSize, codepoints, codepointCount);
```

| Parameter | Type |
|-----------|------|
| `fileType` | `const char*` |
| `fileData` | `const unsigned char*` |
| `dataSize` | `int` |
| `fontSize` | `int` |
| `codepoints` | `const int*` |
| `codepointCount` | `int` |

---

### IsFontValid

```c
bool IsFontValid(Font font)
```

**Koda usage**

```koda
let result = IsFontValid(font);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |

---

### GenImageFontAtlas

```c
Image GenImageFontAtlas(const GlyphInfo* glyphs, Rectangle** glyphRecs, int glyphCount, int fontSize, int padding, int packMethod)
```

**Koda usage**

```koda
let result = GenImageFontAtlas(glyphs, glyphRecs, glyphCount, fontSize, padding, packMethod);
```

| Parameter | Type |
|-----------|------|
| `glyphs` | `const GlyphInfo*` |
| `glyphRecs` | `Rectangle**` |
| `glyphCount` | `int` |
| `fontSize` | `int` |
| `padding` | `int` |
| `packMethod` | `int` |

---

### UnloadFontData

```c
void UnloadFontData(GlyphInfo* glyphs, int glyphCount)
```

**Koda usage**

```koda
UnloadFontData(glyphs, glyphCount);
```

| Parameter | Type |
|-----------|------|
| `glyphs` | `GlyphInfo*` |
| `glyphCount` | `int` |

---

### UnloadFont

```c
void UnloadFont(Font font)
```

**Koda usage**

```koda
UnloadFont(font);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |

---

### ExportFontAsCode

```c
bool ExportFontAsCode(Font font, const char* fileName)
```

**Koda usage**

```koda
let result = ExportFontAsCode(font, fileName);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `fileName` | `const char*` |

---

### DrawFPS

```c
void DrawFPS(int posX, int posY)
```

**Koda usage**

```koda
DrawFPS(posX, posY);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |

---

### DrawText

```c
void DrawText(const char* text, int posX, int posY, int fontSize, Color color)
```

**Koda usage**

```koda
DrawText(text, posX, posY, fontSize, color);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `posX` | `int` |
| `posY` | `int` |
| `fontSize` | `int` |
| `color` | `Color` |

---

### DrawTextEx

```c
void DrawTextEx(Font font, const char* text, Vector2 position, float fontSize, float spacing, Color tint)
```

**Koda usage**

```koda
DrawTextEx(font, text, position, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `text` | `const char*` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### DrawTextPro

```c
void DrawTextPro(Font font, const char* text, Vector2 position, Vector2 origin, float rotation, float fontSize, float spacing, Color tint)
```

**Koda usage**

```koda
DrawTextPro(font, text, position, origin, rotation, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `text` | `const char*` |
| `position` | `Vector2` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### DrawTextCodepoint

```c
void DrawTextCodepoint(Font font, int codepoint, Vector2 position, float fontSize, Color tint)
```

**Koda usage**

```koda
DrawTextCodepoint(font, codepoint, position, fontSize, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `tint` | `Color` |

---

### DrawTextCodepoints

```c
void DrawTextCodepoints(Font font, const int* codepoints, int codepointCount, Vector2 position, float fontSize, float spacing, Color tint)
```

**Koda usage**

```koda
DrawTextCodepoints(font, codepoints, codepointCount, position, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoints` | `const int*` |
| `codepointCount` | `int` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### SetTextLineSpacing

```c
void SetTextLineSpacing(int spacing)
```

**Koda usage**

```koda
SetTextLineSpacing(spacing);
```

| Parameter | Type |
|-----------|------|
| `spacing` | `int` |

---

### MeasureText

```c
int MeasureText(const char* text, int fontSize)
```

**Koda usage**

```koda
let result = MeasureText(text, fontSize);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `fontSize` | `int` |

---

### MeasureTextEx

```c
Vector2 MeasureTextEx(Font font, const char* text, float fontSize, float spacing)
```

**Koda usage**

```koda
let result = MeasureTextEx(font, text, fontSize, spacing);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `text` | `const char*` |
| `fontSize` | `float` |
| `spacing` | `float` |

---

### MeasureTextCodepoints

```c
Vector2 MeasureTextCodepoints(Font font, const int* codepoints, int length, float fontSize, float spacing)
```

**Koda usage**

```koda
let result = MeasureTextCodepoints(font, codepoints, length, fontSize, spacing);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoints` | `const int*` |
| `length` | `int` |
| `fontSize` | `float` |
| `spacing` | `float` |

---

### GetGlyphIndex

```c
int GetGlyphIndex(Font font, int codepoint)
```

**Koda usage**

```koda
let result = GetGlyphIndex(font, codepoint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |

---

### GetGlyphInfo

```c
GlyphInfo GetGlyphInfo(Font font, int codepoint)
```

**Koda usage**

```koda
let result = GetGlyphInfo(font, codepoint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |

---

### GetGlyphAtlasRec

```c
Rectangle GetGlyphAtlasRec(Font font, int codepoint)
```

**Koda usage**

```koda
let result = GetGlyphAtlasRec(font, codepoint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |

---

### UnloadUTF8

```c
void UnloadUTF8(char* text)
```

**Koda usage**

```koda
UnloadUTF8(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `char*` |

---

### UnloadCodepoints

```c
void UnloadCodepoints(int* codepoints)
```

**Koda usage**

```koda
UnloadCodepoints(codepoints);
```

| Parameter | Type |
|-----------|------|
| `codepoints` | `int*` |

---

### GetCodepointCount

```c
int GetCodepointCount(const char* text)
```

**Koda usage**

```koda
let result = GetCodepointCount(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |

---

### GetCodepoint

```c
int GetCodepoint(const char* text, int* codepointSize)
```

**Koda usage**

```koda
let result = GetCodepoint(text, codepointSize);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `codepointSize` | `int*` |

---

### GetCodepointNext

```c
int GetCodepointNext(const char* text, int* codepointSize)
```

**Koda usage**

```koda
let result = GetCodepointNext(text, codepointSize);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `codepointSize` | `int*` |

---

### GetCodepointPrevious

```c
int GetCodepointPrevious(const char* text, int* codepointSize)
```

**Koda usage**

```koda
let result = GetCodepointPrevious(text, codepointSize);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `codepointSize` | `int*` |

---

### UnloadTextLines

```c
void UnloadTextLines(char** text, int lineCount)
```

**Koda usage**

```koda
UnloadTextLines(text, lineCount);
```

| Parameter | Type |
|-----------|------|
| `text` | `char**` |
| `lineCount` | `int` |

---

### TextCopy

```c
int TextCopy(char* dst, const char* src)
```

**Koda usage**

```koda
let result = TextCopy(dst, src);
```

| Parameter | Type |
|-----------|------|
| `dst` | `char*` |
| `src` | `const char*` |

---

### TextIsEqual

```c
bool TextIsEqual(const char* text1, const char* text2)
```

**Koda usage**

```koda
let result = TextIsEqual(text1, text2);
```

| Parameter | Type |
|-----------|------|
| `text1` | `const char*` |
| `text2` | `const char*` |

---

### TextLength

```c
unsigned int TextLength(const char* text)
```

**Koda usage**

```koda
let result = TextLength(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |

---

### TextAppend

```c
void TextAppend(char* text, const char* append, int* position)
```

**Koda usage**

```koda
TextAppend(text, append, position);
```

| Parameter | Type |
|-----------|------|
| `text` | `char*` |
| `append` | `const char*` |
| `position` | `int*` |

---

### TextFindIndex

```c
int TextFindIndex(const char* text, const char* search)
```

**Koda usage**

```koda
let result = TextFindIndex(text, search);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |
| `search` | `const char*` |

---

### TextToInteger

```c
int TextToInteger(const char* text)
```

**Koda usage**

```koda
let result = TextToInteger(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |

---

### TextToFloat

```c
float TextToFloat(const char* text)
```

**Koda usage**

```koda
let result = TextToFloat(text);
```

| Parameter | Type |
|-----------|------|
| `text` | `const char*` |

---

### DrawLine3D

```c
void DrawLine3D(Vector3 startPos, Vector3 endPos, Color color)
```

**Koda usage**

```koda
DrawLine3D(startPos, endPos, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `color` | `Color` |

---

### DrawPoint3D

```c
void DrawPoint3D(Vector3 position, Color color)
```

**Koda usage**

```koda
DrawPoint3D(position, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `color` | `Color` |

---

### DrawCircle3D

```c
void DrawCircle3D(Vector3 center, float radius, Vector3 rotationAxis, float rotationAngle, Color color)
```

**Koda usage**

```koda
DrawCircle3D(center, radius, rotationAxis, rotationAngle, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector3` |
| `radius` | `float` |
| `rotationAxis` | `Vector3` |
| `rotationAngle` | `float` |
| `color` | `Color` |

---

### DrawTriangle3D

```c
void DrawTriangle3D(Vector3 v1, Vector3 v2, Vector3 v3, Color color)
```

**Koda usage**

```koda
DrawTriangle3D(v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector3` |
| `v2` | `Vector3` |
| `v3` | `Vector3` |
| `color` | `Color` |

---

### DrawTriangleStrip3D

```c
void DrawTriangleStrip3D(const Vector3* points, int pointCount, Color color)
```

**Koda usage**

```koda
DrawTriangleStrip3D(points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `points` | `const Vector3*` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawCube

```c
void DrawCube(Vector3 position, float width, float height, float length, Color color)
```

**Koda usage**

```koda
DrawCube(position, width, height, length, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `width` | `float` |
| `height` | `float` |
| `length` | `float` |
| `color` | `Color` |

---

### DrawCubeV

```c
void DrawCubeV(Vector3 position, Vector3 size, Color color)
```

**Koda usage**

```koda
DrawCubeV(position, size, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `size` | `Vector3` |
| `color` | `Color` |

---

### DrawCubeWires

```c
void DrawCubeWires(Vector3 position, float width, float height, float length, Color color)
```

**Koda usage**

```koda
DrawCubeWires(position, width, height, length, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `width` | `float` |
| `height` | `float` |
| `length` | `float` |
| `color` | `Color` |

---

### DrawCubeWiresV

```c
void DrawCubeWiresV(Vector3 position, Vector3 size, Color color)
```

**Koda usage**

```koda
DrawCubeWiresV(position, size, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `size` | `Vector3` |
| `color` | `Color` |

---

### DrawSphere

```c
void DrawSphere(Vector3 centerPos, float radius, Color color)
```

**Koda usage**

```koda
DrawSphere(centerPos, radius, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawSphereEx

```c
void DrawSphereEx(Vector3 centerPos, float radius, int rings, int slices, Color color)
```

**Koda usage**

```koda
DrawSphereEx(centerPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawSphereWires

```c
void DrawSphereWires(Vector3 centerPos, float radius, int rings, int slices, Color color)
```

**Koda usage**

```koda
DrawSphereWires(centerPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCylinder

```c
void DrawCylinder(Vector3 position, float radiusTop, float radiusBottom, float height, int slices, Color color)
```

**Koda usage**

```koda
DrawCylinder(position, radiusTop, radiusBottom, height, slices, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `radiusTop` | `float` |
| `radiusBottom` | `float` |
| `height` | `float` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCylinderEx

```c
void DrawCylinderEx(Vector3 startPos, Vector3 endPos, float startRadius, float endRadius, int sides, Color color)
```

**Koda usage**

```koda
DrawCylinderEx(startPos, endPos, startRadius, endRadius, sides, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `startRadius` | `float` |
| `endRadius` | `float` |
| `sides` | `int` |
| `color` | `Color` |

---

### DrawCylinderWires

```c
void DrawCylinderWires(Vector3 position, float radiusTop, float radiusBottom, float height, int slices, Color color)
```

**Koda usage**

```koda
DrawCylinderWires(position, radiusTop, radiusBottom, height, slices, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `radiusTop` | `float` |
| `radiusBottom` | `float` |
| `height` | `float` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCylinderWiresEx

```c
void DrawCylinderWiresEx(Vector3 startPos, Vector3 endPos, float startRadius, float endRadius, int slices, Color color)
```

**Koda usage**

```koda
DrawCylinderWiresEx(startPos, endPos, startRadius, endRadius, slices, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `startRadius` | `float` |
| `endRadius` | `float` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCapsule

```c
void DrawCapsule(Vector3 startPos, Vector3 endPos, float radius, int rings, int slices, Color color)
```

**Koda usage**

```koda
DrawCapsule(startPos, endPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCapsuleWires

```c
void DrawCapsuleWires(Vector3 startPos, Vector3 endPos, float radius, int rings, int slices, Color color)
```

**Koda usage**

```koda
DrawCapsuleWires(startPos, endPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawPlane

```c
void DrawPlane(Vector3 centerPos, Vector2 size, Color color)
```

**Koda usage**

```koda
DrawPlane(centerPos, size, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `size` | `Vector2` |
| `color` | `Color` |

---

### DrawRay

```c
void DrawRay(Ray ray, Color color)
```

**Koda usage**

```koda
DrawRay(ray, color);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `color` | `Color` |

---

### DrawGrid

```c
void DrawGrid(int slices, float spacing)
```

**Koda usage**

```koda
DrawGrid(slices, spacing);
```

| Parameter | Type |
|-----------|------|
| `slices` | `int` |
| `spacing` | `float` |

---

### LoadModel

```c
Model LoadModel(const char* fileName)
```

**Koda usage**

```koda
let result = LoadModel(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadModelFromMesh

```c
Model LoadModelFromMesh(Mesh mesh)
```

**Koda usage**

```koda
let result = LoadModelFromMesh(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |

---

### IsModelValid

```c
bool IsModelValid(Model model)
```

**Koda usage**

```koda
let result = IsModelValid(model);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |

---

### UnloadModel

```c
void UnloadModel(Model model)
```

**Koda usage**

```koda
UnloadModel(model);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |

---

### GetModelBoundingBox

```c
BoundingBox GetModelBoundingBox(Model model)
```

**Koda usage**

```koda
let result = GetModelBoundingBox(model);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |

---

### DrawModel

```c
void DrawModel(Model model, Vector3 position, float scale, Color tint)
```

**Koda usage**

```koda
DrawModel(model, position, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawModelEx

```c
void DrawModelEx(Model model, Vector3 position, Vector3 rotationAxis, float rotationAngle, Vector3 scale, Color tint)
```

**Koda usage**

```koda
DrawModelEx(model, position, rotationAxis, rotationAngle, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `rotationAxis` | `Vector3` |
| `rotationAngle` | `float` |
| `scale` | `Vector3` |
| `tint` | `Color` |

---

### DrawModelWires

```c
void DrawModelWires(Model model, Vector3 position, float scale, Color tint)
```

**Koda usage**

```koda
DrawModelWires(model, position, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawModelWiresEx

```c
void DrawModelWiresEx(Model model, Vector3 position, Vector3 rotationAxis, float rotationAngle, Vector3 scale, Color tint)
```

**Koda usage**

```koda
DrawModelWiresEx(model, position, rotationAxis, rotationAngle, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `rotationAxis` | `Vector3` |
| `rotationAngle` | `float` |
| `scale` | `Vector3` |
| `tint` | `Color` |

---

### DrawBoundingBox

```c
void DrawBoundingBox(BoundingBox box, Color color)
```

**Koda usage**

```koda
DrawBoundingBox(box, color);
```

| Parameter | Type |
|-----------|------|
| `box` | `BoundingBox` |
| `color` | `Color` |

---

### DrawBillboard

```c
void DrawBillboard(Camera camera, Texture2D texture, Vector3 position, float scale, Color tint)
```

**Koda usage**

```koda
DrawBillboard(camera, texture, position, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |
| `texture` | `Texture2D` |
| `position` | `Vector3` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawBillboardRec

```c
void DrawBillboardRec(Camera camera, Texture2D texture, Rectangle source, Vector3 position, Vector2 size, Color tint)
```

**Koda usage**

```koda
DrawBillboardRec(camera, texture, source, position, size, tint);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `position` | `Vector3` |
| `size` | `Vector2` |
| `tint` | `Color` |

---

### DrawBillboardPro

```c
void DrawBillboardPro(Camera camera, Texture2D texture, Rectangle source, Vector3 position, Vector3 up, Vector2 size, Vector2 origin, float rotation, Color tint)
```

**Koda usage**

```koda
DrawBillboardPro(camera, texture, source, position, up, size, origin, rotation, tint);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `position` | `Vector3` |
| `up` | `Vector3` |
| `size` | `Vector2` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `tint` | `Color` |

---

### UploadMesh

```c
void UploadMesh(Mesh* mesh, bool dynamic)
```

**Koda usage**

```koda
UploadMesh(mesh, dynamic);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh*` |
| `dynamic` | `bool` |

---

### UpdateMeshBuffer

```c
void UpdateMeshBuffer(Mesh mesh, int index, const void* data, int dataSize, int offset)
```

**Koda usage**

```koda
UpdateMeshBuffer(mesh, index, data, dataSize, offset);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `index` | `int` |
| `data` | `const void*` |
| `dataSize` | `int` |
| `offset` | `int` |

---

### UnloadMesh

```c
void UnloadMesh(Mesh mesh)
```

**Koda usage**

```koda
UnloadMesh(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |

---

### DrawMesh

```c
void DrawMesh(Mesh mesh, Material material, Matrix transform)
```

**Koda usage**

```koda
DrawMesh(mesh, material, transform);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `material` | `Material` |
| `transform` | `Matrix` |

---

### DrawMeshInstanced

```c
void DrawMeshInstanced(Mesh mesh, Material material, const Matrix* transforms, int instances)
```

**Koda usage**

```koda
DrawMeshInstanced(mesh, material, transforms, instances);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `material` | `Material` |
| `transforms` | `const Matrix*` |
| `instances` | `int` |

---

### GetMeshBoundingBox

```c
BoundingBox GetMeshBoundingBox(Mesh mesh)
```

**Koda usage**

```koda
let result = GetMeshBoundingBox(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |

---

### GenMeshTangents

```c
void GenMeshTangents(Mesh* mesh)
```

**Koda usage**

```koda
GenMeshTangents(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh*` |

---

### ExportMesh

```c
bool ExportMesh(Mesh mesh, const char* fileName)
```

**Koda usage**

```koda
let result = ExportMesh(mesh, fileName);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `fileName` | `const char*` |

---

### ExportMeshAsCode

```c
bool ExportMeshAsCode(Mesh mesh, const char* fileName)
```

**Koda usage**

```koda
let result = ExportMeshAsCode(mesh, fileName);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `fileName` | `const char*` |

---

### GenMeshPoly

```c
Mesh GenMeshPoly(int sides, float radius)
```

**Koda usage**

```koda
let result = GenMeshPoly(sides, radius);
```

| Parameter | Type |
|-----------|------|
| `sides` | `int` |
| `radius` | `float` |

---

### GenMeshPlane

```c
Mesh GenMeshPlane(float width, float length, int resX, int resZ)
```

**Koda usage**

```koda
let result = GenMeshPlane(width, length, resX, resZ);
```

| Parameter | Type |
|-----------|------|
| `width` | `float` |
| `length` | `float` |
| `resX` | `int` |
| `resZ` | `int` |

---

### GenMeshCube

```c
Mesh GenMeshCube(float width, float height, float length)
```

**Koda usage**

```koda
let result = GenMeshCube(width, height, length);
```

| Parameter | Type |
|-----------|------|
| `width` | `float` |
| `height` | `float` |
| `length` | `float` |

---

### GenMeshSphere

```c
Mesh GenMeshSphere(float radius, int rings, int slices)
```

**Koda usage**

```koda
let result = GenMeshSphere(radius, rings, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |

---

### GenMeshHemiSphere

```c
Mesh GenMeshHemiSphere(float radius, int rings, int slices)
```

**Koda usage**

```koda
let result = GenMeshHemiSphere(radius, rings, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |

---

### GenMeshCylinder

```c
Mesh GenMeshCylinder(float radius, float height, int slices)
```

**Koda usage**

```koda
let result = GenMeshCylinder(radius, height, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `height` | `float` |
| `slices` | `int` |

---

### GenMeshCone

```c
Mesh GenMeshCone(float radius, float height, int slices)
```

**Koda usage**

```koda
let result = GenMeshCone(radius, height, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `height` | `float` |
| `slices` | `int` |

---

### GenMeshTorus

```c
Mesh GenMeshTorus(float radius, float size, int radSeg, int sides)
```

**Koda usage**

```koda
let result = GenMeshTorus(radius, size, radSeg, sides);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `size` | `float` |
| `radSeg` | `int` |
| `sides` | `int` |

---

### GenMeshKnot

```c
Mesh GenMeshKnot(float radius, float size, int radSeg, int sides)
```

**Koda usage**

```koda
let result = GenMeshKnot(radius, size, radSeg, sides);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `size` | `float` |
| `radSeg` | `int` |
| `sides` | `int` |

---

### GenMeshHeightmap

```c
Mesh GenMeshHeightmap(Image heightmap, Vector3 size)
```

**Koda usage**

```koda
let result = GenMeshHeightmap(heightmap, size);
```

| Parameter | Type |
|-----------|------|
| `heightmap` | `Image` |
| `size` | `Vector3` |

---

### GenMeshCubicmap

```c
Mesh GenMeshCubicmap(Image cubicmap, Vector3 cubeSize)
```

**Koda usage**

```koda
let result = GenMeshCubicmap(cubicmap, cubeSize);
```

| Parameter | Type |
|-----------|------|
| `cubicmap` | `Image` |
| `cubeSize` | `Vector3` |

---

### LoadMaterialDefault

```c
Material LoadMaterialDefault()
```

**Koda usage**

```koda
let result = LoadMaterialDefault();
```

---

### IsMaterialValid

```c
bool IsMaterialValid(Material material)
```

**Koda usage**

```koda
let result = IsMaterialValid(material);
```

| Parameter | Type |
|-----------|------|
| `material` | `Material` |

---

### UnloadMaterial

```c
void UnloadMaterial(Material material)
```

**Koda usage**

```koda
UnloadMaterial(material);
```

| Parameter | Type |
|-----------|------|
| `material` | `Material` |

---

### SetMaterialTexture

```c
void SetMaterialTexture(Material* material, int mapType, Texture2D texture)
```

**Koda usage**

```koda
SetMaterialTexture(material, mapType, texture);
```

| Parameter | Type |
|-----------|------|
| `material` | `Material*` |
| `mapType` | `int` |
| `texture` | `Texture2D` |

---

### SetModelMeshMaterial

```c
void SetModelMeshMaterial(Model* model, int meshId, int materialId)
```

**Koda usage**

```koda
SetModelMeshMaterial(model, meshId, materialId);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model*` |
| `meshId` | `int` |
| `materialId` | `int` |

---

### UpdateModelAnimation

```c
void UpdateModelAnimation(Model model, ModelAnimation anim, float frame)
```

**Koda usage**

```koda
UpdateModelAnimation(model, anim, frame);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `anim` | `ModelAnimation` |
| `frame` | `float` |

---

### UpdateModelAnimationEx

```c
void UpdateModelAnimationEx(Model model, ModelAnimation animA, float frameA, ModelAnimation animB, float frameB, float blend)
```

**Koda usage**

```koda
UpdateModelAnimationEx(model, animA, frameA, animB, frameB, blend);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `animA` | `ModelAnimation` |
| `frameA` | `float` |
| `animB` | `ModelAnimation` |
| `frameB` | `float` |
| `blend` | `float` |

---

### UnloadModelAnimations

```c
void UnloadModelAnimations(ModelAnimation* animations, int animCount)
```

**Koda usage**

```koda
UnloadModelAnimations(animations, animCount);
```

| Parameter | Type |
|-----------|------|
| `animations` | `ModelAnimation*` |
| `animCount` | `int` |

---

### IsModelAnimationValid

```c
bool IsModelAnimationValid(Model model, ModelAnimation anim)
```

**Koda usage**

```koda
let result = IsModelAnimationValid(model, anim);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `anim` | `ModelAnimation` |

---

### CheckCollisionSpheres

```c
bool CheckCollisionSpheres(Vector3 center1, float radius1, Vector3 center2, float radius2)
```

**Koda usage**

```koda
let result = CheckCollisionSpheres(center1, radius1, center2, radius2);
```

| Parameter | Type |
|-----------|------|
| `center1` | `Vector3` |
| `radius1` | `float` |
| `center2` | `Vector3` |
| `radius2` | `float` |

---

### CheckCollisionBoxes

```c
bool CheckCollisionBoxes(BoundingBox box1, BoundingBox box2)
```

**Koda usage**

```koda
let result = CheckCollisionBoxes(box1, box2);
```

| Parameter | Type |
|-----------|------|
| `box1` | `BoundingBox` |
| `box2` | `BoundingBox` |

---

### CheckCollisionBoxSphere

```c
bool CheckCollisionBoxSphere(BoundingBox box, Vector3 center, float radius)
```

**Koda usage**

```koda
let result = CheckCollisionBoxSphere(box, center, radius);
```

| Parameter | Type |
|-----------|------|
| `box` | `BoundingBox` |
| `center` | `Vector3` |
| `radius` | `float` |

---

### GetRayCollisionSphere

```c
RayCollision GetRayCollisionSphere(Ray ray, Vector3 center, float radius)
```

**Koda usage**

```koda
let result = GetRayCollisionSphere(ray, center, radius);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `center` | `Vector3` |
| `radius` | `float` |

---

### GetRayCollisionBox

```c
RayCollision GetRayCollisionBox(Ray ray, BoundingBox box)
```

**Koda usage**

```koda
let result = GetRayCollisionBox(ray, box);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `box` | `BoundingBox` |

---

### GetRayCollisionMesh

```c
RayCollision GetRayCollisionMesh(Ray ray, Mesh mesh, Matrix transform)
```

**Koda usage**

```koda
let result = GetRayCollisionMesh(ray, mesh, transform);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `mesh` | `Mesh` |
| `transform` | `Matrix` |

---

### GetRayCollisionTriangle

```c
RayCollision GetRayCollisionTriangle(Ray ray, Vector3 p1, Vector3 p2, Vector3 p3)
```

**Koda usage**

```koda
let result = GetRayCollisionTriangle(ray, p1, p2, p3);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `p1` | `Vector3` |
| `p2` | `Vector3` |
| `p3` | `Vector3` |

---

### GetRayCollisionQuad

```c
RayCollision GetRayCollisionQuad(Ray ray, Vector3 p1, Vector3 p2, Vector3 p3, Vector3 p4)
```

**Koda usage**

```koda
let result = GetRayCollisionQuad(ray, p1, p2, p3, p4);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `p1` | `Vector3` |
| `p2` | `Vector3` |
| `p3` | `Vector3` |
| `p4` | `Vector3` |

---

### void

```c
typedef void()
```

**Koda usage**

```koda
let result = void();
```

---

### InitAudioDevice

```c
void InitAudioDevice()
```

**Koda usage**

```koda
InitAudioDevice();
```

---

### CloseAudioDevice

```c
void CloseAudioDevice()
```

**Koda usage**

```koda
CloseAudioDevice();
```

---

### IsAudioDeviceReady

```c
bool IsAudioDeviceReady()
```

**Koda usage**

```koda
let result = IsAudioDeviceReady();
```

---

### SetMasterVolume

```c
void SetMasterVolume(float volume)
```

**Koda usage**

```koda
SetMasterVolume(volume);
```

| Parameter | Type |
|-----------|------|
| `volume` | `float` |

---

### GetMasterVolume

```c
float GetMasterVolume()
```

**Koda usage**

```koda
let result = GetMasterVolume();
```

---

### LoadWave

```c
Wave LoadWave(const char* fileName)
```

**Koda usage**

```koda
let result = LoadWave(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadWaveFromMemory

```c
Wave LoadWaveFromMemory(const char* fileType, const unsigned char* fileData, int dataSize)
```

**Koda usage**

```koda
let result = LoadWaveFromMemory(fileType, fileData, dataSize);
```

| Parameter | Type |
|-----------|------|
| `fileType` | `const char*` |
| `fileData` | `const unsigned char*` |
| `dataSize` | `int` |

---

### IsWaveValid

```c
bool IsWaveValid(Wave wave)
```

**Koda usage**

```koda
let result = IsWaveValid(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### LoadSound

```c
Sound LoadSound(const char* fileName)
```

**Koda usage**

```koda
let result = LoadSound(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadSoundFromWave

```c
Sound LoadSoundFromWave(Wave wave)
```

**Koda usage**

```koda
let result = LoadSoundFromWave(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### LoadSoundAlias

```c
Sound LoadSoundAlias(Sound source)
```

**Koda usage**

```koda
let result = LoadSoundAlias(source);
```

| Parameter | Type |
|-----------|------|
| `source` | `Sound` |

---

### IsSoundValid

```c
bool IsSoundValid(Sound sound)
```

**Koda usage**

```koda
let result = IsSoundValid(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### UpdateSound

```c
void UpdateSound(Sound sound, const void* data, int frameCount)
```

**Koda usage**

```koda
UpdateSound(sound, data, frameCount);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `data` | `const void*` |
| `frameCount` | `int` |

---

### UnloadWave

```c
void UnloadWave(Wave wave)
```

**Koda usage**

```koda
UnloadWave(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### UnloadSound

```c
void UnloadSound(Sound sound)
```

**Koda usage**

```koda
UnloadSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### UnloadSoundAlias

```c
void UnloadSoundAlias(Sound alias)
```

**Koda usage**

```koda
UnloadSoundAlias(alias);
```

| Parameter | Type |
|-----------|------|
| `alias` | `Sound` |

---

### ExportWave

```c
bool ExportWave(Wave wave, const char* fileName)
```

**Koda usage**

```koda
let result = ExportWave(wave, fileName);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |
| `fileName` | `const char*` |

---

### ExportWaveAsCode

```c
bool ExportWaveAsCode(Wave wave, const char* fileName)
```

**Koda usage**

```koda
let result = ExportWaveAsCode(wave, fileName);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |
| `fileName` | `const char*` |

---

### PlaySound

```c
void PlaySound(Sound sound)
```

**Koda usage**

```koda
PlaySound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### StopSound

```c
void StopSound(Sound sound)
```

**Koda usage**

```koda
StopSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### PauseSound

```c
void PauseSound(Sound sound)
```

**Koda usage**

```koda
PauseSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### ResumeSound

```c
void ResumeSound(Sound sound)
```

**Koda usage**

```koda
ResumeSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### IsSoundPlaying

```c
bool IsSoundPlaying(Sound sound)
```

**Koda usage**

```koda
let result = IsSoundPlaying(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### SetSoundVolume

```c
void SetSoundVolume(Sound sound, float volume)
```

**Koda usage**

```koda
SetSoundVolume(sound, volume);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `volume` | `float` |

---

### SetSoundPitch

```c
void SetSoundPitch(Sound sound, float pitch)
```

**Koda usage**

```koda
SetSoundPitch(sound, pitch);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `pitch` | `float` |

---

### SetSoundPan

```c
void SetSoundPan(Sound sound, float pan)
```

**Koda usage**

```koda
SetSoundPan(sound, pan);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `pan` | `float` |

---

### WaveCopy

```c
Wave WaveCopy(Wave wave)
```

**Koda usage**

```koda
let result = WaveCopy(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### WaveCrop

```c
void WaveCrop(Wave* wave, int initFrame, int finalFrame)
```

**Koda usage**

```koda
WaveCrop(wave, initFrame, finalFrame);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave*` |
| `initFrame` | `int` |
| `finalFrame` | `int` |

---

### WaveFormat

```c
void WaveFormat(Wave* wave, int sampleRate, int sampleSize, int channels)
```

**Koda usage**

```koda
WaveFormat(wave, sampleRate, sampleSize, channels);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave*` |
| `sampleRate` | `int` |
| `sampleSize` | `int` |
| `channels` | `int` |

---

### UnloadWaveSamples

```c
void UnloadWaveSamples(float* samples)
```

**Koda usage**

```koda
UnloadWaveSamples(samples);
```

| Parameter | Type |
|-----------|------|
| `samples` | `float*` |

---

### LoadMusicStream

```c
Music LoadMusicStream(const char* fileName)
```

**Koda usage**

```koda
let result = LoadMusicStream(fileName);
```

| Parameter | Type |
|-----------|------|
| `fileName` | `const char*` |

---

### LoadMusicStreamFromMemory

```c
Music LoadMusicStreamFromMemory(const char* fileType, const unsigned char* data, int dataSize)
```

**Koda usage**

```koda
let result = LoadMusicStreamFromMemory(fileType, data, dataSize);
```

| Parameter | Type |
|-----------|------|
| `fileType` | `const char*` |
| `data` | `const unsigned char*` |
| `dataSize` | `int` |

---

### IsMusicValid

```c
bool IsMusicValid(Music music)
```

**Koda usage**

```koda
let result = IsMusicValid(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### UnloadMusicStream

```c
void UnloadMusicStream(Music music)
```

**Koda usage**

```koda
UnloadMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### PlayMusicStream

```c
void PlayMusicStream(Music music)
```

**Koda usage**

```koda
PlayMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### IsMusicStreamPlaying

```c
bool IsMusicStreamPlaying(Music music)
```

**Koda usage**

```koda
let result = IsMusicStreamPlaying(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### UpdateMusicStream

```c
void UpdateMusicStream(Music music)
```

**Koda usage**

```koda
UpdateMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### StopMusicStream

```c
void StopMusicStream(Music music)
```

**Koda usage**

```koda
StopMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### PauseMusicStream

```c
void PauseMusicStream(Music music)
```

**Koda usage**

```koda
PauseMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### ResumeMusicStream

```c
void ResumeMusicStream(Music music)
```

**Koda usage**

```koda
ResumeMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### SeekMusicStream

```c
void SeekMusicStream(Music music, float position)
```

**Koda usage**

```koda
SeekMusicStream(music, position);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `position` | `float` |

---

### SetMusicVolume

```c
void SetMusicVolume(Music music, float volume)
```

**Koda usage**

```koda
SetMusicVolume(music, volume);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `volume` | `float` |

---

### SetMusicPitch

```c
void SetMusicPitch(Music music, float pitch)
```

**Koda usage**

```koda
SetMusicPitch(music, pitch);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `pitch` | `float` |

---

### SetMusicPan

```c
void SetMusicPan(Music music, float pan)
```

**Koda usage**

```koda
SetMusicPan(music, pan);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `pan` | `float` |

---

### GetMusicTimeLength

```c
float GetMusicTimeLength(Music music)
```

**Koda usage**

```koda
let result = GetMusicTimeLength(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### GetMusicTimePlayed

```c
float GetMusicTimePlayed(Music music)
```

**Koda usage**

```koda
let result = GetMusicTimePlayed(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### LoadAudioStream

```c
AudioStream LoadAudioStream(unsigned int sampleRate, unsigned int sampleSize, unsigned int channels)
```

**Koda usage**

```koda
let result = LoadAudioStream(sampleRate, sampleSize, channels);
```

| Parameter | Type |
|-----------|------|
| `sampleRate` | `unsigned int` |
| `sampleSize` | `unsigned int` |
| `channels` | `unsigned int` |

---

### IsAudioStreamValid

```c
bool IsAudioStreamValid(AudioStream stream)
```

**Koda usage**

```koda
let result = IsAudioStreamValid(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### UnloadAudioStream

```c
void UnloadAudioStream(AudioStream stream)
```

**Koda usage**

```koda
UnloadAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### UpdateAudioStream

```c
void UpdateAudioStream(AudioStream stream, const void* data, int frameCount)
```

**Koda usage**

```koda
UpdateAudioStream(stream, data, frameCount);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `data` | `const void*` |
| `frameCount` | `int` |

---

### IsAudioStreamProcessed

```c
bool IsAudioStreamProcessed(AudioStream stream)
```

**Koda usage**

```koda
let result = IsAudioStreamProcessed(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### PlayAudioStream

```c
void PlayAudioStream(AudioStream stream)
```

**Koda usage**

```koda
PlayAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### PauseAudioStream

```c
void PauseAudioStream(AudioStream stream)
```

**Koda usage**

```koda
PauseAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### ResumeAudioStream

```c
void ResumeAudioStream(AudioStream stream)
```

**Koda usage**

```koda
ResumeAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### IsAudioStreamPlaying

```c
bool IsAudioStreamPlaying(AudioStream stream)
```

**Koda usage**

```koda
let result = IsAudioStreamPlaying(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### StopAudioStream

```c
void StopAudioStream(AudioStream stream)
```

**Koda usage**

```koda
StopAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### SetAudioStreamVolume

```c
void SetAudioStreamVolume(AudioStream stream, float volume)
```

**Koda usage**

```koda
SetAudioStreamVolume(stream, volume);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `volume` | `float` |

---

### SetAudioStreamPitch

```c
void SetAudioStreamPitch(AudioStream stream, float pitch)
```

**Koda usage**

```koda
SetAudioStreamPitch(stream, pitch);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `pitch` | `float` |

---

### SetAudioStreamPan

```c
void SetAudioStreamPan(AudioStream stream, float pan)
```

**Koda usage**

```koda
SetAudioStreamPan(stream, pan);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `pan` | `float` |

---

### SetAudioStreamBufferSizeDefault

```c
void SetAudioStreamBufferSizeDefault(int size)
```

**Koda usage**

```koda
SetAudioStreamBufferSizeDefault(size);
```

| Parameter | Type |
|-----------|------|
| `size` | `int` |

---

### SetAudioStreamCallback

```c
void SetAudioStreamCallback(AudioStream stream, AudioCallback callback)
```

**Koda usage**

```koda
SetAudioStreamCallback(stream, callback);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `callback` | `AudioCallback` |

---

### AttachAudioStreamProcessor

```c
void AttachAudioStreamProcessor(AudioStream stream, AudioCallback processor)
```

**Koda usage**

```koda
AttachAudioStreamProcessor(stream, processor);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `processor` | `AudioCallback` |

---

### DetachAudioStreamProcessor

```c
void DetachAudioStreamProcessor(AudioStream stream, AudioCallback processor)
```

**Koda usage**

```koda
DetachAudioStreamProcessor(stream, processor);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `processor` | `AudioCallback` |

---

### AttachAudioMixedProcessor

```c
void AttachAudioMixedProcessor(AudioCallback processor)
```

**Koda usage**

```koda
AttachAudioMixedProcessor(processor);
```

| Parameter | Type |
|-----------|------|
| `processor` | `AudioCallback` |

---

### DetachAudioMixedProcessor

```c
void DetachAudioMixedProcessor(AudioCallback processor)
```

**Koda usage**

```koda
DetachAudioMixedProcessor(processor);
```

| Parameter | Type |
|-----------|------|
| `processor` | `AudioCallback` |

---

## Structs

### Vector2

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `float` |

---

### Vector3

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `float` |
| `z` | `float` |

---

### Vector4

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `float` |
| `z` | `float` |
| `w` | `float` |

---

### Matrix

| Field | Type |
|-------|------|
| `m12` | `float m0, m4, m8,` |
| `m13` | `float m1, m5, m9,` |
| `m14` | `float m2, m6, m10,` |
| `m15` | `float m3, m7, m11,` |

---

### Color

| Field | Type |
|-------|------|
| `r` | `unsigned char` |
| `g` | `unsigned char` |
| `b` | `unsigned char` |
| `a` | `unsigned char` |

---

### Rectangle

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `float` |
| `width` | `float` |
| `height` | `float` |

---

### Image

| Field | Type |
|-------|------|
| `data` | `void*` |
| `width` | `int` |
| `height` | `int` |
| `mipmaps` | `int` |
| `format` | `int` |

---

### Texture

| Field | Type |
|-------|------|
| `id` | `unsigned int` |
| `width` | `int` |
| `height` | `int` |
| `mipmaps` | `int` |
| `format` | `int` |

---

### RenderTexture

| Field | Type |
|-------|------|
| `id` | `unsigned int` |
| `texture` | `Texture` |
| `depth` | `Texture` |

---

### NPatchInfo

| Field | Type |
|-------|------|
| `source` | `Rectangle` |
| `left` | `int` |
| `top` | `int` |
| `right` | `int` |
| `bottom` | `int` |
| `layout` | `int` |

---

### GlyphInfo

| Field | Type |
|-------|------|
| `value` | `int` |
| `offsetX` | `int` |
| `offsetY` | `int` |
| `advanceX` | `int` |
| `image` | `Image` |

---

### Font

| Field | Type |
|-------|------|
| `baseSize` | `int` |
| `glyphCount` | `int` |
| `glyphPadding` | `int` |
| `texture` | `Texture2D` |
| `recs` | `Rectangle*` |
| `glyphs` | `GlyphInfo*` |

---

### Camera3D

| Field | Type |
|-------|------|
| `position` | `Vector3` |
| `target` | `Vector3` |
| `up` | `Vector3` |
| `fovy` | `float` |
| `projection` | `int` |

---

### Camera2D

| Field | Type |
|-------|------|
| `offset` | `Vector2` |
| `target` | `Vector2` |
| `rotation` | `float` |
| `zoom` | `float` |

---

### Mesh

| Field | Type |
|-------|------|
| `vertexCount` | `int` |
| `triangleCount` | `int` |
| `vertices` | `float*` |
| `texcoords` | `float*` |
| `texcoords2` | `float*` |
| `normals` | `float*` |
| `tangents` | `float*` |
| `colors` | `unsigned char*` |
| `indices` | `unsigned short*` |
| `boneCount` | `int` |
| `boneIndices` | `unsigned char*` |
| `boneWeights` | `float*` |
| `animVertices` | `float*` |
| `animNormals` | `float*` |
| `vaoId` | `unsigned int` |
| `vboId` | `unsigned int*` |

---

### Shader

| Field | Type |
|-------|------|
| `id` | `unsigned int` |
| `locs` | `int*` |

---

### MaterialMap

| Field | Type |
|-------|------|
| `texture` | `Texture2D` |
| `color` | `Color` |
| `value` | `float` |

---

### Material

| Field | Type |
|-------|------|
| `shader` | `Shader` |
| `maps` | `MaterialMap*` |
| `params[4]` | `float` |

---

### Transform

| Field | Type |
|-------|------|
| `translation` | `Vector3` |
| `rotation` | `Quaternion` |
| `scale` | `Vector3` |

---

### BoneInfo

| Field | Type |
|-------|------|
| `name[32]` | `char` |
| `parent` | `int` |

---

### ModelSkeleton

| Field | Type |
|-------|------|
| `boneCount` | `int` |
| `bones` | `BoneInfo*` |
| `bindPose` | `ModelAnimPose` |

---

### Model

| Field | Type |
|-------|------|
| `transform` | `Matrix` |
| `meshCount` | `int` |
| `materialCount` | `int` |
| `meshes` | `Mesh*` |
| `materials` | `Material*` |
| `meshMaterial` | `int*` |
| `skeleton` | `ModelSkeleton` |
| `currentPose` | `ModelAnimPose` |
| `boneMatrices` | `Matrix*` |

---

### ModelAnimation

| Field | Type |
|-------|------|
| `name[32]` | `char` |
| `boneCount` | `int` |
| `keyframeCount` | `int` |
| `keyframePoses` | `ModelAnimPose*` |

---

### Ray

| Field | Type |
|-------|------|
| `position` | `Vector3` |
| `direction` | `Vector3` |

---

### RayCollision

| Field | Type |
|-------|------|
| `hit` | `bool` |
| `distance` | `float` |
| `point` | `Vector3` |
| `normal` | `Vector3` |

---

### BoundingBox

| Field | Type |
|-------|------|
| `min` | `Vector3` |
| `max` | `Vector3` |

---

### Wave

| Field | Type |
|-------|------|
| `frameCount` | `unsigned int` |
| `sampleRate` | `unsigned int` |
| `sampleSize` | `unsigned int` |
| `channels` | `unsigned int` |
| `data` | `void*` |

---

### AudioStream

| Field | Type |
|-------|------|
| `buffer` | `rAudioBuffer*` |
| `processor` | `rAudioProcessor*` |
| `sampleRate` | `unsigned int` |
| `sampleSize` | `unsigned int` |
| `channels` | `unsigned int` |

---

### Sound

| Field | Type |
|-------|------|
| `stream` | `AudioStream` |
| `frameCount` | `unsigned int` |

---

### Music

| Field | Type |
|-------|------|
| `stream` | `AudioStream` |
| `frameCount` | `unsigned int` |
| `looping` | `bool` |
| `ctxType` | `int` |
| `ctxData` | `void*` |

---

### VrDeviceInfo

| Field | Type |
|-------|------|
| `hResolution` | `int` |
| `vResolution` | `int` |
| `hScreenSize` | `float` |
| `vScreenSize` | `float` |
| `eyeToScreenDistance` | `float` |
| `lensSeparationDistance` | `float` |
| `interpupillaryDistance` | `float` |
| `lensDistortionValues[4]` | `float` |
| `chromaAbCorrection[4]` | `float` |

---

### VrStereoConfig

| Field | Type |
|-------|------|
| `projection[2]` | `Matrix` |
| `viewOffset[2]` | `Matrix` |
| `leftLensCenter[2]` | `float` |
| `rightLensCenter[2]` | `float` |
| `leftScreenCenter[2]` | `float` |
| `rightScreenCenter[2]` | `float` |
| `scale[2]` | `float` |
| `scaleIn[2]` | `float` |

---

### FilePathList

| Field | Type |
|-------|------|
| `count` | `unsigned int` |
| `paths` | `char**` |

---

### AutomationEvent

| Field | Type |
|-------|------|
| `frame` | `unsigned int` |
| `type` | `unsigned int` |
| `params[4]` | `int` |

---

### AutomationEventList

| Field | Type |
|-------|------|
| `capacity` | `unsigned int` |
| `count` | `unsigned int` |
| `events` | `AutomationEvent*` |

---

## Enums

### bool

| Name | Value |
|------|-------|
| `false` | `0` |
| `true` | `1` |

---

### ConfigFlags

| Name | Value |
|------|-------|
| `FLAG_VSYNC_HINT` | `0` |
| `FLAG_FULLSCREEN_MODE` | `0` |
| `FLAG_WINDOW_RESIZABLE` | `0` |
| `FLAG_WINDOW_UNDECORATED` | `0` |
| `FLAG_WINDOW_HIDDEN` | `0` |
| `FLAG_WINDOW_MINIMIZED` | `0` |
| `FLAG_WINDOW_MAXIMIZED` | `0` |
| `FLAG_WINDOW_UNFOCUSED` | `0` |
| `FLAG_WINDOW_TOPMOST` | `0` |
| `FLAG_WINDOW_ALWAYS_RUN` | `0` |
| `FLAG_WINDOW_TRANSPARENT` | `0` |
| `FLAG_WINDOW_HIGHDPI` | `0` |
| `FLAG_WINDOW_MOUSE_PASSTHROUGH` | `0` |
| `FLAG_BORDERLESS_WINDOWED_MODE` | `0` |
| `FLAG_MSAA_4X_HINT` | `0` |
| `FLAG_INTERLACED_HINT` | `0` |

---

### TraceLogLevel

| Name | Value |
|------|-------|
| `LOG_ALL` | `0` |
| `LOG_TRACE` | `1` |
| `LOG_DEBUG` | `2` |
| `LOG_INFO` | `3` |
| `LOG_WARNING` | `4` |
| `LOG_ERROR` | `5` |
| `LOG_FATAL` | `6` |
| `LOG_NONE` | `7` |

---

### KeyboardKey

| Name | Value |
|------|-------|
| `KEY_NULL` | `0` |
| `KEY_APOSTROPHE` | `39` |
| `KEY_COMMA` | `44` |
| `KEY_MINUS` | `45` |
| `KEY_PERIOD` | `46` |
| `KEY_SLASH` | `47` |
| `KEY_ZERO` | `48` |
| `KEY_ONE` | `49` |
| `KEY_TWO` | `50` |
| `KEY_THREE` | `51` |
| `KEY_FOUR` | `52` |
| `KEY_FIVE` | `53` |
| `KEY_SIX` | `54` |
| `KEY_SEVEN` | `55` |
| `KEY_EIGHT` | `56` |
| `KEY_NINE` | `57` |
| `KEY_SEMICOLON` | `59` |
| `KEY_EQUAL` | `61` |
| `KEY_A` | `65` |
| `KEY_B` | `66` |
| `KEY_C` | `67` |
| `KEY_D` | `68` |
| `KEY_E` | `69` |
| `KEY_F` | `70` |
| `KEY_G` | `71` |
| `KEY_H` | `72` |
| `KEY_I` | `73` |
| `KEY_J` | `74` |
| `KEY_K` | `75` |
| `KEY_L` | `76` |
| `KEY_M` | `77` |
| `KEY_N` | `78` |
| `KEY_O` | `79` |
| `KEY_P` | `80` |
| `KEY_Q` | `81` |
| `KEY_R` | `82` |
| `KEY_S` | `83` |
| `KEY_T` | `84` |
| `KEY_U` | `85` |
| `KEY_V` | `86` |
| `KEY_W` | `87` |
| `KEY_X` | `88` |
| `KEY_Y` | `89` |
| `KEY_Z` | `90` |
| `KEY_LEFT_BRACKET` | `91` |
| `KEY_BACKSLASH` | `92` |
| `KEY_RIGHT_BRACKET` | `93` |
| `KEY_GRAVE` | `96` |
| `KEY_SPACE` | `32` |
| `KEY_ESCAPE` | `256` |
| `KEY_ENTER` | `257` |
| `KEY_TAB` | `258` |
| `KEY_BACKSPACE` | `259` |
| `KEY_INSERT` | `260` |
| `KEY_DELETE` | `261` |
| `KEY_RIGHT` | `262` |
| `KEY_LEFT` | `263` |
| `KEY_DOWN` | `264` |
| `KEY_UP` | `265` |
| `KEY_PAGE_UP` | `266` |
| `KEY_PAGE_DOWN` | `267` |
| `KEY_HOME` | `268` |
| `KEY_END` | `269` |
| `KEY_CAPS_LOCK` | `280` |
| `KEY_SCROLL_LOCK` | `281` |
| `KEY_NUM_LOCK` | `282` |
| `KEY_PRINT_SCREEN` | `283` |
| `KEY_PAUSE` | `284` |
| `KEY_F1` | `290` |
| `KEY_F2` | `291` |
| `KEY_F3` | `292` |
| `KEY_F4` | `293` |
| `KEY_F5` | `294` |
| `KEY_F6` | `295` |
| `KEY_F7` | `296` |
| `KEY_F8` | `297` |
| `KEY_F9` | `298` |
| `KEY_F10` | `299` |
| `KEY_F11` | `300` |
| `KEY_F12` | `301` |
| `KEY_LEFT_SHIFT` | `340` |
| `KEY_LEFT_CONTROL` | `341` |
| `KEY_LEFT_ALT` | `342` |
| `KEY_LEFT_SUPER` | `343` |
| `KEY_RIGHT_SHIFT` | `344` |
| `KEY_RIGHT_CONTROL` | `345` |
| `KEY_RIGHT_ALT` | `346` |
| `KEY_RIGHT_SUPER` | `347` |
| `KEY_KB_MENU` | `348` |
| `KEY_KP_0` | `320` |
| `KEY_KP_1` | `321` |
| `KEY_KP_2` | `322` |
| `KEY_KP_3` | `323` |
| `KEY_KP_4` | `324` |
| `KEY_KP_5` | `325` |
| `KEY_KP_6` | `326` |
| `KEY_KP_7` | `327` |
| `KEY_KP_8` | `328` |
| `KEY_KP_9` | `329` |
| `KEY_KP_DECIMAL` | `330` |
| `KEY_KP_DIVIDE` | `331` |
| `KEY_KP_MULTIPLY` | `332` |
| `KEY_KP_SUBTRACT` | `333` |
| `KEY_KP_ADD` | `334` |
| `KEY_KP_ENTER` | `335` |
| `KEY_KP_EQUAL` | `336` |
| `KEY_BACK` | `4` |
| `KEY_MENU` | `5` |
| `KEY_VOLUME_UP` | `24` |
| `KEY_VOLUME_DOWN` | `25` |

---

### MouseButton

| Name | Value |
|------|-------|
| `MOUSE_BUTTON_LEFT` | `0` |
| `MOUSE_BUTTON_RIGHT` | `1` |
| `MOUSE_BUTTON_MIDDLE` | `2` |
| `MOUSE_BUTTON_SIDE` | `3` |
| `MOUSE_BUTTON_EXTRA` | `4` |
| `MOUSE_BUTTON_FORWARD` | `5` |
| `MOUSE_BUTTON_BACK` | `6` |

---

### MouseCursor

| Name | Value |
|------|-------|
| `MOUSE_CURSOR_DEFAULT` | `0` |
| `MOUSE_CURSOR_ARROW` | `1` |
| `MOUSE_CURSOR_IBEAM` | `2` |
| `MOUSE_CURSOR_CROSSHAIR` | `3` |
| `MOUSE_CURSOR_POINTING_HAND` | `4` |
| `MOUSE_CURSOR_RESIZE_EW` | `5` |
| `MOUSE_CURSOR_RESIZE_NS` | `6` |
| `MOUSE_CURSOR_RESIZE_NWSE` | `7` |
| `MOUSE_CURSOR_RESIZE_NESW` | `8` |
| `MOUSE_CURSOR_RESIZE_ALL` | `9` |
| `MOUSE_CURSOR_NOT_ALLOWED` | `10` |

---

### GamepadButton

| Name | Value |
|------|-------|
| `GAMEPAD_BUTTON_UNKNOWN` | `0` |
| `GAMEPAD_BUTTON_LEFT_FACE_UP` | `1` |
| `GAMEPAD_BUTTON_LEFT_FACE_RIGHT` | `2` |
| `GAMEPAD_BUTTON_LEFT_FACE_DOWN` | `3` |
| `GAMEPAD_BUTTON_LEFT_FACE_LEFT` | `4` |
| `GAMEPAD_BUTTON_RIGHT_FACE_UP` | `5` |
| `GAMEPAD_BUTTON_RIGHT_FACE_RIGHT` | `6` |
| `GAMEPAD_BUTTON_RIGHT_FACE_DOWN` | `7` |
| `GAMEPAD_BUTTON_RIGHT_FACE_LEFT` | `8` |
| `GAMEPAD_BUTTON_LEFT_TRIGGER_1` | `9` |
| `GAMEPAD_BUTTON_LEFT_TRIGGER_2` | `10` |
| `GAMEPAD_BUTTON_RIGHT_TRIGGER_1` | `11` |
| `GAMEPAD_BUTTON_RIGHT_TRIGGER_2` | `12` |
| `GAMEPAD_BUTTON_MIDDLE_LEFT` | `13` |
| `GAMEPAD_BUTTON_MIDDLE` | `14` |
| `GAMEPAD_BUTTON_MIDDLE_RIGHT` | `15` |
| `GAMEPAD_BUTTON_LEFT_THUMB` | `16` |
| `GAMEPAD_BUTTON_RIGHT_THUMB` | `17` |

---

### GamepadAxis

| Name | Value |
|------|-------|
| `GAMEPAD_AXIS_LEFT_X` | `0` |
| `GAMEPAD_AXIS_LEFT_Y` | `1` |
| `GAMEPAD_AXIS_RIGHT_X` | `2` |
| `GAMEPAD_AXIS_RIGHT_Y` | `3` |
| `GAMEPAD_AXIS_LEFT_TRIGGER` | `4` |
| `GAMEPAD_AXIS_RIGHT_TRIGGER` | `5` |

---

### MaterialMapIndex

| Name | Value |
|------|-------|
| `MATERIAL_MAP_ALBEDO` | `0` |
| `MATERIAL_MAP_METALNESS` | `1` |
| `MATERIAL_MAP_NORMAL` | `2` |
| `MATERIAL_MAP_ROUGHNESS` | `3` |
| `MATERIAL_MAP_OCCLUSION` | `4` |
| `MATERIAL_MAP_EMISSION` | `5` |
| `MATERIAL_MAP_HEIGHT` | `6` |
| `MATERIAL_MAP_CUBEMAP` | `7` |
| `MATERIAL_MAP_IRRADIANCE` | `8` |
| `MATERIAL_MAP_PREFILTER` | `9` |
| `MATERIAL_MAP_BRDF` | `10` |

---

### ShaderLocationIndex

| Name | Value |
|------|-------|
| `SHADER_LOC_VERTEX_POSITION` | `0` |
| `SHADER_LOC_VERTEX_TEXCOORD01` | `1` |
| `SHADER_LOC_VERTEX_TEXCOORD02` | `2` |
| `SHADER_LOC_VERTEX_NORMAL` | `3` |
| `SHADER_LOC_VERTEX_TANGENT` | `4` |
| `SHADER_LOC_VERTEX_COLOR` | `5` |
| `SHADER_LOC_MATRIX_MVP` | `6` |
| `SHADER_LOC_MATRIX_VIEW` | `7` |
| `SHADER_LOC_MATRIX_PROJECTION` | `8` |
| `SHADER_LOC_MATRIX_MODEL` | `9` |
| `SHADER_LOC_MATRIX_NORMAL` | `10` |
| `SHADER_LOC_VECTOR_VIEW` | `11` |
| `SHADER_LOC_COLOR_DIFFUSE` | `12` |
| `SHADER_LOC_COLOR_SPECULAR` | `13` |
| `SHADER_LOC_COLOR_AMBIENT` | `14` |
| `SHADER_LOC_MAP_ALBEDO` | `15` |
| `SHADER_LOC_MAP_METALNESS` | `16` |
| `SHADER_LOC_MAP_NORMAL` | `17` |
| `SHADER_LOC_MAP_ROUGHNESS` | `18` |
| `SHADER_LOC_MAP_OCCLUSION` | `19` |
| `SHADER_LOC_MAP_EMISSION` | `20` |
| `SHADER_LOC_MAP_HEIGHT` | `21` |
| `SHADER_LOC_MAP_CUBEMAP` | `22` |
| `SHADER_LOC_MAP_IRRADIANCE` | `23` |
| `SHADER_LOC_MAP_PREFILTER` | `24` |
| `SHADER_LOC_MAP_BRDF` | `25` |
| `SHADER_LOC_VERTEX_BONEIDS` | `26` |
| `SHADER_LOC_VERTEX_BONEWEIGHTS` | `27` |
| `SHADER_LOC_MATRIX_BONETRANSFORMS` | `28` |
| `SHADER_LOC_VERTEX_INSTANCETRANSFORM` | `29` |

---

### ShaderUniformDataType

| Name | Value |
|------|-------|
| `SHADER_UNIFORM_FLOAT` | `0` |
| `SHADER_UNIFORM_VEC2` | `1` |
| `SHADER_UNIFORM_VEC3` | `2` |
| `SHADER_UNIFORM_VEC4` | `3` |
| `SHADER_UNIFORM_INT` | `4` |
| `SHADER_UNIFORM_IVEC2` | `5` |
| `SHADER_UNIFORM_IVEC3` | `6` |
| `SHADER_UNIFORM_IVEC4` | `7` |
| `SHADER_UNIFORM_UINT` | `8` |
| `SHADER_UNIFORM_UIVEC2` | `9` |
| `SHADER_UNIFORM_UIVEC3` | `10` |
| `SHADER_UNIFORM_UIVEC4` | `11` |
| `SHADER_UNIFORM_SAMPLER2D` | `12` |

---

### ShaderAttributeDataType

| Name | Value |
|------|-------|
| `SHADER_ATTRIB_FLOAT` | `0` |
| `SHADER_ATTRIB_VEC2` | `1` |
| `SHADER_ATTRIB_VEC3` | `2` |
| `SHADER_ATTRIB_VEC4` | `3` |

---

### PixelFormat

| Name | Value |
|------|-------|
| `PIXELFORMAT_UNCOMPRESSED_GRAYSCALE` | `1` |
| `PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA` | `2` |
| `PIXELFORMAT_UNCOMPRESSED_R5G6B5` | `3` |
| `PIXELFORMAT_UNCOMPRESSED_R8G8B8` | `4` |
| `PIXELFORMAT_UNCOMPRESSED_R5G5B5A1` | `5` |
| `PIXELFORMAT_UNCOMPRESSED_R4G4B4A4` | `6` |
| `PIXELFORMAT_UNCOMPRESSED_R8G8B8A8` | `7` |
| `PIXELFORMAT_UNCOMPRESSED_R32` | `8` |
| `PIXELFORMAT_UNCOMPRESSED_R32G32B32` | `9` |
| `PIXELFORMAT_UNCOMPRESSED_R32G32B32A32` | `10` |
| `PIXELFORMAT_UNCOMPRESSED_R16` | `11` |
| `PIXELFORMAT_UNCOMPRESSED_R16G16B16` | `12` |
| `PIXELFORMAT_UNCOMPRESSED_R16G16B16A16` | `13` |
| `PIXELFORMAT_COMPRESSED_DXT1_RGB` | `14` |
| `PIXELFORMAT_COMPRESSED_DXT1_RGBA` | `15` |
| `PIXELFORMAT_COMPRESSED_DXT3_RGBA` | `16` |
| `PIXELFORMAT_COMPRESSED_DXT5_RGBA` | `17` |
| `PIXELFORMAT_COMPRESSED_ETC1_RGB` | `18` |
| `PIXELFORMAT_COMPRESSED_ETC2_RGB` | `19` |
| `PIXELFORMAT_COMPRESSED_ETC2_EAC_RGBA` | `20` |
| `PIXELFORMAT_COMPRESSED_PVRT_RGB` | `21` |
| `PIXELFORMAT_COMPRESSED_PVRT_RGBA` | `22` |
| `PIXELFORMAT_COMPRESSED_ASTC_4x4_RGBA` | `23` |
| `PIXELFORMAT_COMPRESSED_ASTC_8x8_RGBA` | `24` |

---

### TextureFilter

| Name | Value |
|------|-------|
| `TEXTURE_FILTER_POINT` | `0` |
| `TEXTURE_FILTER_BILINEAR` | `1` |
| `TEXTURE_FILTER_TRILINEAR` | `2` |
| `TEXTURE_FILTER_ANISOTROPIC_4X` | `3` |
| `TEXTURE_FILTER_ANISOTROPIC_8X` | `4` |
| `TEXTURE_FILTER_ANISOTROPIC_16X` | `5` |

---

### TextureWrap

| Name | Value |
|------|-------|
| `TEXTURE_WRAP_REPEAT` | `0` |
| `TEXTURE_WRAP_CLAMP` | `1` |
| `TEXTURE_WRAP_MIRROR_REPEAT` | `2` |
| `TEXTURE_WRAP_MIRROR_CLAMP` | `3` |

---

### CubemapLayout

| Name | Value |
|------|-------|
| `CUBEMAP_LAYOUT_AUTO_DETECT` | `0` |
| `CUBEMAP_LAYOUT_LINE_VERTICAL` | `1` |
| `CUBEMAP_LAYOUT_LINE_HORIZONTAL` | `2` |
| `CUBEMAP_LAYOUT_CROSS_THREE_BY_FOUR` | `3` |
| `CUBEMAP_LAYOUT_CROSS_FOUR_BY_THREE` | `4` |

---

### FontType

| Name | Value |
|------|-------|
| `FONT_DEFAULT` | `0` |
| `FONT_BITMAP` | `1` |
| `FONT_SDF` | `2` |

---

### BlendMode

| Name | Value |
|------|-------|
| `BLEND_ALPHA` | `0` |
| `BLEND_ADDITIVE` | `1` |
| `BLEND_MULTIPLIED` | `2` |
| `BLEND_ADD_COLORS` | `3` |
| `BLEND_SUBTRACT_COLORS` | `4` |
| `BLEND_ALPHA_PREMULTIPLY` | `5` |
| `BLEND_CUSTOM` | `6` |
| `BLEND_CUSTOM_SEPARATE` | `7` |

---

### Gesture

| Name | Value |
|------|-------|
| `GESTURE_NONE` | `0` |
| `GESTURE_TAP` | `1` |
| `GESTURE_DOUBLETAP` | `2` |
| `GESTURE_HOLD` | `4` |
| `GESTURE_DRAG` | `8` |
| `GESTURE_SWIPE_RIGHT` | `16` |
| `GESTURE_SWIPE_LEFT` | `32` |
| `GESTURE_SWIPE_UP` | `64` |
| `GESTURE_SWIPE_DOWN` | `128` |
| `GESTURE_PINCH_IN` | `256` |
| `GESTURE_PINCH_OUT` | `512` |

---

### CameraMode

| Name | Value |
|------|-------|
| `CAMERA_CUSTOM` | `0` |
| `CAMERA_FREE` | `1` |
| `CAMERA_ORBITAL` | `2` |
| `CAMERA_FIRST_PERSON` | `3` |
| `CAMERA_THIRD_PERSON` | `4` |

---

### CameraProjection

| Name | Value |
|------|-------|
| `CAMERA_PERSPECTIVE` | `0` |
| `CAMERA_ORTHOGRAPHIC` | `1` |

---

### NPatchLayout

| Name | Value |
|------|-------|
| `NPATCH_NINE_PATCH` | `0` |
| `NPATCH_THREE_PATCH_VERTICAL` | `1` |
| `NPATCH_THREE_PATCH_HORIZONTAL` | `2` |

---

## Macros

### RAYLIB_H

```c
#define RAYLIB_H 
```

---

### RAYLIB_VERSION_MAJOR

```c
#define RAYLIB_VERSION_MAJOR 
```

---

### RAYLIB_VERSION_MINOR

```c
#define RAYLIB_VERSION_MINOR 
```

---

### RAYLIB_VERSION_PATCH

```c
#define RAYLIB_VERSION_PATCH 
```

---

### RAYLIB_VERSION

```c
#define RAYLIB_VERSION 
```

---

### RL_COLOR_TYPE

```c
#define RL_COLOR_TYPE 
```

---

### RL_VECTOR2_TYPE

```c
#define RL_VECTOR2_TYPE 
```

---

### RL_VECTOR4_TYPE

```c
#define RL_VECTOR4_TYPE 
```

---

### RL_MATRIX_TYPE

```c
#define RL_MATRIX_TYPE 
```

---

### GRAY

```c
#define GRAY 
```

---

### DARKGRAY

```c
#define DARKGRAY 
```

---

### YELLOW

```c
#define YELLOW 
```

---

### GOLD

```c
#define GOLD 
```

---

### ORANGE

```c
#define ORANGE 
```

---

### PINK

```c
#define PINK 
```

---

### RED

```c
#define RED 
```

---

### MAROON

```c
#define MAROON 
```

---

### GREEN

```c
#define GREEN 
```

---

### LIME

```c
#define LIME 
```

---

### DARKGREEN

```c
#define DARKGREEN 
```

---

### SKYBLUE

```c
#define SKYBLUE 
```

---

### BLUE

```c
#define BLUE 
```

---

### DARKBLUE

```c
#define DARKBLUE 
```

---

### PURPLE

```c
#define PURPLE 
```

---

### VIOLET

```c
#define VIOLET 
```

---

### DARKPURPLE

```c
#define DARKPURPLE 
```

---

### BEIGE

```c
#define BEIGE 
```

---

### BROWN

```c
#define BROWN 
```

---

### DARKBROWN

```c
#define DARKBROWN 
```

---

### WHITE

```c
#define WHITE 
```

---

### BLACK

```c
#define BLACK 
```

---

### BLANK

```c
#define BLANK 
```

---

### MAGENTA

```c
#define MAGENTA 
```

---

### RAYWHITE

```c
#define RAYWHITE 
```

---

### MOUSE_LEFT_BUTTON

```c
#define MOUSE_LEFT_BUTTON 
```

---

### MOUSE_RIGHT_BUTTON

```c
#define MOUSE_RIGHT_BUTTON 
```

---

### MOUSE_MIDDLE_BUTTON

```c
#define MOUSE_MIDDLE_BUTTON 
```

---

### MATERIAL_MAP_DIFFUSE

```c
#define MATERIAL_MAP_DIFFUSE 
```

---

### MATERIAL_MAP_SPECULAR

```c
#define MATERIAL_MAP_SPECULAR 
```

---

### SHADER_LOC_MAP_DIFFUSE

```c
#define SHADER_LOC_MAP_DIFFUSE 
```

---

### SHADER_LOC_MAP_SPECULAR

```c
#define SHADER_LOC_MAP_SPECULAR 
```

---

### GetMouseRay

```c
#define GetMouseRay 
```

---

