from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer

app = FastAPI()
model = SentenceTransformer("all-MiniLM-L6-v2")


class EmbedRequest(BaseModel):
    text: str


@app.post("/embed")
def embed(req: EmbedRequest):
    if not req.text:
        raise HTTPException(status_code=400, detail="No text provided")

    try:
        vector = model.encode([req.text])[0].tolist()
        return {"embedding": vector}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Model error: {str(e)}")
