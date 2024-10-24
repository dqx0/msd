package reader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	core "github.com/dqx0/msd/pkg/core"
)

func GetTracers(root string) ([]core.IParticle, int, error) {
	return getTracers(root)
}

func readTracerFileHeader(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to read header")
	}
	line := strings.TrimSpace(scanner.Text())
	return strconv.Atoi(line)
}

func parseCoordinates(line string) (float64, float64, float64, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid line format: %s", line)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse x value: %v", err)
	}
	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse y value: %v", err)
	}
	z, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse z value: %v", err)
	}

	return x, y, z, nil
}

func readTracerFile(filePath string, particles []core.IParticle) ([]core.IParticle, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// バッファサイズを最適化
	const bufferSize = 64 * 1024 // 64KB
	scanner := bufio.NewScanner(file)
	buf := make([]byte, bufferSize)
	scanner.Buffer(buf, bufferSize)

	// ヘッダーを読み取り
	lineNum, err := readTracerFileHeader(scanner)
	if err != nil {
		return nil, err
	}

	// 初回のみパーティクルスライスを初期化
	if particles == nil {
		particles = make([]core.IParticle, lineNum)
		for i := range particles {
			particles[i] = &core.ParticlePath{}
		}
	}

	// データ行を読み取り
	index := 0
	for scanner.Scan() && index < lineNum {
		x, y, z, err := parseCoordinates(scanner.Text())
		if err != nil {
			return nil, err
		}
		particles[index].Append(x, y, z)
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return particles, nil
}

func getTracers(root string) ([]core.IParticle, int, error) {
	var particles []core.IParticle
	i := 0

	// ファイル名のプレフィックスを事前に作成
	const batchSize = 100
	prefix := root + "\\"

	log.Printf("Reading tracers from %s", root)
	for i = 1; ; i++ {
		path := prefix + strconv.Itoa(i*batchSize) + ".tracers"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}

		var err error
		particles, err = readTracerFile(path, particles)
		if err != nil {
			return nil, 0, err
		}
	}

	return particles, (i - 1) * 100, nil
}
