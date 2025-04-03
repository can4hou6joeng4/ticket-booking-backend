#!/bin/bash

# 更新 Swagger 文档
echo "正在更新 Swagger 文档..."
~/go/bin/swag init -g cmd/api/main.go

# 检查是否成功生成文档
if [ $? -eq 0 ]; then
    echo "Swagger 文档更新成功！"
    echo "您可以通过访问 http://localhost:8081/swagger/ 查看更新后的文档"
else
    echo "Swagger 文档更新失败，请检查错误信息"
fi 