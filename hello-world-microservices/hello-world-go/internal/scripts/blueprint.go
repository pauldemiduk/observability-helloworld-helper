package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type ParseContext struct {
	CurrentParent *ItemWithLine
}

// ItemWithLine represents an entity (function, struct, comment, etc.)
// that has a name and an associated line number in the source code.
type ItemWithLine struct {
	Name       string `json:"name" yaml:"name"`               // Name of the item (e.g., function name, struct name, log call)
	LineNumber int    `json:"line_number" yaml:"line_number"` // Line number where the item is found in the source file
}

// FileMetadata stores all extracted metadata for a single Go source file.
type FileMetadata struct {
	FilePath           string                    `json:"file_path" yaml:"file_path"`                     // Absolute path of the Go source file
	Package            string                    `json:"package" yaml:"package"`                         // Package name declared in the file
	InternalFunctions  []ItemWithLine            `json:"internal_functions" yaml:"internal_functions"`   // List of internal (unexported) functions
	ExternalFunctions  []ItemWithLine            `json:"external_functions" yaml:"external_functions"`   // List of external (exported) functions
	CallMap            map[string][]ItemWithLine `json:"call_map" yaml:"call_map"`                       // Tracks function call dependencies (future use)
	ObservabilityHints map[string][]ItemWithLine `json:"observability_hints" yaml:"observability_hints"` // Tracks log, metric, and trace occurrences
	Structs            []ItemWithLine            `json:"structs" yaml:"structs"`                         // List of struct definitions
	Interfaces         []ItemWithLine            `json:"interfaces" yaml:"interfaces"`                   // List of interface definitions
	Imports            []string                  `json:"imports" yaml:"imports"`                         // List of imported packages
	Comments           map[int][]ItemWithLine    `json:"comments" yaml:"comments"`                       // Captured inline and block comments, indexed by line number
}

// BlueprintConfig represents the configuration structure for the blueprint utility.
// It is loaded from a YAML file and controls how the utility scans, analyzes, and outputs metadata.
type BlueprintConfig struct {
	Settings struct {
		ConfigFilePath string `yaml:"config_file_path"` // Path to the main configuration file
		CollectionMode string `yaml:"collection_mode"`  // Mode of operation: "collect_first" (batch processing) or "write_directly" (real-time)
	} `yaml:"settings"`
	DefaultApproach struct {
		// OutputManagement controls how the generated blueprint output is handled
		OutputManagement struct {
			LatestFile      string `yaml:"latest_file"`      // Filename for the latest blueprint output
			TimestampFormat string `yaml:"timestamp_format"` // Format for timestamps used in file versioning
			EnableRotation  bool   `yaml:"enable_rotation"`  // Whether old blueprint versions should be rotated
			RenameExisting  struct {
				Pattern             string `yaml:"pattern"`               // Filename pattern for rotated files (e.g., "code-blueprint-[{timestamp}].txt")
				PreserveOldVersions bool   `yaml:"preserve_old_versions"` // If true, retains older versions instead of overwriting
			} `yaml:"rename_existing"`
			LatestSymlink bool `yaml:"latest_symlink"` // Whether to create a symlink pointing to the latest output file
			MaxHistory    int  `yaml:"max_history"`    // Maximum number of historical blueprint files to retain
		} `yaml:"output_management"`

		// CodeScanToggles enable or disable various types of code structure analysis.
		CodeScanToggles struct {
			IncludeFunctions         bool `yaml:"include_functions"`          // Extract both internal and external functions
			IncludeStructs           bool `yaml:"include_structs"`            // Extract struct definitions
			IncludeInterfaces        bool `yaml:"include_interfaces"`         // Extract interface definitions
			IncludeImports           bool `yaml:"include_imports"`            // Extract imported packages
			IncludeComments          bool `yaml:"include_comments"`           // Capture inline and block comments
			IncludeInternalFunctions bool `yaml:"include_internal_functions"` // More granular function filtering
			IncludeExternalFunctions bool `yaml:"include_external_functions"` // More granular function filtering
			IncludeCallMap           bool `yaml:"include_call_map"`           // Future use for call dependency mapping
			IncludeObservability     bool `yaml:"include_observability"`      // Enable scanning for observability patterns (logs, metrics, traces)
		} `yaml:"code_scan_toggles"`

		// ObservabilityScanToggles enable or disable different observability signals.
		ObservabilityScanToggles struct {
			Logs      bool `yaml:"logs"`      // Capture log statements (e.g., log.Println, logrus)
			Metrics   bool `yaml:"metrics"`   // Capture metric collection (e.g., Prometheus, StatsD)
			Traces    bool `yaml:"traces"`    // Capture distributed tracing (e.g., OpenTelemetry, Jaeger)
			Debug     bool `yaml:"debug"`     // Capture debug-level log statements
			Info      bool `yaml:"info"`      // Capture info-level log statements
			Error     bool `yaml:"error"`     // Capture error-level log statements
			Fatal     bool `yaml:"fatal"`     // Capture fatal-level log statements
			Alerts    bool `yaml:"alerts"`    // Capture alert notifications (e.g., PagerDuty, Slack alerts)
			Incidents bool `yaml:"incidents"` // Capture incident handling patterns
		} `yaml:"observability_scan_toggles"`

		// Exclusions define which directories and files should be ignored during scanning.
		Exclusions struct {
			SkipDirectories []string `yaml:"skip_directories"` // List of directories to exclude from scanning
			SkipFiles       []string `yaml:"skip_files"`       // List of specific files to exclude
		} `yaml:"exclusions"`

		// CodeScanExpressions define regex patterns used for identifying code structures.
		CodeScanExpressions map[string]string `yaml:"code_scan_expressions"`

		// ObservabilityScanExpressions define regex patterns used for identifying observability-related patterns.
		ObservabilityScanExpressions map[string]string `yaml:"observability_scan_expressions"`
	} `yaml:"default_approach"`
}

