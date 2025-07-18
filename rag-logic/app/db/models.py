from typing import Any

from pgvector.sqlalchemy import Vector  # type: ignore
from sqlalchemy import Column, String, Text
from sqlalchemy.ext.declarative import declarative_base

Base: Any = declarative_base()


class File(Base):
    __tablename__ = "files"

    filename = Column(String, primary_key=True)
    content = Column(Text)
    embedding: Any = Column(Vector(384))
