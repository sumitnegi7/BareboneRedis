# Barebone Redis Implementation 🚀

Welcome to a minimal Redis-like server implemented in Golang! 🌟 This lightweight server supports essential Redis commands such as `SET`, `GET`, `PING`, `ECHO`, and `SET` with expiry options. You can connect using `redis-cli` or any Redis-compatible client. 🛠️

---

## Features 🎯

- **`SET`**: Store a key-value pair 📥
- **`GET`**: Retrieve the value of a key 🔍
- **`PING`**: Check connectivity 🏓
- **`ECHO`**: Return your input as-is 🗣️
- **`SET` with options**: 
  - `PX <milliseconds>`: Set a key with a time-to-live ⏲️

---

## Commands Overview 📜

### 1. **SET** 📝
Store a key-value pair.  
**Syntax**: `SET key value [NX|XX] [PX milliseconds]`  
**Examples**:
- `SET mykey myvalue`  
- `SET mykey myvalue NX`  
- `SET mykey myvalue PX 5000`  

---

### 2. **GET** 📤
Retrieve the value associated with a key.  
**Syntax**: `GET key`  
**Example**:
- `GET mykey`  

---

### 3. **PING** 🏓
Test server connectivity.  
**Syntax**: `PING`  
**Example**:
- `PING`  
  **Response**: `PONG`  

---

### 4. **ECHO** 🗣️
Return the input string.  
**Syntax**: `ECHO message`  
**Example**:
- `ECHO "Hello World"`  
  **Response**: `Hello World`  

---

### 5. **SET with Expiry** ⏳
Store a key-value pair with a time-to-live.  
**Syntax**: `SET key value PX milliseconds`  
**Example**:
- `SET mykey myvalue PX 5000`  
  *(Key will expire after 5000 milliseconds)*  

---

## Running the Server 🖥️

1. Clone the repository:  
   ```bash
   git clone https://github.com/sumitnegi7/BareboneRedis/
   cd BareboneRedis/app
   ```

2. Build and run the server:  
   ```bash
   go run .
   ```

3. Use `redis-cli` or any Redis-compatible client to interact with the server:  
   ```bash
   redis-cli -h localhost -p 6379
   ```

---

## Notes 🗒️

- ⚡ This implementation is **in-memory** and does not persist data as of now.
- 🕒 TTL-based expiry is automatically handled by the server.
- 🤝 Fully compatible with the RESP protocol and Redis clients.

---

## License 📜

This project is licensed under the MIT License. See `LICENSE` for details.  
