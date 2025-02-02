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

	"github.com/canta2899/logo-ls/ctw"
	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/model"
)

// DummyTimeFormatter implements the minimal time formatter needed for tests.
type DummyTimeFormatter struct{}

func (d DummyTimeFormatter) Format(t *time.Time) string {
	// Fixed layout for testing.
	return t.Format("2006-01-02 15:04:05")
}

// dummyIcon is used in place of a real icon when testing PrintDirectory.
type dummyIcon struct{}

func (d dummyIcon) GetColor() string { return "" }
func (d dummyIcon) GetGlyph() string { return "" }

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

	conf := &Config{
		FileList:      []string{nonExistent, tempDir, tempFileName},
		TimeFormatter: DummyTimeFormatter{},
	}

	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	appInstance := &App{
		Config: conf,
		Writer: new(bytes.Buffer),
		Logger: logger,
	}

	args := appInstance.GetArguments()
	if len(args.Files) != 1 {
		t.Errorf("Expected 1 file entry, got %d", len(args.Files))
	}
	if len(args.Dirs) != 1 {
		t.Errorf("Expected 1 directory entry, got %d", len(args.Dirs))
	}

	// Check that a log message was emitted for the non-existent entry.
	if !strings.Contains(logBuf.String(), "cannot access") && !strings.Contains(logBuf.String(), "cannot get absolute path") {
		t.Error("Expected log output for non-existent path, got none")
	}

	// Check that the exit code was flagged (assuming zero means success).
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

	fi, err := os.Stat(tempFileName)
	if err != nil {
		t.Fatalf("Failed to stat temporary file: %v", err)
	}

	fileEntry := model.FileEntry{
		FileInfo: fi,
		AbsPath:  tempFileName,
	}

	conf := &Config{
		LongListingMode: model.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
		DisableIcon:     true,
	}
	appInstance := &App{
		Config: conf,
		Writer: new(bytes.Buffer),
		Logger: log.New(io.Discard, "", 0),
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
	err = os.WriteFile(filePath, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file in temporary directory: %v", err)
	}

	f, err := os.Open(tempDir)
	if err != nil {
		t.Fatalf("Failed to open temporary directory: %v", err)
	}

	// ProcessDirectory will call Close() on the directory.
	dirEntry := &model.DirectoryEntry{
		File:    *f,
		AbsPath: tempDir,
	}

	conf := &Config{
		AllMode:         model.IncludeDefault,
		LongListingMode: model.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
		DisableIcon:     true,
	}
	appInstance := &App{
		Config: conf,
		Writer: new(bytes.Buffer),
		Logger: log.New(io.Discard, "", 0),
	}

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

	fi, err := os.Stat(tempFileName)
	if err != nil {
		t.Fatalf("Failed to stat temporary file: %v", err)
	}

	conf := &Config{
		LongListingMode: model.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
		ShowInodeNumber: false,
		ShowBlockSize:   false,
		DisableIcon:     true,
	}
	appInstance := &App{
		Config: conf,
		Writer: new(bytes.Buffer),
		Logger: log.New(io.Discard, "", 0),
	}

	entry := appInstance.buildEntry(tempFileName, fi, false)
	if entry.Name == "" {
		t.Error("buildEntry returned an empty name")
	}
	// Expect no extension if the file name has none.
	if entry.Ext != "" {
		t.Errorf("Expected empty extension, got %q", entry.Ext)
	}
}

// TestPrintDirectory builds a dummy directory model and then verifies that PrintDirectory writes output.
func TestPrintDirectory(t *testing.T) {
	// Create a dummy file entry.
	dummyEntry := &model.Entry{
		Icon:      &icons.IconInfo{Glyph: "dummy", Color: [3]uint8{0, 0, 0}, IsExecutable: false},
		Name:      "dummy",
		Ext:       ".txt",
		Indicator: "",
		Size:      456,
		ModTime:   time.Now(),
	}
	// Create a dummy directory model.
	dirModel := &model.Directory{
		Files: []*model.Entry{dummyEntry},
	}

	conf := &Config{
		LongListingMode: model.LongListingNone,
		OneFilePerLine:  false,
		DisableIcon:     true,
		HumanReadable:   true,
		TimeFormatter:   DummyTimeFormatter{},
	}
	var buf bytes.Buffer
	appInstance := &App{
		Config: conf,
		Writer: &buf,
		Logger: log.New(io.Discard, "", 0),
	}

	appInstance.PrintDirectory(dirModel)
	if buf.Len() == 0 {
		t.Error("PrintDirectory did not produce any output")
	}
}

// TestBlockSizeWithInode verifies that blockSizeWithInode produces a string containing the inode.
func TestBlockSizeWithInode(t *testing.T) {
	conf := &Config{
		ShowInodeNumber: true,
		ShowBlockSize:   true,
	}
	appInstance := &App{
		Config: conf,
	}

	// Create a dummy entry with an inode number and block count.
	dummyEntry := &model.Entry{
		InodeNumber: "98765",
		Blocks:      16, // Assuming Blocks is an integer type.
	}

	result := appInstance.blockSizeWithInode(dummyEntry)
	if !strings.Contains(result, "98765") {
		t.Errorf("Expected inode number in output, got %q", result)
	}
}

// TestGetCTW verifies that getCTW returns a non-nil instance.
func TestGetCTW(t *testing.T) {
	conf := &Config{
		LongListingMode: model.LongListingNone,
		OneFilePerLine:  false,
		DisableIcon:     false,
	}
	appInstance := &App{
		Config: conf,
	}
	ctwInstance := appInstance.getCTW()
	if ctwInstance == nil {
		t.Error("getCTW returned nil")
	}

	if _, ok := ctwInstance.(ctw.CTW); !ok {
		t.Error("getCTW did not return an instance that implements ctw.CTW")
	}
}
