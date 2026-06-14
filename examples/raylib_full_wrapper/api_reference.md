# raylib API Reference

## Table of Contents

- [Functions](#functions)
- [Structs](#structs)
- [Enums](#enums)
- [Macros](#macros)

## Functions

### void

**Signature:**
```c
typedef void()
```

**Koda Binding:**
```koda
native void();
```

### bool

**Signature:**
```c
typedef bool()
```

**Koda Binding:**
```koda
native bool();
```

### bool

**Signature:**
```c
typedef bool()
```

**Koda Binding:**
```koda
native bool();
```

### InitWindow

**Signature:**
```c
RLAPI void InitWindow(int width, int height, const char *title)
```

**Koda Binding:**
```koda
native InitWindow(width, height, *title);
```

### CloseWindow

**Signature:**
```c
RLAPI void CloseWindow()
```

**Koda Binding:**
```koda
native CloseWindow();
```

### WindowShouldClose

**Signature:**
```c
RLAPI bool WindowShouldClose()
```

**Koda Binding:**
```koda
native WindowShouldClose();
```

### IsWindowReady

**Signature:**
```c
RLAPI bool IsWindowReady()
```

**Koda Binding:**
```koda
native IsWindowReady();
```

### IsWindowFullscreen

**Signature:**
```c
RLAPI bool IsWindowFullscreen()
```

**Koda Binding:**
```koda
native IsWindowFullscreen();
```

### IsWindowHidden

**Signature:**
```c
RLAPI bool IsWindowHidden()
```

**Koda Binding:**
```koda
native IsWindowHidden();
```

### IsWindowMinimized

**Signature:**
```c
RLAPI bool IsWindowMinimized()
```

**Koda Binding:**
```koda
native IsWindowMinimized();
```

### IsWindowMaximized

**Signature:**
```c
RLAPI bool IsWindowMaximized()
```

**Koda Binding:**
```koda
native IsWindowMaximized();
```

### IsWindowFocused

**Signature:**
```c
RLAPI bool IsWindowFocused()
```

**Koda Binding:**
```koda
native IsWindowFocused();
```

### IsWindowResized

**Signature:**
```c
RLAPI bool IsWindowResized()
```

**Koda Binding:**
```koda
native IsWindowResized();
```

### IsWindowState

**Signature:**
```c
RLAPI bool IsWindowState(unsigned int flag)
```

**Koda Binding:**
```koda
native IsWindowState(flag);
```

### SetWindowState

**Signature:**
```c
RLAPI void SetWindowState(unsigned int flags)
```

**Koda Binding:**
```koda
native SetWindowState(flags);
```

### ClearWindowState

**Signature:**
```c
RLAPI void ClearWindowState(unsigned int flags)
```

**Koda Binding:**
```koda
native ClearWindowState(flags);
```

### ToggleFullscreen

**Signature:**
```c
RLAPI void ToggleFullscreen()
```

**Koda Binding:**
```koda
native ToggleFullscreen();
```

### ToggleBorderlessWindowed

**Signature:**
```c
RLAPI void ToggleBorderlessWindowed()
```

**Koda Binding:**
```koda
native ToggleBorderlessWindowed();
```

### MaximizeWindow

**Signature:**
```c
RLAPI void MaximizeWindow()
```

**Koda Binding:**
```koda
native MaximizeWindow();
```

### MinimizeWindow

**Signature:**
```c
RLAPI void MinimizeWindow()
```

**Koda Binding:**
```koda
native MinimizeWindow();
```

### RestoreWindow

**Signature:**
```c
RLAPI void RestoreWindow()
```

**Koda Binding:**
```koda
native RestoreWindow();
```

### SetWindowIcon

**Signature:**
```c
RLAPI void SetWindowIcon(Image image)
```

**Koda Binding:**
```koda
native SetWindowIcon(image);
```

### SetWindowIcons

**Signature:**
```c
RLAPI void SetWindowIcons(Image *images, int count)
```

**Koda Binding:**
```koda
native SetWindowIcons(*images, count);
```

### SetWindowTitle

**Signature:**
```c
RLAPI void SetWindowTitle(const char *title)
```

**Koda Binding:**
```koda
native SetWindowTitle(*title);
```

### SetWindowPosition

**Signature:**
```c
RLAPI void SetWindowPosition(int x, int y)
```

**Koda Binding:**
```koda
native SetWindowPosition(x, y);
```

### SetWindowMonitor

**Signature:**
```c
RLAPI void SetWindowMonitor(int monitor)
```

**Koda Binding:**
```koda
native SetWindowMonitor(monitor);
```

### SetWindowMinSize

**Signature:**
```c
RLAPI void SetWindowMinSize(int width, int height)
```

**Koda Binding:**
```koda
native SetWindowMinSize(width, height);
```

### SetWindowMaxSize

**Signature:**
```c
RLAPI void SetWindowMaxSize(int width, int height)
```

**Koda Binding:**
```koda
native SetWindowMaxSize(width, height);
```

### SetWindowSize

**Signature:**
```c
RLAPI void SetWindowSize(int width, int height)
```

**Koda Binding:**
```koda
native SetWindowSize(width, height);
```

### SetWindowOpacity

**Signature:**
```c
RLAPI void SetWindowOpacity(float opacity)
```

**Koda Binding:**
```koda
native SetWindowOpacity(opacity);
```

### SetWindowFocused

**Signature:**
```c
RLAPI void SetWindowFocused()
```

**Koda Binding:**
```koda
native SetWindowFocused();
```

### GetScreenWidth

**Signature:**
```c
RLAPI int GetScreenWidth()
```

**Koda Binding:**
```koda
native GetScreenWidth();
```

### GetScreenHeight

**Signature:**
```c
RLAPI int GetScreenHeight()
```

**Koda Binding:**
```koda
native GetScreenHeight();
```

### GetRenderWidth

**Signature:**
```c
RLAPI int GetRenderWidth()
```

**Koda Binding:**
```koda
native GetRenderWidth();
```

### GetRenderHeight

**Signature:**
```c
RLAPI int GetRenderHeight()
```

**Koda Binding:**
```koda
native GetRenderHeight();
```

### GetMonitorCount

**Signature:**
```c
RLAPI int GetMonitorCount()
```

**Koda Binding:**
```koda
native GetMonitorCount();
```

### GetCurrentMonitor

**Signature:**
```c
RLAPI int GetCurrentMonitor()
```

**Koda Binding:**
```koda
native GetCurrentMonitor();
```

### GetMonitorPosition

**Signature:**
```c
RLAPI Vector2 GetMonitorPosition(int monitor)
```

**Koda Binding:**
```koda
native GetMonitorPosition(monitor);
```

### GetMonitorWidth

**Signature:**
```c
RLAPI int GetMonitorWidth(int monitor)
```

**Koda Binding:**
```koda
native GetMonitorWidth(monitor);
```

### GetMonitorHeight

**Signature:**
```c
RLAPI int GetMonitorHeight(int monitor)
```

**Koda Binding:**
```koda
native GetMonitorHeight(monitor);
```

### GetMonitorPhysicalWidth

**Signature:**
```c
RLAPI int GetMonitorPhysicalWidth(int monitor)
```

**Koda Binding:**
```koda
native GetMonitorPhysicalWidth(monitor);
```

### GetMonitorPhysicalHeight

**Signature:**
```c
RLAPI int GetMonitorPhysicalHeight(int monitor)
```

**Koda Binding:**
```koda
native GetMonitorPhysicalHeight(monitor);
```

### GetMonitorRefreshRate

**Signature:**
```c
RLAPI int GetMonitorRefreshRate(int monitor)
```

**Koda Binding:**
```koda
native GetMonitorRefreshRate(monitor);
```

### GetWindowPosition

**Signature:**
```c
RLAPI Vector2 GetWindowPosition()
```

**Koda Binding:**
```koda
native GetWindowPosition();
```

### GetWindowScaleDPI

**Signature:**
```c
RLAPI Vector2 GetWindowScaleDPI()
```

**Koda Binding:**
```koda
native GetWindowScaleDPI();
```

### SetClipboardText

**Signature:**
```c
RLAPI void SetClipboardText(const char *text)
```

**Koda Binding:**
```koda
native SetClipboardText(*text);
```

### GetClipboardImage

**Signature:**
```c
RLAPI Image GetClipboardImage()
```

**Koda Binding:**
```koda
native GetClipboardImage();
```

### EnableEventWaiting

**Signature:**
```c
RLAPI void EnableEventWaiting()
```

**Koda Binding:**
```koda
native EnableEventWaiting();
```

### DisableEventWaiting

**Signature:**
```c
RLAPI void DisableEventWaiting()
```

**Koda Binding:**
```koda
native DisableEventWaiting();
```

### ShowCursor

**Signature:**
```c
RLAPI void ShowCursor()
```

**Koda Binding:**
```koda
native ShowCursor();
```

### HideCursor

**Signature:**
```c
RLAPI void HideCursor()
```

**Koda Binding:**
```koda
native HideCursor();
```

### IsCursorHidden

**Signature:**
```c
RLAPI bool IsCursorHidden()
```

**Koda Binding:**
```koda
native IsCursorHidden();
```

### EnableCursor

**Signature:**
```c
RLAPI void EnableCursor()
```

**Koda Binding:**
```koda
native EnableCursor();
```

### DisableCursor

**Signature:**
```c
RLAPI void DisableCursor()
```

**Koda Binding:**
```koda
native DisableCursor();
```

### IsCursorOnScreen

**Signature:**
```c
RLAPI bool IsCursorOnScreen()
```

**Koda Binding:**
```koda
native IsCursorOnScreen();
```

### ClearBackground

**Signature:**
```c
RLAPI void ClearBackground(Color color)
```

**Koda Binding:**
```koda
native ClearBackground(color);
```

### BeginDrawing

**Signature:**
```c
RLAPI void BeginDrawing()
```

**Koda Binding:**
```koda
native BeginDrawing();
```

### EndDrawing

**Signature:**
```c
RLAPI void EndDrawing()
```

**Koda Binding:**
```koda
native EndDrawing();
```

### BeginMode2D

**Signature:**
```c
RLAPI void BeginMode2D(Camera2D camera)
```

**Koda Binding:**
```koda
native BeginMode2D(camera);
```

### EndMode2D

**Signature:**
```c
RLAPI void EndMode2D()
```

**Koda Binding:**
```koda
native EndMode2D();
```

### BeginMode3D

**Signature:**
```c
RLAPI void BeginMode3D(Camera3D camera)
```

**Koda Binding:**
```koda
native BeginMode3D(camera);
```

### EndMode3D

**Signature:**
```c
RLAPI void EndMode3D()
```

**Koda Binding:**
```koda
native EndMode3D();
```

### BeginTextureMode

**Signature:**
```c
RLAPI void BeginTextureMode(RenderTexture2D target)
```

**Koda Binding:**
```koda
native BeginTextureMode(target);
```

### EndTextureMode

**Signature:**
```c
RLAPI void EndTextureMode()
```

**Koda Binding:**
```koda
native EndTextureMode();
```

### BeginShaderMode

**Signature:**
```c
RLAPI void BeginShaderMode(Shader shader)
```

**Koda Binding:**
```koda
native BeginShaderMode(shader);
```

### EndShaderMode

**Signature:**
```c
RLAPI void EndShaderMode()
```

**Koda Binding:**
```koda
native EndShaderMode();
```

### BeginBlendMode

**Signature:**
```c
RLAPI void BeginBlendMode(int mode)
```

**Koda Binding:**
```koda
native BeginBlendMode(mode);
```

### EndBlendMode

**Signature:**
```c
RLAPI void EndBlendMode()
```

**Koda Binding:**
```koda
native EndBlendMode();
```

### BeginScissorMode

**Signature:**
```c
RLAPI void BeginScissorMode(int x, int y, int width, int height)
```

**Koda Binding:**
```koda
native BeginScissorMode(x, y, width, height);
```

### EndScissorMode

**Signature:**
```c
RLAPI void EndScissorMode()
```

**Koda Binding:**
```koda
native EndScissorMode();
```

### BeginVrStereoMode

**Signature:**
```c
RLAPI void BeginVrStereoMode(VrStereoConfig config)
```

**Koda Binding:**
```koda
native BeginVrStereoMode(config);
```

### EndVrStereoMode

**Signature:**
```c
RLAPI void EndVrStereoMode()
```

**Koda Binding:**
```koda
native EndVrStereoMode();
```

### LoadVrStereoConfig

**Signature:**
```c
RLAPI VrStereoConfig LoadVrStereoConfig(VrDeviceInfo device)
```

**Koda Binding:**
```koda
native LoadVrStereoConfig(device);
```

### UnloadVrStereoConfig

**Signature:**
```c
RLAPI void UnloadVrStereoConfig(VrStereoConfig config)
```

**Koda Binding:**
```koda
native UnloadVrStereoConfig(config);
```

### LoadShader

**Signature:**
```c
RLAPI Shader LoadShader(const char *vsFileName, const char *fsFileName)
```

**Koda Binding:**
```koda
native LoadShader(*vsFileName, *fsFileName);
```

### LoadShaderFromMemory

**Signature:**
```c
RLAPI Shader LoadShaderFromMemory(const char *vsCode, const char *fsCode)
```

**Koda Binding:**
```koda
native LoadShaderFromMemory(*vsCode, *fsCode);
```

### IsShaderValid

**Signature:**
```c
RLAPI bool IsShaderValid(Shader shader)
```

**Koda Binding:**
```koda
native IsShaderValid(shader);
```

### GetShaderLocation

**Signature:**
```c
RLAPI int GetShaderLocation(Shader shader, const char *uniformName)
```

**Koda Binding:**
```koda
native GetShaderLocation(shader, *uniformName);
```

### GetShaderLocationAttrib

**Signature:**
```c
RLAPI int GetShaderLocationAttrib(Shader shader, const char *attribName)
```

**Koda Binding:**
```koda
native GetShaderLocationAttrib(shader, *attribName);
```

### SetShaderValue

**Signature:**
```c
RLAPI void SetShaderValue(Shader shader, int locIndex, const void *value, int uniformType)
```

**Koda Binding:**
```koda
native SetShaderValue(shader, locIndex, *value, uniformType);
```

### SetShaderValueV

**Signature:**
```c
RLAPI void SetShaderValueV(Shader shader, int locIndex, const void *value, int uniformType, int count)
```

**Koda Binding:**
```koda
native SetShaderValueV(shader, locIndex, *value, uniformType, count);
```

### SetShaderValueMatrix

**Signature:**
```c
RLAPI void SetShaderValueMatrix(Shader shader, int locIndex, Matrix mat)
```

**Koda Binding:**
```koda
native SetShaderValueMatrix(shader, locIndex, mat);
```

### SetShaderValueTexture

**Signature:**
```c
RLAPI void SetShaderValueTexture(Shader shader, int locIndex, Texture2D texture)
```

**Koda Binding:**
```koda
native SetShaderValueTexture(shader, locIndex, texture);
```

### UnloadShader

**Signature:**
```c
RLAPI void UnloadShader(Shader shader)
```

**Koda Binding:**
```koda
native UnloadShader(shader);
```

### GetScreenToWorldRay

**Signature:**
```c
RLAPI Ray GetScreenToWorldRay(Vector2 position, Camera camera)
```

**Koda Binding:**
```koda
native GetScreenToWorldRay(position, camera);
```

### GetScreenToWorldRayEx

**Signature:**
```c
RLAPI Ray GetScreenToWorldRayEx(Vector2 position, Camera camera, int width, int height)
```

**Koda Binding:**
```koda
native GetScreenToWorldRayEx(position, camera, width, height);
```

### GetWorldToScreen

**Signature:**
```c
RLAPI Vector2 GetWorldToScreen(Vector3 position, Camera camera)
```

**Koda Binding:**
```koda
native GetWorldToScreen(position, camera);
```

### GetWorldToScreenEx

**Signature:**
```c
RLAPI Vector2 GetWorldToScreenEx(Vector3 position, Camera camera, int width, int height)
```

**Koda Binding:**
```koda
native GetWorldToScreenEx(position, camera, width, height);
```

### GetWorldToScreen2D

**Signature:**
```c
RLAPI Vector2 GetWorldToScreen2D(Vector2 position, Camera2D camera)
```

**Koda Binding:**
```koda
native GetWorldToScreen2D(position, camera);
```

### GetScreenToWorld2D

**Signature:**
```c
RLAPI Vector2 GetScreenToWorld2D(Vector2 position, Camera2D camera)
```

**Koda Binding:**
```koda
native GetScreenToWorld2D(position, camera);
```

### GetCameraMatrix

**Signature:**
```c
RLAPI Matrix GetCameraMatrix(Camera camera)
```

**Koda Binding:**
```koda
native GetCameraMatrix(camera);
```

### GetCameraMatrix2D

**Signature:**
```c
RLAPI Matrix GetCameraMatrix2D(Camera2D camera)
```

**Koda Binding:**
```koda
native GetCameraMatrix2D(camera);
```

### SetTargetFPS

**Signature:**
```c
RLAPI void SetTargetFPS(int fps)
```

**Koda Binding:**
```koda
native SetTargetFPS(fps);
```

### GetFrameTime

**Signature:**
```c
RLAPI float GetFrameTime()
```

**Koda Binding:**
```koda
native GetFrameTime();
```

### GetTime

**Signature:**
```c
RLAPI double GetTime()
```

**Koda Binding:**
```koda
native GetTime();
```

### GetFPS

**Signature:**
```c
RLAPI int GetFPS()
```

**Koda Binding:**
```koda
native GetFPS();
```

### SwapScreenBuffer

**Signature:**
```c
RLAPI void SwapScreenBuffer()
```

**Koda Binding:**
```koda
native SwapScreenBuffer();
```

### PollInputEvents

**Signature:**
```c
RLAPI void PollInputEvents()
```

**Koda Binding:**
```koda
native PollInputEvents();
```

### WaitTime

**Signature:**
```c
RLAPI void WaitTime(double seconds)
```

**Koda Binding:**
```koda
native WaitTime(seconds);
```

### SetRandomSeed

**Signature:**
```c
RLAPI void SetRandomSeed(unsigned int seed)
```

**Koda Binding:**
```koda
native SetRandomSeed(seed);
```

### GetRandomValue

**Signature:**
```c
RLAPI int GetRandomValue(int min, int max)
```

**Koda Binding:**
```koda
native GetRandomValue(min, max);
```

### UnloadRandomSequence

**Signature:**
```c
RLAPI void UnloadRandomSequence(int *sequence)
```

**Koda Binding:**
```koda
native UnloadRandomSequence(*sequence);
```

### TakeScreenshot

**Signature:**
```c
RLAPI void TakeScreenshot(const char *fileName)
```

**Koda Binding:**
```koda
native TakeScreenshot(*fileName);
```

### SetConfigFlags

**Signature:**
```c
RLAPI void SetConfigFlags(unsigned int flags)
```

**Koda Binding:**
```koda
native SetConfigFlags(flags);
```

### OpenURL

**Signature:**
```c
RLAPI void OpenURL(const char *url)
```

**Koda Binding:**
```koda
native OpenURL(*url);
```

### SetTraceLogLevel

**Signature:**
```c
RLAPI void SetTraceLogLevel(int logLevel)
```

**Koda Binding:**
```koda
native SetTraceLogLevel(logLevel);
```

### TraceLog

**Signature:**
```c
RLAPI void TraceLog(int logLevel, const char *text)
```

**Koda Binding:**
```koda
native TraceLog(logLevel, *text);
```

### SetTraceLogCallback

**Signature:**
```c
RLAPI void SetTraceLogCallback(TraceLogCallback callback)
```

**Koda Binding:**
```koda
native SetTraceLogCallback(callback);
```

### MemFree

**Signature:**
```c
RLAPI void MemFree(void *ptr)
```

**Koda Binding:**
```koda
native MemFree(*ptr);
```

### UnloadFileData

**Signature:**
```c
RLAPI void UnloadFileData(unsigned char *data)
```

**Koda Binding:**
```koda
native UnloadFileData(*data);
```

### SaveFileData

**Signature:**
```c
RLAPI bool SaveFileData(const char *fileName, const void *data, int dataSize)
```

**Koda Binding:**
```koda
native SaveFileData(*fileName, *data, dataSize);
```

### ExportDataAsCode

**Signature:**
```c
RLAPI bool ExportDataAsCode(const unsigned char *data, int dataSize, const char *fileName)
```

**Koda Binding:**
```koda
native ExportDataAsCode(*data, dataSize, *fileName);
```

### UnloadFileText

**Signature:**
```c
RLAPI void UnloadFileText(char *text)
```

**Koda Binding:**
```koda
native UnloadFileText(*text);
```

### SaveFileText

**Signature:**
```c
RLAPI bool SaveFileText(const char *fileName, const char *text)
```

**Koda Binding:**
```koda
native SaveFileText(*fileName, *text);
```

### SetLoadFileDataCallback

**Signature:**
```c
RLAPI void SetLoadFileDataCallback(LoadFileDataCallback callback)
```

**Koda Binding:**
```koda
native SetLoadFileDataCallback(callback);
```

### SetSaveFileDataCallback

**Signature:**
```c
RLAPI void SetSaveFileDataCallback(SaveFileDataCallback callback)
```

**Koda Binding:**
```koda
native SetSaveFileDataCallback(callback);
```

### SetLoadFileTextCallback

**Signature:**
```c
RLAPI void SetLoadFileTextCallback(LoadFileTextCallback callback)
```

**Koda Binding:**
```koda
native SetLoadFileTextCallback(callback);
```

### SetSaveFileTextCallback

**Signature:**
```c
RLAPI void SetSaveFileTextCallback(SaveFileTextCallback callback)
```

**Koda Binding:**
```koda
native SetSaveFileTextCallback(callback);
```

### FileRename

**Signature:**
```c
RLAPI int FileRename(const char *fileName, const char *fileRename)
```

**Koda Binding:**
```koda
native FileRename(*fileName, *fileRename);
```

### FileRemove

**Signature:**
```c
RLAPI int FileRemove(const char *fileName)
```

**Koda Binding:**
```koda
native FileRemove(*fileName);
```

### FileCopy

**Signature:**
```c
RLAPI int FileCopy(const char *srcPath, const char *dstPath)
```

**Koda Binding:**
```koda
native FileCopy(*srcPath, *dstPath);
```

### FileMove

**Signature:**
```c
RLAPI int FileMove(const char *srcPath, const char *dstPath)
```

**Koda Binding:**
```koda
native FileMove(*srcPath, *dstPath);
```

### FileTextReplace

**Signature:**
```c
RLAPI int FileTextReplace(const char *fileName, const char *search, const char *replacement)
```

**Koda Binding:**
```koda
native FileTextReplace(*fileName, *search, *replacement);
```

### FileTextFindIndex

**Signature:**
```c
RLAPI int FileTextFindIndex(const char *fileName, const char *search)
```

**Koda Binding:**
```koda
native FileTextFindIndex(*fileName, *search);
```

### FileExists

**Signature:**
```c
RLAPI bool FileExists(const char *fileName)
```

**Koda Binding:**
```koda
native FileExists(*fileName);
```

### DirectoryExists

**Signature:**
```c
RLAPI bool DirectoryExists(const char *dirPath)
```

**Koda Binding:**
```koda
native DirectoryExists(*dirPath);
```

### IsFileExtension

**Signature:**
```c
RLAPI bool IsFileExtension(const char *fileName, const char *ext)
```

**Koda Binding:**
```koda
native IsFileExtension(*fileName, *ext);
```

### GetFileLength

**Signature:**
```c
RLAPI int GetFileLength(const char *fileName)
```

**Koda Binding:**
```koda
native GetFileLength(*fileName);
```

### GetFileModTime

**Signature:**
```c
RLAPI long GetFileModTime(const char *fileName)
```

**Koda Binding:**
```koda
native GetFileModTime(*fileName);
```

### MakeDirectory

**Signature:**
```c
RLAPI int MakeDirectory(const char *dirPath)
```

**Koda Binding:**
```koda
native MakeDirectory(*dirPath);
```

### ChangeDirectory

**Signature:**
```c
RLAPI bool ChangeDirectory(const char *dirPath)
```

**Koda Binding:**
```koda
native ChangeDirectory(*dirPath);
```

### IsPathFile

**Signature:**
```c
RLAPI bool IsPathFile(const char *path)
```

**Koda Binding:**
```koda
native IsPathFile(*path);
```

### IsFileNameValid

**Signature:**
```c
RLAPI bool IsFileNameValid(const char *fileName)
```

**Koda Binding:**
```koda
native IsFileNameValid(*fileName);
```

### LoadDirectoryFiles

**Signature:**
```c
RLAPI FilePathList LoadDirectoryFiles(const char *dirPath)
```

**Koda Binding:**
```koda
native LoadDirectoryFiles(*dirPath);
```

### LoadDirectoryFilesEx

**Signature:**
```c
RLAPI FilePathList LoadDirectoryFilesEx(const char *basePath, const char *filter, bool scanSubdirs)
```

**Koda Binding:**
```koda
native LoadDirectoryFilesEx(*basePath, *filter, scanSubdirs);
```

### UnloadDirectoryFiles

**Signature:**
```c
RLAPI void UnloadDirectoryFiles(FilePathList files)
```

**Koda Binding:**
```koda
native UnloadDirectoryFiles(files);
```

### IsFileDropped

**Signature:**
```c
RLAPI bool IsFileDropped()
```

**Koda Binding:**
```koda
native IsFileDropped();
```

### LoadDroppedFiles

**Signature:**
```c
RLAPI FilePathList LoadDroppedFiles()
```

**Koda Binding:**
```koda
native LoadDroppedFiles();
```

### UnloadDroppedFiles

**Signature:**
```c
RLAPI void UnloadDroppedFiles(FilePathList files)
```

**Koda Binding:**
```koda
native UnloadDroppedFiles(files);
```

### GetDirectoryFileCount

**Signature:**
```c
RLAPI unsigned int GetDirectoryFileCount(const char *dirPath)
```

**Koda Binding:**
```koda
native GetDirectoryFileCount(*dirPath);
```

### GetDirectoryFileCountEx

**Signature:**
```c
RLAPI unsigned int GetDirectoryFileCountEx(const char *basePath, const char *filter, bool scanSubdirs)
```

**Koda Binding:**
```koda
native GetDirectoryFileCountEx(*basePath, *filter, scanSubdirs);
```

### ComputeCRC32

**Signature:**
```c
RLAPI unsigned int ComputeCRC32(const unsigned char *data, int dataSize)
```

**Koda Binding:**
```koda
native ComputeCRC32(*data, dataSize);
```

### LoadAutomationEventList

**Signature:**
```c
RLAPI AutomationEventList LoadAutomationEventList(const char *fileName)
```

**Koda Binding:**
```koda
native LoadAutomationEventList(*fileName);
```

### UnloadAutomationEventList

**Signature:**
```c
RLAPI void UnloadAutomationEventList(AutomationEventList list)
```

**Koda Binding:**
```koda
native UnloadAutomationEventList(list);
```

### ExportAutomationEventList

**Signature:**
```c
RLAPI bool ExportAutomationEventList(AutomationEventList list, const char *fileName)
```

**Koda Binding:**
```koda
native ExportAutomationEventList(list, *fileName);
```

### SetAutomationEventList

**Signature:**
```c
RLAPI void SetAutomationEventList(AutomationEventList *list)
```

**Koda Binding:**
```koda
native SetAutomationEventList(*list);
```

### SetAutomationEventBaseFrame

**Signature:**
```c
RLAPI void SetAutomationEventBaseFrame(int frame)
```

**Koda Binding:**
```koda
native SetAutomationEventBaseFrame(frame);
```

### StartAutomationEventRecording

**Signature:**
```c
RLAPI void StartAutomationEventRecording()
```

**Koda Binding:**
```koda
native StartAutomationEventRecording();
```

### StopAutomationEventRecording

**Signature:**
```c
RLAPI void StopAutomationEventRecording()
```

**Koda Binding:**
```koda
native StopAutomationEventRecording();
```

### PlayAutomationEvent

**Signature:**
```c
RLAPI void PlayAutomationEvent(AutomationEvent event)
```

**Koda Binding:**
```koda
native PlayAutomationEvent(event);
```

### IsKeyPressed

**Signature:**
```c
RLAPI bool IsKeyPressed(int key)
```

**Koda Binding:**
```koda
native IsKeyPressed(key);
```

### IsKeyPressedRepeat

**Signature:**
```c
RLAPI bool IsKeyPressedRepeat(int key)
```

**Koda Binding:**
```koda
native IsKeyPressedRepeat(key);
```

### IsKeyDown

**Signature:**
```c
RLAPI bool IsKeyDown(int key)
```

**Koda Binding:**
```koda
native IsKeyDown(key);
```

### IsKeyReleased

**Signature:**
```c
RLAPI bool IsKeyReleased(int key)
```

**Koda Binding:**
```koda
native IsKeyReleased(key);
```

### IsKeyUp

**Signature:**
```c
RLAPI bool IsKeyUp(int key)
```

**Koda Binding:**
```koda
native IsKeyUp(key);
```

### GetKeyPressed

**Signature:**
```c
RLAPI int GetKeyPressed()
```

**Koda Binding:**
```koda
native GetKeyPressed();
```

### GetCharPressed

**Signature:**
```c
RLAPI int GetCharPressed()
```

**Koda Binding:**
```koda
native GetCharPressed();
```

### SetExitKey

**Signature:**
```c
RLAPI void SetExitKey(int key)
```

**Koda Binding:**
```koda
native SetExitKey(key);
```

### IsGamepadAvailable

**Signature:**
```c
RLAPI bool IsGamepadAvailable(int gamepad)
```

**Koda Binding:**
```koda
native IsGamepadAvailable(gamepad);
```

### IsGamepadButtonPressed

**Signature:**
```c
RLAPI bool IsGamepadButtonPressed(int gamepad, int button)
```

**Koda Binding:**
```koda
native IsGamepadButtonPressed(gamepad, button);
```

### IsGamepadButtonDown

**Signature:**
```c
RLAPI bool IsGamepadButtonDown(int gamepad, int button)
```

**Koda Binding:**
```koda
native IsGamepadButtonDown(gamepad, button);
```

### IsGamepadButtonReleased

**Signature:**
```c
RLAPI bool IsGamepadButtonReleased(int gamepad, int button)
```

**Koda Binding:**
```koda
native IsGamepadButtonReleased(gamepad, button);
```

### IsGamepadButtonUp

**Signature:**
```c
RLAPI bool IsGamepadButtonUp(int gamepad, int button)
```

**Koda Binding:**
```koda
native IsGamepadButtonUp(gamepad, button);
```

### GetGamepadButtonPressed

**Signature:**
```c
RLAPI int GetGamepadButtonPressed()
```

**Koda Binding:**
```koda
native GetGamepadButtonPressed();
```

### GetGamepadAxisCount

**Signature:**
```c
RLAPI int GetGamepadAxisCount(int gamepad)
```

**Koda Binding:**
```koda
native GetGamepadAxisCount(gamepad);
```

### GetGamepadAxisMovement

**Signature:**
```c
RLAPI float GetGamepadAxisMovement(int gamepad, int axis)
```

**Koda Binding:**
```koda
native GetGamepadAxisMovement(gamepad, axis);
```

### SetGamepadMappings

**Signature:**
```c
RLAPI int SetGamepadMappings(const char *mappings)
```

**Koda Binding:**
```koda
native SetGamepadMappings(*mappings);
```

### SetGamepadVibration

**Signature:**
```c
RLAPI void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
```

**Koda Binding:**
```koda
native SetGamepadVibration(gamepad, leftMotor, rightMotor, duration);
```

### IsMouseButtonPressed

**Signature:**
```c
RLAPI bool IsMouseButtonPressed(int button)
```

**Koda Binding:**
```koda
native IsMouseButtonPressed(button);
```

### IsMouseButtonDown

**Signature:**
```c
RLAPI bool IsMouseButtonDown(int button)
```

**Koda Binding:**
```koda
native IsMouseButtonDown(button);
```

### IsMouseButtonReleased

**Signature:**
```c
RLAPI bool IsMouseButtonReleased(int button)
```

**Koda Binding:**
```koda
native IsMouseButtonReleased(button);
```

### IsMouseButtonUp

**Signature:**
```c
RLAPI bool IsMouseButtonUp(int button)
```

**Koda Binding:**
```koda
native IsMouseButtonUp(button);
```

### GetMouseX

**Signature:**
```c
RLAPI int GetMouseX()
```

**Koda Binding:**
```koda
native GetMouseX();
```

### GetMouseY

**Signature:**
```c
RLAPI int GetMouseY()
```

**Koda Binding:**
```koda
native GetMouseY();
```

### GetMousePosition

**Signature:**
```c
RLAPI Vector2 GetMousePosition()
```

**Koda Binding:**
```koda
native GetMousePosition();
```

### GetMouseDelta

**Signature:**
```c
RLAPI Vector2 GetMouseDelta()
```

**Koda Binding:**
```koda
native GetMouseDelta();
```

### SetMousePosition

**Signature:**
```c
RLAPI void SetMousePosition(int x, int y)
```

**Koda Binding:**
```koda
native SetMousePosition(x, y);
```

### SetMouseOffset

**Signature:**
```c
RLAPI void SetMouseOffset(int offsetX, int offsetY)
```

**Koda Binding:**
```koda
native SetMouseOffset(offsetX, offsetY);
```

### SetMouseScale

**Signature:**
```c
RLAPI void SetMouseScale(float scaleX, float scaleY)
```

**Koda Binding:**
```koda
native SetMouseScale(scaleX, scaleY);
```

### GetMouseWheelMove

**Signature:**
```c
RLAPI float GetMouseWheelMove()
```

**Koda Binding:**
```koda
native GetMouseWheelMove();
```

### GetMouseWheelMoveV

**Signature:**
```c
RLAPI Vector2 GetMouseWheelMoveV()
```

**Koda Binding:**
```koda
native GetMouseWheelMoveV();
```

### SetMouseCursor

**Signature:**
```c
RLAPI void SetMouseCursor(int cursor)
```

**Koda Binding:**
```koda
native SetMouseCursor(cursor);
```

### GetTouchX

**Signature:**
```c
RLAPI int GetTouchX()
```

**Koda Binding:**
```koda
native GetTouchX();
```

### GetTouchY

**Signature:**
```c
RLAPI int GetTouchY()
```

**Koda Binding:**
```koda
native GetTouchY();
```

### GetTouchPosition

**Signature:**
```c
RLAPI Vector2 GetTouchPosition(int index)
```

**Koda Binding:**
```koda
native GetTouchPosition(index);
```

### GetTouchPointId

**Signature:**
```c
RLAPI int GetTouchPointId(int index)
```

**Koda Binding:**
```koda
native GetTouchPointId(index);
```

### GetTouchPointCount

**Signature:**
```c
RLAPI int GetTouchPointCount()
```

**Koda Binding:**
```koda
native GetTouchPointCount();
```

### SetGesturesEnabled

**Signature:**
```c
RLAPI void SetGesturesEnabled(unsigned int flags)
```

**Koda Binding:**
```koda
native SetGesturesEnabled(flags);
```

### IsGestureDetected

**Signature:**
```c
RLAPI bool IsGestureDetected(unsigned int gesture)
```

**Koda Binding:**
```koda
native IsGestureDetected(gesture);
```

### GetGestureDetected

**Signature:**
```c
RLAPI int GetGestureDetected()
```

**Koda Binding:**
```koda
native GetGestureDetected();
```

### GetGestureHoldDuration

**Signature:**
```c
RLAPI float GetGestureHoldDuration()
```

**Koda Binding:**
```koda
native GetGestureHoldDuration();
```

### GetGestureDragVector

**Signature:**
```c
RLAPI Vector2 GetGestureDragVector()
```

**Koda Binding:**
```koda
native GetGestureDragVector();
```

### GetGestureDragAngle

**Signature:**
```c
RLAPI float GetGestureDragAngle()
```

**Koda Binding:**
```koda
native GetGestureDragAngle();
```

### GetGesturePinchVector

**Signature:**
```c
RLAPI Vector2 GetGesturePinchVector()
```

**Koda Binding:**
```koda
native GetGesturePinchVector();
```

### GetGesturePinchAngle

**Signature:**
```c
RLAPI float GetGesturePinchAngle()
```

**Koda Binding:**
```koda
native GetGesturePinchAngle();
```

### UpdateCamera

**Signature:**
```c
RLAPI void UpdateCamera(Camera *camera, int mode)
```

**Koda Binding:**
```koda
native UpdateCamera(*camera, mode);
```

### UpdateCameraPro

**Signature:**
```c
RLAPI void UpdateCameraPro(Camera *camera, Vector3 movement, Vector3 rotation, float zoom)
```

**Koda Binding:**
```koda
native UpdateCameraPro(*camera, movement, rotation, zoom);
```

### SetShapesTexture

**Signature:**
```c
RLAPI void SetShapesTexture(Texture2D texture, Rectangle source)
```

**Koda Binding:**
```koda
native SetShapesTexture(texture, source);
```

### GetShapesTexture

**Signature:**
```c
RLAPI Texture2D GetShapesTexture()
```

**Koda Binding:**
```koda
native GetShapesTexture();
```

### GetShapesTextureRectangle

**Signature:**
```c
RLAPI Rectangle GetShapesTextureRectangle()
```

**Koda Binding:**
```koda
native GetShapesTextureRectangle();
```

### DrawPixel

**Signature:**
```c
RLAPI void DrawPixel(int posX, int posY, Color color)
```

**Koda Binding:**
```koda
native DrawPixel(posX, posY, color);
```

### DrawPixelV

**Signature:**
```c
RLAPI void DrawPixelV(Vector2 position, Color color)
```

**Koda Binding:**
```koda
native DrawPixelV(position, color);
```

### DrawLine

**Signature:**
```c
RLAPI void DrawLine(int startPosX, int startPosY, int endPosX, int endPosY, Color color)
```

**Koda Binding:**
```koda
native DrawLine(startPosX, startPosY, endPosX, endPosY, color);
```

### DrawLineV

**Signature:**
```c
RLAPI void DrawLineV(Vector2 startPos, Vector2 endPos, Color color)
```

**Koda Binding:**
```koda
native DrawLineV(startPos, endPos, color);
```

### DrawLineEx

**Signature:**
```c
RLAPI void DrawLineEx(Vector2 startPos, Vector2 endPos, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawLineEx(startPos, endPos, thick, color);
```

### DrawLineStrip

**Signature:**
```c
RLAPI void DrawLineStrip(const Vector2 *points, int pointCount, Color color)
```

**Koda Binding:**
```koda
native DrawLineStrip(*points, pointCount, color);
```

### DrawLineBezier

**Signature:**
```c
RLAPI void DrawLineBezier(Vector2 startPos, Vector2 endPos, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawLineBezier(startPos, endPos, thick, color);
```

### DrawLineDashed

**Signature:**
```c
RLAPI void DrawLineDashed(Vector2 startPos, Vector2 endPos, int dashSize, int spaceSize, Color color)
```

**Koda Binding:**
```koda
native DrawLineDashed(startPos, endPos, dashSize, spaceSize, color);
```

### DrawCircle

**Signature:**
```c
RLAPI void DrawCircle(int centerX, int centerY, float radius, Color color)
```

**Koda Binding:**
```koda
native DrawCircle(centerX, centerY, radius, color);
```

### DrawCircleV

**Signature:**
```c
RLAPI void DrawCircleV(Vector2 center, float radius, Color color)
```

**Koda Binding:**
```koda
native DrawCircleV(center, radius, color);
```

### DrawCircleGradient

**Signature:**
```c
RLAPI void DrawCircleGradient(Vector2 center, float radius, Color inner, Color outer)
```

**Koda Binding:**
```koda
native DrawCircleGradient(center, radius, inner, outer);
```

### DrawCircleSector

**Signature:**
```c
RLAPI void DrawCircleSector(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color)
```

**Koda Binding:**
```koda
native DrawCircleSector(center, radius, startAngle, endAngle, segments, color);
```

### DrawCircleSectorLines

**Signature:**
```c
RLAPI void DrawCircleSectorLines(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color)
```

**Koda Binding:**
```koda
native DrawCircleSectorLines(center, radius, startAngle, endAngle, segments, color);
```

### DrawCircleLines

**Signature:**
```c
RLAPI void DrawCircleLines(int centerX, int centerY, float radius, Color color)
```

**Koda Binding:**
```koda
native DrawCircleLines(centerX, centerY, radius, color);
```

### DrawCircleLinesV

**Signature:**
```c
RLAPI void DrawCircleLinesV(Vector2 center, float radius, Color color)
```

**Koda Binding:**
```koda
native DrawCircleLinesV(center, radius, color);
```

### DrawEllipse

**Signature:**
```c
RLAPI void DrawEllipse(int centerX, int centerY, float radiusH, float radiusV, Color color)
```

**Koda Binding:**
```koda
native DrawEllipse(centerX, centerY, radiusH, radiusV, color);
```

### DrawEllipseV

**Signature:**
```c
RLAPI void DrawEllipseV(Vector2 center, float radiusH, float radiusV, Color color)
```

**Koda Binding:**
```koda
native DrawEllipseV(center, radiusH, radiusV, color);
```

### DrawEllipseLines

**Signature:**
```c
RLAPI void DrawEllipseLines(int centerX, int centerY, float radiusH, float radiusV, Color color)
```

**Koda Binding:**
```koda
native DrawEllipseLines(centerX, centerY, radiusH, radiusV, color);
```

### DrawEllipseLinesV

**Signature:**
```c
RLAPI void DrawEllipseLinesV(Vector2 center, float radiusH, float radiusV, Color color)
```

**Koda Binding:**
```koda
native DrawEllipseLinesV(center, radiusH, radiusV, color);
```

### DrawRing

**Signature:**
```c
RLAPI void DrawRing(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color)
```

**Koda Binding:**
```koda
native DrawRing(center, innerRadius, outerRadius, startAngle, endAngle, segments, color);
```

### DrawRingLines

**Signature:**
```c
RLAPI void DrawRingLines(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color)
```

**Koda Binding:**
```koda
native DrawRingLines(center, innerRadius, outerRadius, startAngle, endAngle, segments, color);
```

### DrawRectangle

**Signature:**
```c
RLAPI void DrawRectangle(int posX, int posY, int width, int height, Color color)
```

**Koda Binding:**
```koda
native DrawRectangle(posX, posY, width, height, color);
```

### DrawRectangleV

**Signature:**
```c
RLAPI void DrawRectangleV(Vector2 position, Vector2 size, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleV(position, size, color);
```

### DrawRectangleRec

**Signature:**
```c
RLAPI void DrawRectangleRec(Rectangle rec, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleRec(rec, color);
```

### DrawRectanglePro

**Signature:**
```c
RLAPI void DrawRectanglePro(Rectangle rec, Vector2 origin, float rotation, Color color)
```

**Koda Binding:**
```koda
native DrawRectanglePro(rec, origin, rotation, color);
```

### DrawRectangleGradientV

**Signature:**
```c
RLAPI void DrawRectangleGradientV(int posX, int posY, int width, int height, Color top, Color bottom)
```

**Koda Binding:**
```koda
native DrawRectangleGradientV(posX, posY, width, height, top, bottom);
```

### DrawRectangleGradientH

**Signature:**
```c
RLAPI void DrawRectangleGradientH(int posX, int posY, int width, int height, Color left, Color right)
```

**Koda Binding:**
```koda
native DrawRectangleGradientH(posX, posY, width, height, left, right);
```

### DrawRectangleGradientEx

**Signature:**
```c
RLAPI void DrawRectangleGradientEx(Rectangle rec, Color topLeft, Color bottomLeft, Color bottomRight, Color topRight)
```

**Koda Binding:**
```koda
native DrawRectangleGradientEx(rec, topLeft, bottomLeft, bottomRight, topRight);
```

### DrawRectangleLines

**Signature:**
```c
RLAPI void DrawRectangleLines(int posX, int posY, int width, int height, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleLines(posX, posY, width, height, color);
```

### DrawRectangleLinesEx

**Signature:**
```c
RLAPI void DrawRectangleLinesEx(Rectangle rec, float lineThick, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleLinesEx(rec, lineThick, color);
```

### DrawRectangleRounded

**Signature:**
```c
RLAPI void DrawRectangleRounded(Rectangle rec, float roundness, int segments, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleRounded(rec, roundness, segments, color);
```

### DrawRectangleRoundedLines

**Signature:**
```c
RLAPI void DrawRectangleRoundedLines(Rectangle rec, float roundness, int segments, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleRoundedLines(rec, roundness, segments, color);
```

### DrawRectangleRoundedLinesEx

**Signature:**
```c
RLAPI void DrawRectangleRoundedLinesEx(Rectangle rec, float roundness, int segments, float lineThick, Color color)
```

**Koda Binding:**
```koda
native DrawRectangleRoundedLinesEx(rec, roundness, segments, lineThick, color);
```

### DrawTriangle

**Signature:**
```c
RLAPI void DrawTriangle(Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda Binding:**
```koda
native DrawTriangle(v1, v2, v3, color);
```

### DrawTriangleGradient

**Signature:**
```c
RLAPI void DrawTriangleGradient(Vector2 v1, Vector2 v2, Vector2 v3, Color c1, Color c2, Color c3)
```

**Koda Binding:**
```koda
native DrawTriangleGradient(v1, v2, v3, c1, c2, c3);
```

### DrawTriangleLines

**Signature:**
```c
RLAPI void DrawTriangleLines(Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda Binding:**
```koda
native DrawTriangleLines(v1, v2, v3, color);
```

### DrawTriangleFan

**Signature:**
```c
RLAPI void DrawTriangleFan(const Vector2 *points, int pointCount, Color color)
```

**Koda Binding:**
```koda
native DrawTriangleFan(*points, pointCount, color);
```

### DrawTriangleStrip

**Signature:**
```c
RLAPI void DrawTriangleStrip(const Vector2 *points, int pointCount, Color color)
```

**Koda Binding:**
```koda
native DrawTriangleStrip(*points, pointCount, color);
```

### DrawPoly

**Signature:**
```c
RLAPI void DrawPoly(Vector2 center, int sides, float radius, float rotation, Color color)
```

**Koda Binding:**
```koda
native DrawPoly(center, sides, radius, rotation, color);
```

### DrawPolyLines

**Signature:**
```c
RLAPI void DrawPolyLines(Vector2 center, int sides, float radius, float rotation, Color color)
```

**Koda Binding:**
```koda
native DrawPolyLines(center, sides, radius, rotation, color);
```

### DrawPolyLinesEx

**Signature:**
```c
RLAPI void DrawPolyLinesEx(Vector2 center, int sides, float radius, float rotation, float lineThick, Color color)
```

**Koda Binding:**
```koda
native DrawPolyLinesEx(center, sides, radius, rotation, lineThick, color);
```

### DrawSplineLinear

**Signature:**
```c
RLAPI void DrawSplineLinear(const Vector2 *points, int pointCount, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineLinear(*points, pointCount, thick, color);
```

### DrawSplineBasis

**Signature:**
```c
RLAPI void DrawSplineBasis(const Vector2 *points, int pointCount, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineBasis(*points, pointCount, thick, color);
```

### DrawSplineCatmullRom

**Signature:**
```c
RLAPI void DrawSplineCatmullRom(const Vector2 *points, int pointCount, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineCatmullRom(*points, pointCount, thick, color);
```

### DrawSplineBezierQuadratic

**Signature:**
```c
RLAPI void DrawSplineBezierQuadratic(const Vector2 *points, int pointCount, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineBezierQuadratic(*points, pointCount, thick, color);
```

### DrawSplineBezierCubic

**Signature:**
```c
RLAPI void DrawSplineBezierCubic(const Vector2 *points, int pointCount, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineBezierCubic(*points, pointCount, thick, color);
```

### DrawSplineSegmentLinear

**Signature:**
```c
RLAPI void DrawSplineSegmentLinear(Vector2 p1, Vector2 p2, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineSegmentLinear(p1, p2, thick, color);
```

### DrawSplineSegmentBasis

**Signature:**
```c
RLAPI void DrawSplineSegmentBasis(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineSegmentBasis(p1, p2, p3, p4, thick, color);
```

### DrawSplineSegmentCatmullRom

**Signature:**
```c
RLAPI void DrawSplineSegmentCatmullRom(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineSegmentCatmullRom(p1, p2, p3, p4, thick, color);
```

### DrawSplineSegmentBezierQuadratic

**Signature:**
```c
RLAPI void DrawSplineSegmentBezierQuadratic(Vector2 p1, Vector2 c2, Vector2 p3, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineSegmentBezierQuadratic(p1, c2, p3, thick, color);
```

### DrawSplineSegmentBezierCubic

**Signature:**
```c
RLAPI void DrawSplineSegmentBezierCubic(Vector2 p1, Vector2 c2, Vector2 c3, Vector2 p4, float thick, Color color)
```

**Koda Binding:**
```koda
native DrawSplineSegmentBezierCubic(p1, c2, c3, p4, thick, color);
```

### GetSplinePointLinear

**Signature:**
```c
RLAPI Vector2 GetSplinePointLinear(Vector2 startPos, Vector2 endPos, float t)
```

**Koda Binding:**
```koda
native GetSplinePointLinear(startPos, endPos, t);
```

### GetSplinePointBasis

**Signature:**
```c
RLAPI Vector2 GetSplinePointBasis(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float t)
```

**Koda Binding:**
```koda
native GetSplinePointBasis(p1, p2, p3, p4, t);
```

### GetSplinePointCatmullRom

**Signature:**
```c
RLAPI Vector2 GetSplinePointCatmullRom(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float t)
```

**Koda Binding:**
```koda
native GetSplinePointCatmullRom(p1, p2, p3, p4, t);
```

### GetSplinePointBezierQuadratic

**Signature:**
```c
RLAPI Vector2 GetSplinePointBezierQuadratic(Vector2 p1, Vector2 c2, Vector2 p3, float t)
```

**Koda Binding:**
```koda
native GetSplinePointBezierQuadratic(p1, c2, p3, t);
```

### GetSplinePointBezierCubic

**Signature:**
```c
RLAPI Vector2 GetSplinePointBezierCubic(Vector2 p1, Vector2 c2, Vector2 c3, Vector2 p4, float t)
```

**Koda Binding:**
```koda
native GetSplinePointBezierCubic(p1, c2, c3, p4, t);
```

### CheckCollisionRecs

**Signature:**
```c
RLAPI bool CheckCollisionRecs(Rectangle rec1, Rectangle rec2)
```

**Koda Binding:**
```koda
native CheckCollisionRecs(rec1, rec2);
```

### CheckCollisionCircles

**Signature:**
```c
RLAPI bool CheckCollisionCircles(Vector2 center1, float radius1, Vector2 center2, float radius2)
```

**Koda Binding:**
```koda
native CheckCollisionCircles(center1, radius1, center2, radius2);
```

### CheckCollisionCircleRec

**Signature:**
```c
RLAPI bool CheckCollisionCircleRec(Vector2 center, float radius, Rectangle rec)
```

**Koda Binding:**
```koda
native CheckCollisionCircleRec(center, radius, rec);
```

### CheckCollisionCircleLine

**Signature:**
```c
RLAPI bool CheckCollisionCircleLine(Vector2 center, float radius, Vector2 p1, Vector2 p2)
```

**Koda Binding:**
```koda
native CheckCollisionCircleLine(center, radius, p1, p2);
```

### CheckCollisionPointRec

**Signature:**
```c
RLAPI bool CheckCollisionPointRec(Vector2 point, Rectangle rec)
```

**Koda Binding:**
```koda
native CheckCollisionPointRec(point, rec);
```

### CheckCollisionPointCircle

**Signature:**
```c
RLAPI bool CheckCollisionPointCircle(Vector2 point, Vector2 center, float radius)
```

**Koda Binding:**
```koda
native CheckCollisionPointCircle(point, center, radius);
```

### CheckCollisionPointTriangle

**Signature:**
```c
RLAPI bool CheckCollisionPointTriangle(Vector2 point, Vector2 p1, Vector2 p2, Vector2 p3)
```

**Koda Binding:**
```koda
native CheckCollisionPointTriangle(point, p1, p2, p3);
```

### CheckCollisionPointLine

**Signature:**
```c
RLAPI bool CheckCollisionPointLine(Vector2 point, Vector2 p1, Vector2 p2, int threshold)
```

**Koda Binding:**
```koda
native CheckCollisionPointLine(point, p1, p2, threshold);
```

### CheckCollisionPointPoly

**Signature:**
```c
RLAPI bool CheckCollisionPointPoly(Vector2 point, const Vector2 *points, int pointCount)
```

**Koda Binding:**
```koda
native CheckCollisionPointPoly(point, *points, pointCount);
```

### CheckCollisionLines

**Signature:**
```c
RLAPI bool CheckCollisionLines(Vector2 startPos1, Vector2 endPos1, Vector2 startPos2, Vector2 endPos2, Vector2 *collisionPoint)
```

**Koda Binding:**
```koda
native CheckCollisionLines(startPos1, endPos1, startPos2, endPos2, *collisionPoint);
```

### GetCollisionRec

**Signature:**
```c
RLAPI Rectangle GetCollisionRec(Rectangle rec1, Rectangle rec2)
```

**Koda Binding:**
```koda
native GetCollisionRec(rec1, rec2);
```

### LoadImage

**Signature:**
```c
RLAPI Image LoadImage(const char *fileName)
```

**Koda Binding:**
```koda
native LoadImage(*fileName);
```

### LoadImageRaw

**Signature:**
```c
RLAPI Image LoadImageRaw(const char *fileName, int width, int height, int format, int headerSize)
```

**Koda Binding:**
```koda
native LoadImageRaw(*fileName, width, height, format, headerSize);
```

### LoadImageAnim

**Signature:**
```c
RLAPI Image LoadImageAnim(const char *fileName, int *frames)
```

**Koda Binding:**
```koda
native LoadImageAnim(*fileName, *frames);
```

### LoadImageAnimFromMemory

**Signature:**
```c
RLAPI Image LoadImageAnimFromMemory(const char *fileType, const unsigned char *fileData, int dataSize, int *frames)
```

**Koda Binding:**
```koda
native LoadImageAnimFromMemory(*fileType, *fileData, dataSize, *frames);
```

### LoadImageFromMemory

**Signature:**
```c
RLAPI Image LoadImageFromMemory(const char *fileType, const unsigned char *fileData, int dataSize)
```

**Koda Binding:**
```koda
native LoadImageFromMemory(*fileType, *fileData, dataSize);
```

### LoadImageFromTexture

**Signature:**
```c
RLAPI Image LoadImageFromTexture(Texture2D texture)
```

**Koda Binding:**
```koda
native LoadImageFromTexture(texture);
```

### LoadImageFromScreen

**Signature:**
```c
RLAPI Image LoadImageFromScreen()
```

**Koda Binding:**
```koda
native LoadImageFromScreen();
```

### IsImageValid

**Signature:**
```c
RLAPI bool IsImageValid(Image image)
```

**Koda Binding:**
```koda
native IsImageValid(image);
```

### UnloadImage

**Signature:**
```c
RLAPI void UnloadImage(Image image)
```

**Koda Binding:**
```koda
native UnloadImage(image);
```

### ExportImage

**Signature:**
```c
RLAPI bool ExportImage(Image image, const char *fileName)
```

**Koda Binding:**
```koda
native ExportImage(image, *fileName);
```

### ExportImageAsCode

**Signature:**
```c
RLAPI bool ExportImageAsCode(Image image, const char *fileName)
```

**Koda Binding:**
```koda
native ExportImageAsCode(image, *fileName);
```

### GenImageColor

**Signature:**
```c
RLAPI Image GenImageColor(int width, int height, Color color)
```

**Koda Binding:**
```koda
native GenImageColor(width, height, color);
```

### GenImageGradientLinear

**Signature:**
```c
RLAPI Image GenImageGradientLinear(int width, int height, int direction, Color start, Color end)
```

**Koda Binding:**
```koda
native GenImageGradientLinear(width, height, direction, start, end);
```

### GenImageGradientRadial

**Signature:**
```c
RLAPI Image GenImageGradientRadial(int width, int height, float density, Color inner, Color outer)
```

**Koda Binding:**
```koda
native GenImageGradientRadial(width, height, density, inner, outer);
```

### GenImageGradientSquare

**Signature:**
```c
RLAPI Image GenImageGradientSquare(int width, int height, float density, Color inner, Color outer)
```

**Koda Binding:**
```koda
native GenImageGradientSquare(width, height, density, inner, outer);
```

### GenImageChecked

**Signature:**
```c
RLAPI Image GenImageChecked(int width, int height, int checksX, int checksY, Color col1, Color col2)
```

**Koda Binding:**
```koda
native GenImageChecked(width, height, checksX, checksY, col1, col2);
```

### GenImageWhiteNoise

**Signature:**
```c
RLAPI Image GenImageWhiteNoise(int width, int height, float factor)
```

**Koda Binding:**
```koda
native GenImageWhiteNoise(width, height, factor);
```

### GenImagePerlinNoise

**Signature:**
```c
RLAPI Image GenImagePerlinNoise(int width, int height, int offsetX, int offsetY, float scale)
```

**Koda Binding:**
```koda
native GenImagePerlinNoise(width, height, offsetX, offsetY, scale);
```

### GenImageCellular

**Signature:**
```c
RLAPI Image GenImageCellular(int width, int height, int tileSize)
```

**Koda Binding:**
```koda
native GenImageCellular(width, height, tileSize);
```

### GenImageText

**Signature:**
```c
RLAPI Image GenImageText(int width, int height, const char *text)
```

**Koda Binding:**
```koda
native GenImageText(width, height, *text);
```

### ImageCopy

**Signature:**
```c
RLAPI Image ImageCopy(Image image)
```

**Koda Binding:**
```koda
native ImageCopy(image);
```

### ImageFromImage

**Signature:**
```c
RLAPI Image ImageFromImage(Image image, Rectangle rec)
```

**Koda Binding:**
```koda
native ImageFromImage(image, rec);
```

### ImageFromChannel

**Signature:**
```c
RLAPI Image ImageFromChannel(Image image, int selectedChannel)
```

**Koda Binding:**
```koda
native ImageFromChannel(image, selectedChannel);
```

### ImageText

**Signature:**
```c
RLAPI Image ImageText(const char *text, int fontSize, Color color)
```

**Koda Binding:**
```koda
native ImageText(*text, fontSize, color);
```

### ImageTextEx

**Signature:**
```c
RLAPI Image ImageTextEx(Font font, const char *text, float fontSize, float spacing, Color tint)
```

**Koda Binding:**
```koda
native ImageTextEx(font, *text, fontSize, spacing, tint);
```

### ImageFormat

**Signature:**
```c
RLAPI void ImageFormat(Image *image, int newFormat)
```

**Koda Binding:**
```koda
native ImageFormat(*image, newFormat);
```

### ImageToPOT

**Signature:**
```c
RLAPI void ImageToPOT(Image *image, Color fill)
```

**Koda Binding:**
```koda
native ImageToPOT(*image, fill);
```

### ImageCrop

**Signature:**
```c
RLAPI void ImageCrop(Image *image, Rectangle crop)
```

**Koda Binding:**
```koda
native ImageCrop(*image, crop);
```

### ImageAlphaCrop

**Signature:**
```c
RLAPI void ImageAlphaCrop(Image *image, float threshold)
```

**Koda Binding:**
```koda
native ImageAlphaCrop(*image, threshold);
```

### ImageAlphaClear

**Signature:**
```c
RLAPI void ImageAlphaClear(Image *image, Color color, float threshold)
```

**Koda Binding:**
```koda
native ImageAlphaClear(*image, color, threshold);
```

### ImageAlphaMask

**Signature:**
```c
RLAPI void ImageAlphaMask(Image *image, Image alphaMask)
```

**Koda Binding:**
```koda
native ImageAlphaMask(*image, alphaMask);
```

### ImageAlphaPremultiply

**Signature:**
```c
RLAPI void ImageAlphaPremultiply(Image *image)
```

**Koda Binding:**
```koda
native ImageAlphaPremultiply(*image);
```

### ImageBlurGaussian

**Signature:**
```c
RLAPI void ImageBlurGaussian(Image *image, int blurSize)
```

**Koda Binding:**
```koda
native ImageBlurGaussian(*image, blurSize);
```

### ImageKernelConvolution

**Signature:**
```c
RLAPI void ImageKernelConvolution(Image *image, const float *kernel, int kernelSize)
```

**Koda Binding:**
```koda
native ImageKernelConvolution(*image, *kernel, kernelSize);
```

### ImageResize

**Signature:**
```c
RLAPI void ImageResize(Image *image, int newWidth, int newHeight)
```

**Koda Binding:**
```koda
native ImageResize(*image, newWidth, newHeight);
```

### ImageResizeNN

**Signature:**
```c
RLAPI void ImageResizeNN(Image *image, int newWidth, int newHeight)
```

**Koda Binding:**
```koda
native ImageResizeNN(*image, newWidth, newHeight);
```

### ImageResizeCanvas

**Signature:**
```c
RLAPI void ImageResizeCanvas(Image *image, int newWidth, int newHeight, int offsetX, int offsetY, Color fill)
```

**Koda Binding:**
```koda
native ImageResizeCanvas(*image, newWidth, newHeight, offsetX, offsetY, fill);
```

### ImageMipmaps

**Signature:**
```c
RLAPI void ImageMipmaps(Image *image)
```

**Koda Binding:**
```koda
native ImageMipmaps(*image);
```

### ImageDither

**Signature:**
```c
RLAPI void ImageDither(Image *image, int rBpp, int gBpp, int bBpp, int aBpp)
```

**Koda Binding:**
```koda
native ImageDither(*image, rBpp, gBpp, bBpp, aBpp);
```

### ImageFlipVertical

**Signature:**
```c
RLAPI void ImageFlipVertical(Image *image)
```

**Koda Binding:**
```koda
native ImageFlipVertical(*image);
```

### ImageFlipHorizontal

**Signature:**
```c
RLAPI void ImageFlipHorizontal(Image *image)
```

**Koda Binding:**
```koda
native ImageFlipHorizontal(*image);
```

### ImageRotate

**Signature:**
```c
RLAPI void ImageRotate(Image *image, int degrees)
```

**Koda Binding:**
```koda
native ImageRotate(*image, degrees);
```

### ImageRotateCW

**Signature:**
```c
RLAPI void ImageRotateCW(Image *image)
```

**Koda Binding:**
```koda
native ImageRotateCW(*image);
```

### ImageRotateCCW

**Signature:**
```c
RLAPI void ImageRotateCCW(Image *image)
```

**Koda Binding:**
```koda
native ImageRotateCCW(*image);
```

### ImageColorTint

**Signature:**
```c
RLAPI void ImageColorTint(Image *image, Color color)
```

**Koda Binding:**
```koda
native ImageColorTint(*image, color);
```

### ImageColorInvert

**Signature:**
```c
RLAPI void ImageColorInvert(Image *image)
```

**Koda Binding:**
```koda
native ImageColorInvert(*image);
```

### ImageColorGrayscale

**Signature:**
```c
RLAPI void ImageColorGrayscale(Image *image)
```

**Koda Binding:**
```koda
native ImageColorGrayscale(*image);
```

### ImageColorContrast

**Signature:**
```c
RLAPI void ImageColorContrast(Image *image, int contrast)
```

**Koda Binding:**
```koda
native ImageColorContrast(*image, contrast);
```

### ImageColorBrightness

**Signature:**
```c
RLAPI void ImageColorBrightness(Image *image, int brightness)
```

**Koda Binding:**
```koda
native ImageColorBrightness(*image, brightness);
```

### ImageColorReplace

**Signature:**
```c
RLAPI void ImageColorReplace(Image *image, Color color, Color replace)
```

**Koda Binding:**
```koda
native ImageColorReplace(*image, color, replace);
```

### UnloadImageColors

**Signature:**
```c
RLAPI void UnloadImageColors(Color *colors)
```

**Koda Binding:**
```koda
native UnloadImageColors(*colors);
```

### UnloadImagePalette

**Signature:**
```c
RLAPI void UnloadImagePalette(Color *colors)
```

**Koda Binding:**
```koda
native UnloadImagePalette(*colors);
```

### GetImageAlphaBorder

**Signature:**
```c
RLAPI Rectangle GetImageAlphaBorder(Image image, float threshold)
```

**Koda Binding:**
```koda
native GetImageAlphaBorder(image, threshold);
```

### GetImageColor

**Signature:**
```c
RLAPI Color GetImageColor(Image image, int x, int y)
```

**Koda Binding:**
```koda
native GetImageColor(image, x, y);
```

### ImageClearBackground

**Signature:**
```c
RLAPI void ImageClearBackground(Image *dst, Color color)
```

**Koda Binding:**
```koda
native ImageClearBackground(*dst, color);
```

### ImageDrawPixel

**Signature:**
```c
RLAPI void ImageDrawPixel(Image *dst, int posX, int posY, Color color)
```

**Koda Binding:**
```koda
native ImageDrawPixel(*dst, posX, posY, color);
```

### ImageDrawPixelV

**Signature:**
```c
RLAPI void ImageDrawPixelV(Image *dst, Vector2 position, Color color)
```

**Koda Binding:**
```koda
native ImageDrawPixelV(*dst, position, color);
```

### ImageDrawLine

**Signature:**
```c
RLAPI void ImageDrawLine(Image *dst, int startPosX, int startPosY, int endPosX, int endPosY, Color color)
```

**Koda Binding:**
```koda
native ImageDrawLine(*dst, startPosX, startPosY, endPosX, endPosY, color);
```

### ImageDrawLineV

**Signature:**
```c
RLAPI void ImageDrawLineV(Image *dst, Vector2 start, Vector2 end, Color color)
```

**Koda Binding:**
```koda
native ImageDrawLineV(*dst, start, end, color);
```

### ImageDrawLineEx

**Signature:**
```c
RLAPI void ImageDrawLineEx(Image *dst, Vector2 start, Vector2 end, int thick, Color color)
```

**Koda Binding:**
```koda
native ImageDrawLineEx(*dst, start, end, thick, color);
```

### ImageDrawCircle

**Signature:**
```c
RLAPI void ImageDrawCircle(Image *dst, int centerX, int centerY, int radius, Color color)
```

**Koda Binding:**
```koda
native ImageDrawCircle(*dst, centerX, centerY, radius, color);
```

### ImageDrawCircleV

**Signature:**
```c
RLAPI void ImageDrawCircleV(Image *dst, Vector2 center, int radius, Color color)
```

**Koda Binding:**
```koda
native ImageDrawCircleV(*dst, center, radius, color);
```

### ImageDrawCircleLines

**Signature:**
```c
RLAPI void ImageDrawCircleLines(Image *dst, int centerX, int centerY, int radius, Color color)
```

**Koda Binding:**
```koda
native ImageDrawCircleLines(*dst, centerX, centerY, radius, color);
```

### ImageDrawCircleLinesV

**Signature:**
```c
RLAPI void ImageDrawCircleLinesV(Image *dst, Vector2 center, int radius, Color color)
```

**Koda Binding:**
```koda
native ImageDrawCircleLinesV(*dst, center, radius, color);
```

### ImageDrawRectangle

**Signature:**
```c
RLAPI void ImageDrawRectangle(Image *dst, int posX, int posY, int width, int height, Color color)
```

**Koda Binding:**
```koda
native ImageDrawRectangle(*dst, posX, posY, width, height, color);
```

### ImageDrawRectangleV

**Signature:**
```c
RLAPI void ImageDrawRectangleV(Image *dst, Vector2 position, Vector2 size, Color color)
```

**Koda Binding:**
```koda
native ImageDrawRectangleV(*dst, position, size, color);
```

### ImageDrawRectangleRec

**Signature:**
```c
RLAPI void ImageDrawRectangleRec(Image *dst, Rectangle rec, Color color)
```

**Koda Binding:**
```koda
native ImageDrawRectangleRec(*dst, rec, color);
```

### ImageDrawRectangleLines

**Signature:**
```c
RLAPI void ImageDrawRectangleLines(Image *dst, int posX, int posY, int width, int height, Color color)
```

**Koda Binding:**
```koda
native ImageDrawRectangleLines(*dst, posX, posY, width, height, color);
```

### ImageDrawRectangleLinesEx

**Signature:**
```c
RLAPI void ImageDrawRectangleLinesEx(Image *dst, Rectangle rec, int thick, Color color)
```

**Koda Binding:**
```koda
native ImageDrawRectangleLinesEx(*dst, rec, thick, color);
```

### ImageDrawTriangle

**Signature:**
```c
RLAPI void ImageDrawTriangle(Image *dst, Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda Binding:**
```koda
native ImageDrawTriangle(*dst, v1, v2, v3, color);
```

### ImageDrawTriangleGradient

**Signature:**
```c
RLAPI void ImageDrawTriangleGradient(Image *dst, Vector2 v1, Vector2 v2, Vector2 v3, Color c1, Color c2, Color c3)
```

**Koda Binding:**
```koda
native ImageDrawTriangleGradient(*dst, v1, v2, v3, c1, c2, c3);
```

### ImageDrawTriangleLines

**Signature:**
```c
RLAPI void ImageDrawTriangleLines(Image *dst, Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Koda Binding:**
```koda
native ImageDrawTriangleLines(*dst, v1, v2, v3, color);
```

### ImageDrawTriangleFan

**Signature:**
```c
RLAPI void ImageDrawTriangleFan(Image *dst, const Vector2 *points, int pointCount, Color color)
```

**Koda Binding:**
```koda
native ImageDrawTriangleFan(*dst, *points, pointCount, color);
```

### ImageDrawTriangleStrip

**Signature:**
```c
RLAPI void ImageDrawTriangleStrip(Image *dst, const Vector2 *points, int pointCount, Color color)
```

**Koda Binding:**
```koda
native ImageDrawTriangleStrip(*dst, *points, pointCount, color);
```

### ImageDraw

**Signature:**
```c
RLAPI void ImageDraw(Image *dst, Image src, Rectangle srcRec, Rectangle dstRec, Color tint)
```

**Koda Binding:**
```koda
native ImageDraw(*dst, src, srcRec, dstRec, tint);
```

### ImageDrawText

**Signature:**
```c
RLAPI void ImageDrawText(Image *dst, const char *text, int posX, int posY, int fontSize, Color color)
```

**Koda Binding:**
```koda
native ImageDrawText(*dst, *text, posX, posY, fontSize, color);
```

### ImageDrawTextEx

**Signature:**
```c
RLAPI void ImageDrawTextEx(Image *dst, Font font, const char *text, Vector2 position, float fontSize, float spacing, Color tint)
```

**Koda Binding:**
```koda
native ImageDrawTextEx(*dst, font, *text, position, fontSize, spacing, tint);
```

### LoadTexture

**Signature:**
```c
RLAPI Texture2D LoadTexture(const char *fileName)
```

**Koda Binding:**
```koda
native LoadTexture(*fileName);
```

### LoadTextureFromImage

**Signature:**
```c
RLAPI Texture2D LoadTextureFromImage(Image image)
```

**Koda Binding:**
```koda
native LoadTextureFromImage(image);
```

### LoadTextureCubemap

**Signature:**
```c
RLAPI TextureCubemap LoadTextureCubemap(Image image, int layout)
```

**Koda Binding:**
```koda
native LoadTextureCubemap(image, layout);
```

### LoadRenderTexture

**Signature:**
```c
RLAPI RenderTexture2D LoadRenderTexture(int width, int height)
```

**Koda Binding:**
```koda
native LoadRenderTexture(width, height);
```

### IsTextureValid

**Signature:**
```c
RLAPI bool IsTextureValid(Texture2D texture)
```

**Koda Binding:**
```koda
native IsTextureValid(texture);
```

### UnloadTexture

**Signature:**
```c
RLAPI void UnloadTexture(Texture2D texture)
```

**Koda Binding:**
```koda
native UnloadTexture(texture);
```

### IsRenderTextureValid

**Signature:**
```c
RLAPI bool IsRenderTextureValid(RenderTexture2D target)
```

**Koda Binding:**
```koda
native IsRenderTextureValid(target);
```

### UnloadRenderTexture

**Signature:**
```c
RLAPI void UnloadRenderTexture(RenderTexture2D target)
```

**Koda Binding:**
```koda
native UnloadRenderTexture(target);
```

### UpdateTexture

**Signature:**
```c
RLAPI void UpdateTexture(Texture2D texture, const void *pixels)
```

**Koda Binding:**
```koda
native UpdateTexture(texture, *pixels);
```

### UpdateTextureRec

**Signature:**
```c
RLAPI void UpdateTextureRec(Texture2D texture, Rectangle rec, const void *pixels)
```

**Koda Binding:**
```koda
native UpdateTextureRec(texture, rec, *pixels);
```

### GenTextureMipmaps

**Signature:**
```c
RLAPI void GenTextureMipmaps(Texture2D *texture)
```

**Koda Binding:**
```koda
native GenTextureMipmaps(*texture);
```

### SetTextureFilter

**Signature:**
```c
RLAPI void SetTextureFilter(Texture2D texture, int filter)
```

**Koda Binding:**
```koda
native SetTextureFilter(texture, filter);
```

### SetTextureWrap

**Signature:**
```c
RLAPI void SetTextureWrap(Texture2D texture, int wrap)
```

**Koda Binding:**
```koda
native SetTextureWrap(texture, wrap);
```

### DrawTexture

**Signature:**
```c
RLAPI void DrawTexture(Texture2D texture, int posX, int posY, Color tint)
```

**Koda Binding:**
```koda
native DrawTexture(texture, posX, posY, tint);
```

### DrawTextureV

**Signature:**
```c
RLAPI void DrawTextureV(Texture2D texture, Vector2 position, Color tint)
```

**Koda Binding:**
```koda
native DrawTextureV(texture, position, tint);
```

### DrawTextureEx

**Signature:**
```c
RLAPI void DrawTextureEx(Texture2D texture, Vector2 position, float rotation, float scale, Color tint)
```

**Koda Binding:**
```koda
native DrawTextureEx(texture, position, rotation, scale, tint);
```

### DrawTextureRec

**Signature:**
```c
RLAPI void DrawTextureRec(Texture2D texture, Rectangle source, Vector2 position, Color tint)
```

**Koda Binding:**
```koda
native DrawTextureRec(texture, source, position, tint);
```

### DrawTexturePro

**Signature:**
```c
RLAPI void DrawTexturePro(Texture2D texture, Rectangle source, Rectangle dest, Vector2 origin, float rotation, Color tint)
```

**Koda Binding:**
```koda
native DrawTexturePro(texture, source, dest, origin, rotation, tint);
```

### DrawTextureNPatch

**Signature:**
```c
RLAPI void DrawTextureNPatch(Texture2D texture, NPatchInfo nPatchInfo, Rectangle dest, Vector2 origin, float rotation, Color tint)
```

**Koda Binding:**
```koda
native DrawTextureNPatch(texture, nPatchInfo, dest, origin, rotation, tint);
```

### ColorIsEqual

**Signature:**
```c
RLAPI bool ColorIsEqual(Color col1, Color col2)
```

**Koda Binding:**
```koda
native ColorIsEqual(col1, col2);
```

### Fade

**Signature:**
```c
RLAPI Color Fade(Color color, float alpha)
```

**Koda Binding:**
```koda
native Fade(color, alpha);
```

### ColorToInt

**Signature:**
```c
RLAPI int ColorToInt(Color color)
```

**Koda Binding:**
```koda
native ColorToInt(color);
```

### ColorNormalize

**Signature:**
```c
RLAPI Vector4 ColorNormalize(Color color)
```

**Koda Binding:**
```koda
native ColorNormalize(color);
```

### ColorFromNormalized

**Signature:**
```c
RLAPI Color ColorFromNormalized(Vector4 normalized)
```

**Koda Binding:**
```koda
native ColorFromNormalized(normalized);
```

### ColorToHSV

**Signature:**
```c
RLAPI Vector3 ColorToHSV(Color color)
```

**Koda Binding:**
```koda
native ColorToHSV(color);
```

### ColorFromHSV

**Signature:**
```c
RLAPI Color ColorFromHSV(float hue, float saturation, float value)
```

**Koda Binding:**
```koda
native ColorFromHSV(hue, saturation, value);
```

### ColorTint

**Signature:**
```c
RLAPI Color ColorTint(Color color, Color tint)
```

**Koda Binding:**
```koda
native ColorTint(color, tint);
```

### ColorBrightness

**Signature:**
```c
RLAPI Color ColorBrightness(Color color, float factor)
```

**Koda Binding:**
```koda
native ColorBrightness(color, factor);
```

### ColorContrast

**Signature:**
```c
RLAPI Color ColorContrast(Color color, float contrast)
```

**Koda Binding:**
```koda
native ColorContrast(color, contrast);
```

### ColorAlpha

**Signature:**
```c
RLAPI Color ColorAlpha(Color color, float alpha)
```

**Koda Binding:**
```koda
native ColorAlpha(color, alpha);
```

### ColorAlphaBlend

**Signature:**
```c
RLAPI Color ColorAlphaBlend(Color dst, Color src, Color tint)
```

**Koda Binding:**
```koda
native ColorAlphaBlend(dst, src, tint);
```

### ColorLerp

**Signature:**
```c
RLAPI Color ColorLerp(Color color1, Color color2, float factor)
```

**Koda Binding:**
```koda
native ColorLerp(color1, color2, factor);
```

### GetColor

**Signature:**
```c
RLAPI Color GetColor(unsigned int hexValue)
```

**Koda Binding:**
```koda
native GetColor(hexValue);
```

### GetPixelColor

**Signature:**
```c
RLAPI Color GetPixelColor(void *srcPtr, int format)
```

**Koda Binding:**
```koda
native GetPixelColor(*srcPtr, format);
```

### SetPixelColor

**Signature:**
```c
RLAPI void SetPixelColor(void *dstPtr, Color color, int format)
```

**Koda Binding:**
```koda
native SetPixelColor(*dstPtr, color, format);
```

### GetPixelDataSize

**Signature:**
```c
RLAPI int GetPixelDataSize(int width, int height, int format)
```

**Koda Binding:**
```koda
native GetPixelDataSize(width, height, format);
```

### GetFontDefault

**Signature:**
```c
RLAPI Font GetFontDefault()
```

**Koda Binding:**
```koda
native GetFontDefault();
```

### LoadFont

**Signature:**
```c
RLAPI Font LoadFont(const char *fileName)
```

**Koda Binding:**
```koda
native LoadFont(*fileName);
```

### LoadFontEx

**Signature:**
```c
RLAPI Font LoadFontEx(const char *fileName, int fontSize, const int *codepoints, int codepointCount)
```

**Koda Binding:**
```koda
native LoadFontEx(*fileName, fontSize, *codepoints, codepointCount);
```

### LoadFontFromImage

**Signature:**
```c
RLAPI Font LoadFontFromImage(Image image, Color key, int firstChar)
```

**Koda Binding:**
```koda
native LoadFontFromImage(image, key, firstChar);
```

### LoadFontFromMemory

**Signature:**
```c
RLAPI Font LoadFontFromMemory(const char *fileType, const unsigned char *fileData, int dataSize, int fontSize, const int *codepoints, int codepointCount)
```

**Koda Binding:**
```koda
native LoadFontFromMemory(*fileType, *fileData, dataSize, fontSize, *codepoints, codepointCount);
```

### IsFontValid

**Signature:**
```c
RLAPI bool IsFontValid(Font font)
```

**Koda Binding:**
```koda
native IsFontValid(font);
```

### GenImageFontAtlas

**Signature:**
```c
RLAPI Image GenImageFontAtlas(const GlyphInfo *glyphs, Rectangle **glyphRecs, int glyphCount, int fontSize, int padding, int packMethod)
```

**Koda Binding:**
```koda
native GenImageFontAtlas(*glyphs, **glyphRecs, glyphCount, fontSize, padding, packMethod);
```

### UnloadFontData

**Signature:**
```c
RLAPI void UnloadFontData(GlyphInfo *glyphs, int glyphCount)
```

**Koda Binding:**
```koda
native UnloadFontData(*glyphs, glyphCount);
```

### UnloadFont

**Signature:**
```c
RLAPI void UnloadFont(Font font)
```

**Koda Binding:**
```koda
native UnloadFont(font);
```

### ExportFontAsCode

**Signature:**
```c
RLAPI bool ExportFontAsCode(Font font, const char *fileName)
```

**Koda Binding:**
```koda
native ExportFontAsCode(font, *fileName);
```

### DrawFPS

**Signature:**
```c
RLAPI void DrawFPS(int posX, int posY)
```

**Koda Binding:**
```koda
native DrawFPS(posX, posY);
```

### DrawText

**Signature:**
```c
RLAPI void DrawText(const char *text, int posX, int posY, int fontSize, Color color)
```

**Koda Binding:**
```koda
native DrawText(*text, posX, posY, fontSize, color);
```

### DrawTextEx

**Signature:**
```c
RLAPI void DrawTextEx(Font font, const char *text, Vector2 position, float fontSize, float spacing, Color tint)
```

**Koda Binding:**
```koda
native DrawTextEx(font, *text, position, fontSize, spacing, tint);
```

### DrawTextPro

**Signature:**
```c
RLAPI void DrawTextPro(Font font, const char *text, Vector2 position, Vector2 origin, float rotation, float fontSize, float spacing, Color tint)
```

**Koda Binding:**
```koda
native DrawTextPro(font, *text, position, origin, rotation, fontSize, spacing, tint);
```

### DrawTextCodepoint

**Signature:**
```c
RLAPI void DrawTextCodepoint(Font font, int codepoint, Vector2 position, float fontSize, Color tint)
```

**Koda Binding:**
```koda
native DrawTextCodepoint(font, codepoint, position, fontSize, tint);
```

### DrawTextCodepoints

**Signature:**
```c
RLAPI void DrawTextCodepoints(Font font, const int *codepoints, int codepointCount, Vector2 position, float fontSize, float spacing, Color tint)
```

**Koda Binding:**
```koda
native DrawTextCodepoints(font, *codepoints, codepointCount, position, fontSize, spacing, tint);
```

### SetTextLineSpacing

**Signature:**
```c
RLAPI void SetTextLineSpacing(int spacing)
```

**Koda Binding:**
```koda
native SetTextLineSpacing(spacing);
```

### MeasureText

**Signature:**
```c
RLAPI int MeasureText(const char *text, int fontSize)
```

**Koda Binding:**
```koda
native MeasureText(*text, fontSize);
```

### MeasureTextEx

**Signature:**
```c
RLAPI Vector2 MeasureTextEx(Font font, const char *text, float fontSize, float spacing)
```

**Koda Binding:**
```koda
native MeasureTextEx(font, *text, fontSize, spacing);
```

### MeasureTextCodepoints

**Signature:**
```c
RLAPI Vector2 MeasureTextCodepoints(Font font, const int *codepoints, int length, float fontSize, float spacing)
```

**Koda Binding:**
```koda
native MeasureTextCodepoints(font, *codepoints, length, fontSize, spacing);
```

### GetGlyphIndex

**Signature:**
```c
RLAPI int GetGlyphIndex(Font font, int codepoint)
```

**Koda Binding:**
```koda
native GetGlyphIndex(font, codepoint);
```

### GetGlyphInfo

**Signature:**
```c
RLAPI GlyphInfo GetGlyphInfo(Font font, int codepoint)
```

**Koda Binding:**
```koda
native GetGlyphInfo(font, codepoint);
```

### GetGlyphAtlasRec

**Signature:**
```c
RLAPI Rectangle GetGlyphAtlasRec(Font font, int codepoint)
```

**Koda Binding:**
```koda
native GetGlyphAtlasRec(font, codepoint);
```

### UnloadUTF8

**Signature:**
```c
RLAPI void UnloadUTF8(char *text)
```

**Koda Binding:**
```koda
native UnloadUTF8(*text);
```

### UnloadCodepoints

**Signature:**
```c
RLAPI void UnloadCodepoints(int *codepoints)
```

**Koda Binding:**
```koda
native UnloadCodepoints(*codepoints);
```

### GetCodepointCount

**Signature:**
```c
RLAPI int GetCodepointCount(const char *text)
```

**Koda Binding:**
```koda
native GetCodepointCount(*text);
```

### GetCodepoint

**Signature:**
```c
RLAPI int GetCodepoint(const char *text, int *codepointSize)
```

**Koda Binding:**
```koda
native GetCodepoint(*text, *codepointSize);
```

### GetCodepointNext

**Signature:**
```c
RLAPI int GetCodepointNext(const char *text, int *codepointSize)
```

**Koda Binding:**
```koda
native GetCodepointNext(*text, *codepointSize);
```

### GetCodepointPrevious

**Signature:**
```c
RLAPI int GetCodepointPrevious(const char *text, int *codepointSize)
```

**Koda Binding:**
```koda
native GetCodepointPrevious(*text, *codepointSize);
```

### UnloadTextLines

**Signature:**
```c
RLAPI void UnloadTextLines(char **text, int lineCount)
```

**Koda Binding:**
```koda
native UnloadTextLines(**text, lineCount);
```

### TextCopy

**Signature:**
```c
RLAPI int TextCopy(char *dst, const char *src)
```

**Koda Binding:**
```koda
native TextCopy(*dst, *src);
```

### TextIsEqual

**Signature:**
```c
RLAPI bool TextIsEqual(const char *text1, const char *text2)
```

**Koda Binding:**
```koda
native TextIsEqual(*text1, *text2);
```

### TextLength

**Signature:**
```c
RLAPI unsigned int TextLength(const char *text)
```

**Koda Binding:**
```koda
native TextLength(*text);
```

### TextAppend

**Signature:**
```c
RLAPI void TextAppend(char *text, const char *append, int *position)
```

**Koda Binding:**
```koda
native TextAppend(*text, *append, *position);
```

### TextFindIndex

**Signature:**
```c
RLAPI int TextFindIndex(const char *text, const char *search)
```

**Koda Binding:**
```koda
native TextFindIndex(*text, *search);
```

### TextToInteger

**Signature:**
```c
RLAPI int TextToInteger(const char *text)
```

**Koda Binding:**
```koda
native TextToInteger(*text);
```

### TextToFloat

**Signature:**
```c
RLAPI float TextToFloat(const char *text)
```

**Koda Binding:**
```koda
native TextToFloat(*text);
```

### DrawLine3D

**Signature:**
```c
RLAPI void DrawLine3D(Vector3 startPos, Vector3 endPos, Color color)
```

**Koda Binding:**
```koda
native DrawLine3D(startPos, endPos, color);
```

### DrawPoint3D

**Signature:**
```c
RLAPI void DrawPoint3D(Vector3 position, Color color)
```

**Koda Binding:**
```koda
native DrawPoint3D(position, color);
```

### DrawCircle3D

**Signature:**
```c
RLAPI void DrawCircle3D(Vector3 center, float radius, Vector3 rotationAxis, float rotationAngle, Color color)
```

**Koda Binding:**
```koda
native DrawCircle3D(center, radius, rotationAxis, rotationAngle, color);
```

### DrawTriangle3D

**Signature:**
```c
RLAPI void DrawTriangle3D(Vector3 v1, Vector3 v2, Vector3 v3, Color color)
```

**Koda Binding:**
```koda
native DrawTriangle3D(v1, v2, v3, color);
```

### DrawTriangleStrip3D

**Signature:**
```c
RLAPI void DrawTriangleStrip3D(const Vector3 *points, int pointCount, Color color)
```

**Koda Binding:**
```koda
native DrawTriangleStrip3D(*points, pointCount, color);
```

### DrawCube

**Signature:**
```c
RLAPI void DrawCube(Vector3 position, float width, float height, float length, Color color)
```

**Koda Binding:**
```koda
native DrawCube(position, width, height, length, color);
```

### DrawCubeV

**Signature:**
```c
RLAPI void DrawCubeV(Vector3 position, Vector3 size, Color color)
```

**Koda Binding:**
```koda
native DrawCubeV(position, size, color);
```

### DrawCubeWires

**Signature:**
```c
RLAPI void DrawCubeWires(Vector3 position, float width, float height, float length, Color color)
```

**Koda Binding:**
```koda
native DrawCubeWires(position, width, height, length, color);
```

### DrawCubeWiresV

**Signature:**
```c
RLAPI void DrawCubeWiresV(Vector3 position, Vector3 size, Color color)
```

**Koda Binding:**
```koda
native DrawCubeWiresV(position, size, color);
```

### DrawSphere

**Signature:**
```c
RLAPI void DrawSphere(Vector3 centerPos, float radius, Color color)
```

**Koda Binding:**
```koda
native DrawSphere(centerPos, radius, color);
```

### DrawSphereEx

**Signature:**
```c
RLAPI void DrawSphereEx(Vector3 centerPos, float radius, int rings, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawSphereEx(centerPos, radius, rings, slices, color);
```

### DrawSphereWires

**Signature:**
```c
RLAPI void DrawSphereWires(Vector3 centerPos, float radius, int rings, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawSphereWires(centerPos, radius, rings, slices, color);
```

### DrawCylinder

**Signature:**
```c
RLAPI void DrawCylinder(Vector3 position, float radiusTop, float radiusBottom, float height, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawCylinder(position, radiusTop, radiusBottom, height, slices, color);
```

### DrawCylinderEx

**Signature:**
```c
RLAPI void DrawCylinderEx(Vector3 startPos, Vector3 endPos, float startRadius, float endRadius, int sides, Color color)
```

**Koda Binding:**
```koda
native DrawCylinderEx(startPos, endPos, startRadius, endRadius, sides, color);
```

### DrawCylinderWires

**Signature:**
```c
RLAPI void DrawCylinderWires(Vector3 position, float radiusTop, float radiusBottom, float height, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawCylinderWires(position, radiusTop, radiusBottom, height, slices, color);
```

### DrawCylinderWiresEx

**Signature:**
```c
RLAPI void DrawCylinderWiresEx(Vector3 startPos, Vector3 endPos, float startRadius, float endRadius, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawCylinderWiresEx(startPos, endPos, startRadius, endRadius, slices, color);
```

### DrawCapsule

**Signature:**
```c
RLAPI void DrawCapsule(Vector3 startPos, Vector3 endPos, float radius, int rings, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawCapsule(startPos, endPos, radius, rings, slices, color);
```

### DrawCapsuleWires

**Signature:**
```c
RLAPI void DrawCapsuleWires(Vector3 startPos, Vector3 endPos, float radius, int rings, int slices, Color color)
```

**Koda Binding:**
```koda
native DrawCapsuleWires(startPos, endPos, radius, rings, slices, color);
```

### DrawPlane

**Signature:**
```c
RLAPI void DrawPlane(Vector3 centerPos, Vector2 size, Color color)
```

**Koda Binding:**
```koda
native DrawPlane(centerPos, size, color);
```

### DrawRay

**Signature:**
```c
RLAPI void DrawRay(Ray ray, Color color)
```

**Koda Binding:**
```koda
native DrawRay(ray, color);
```

### DrawGrid

**Signature:**
```c
RLAPI void DrawGrid(int slices, float spacing)
```

**Koda Binding:**
```koda
native DrawGrid(slices, spacing);
```

### LoadModel

**Signature:**
```c
RLAPI Model LoadModel(const char *fileName)
```

**Koda Binding:**
```koda
native LoadModel(*fileName);
```

### LoadModelFromMesh

**Signature:**
```c
RLAPI Model LoadModelFromMesh(Mesh mesh)
```

**Koda Binding:**
```koda
native LoadModelFromMesh(mesh);
```

### IsModelValid

**Signature:**
```c
RLAPI bool IsModelValid(Model model)
```

**Koda Binding:**
```koda
native IsModelValid(model);
```

### UnloadModel

**Signature:**
```c
RLAPI void UnloadModel(Model model)
```

**Koda Binding:**
```koda
native UnloadModel(model);
```

### GetModelBoundingBox

**Signature:**
```c
RLAPI BoundingBox GetModelBoundingBox(Model model)
```

**Koda Binding:**
```koda
native GetModelBoundingBox(model);
```

### DrawModel

**Signature:**
```c
RLAPI void DrawModel(Model model, Vector3 position, float scale, Color tint)
```

**Koda Binding:**
```koda
native DrawModel(model, position, scale, tint);
```

### DrawModelEx

**Signature:**
```c
RLAPI void DrawModelEx(Model model, Vector3 position, Vector3 rotationAxis, float rotationAngle, Vector3 scale, Color tint)
```

**Koda Binding:**
```koda
native DrawModelEx(model, position, rotationAxis, rotationAngle, scale, tint);
```

### DrawModelWires

**Signature:**
```c
RLAPI void DrawModelWires(Model model, Vector3 position, float scale, Color tint)
```

**Koda Binding:**
```koda
native DrawModelWires(model, position, scale, tint);
```

### DrawModelWiresEx

**Signature:**
```c
RLAPI void DrawModelWiresEx(Model model, Vector3 position, Vector3 rotationAxis, float rotationAngle, Vector3 scale, Color tint)
```

**Koda Binding:**
```koda
native DrawModelWiresEx(model, position, rotationAxis, rotationAngle, scale, tint);
```

### DrawBoundingBox

**Signature:**
```c
RLAPI void DrawBoundingBox(BoundingBox box, Color color)
```

**Koda Binding:**
```koda
native DrawBoundingBox(box, color);
```

### DrawBillboard

**Signature:**
```c
RLAPI void DrawBillboard(Camera camera, Texture2D texture, Vector3 position, float scale, Color tint)
```

**Koda Binding:**
```koda
native DrawBillboard(camera, texture, position, scale, tint);
```

### DrawBillboardRec

**Signature:**
```c
RLAPI void DrawBillboardRec(Camera camera, Texture2D texture, Rectangle source, Vector3 position, Vector2 size, Color tint)
```

**Koda Binding:**
```koda
native DrawBillboardRec(camera, texture, source, position, size, tint);
```

### DrawBillboardPro

**Signature:**
```c
RLAPI void DrawBillboardPro(Camera camera, Texture2D texture, Rectangle source, Vector3 position, Vector3 up, Vector2 size, Vector2 origin, float rotation, Color tint)
```

**Koda Binding:**
```koda
native DrawBillboardPro(camera, texture, source, position, up, size, origin, rotation, tint);
```

### UploadMesh

**Signature:**
```c
RLAPI void UploadMesh(Mesh *mesh, bool dynamic)
```

**Koda Binding:**
```koda
native UploadMesh(*mesh, dynamic);
```

### UpdateMeshBuffer

**Signature:**
```c
RLAPI void UpdateMeshBuffer(Mesh mesh, int index, const void *data, int dataSize, int offset)
```

**Koda Binding:**
```koda
native UpdateMeshBuffer(mesh, index, *data, dataSize, offset);
```

### UnloadMesh

**Signature:**
```c
RLAPI void UnloadMesh(Mesh mesh)
```

**Koda Binding:**
```koda
native UnloadMesh(mesh);
```

### DrawMesh

**Signature:**
```c
RLAPI void DrawMesh(Mesh mesh, Material material, Matrix transform)
```

**Koda Binding:**
```koda
native DrawMesh(mesh, material, transform);
```

### DrawMeshInstanced

**Signature:**
```c
RLAPI void DrawMeshInstanced(Mesh mesh, Material material, const Matrix *transforms, int instances)
```

**Koda Binding:**
```koda
native DrawMeshInstanced(mesh, material, *transforms, instances);
```

### GetMeshBoundingBox

**Signature:**
```c
RLAPI BoundingBox GetMeshBoundingBox(Mesh mesh)
```

**Koda Binding:**
```koda
native GetMeshBoundingBox(mesh);
```

### GenMeshTangents

**Signature:**
```c
RLAPI void GenMeshTangents(Mesh *mesh)
```

**Koda Binding:**
```koda
native GenMeshTangents(*mesh);
```

### ExportMesh

**Signature:**
```c
RLAPI bool ExportMesh(Mesh mesh, const char *fileName)
```

**Koda Binding:**
```koda
native ExportMesh(mesh, *fileName);
```

### ExportMeshAsCode

**Signature:**
```c
RLAPI bool ExportMeshAsCode(Mesh mesh, const char *fileName)
```

**Koda Binding:**
```koda
native ExportMeshAsCode(mesh, *fileName);
```

### GenMeshPoly

**Signature:**
```c
RLAPI Mesh GenMeshPoly(int sides, float radius)
```

**Koda Binding:**
```koda
native GenMeshPoly(sides, radius);
```

### GenMeshPlane

**Signature:**
```c
RLAPI Mesh GenMeshPlane(float width, float length, int resX, int resZ)
```

**Koda Binding:**
```koda
native GenMeshPlane(width, length, resX, resZ);
```

### GenMeshCube

**Signature:**
```c
RLAPI Mesh GenMeshCube(float width, float height, float length)
```

**Koda Binding:**
```koda
native GenMeshCube(width, height, length);
```

### GenMeshSphere

**Signature:**
```c
RLAPI Mesh GenMeshSphere(float radius, int rings, int slices)
```

**Koda Binding:**
```koda
native GenMeshSphere(radius, rings, slices);
```

### GenMeshHemiSphere

**Signature:**
```c
RLAPI Mesh GenMeshHemiSphere(float radius, int rings, int slices)
```

**Koda Binding:**
```koda
native GenMeshHemiSphere(radius, rings, slices);
```

### GenMeshCylinder

**Signature:**
```c
RLAPI Mesh GenMeshCylinder(float radius, float height, int slices)
```

**Koda Binding:**
```koda
native GenMeshCylinder(radius, height, slices);
```

### GenMeshCone

**Signature:**
```c
RLAPI Mesh GenMeshCone(float radius, float height, int slices)
```

**Koda Binding:**
```koda
native GenMeshCone(radius, height, slices);
```

### GenMeshTorus

**Signature:**
```c
RLAPI Mesh GenMeshTorus(float radius, float size, int radSeg, int sides)
```

**Koda Binding:**
```koda
native GenMeshTorus(radius, size, radSeg, sides);
```

### GenMeshKnot

**Signature:**
```c
RLAPI Mesh GenMeshKnot(float radius, float size, int radSeg, int sides)
```

**Koda Binding:**
```koda
native GenMeshKnot(radius, size, radSeg, sides);
```

### GenMeshHeightmap

**Signature:**
```c
RLAPI Mesh GenMeshHeightmap(Image heightmap, Vector3 size)
```

**Koda Binding:**
```koda
native GenMeshHeightmap(heightmap, size);
```

### GenMeshCubicmap

**Signature:**
```c
RLAPI Mesh GenMeshCubicmap(Image cubicmap, Vector3 cubeSize)
```

**Koda Binding:**
```koda
native GenMeshCubicmap(cubicmap, cubeSize);
```

### LoadMaterialDefault

**Signature:**
```c
RLAPI Material LoadMaterialDefault()
```

**Koda Binding:**
```koda
native LoadMaterialDefault();
```

### IsMaterialValid

**Signature:**
```c
RLAPI bool IsMaterialValid(Material material)
```

**Koda Binding:**
```koda
native IsMaterialValid(material);
```

### UnloadMaterial

**Signature:**
```c
RLAPI void UnloadMaterial(Material material)
```

**Koda Binding:**
```koda
native UnloadMaterial(material);
```

### SetMaterialTexture

**Signature:**
```c
RLAPI void SetMaterialTexture(Material *material, int mapType, Texture2D texture)
```

**Koda Binding:**
```koda
native SetMaterialTexture(*material, mapType, texture);
```

### SetModelMeshMaterial

**Signature:**
```c
RLAPI void SetModelMeshMaterial(Model *model, int meshId, int materialId)
```

**Koda Binding:**
```koda
native SetModelMeshMaterial(*model, meshId, materialId);
```

### UpdateModelAnimation

**Signature:**
```c
RLAPI void UpdateModelAnimation(Model model, ModelAnimation anim, float frame)
```

**Koda Binding:**
```koda
native UpdateModelAnimation(model, anim, frame);
```

### UpdateModelAnimationEx

**Signature:**
```c
RLAPI void UpdateModelAnimationEx(Model model, ModelAnimation animA, float frameA, ModelAnimation animB, float frameB, float blend)
```

**Koda Binding:**
```koda
native UpdateModelAnimationEx(model, animA, frameA, animB, frameB, blend);
```

### UnloadModelAnimations

**Signature:**
```c
RLAPI void UnloadModelAnimations(ModelAnimation *animations, int animCount)
```

**Koda Binding:**
```koda
native UnloadModelAnimations(*animations, animCount);
```

### IsModelAnimationValid

**Signature:**
```c
RLAPI bool IsModelAnimationValid(Model model, ModelAnimation anim)
```

**Koda Binding:**
```koda
native IsModelAnimationValid(model, anim);
```

### CheckCollisionSpheres

**Signature:**
```c
RLAPI bool CheckCollisionSpheres(Vector3 center1, float radius1, Vector3 center2, float radius2)
```

**Koda Binding:**
```koda
native CheckCollisionSpheres(center1, radius1, center2, radius2);
```

### CheckCollisionBoxes

**Signature:**
```c
RLAPI bool CheckCollisionBoxes(BoundingBox box1, BoundingBox box2)
```

**Koda Binding:**
```koda
native CheckCollisionBoxes(box1, box2);
```

### CheckCollisionBoxSphere

**Signature:**
```c
RLAPI bool CheckCollisionBoxSphere(BoundingBox box, Vector3 center, float radius)
```

**Koda Binding:**
```koda
native CheckCollisionBoxSphere(box, center, radius);
```

### GetRayCollisionSphere

**Signature:**
```c
RLAPI RayCollision GetRayCollisionSphere(Ray ray, Vector3 center, float radius)
```

**Koda Binding:**
```koda
native GetRayCollisionSphere(ray, center, radius);
```

### GetRayCollisionBox

**Signature:**
```c
RLAPI RayCollision GetRayCollisionBox(Ray ray, BoundingBox box)
```

**Koda Binding:**
```koda
native GetRayCollisionBox(ray, box);
```

### GetRayCollisionMesh

**Signature:**
```c
RLAPI RayCollision GetRayCollisionMesh(Ray ray, Mesh mesh, Matrix transform)
```

**Koda Binding:**
```koda
native GetRayCollisionMesh(ray, mesh, transform);
```

### GetRayCollisionTriangle

**Signature:**
```c
RLAPI RayCollision GetRayCollisionTriangle(Ray ray, Vector3 p1, Vector3 p2, Vector3 p3)
```

**Koda Binding:**
```koda
native GetRayCollisionTriangle(ray, p1, p2, p3);
```

### GetRayCollisionQuad

**Signature:**
```c
RLAPI RayCollision GetRayCollisionQuad(Ray ray, Vector3 p1, Vector3 p2, Vector3 p3, Vector3 p4)
```

**Koda Binding:**
```koda
native GetRayCollisionQuad(ray, p1, p2, p3, p4);
```

### void

**Signature:**
```c
typedef void()
```

**Koda Binding:**
```koda
native void();
```

### InitAudioDevice

**Signature:**
```c
RLAPI void InitAudioDevice()
```

**Koda Binding:**
```koda
native InitAudioDevice();
```

### CloseAudioDevice

**Signature:**
```c
RLAPI void CloseAudioDevice()
```

**Koda Binding:**
```koda
native CloseAudioDevice();
```

### IsAudioDeviceReady

**Signature:**
```c
RLAPI bool IsAudioDeviceReady()
```

**Koda Binding:**
```koda
native IsAudioDeviceReady();
```

### SetMasterVolume

**Signature:**
```c
RLAPI void SetMasterVolume(float volume)
```

**Koda Binding:**
```koda
native SetMasterVolume(volume);
```

### GetMasterVolume

**Signature:**
```c
RLAPI float GetMasterVolume()
```

**Koda Binding:**
```koda
native GetMasterVolume();
```

### LoadWave

**Signature:**
```c
RLAPI Wave LoadWave(const char *fileName)
```

**Koda Binding:**
```koda
native LoadWave(*fileName);
```

### LoadWaveFromMemory

**Signature:**
```c
RLAPI Wave LoadWaveFromMemory(const char *fileType, const unsigned char *fileData, int dataSize)
```

**Koda Binding:**
```koda
native LoadWaveFromMemory(*fileType, *fileData, dataSize);
```

### IsWaveValid

**Signature:**
```c
RLAPI bool IsWaveValid(Wave wave)
```

**Koda Binding:**
```koda
native IsWaveValid(wave);
```

### LoadSound

**Signature:**
```c
RLAPI Sound LoadSound(const char *fileName)
```

**Koda Binding:**
```koda
native LoadSound(*fileName);
```

### LoadSoundFromWave

**Signature:**
```c
RLAPI Sound LoadSoundFromWave(Wave wave)
```

**Koda Binding:**
```koda
native LoadSoundFromWave(wave);
```

### LoadSoundAlias

**Signature:**
```c
RLAPI Sound LoadSoundAlias(Sound source)
```

**Koda Binding:**
```koda
native LoadSoundAlias(source);
```

### IsSoundValid

**Signature:**
```c
RLAPI bool IsSoundValid(Sound sound)
```

**Koda Binding:**
```koda
native IsSoundValid(sound);
```

### UpdateSound

**Signature:**
```c
RLAPI void UpdateSound(Sound sound, const void *data, int frameCount)
```

**Koda Binding:**
```koda
native UpdateSound(sound, *data, frameCount);
```

### UnloadWave

**Signature:**
```c
RLAPI void UnloadWave(Wave wave)
```

**Koda Binding:**
```koda
native UnloadWave(wave);
```

### UnloadSound

**Signature:**
```c
RLAPI void UnloadSound(Sound sound)
```

**Koda Binding:**
```koda
native UnloadSound(sound);
```

### UnloadSoundAlias

**Signature:**
```c
RLAPI void UnloadSoundAlias(Sound alias)
```

**Koda Binding:**
```koda
native UnloadSoundAlias(alias);
```

### ExportWave

**Signature:**
```c
RLAPI bool ExportWave(Wave wave, const char *fileName)
```

**Koda Binding:**
```koda
native ExportWave(wave, *fileName);
```

### ExportWaveAsCode

**Signature:**
```c
RLAPI bool ExportWaveAsCode(Wave wave, const char *fileName)
```

**Koda Binding:**
```koda
native ExportWaveAsCode(wave, *fileName);
```

### PlaySound

**Signature:**
```c
RLAPI void PlaySound(Sound sound)
```

**Koda Binding:**
```koda
native PlaySound(sound);
```

### StopSound

**Signature:**
```c
RLAPI void StopSound(Sound sound)
```

**Koda Binding:**
```koda
native StopSound(sound);
```

### PauseSound

**Signature:**
```c
RLAPI void PauseSound(Sound sound)
```

**Koda Binding:**
```koda
native PauseSound(sound);
```

### ResumeSound

**Signature:**
```c
RLAPI void ResumeSound(Sound sound)
```

**Koda Binding:**
```koda
native ResumeSound(sound);
```

### IsSoundPlaying

**Signature:**
```c
RLAPI bool IsSoundPlaying(Sound sound)
```

**Koda Binding:**
```koda
native IsSoundPlaying(sound);
```

### SetSoundVolume

**Signature:**
```c
RLAPI void SetSoundVolume(Sound sound, float volume)
```

**Koda Binding:**
```koda
native SetSoundVolume(sound, volume);
```

### SetSoundPitch

**Signature:**
```c
RLAPI void SetSoundPitch(Sound sound, float pitch)
```

**Koda Binding:**
```koda
native SetSoundPitch(sound, pitch);
```

### SetSoundPan

**Signature:**
```c
RLAPI void SetSoundPan(Sound sound, float pan)
```

**Koda Binding:**
```koda
native SetSoundPan(sound, pan);
```

### WaveCopy

**Signature:**
```c
RLAPI Wave WaveCopy(Wave wave)
```

**Koda Binding:**
```koda
native WaveCopy(wave);
```

### WaveCrop

**Signature:**
```c
RLAPI void WaveCrop(Wave *wave, int initFrame, int finalFrame)
```

**Koda Binding:**
```koda
native WaveCrop(*wave, initFrame, finalFrame);
```

### WaveFormat

**Signature:**
```c
RLAPI void WaveFormat(Wave *wave, int sampleRate, int sampleSize, int channels)
```

**Koda Binding:**
```koda
native WaveFormat(*wave, sampleRate, sampleSize, channels);
```

### UnloadWaveSamples

**Signature:**
```c
RLAPI void UnloadWaveSamples(float *samples)
```

**Koda Binding:**
```koda
native UnloadWaveSamples(*samples);
```

### LoadMusicStream

**Signature:**
```c
RLAPI Music LoadMusicStream(const char *fileName)
```

**Koda Binding:**
```koda
native LoadMusicStream(*fileName);
```

### LoadMusicStreamFromMemory

**Signature:**
```c
RLAPI Music LoadMusicStreamFromMemory(const char *fileType, const unsigned char *data, int dataSize)
```

**Koda Binding:**
```koda
native LoadMusicStreamFromMemory(*fileType, *data, dataSize);
```

### IsMusicValid

**Signature:**
```c
RLAPI bool IsMusicValid(Music music)
```

**Koda Binding:**
```koda
native IsMusicValid(music);
```

### UnloadMusicStream

**Signature:**
```c
RLAPI void UnloadMusicStream(Music music)
```

**Koda Binding:**
```koda
native UnloadMusicStream(music);
```

### PlayMusicStream

**Signature:**
```c
RLAPI void PlayMusicStream(Music music)
```

**Koda Binding:**
```koda
native PlayMusicStream(music);
```

### IsMusicStreamPlaying

**Signature:**
```c
RLAPI bool IsMusicStreamPlaying(Music music)
```

**Koda Binding:**
```koda
native IsMusicStreamPlaying(music);
```

### UpdateMusicStream

**Signature:**
```c
RLAPI void UpdateMusicStream(Music music)
```

**Koda Binding:**
```koda
native UpdateMusicStream(music);
```

### StopMusicStream

**Signature:**
```c
RLAPI void StopMusicStream(Music music)
```

**Koda Binding:**
```koda
native StopMusicStream(music);
```

### PauseMusicStream

**Signature:**
```c
RLAPI void PauseMusicStream(Music music)
```

**Koda Binding:**
```koda
native PauseMusicStream(music);
```

### ResumeMusicStream

**Signature:**
```c
RLAPI void ResumeMusicStream(Music music)
```

**Koda Binding:**
```koda
native ResumeMusicStream(music);
```

### SeekMusicStream

**Signature:**
```c
RLAPI void SeekMusicStream(Music music, float position)
```

**Koda Binding:**
```koda
native SeekMusicStream(music, position);
```

### SetMusicVolume

**Signature:**
```c
RLAPI void SetMusicVolume(Music music, float volume)
```

**Koda Binding:**
```koda
native SetMusicVolume(music, volume);
```

### SetMusicPitch

**Signature:**
```c
RLAPI void SetMusicPitch(Music music, float pitch)
```

**Koda Binding:**
```koda
native SetMusicPitch(music, pitch);
```

### SetMusicPan

**Signature:**
```c
RLAPI void SetMusicPan(Music music, float pan)
```

**Koda Binding:**
```koda
native SetMusicPan(music, pan);
```

### GetMusicTimeLength

**Signature:**
```c
RLAPI float GetMusicTimeLength(Music music)
```

**Koda Binding:**
```koda
native GetMusicTimeLength(music);
```

### GetMusicTimePlayed

**Signature:**
```c
RLAPI float GetMusicTimePlayed(Music music)
```

**Koda Binding:**
```koda
native GetMusicTimePlayed(music);
```

### LoadAudioStream

**Signature:**
```c
RLAPI AudioStream LoadAudioStream(unsigned int sampleRate, unsigned int sampleSize, unsigned int channels)
```

**Koda Binding:**
```koda
native LoadAudioStream(sampleRate, sampleSize, channels);
```

### IsAudioStreamValid

**Signature:**
```c
RLAPI bool IsAudioStreamValid(AudioStream stream)
```

**Koda Binding:**
```koda
native IsAudioStreamValid(stream);
```

### UnloadAudioStream

**Signature:**
```c
RLAPI void UnloadAudioStream(AudioStream stream)
```

**Koda Binding:**
```koda
native UnloadAudioStream(stream);
```

### UpdateAudioStream

**Signature:**
```c
RLAPI void UpdateAudioStream(AudioStream stream, const void *data, int frameCount)
```

**Koda Binding:**
```koda
native UpdateAudioStream(stream, *data, frameCount);
```

### IsAudioStreamProcessed

**Signature:**
```c
RLAPI bool IsAudioStreamProcessed(AudioStream stream)
```

**Koda Binding:**
```koda
native IsAudioStreamProcessed(stream);
```

### PlayAudioStream

**Signature:**
```c
RLAPI void PlayAudioStream(AudioStream stream)
```

**Koda Binding:**
```koda
native PlayAudioStream(stream);
```

### PauseAudioStream

**Signature:**
```c
RLAPI void PauseAudioStream(AudioStream stream)
```

**Koda Binding:**
```koda
native PauseAudioStream(stream);
```

### ResumeAudioStream

**Signature:**
```c
RLAPI void ResumeAudioStream(AudioStream stream)
```

**Koda Binding:**
```koda
native ResumeAudioStream(stream);
```

### IsAudioStreamPlaying

**Signature:**
```c
RLAPI bool IsAudioStreamPlaying(AudioStream stream)
```

**Koda Binding:**
```koda
native IsAudioStreamPlaying(stream);
```

### StopAudioStream

**Signature:**
```c
RLAPI void StopAudioStream(AudioStream stream)
```

**Koda Binding:**
```koda
native StopAudioStream(stream);
```

### SetAudioStreamVolume

**Signature:**
```c
RLAPI void SetAudioStreamVolume(AudioStream stream, float volume)
```

**Koda Binding:**
```koda
native SetAudioStreamVolume(stream, volume);
```

### SetAudioStreamPitch

**Signature:**
```c
RLAPI void SetAudioStreamPitch(AudioStream stream, float pitch)
```

**Koda Binding:**
```koda
native SetAudioStreamPitch(stream, pitch);
```

### SetAudioStreamPan

**Signature:**
```c
RLAPI void SetAudioStreamPan(AudioStream stream, float pan)
```

**Koda Binding:**
```koda
native SetAudioStreamPan(stream, pan);
```

### SetAudioStreamBufferSizeDefault

**Signature:**
```c
RLAPI void SetAudioStreamBufferSizeDefault(int size)
```

**Koda Binding:**
```koda
native SetAudioStreamBufferSizeDefault(size);
```

### SetAudioStreamCallback

**Signature:**
```c
RLAPI void SetAudioStreamCallback(AudioStream stream, AudioCallback callback)
```

**Koda Binding:**
```koda
native SetAudioStreamCallback(stream, callback);
```

### AttachAudioStreamProcessor

**Signature:**
```c
RLAPI void AttachAudioStreamProcessor(AudioStream stream, AudioCallback processor)
```

**Koda Binding:**
```koda
native AttachAudioStreamProcessor(stream, processor);
```

### DetachAudioStreamProcessor

**Signature:**
```c
RLAPI void DetachAudioStreamProcessor(AudioStream stream, AudioCallback processor)
```

**Koda Binding:**
```koda
native DetachAudioStreamProcessor(stream, processor);
```

### AttachAudioMixedProcessor

**Signature:**
```c
RLAPI void AttachAudioMixedProcessor(AudioCallback processor)
```

**Koda Binding:**
```koda
native AttachAudioMixedProcessor(processor);
```

### DetachAudioMixedProcessor

**Signature:**
```c
RLAPI void DetachAudioMixedProcessor(AudioCallback processor)
```

**Koda Binding:**
```koda
native DetachAudioMixedProcessor(processor);
```


## Structs

### Vector2

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| x | float | | 
| y | // Vector x component float | | 
| component | // Vector y | | 

### Vector3

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| x | float | | 
| y | // Vector x component float | | 
| z | // Vector y component float | | 
| component | // Vector z | | 

### Vector4

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| x | float | | 
| y | // Vector x component float | | 
| z | // Vector y component float | | 
| w | // Vector z component float | | 
| component | // Vector w | | 

### Matrix

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| m12 | float m0, m4, m8, | | 
| m13 | // Matrix first row (4 components) float m1, m5, m9, | | 
| m14 | // Matrix second row (4 components) float m2, m6, m10, | | 
| m15 | // Matrix third row (4 components) float m3, m7, m11, | | 
| components) | // Matrix fourth row (4 | | 

### Color

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| r | unsigned char | | 
| g | // Color red value unsigned char | | 
| b | // Color green value unsigned char | | 
| a | // Color blue value unsigned char | | 
| value | // Color alpha | | 

### Rectangle

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| x | float | | 
| y | // Rectangle top-left corner position x float | | 
| width | // Rectangle top-left corner position y float | | 
| height | // Rectangle width float | | 
| height | // Rectangle | | 

### Image

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| *data | void | | 
| width | // Image raw data int | | 
| height | // Image base width int | | 
| mipmaps | // Image base height int | | 
| format | // Mipmap levels, 1 by default int | | 
| type) | // Data format (PixelFormat | | 

### Texture

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| id | unsigned int | | 
| width | // OpenGL texture id int | | 
| height | // Texture base width int | | 
| mipmaps | // Texture base height int | | 
| format | // Mipmap levels, 1 by default int | | 
| type) | // Data format (PixelFormat | | 

### RenderTexture

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| id | unsigned int | | 
| texture | // OpenGL framebuffer object id Texture | | 
| depth | // Color buffer attachment texture Texture | | 
| texture | // Depth buffer attachment | | 

### NPatchInfo

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| source | Rectangle | | 
| left | // Texture source rectangle int | | 
| top | // Left border offset int | | 
| right | // Top border offset int | | 
| bottom | // Right border offset int | | 
| layout | // Bottom border offset int | | 
| 3x1 | // Layout of the n-patch: 3x3, 1x3 or | | 

### GlyphInfo

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| value | int | | 
| offsetX | // Character value (Unicode) int | | 
| offsetY | // Character offset X when drawing int | | 
| advanceX | // Character offset Y when drawing int | | 
| image | // Character advance position X Image | | 
| data | // Character image | | 

### Font

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| baseSize | int | | 
| glyphCount | // Base size (default chars height) int | | 
| glyphPadding | // Number of glyph characters int | | 
| texture | // Padding around the glyph characters Texture2D | | 
| *recs | // Texture atlas containing the glyphs Rectangle | | 
| *glyphs | // Rectangles in texture for the glyphs GlyphInfo | | 
| data | // Glyphs info | | 

### Camera3D

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| position | Vector3 | | 
| target | // Camera position Vector3 | | 
| up | // Camera target it looks-at Vector3 | | 
| fovy | // Camera up vector (rotation over its axis) float | | 
| projection | // Camera field-of-view aperture in Y (degrees) in perspective, used as near plane height in world units in orthographic int | | 
| CAMERA_ORTHOGRAPHIC | // Camera projection: CAMERA_PERSPECTIVE or | | 

### Camera2D

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| offset | Vector2 | | 
| target | // Camera offset (screen space offset from window origin) Vector2 | | 
| rotation | // Camera target (world space target point that is mapped to screen space offset) float | | 
| zoom | // Camera rotation in degrees (pivots around target) float | | 
| scale | // Camera zoom (scaling around target), must not be set to 0, set to 1.0f for no | | 

### Mesh

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| vertexCount | int | | 
| triangleCount | // Number of vertices stored in arrays int | | 
| *vertices | // Number of triangles stored (indexed or not) // Vertex attributes data float | | 
| *texcoords | // Vertex position (XYZ - 3 components per vertex) (shader-location = 0) float | | 
| *texcoords2 | // Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1) float | | 
| *normals | // Vertex texture second coordinates (UV - 2 components per vertex) (shader-location = 5) float | | 
| *tangents | // Vertex normals (XYZ - 3 components per vertex) (shader-location = 2) float | | 
| *colors | // Vertex tangents (XYZW - 4 components per vertex) (shader-location = 4) unsigned char | | 
| *indices | // Vertex colors (RGBA - 4 components per vertex) (shader-location = 3) unsigned short | | 
| boneCount | // Vertex indices (in case vertex data comes indexed) // Skin data for animation int | | 
| *boneIndices | // Number of bones (MAX: 256 bones) unsigned char | | 
| *boneWeights | // Vertex bone indices, up to 4 bones influence by vertex (skinning) (shader-location = 6) float | | 
| *animVertices | // Vertex bone weight, up to 4 bones influence by vertex (skinning) (shader-location = 7) // Runtime animation vertex data (CPU skinning) // NOTE: In case of GPU skinning, not used, pointers are NULL float | | 
| *animNormals | // Animated vertex positions (after bones transformations) float | | 
| vaoId | // Animated normals (after bones transformations) // OpenGL identifiers unsigned int | | 
| *vboId | // OpenGL Vertex Array Object id unsigned int | | 
| data) | // OpenGL Vertex Buffer Objects id (default vertex | | 

### Shader

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| id | unsigned int | | 
| *locs | // Shader program id int | | 
| (RL_MAX_SHADER_LOCATIONS) | // Shader locations array | | 

### MaterialMap

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| texture | Texture2D | | 
| color | // Material map texture Color | | 
| value | // Material map color float | | 
| value | // Material map | | 

### Material

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| shader | Shader | | 
| *maps | // Material shader MaterialMap | | 
| params[4] | // Material maps array (MAX_MATERIAL_MAPS) float | | 
| required) | // Material generic parameters (if | | 

### Transform

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| translation | Vector3 | | 
| rotation | // Translation Quaternion | | 
| scale | // Rotation Vector3 | | 
| Scale | // | | 

### BoneInfo

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| name[32] | char | | 
| parent | // Bone name int | | 
| parent | // Bone | | 

### ModelSkeleton

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| boneCount | int | | 
| *bones | // Number of bones BoneInfo | | 
| bindPose | // Bones information (skeleton) ModelAnimPose | | 
| (Transform[]) | // Bones base transformation | | 

### Model

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| transform | Matrix | | 
| meshCount | // Local transform matrix int | | 
| materialCount | // Number of meshes int | | 
| *meshes | // Number of materials Mesh | | 
| *materials | // Meshes array Material | | 
| *meshMaterial | // Materials array int | | 
| skeleton | // Mesh material number // Animation data ModelSkeleton | | 
| currentPose | // Skeleton for animation // Runtime animation data (CPU/GPU skinning) ModelAnimPose | | 
| *boneMatrices | // Current animation pose (Transform[]) Matrix | | 
| matrices | // Bones animated transformation | | 

### ModelAnimation

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| name[32] | char | | 
| boneCount | // Animation name int | | 
| keyframeCount | // Number of bones (per pose) int | | 
| *keyframePoses | // Number of animation key frames ModelAnimPose | | 
| [keyframe][pose] | // Animation sequence keyframe poses | | 

### Ray

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| position | Vector3 | | 
| direction | // Ray position (origin) Vector3 | | 
| (normalized) | // Ray direction | | 

### RayCollision

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| hit | bool | | 
| distance | // Did the ray hit something? float | | 
| point | // Distance to the nearest hit Vector3 | | 
| normal | // Point of the nearest hit Vector3 | | 
| hit | // Surface normal of | | 

### BoundingBox

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| min | Vector3 | | 
| max | // Minimum vertex box-corner Vector3 | | 
| box-corner | // Maximum vertex | | 

### Wave

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| frameCount | unsigned int | | 
| sampleRate | // Total number of frames (considering channels) unsigned int | | 
| sampleSize | // Frequency (samples per second) unsigned int | | 
| channels | // Bit depth (bits per sample): 8, 16, 32 (24 not supported) unsigned int | | 
| *data | // Number of channels (1-mono, 2-stereo, ...) void | | 
| pointer | // Buffer data | | 

### AudioStream

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| *buffer | rAudioBuffer | | 
| *processor | // Pointer to internal data used by the audio system rAudioProcessor | | 
| sampleRate | // Pointer to internal data processor, useful for audio effects unsigned int | | 
| sampleSize | // Frequency (samples per second) unsigned int | | 
| channels | // Bit depth (bits per sample): 8, 16, 32 (24 not supported) unsigned int | | 
| ...) | // Number of channels (1-mono, 2-stereo, | | 

