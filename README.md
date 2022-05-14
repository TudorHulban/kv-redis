# Interacting with KV stores: Redis DB
## Rollout Redis container
```sh
docker run --name redis-test -p 6379:6379  -d redis 
```
## Connect to container
```sh
docker exec -it <container id> bash
```
## Connect to the Redis instance
Use the RedisInsight application.
## Resources
```
https://www.tutorialspoint.com/redis/index.htm
https://developer.redis.com/develop/golang/
```