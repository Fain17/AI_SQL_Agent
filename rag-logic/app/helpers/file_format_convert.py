from pathlib import Path
from fastapi import UploadFile
from pdfminer.high_level import extract_text

async def to_text(file: UploadFile) -> str:
    ext = Path(file.filename).suffix.lower()
    data = await file.read()

    if ext == '.pdf':
        tmp = f"/tmp/{file.filename}"
        with open(tmp, 'wb') as f:
            f.write(data)
        return extract_text(tmp)

    return data.decode('utf-8')
