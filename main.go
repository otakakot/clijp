package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	copilot "github.com/github/copilot-sdk/go"
)

const (
	prompt  = "コマンドラインで受け取った内容を日本語に翻訳してください。コードや記号はそのまま保持してください。結果のみを翻訳した状態で、元の形式を保持したまま出力してください。コードブロックなども不要です。：\n\n%s"
	model   = "gpt-5-mini" // ref: https://docs.github.com/en/copilot/concepts/billing/copilot-requests
	version = "0.0.2"
	help    = "clijp v" + version + " - 標準入力で受け取った内容を Copilot SDK を使って日本語に翻訳するツール"
)

func showHelp() {
	fmt.Fprintln(os.Stderr, help)
}

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		showHelp()
		log.Fatalf("標準入力の情報取得に失敗しました: %v", err)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		showHelp()
		return
	}

	reader := bufio.NewReader(os.Stdin)

	builder := strings.Builder{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				builder.WriteString(line)
				break
			}
			showHelp()
			log.Fatalf("標準入力の読み取りに失敗しました: %v", err)
		}
		builder.WriteString(line)
	}

	input := builder.String()

	if strings.TrimSpace(input) == "" {
		showHelp()
		return
	}

	fmt.Print(strings.TrimSuffix(input, "\n"))

	client := copilot.NewClient(&copilot.ClientOptions{
		LogLevel: "error",
	})

	ctx := context.Background()

	if err := client.Start(ctx); err != nil {
		log.Fatalf("Copilot クライアントの開始に失敗しました: %v", err)
	}
	defer client.Stop()

	session, err := client.CreateSession(ctx, &copilot.SessionConfig{
		Model: model,
	})
	if err != nil {
		log.Fatalf("セッションの作成に失敗しました: %v", err)
	}
	defer session.Destroy()

	translation := strings.Builder{}

	done := make(chan bool)

	loadingDone := make(chan bool)

	session.On(func(event copilot.SessionEvent) {
		if event.Type == "assistant.message" {
			if event.Data.Content != nil {
				translation.WriteString(*event.Data.Content)
			}
		}
		if event.Type == "session.idle" {
			close(done)
		}
	})

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-loadingDone:
				fmt.Fprint(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				fmt.Fprintf(os.Stderr, "\r%s 翻訳中...", frames[i%len(frames)])
				i++
			case <-ctx.Done():
				return
			}
		}
	}()

	translationPrompt := fmt.Sprintf(prompt, input)

	if _, err = session.Send(ctx, copilot.MessageOptions{
		Prompt: translationPrompt,
	}); err != nil {
		close(loadingDone)
		log.Fatalf("メッセージの送信に失敗しました: %v", err)
	}

	<-done
	close(loadingDone)
	time.Sleep(50 * time.Millisecond)

	fmt.Print("\n=== 日本語翻訳 ===\n\n")
	fmt.Print(translation.String())
}