### Sound

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| stream | AudioStream | | 
| frameCount | // Audio stream unsigned int | | 
| channels) | // Total number of frames (considering | | 

### Music

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| stream | AudioStream | | 
| frameCount | // Audio stream unsigned int | | 
| looping | // Total number of frames (considering channels) bool | | 
| ctxType | // Music looping enable int | | 
| *ctxData | // Type of music context (audio filetype) void | | 
| type | // Audio context data, depends on | | 

### VrDeviceInfo

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| hResolution | int | | 
| vResolution | // Horizontal resolution in pixels int | | 
| hScreenSize | // Vertical resolution in pixels float | | 
| vScreenSize | // Horizontal size in meters float | | 
| eyeToScreenDistance | // Vertical size in meters float | | 
| lensSeparationDistance | // Distance between eye and display in meters float | | 
| interpupillaryDistance | // Lens separation distance in meters float | | 
| lensDistortionValues[4] | // IPD (distance between pupils) in meters float | | 
| chromaAbCorrection[4] | // Lens distortion constant parameters float | | 
| parameters | // Chromatic aberration correction | | 

### VrStereoConfig

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| projection[2] | Matrix | | 
| viewOffset[2] | // VR projection matrices (per eye) Matrix | | 
| leftLensCenter[2] | // VR view offset matrices (per eye) float | | 
| rightLensCenter[2] | // VR left lens center float | | 
| leftScreenCenter[2] | // VR right lens center float | | 
| rightScreenCenter[2] | // VR left screen center float | | 
| scale[2] | // VR right screen center float | | 
| scaleIn[2] | // VR distortion scale float | | 
| in | // VR distortion scale | | 

