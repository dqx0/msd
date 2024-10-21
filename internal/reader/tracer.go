package reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	core "github.com/dqx0/msd/pkg/core"
)

// 指定されたディレクトリ内の.tracerファイルを全て読み取る関数
func GetTracers(root string) ([]core.IParticle, error) {
	var particles []core.IParticle

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tracer") {
			fmt.Printf("Reading file: %s\n", path)

			particle, err := ReadTracerFile(path)
			if err != nil {
				return err
			}
			particles = append(particles, particle)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return particles, nil
}

// .tracerファイルの内容を読み取って表示する関数
func ReadTracerFile(filePath string) (core.IParticle, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 1行目を読み取り、行数を取得
	scanner.Scan()
	line := scanner.Text()
	numLines, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return nil, fmt.Errorf("failed to parse number of lines: %v", err)
	}

	// x, y, zのスライスを作成
	x := make([]float64, 0, numLines)
	y := make([]float64, 0, numLines)
	z := make([]float64, 0, numLines)

	// 2行目以降を読み取り、x, y, zに格納
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		// x, y, zの各値を取得して格納
		xVal, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse x value: %v", err)
		}
		yVal, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse y value: %v", err)
		}
		zVal, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse z value: %v", err)
		}

		x = append(x, xVal)
		y = append(y, yVal)
		z = append(z, zVal)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	particle := core.NewParticle()
	particle.Append(core.ParticlePath{
		X: x,
		Y: y,
		Z: z,
	})
	return particle, nil
}
