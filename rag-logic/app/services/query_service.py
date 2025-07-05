from app.services.embedding import get_embedding
from app.services.llm_chain import chain
from app.models.schemas import FileData
from sqlalchemy.orm import Session
from sqlalchemy import text, bindparam
from app.models.schemas import FileData
from app.db.models import File


async def fetch_similar_files_pgvector(embedding: list[float], db: Session, top_k: int = 5) -> list[FileData]:
    # query = text("""
    #     SELECT filename, content, embedding <=> :embedding::vector AS similarity
    #     FROM files
    #     ORDER BY embedding <=> :embedding::vector
    #     LIMIT :top_k
    # """)
    
    # query = text("""
    # SELECT filename, content, embedding <=> :embedding::vector AS similarity
    # FROM files
    # ORDER BY embedding <=> :embedding::vector
    # LIMIT :top_k
    # """).bindparams(
    #         bindparam("embedding"),  # ⬅ no type here!
    #         bindparam("top_k")
    #     )

    # result = db.execute(query, {"embedding": embedding, "top_k": top_k})
    # rows = result.fetchall()
    
    # Turn list into Postgres vector format: '[0.1, 0.2, ...]'
    
    embedding_str = f"[{', '.join(map(str, embedding))}]"

    query = text("""
        SELECT filename, content, embedding <=> CAST(:embedding AS vector) AS similarity
        FROM files
        ORDER BY embedding <=> CAST(:embedding AS vector)
        LIMIT :top_k;
    """)

    result = db.execute(query, {"embedding": embedding_str, "top_k": top_k})
    
    rows  = result.fetchall()
    
    return [FileData(filename=row[0], content=row[1], similarity=row[2]) for row in rows]


async def run_query_pipeline(prompt: str, db: Session) -> tuple[list[FileData], str]:
    embedding = await get_embedding(prompt)
    files = await fetch_similar_files_pgvector(embedding, db)

    context = "\n\n".join(f"{f.filename}:\n{f.content}" for f in files)
    answer = chain.invoke({"context": context, "question": prompt})

    return files, answer