### FilePathList

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| count | unsigned int | | 
| **paths | // Filepaths entries count char | | 
| entries | // Filepaths | | 

### AutomationEvent

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| frame | unsigned int | | 
| type | // Event frame unsigned int | | 
| params[4] | // Event type (AutomationEventType) int | | 
| required) | // Event parameters (if | | 

### AutomationEventList

**Fields:**

| Name | Type | Description |
|------|------|-------------|
| capacity | unsigned int | | 
| count | // Events max entries (MAX_AUTOMATION_EVENTS) unsigned int | | 
| *events | // Events entries count AutomationEvent | | 
| entries | // Events | | 


## Enums

### bool

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| false | 0 | | 
| true | 1 | | 

### ConfigFlags

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| FLAG_VSYNC_HINT | 0 | | 
| // Set to try enabling V-Sync on GPU
    FLAG_FULLSCREEN_MODE | 0 | | 
| // Set to run program in fullscreen
    FLAG_WINDOW_RESIZABLE | 0 | | 
| // Set to allow resizable window
    FLAG_WINDOW_UNDECORATED | 0 | | 
| // Set to disable window decoration (frame and buttons)
    FLAG_WINDOW_HIDDEN | 0 | | 
| // Set to hide window
    FLAG_WINDOW_MINIMIZED | 0 | | 