// loadConfig reads a YAML file and unmarshals it into a BlueprintConfig struct.
func loadConfig(configPath string) (*BlueprintConfig, error) {
	log.Printf("[DEBUG] Loading configuration from: %s", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Return an error if reading the file fails.
		return nil, fmt.Errorf("loadConfig: failed to read config file %q: %w", configPath, err)
	}

	// Unmarshal the YAML data into a BlueprintConfig struct.
	var blueprintConfig BlueprintConfig
	if err := yaml.Unmarshal(data, &blueprintConfig); err != nil {
		// Return an error if unmarshaling fails.
		return nil, fmt.Errorf("loadConfig: failed to parse YAML from %q: %w", configPath, err)
	}

	// Return the populated BlueprintConfig and nil error.
	return &blueprintConfig, nil
}

// findGoFiles traverses the given root directory recursively,
// applying configured exclusion rules to skip unwanted paths.
// It returns a list of valid .go file paths.
func findGoFiles(root string, config *BlueprintConfig) ([]string, error) {
	log.Println("[DEBUG] findGoFiles: Starting traversal from", root)

	var goFiles []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[ERROR] findGoFiles: Error accessing path %q: %v", path, err)
			return nil // continue despite the error
		}

		// Check if the path should be skipped based on exclusions.
		if shouldSkipPath(path, config) {
			if info.IsDir() {
				log.Println("[DEBUG] Skipping directory:", path)
				return filepath.SkipDir
			}
			log.Println("[DEBUG] Skipping file:", path)
			return nil
		}

		// Collect only .go files, excluding test files.
		if !info.IsDir() && filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			log.Println("[DEBUG] findGoFiles: Found Go file:", path)
			goFiles = append(goFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("findGoFiles: failed directory traversal: %w", err)
	}

	log.Printf("[DEBUG] findGoFiles: Found total %d Go files", len(goFiles))
	return goFiles, nil
}

