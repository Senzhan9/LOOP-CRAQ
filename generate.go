package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ZipfianGenerator 生成 Zipf 分布下的随机数
type ZipfianGenerator struct {
	N    int
	skew float64
	cdf  []float64
}

func NewZipfianGenerator(n int, skew float64) *ZipfianGenerator {
	z := &ZipfianGenerator{N: n, skew: skew}
	z.buildCDF()
	return z
}

func (z *ZipfianGenerator) buildCDF() {
	z.cdf = make([]float64, z.N)
	var norm float64
	for i := 1; i <= z.N; i++ {
		norm += 1.0 / math.Pow(float64(i), z.skew)
	}
	sum := 0.0
	for i := 1; i <= z.N; i++ {
		sum += 1.0 / math.Pow(float64(i), z.skew)
		z.cdf[i-1] = sum / norm
	}
}

func (z *ZipfianGenerator) Next(r *rand.Rand) int {
	rnd := r.Float64()
	for i, p := range z.cdf {
		if p >= rnd {
			return i
		}
	}
	return z.N - 1
}

// 随机字符串生成
func generateRandomValue(length int, r *rand.Rand) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteByte(charset[r.Intn(len(charset))])
	}
	return b.String()
}

func main() {
	var (
		keyCount       = flag.Int("key_count", 10, "Total number of keys")
		readProportion = flag.Float64("read_proportion", 0.5, "Proportion of read operations (0~1)")
		valueLength    = flag.Int("value_length", 4, "Length of the value string for write ops")
		distribution   = flag.String("distribution", "uniform", "Key distribution: uniform or zipf")
		totalOps       = flag.Int("total_operations", 1000, "Total number of operations to generate")
	)

	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Usage: ./generate --key_count=N --read_proportion=x --value_length=L --distribution=uniform|zipf output.txt")
		os.Exit(1)
	}
	outputFile := args[0]

	if *readProportion < 0 || *readProportion > 1 {
		fmt.Println("Error: read_proportion must be between 0 and 1")
		os.Exit(1)
	}
	if *valueLength <= 0 {
		fmt.Println("Error: value_length must be positive")
		os.Exit(1)
	}
	if *distribution != "uniform" && *distribution != "zipf" {
		fmt.Println("Error: distribution must be 'uniform' or 'zipf'")
		os.Exit(1)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		os.Exit(1)
	}
	defer file.Close()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var zipfGen *ZipfianGenerator
	if *distribution == "zipf" {
		zipfGen = NewZipfianGenerator(*keyCount, 1.2)
	}

	for i := 0; i < *totalOps; i++ {
		isRead := r.Float64() < *readProportion

		var keyIndex int
		if *distribution == "uniform" {
			keyIndex = r.Intn(*keyCount)
		} else {
			keyIndex = zipfGen.Next(r)
		}

		key := fmt.Sprintf("key%d", keyIndex)
		if isRead {
			fmt.Fprintf(file, "/opt/craq/craq-client -c l-coord:1234 read %s\n", key)
		} else {
			value := generateRandomValue(*valueLength, r)
			fmt.Fprintf(file, "/opt/craq/craq-client -c l-coord:1234 write %s %s\n", key, value)
		}
	}

	fmt.Printf("Generated %d operations into %s\n", *totalOps, outputFile)
}
