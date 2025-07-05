from fastapi import HTTPException, UploadFile, File
from sqlalchemy.orm import Session
from pathlib import Path
from app.helpers.file_format_convert import to_text
from app.services.embedding import embed_text
import httpx
import requests

GO_BACKEND_URL = "http://127.0.0.1:8080"
ALLOWED = {'.txt', '.md', '.pdf'}

async def upload_file_service(file: UploadFile = File(...)):
    
    try:
        ext = Path(file.filename).suffix.lower()
        if ext not in ALLOWED:
            raise HTTPException(400, f"Unsupported file type: {ext}")

        file_content = await to_text(file)
        embedding = await embed_text(file_content)

        payload = {
            "filename": file.filename,
            "content": file_content,
            "embedding": embedding
        }

        async with httpx.AsyncClient() as client:
            resp = await client.post(f"{GO_BACKEND_URL}/files/upload", json=payload)
            resp.raise_for_status()
            return resp.json()

    except httpx.RequestError as e:
        raise HTTPException(status_code=502, detail=f"Failed to reach backend: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Unexpected error: {str(e)}")

async def delete_file_service(file_id: str):
    try:
        async with httpx.AsyncClient() as client:
            resp = await client.delete(f"{GO_BACKEND_URL}/files/{file_id}")
            if resp.status_code == 200:
                return {"message": f"File {file_id} deleted successfully."}
            raise HTTPException(status_code=resp.status_code, detail="Failed to delete file.")
    except httpx.RequestError as e:
        raise HTTPException(status_code=500, detail=f"Connection error: {str(e)}")

async def soft_delete_file_service(file_id: str):
    try:
        async with httpx.AsyncClient() as client:
            resp = await client.patch(
                f"{GO_BACKEND_URL}/files/{file_id}/soft-delete",
                json={"deleted": True}
            )
            if resp.status_code == 200:
                return {"message": f"File {file_id} soft-deleted successfully."}
            raise HTTPException(status_code=resp.status_code, detail="Failed to soft-delete file.")
    except httpx.RequestError as e:
        raise HTTPException(status_code=500, detail=f"Connection error: {str(e)}")
    
async def restore_file_service(file_id: str):
    try:
        async with httpx.AsyncClient() as client:
            resp = await client.patch(
                f"{GO_BACKEND_URL}/files/{file_id}/restore",
                json={"deleted": False}
            )
            if resp.status_code == 200:
                return {"message": f"File {file_id} restored successfully."}
            raise HTTPException(status_code=resp.status_code, detail="Failed to restore file.")
    except httpx.RequestError as e:
        raise HTTPException(status_code=500, detail=f"Connection error: {str(e)}")
    
async def get_soft_deleted_files_service():
    try:
        async with httpx.AsyncClient() as client:
            resp = await client.get(f"{GO_BACKEND_URL}/files/recycle-bin")
            if resp.status_code == 200:
                return resp.json()
            raise HTTPException(status_code=resp.status_code, detail="Failed to retrieve soft-deleted files.")
    except httpx.RequestError as e:
        raise HTTPException(status_code=500, detail=f"Connection error: {str(e)}")



async def update_file_service(file_id: str, file: UploadFile):
    try:
        ext = Path(file.filename).suffix.lower()
        print(ext)
        if ext not in ALLOWED:
            raise HTTPException(status_code=400, detail=f"Unsupported file type: {ext}")

        file_content = await to_text(file)
        embedding = await embed_text(file_content)
        
        payload = {
            "filename": file.filename,
            "content": file_content,
            "embedding": embedding
        }

        async with httpx.AsyncClient() as client:
            resp = await client.put(f"{GO_BACKEND_URL}/files/{file_id}", json=payload)
            if resp.status_code == 200:
                return {"message": "File updated successfully."}
            raise HTTPException(status_code=resp.status_code, detail="Failed to update file.")

    except httpx.RequestError as e:
        raise HTTPException(status_code=502, detail=f"Connection error: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Unexpected error: {str(e)}")
