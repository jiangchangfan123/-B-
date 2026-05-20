#!/bin/bash
# 同时启动后端 + 前端（不依赖 concurrently，适配 WSL+NTFS）

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/../bilibili-backend"
FRONTEND_DIR="$SCRIPT_DIR"

echo "========================================"
echo "  BILIBILI CLONE - DEV STARTER"
echo "========================================"

# 启动后端
cd "$BACKEND_DIR"
echo "[BACKEND] Starting Go server..."
if [ -f ./server ]; then
    ./server > "$BACKEND_DIR/server.log" 2>&1 &
else
    echo "[BACKEND] ERROR: ./server not found. Run: go build -o server main.go"
    exit 1
fi
BACKEND_PID=$!
echo "[BACKEND] PID: $BACKEND_PID"

# 等待后端就绪
sleep 2
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "[BACKEND] ERROR: Server failed to start. Check $BACKEND_DIR/server.log"
    exit 1
fi

# 启动前端
cd "$FRONTEND_DIR"
echo "[FRONTEND] Starting Vite dev server..."
npx vite --host 0.0.0.0 &
FRONTEND_PID=$!
echo "[FRONTEND] PID: $FRONTEND_PID"

echo ""
echo "========================================"
echo "  Both servers started!"
echo "  Backend : http://localhost:8080"
echo "  Frontend: http://localhost:5173"
echo "========================================"
echo ""
echo "Press Ctrl+C to stop both servers"
echo ""

# 捕获 Ctrl+C 信号关闭两个进程
cleanup() {
    echo ""
    echo "[STOP] Shutting down servers..."
    kill $FRONTEND_PID 2>/dev/null || true
    kill $BACKEND_PID 2>/dev/null || true
    wait $FRONTEND_PID 2>/dev/null || true
    wait $BACKEND_PID 2>/dev/null || true
    echo "[STOP] Done."
    exit 0
}
trap cleanup INT TERM

# 等待前端进程
wait $FRONTEND_PID
