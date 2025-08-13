package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"shell/shell"
	"strings"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		for sig := range sigs {
			if shell.CurrentCmd != nil && shell.CurrentCmd.Process != nil {
				_ = shell.CurrentCmd.Process.Signal(sig)
				if sig == syscall.SIGQUIT {
					time.Sleep(100 * time.Millisecond)
					if shell.CurrentCmd.Process != nil {
						_ = shell.CurrentCmd.Process.Kill()
					}
				}
			} else {
				fmt.Print("\nshell> ")
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("shell> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nВыход из shell.")
				return
			}
			fmt.Println("\nОшибка чтения:", err)
			continue
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			continue
		}

		line = shell.ExpandVariables(line)

		cmdSeq := shell.ParseCommandLine(line)

		if err := shell.ExecuteLogicalSequence(cmdSeq); err != nil {
			fmt.Fprintln(os.Stderr, "ошибка:", err)
		}
	}
}
