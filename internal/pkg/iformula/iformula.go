package iformula

import (
	"fmt"
	"github.com/expr-lang/expr"
	"math"
)

type ScoreFormulaConfig struct {
	Expression string
	Rounding   int32
}

// ValidateExpression 检查给定的表达式是否合法，并可选地验证变量是否存在
func ValidateExpression(expression string, allowedVariables map[string]interface{}) error {
	// 编译表达式
	_, err := expr.Compile(expression, expr.Env(allowedVariables))
	if err != nil {
		return fmt.Errorf("表达式无效: %w", err)
	}
	return nil
}

// Evaluate 执行公式并返回结果
func Evaluate(config *ScoreFormulaConfig, rawScore, averageMark, standardMark float64) (float64, error) {
	env := map[string]interface{}{
		"raw_score":     rawScore,
		"average_mark":  averageMark,
		"standard_mark": standardMark,
	}

	program, err := expr.Compile(config.Expression, expr.Env(env))
	if err != nil {
		return 0, fmt.Errorf("compile failed: %w", err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return 0, fmt.Errorf("run failed: %w", err)
	}

	result, ok := output.(float64)
	if !ok {
		return 0, fmt.Errorf("result is not a number")
	}

	return math.Round(result*math.Pow(10, float64(config.Rounding))) / math.Pow(10, float64(config.Rounding)), nil
}
