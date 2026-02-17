package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	copilot "github.com/github/copilot-sdk/go"
)

const (
	prompt  = "コマンドラインで受け取った内容を日本語に翻訳してください。コードや記号はそのまま保持してください。結果のみを翻訳した状態で、元の形式を保持したまま出力してください。コードブロックなども不要です。：\n\n%s"
	model   = "gpt-5-mini" // ref: https://docs.github.com/en/copilot/concepts/billing/copilot-requests
	version = "0.0.4"
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

	ctx := context.Background()

	// Prepare cache path: ~/.cache/clijp/<sha256>.txt
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("ホームディレクトリの取得に失敗しました: %v", err)
	}

	cacheDir := filepath.Join(home, ".cache", "clijp")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		log.Printf("キャッシュディレクトリの作成に失敗しました: %v", err)
	}

	h := sha256.Sum256([]byte(input))
	fname := hex.EncodeToString(h[:]) + ".txt"
	cachePath := filepath.Join(cacheDir, fname)

	if data, err := os.ReadFile(cachePath); err == nil {
		time.Sleep(50 * time.Millisecond)
		fmt.Print("\n=== 日本語翻訳 (キャッシュ) ===\n\n")
		fmt.Print(string(data))
		return
	}

	done := make(chan struct{})
	cleared := make(chan struct{})
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Fprintf(os.Stderr, "\r\033[K")
				close(cleared)
				return
			default:
				fmt.Fprintf(os.Stderr, "\r%s 翻訳中...", frames[i%len(frames)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	translated, err := TranslateJP(ctx, input)

	close(done)
	<-cleared

	if err != nil {
		log.Fatalf("翻訳に失敗しました: %v", err)
	}

	if err := os.WriteFile(cachePath, []byte(translated), 0o644); err != nil {
		log.Printf("キャッシュの書き込みに失敗しました: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	fmt.Print("\n=== 日本語翻訳 ===\n\n")
	fmt.Print(translated)
}

func TranslateJP(ctx context.Context, input string) (string, error) {
	client := copilot.NewClient(&copilot.ClientOptions{
		LogLevel: "error",
	})

	if err := client.Start(ctx); err != nil {
		return "", fmt.Errorf("Copilot クライアントの開始に失敗しました: %w", err)
	}
	defer client.Stop()

	session, err := client.CreateSession(ctx, &copilot.SessionConfig{Model: model})
	if err != nil {
		return "", fmt.Errorf("セッションの作成に失敗しました: %w", err)
	}
	defer session.Destroy()

	resp, err := session.SendAndWait(ctx, copilot.MessageOptions{
		Prompt: fmt.Sprintf(prompt, input),
	})
	if err != nil {
		return "", fmt.Errorf("メッセージの送信に失敗しました: %w", err)
	}

	if resp == nil || resp.Data.Content == nil {
		return "", fmt.Errorf("Copilot から有効なレスポンスが返ってきませんでした")
	}

	return *resp.Data.Content, nil
}
