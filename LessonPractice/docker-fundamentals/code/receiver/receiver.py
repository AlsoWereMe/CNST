from http.server import BaseHTTPRequestHandler, HTTPServer
import socketserver
import io

# 处理请求的类
class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # 获取请求的路径和参数
        ca = self.client_address
        
        # 创建一个响应内容的字符串
        response_content = f"Received request from {ca}"
        
        # 设置响应头
        self.send_response(200)
        self.send_header('Content-type', 'text/plain; charset=utf-8')
        self.end_headers()
        
        # 将响应内容写入文件
        with open('/file/request_log.txt', 'a', encoding='utf-8') as file:
            file.write(response_content + '\n')
        
        # 发送响应内容给客户端
        self.wfile.write(response_content.encode())

# 创建HTTP服务器
with socketserver.TCPServer(("", 8081), RequestHandler) as httpd:
    print("Server started at localhost:8081")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("Server stopped by user.")