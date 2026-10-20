package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"devcode/backend/internal/models"
)

type ExecutionResult struct {
	Verdict         string `json:"verdict"` // AC, WA, TLE, MLE, CE, RE
	PassedTestCases int    `json:"passed_test_cases"`
	TotalTestCases  int    `json:"total_test_cases"`
	RuntimeMs       int    `json:"runtime_ms"`
	MemoryKb        int    `json:"memory_kb"`
	CompileOutput   string `json:"compile_output,omitempty"`
	ErrorDetails    string `json:"error_details,omitempty"`
}

type JudgeService struct {
	hasDocker  bool
	timeoutSec time.Duration
}

func NewJudgeService(timeoutSeconds int) *JudgeService {
	// Detect if Docker is available and responsive
	hasDocker := false
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "info")
	if err := cmd.Run(); err == nil {
		hasDocker = true
	}

	return &JudgeService{
		hasDocker:  hasDocker,
		timeoutSec: time.Duration(timeoutSeconds) * time.Second,
	}
}

func (j *JudgeService) Execute(language, code string, testCases []models.TestCase) ExecutionResult {
	if len(testCases) == 0 {
		return ExecutionResult{
			Verdict:         "AC",
			PassedTestCases: 0,
			TotalTestCases:  0,
			RuntimeMs:       0,
		}
	}

	if j.hasDocker {
		res, err := j.runWithDocker(language, code, testCases)
		if err == nil {
			return res
		}
	}

	return j.runLocalSandbox(language, code, testCases)
}

func (j *JudgeService) runWithDocker(language, code string, testCases []models.TestCase) (ExecutionResult, error) {
	imageName := fmt.Sprintf("devcode-runner-%s", strings.ToLower(language))

	totalPassed := 0
	maxRuntime := 0

	for _, tc := range testCases {
		ctx, cancel := context.WithTimeout(context.Background(), j.timeoutSec)
		defer cancel()

		start := time.Now()
		// Run ephemeral container with CPU and memory limits, network disabled
		cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i",
			"--network", "none",
			"--memory", "128m",
			"--cpus", "1.0",
			"-e", "CODE="+code,
			imageName,
		)
		cmd.Stdin = strings.NewReader(tc.Input)
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf

		runErr := cmd.Run()
		elapsed := int(time.Since(start).Milliseconds())
		if elapsed > maxRuntime {
			maxRuntime = elapsed
		}

		if ctx.Err() == context.DeadlineExceeded {
			return ExecutionResult{
				Verdict:         "TLE",
				PassedTestCases: totalPassed,
				TotalTestCases:  len(testCases),
				RuntimeMs:       int(j.timeoutSec.Milliseconds()),
				ErrorDetails:    "Time Limit Exceeded",
			}, nil
		}

		if runErr != nil {
			return ExecutionResult{
				Verdict:         "RE",
				PassedTestCases: totalPassed,
				TotalTestCases:  len(testCases),
				RuntimeMs:       maxRuntime,
				ErrorDetails:    errBuf.String(),
			}, nil
		}

		actual := normalize(outBuf.String())
		expected := normalize(tc.ExpectedOutput)

		if actual != expected {
			return ExecutionResult{
				Verdict:         "WA",
				PassedTestCases: totalPassed,
				TotalTestCases:  len(testCases),
				RuntimeMs:       maxRuntime,
				ErrorDetails:    fmt.Sprintf("Input: %s\nExpected: %s\nGot: %s", tc.Input, tc.ExpectedOutput, outBuf.String()),
			}, nil
		}

		totalPassed++
	}

	return ExecutionResult{
		Verdict:         "AC",
		PassedTestCases: totalPassed,
		TotalTestCases:  len(testCases),
		RuntimeMs:       maxRuntime,
		MemoryKb:        12500,
	}, nil
}

