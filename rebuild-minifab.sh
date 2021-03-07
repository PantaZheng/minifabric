docker build -t registry.cn-qingdao.aliyuncs.com/pantazheng/minifab:latest .
# docker push registry.cn-qingdao.aliyuncs.com/pantazheng/minifab:latest
./minifab cleanup
./minifab up
./minifab invoke -p '"AddRecord"' -t '{"record":"eyJ0aW1lc3RhbXAiOiIxNjE1MDQxMTM2IiwiZGV2aWNlX2lkIjoiMDEiLCJ0ZW1wZXJhdHVyZSI6IDAuMX0K"}'