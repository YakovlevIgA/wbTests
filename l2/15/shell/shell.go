package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var CurrentCmd *exec.Cmd

func ExpandVariables(input string) string {
	return os.Expand(input, func(key string) string {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		if key == "HOME" {
			return os.Getenv("USERPROFILE")
		}
		return ""
	})
}

func ParseCommandLine(input string) []string {
	input = ExpandVariables(input)
	return tokenize(input)
}

func tokenize(input string) []string {
	var tokens []string
	var buf strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch {
		case c == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case c == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case !inSingleQuote && !inDoubleQuote && (strings.HasPrefix(input[i:], "&&") || strings.HasPrefix(input[i:], "||")):
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			tokens = append(tokens, input[i:i+2])
			i++
		case !inSingleQuote && !inDoubleQuote && (c == '|' || c == '>' || c == '<'):
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			tokens = append(tokens, string(c))
		default:
			buf.WriteByte(c)
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}

	return mergeOps(tokens)
}

func mergeOps(tokens []string) []string {
	res := []string{}
	curr := ""

	for _, t := range tokens {
		if t == "&&" || t == "||" || t == "|" || t == ">" || t == "<" || t == ">>" {
			if strings.TrimSpace(curr) != "" {
				res = append(res, strings.TrimSpace(curr))
				curr = ""
			}
			res = append(res, t)
		} else {
			curr += t
		}
	}

	if strings.TrimSpace(curr) != "" {
		res = append(res, strings.TrimSpace(curr))
	}

	return res
}

func IsBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "pwd", "echo", "kill", "ps", "true", "false", "tr", "tee":
		return true
	default:
		return false
	}
}

func RunBuiltin(cmd string, args []string, stdin io.Reader, stdout io.Writer) error {
	switch cmd {
	case "cd":
		if len(args) < 1 {
			return os.Chdir(os.Getenv("HOME"))
		}
		return os.Chdir(args[0])
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, dir)
		return nil
	case "echo":
		fmt.Fprintln(stdout, strings.Join(args, " "))
		return nil
	case "kill":
		if len(args) < 1 {
			return errors.New("укажите PID")
		}
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		return proc.Kill()
	case "ps":
		var cmdExec *exec.Cmd
		if runtime.GOOS == "windows" {
			cmdExec = exec.Command("tasklist")
		} else {
			cmdExec = exec.Command("ps", "aux")
		}
		cmdExec.Stdout = stdout
		cmdExec.Stderr = os.Stderr
		return cmdExec.Run()
	case "true":
		return nil
	case "false":
		return errors.New("команда false возвращена с ошибкой")
	case "tr":
		if len(args) < 2 {
			return errors.New("необходимо указать from-set to-set")
		}
		input, err := io.ReadAll(stdin)
		if err != nil {
			return err
		}

		from := args[0]
		to := args[1]

		if from == "a-z" && to == "A-Z" {
			result := strings.ToUpper(string(input))
			fmt.Fprint(stdout, result)
			return nil
		}
		return fmt.Errorf("неподдерживаемые наборы символов: %s %s", from, to)
	case "tee":
		if len(args) < 1 {
			return errors.New("укажите файл")
		}
		file, err := os.Create(args[0])
		if err != nil {
			return err
		}
		defer file.Close()

		input, err := io.ReadAll(stdin)
		if err != nil {
			return err
		}

		_, err = file.Write(input)
		if err != nil {
			return err
		}

		_, err = stdout.Write(input)
		return err
	}
	return fmt.Errorf("неизвестная команда: %s", cmd)
}

func processRedirects(args []string) ([]string, io.Reader, io.Writer, error) {
	var stdin io.Reader = os.Stdin
	var stdout io.Writer = os.Stdout
	cleanArgs := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == ">":
			if i+1 >= len(args) {
				return nil, nil, nil, errors.New("не указан файл для перенаправления")
			}
			file, err := os.Create(args[i+1])
			if err != nil {
				return nil, nil, nil, err
			}
			stdout = file
			i++
		case arg == "<":
			if i+1 >= len(args) {
				return nil, nil, nil, errors.New("не указан файл для ввода")
			}
			file, err := os.Open(args[i+1])
			if err != nil {
				return nil, nil, nil, err
			}
			stdin = file
			i++
		default:
			cleanArgs = append(cleanArgs, arg)
		}
	}

	return cleanArgs, stdin, stdout, nil
}

func ExecutePipeline(commands []string) error {
	if len(commands) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(commands))
	pipes := make([]*io.PipeReader, len(commands)-1)

	for i, cmdStr := range commands {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue
		}

		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		cmdName := parts[0]
		args := parts[1:]

		args, stdin, stdout, err := processRedirects(args)
		if err != nil {
			return err
		}

		if i > 0 {
			stdin = pipes[i-1]
		}
		if i < len(commands)-1 {
			r, w := io.Pipe()
			if stdout == os.Stdout {
				stdout = w
			}
			pipes[i] = r
			defer r.Close()
		}

		wg.Add(1)
		go func(cmdName string, args []string, stdin io.Reader, stdout io.Writer, idx int) {
			defer wg.Done()

			if closer, ok := stdin.(io.Closer); ok && idx > 0 {
				defer closer.Close()
			}

			if IsBuiltin(cmdName) {
				if err := RunBuiltin(cmdName, args, stdin, stdout); err != nil {
					errChan <- err
				}
				return
			}

			cmd := exec.Command(cmdName, args...)
			CurrentCmd = cmd
			cmd.Stdin = stdin
			cmd.Stdout = stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				errChan <- err
			}
			CurrentCmd = nil
		}(cmdName, args, stdin, stdout, i)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func ExecuteLogicalSequence(sequence []string) error {
	skipNext := false
	lastSuccess := true

	for i := 0; i < len(sequence); i++ {
		cmd := sequence[i]

		if cmd == "&&" {
			if !lastSuccess {
				skipNext = true
			}
			continue
		} else if cmd == "||" {
			if lastSuccess {
				skipNext = true
			}
			continue
		}

		if skipNext {
			skipNext = false
			continue
		}

		commands := strings.Split(cmd, "|")
		err := ExecutePipeline(commands)
		lastSuccess = (err == nil)
	}

	return nil
}