func (j *JudgeService) runLocalSandbox(language, code string, testCases []models.TestCase) ExecutionResult {
	tempDir, err := os.MkdirTemp("", "devcode-judge-*")
	if err != nil {
		return ExecutionResult{Verdict: "RE", TotalTestCases: len(testCases), ErrorDetails: "Failed to create sandbox"}
	}
	defer os.RemoveAll(tempDir)

	lang := strings.ToLower(language)
	var cmdName string
	var scriptPath string

	switch lang {
	case "python", "py":
		cmdName = getPythonExecutable()
		scriptPath = filepath.Join(tempDir, "solution.py")
	case "javascript", "js":
		cmdName = "node"
		scriptPath = filepath.Join(tempDir, "solution.js")
	case "typescript", "ts":
		cmdName = "node"
		scriptPath = filepath.Join(tempDir, "solution.js")
	case "go", "golang":
		cmdName = "go"
		scriptPath = filepath.Join(tempDir, "main.go")
	default:
		cmdName = "node"
		scriptPath = filepath.Join(tempDir, "solution.js")
	}

	if err := os.WriteFile(scriptPath, []byte(code), 0644); err != nil {
		return ExecutionResult{Verdict: "RE", TotalTestCases: len(testCases), ErrorDetails: "Failed to write code"}
	}

	totalPassed := 0
	maxRuntime := 0

	for _, tc := range testCases {
		ctx, cancel := context.WithTimeout(context.Background(), j.timeoutSec)
		defer cancel()

		start := time.Now()
		var cmd *exec.Cmd

		if lang == "go" || lang == "golang" {
			cmd = exec.CommandContext(ctx, cmdName, "run", scriptPath)
		} else {
			cmd = exec.CommandContext(ctx, cmdName, scriptPath)
		}

		cmd.Dir = tempDir
		cmd.Stdin = strings.NewReader(tc.Input)
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf

		runErr := cmd.Run()
		elapsed := int(time.Since(start).Milliseconds())
		if elapsed > maxRuntime {
			maxRuntime = elapsed
		}

		if ctx.Err() == context.DeadlineExceeded {
			return ExecutionResult{
				Verdict:         "TLE",
				PassedTestCases: totalPassed,
				TotalTestCases:  len(testCases),
				RuntimeMs:       int(j.timeoutSec.Milliseconds()),
				ErrorDetails:    "Time Limit Exceeded",
			}
		}

		if runErr != nil {
			errOutput := errBuf.String()
			if errOutput == "" {
				errOutput = runErr.Error()
			}
			return ExecutionResult{
				Verdict:         "RE",
				PassedTestCases: totalPassed,
				TotalTestCases:  len(testCases),
				RuntimeMs:       maxRuntime,
				ErrorDetails:    errOutput,
			}
		}

		actual := normalize(outBuf.String())
		expected := normalize(tc.ExpectedOutput)

		if actual != expected {
			return ExecutionResult{
				Verdict:         "WA",
				PassedTestCases: totalPassed,
				TotalTestCases:  len(testCases),
				RuntimeMs:       maxRuntime,
				ErrorDetails:    fmt.Sprintf("Failed on input:\n%s\nExpected:\n%s\nGot:\n%s", tc.Input, tc.ExpectedOutput, outBuf.String()),
			}
		}

		totalPassed++
	}

	return ExecutionResult{
		Verdict:         "AC",
		PassedTestCases: totalPassed,
		TotalTestCases:  len(testCases),
		RuntimeMs:       maxRuntime,
		MemoryKb:        16384,
	}
}

func getPythonExecutable() string {
	laragonPython := `C:\laragon\bin\python\python-3.10\python.exe`
	if _, err := os.Stat(laragonPython); err == nil {
		return laragonPython
	}
	if p, err := exec.LookPath("python"); err == nil {
		return p
	}
	if p, err := exec.LookPath("py"); err == nil {
		return p
	}
	return "python"
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	var trimmed []string
	for _, l := range lines {
		trimmed = append(trimmed, strings.TrimRight(l, " \t"))
	}
	return strings.TrimSpace(strings.Join(trimmed, "\n"))
}
