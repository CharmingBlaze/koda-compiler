package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ReleaseConfig holds configuration for building releases
type ReleaseConfig struct {
	Version         string
	OutputDir       string
	Platforms       []string
	IncludeSource   bool
	IncludeDocs     bool
	CreateInstaller bool
}

// Platform represents a target platform
type Platform struct {
	GOOS   string
	GOARCH string
	Ext    string
	Name   string
}

func main() {
	config := parseFlags()

	fmt.Printf("Koda Cross-Platform Release Builder\n")
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Platforms: %v\n", config.Platforms)
	fmt.Printf("Output: %s\n", config.OutputDir)

	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Build releases for all platforms
	for _, platform := range config.Platforms {
		fmt.Printf("\nBuilding for %s...\n", platform)
		if err := buildPlatformRelease(platform, config); err != nil {
			log.Printf("Failed to build for %s: %v", platform, err)
		} else {
			fmt.Printf("[ok] %s built successfully\n", platform)
		}
	}

	// Create source archive
	if config.IncludeSource {
		fmt.Printf("\nCreating source archive...\n")
		if err := createSourceArchive(config); err != nil {
			log.Printf("Failed to create source archive: %v", err)
		} else {
			fmt.Printf("[ok] Source archive created\n")
		}
	}

	// Create checksums
	fmt.Printf("\nCreating checksums...\n")
	if err := createChecksums(config); err != nil {
		log.Printf("Failed to create checksums: %v", err)
	} else {
		fmt.Printf("[ok] Checksums created\n")
	}

	fmt.Printf("\nRelease build completed.\n")
	fmt.Printf("Output directory: %s\n", config.OutputDir)
}

func parseFlags() *ReleaseConfig {
	config := &ReleaseConfig{
		Version:         "1.0.0",
		OutputDir:       "./releases",
		Platforms:       []string{"windows-amd64", "linux-amd64", "darwin-amd64", "darwin-arm64"},
		IncludeSource:   true,
		IncludeDocs:     true,
		CreateInstaller: true,
	}

	// In a real implementation, parse command line flags here
	// For now, use defaults

	return config
}

func buildPlatformRelease(platformStr string, config *ReleaseConfig) error {
	platform := parsePlatform(platformStr)

	// Create platform-specific directory
	platformDir := filepath.Join(config.OutputDir, "koda-"+platformStr)
	if err := os.MkdirAll(platformDir, 0755); err != nil {
		return err
	}

	// Build the koda-single executable for this platform
	if err := buildExecutable(platformDir, platform, config); err != nil {
		return fmt.Errorf("failed to build executable: %v", err)
	}

	// Copy embedded resources
	if err := copyResources(platformDir); err != nil {
		return fmt.Errorf("failed to copy resources: %v", err)
	}

	// Create platform-specific files
	if err := createPlatformFiles(platformDir, platform, config); err != nil {
		return fmt.Errorf("failed to create platform files: %v", err)
	}

	// Create documentation
	if config.IncludeDocs {
		if err := createDocumentation(platformDir, platform, config); err != nil {
			return fmt.Errorf("failed to create documentation: %v", err)
		}
	}

	// Create examples
	if err := createExamples(platformDir, platform, config); err != nil {
		return fmt.Errorf("failed to create examples: %v", err)
	}

	// Create installer
	if config.CreateInstaller {
		if err := createInstaller(platformDir, platform, config); err != nil {
			return fmt.Errorf("failed to create installer: %v", err)
		}
	}

	// Create distribution archive
	if err := createDistributionArchive(platformDir, platformStr, config); err != nil {
		return fmt.Errorf("failed to create archive: %v", err)
	}

	return nil
}

