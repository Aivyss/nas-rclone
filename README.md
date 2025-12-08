# 사전준비
- `rclone`(리눅스 아키텍처 arm64)의 바이너리 파일이 이 프로젝트의 루트 폴더 안에 있어야 합니다.
- `.config/rclone/rclone.conf`설정 파일을 설정하셔야 합니다.
- `.config/.config.json`설정 파일을 설정하셔야 합니다.
- (아마도) `docker-compose`가 실행가능한 환경이어야 합니다

# `.config/.config.json` 설정방법
docker-compose.yml의 volumes 설정에는 `DO NOT FIX THIS` 코멘트가 있는 설정이 있는데 이것은 건드리지 마십시오. 
해당 부분을 건드리지 않고 도커 컨테이너와 바인딩할 (sync할) 디렉토리를 추가하세요.

```yaml
    volumes:
      - '.:/app' # DO NOT FIX THIS (you need your rclone binary file and modified .config.json)
      - "./.config:/root/.config" # DO NOT FIX THIS
      - '/local/path/to/local_folder_a:/remote_folder_a' # 'source (in NAS):destination (in CloudStorage)'
      - '/local/path/to/local_folder_b:/remote_folder_b' # 'source (in NAS):destination (in CloudStorage)'
```

그 다음 `.config/.config.json` 설정 파일에도 동일하게 설정하셔야 합니다.

```json
{
  "workerConfigurations": [
    {
      "syncType": "copy",
      "alias": "cloud_a",
      "workerName": "worker1",
      "destinationPath": "/local_folder_a",
      "sourcePath": "/remote_folder_a",
      "cron": "* * * * *",
      "transfers": 8
    },
    {
      "syncType": "sync",
      "alias": "cloud_a",
      "workerName": "worker2",
      "destinationPath": "/local_folder_b",
      "sourcePath": "/remote_folder_b",
      "cron": "* * * * *",
      "transfers": 16
    }
  ]
}
```

# 실행
리눅스 기반에서는 다음과 같이 실행하세요
```bash
$ docker-compose up --build -d
```
컨테이너가 이미 만들어졌다면 `--build -d` 옵션을 붙이지 않아도 됩니다.