// shouldSkipPath - Determines if a path should be skipped based on config
func shouldSkipPath(path string, config *BlueprintConfig) bool {
	log.Println("[DEBUG] shouldSkipPath: Checking if path should be skipped", path)
	for _, dir := range config.DefaultApproach.Exclusions.SkipDirectories {
		if strings.Contains(path, dir) {
			log.Println("[DEBUG] Skipping directory:", path)
			return true
		}
	}
	for _, file := range config.DefaultApproach.Exclusions.SkipFiles {
		if strings.HasSuffix(path, file) {
			log.Println("[DEBUG] Skipping file:", path)
			return true
		}
	}
	return false
}

// parseGoFile opens and reads a Go source file, then processes each line.
// Returns an error if the file can't be opened or read.
func parseGoFile(filePath string, metadata *FileMetadata, config *BlueprintConfig) error {
	log.Println("[DEBUG] parseGoFile: Processing file", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("parseGoFile: unable to open file %q: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	ctx := &ParseContext{}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		searchFileLines(line, lineNumber, filePath, metadata, ctx, config)
	}

	return scanner.Err()
}

// searchFileLines processes a single line from a Go source file,
// passing it to matchPatterns for regex-based parsing.
func searchFileLines(line string, lineNumber int, filePath string, metadata *FileMetadata, ctx *ParseContext, config *BlueprintConfig) {
	log.Printf("[DEBUG] searchFileLines: Analyzing line %d in file %s", lineNumber, filePath)

	matches := matchPatterns(line, config)
	if len(matches) > 0 {
		extractMetadata(matches, lineNumber, filePath, metadata, ctx)
	}
}

// matchPatterns applies regex patterns to a given line and returns any matches.
// Patterns are configured based on the blueprint configuration.
func matchPatterns(line string, config *BlueprintConfig) map[string]string {
	log.Println("[DEBUG] matchPatterns: Matching patterns in line")

	matches := make(map[string]string)

	// Core code structure patterns
	for key, patternStr := range config.DefaultApproach.CodeScanExpressions {
		re := regexp.MustCompile(patternStr)
		if match := re.FindStringSubmatch(line); len(match) > 0 {
			matches[key] = match[len(match)-1]
		}
	}

	// Observability patterns
	for key, enabled := range map[string]bool{
		"logs":      config.DefaultApproach.ObservabilityScanToggles.Logs,
		"metrics":   config.DefaultApproach.ObservabilityScanToggles.Metrics,
		"traces":    config.DefaultApproach.ObservabilityScanToggles.Traces,
		"debug":     config.DefaultApproach.ObservabilityScanToggles.Debug,
		"info":      config.DefaultApproach.ObservabilityScanToggles.Info,
		"error":     config.DefaultApproach.ObservabilityScanToggles.Error,
		"fatal":     config.DefaultApproach.ObservabilityScanToggles.Fatal,
		"alerts":    config.DefaultApproach.ObservabilityScanToggles.Alerts,
		"incidents": config.DefaultApproach.ObservabilityScanToggles.Incidents,
	} {
		if enabled {
			pattern := config.DefaultApproach.ObservabilityScanExpressions[key]
			re := regexp.MustCompile(pattern)
			if re.MatchString(line) {
				matches[key] = strings.TrimSpace(line)
			}
		}
	}

	return matches
}

// extractMetadata processes regex matches and correctly associates comments with the nearest entity.
func extractMetadata(matches map[string]string, lineNumber int, filePath string, metadata *FileMetadata, ctx *ParseContext) {
	log.Printf("[DEBUG] extractMetadata: Processing line %d in file %s", lineNumber, filePath)

	if metadata.Comments == nil {
		metadata.Comments = make(map[int][]ItemWithLine)
	}

	for key, value := range matches {
		switch key {
		case "package":
			metadata.Package = value

		case "function":
			parent := ItemWithLine{Name: value, LineNumber: lineNumber}
			if strings.Title(value) == value {
				metadata.ExternalFunctions = append(metadata.ExternalFunctions, parent)
			} else {
				metadata.InternalFunctions = append(metadata.InternalFunctions, parent)
			}
			ctx.CurrentParent = &parent

		case "struct":
			parent := ItemWithLine{Name: value, LineNumber: lineNumber}
			metadata.Structs = append(metadata.Structs, parent)
			ctx.CurrentParent = &parent

		case "interface":
			parent := ItemWithLine{Name: value, LineNumber: lineNumber}
			metadata.Interfaces = append(metadata.Interfaces, parent)
			ctx.CurrentParent = &parent

		case "import":
			metadata.Imports = append(metadata.Imports, value)

		case "comment":
			clean := strings.TrimSpace(value)
			if clean == "" {
				log.Printf("[DEBUG] Skipping empty comment at [%d]", lineNumber)
				continue
			}
			commentItem := ItemWithLine{Name: clean, LineNumber: lineNumber}

			if ctx.CurrentParent != nil {
				metadata.Comments[ctx.CurrentParent.LineNumber] = append(metadata.Comments[ctx.CurrentParent.LineNumber], commentItem)
				log.Printf("[DEBUG] Linked comment to parent: %s at [%d]", ctx.CurrentParent.Name, ctx.CurrentParent.LineNumber)
			} else {
				metadata.Comments[lineNumber] = append(metadata.Comments[lineNumber], commentItem)
				log.Printf("[DEBUG] Stored standalone comment at [%d]", lineNumber)
			}
		case "logs", "metrics", "traces", "debug", "info", "warn", "error", "fatal", "alerts", "incidents":
			if metadata.ObservabilityHints == nil {
				metadata.ObservabilityHints = make(map[string][]ItemWithLine)
			}
			metadata.ObservabilityHints[key] = append(metadata.ObservabilityHints[key], ItemWithLine{
				Name:       value,
				LineNumber: lineNumber,
			})
			log.Printf("[DEBUG] Captured observability signal [%s] at line %d", key, lineNumber)

		}
	}
}

func organizeMetadata(metadata []FileMetadata) []FileMetadata {
	log.Println("[DEBUG] organizeMetadata: Organizing collected metadata")
	return metadata // Return processed data
}

// summarizeMetadata - Generates summary data
func summarizeMetadata() {
	log.Println("[DEBUG] summarizeMetadata: Summarizing metadata")
}

// collectMetadata walks through files and collects metadata, ensuring comments are correctly linked.
func collectMetadata(files []string, config *BlueprintConfig) []FileMetadata {
	log.Println("[DEBUG] collectMetadata: Starting metadata collection")
	var collectedMetadata []FileMetadata

	for _, file := range files {
		log.Printf("[DEBUG] Processing file: %s", file)
		metadata := FileMetadata{
			FilePath:           file,
			CallMap:            make(map[string][]ItemWithLine),
			ObservabilityHints: make(map[string][]ItemWithLine),
			Comments:           make(map[int][]ItemWithLine),
		}

		err := parseGoFile(file, &metadata, config)
		if err != nil {
			log.Printf("[ERROR] Failed parsing file %q: %v", file, err)
			continue
		}

		collectedMetadata = append(collectedMetadata, metadata)
	}

	log.Println("[DEBUG] collectMetadata: Metadata collection complete")
	return collectedMetadata
}

// generateOutput processes collected details and generates output, structuring comments correctly.
func generateOutput(metadata []FileMetadata, config *BlueprintConfig) {
	outputFile := config.DefaultApproach.OutputManagement.LatestFile
	log.Println("[DEBUG] generateOutput: Writing metadata to output file")

	if shouldRotateFile(config) {
		exists, err := fileExists(outputFile)
		if err != nil {
			log.Printf("[ERROR] generateOutput: error checking if file exists %q: %v", outputFile, err)
			return
		}
		if exists {
			timestamp := time.Now().Format(config.DefaultApproach.OutputManagement.TimestampFormat)
			rotatedFilename := strings.Replace(config.DefaultApproach.OutputManagement.RenameExisting.Pattern, "{timestamp}", timestamp, -1)
			rotateFile(outputFile, rotatedFilename)
		}
	}

	if handled := handleMissingOrEmptyFile(outputFile); handled {
		log.Printf("[DEBUG] generateOutput: handled missing or empty file: %q", outputFile)
	}

	file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("[ERROR] Failed to open output file: %v", err)
	}
	defer file.Close()

	for _, data := range metadata {
		fmt.Fprintf(file, "file: %s\n", data.FilePath)
		fmt.Fprintf(file, "package: %s\n", data.Package)

		// Imports
		fmt.Fprintln(file, "imports:")
		if len(data.Imports) == 0 {
			fmt.Fprintln(file, "- (none)")
		} else {
			for _, imp := range data.Imports {
				fmt.Fprintf(file, "- %s\n", imp)
			}
		}

		// Structs
		fmt.Fprintln(file, "structs:")
		if len(data.Structs) == 0 {
			fmt.Fprintln(file, "- (none)")
		} else {
			for _, s := range data.Structs {
				fmt.Fprintf(file, "- %s [%d]\n", s.Name, s.LineNumber)
				printCommentsWithLabel(file, data.Comments, s.LineNumber, "struct-comments")
			}
		}

		// Interfaces
		fmt.Fprintln(file, "interfaces:")
		if len(data.Interfaces) == 0 {
			fmt.Fprintln(file, "- (none)")
		} else {
			for _, iface := range data.Interfaces {
				fmt.Fprintf(file, "- %s [%d]\n", iface.Name, iface.LineNumber)
				printCommentsWithLabel(file, data.Comments, iface.LineNumber, "interface-comments")
			}
		}

		// Internal Functions
		fmt.Fprintln(file, "internal-funcs:")
		if len(data.InternalFunctions) == 0 {
			fmt.Fprintln(file, "- (none)")
		} else {
			for _, fn := range data.InternalFunctions {
				fmt.Fprintf(file, "- %s [%d]\n", fn.Name, fn.LineNumber)
				printCommentsWithLabel(file, data.Comments, fn.LineNumber, "func-comments")
			}
		}

		// External Functions
		fmt.Fprintln(file, "external-funcs:")
		if len(data.ExternalFunctions) == 0 {
			fmt.Fprintln(file, "- (none)")
		} else {
			for _, fn := range data.ExternalFunctions {
				fmt.Fprintf(file, "- %s [%d]\n", fn.Name, fn.LineNumber)
				printCommentsWithLabel(file, data.Comments, fn.LineNumber, "func-comments")
			}
		}

		// Observability
		fmt.Fprintln(file, "observability-hints:")
		if len(data.ObservabilityHints) == 0 {
			fmt.Fprintln(file, "- (none)")
		} else {
			for key, items := range data.ObservabilityHints {
				fmt.Fprintf(file, "  %s [%d]:\n", key, len(items))
				for _, item := range items {
					fmt.Fprintf(file, "    - %s [%d]\n", item.Name, item.LineNumber)
				}
			}
			fmt.Fprintln(file, "----------------------------------------")

		}

		writeFileSummary(file, data)
		// File separator
		fmt.Fprintln(file, "====================================")
	}

	writeGlobalScanSummary(file, metadata)

	log.Println("[DEBUG] generateOutput: Output successfully written to", outputFile)
}

