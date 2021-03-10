# docker build -t registry.cn-qingdao.aliyuncs.com/pantazheng/minifab:latest .
# # # docker push registry.cn-qingdao.aliyuncs.com/pantazheng/minifab:latest
# ./minifab cleanup
# ./minifab up
# ./minifab invoke -p '"AddRecord"' -t '{"record":"'$(echo '{"timestamp":"'$(date '+%s')'","device_id":"01","temperature": 0.1}' | base64 | tr -d \\n)'"}'
./minifab invoke -p '"qscc"' -t '{"record":"'$(echo '{"timestamp":"'$(date '+%s')'","device_id":"01","temperature": 0.1}' | base64 | tr -d \\n)'"}'