package app

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/canta2899/logo-ls/fs/osfs"
	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/internal/inspect"
	"github.com/canta2899/logo-ls/model"
)

// DummyTimeFormatter implements the minimal time formatter needed for tests.
type DummyTimeFormatter struct{}

func (d DummyTimeFormatter) Format(t *time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func newTestApp(conf *cli.Config, logger *log.Logger, writer io.Writer) *App {
	return &App{
		Config: conf,
		Writer: writer,
		Logger: logger,
		FS:     osfs.New(),
	}
}

// TestGetArguments creates a temporary file, a temporary directory,
// and passes a non-existent path. It then checks that GetArguments
// correctly categorizes valid entries and logs errors for the bad one.
func TestGetArguments(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFileName)

	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	nonExistent := filepath.Join(os.TempDir(), "nonexistent_path")

	conf := &cli.Config{
		FileList:      []string{nonExistent, tempDir, tempFileName},
		TimeFormatter: DummyTimeFormatter{},
	}

	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)
	appInstance := newTestApp(conf, logger, new(bytes.Buffer))

	args := appInstance.GetArguments()
	if len(args.Files) != 1 {
		t.Errorf("Expected 1 file entry, got %d", len(args.Files))
	}
	if len(args.Dirs) != 1 {
		t.Errorf("Expected 1 directory entry, got %d", len(args.Dirs))
	}

	if !strings.Contains(logBuf.String(), "cannot access") && !strings.Contains(logBuf.String(), "cannot get absolute path") {
		t.Error("Expected log output for non-existent path, got none")
	}

	if int(appInstance.ExitCode) == 0 {
		t.Error("Expected non-zero exit code due to error, got 0")
	}
}

// TestProcessFiles verifies that ProcessFiles returns a directory model with the file entry.
func TestProcessFiles(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFileName)

	conf := &cli.Config{
		LongListingMode: model.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
		DisableIcon:     true,
	}
	appInstance := newTestApp(conf, log.New(io.Discard, "", 0), new(bytes.Buffer))

	fi, err := appInstance.FS.Stat(tempFileName)
	if err != nil {
		t.Fatalf("Failed to stat temporary file: %v", err)
	}

	fileEntry := model.FileEntry{
		Info:    fi,
		AbsPath: tempFileName,
	}

	dirModel := appInstance.ProcessFiles([]model.FileEntry{fileEntry})
	if len(dirModel.Files) != 1 {
		t.Errorf("Expected 1 file in processed directory, got %d", len(dirModel.Files))
	}
}

// TestProcessDirectory creates a temporary directory with one file inside,
// calls ProcessDirectory, and then verifies that the file appears in the output.
func TestProcessDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("Failed to create file in temporary directory: %v", err)
	}

	conf := &cli.Config{
		AllMode:         model.IncludeDefault,
		LongListingMode: model.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
		DisableIcon:     true,
	}
	appInstance := newTestApp(conf, log.New(io.Discard, "", 0), new(bytes.Buffer))

	f, err := appInstance.FS.Open(tempDir)
	if err != nil {
		t.Fatalf("Failed to open temporary directory: %v", err)
	}
	dirEntry := &model.DirectoryEntry{File: f, AbsPath: tempDir}

	dirModel, err := appInstance.ProcessDirectory(dirEntry)
	if err != nil {
		t.Fatalf("ProcessDirectory returned error: %v", err)
	}
	found := false
	for _, entry := range dirModel.Files {
		if strings.Contains(entry.Name, "file") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected file.txt to be listed in directory contents")
	}
}

// TestBuildEntry verifies that buildEntry properly extracts file name and extension.
func TestBuildEntry(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFileName)

	conf := &cli.Config{
		LongListingMode: model.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
		ShowInodeNumber: false,
		ShowBlockSize:   false,
		DisableIcon:     true,
	}
	appInstance := newTestApp(conf, log.New(io.Discard, "", 0), new(bytes.Buffer))

	fi, err := appInstance.FS.Stat(tempFileName)
	if err != nil {
		t.Fatalf("Failed to stat temporary file: %v", err)
	}

	entry := appInstance.buildEntry(tempFileName, fi, false)
	if entry.Name == "" {
		t.Error("buildEntry returned an empty name")
	}
	if entry.Ext != "" {
		t.Errorf("Expected empty extension, got %q", entry.Ext)
	}
}

// TestPrintDirectory builds a dummy directory model and then verifies that PrintDirectory writes output.
func TestPrintDirectory(t *testing.T) {
	dummyEntry := &inspect.InspectedEntry{
		Icon:      &icons.IconInfo{Glyph: "dummy", Color: [3]uint8{0, 0, 0}, IsExecutable: false},
		Name:      "dummy.txt",
		Base:      "dummy",
		Ext:       ".txt",
		Indicator: "",
		Size:      456,
		ModTime:   time.Now(),
	}
	dirModel := &model.Directory{
		Files: []*inspect.InspectedEntry{dummyEntry},
	}

	conf := &cli.Config{
		LongListingMode: model.LongListingNone,
		OneFilePerLine:  false,
		DisableIcon:     true,
		HumanReadable:   true,
		TimeFormatter:   DummyTimeFormatter{},
	}
	var buf bytes.Buffer
	appInstance := newTestApp(conf, log.New(io.Discard, "", 0), &buf)

	appInstance.PrintDirectory(dirModel)
	if buf.Len() == 0 {
		t.Error("PrintDirectory did not produce any output")
	}
}