// printComments writes comments associated with a parent entity
func printComments(file *os.File, comments map[int][]ItemWithLine, parentLine int) {
	for _, comment := range comments[parentLine] {
		fmt.Fprintf(file, "  - %s [%d, %d]\n", strings.TrimSpace(comment.Name), parentLine, comment.LineNumber)
	}
}

func writeFileSummary(file *os.File, data FileMetadata) {
	fmt.Fprintln(file, "summary:")
	fmt.Fprintf(file, "  total-structs: %d\n", len(data.Structs))
	fmt.Fprintf(file, "  total-interfaces: %d\n", len(data.Interfaces))
	fmt.Fprintf(file, "  total-internal-funcs: %d\n", len(data.InternalFunctions))
	fmt.Fprintf(file, "  total-external-funcs: %d\n", len(data.ExternalFunctions))

	totalComments := 0
	for _, list := range data.Comments {
		totalComments += len(list)
	}
	fmt.Fprintf(file, "  total-comments: %d\n", totalComments)

	if len(data.ObservabilityHints) == 0 {
		fmt.Fprintf(file, "  total-observability: 0\n")
	} else {
		for key, list := range data.ObservabilityHints {
			fmt.Fprintf(file, "  total-%s: %d\n", key, len(list))
		}
	}
}

