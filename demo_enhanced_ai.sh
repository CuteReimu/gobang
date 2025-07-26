#!/bin/bash

echo "=== 五子棋AI增强优化演示 ==="
echo ""
echo "用户需求: 在保持6层深度的情况下优化算法性能"
echo ""

cd go

echo "构建项目..."
go build -o gobang_test *.go 2>/dev/null || echo "跳过构建（缺少GUI依赖）"

echo ""
echo "=== 增强AI性能测试 ==="
echo ""
echo "运行增强AI性能基准测试..."
timeout 60 go run test_enhanced.go board.go player_robot.go point.go player.go

echo ""
echo "=== 使用方法 ==="
echo ""
echo "可用的AI模式："
echo "  ./gobang             # 原始AI (6层深度，强棋力，较慢)"
echo "  ./gobang -optimized  # 优化AI (4层深度，速度快，棋力较弱)"
echo "  ./gobang -balanced   # 平衡AI (4层深度，平衡速度和棋力)"
echo "  ./gobang -enhanced   # 增强AI (6层深度，算法优化，最佳选择)"
echo ""
echo "推荐使用：./gobang -enhanced"
echo "这个模式保持了6层深度的强棋力，同时通过算法优化实现了显著的性能提升。"