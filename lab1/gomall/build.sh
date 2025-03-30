#!/bin/bash
# build-and-run-v2.sh

set -eo pipefail

# 微服务列表
declare -a services=("cart" "checkout" "email" "frontend" "order" "payment" "product" "user")

# 主流程
for service in "${services[@]}"; do
    echo "🔄 处理服务: $service"
    
    # 清理旧镜像（强制删除所有标签）
    if docker images -q "${service}-service" | grep -q .; then
        echo "🧹 清理旧镜像: ${service}-service"
        docker rmi -f $(docker images -q "${service}-service") 2>/dev/null || true
    fi

    # 构建镜像
    docker build -f "app/$service/Dockerfile" -t "${service}-service" . || {
        echo "❌ Build failed for $service"
        continue
    }

    # 删除临时的<none>镜像
    echo "🧽 清理中间镜像"
    docker image prune -f --filter "label=stage=builder"
done