func writeGlobalScanSummary(file *os.File, metadata []FileMetadata) {
	fmt.Fprintln(file, "code-scan-summary:")
	fmt.Fprintf(file, "  total-files: %d\n", len(metadata))

	totalStructs := 0
	totalInterfaces := 0
	totalInternalFuncs := 0
	totalExternalFuncs := 0
	totalComments := 0
	observabilityTotals := make(map[string]int)

	for _, data := range metadata {
		totalStructs += len(data.Structs)
		totalInterfaces += len(data.Interfaces)
		totalInternalFuncs += len(data.InternalFunctions)
		totalExternalFuncs += len(data.ExternalFunctions)
		for _, list := range data.Comments {
			totalComments += len(list)
		}
		for key, list := range data.ObservabilityHints {
			observabilityTotals[key] += len(list)
		}
	}

	fmt.Fprintf(file, "  total-structs: %d\n", totalStructs)
	fmt.Fprintf(file, "  total-interfaces: %d\n", totalInterfaces)
	fmt.Fprintf(file, "  total-internal-funcs: %d\n", totalInternalFuncs)
	fmt.Fprintf(file, "  total-external-funcs: %d\n", totalExternalFuncs)
	fmt.Fprintf(file, "  total-comments: %d\n", totalComments)

	for key, count := range observabilityTotals {
		fmt.Fprintf(file, "  total-%s: %d\n", key, count)
	}
}

