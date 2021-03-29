# gRPC-Service

gRPC架構實作與Protobuf撰寫

## Descript

1. Protobuf撰寫並利用protoc與protoc-gen-go產生對應的.pb.go檔案
2. 使用grpc package實作RPC架構，
   - Unary RPC架構
   - Server-side streaming RPC架構
   - Client-side streaming RPC架構
   - Bidirectional streaming RPC架構
3. 實作gRPC Interceptor，進行逾時控制
4. 使用grpcurl工具，查詢該RPC服務的方法與偵錯
