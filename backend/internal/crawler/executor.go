package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ExecutePython 执行Python脚本并返回结果
func ExecutePython(script string, keyword string) ([]map[string]interface{}, string, error) {
	// 直接执行脚本内容
	cmd := exec.Command("python3", "-c", script)
	cmd.Env = append(cmd.Env, fmt.Sprintf("KEYWORD=%s", keyword))

	// 设置超时
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	// 获取stderr作为日志
	logOutput := stderr.String()

	if err != nil {
		return nil, logOutput, fmt.Errorf("执行失败: %v, 耗时: %v", err, duration)
	}

	// 解析stdout中的JSON
	output := stdout.String()
	output = strings.TrimSpace(output)

	// 尝试解析JSON
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		// 如果不是数组，尝试解析单个对象
		var singleResult map[string]interface{}
		if err := json.Unmarshal([]byte(output), &singleResult); err != nil {
			return nil, logOutput + "\n" + output, nil
		}
		result = []map[string]interface{}{singleResult}
	}

	return result, fmt.Sprintf("%s\n执行成功，耗时: %v", logOutput, duration), nil
}

// ValidateScript 验证脚本安全性
func ValidateScript(script string) error {
	// 简单的安全检查
	dangerous := []string{"os.system", "subprocess.call", "subprocess.run", "eval(", "exec("}
	for _, d := range dangerous {
		if strings.Contains(script, d) {
			return fmt.Errorf("脚本包含不安全的代码: %s", d)
		}
	}
	return nil
}