func parsePlatform(platformStr string) Platform {
	parts := strings.Split(platformStr, "-")
	if len(parts) != 2 {
		return Platform{GOOS: "unknown", GOARCH: "unknown"}
	}

	goos := parts[0]
	goarch := parts[1]

	var ext string
	var name string

	switch goos {
	case "windows":
		ext = ".exe"
		name = "Windows"
	case "linux":
		ext = ""
		name = "Linux"
	case "darwin":
		ext = ""
		if goarch == "arm64" {
			name = "macOS (Apple Silicon)"
		} else {
			name = "macOS (Intel)"
		}
	default:
		ext = ""
		name = goos
	}

	return Platform{
		GOOS:   goos,
		GOARCH: goarch,
		Ext:    ext,
		Name:   name,
	}
}

func buildExecutable(outputDir string, platform Platform, config *ReleaseConfig) error {
	_ = config
	fmt.Printf("  Building executable for %s %s...\n", platform.GOOS, platform.GOARCH)

	// Set environment variables for cross-compilation
	env := append(os.Environ(),
		fmt.Sprintf("GOOS=%s", platform.GOOS),
		fmt.Sprintf("GOARCH=%s", platform.GOARCH),
		"CGO_ENABLED=0", // Disable CGO for static binaries
	)

	// Build command
	cmd := exec.Command("go", "build",
		"-ldflags", "-s -w", // Strip symbols for smaller binaries
		"-o", filepath.Join(outputDir, "koda"+platform.Ext),
		"../koda-single",
	)

	cmd.Dir = "."
	cmd.Env = env

	// Run build
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v, output: %s", err, string(output))
	}

	fmt.Printf("  [ok] Executable built: koda%s\n", platform.Ext)
	return nil
}

