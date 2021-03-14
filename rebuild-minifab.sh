docker build -t registry.cn-qingdao.aliyuncs.com/pantazheng/minifab:latest .
# # # docker push registry.cn-qingdao.aliyuncs.com/pantazheng/minifab:latest
./minifab cleanup
./minifab up
# ./minifab up
# ./minifab invoke -p '"AddPublicRecord","'$(date '+%s')'","01","0.1"'