func printCommentsWithLabel(file *os.File, comments map[int][]ItemWithLine, parentLine int, label string) {
	if commentList, ok := comments[parentLine]; ok && len(commentList) > 0 {
		fmt.Fprintf(file, "  %s:\n", label)
		for _, c := range commentList {
			fmt.Fprintf(file, "    - %s [%d, %d]\n", c.Name, parentLine, c.LineNumber)
		}
	}
}

func main() {
	log.Println("\n====== Blueprint utility started ======")

	// Load initial config to obtain the real config file path.
	initialConfigPath := "../../configs/code-blueprint.yaml"
	tempConfig, err := loadConfig(initialConfigPath)
	if err != nil {
		log.Fatalf("[ERROR] main: Failed loading initial config: %v", err)
	}

	// Use actual config file path from initial config.
	configPath := tempConfig.Settings.ConfigFilePath
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("[ERROR] main: Failed loading config from %q: %v", configPath, err)
	}

	log.Printf("[DEBUG] main: Loaded configuration from %s", configPath)

	// Set output file name from config.
	outputFile := config.DefaultApproach.OutputManagement.LatestFile

	// Determine absolute root directory path.
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		log.Fatalf("[ERROR] main: Unable to determine absolute path: %v", err)
	}
	log.Printf("[DEBUG] main: Scanning directory %s", rootDir)

	// Find Go files in directory.
	files, err := findGoFiles(rootDir, config)
	if err != nil {
		log.Fatalf("[ERROR] main: Failed finding Go files: %v", err)
	}

	var collectedMetadata []FileMetadata
	for _, file := range files {
		metadata := FileMetadata{
			FilePath:           file,
			InternalFunctions:  []ItemWithLine{},
			ExternalFunctions:  []ItemWithLine{},
			Structs:            []ItemWithLine{},
			Interfaces:         []ItemWithLine{},
			Imports:            []string{},
			CallMap:            make(map[string][]ItemWithLine),
			ObservabilityHints: make(map[string][]ItemWithLine),
			Comments:           make(map[int][]ItemWithLine),
		}

		err := parseGoFile(file, &metadata, config)

		if err != nil {
			log.Printf("[ERROR] Failed parsing file %q: %v", file, err)
			continue
		}

		collectedMetadata = append(collectedMetadata, metadata)
	}

	// Organize metadata and generate output.
	// organizedMetadata := organizeMetadata(collectedMetadata)
	generateOutput(collectedMetadata, config)

	log.Printf("[DEBUG] main: Summary successfully saved to %s", outputFile)
	log.Println("\n====== Blueprint utility completed ======")
}

