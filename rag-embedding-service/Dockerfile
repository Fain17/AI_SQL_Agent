FROM python:3.12-slim-bookworm

# Set working directory
WORKDIR /app

# Install system dependencies (for Debian/Ubuntu)
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    make \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python packages
RUN pip install --no-cache-dir --upgrade pip setuptools wheel && \
    pip install --no-cache-dir -r requirements.txt

# Copy the app code
COPY app/ ./app/

# Expose port
EXPOSE 8001

# Run the embedder FastAPI app
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8001"]
