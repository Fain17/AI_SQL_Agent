FROM python:3.12-alpine

# Set working directory
WORKDIR /app

# Install system dependencies (for pdfminer and model ops)
RUN apk add --no-cache gcc g++ musl-dev libffi-dev

# Copy requirements
COPY requirements.txt .

# Install Python packages
RUN pip install --no-cache-dir -r requirements.txt

# Copy the app code
COPY app/ ./app/

# Expose port
EXPOSE 8000

# Run the embedder FastAPI app
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