// Helper function to check if a file exists.
func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// isValidFile ensures only Go source files (excluding test files) are processed.
func isValidFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".go") && !strings.HasSuffix(fileName, "_test.go")
}

// createBlueprintFile creates a new file at the given filename.
// It returns an error if the file creation fails.
func createBlueprintFile(filename string) error {
	// Create the file. os.Create truncates the file if it already exists.
	file, err := os.Create(filename)
	if err != nil {
		// Return the error directly. No need for wrapping in this case.
		return err
	}
	// Ensure the file is closed, even if an error occurs later.
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the close error. Consider returning it if critical.
			log.Printf("createBlueprintFile: error closing file %q: %v", filename, closeErr)
		}
	}()

	// Return nil to indicate successful creation.
	return nil
}

// handleMissingOrEmptyFile checks if a file exists and is not empty.
// If the file does not exist, it creates a new file.
// If the file exists but is empty, it removes and recreates the file.
// It returns true if the file was handled (created or recreated), false otherwise.
func handleMissingOrEmptyFile(file string) bool {
	fileInfo, err := os.Stat(file)

	if os.IsNotExist(err) {
		// File does not exist. Create a new file.
		log.Printf("handleMissingOrEmptyFile: file %q does not exist, creating", file)
		if err := createBlueprintFile(file); err != nil {
			log.Printf("handleMissingOrEmptyFile: error creating file %q: %v", file, err)
			return false // Creation failed, but we should not say that the file was handled.
		}
		return true
	}

	if err != nil {
		// os.Stat returned an error other than IsNotExist.
		log.Printf("handleMissingOrEmptyFile: error getting file info for %q: %v", file, err)
		return false // Getting file info failed, we should not say that the file was handled.
	}

	if fileInfo.Size() == 0 {
		// File exists but is empty. Remove and recreate.
		log.Printf("handleMissingOrEmptyFile: file %q is empty, removing and recreating", file)
		if err := os.Remove(file); err != nil {
			log.Printf("handleMissingOrEmptyFile: error removing file %q: %v", file, err)
			return false // Removal failed, we should not say that the file was handled.
		}

		if err := createBlueprintFile(file); err != nil {
			log.Printf("handleMissingOrEmptyFile: error recreating file %q: %v", file, err)
			return false // Creation failed, we should not say that the file was handled.
		}
		return true
	}

	// File exists and is not empty. No action taken.
	return false
}

func shouldRotateFile(config *BlueprintConfig) bool {
	return config.DefaultApproach.OutputManagement.EnableRotation &&
		config.DefaultApproach.OutputManagement.RenameExisting.PreserveOldVersions
}

// rotateFile renames oldFile to newFile if newFile does not exist.
// It logs warnings and errors during the rotation process.
func rotateFile(oldFile, newFile string) {
	_, err := os.Stat(newFile) // Get os.Stat results.

	if os.IsNotExist(err) {
		// newFile does not exist, proceed with renaming.
		if err := os.Rename(oldFile, newFile); err != nil {
			log.Printf("rotateFile: warning: failed to rotate %q to %q: %v", oldFile, newFile, err)
			return // Exit function on error.
		}

		log.Printf("rotateFile: rotated %q to %q", oldFile, newFile)
		return
	}

	if err != nil {
		// os.Stat returned an error other than IsNotExist.
		log.Printf("rotateFile: warning: failed to stat %q : %v", newFile, err)
		return
	}

	// newFile already exists, skip rotation.
	log.Printf("rotateFile: skipping rotation: %q already exists", newFile)
}

// createNewBlueprintFile creates a new blueprint file with the given filename.
// It logs an error if the creation fails, and a success message otherwise.
func createNewBlueprintFile(filename string) {
	if err := createBlueprintFile(filename); err != nil {
		log.Printf("createNewBlueprintFile: error: failed to create file %q: %v", filename, err)
		return
	}
	log.Printf("createNewBlueprintFile: created new blueprint file: %q", filename)
}