| // Set to minimize window (iconify)
    FLAG_WINDOW_MAXIMIZED | 0 | | 
| // Set to maximize window (expanded to monitor)
    FLAG_WINDOW_UNFOCUSED | 0 | | 
| // Set to window non focused
    FLAG_WINDOW_TOPMOST | 0 | | 
| // Set to window always on top
    FLAG_WINDOW_ALWAYS_RUN | 0 | | 
| // Set to allow windows running while minimized
    FLAG_WINDOW_TRANSPARENT | 0 | | 
| // Set to allow transparent framebuffer
    FLAG_WINDOW_HIGHDPI | 0 | | 
| // Set to support HighDPI
    FLAG_WINDOW_MOUSE_PASSTHROUGH | 0 | | 
| // Set to support mouse passthrough | 13 | | 
| only supported when FLAG_WINDOW_UNDECORATED
    FLAG_BORDERLESS_WINDOWED_MODE | 0 | | 
| // Set to run program in borderless windowed mode
    FLAG_MSAA_4X_HINT | 0 | | 
| // Set to try enabling MSAA 4X
    FLAG_INTERLACED_HINT | 0 | | 

### TraceLogLevel

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| LOG_ALL | 0 | | 
| // Display all logs
    LOG_TRACE | 1 | | 
| // Trace logging | 2 | | 
| intended for internal use only
    LOG_DEBUG | 3 | | 
