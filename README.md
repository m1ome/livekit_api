# Livekit Token Api
> Basic and simple livekit public api token generator

## Deployment
```shell
docker run -d -p 80:3000 w1n2k/livekit_api:1.0 -api_key=<your key> -api_secret=<your_secret>
```

## Token Request
```shell
curl -X "POST" "http://localhost:3000/tokens" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "roomName": "random",
  "metadata": "metadata",
  "userName": "username"
}'
```