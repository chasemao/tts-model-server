# TTS Model Server  

TTS Model Server is a `gin`-based API server for text-to-speech (TTS) with multiple model integrations.  

## Features  

- **API Access for TTS Models**: Use the API to convert text into speech.  
- **Web UI**: An interactive interface to test and use different TTS models.  
- **Legado App Integration**: Generate subscription URLs for `legado` app and use it as a TTS engine.  

## Public Demo  

You can test the public server here: [TTS Model Server Web UI](https://chasemao.com/tts/webui/)  
*Please use it responsibly and avoid excessive requests.*  

## Installation  

You can install the server in two ways:  

1. **Build from Source**  
   ```bash
   git clone https://github.com/chasemao/tts-model-server.git
   cd tts-model-server
   go build -o tts-model-server
   ```  

2. **Download the Latest Release**  
   Grab the latest pre-built binary from the [Releases Page](https://github.com/chasemao/tts-model-server/releases).  

## Running the Server  

Start the server with the default settings:  
```bash
./tts-model-server
```  
By default, the server listens on `0.0.0.0:1233`.  

### Command-line Flags  

| Flag  | Description | Default |
|-------|------------|---------|
| `--token` | API access token (optional) | *empty* |
| `--ip` | Server IP to bind to | `0.0.0.0` |
| `--port` | Port to listen on | `1233` |

## Web UI  

Visit: [http://localhost:1233/tts/webui](http://localhost:1233/tts/webui)  
- Select a TTS model  
- Input text and generate speech  
- Copy the subscription link for `legado` app  

## Legado Integration  

1. Download the `legado` app from [GitHub](https://github.com/gedoor/legado).  
2. Generate a subscription URL using the Web UI.  
3. Import the subscription URL into the `legado` app.  

### Steps to Import in Legado  

1. Open `legado` settings  
2. Select "Speech Engine"  
3. Choose "Import Engine"  
4. Paste the generated subscription URL  

![Open Speech](./doc/open%20speech.png)  
![Open Speech Setting](./doc/open%20speech%20setting.png)  
![Choose Speak Engine](./doc/choose%20speak%20engine.png)  
![Import Engine](./doc/import%20engine.png)  

## Supported Models  

| Model | Description |
|-------|-------------|
| **Edge Read Aloud API** | Uses Microsoft Edge TTS API |
| **[Coqui-ai/TTS](https://github.com/coqui-ai/TTS)** | Requires installation in a `conda` base environment |

## API Endpoints  

| Endpoint | Description |
|----------|------------|
| `/tts/webui/` | Web UI for selecting models, testing TTS, and generating subscription links |
| `/tts/api/invoke` | Main API for TTS conversion (`text` + `model` → speech file) |
| `/tts/api/fields` | Retrieves required parameters for each TTS model |
| `/tts/api/subscribe` | Generates a `legado` subscription configuration |