| // Debug logging | 4 | | 
| used for internal debugging | 5 | | 
| it should be disabled on release builds
    LOG_INFO | 6 | | 
| // Info logging | 7 | | 
| used for program execution info
    LOG_WARNING | 8 | | 
| // Warning logging | 9 | | 
| used on recoverable failures
    LOG_ERROR | 10 | | 
| // Error logging | 11 | | 
| used on unrecoverable failures
    LOG_FATAL | 12 | | 
| // Fatal logging | 13 | | 
| used to abort program: exit(EXIT_FAILURE)
    LOG_NONE            // Disable logging | 14 | | 

### KeyboardKey

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| KEY_NULL | 0 | | 
| // Key: NULL | 1 | | 
| used for no key pressed
    // Alphanumeric keys
    KEY_APOSTROPHE | 39 | | 
| // Key: '
    KEY_COMMA | 44 | | 
| // Key: | 4 | | 
| KEY_MINUS | 45 | | 
| // Key: -
    KEY_PERIOD | 46 | | 
| // Key: .
    KEY_SLASH | 47 | | 
| // Key: /
    KEY_ZERO | 48 | | 
| // Key: 0
    KEY_ONE | 49 | | 
| // Key: 1
    KEY_TWO | 50 | | 
| // Key: 2
    KEY_THREE | 51 | | 
