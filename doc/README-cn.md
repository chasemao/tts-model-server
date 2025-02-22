# TTS Model Server  

TTS Model Server 是一个基于 `gin` 的文本转语音（TTS）API 服务器，支持多种 TTS 模型。  

## 功能特点  

- **TTS API 接口**：通过 API 直接调用 TTS 功能，实现文本转语音。  
- **Web UI 交互界面**：提供可视化界面，选择模型并试听生成的语音。  
- **Legado（阅读）App 订阅引擎**：可生成订阅 URL，在 `legado` App 中作为 TTS 引擎使用。  

## 公开服务器  

你可以访问公开服务器进行测试：[TTS Model Server Web UI](https://chasemao.com/tts/webui/)  
*请勿滥用或过度请求，以免影响服务质量。*  

## 安装方法  

你可以选择以下两种方式安装服务器：  

### 1. 从源码构建  

```bash
git clone https://github.com/chasemao/tts-model-server.git
cd tts-model-server
go build -o tts-model-server
```  

### 2. 下载最新版本  

从 [Releases](https://github.com/chasemao/tts-model-server/releases) 页面下载最新的预编译二进制文件。  

## 启动服务器  

使用默认配置启动服务器：  

```bash
./tts-model-server
```  

服务器默认监听 `0.0.0.0:1233` 端口。  

### 命令行参数  

| 参数 | 说明 | 默认值 |
|------|------|-------|
| `--token` | API 访问所需的令牌（可选） | *空* |
| `--ip` | 服务器监听的 IP 地址 | `0.0.0.0` |
| `--port` | 服务器监听的端口号 | `1233` |

## Web UI 使用方法  

访问：[http://localhost:1233/tts/webui](http://localhost:1233/tts/webui)  

- 选择 TTS 模型  
- 输入文本，生成语音  
- 复制 `legado` 订阅链接  

## `legado`（阅读）App 订阅配置  

1. 从 [GitHub](https://github.com/gedoor/legado) 下载 `legado` App。  
2. 通过 Web UI 生成 TTS 订阅 URL。  
3. 在 `legado` App 中导入订阅链接。  

### `legado` 订阅导入步骤  

1. 打开 `legado` 设置  
2. 选择 “朗读引擎”  
3. 选择 “导入引擎”  
4. 输入生成的订阅 URL  

![打开朗读设置](./open%20speech.png)  
![朗读设置](./open%20speech%20setting.png)  
![选择朗读引擎](./choose%20speak%20engine.png)  
![导入朗读引擎](./import%20engine.png)  

## 支持的 TTS 模型  

| 模型 | 说明 |
|------|------|
| **Edge Read Aloud API** | 使用微软 Edge TTS API |
| **[Coqui-ai/TTS](https://github.com/coqui-ai/TTS)** | 需要在 `conda` 基础环境中安装 |

## API 接口  

| API 端点 | 说明 |
|----------|------|
| `/tts/webui/` | 访问 Web UI，选择模型，试听 TTS 语音，生成 `legado` 订阅链接 |
| `/tts/api/invoke` | 主要 API 接口，输入 `text` 和 `model`，返回语音文件 |
| `/tts/api/fields` | 获取各个 TTS 模型所需的参数 |
| `/tts/api/subscribe` | 生成 `legado` TTS 订阅配置 |