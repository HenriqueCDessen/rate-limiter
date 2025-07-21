package main_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	baseURL          = "http://localhost:8080"
	requestsPerIP    = 12
	requestsPerToken = 102
	testToken        = "abc123"
	redisAddr        = "localhost:6379"
	redisPassword    = "your_redis_password"
	redisDB          = 0
)

var (
	redisClient *redis.Client
)

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
}

func clearRedisKeys() {
	ctx := context.Background()
	iter := redisClient.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		redisClient.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		fmt.Printf("Erro ao limpar Redis: %v\n", err)
	}
}

func makeRequest(client *http.Client, url string, headers map[string]string) (int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func testIPRateLimiting() {
	fmt.Println("\n=== Testando limitação por IP ===")
	clearRedisKeys() // Limpa estado anterior

	client := &http.Client{}
	results := make(chan int, requestsPerIP)
	var wg sync.WaitGroup

	for i := 1; i <= requestsPerIP; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status, err := makeRequest(client, baseURL, nil)
			if err != nil {
				fmt.Printf("Erro na requisição %d: %v\n", i, err)
				return
			}
			results <- status
			fmt.Printf("Requisição %d: Status %d\n", i, status)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	wg.Wait()
	close(results)

	successCount := 0
	rateLimitedCount := 0
	for status := range results {
		if status == http.StatusOK {
			successCount++
		} else if status == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	fmt.Printf("\nResumo IP:\n- Sucessos: %d (esperado: 10)\n- Limites excedidos: %d (esperado: 2)\n",
		successCount, rateLimitedCount)
}

func testTokenRateLimiting() {
	fmt.Println("\n=== Testando limitação por Token ===")
	clearRedisKeys()

	client := &http.Client{}
	results := make(chan int, requestsPerToken)
	var wg sync.WaitGroup

	for i := 1; i <= requestsPerToken; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			status, err := makeRequest(client, baseURL, map[string]string{
				"API_KEY": testToken,
			})
			if err != nil {
				fmt.Printf("Erro na requisição %d: %v\n", i, err)
				return
			}
			results <- status
			fmt.Printf("Requisição %d: Status %d\n", i, status)
			time.Sleep(20 * time.Millisecond)
		}(i)
	}

	wg.Wait()
	close(results)

	successCount := 0
	rateLimitedCount := 0
	for status := range results {
		if status == http.StatusOK {
			successCount++
		} else if status == http.StatusTooManyRequests {
			rateLimitedCount++
		}
	}

	fmt.Printf("\nResumo Token:\n- Sucessos: %d (esperado: 100)\n- Limites excedidos: %d (esperado: 2)\n",
		successCount, rateLimitedCount)
}

func testBlockTime() {
	fmt.Println("\n=== Testando tempo de bloqueio ===")
	clearRedisKeys()

	client := &http.Client{}

	status, err := makeRequest(client, baseURL, nil)
	if err != nil {
		fmt.Printf("Erro na requisição inicial: %v\n", err)
		return
	}
	fmt.Printf("Requisição inicial: Status %d (esperado: 200)\n", status)

	for i := 1; i <= 11; i++ {
		makeRequest(client, baseURL, nil)
	}

	status, err = makeRequest(client, baseURL, nil)
	if err != nil {
		fmt.Printf("Erro na requisição pós-excesso: %v\n", err)
		return
	}
	fmt.Printf("Requisição pós-excesso: Status %d (esperado: 429)\n", status)

	blockTime := 65 * time.Second
	fmt.Printf("Aguardando %v para o bloqueio expirar...\n", blockTime)
	time.Sleep(blockTime)

	status, err = makeRequest(client, baseURL, nil)
	if err != nil {
		fmt.Printf("Erro na requisição pós-bloqueio: %v\n", err)
		return
	}
	fmt.Printf("Requisição pós-bloqueio: Status %d (esperado: 200)\n", status)
}

func main() {
	initRedis()
	defer redisClient.Close()

	fmt.Println("Iniciando testes do Rate Limiter...")

	testIPRateLimiting()
	testTokenRateLimiting()
	testBlockTime()

	fmt.Println("\nTestes concluídos!")
}