| // Key: 3
    KEY_FOUR | 52 | | 
| // Key: 4
    KEY_FIVE | 53 | | 
| // Key: 5
    KEY_SIX | 54 | | 
| // Key: 6
    KEY_SEVEN | 55 | | 
| // Key: 7
    KEY_EIGHT | 56 | | 
| // Key: 8
    KEY_NINE | 57 | | 
| // Key: 9
    KEY_SEMICOLON | 59 | | 
| // Key: ;
    KEY_EQUAL | 61 | | 
| // Key: | 20 | | 
| // Key: A | a
    KEY_B | 66 | | 
| // Key: B | b
    KEY_C | 67 | | 
| // Key: C | c
    KEY_D | 68 | | 
| // Key: D | d
    KEY_E | 69 | | 
| // Key: E | e
    KEY_F | 70 | | 
| // Key: F | f
    KEY_G | 71 | | 
| // Key: G | g
    KEY_H | 72 | | 
| // Key: H | h
    KEY_I | 73 | | 
| // Key: I | i
    KEY_J | 74 | | 
| // Key: J | j
    KEY_K | 75 | | 
| // Key: K | k
    KEY_L | 76 | | 
| // Key: L | l
    KEY_M | 77 | | 
| // Key: M | m
    KEY_N | 78 | | 
