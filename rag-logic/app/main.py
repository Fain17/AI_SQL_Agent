from fastapi import FastAPI
from app.routes import file_routes

app = FastAPI(title="Business Logic API", version="1.0")
app.include_router(file_routes.router)
