#!/bin/bash

cd /opt/api/exam_api || exit

echo "📥 正在拉取最新代码..."
git fetch origin
git reset --hard origin/main

# 清除可能干扰的环境变量
unset DOCKER_HOST

echo "🐳 正在构建并启动 Docker 容器..."

# -------------------------------
# ✅ 关键步骤：强制删除旧容器（无论是否在运行）
# -------------------------------
if docker ps -a --format '{{.Names}}' | grep -Eq "^exam_api$"; then
    echo "⚠️  发现旧容器 'exam_api'，正在停止并删除..."
    docker stop exam_api >/dev/null 2>&1 || true
    docker rm exam_api >/dev/null 2>&1
    echo "🗑️  旧容器 'exam_api' 已删除"
fi

# -------------------------------
# ✅ 构建并启动新容器（使用固定名称）
# -------------------------------
echo "🏗️  正在构建镜像..."
docker compose build --no-cache

echo "🚀 启动新容器 'exam_api'..."
docker compose up -d

# -------------------------------
# ✅ 额外清理：删除悬空镜像（可选）
# -------------------------------
docker image prune -f --filter "label=stage=builder" >/dev/null 2>&1

echo "✅ 部署完成！新容器名为: exam_api"