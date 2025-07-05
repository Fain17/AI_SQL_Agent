from fastapi import APIRouter, UploadFile, File, HTTPException, Depends
from app.helpers.file_format_convert import to_text
from app.services.embedding import embed_text
from pathlib import Path
from app.models.schemas import QueryRequest, QueryResponse
from app.services.query_service import run_query_pipeline
import httpx
from sqlalchemy.orm import Session
from app.db.session import get_db_session
import app.services.file_operations as fo

router = APIRouter(prefix="/file", tags=["file"])
ALLOWED = {'.txt', '.md', '.pdf'}

@router.post("/upload")
async def upload_file(file: UploadFile = File(...)):
    return await fo.upload_file_service(file)

@router.put("/{file_id}")
async def update_file(file_id: str, file: UploadFile = File(...), db: Session = Depends(get_db_session)):
    return await fo.update_file_service(file_id, file)

@router.post("/query", response_model=QueryResponse)
async def query_route(req: QueryRequest, db: Session = Depends(get_db_session)):
    files, answer = await run_query_pipeline(req.prompt, db)
    return QueryResponse(matches=files, answer=answer)

@router.delete("/{file_id}")
async def delete_file(file_id: str, db: Session = Depends(get_db_session)):
    return await fo.delete_file_service(file_id)

@router.patch("/{file_id}/soft-delete")
async def soft_delete_file(file_id: str):
    return await fo.soft_delete_file_service(file_id)

@router.patch("/{file_id}/restore")
async def restore_file(file_id: str):
    return await fo.restore_file_service(file_id)

@router.get("/recycle-bin")
async def get_soft_deleted_files():
    return await fo.get_soft_deleted_files_service()

