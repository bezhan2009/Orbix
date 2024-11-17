package PiCalculation

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// calculatePi теперь принимает контекст для прерывания вычислений
func calculatePi(ctx context.Context, precision int) (*big.Float, error) {
	pi := new(big.Float).SetPrec(uint(precision)).SetFloat64(0)
	k := new(big.Float).SetPrec(uint(precision))
	n := new(big.Float).SetPrec(uint(precision))
	one := new(big.Float).SetPrec(uint(precision)).SetInt64(1)
	four := new(big.Float).SetPrec(uint(precision)).SetInt64(4)

	for i := int64(0); i < int64(precision); i++ {
		select {
		// Проверяем, не был ли отменён контекст
		case <-ctx.Done():
			// Если контекст отменён, возвращаем ошибку
			return nil, ctx.Err()
		default:
			// Продолжаем вычисления
			k.SetInt64(2 * i)
			n.Quo(one, new(big.Float).SetPrec(uint(precision)).SetInt64(2*i+1))
			if i%2 == 0 {
				pi.Add(pi, n)
			} else {
				pi.Sub(pi, n)
			}
		}
	}

	pi.Mul(pi, four)
	return pi, nil
}

func PiCalcCommand(precision int) {
	green := color.New(color.FgGreen).SprintFunc()

	// Создаём контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())

	// Обработка сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Ожидание сигнала
		<-sigs
		fmt.Println("\nA termination signal has been received, interrupting calculations...")
		cancel() // Отмена контекста при получении сигнала
	}()

	start := time.Now()

	// Вызываем calculatePi с контекстом
	pi, err := calculatePi(ctx, precision)
	elapsed := time.Since(start)

	// Проверяем, была ли ошибка из-за отмены контекста
	if err != nil {
		fmt.Println("The calculations were interrupted:", err)
		return
	}

	// Выводим результат, если вычисления завершились успешно
	printResult := fmt.Sprintf("Calculated pi to %d digits in %s\n", precision, elapsed)
	fmt.Printf(green(printResult))
	fmt.Println(green("Pi:", pi))
}
