# 🧠 Sentence Embedding API

A lightweight FastAPI service that exposes a REST endpoint to generate **sentence embeddings** using the [`sentence-transformers`](https://www.sbert.net/) model **`all-MiniLM-L6-v2`**.

This service is useful in NLP pipelines, RAG architectures, semantic search, or any system that needs to convert user input into embedding vectors.

---

## 🚀 Features

- 🧠 Uses the `all-MiniLM-L6-v2` model (384-dimensional embeddings)
- ⚡ Powered by [FastAPI](https://fastapi.tiangolo.com/)
- 🔌 Exposes a simple `/embed` POST endpoint
- 🐳 Easily containerizable via Docker

---

## 📡 API Endpoint

### `POST /embed`

Generate sentence embedding for a given string.

#### 🔸 Request

```json
{
  "text": "This is an example sentence"
}

```

#### 🔸 Response

```json

{
  "embedding": [0.123, -0.456, 0.789, ...]   // 384-dimensional float array
}

```

####  🔸 Docker pull and run Command - Docker images can be found in packages

```
docker pull ghcr.io/fain17/ai_rag_agent/rag-embedding-service:dev
```

```
docker run -p 8001:8001 rag-embedding-service:dev
```