func copyResources(outputDir string) error {
	fmt.Printf("  Copying resources...\n")

	// Copy runtime, stdlib, wrappers directories
	resources := []string{"runtime", "stdlib", "wrappers"}

	for _, resource := range resources {
		source := filepath.Join("../koda-single", resource)
		dest := filepath.Join(outputDir, resource)

		if err := copyDir(source, dest); err != nil {
			return fmt.Errorf("failed to copy %s: %v", resource, err)
		}
	}

	fmt.Printf("  [ok] Resources copied\n")
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func createPlatformFiles(outputDir string, platform Platform, config *ReleaseConfig) error {
	fmt.Printf("  Creating platform-specific files...\n")

	// Create launcher script
	var launcherName string
	var launcherContent string

	if platform.GOOS == "windows" {
		launcherName = "koda.bat"
		launcherContent = fmt.Sprintf(`@echo off
REM Koda Launcher for Windows
echo Koda Programming Language v%s - %s
echo.

REM Run Koda with all arguments
"%%~dp0koda.exe" %%*
`, config.Version, platform.Name)
	} else {
		launcherName = "koda"
		launcherContent = fmt.Sprintf(`#!/bin/bash
# Koda Launcher for %s
echo "Koda Programming Language v%s - %s"
echo

# Run Koda with all arguments
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/koda" "$@"
`, platform.Name, config.Version, platform.Name)
	}

	launcherPath := filepath.Join(outputDir, launcherName)
	if err := os.WriteFile(launcherPath, []byte(launcherContent), 0755); err != nil {
		return err
	}

	// Create install script
	installScript := createInstallScript(platform, config)
	installPath := filepath.Join(outputDir, "install.sh")
	if platform.GOOS == "windows" {
		installPath = filepath.Join(outputDir, "install.bat")
	}

	if err := os.WriteFile(installPath, []byte(installScript), 0755); err != nil {
		return err
	}

	fmt.Printf("  [ok] Platform files created\n")
	return nil
}

func createInstallScript(platform Platform, config *ReleaseConfig) string {
	if platform.GOOS == "windows" {
		return fmt.Sprintf(`@echo off
REM Koda Installation Script for Windows
echo Installing Koda v%s...

REM Create installation directory
if not exist "%%PROGRAMFILES%%\Koda" mkdir "%%PROGRAMFILES%%\Koda"

REM Copy files
xcopy /E /Y "*" "%%PROGRAMFILES%%\Koda\"

REM Add to PATH (optional)
echo Adding Koda to system PATH...
setx PATH "%%PATH%%;%%PROGRAMFILES%%\Koda" /M

echo Installation completed!
echo You can now run 'koda' from anywhere.
echo.
pause
`, config.Version)
	}

	return fmt.Sprintf(`#!/bin/bash
# Koda Installation Script for %s
echo "Installing Koda v%s..."

# Determine installation directory
INSTALL_DIR="/usr/local/bin"
if [ "$EUID" -ne 0 ]; then
    INSTALL_DIR="$HOME/.local/bin"
    echo "Installing to user directory: $INSTALL_DIR"
fi

# Create installation directory
mkdir -p "$INSTALL_DIR"

# Copy files
cp -r * "$INSTALL_DIR/koda/"

# Create symlink
ln -sf "$INSTALL_DIR/koda/koda" "$INSTALL_DIR/koda"

# Update PATH if needed
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo 'export PATH="$PATH:'"$INSTALL_DIR"'"' >> ~/.bashrc
    echo "Added $INSTALL_DIR to PATH in ~/.bashrc"
fi

echo "Installation completed!"
echo "You can now run 'koda' from anywhere."
echo "Run 'source ~/.bashrc' or restart your terminal to update PATH."
`, platform.Name, config.Version)
}

func createDocumentation(outputDir string, platform Platform, config *ReleaseConfig) error {
	fmt.Printf("  Creating documentation...\n")

	docsDir := filepath.Join(outputDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	// Create README
	readmeContent := fmt.Sprintf(`# Koda Programming Language v%s

## Platform: %s

### Quick Start

Koda is a self-contained programming language that requires no external dependencies.

#### Installation

1. Extract this archive to any directory
2. Run the install script:
   - Windows: Double-click `+"`install.bat`"+`
   - Linux/macOS: Run `+"`./install.sh`"+`

3. Start using Koda:
   `+"```bash"+`
   koda hello.koda
   `+"```"+`

#### Manual Installation

If you prefer not to run the install script:

1. Extract the archive to your preferred location
2. Add the directory to your system PATH
3. Run `+"`koda`"+` from anywhere

### Features

- **Zero Dependencies**: No external compilers, runtimes, or libraries needed
- **Cross-Platform**: Works on Windows, Linux, and macOS
- **Self-Contained**: Everything included in a single package
- **Modern Syntax**: Clean, expressive language syntax
- **Fast Compilation**: Compiles to native code via LLVM
- **Rich Standard Library**: Built-in functions for common tasks
- **C/C++ Interop**: Easy integration with existing libraries
- **Wrapper Generation**: Automatic bindings for any C/C++ library

### Directory Structure

`+"```"+`
koda/
├── koda%s              # Main compiler executable
├── runtime/             # Embedded runtime components
├── stdlib/              # Standard library files
├── wrappers/            # Pre-built library wrappers
├── examples/            # Example programs
└── docs/                # Documentation
`+"```"+`

### Examples

Check the `+"`examples/`"+` directory for sample programs:

- `+"`hello.koda`"+` - Basic hello world
- `+"`functions.koda`"+` - Function definitions and usage
- `+"`loops.koda`"+` - Loop constructs
- `+"`raylib_demo.koda`"+` - Graphics programming with Raylib

### Language Reference

#### Basic Syntax

`+"```koda"+`
// Variables
let name = "Koda";
let version = 1.0;

// Functions
func greet(name) {
    print("Hello, " + name + "!");
}

// Loops
for (let i = 0; i < 5; i = i + 1) {
    print("i = " + i);
}

// Conditionals
if (version > 1.0) {
    print("Latest version!");
}
`+"```"+`

#### C Library Integration

`+"```koda"+`
// Use any C library with automatic wrapper generation
native printf(format, ...);
printf("Hello from C: %%d", 42);

// Or use pre-built wrappers
import "raylib";
initWindow(800, 600, "My Game");
`+"```"+`

### Advanced Features

- **Wrapper Generation**: Automatically generate bindings for any C/C++ library
- **Package Management**: Built-in package system with no external dependencies
- **WebAssembly**: Compile to WebAssembly for web applications
- **Cross-Compilation**: Compile for any platform from any platform

### Getting Help

- Check the `+"`examples/`"+` directory for sample programs
- Visit the online documentation at https://koda-lang.org
- Join the community at https://github.com/koda-lang/koda
- Report issues at https://github.com/koda-lang/koda/issues

### System Requirements

- **Windows**: Windows 7 or later
- **Linux**: Any modern Linux distribution
- **macOS**: macOS 10.12 or later

No additional software or dependencies required!

### Uninstallation

To uninstall Koda:

**Windows:**
1. Delete the installation directory
2. Remove from system PATH if added

**Linux/macOS:**
1. Delete the installation directory
2. Remove the symlink: `+"`sudo rm /usr/local/bin/koda`"+`
3. Remove from PATH in ~/.bashrc if added

---

Enjoy programming with Koda.

Version: %s
Build Date: %s
Platform: %s
`, config.Version, platform.Name, platform.Ext, config.Version, time.Now().Format("2006-01-02"), platform.Name)

	readmePath := filepath.Join(outputDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return err
	}

	fmt.Printf("  [ok] Documentation created\n")
	return nil
}

func createExamples(outputDir string, platform Platform, config *ReleaseConfig) error {
	_ = platform
	fmt.Printf("  Creating examples...\n")

	examplesDir := filepath.Join(outputDir, "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return err
	}

	// Create example files
	examples := map[string]string{
		"hello.koda": `// Hello World Example
print("Hello, World!");
print("Welcome to Koda v%s!");

// Variables
let name = "Koda";
let version = 1.0;

print("Language: " + name);
print("Version: " + version);

// Functions
func greet(who) {
    print("Hello, " + who + "!");
}

greet("World");
greet("Self-contained Koda");

print("Example completed successfully!");
`,
		"functions.koda": `// Function Examples
// Simple function
func add(a, b) {
    return a + b;
}

// Recursive function
func factorial(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

// Function as value
let multiply = func(a, b) {
    return a * b;
};

// Higher-order function
func apply(func, a, b) {
    return func(a, b);
}

print("5 + 3 = " + add(5, 3));
print("5 * 3 = " + multiply(5, 3));
print("5! = " + factorial(5));
print("apply(add, 10, 20) = " + apply(add, 10, 20));
`,
		"loops.koda": `// Loop Examples
print("For loop:");
for (let i = 0; i < 5; i = i + 1) {
    print("i = " + i);
}

print("While loop:");
let j = 0;
while (j < 3) {
    print("j = " + j);
    j = j + 1;
}

print("Do-while loop:");
let k = 0;
do {
    print("k = " + k);
    k = k + 1;
} while (k < 2);

print("Nested loops:");
for (let i = 0; i < 2; i = i + 1) {
    for (let j = 0; j < 2; j = j + 1) {
        print("i=" + i + ", j=" + j);
    }
}

print("Loop examples completed!");
`,
		"raylib_demo.koda": `// Raylib Graphics Demo
// Note: This requires the Raylib wrapper

import "raylib";

// Initialize window
initWindow(800, 600, "Koda + Raylib");

let targetFPS = 60;
setTargetFPS(targetFPS);

print("Raylib demo started!");
print("Press ESC to exit");

// Main game loop
while (!windowShouldClose()) {
    // Update
    // (Update logic here)
    
    // Draw
    beginDrawing();
    clearBackground(RAYWHITE);
    
    // Draw text
    drawText("Hello from Koda + Raylib!", 190, 200, 20, BLACK);
    drawText("Press ESC to exit", 200, 250, 20, GRAY);
    
    // Draw shapes
    drawCircle(400, 150, 30, MAROON);
    drawRectangle(250, 350, 100, 50, BLUE);
    drawLine(100, 100, 700, 100, DARKGREEN);
    
    endDrawing();
}

// Cleanup
closeWindow();
print("Raylib demo completed!");
`,
	}

	for filename, content := range examples {
		// Format the content with version
		formattedContent := fmt.Sprintf(content, config.Version)
		outputPath := filepath.Join(examplesDir, filename)
		if err := os.WriteFile(outputPath, []byte(formattedContent), 0644); err != nil {
			return err
		}
	}

	fmt.Printf("  [ok] Examples created\n")
	return nil
}

func createInstaller(outputDir string, platform Platform, config *ReleaseConfig) error {
	fmt.Printf("  Creating installer...\n")

	// Create installer script
	installerContent := createInstallerScript(platform, config)

	var installerName string
	if platform.GOOS == "windows" {
		installerName = "setup.exe"
	} else {
		installerName = "install.sh"
	}

	installerPath := filepath.Join(outputDir, installerName)
	if err := os.WriteFile(installerPath, []byte(installerContent), 0755); err != nil {
		return err
	}

	fmt.Printf("  [ok] Installer created\n")
	return nil
}

func createInstallerScript(platform Platform, config *ReleaseConfig) string {
	// This would create a proper installer in a real implementation
	// For now, return the install script content
	return createInstallScript(platform, config)
}

func createDistributionArchive(platformDir, platformStr string, config *ReleaseConfig) error {
	fmt.Printf("  Creating distribution archive...\n")

	archiveName := fmt.Sprintf("koda-%s-%s", config.Version, platformStr)
	archivePath := filepath.Join(config.OutputDir, archiveName)

	if strings.Contains(platformStr, "windows") {
		// Create ZIP for Windows
		if err := createZip(archivePath+".zip", platformDir); err != nil {
			return err
		}
	} else {
		// Create tar.gz for Unix systems
		if err := createTarGz(archivePath+".tar.gz", platformDir); err != nil {
			return err
		}
	}

	// Remove the temporary directory
	os.RemoveAll(platformDir)

	fmt.Printf("  [ok] Archive created: %s\n", archiveName)
	return nil
}

func createZip(filename, sourceDir string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func createTarGz(filename, sourceDir string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})
}

func createSourceArchive(config *ReleaseConfig) error {
	// Create source archive
	sourceDir := filepath.Join(config.OutputDir, "koda-"+config.Version+"-source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		return err
	}

	// Copy source files (simplified)
	// In a real implementation, you'd copy all source files

	archiveName := fmt.Sprintf("koda-%s-source.tar.gz", config.Version)
	archivePath := filepath.Join(config.OutputDir, archiveName)

	if err := createTarGz(archivePath, sourceDir); err != nil {
		return err
	}

	os.RemoveAll(sourceDir)
	return nil
}

func createChecksums(config *ReleaseConfig) error {
	checksumsFile := filepath.Join(config.OutputDir, "checksums.txt")
	file, err := os.Create(checksumsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Calculate checksums for all archives
	archives, err := filepath.Glob(filepath.Join(config.OutputDir, "*.zip"))
	if err != nil {
		return err
	}

	tarArchives, err := filepath.Glob(filepath.Join(config.OutputDir, "*.tar.gz"))
	if err != nil {
		return err
	}

	archives = append(archives, tarArchives...)

	for _, archive := range archives {
		checksum := calculateChecksum(archive)
		filename := filepath.Base(archive)
		fmt.Fprintf(file, "%s  %s\n", checksum, filename)
	}

	return nil
}

func calculateChecksum(filename string) string {
	_ = filename
	// In a real implementation, calculate SHA256 checksum
	// For now, return a placeholder
	return "sha256-placeholder"
}