| // Key: N | n
    KEY_O | 79 | | 
| // Key: O | o
    KEY_P | 80 | | 
| // Key: P | p
    KEY_Q | 81 | | 
| // Key: Q | q
    KEY_R | 82 | | 
| // Key: R | r
    KEY_S | 83 | | 
| // Key: S | s
    KEY_T | 84 | | 
| // Key: T | t
    KEY_U | 85 | | 
| // Key: U | u
    KEY_V | 86 | | 
| // Key: V | v
    KEY_W | 87 | | 
| // Key: W | w
    KEY_X | 88 | | 
| // Key: X | x
    KEY_Y | 89 | | 
| // Key: Y | y
    KEY_Z | 90 | | 
| // Key: Z | z
    KEY_LEFT_BRACKET | 91 | | 
| // Key: [
    KEY_BACKSLASH | 92 | | 
| // Key: '\'
    KEY_RIGHT_BRACKET | 93 | | 
| // Key: ]
    KEY_GRAVE | 96 | | 
| // Key: `
    // Function keys
    KEY_SPACE | 32 | | 
| // Key: Space
    KEY_ESCAPE | 256 | | 
| // Key: Esc
    KEY_ENTER | 257 | | 
| // Key: Enter
    KEY_TAB | 258 | | 
| // Key: Tab
    KEY_BACKSPACE | 259 | | 
| // Key: Backspace
    KEY_INSERT | 260 | | 
| // Key: Ins
    KEY_DELETE | 261 | | 
| // Key: Del
    KEY_RIGHT | 262 | | 
| // Key: Cursor right
    KEY_LEFT | 263 | | 
| // Key: Cursor left
    KEY_DOWN | 264 | | 
| // Key: Cursor down
    KEY_UP | 265 | | 
| // Key: Cursor up
    KEY_PAGE_UP | 266 | | 
| // Key: Page up
    KEY_PAGE_DOWN | 267 | | 
| // Key: Page down
    KEY_HOME | 268 | | 
| // Key: Home
    KEY_END | 269 | | 
| // Key: End
    KEY_CAPS_LOCK | 280 | | 
| // Key: Caps lock
    KEY_SCROLL_LOCK | 281 | | 
| // Key: Scroll down
    KEY_NUM_LOCK | 282 | | 
| // Key: Num lock
    KEY_PRINT_SCREEN | 283 | | 
| // Key: Print screen
    KEY_PAUSE | 284 | | 
| // Key: Pause
    KEY_F1 | 290 | | 
| // Key: F1
    KEY_F2 | 291 | | 
| // Key: F2
    KEY_F3 | 292 | | 
| // Key: F3
    KEY_F4 | 293 | | 
| // Key: F4
    KEY_F5 | 294 | | 
| // Key: F5
    KEY_F6 | 295 | | 
| // Key: F6
    KEY_F7 | 296 | | 
| // Key: F7
    KEY_F8 | 297 | | 
| // Key: F8
    KEY_F9 | 298 | | 
| // Key: F9
    KEY_F10 | 299 | | 
| // Key: F10
    KEY_F11 | 300 | | 
| // Key: F11
    KEY_F12 | 301 | | 
| // Key: F12
    KEY_LEFT_SHIFT | 340 | | 
| // Key: Shift left
    KEY_LEFT_CONTROL | 341 | | 
| // Key: Control left
    KEY_LEFT_ALT | 342 | | 
| // Key: Alt left
    KEY_LEFT_SUPER | 343 | | 
| // Key: Super left
    KEY_RIGHT_SHIFT | 344 | | 
| // Key: Shift right
    KEY_RIGHT_CONTROL | 345 | | 
| // Key: Control right
    KEY_RIGHT_ALT | 346 | | 
| // Key: Alt right
    KEY_RIGHT_SUPER | 347 | | 
| // Key: Super right
    KEY_KB_MENU | 348 | | 
| // Key: KB menu
    // Keypad keys
    KEY_KP_0 | 320 | | 
| // Key: Keypad 0
    KEY_KP_1 | 321 | | 
| // Key: Keypad 1
    KEY_KP_2 | 322 | | 
| // Key: Keypad 2
    KEY_KP_3 | 323 | | 
| // Key: Keypad 3
    KEY_KP_4 | 324 | | 
| // Key: Keypad 4
    KEY_KP_5 | 325 | | 
| // Key: Keypad 5
    KEY_KP_6 | 326 | | 
| // Key: Keypad 6
    KEY_KP_7 | 327 | | 
| // Key: Keypad 7
    KEY_KP_8 | 328 | | 
| // Key: Keypad 8
    KEY_KP_9 | 329 | | 
| // Key: Keypad 9
    KEY_KP_DECIMAL | 330 | | 
| // Key: Keypad .
    KEY_KP_DIVIDE | 331 | | 
| // Key: Keypad /
    KEY_KP_MULTIPLY | 332 | | 
| // Key: Keypad *
    KEY_KP_SUBTRACT | 333 | | 
| // Key: Keypad -
    KEY_KP_ADD | 334 | | 
| // Key: Keypad +
    KEY_KP_ENTER | 335 | | 
| // Key: Keypad Enter
    KEY_KP_EQUAL | 336 | | 
| // Key: Keypad | 108 | | 
| // Key: Android back button
    KEY_MENU | 5 | | 
| // Key: Android menu button
    KEY_VOLUME_UP | 24 | | 
| // Key: Android volume up button
    KEY_VOLUME_DOWN | 25 | | 

### MouseButton

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| MOUSE_BUTTON_LEFT | 0 | | 
| // Mouse button left
    MOUSE_BUTTON_RIGHT | 1 | | 
| // Mouse button right
    MOUSE_BUTTON_MIDDLE | 2 | | 
| // Mouse button middle (pressed wheel)
    MOUSE_BUTTON_SIDE | 3 | | 
| // Mouse button side (advanced mouse device)
    MOUSE_BUTTON_EXTRA | 4 | | 
| // Mouse button extra (advanced mouse device)
    MOUSE_BUTTON_FORWARD | 5 | | 
| // Mouse button forward (advanced mouse device)
    MOUSE_BUTTON_BACK | 6 | | 
| // Mouse button back (advanced mouse device) | 7 | | 

### MouseCursor

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| MOUSE_CURSOR_DEFAULT | 0 | | 
| // Default pointer shape
    MOUSE_CURSOR_ARROW | 1 | | 
| // Arrow shape
    MOUSE_CURSOR_IBEAM | 2 | | 
| // Text writing cursor shape
    MOUSE_CURSOR_CROSSHAIR | 3 | | 
| // Cross shape
    MOUSE_CURSOR_POINTING_HAND | 4 | | 
| // Pointing hand cursor
    MOUSE_CURSOR_RESIZE_EW | 5 | | 
| // Horizontal resize/move arrow shape
    MOUSE_CURSOR_RESIZE_NS | 6 | | 
| // Vertical resize/move arrow shape
    MOUSE_CURSOR_RESIZE_NWSE | 7 | | 
| // Top-left to bottom-right diagonal resize/move arrow shape
    MOUSE_CURSOR_RESIZE_NESW | 8 | | 
| // The top-right to bottom-left diagonal resize/move arrow shape
    MOUSE_CURSOR_RESIZE_ALL | 9 | | 
| // The omnidirectional resize/move cursor shape
    MOUSE_CURSOR_NOT_ALLOWED | 10 | | 

### GamepadButton

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| GAMEPAD_BUTTON_UNKNOWN | 0 | | 
| // Unknown button | 1 | | 
| for error checking
    GAMEPAD_BUTTON_LEFT_FACE_UP | 2 | | 
| // Gamepad left DPAD up button
    GAMEPAD_BUTTON_LEFT_FACE_RIGHT | 3 | | 
| // Gamepad left DPAD right button
    GAMEPAD_BUTTON_LEFT_FACE_DOWN | 4 | | 
| // Gamepad left DPAD down button
    GAMEPAD_BUTTON_LEFT_FACE_LEFT | 5 | | 
| // Gamepad left DPAD left button
    GAMEPAD_BUTTON_RIGHT_FACE_UP | 6 | | 
| // Gamepad right button up (i.e. PS3: Triangle | 7 | | 
| Xbox: Y)
    GAMEPAD_BUTTON_RIGHT_FACE_RIGHT | 8 | | 
| // Gamepad right button right (i.e. PS3: Circle | 9 | | 
| Xbox: B)
    GAMEPAD_BUTTON_RIGHT_FACE_DOWN | 10 | | 
| // Gamepad right button down (i.e. PS3: Cross | 11 | | 
| Xbox: A)
    GAMEPAD_BUTTON_RIGHT_FACE_LEFT | 12 | | 
| // Gamepad right button left (i.e. PS3: Square | 13 | | 
| Xbox: X)
    GAMEPAD_BUTTON_LEFT_TRIGGER_1 | 14 | | 
| // Gamepad top/back trigger left (first) | 15 | | 
| it could be a trailing button
    GAMEPAD_BUTTON_LEFT_TRIGGER_2 | 16 | | 
| // Gamepad top/back trigger left (second) | 17 | | 
| it could be a trailing button
    GAMEPAD_BUTTON_RIGHT_TRIGGER_1 | 18 | | 
| // Gamepad top/back trigger right (first) | 19 | | 
| it could be a trailing button
    GAMEPAD_BUTTON_RIGHT_TRIGGER_2 | 20 | | 
| // Gamepad top/back trigger right (second) | 21 | | 
| it could be a trailing button
    GAMEPAD_BUTTON_MIDDLE_LEFT | 22 | | 
