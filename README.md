# LRU Cache

### How to start
---
- Clone the repository
- Run 
```bash
go build .
```
- Run 
```bash
./lru-cache
```


### How to use
---
- Use Frontend to interact with the LRU Cache
  - [Frontend](https://github.com/GauravMakhijani/LRU-cache-frontend)

- Steps : 
  - First Initailize the Cache with the capacity, capacity should be greater than 0
  - Add key-value-expiry pairs to the cache; expiry should be greater than current time
  - after adding the key-value pairs, websocket will notify the frontend about the changes in the cache
  - use get cache endpoint to get the value of the key
  - use delete cache endpoint to delete the key from the cache
  - use websocket to get the changes in the cache

### About
---
- This is a simple LRU Cache implementation in Golang
- It uses a doubly linked list and a hashmap to store the key-value pairs
- The cache capacity can be set by the user
- The cache is thread-safe
- cache has websocket support to notify the frontend about the changes in the cache
- The cache has a REST API to interact with the cache
- The cache uses cron to remove expired keys from the cache