| // Gamepad center buttons | 23 | | 
| left one (i.e. PS3: Select)
    GAMEPAD_BUTTON_MIDDLE | 24 | | 
| // Gamepad center buttons | 25 | | 
| middle one (i.e. PS3: PS | 26 | | 
| Xbox: XBOX)
    GAMEPAD_BUTTON_MIDDLE_RIGHT | 27 | | 
| // Gamepad center buttons | 28 | | 
| right one (i.e. PS3: Start)
    GAMEPAD_BUTTON_LEFT_THUMB | 29 | | 
| // Gamepad joystick pressed button left
    GAMEPAD_BUTTON_RIGHT_THUMB          // Gamepad joystick pressed button right | 30 | | 

### GamepadAxis

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| GAMEPAD_AXIS_LEFT_X | 0 | | 
| // Gamepad left stick X axis
    GAMEPAD_AXIS_LEFT_Y | 1 | | 
| // Gamepad left stick Y axis
    GAMEPAD_AXIS_RIGHT_X | 2 | | 
| // Gamepad right stick X axis
    GAMEPAD_AXIS_RIGHT_Y | 3 | | 
| // Gamepad right stick Y axis
    GAMEPAD_AXIS_LEFT_TRIGGER | 4 | | 
| // Gamepad back trigger left | 5 | | 
| pressure level: [1..-1]
    GAMEPAD_AXIS_RIGHT_TRIGGER | 5 | | 
| pressure level: [1..-1] | 7 | | 

### MaterialMapIndex

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| MATERIAL_MAP_ALBEDO | 0 | | 
| // Albedo material (same as: MATERIAL_MAP_DIFFUSE)
    MATERIAL_MAP_METALNESS | 1 | | 
| // Metalness material (same as: MATERIAL_MAP_SPECULAR)
    MATERIAL_MAP_NORMAL | 2 | | 
| // Normal material
    MATERIAL_MAP_ROUGHNESS | 3 | | 
| // Roughness material
    MATERIAL_MAP_OCCLUSION | 4 | | 
| // Ambient occlusion material
    MATERIAL_MAP_EMISSION | 5 | | 
| // Emission material
    MATERIAL_MAP_HEIGHT | 6 | | 
| // Heightmap material
    MATERIAL_MAP_CUBEMAP | 7 | | 
| // Cubemap material (NOTE: Uses GL_TEXTURE_CUBE_MAP)
    MATERIAL_MAP_IRRADIANCE | 8 | | 
| // Irradiance material (NOTE: Uses GL_TEXTURE_CUBE_MAP)
    MATERIAL_MAP_PREFILTER | 9 | | 
| // Prefilter material (NOTE: Uses GL_TEXTURE_CUBE_MAP)
    MATERIAL_MAP_BRDF               // Brdf material | 10 | | 

### ShaderLocationIndex

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| SHADER_LOC_VERTEX_POSITION | 0 | | 
| // Shader location: vertex attribute: position
    SHADER_LOC_VERTEX_TEXCOORD01 | 1 | | 
| // Shader location: vertex attribute: texcoord01
    SHADER_LOC_VERTEX_TEXCOORD02 | 2 | | 
| // Shader location: vertex attribute: texcoord02
    SHADER_LOC_VERTEX_NORMAL | 3 | | 
| // Shader location: vertex attribute: normal
    SHADER_LOC_VERTEX_TANGENT | 4 | | 
| // Shader location: vertex attribute: tangent
    SHADER_LOC_VERTEX_COLOR | 5 | | 
| // Shader location: vertex attribute: color
    SHADER_LOC_MATRIX_MVP | 6 | | 
| // Shader location: matrix uniform: model-view-projection
    SHADER_LOC_MATRIX_VIEW | 7 | | 
| // Shader location: matrix uniform: view (camera transform)
    SHADER_LOC_MATRIX_PROJECTION | 8 | | 
| // Shader location: matrix uniform: projection
    SHADER_LOC_MATRIX_MODEL | 9 | | 
| // Shader location: matrix uniform: model (transform)
    SHADER_LOC_MATRIX_NORMAL | 10 | | 
| // Shader location: matrix uniform: normal
    SHADER_LOC_VECTOR_VIEW | 11 | | 
| // Shader location: vector uniform: view
    SHADER_LOC_COLOR_DIFFUSE | 12 | | 
| // Shader location: vector uniform: diffuse color
    SHADER_LOC_COLOR_SPECULAR | 13 | | 
| // Shader location: vector uniform: specular color
    SHADER_LOC_COLOR_AMBIENT | 14 | | 
| // Shader location: vector uniform: ambient color
    SHADER_LOC_MAP_ALBEDO | 15 | | 
| // Shader location: sampler2d texture: albedo (same as: SHADER_LOC_MAP_DIFFUSE)
    SHADER_LOC_MAP_METALNESS | 16 | | 
| // Shader location: sampler2d texture: metalness (same as: SHADER_LOC_MAP_SPECULAR)
    SHADER_LOC_MAP_NORMAL | 17 | | 
| // Shader location: sampler2d texture: normal
    SHADER_LOC_MAP_ROUGHNESS | 18 | | 
| // Shader location: sampler2d texture: roughness
    SHADER_LOC_MAP_OCCLUSION | 19 | | 
| // Shader location: sampler2d texture: occlusion
    SHADER_LOC_MAP_EMISSION | 20 | | 
| // Shader location: sampler2d texture: emission
    SHADER_LOC_MAP_HEIGHT | 21 | | 
| // Shader location: sampler2d texture: heightmap
    SHADER_LOC_MAP_CUBEMAP | 22 | | 
| // Shader location: samplerCube texture: cubemap
    SHADER_LOC_MAP_IRRADIANCE | 23 | | 
| // Shader location: samplerCube texture: irradiance
    SHADER_LOC_MAP_PREFILTER | 24 | | 
| // Shader location: samplerCube texture: prefilter
    SHADER_LOC_MAP_BRDF | 25 | | 
| // Shader location: sampler2d texture: brdf
    SHADER_LOC_VERTEX_BONEIDS | 26 | | 
| // Shader location: vertex attribute: bone indices
    SHADER_LOC_VERTEX_BONEWEIGHTS | 27 | | 
| // Shader location: vertex attribute: bone weights
    SHADER_LOC_MATRIX_BONETRANSFORMS | 28 | | 
| // Shader location: matrix attribute: bone transforms (animation)
    SHADER_LOC_VERTEX_INSTANCETRANSFORM // Shader location: vertex attribute: instance transforms | 29 | | 

### ShaderUniformDataType

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| SHADER_UNIFORM_FLOAT | 0 | | 
| // Shader uniform type: float
    SHADER_UNIFORM_VEC2 | 1 | | 
| // Shader uniform type: vec2 (2 float)
    SHADER_UNIFORM_VEC3 | 2 | | 
| // Shader uniform type: vec3 (3 float)
    SHADER_UNIFORM_VEC4 | 3 | | 
| // Shader uniform type: vec4 (4 float)
    SHADER_UNIFORM_INT | 4 | | 
| // Shader uniform type: int
    SHADER_UNIFORM_IVEC2 | 5 | | 
| // Shader uniform type: ivec2 (2 int)
    SHADER_UNIFORM_IVEC3 | 6 | | 
| // Shader uniform type: ivec3 (3 int)
    SHADER_UNIFORM_IVEC4 | 7 | | 
| // Shader uniform type: ivec4 (4 int)
    SHADER_UNIFORM_UINT | 8 | | 
| // Shader uniform type: unsigned int
    SHADER_UNIFORM_UIVEC2 | 9 | | 
| // Shader uniform type: uivec2 (2 unsigned int)
    SHADER_UNIFORM_UIVEC3 | 10 | | 
| // Shader uniform type: uivec3 (3 unsigned int)
    SHADER_UNIFORM_UIVEC4 | 11 | | 
| // Shader uniform type: uivec4 (4 unsigned int)
    SHADER_UNIFORM_SAMPLER2D        // Shader uniform type: sampler2d | 12 | | 

### ShaderAttributeDataType

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| SHADER_ATTRIB_FLOAT | 0 | | 
| // Shader attribute type: float
    SHADER_ATTRIB_VEC2 | 1 | | 
| // Shader attribute type: vec2 (2 float)
    SHADER_ATTRIB_VEC3 | 2 | | 
| // Shader attribute type: vec3 (3 float)
    SHADER_ATTRIB_VEC4              // Shader attribute type: vec4 (4 float) | 3 | | 

### PixelFormat

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| PIXELFORMAT_UNCOMPRESSED_GRAYSCALE | 1 | | 
| // 8 bit per pixel (no alpha)
    PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA | 1 | | 
| // 8*2 bpp (2 channels)
    PIXELFORMAT_UNCOMPRESSED_R5G6B5 | 2 | | 
| // 16 bpp
    PIXELFORMAT_UNCOMPRESSED_R8G8B8 | 3 | | 
| // 24 bpp
    PIXELFORMAT_UNCOMPRESSED_R5G5B5A1 | 4 | | 
| // 16 bpp (1 bit alpha)
    PIXELFORMAT_UNCOMPRESSED_R4G4B4A4 | 5 | | 
| // 16 bpp (4 bit alpha)
    PIXELFORMAT_UNCOMPRESSED_R8G8B8A8 | 6 | | 
| // 32 bpp
    PIXELFORMAT_UNCOMPRESSED_R32 | 7 | | 
| // 32 bpp (1 channel - float)
    PIXELFORMAT_UNCOMPRESSED_R32G32B32 | 8 | | 
| // 32*3 bpp (3 channels - float)
    PIXELFORMAT_UNCOMPRESSED_R32G32B32A32 | 9 | | 
| // 32*4 bpp (4 channels - float)
    PIXELFORMAT_UNCOMPRESSED_R16 | 10 | | 
| // 16 bpp (1 channel - half float)
    PIXELFORMAT_UNCOMPRESSED_R16G16B16 | 11 | | 
| // 16*3 bpp (3 channels - half float)
    PIXELFORMAT_UNCOMPRESSED_R16G16B16A16 | 12 | | 
| // 16*4 bpp (4 channels - half float)
    PIXELFORMAT_COMPRESSED_DXT1_RGB | 13 | | 
| // 4 bpp (no alpha)
    PIXELFORMAT_COMPRESSED_DXT1_RGBA | 14 | | 
| // 4 bpp (1 bit alpha)
    PIXELFORMAT_COMPRESSED_DXT3_RGBA | 15 | | 
| // 8 bpp
    PIXELFORMAT_COMPRESSED_DXT5_RGBA | 16 | | 
| // 8 bpp
    PIXELFORMAT_COMPRESSED_ETC1_RGB | 17 | | 
| // 4 bpp
    PIXELFORMAT_COMPRESSED_ETC2_RGB | 18 | | 
| // 4 bpp
    PIXELFORMAT_COMPRESSED_ETC2_EAC_RGBA | 19 | | 
| // 8 bpp
    PIXELFORMAT_COMPRESSED_PVRT_RGB | 20 | | 
| // 4 bpp
    PIXELFORMAT_COMPRESSED_PVRT_RGBA | 21 | | 
| // 4 bpp
    PIXELFORMAT_COMPRESSED_ASTC_4x4_RGBA | 22 | | 
| // 8 bpp
    PIXELFORMAT_COMPRESSED_ASTC_8x8_RGBA    // 2 bpp | 23 | | 

### TextureFilter

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| TEXTURE_FILTER_POINT | 0 | | 
| // No filter | 1 | | 
| pixel approximation
    TEXTURE_FILTER_BILINEAR | 2 | | 
| // Linear filtering
    TEXTURE_FILTER_TRILINEAR | 3 | | 
| // Trilinear filtering (linear with mipmaps)
    TEXTURE_FILTER_ANISOTROPIC_4X | 4 | | 
| // Anisotropic filtering 4x
    TEXTURE_FILTER_ANISOTROPIC_8X | 5 | | 
| // Anisotropic filtering 8x
    TEXTURE_FILTER_ANISOTROPIC_16X | 6 | | 
| // Anisotropic filtering 16x | 7 | | 

### TextureWrap

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| TEXTURE_WRAP_REPEAT | 0 | | 
| // Repeats texture in tiled mode
    TEXTURE_WRAP_CLAMP | 1 | | 
| // Clamps texture to edge pixel in tiled mode
    TEXTURE_WRAP_MIRROR_REPEAT | 2 | | 
| // Mirrors and repeats the texture in tiled mode
    TEXTURE_WRAP_MIRROR_CLAMP               // Mirrors and clamps to border the texture in tiled mode | 3 | | 

### CubemapLayout

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| CUBEMAP_LAYOUT_AUTO_DETECT | 0 | | 
| // Automatically detect layout type
    CUBEMAP_LAYOUT_LINE_VERTICAL | 1 | | 
| // Layout is defined by a vertical line with faces
    CUBEMAP_LAYOUT_LINE_HORIZONTAL | 2 | | 
| // Layout is defined by a horizontal line with faces
    CUBEMAP_LAYOUT_CROSS_THREE_BY_FOUR | 3 | | 
| // Layout is defined by a 3x4 cross with cubemap faces
    CUBEMAP_LAYOUT_CROSS_FOUR_BY_THREE     // Layout is defined by a 4x3 cross with cubemap faces | 4 | | 

### FontType

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| FONT_DEFAULT | 0 | | 
| // Default font generation | 1 | | 
| anti-aliased
    FONT_BITMAP | 2 | | 
| // Bitmap font generation | 3 | | 
| no anti-aliasing
    FONT_SDF                        // SDF font generation | 4 | | 
| requires external shader | 5 | | 

### BlendMode

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| BLEND_ALPHA | 0 | | 
| // Blend textures considering alpha (default)
    BLEND_ADDITIVE | 1 | | 
| // Blend textures adding colors
    BLEND_MULTIPLIED | 2 | | 
| // Blend textures multiplying colors
    BLEND_ADD_COLORS | 3 | | 
| // Blend textures adding colors (alternative)
    BLEND_SUBTRACT_COLORS | 4 | | 
| // Blend textures subtracting colors (alternative)
    BLEND_ALPHA_PREMULTIPLY | 5 | | 
| // Blend premultiplied textures considering alpha
    BLEND_CUSTOM | 6 | | 
| // Blend textures using custom src/dst factors (use rlSetBlendFactors())
    BLEND_CUSTOM_SEPARATE           // Blend textures using custom rgb/alpha separate src/dst factors (use rlSetBlendFactorsSeparate()) | 7 | | 

### Gesture

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| GESTURE_NONE | 0 | | 
| // No gesture
    GESTURE_TAP | 1 | | 
| // Tap gesture
    GESTURE_DOUBLETAP | 2 | | 
| // Double tap gesture
    GESTURE_HOLD | 4 | | 
| // Hold gesture
    GESTURE_DRAG | 8 | | 
| // Drag gesture
    GESTURE_SWIPE_RIGHT | 16 | | 
| // Swipe right gesture
    GESTURE_SWIPE_LEFT | 32 | | 
| // Swipe left gesture
    GESTURE_SWIPE_UP | 64 | | 
| // Swipe up gesture
    GESTURE_SWIPE_DOWN | 128 | | 
| // Swipe down gesture
    GESTURE_PINCH_IN | 256 | | 
| // Pinch in gesture
    GESTURE_PINCH_OUT | 512 | | 

### CameraMode

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| CAMERA_CUSTOM | 0 | | 
| // Camera custom | 1 | | 
| controlled by user (UpdateCamera() does nothing)
    CAMERA_FREE | 2 | | 
| // Camera free mode
    CAMERA_ORBITAL | 3 | | 
| // Camera orbital | 4 | | 
| around target | 5 | | 
| zoom supported
    CAMERA_FIRST_PERSON | 6 | | 
| // Camera first person
    CAMERA_THIRD_PERSON             // Camera third person | 7 | | 

### CameraProjection

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| CAMERA_PERSPECTIVE | 0 | | 
| // Perspective projection
    CAMERA_ORTHOGRAPHIC             // Orthographic projection | 1 | | 

### NPatchLayout

| Value | Numeric Value | Description |
|-------|---------------|-------------|
| NPATCH_NINE_PATCH | 0 | | 
| // Npatch layout: 3x3 tiles
    NPATCH_THREE_PATCH_VERTICAL | 1 | | 
| // Npatch layout: 1x3 tiles
    NPATCH_THREE_PATCH_HORIZONTAL   // Npatch layout: 3x1 tiles | 2 | | 


## Macros

### RAYLIB_H

**Definition:**
```c
#define RAYLIB_H 
```

### RAYLIB_VERSION_MAJOR

**Definition:**
```c
#define RAYLIB_VERSION_MAJOR 
```

### RAYLIB_VERSION_MINOR

**Definition:**
```c
#define RAYLIB_VERSION_MINOR 
```

### RAYLIB_VERSION_PATCH

**Definition:**
```c
#define RAYLIB_VERSION_PATCH 
```

### RAYLIB_VERSION

**Definition:**
```c
#define RAYLIB_VERSION 
```

### RL_COLOR_TYPE

**Definition:**
```c
#define RL_COLOR_TYPE 
```

### RL_VECTOR2_TYPE

**Definition:**
```c
#define RL_VECTOR2_TYPE 
```

### RL_VECTOR4_TYPE

**Definition:**
```c
#define RL_VECTOR4_TYPE 
```

### RL_MATRIX_TYPE

**Definition:**
```c
#define RL_MATRIX_TYPE 
```

### LIGHTGRAY

**Definition:**
```c
#define LIGHTGRAY 
```

### GRAY

**Definition:**
```c
#define GRAY 
```

### DARKGRAY

**Definition:**
```c
#define DARKGRAY 
```

### YELLOW

**Definition:**
```c
#define YELLOW 
```

### GOLD

**Definition:**
```c
#define GOLD 
```

### ORANGE

**Definition:**
```c
#define ORANGE 
```

### PINK

**Definition:**
```c
#define PINK 
```

### RED

**Definition:**
```c
#define RED 
```

### MAROON

**Definition:**
```c
#define MAROON 
```

### GREEN

**Definition:**
```c
#define GREEN 
```

### LIME

**Definition:**
```c
#define LIME 
```

### DARKGREEN

**Definition:**
```c
#define DARKGREEN 
```

### SKYBLUE

**Definition:**
```c
#define SKYBLUE 
```

### BLUE

**Definition:**
```c
#define BLUE 
```

### DARKBLUE

**Definition:**
```c
#define DARKBLUE 
```

### PURPLE

**Definition:**
```c
#define PURPLE 
```

### VIOLET

**Definition:**
```c
#define VIOLET 
```

### DARKPURPLE

**Definition:**
```c
#define DARKPURPLE 
```

### BEIGE

**Definition:**
```c
#define BEIGE 
```

### BROWN

**Definition:**
```c
#define BROWN 
```

### DARKBROWN

**Definition:**
```c
#define DARKBROWN 
```

### WHITE

**Definition:**
```c
#define WHITE 
```

### BLACK

**Definition:**
```c
#define BLACK 
```

### BLANK

**Definition:**
```c
#define BLANK 
```

### MAGENTA

**Definition:**
```c
#define MAGENTA 
```

### RAYWHITE

**Definition:**
```c
#define RAYWHITE 
```

### MOUSE_LEFT_BUTTON

**Definition:**
```c
#define MOUSE_LEFT_BUTTON 
```

### MOUSE_RIGHT_BUTTON

**Definition:**
```c
#define MOUSE_RIGHT_BUTTON 
```

### MOUSE_MIDDLE_BUTTON

**Definition:**
```c
#define MOUSE_MIDDLE_BUTTON 
```

### MATERIAL_MAP_DIFFUSE

**Definition:**
```c
#define MATERIAL_MAP_DIFFUSE 
```

### MATERIAL_MAP_SPECULAR

**Definition:**
```c
#define MATERIAL_MAP_SPECULAR 
```

### SHADER_LOC_MAP_DIFFUSE

**Definition:**
```c
#define SHADER_LOC_MAP_DIFFUSE 
```

### SHADER_LOC_MAP_SPECULAR

**Definition:**
```c
#define SHADER_LOC_MAP_SPECULAR 
```

### GetMouseRay

**Definition:**
```c
#define GetMouseRay 
```

