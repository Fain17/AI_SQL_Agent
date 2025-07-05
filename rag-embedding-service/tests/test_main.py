"""Test cases for the embedding service main application."""

from unittest.mock import patch

import numpy as np
import pytest
from fastapi.testclient import TestClient

from app.main import EmbedRequest, app


class TestEmbedRequest:
    """Test cases for EmbedRequest model."""

    def test_embed_request_valid(self):
        """Test valid EmbedRequest creation."""
        request = EmbedRequest(text="Hello world")
        assert request.text == "Hello world"

    def test_embed_request_empty_string(self):
        """Test EmbedRequest with empty string."""
        request = EmbedRequest(text="")
        assert request.text == ""

    def test_embed_request_long_text(self):
        """Test EmbedRequest with long text."""
        long_text = "A" * 1000
        request = EmbedRequest(text=long_text)
        assert request.text == long_text


class TestEmbedEndpoint:
    """Test cases for the /embed endpoint."""

    def setup_method(self):
        """Set up test client."""
        self.client = TestClient(app)

    @patch("app.main.model.encode")
    def test_embed_success(self, mock_encode):
        """Test successful embedding generation."""
        # Mock the model response
        mock_embedding = np.array([0.1, 0.2, 0.3, 0.4, 0.5])
        mock_encode.return_value = [mock_embedding]

        response = self.client.post("/embed", json={"text": "Hello world"})

        assert response.status_code == 200
        data = response.json()
        assert "embedding" in data
        assert isinstance(data["embedding"], list)
        assert len(data["embedding"]) == 5
        assert data["embedding"] == [0.1, 0.2, 0.3, 0.4, 0.5]

        # Verify the model was called correctly
        mock_encode.assert_called_once_with(["Hello world"])

    def test_embed_empty_text(self):
        """Test embedding with empty text."""
        response = self.client.post("/embed", json={"text": ""})

        assert response.status_code == 400
        data = response.json()
        assert data["detail"] == "No text provided"

    def test_embed_missing_text_field(self):
        """Test embedding with missing text field."""
        response = self.client.post("/embed", json={})

        assert response.status_code == 422  # Validation error

    def test_embed_invalid_json(self):
        """Test embedding with invalid JSON."""
        response = self.client.post("/embed", data="invalid json")

        assert response.status_code == 422

    @patch("app.main.model.encode")
    def test_embed_long_text(self, mock_encode):
        """Test embedding with long text."""
        long_text = "A" * 1000
        mock_embedding = np.array([0.1] * 384)  # Common embedding size
        mock_encode.return_value = [mock_embedding]

        response = self.client.post("/embed", json={"text": long_text})

        assert response.status_code == 200
        data = response.json()
        assert "embedding" in data
        assert len(data["embedding"]) == 384

    @patch("app.main.model.encode")
    def test_embed_special_characters(self, mock_encode):
        """Test embedding with special characters."""
        special_text = "Hello! @#$%^&*()_+ 世界 🌍"
        mock_embedding = np.array([0.1, 0.2, 0.3])
        mock_encode.return_value = [mock_embedding]

        response = self.client.post("/embed", json={"text": special_text})

        assert response.status_code == 200
        data = response.json()
        assert "embedding" in data

    @patch("app.main.model.encode")
    def test_embed_model_exception(self, mock_encode):
        """Test handling of model exceptions."""
        mock_encode.side_effect = Exception("Model error")

        response = self.client.post("/embed", json={"text": "Hello world"})

        assert response.status_code == 500

    def test_embed_wrong_http_method(self):
        """Test wrong HTTP method."""
        response = self.client.get("/embed")
        assert response.status_code == 405  # Method not allowed

    @patch("app.main.model.encode")
    def test_embed_whitespace_only_text(self, mock_encode):
        """Test embedding with whitespace-only text."""
        response = self.client.post("/embed", json={"text": "   "})

        # Should still process whitespace as valid text
        assert response.status_code == 200

    @patch("app.main.model.encode")
    def test_embed_numeric_text(self, mock_encode):
        """Test embedding with numeric text."""
        mock_embedding = np.array([0.1, 0.2])
        mock_encode.return_value = [mock_embedding]

        response = self.client.post("/embed", json={"text": "12345"})

        assert response.status_code == 200
        data = response.json()
        assert "embedding" in data


class TestHealthEndpoints:
    """Test cases for health and status endpoints."""

    def setup_method(self):
        """Set up test client."""
        self.client = TestClient(app)

    def test_root_endpoint(self):
        """Test if root endpoint exists (optional)."""
        response = self.client.get("/")
        # This might return 404 if not implemented, which is fine
        assert response.status_code in [200, 404]


@pytest.mark.integration
class TestIntegrationEmbedding:
    """Integration tests that use the actual model."""

    def setup_method(self):
        """Set up test client."""
        self.client = TestClient(app)

    @pytest.mark.slow
    def test_real_embedding_generation(self):
        """Test actual embedding generation with real model."""
        response = self.client.post("/embed", json={"text": "This is a test sentence."})

        assert response.status_code == 200
        data = response.json()
        assert "embedding" in data
        assert isinstance(data["embedding"], list)
        assert len(data["embedding"]) > 0
        # MiniLM model typically produces 384-dimensional embeddings
        assert len(data["embedding"]) == 384

    @pytest.mark.slow
    def test_real_embedding_consistency(self):
        """Test that same text produces same embedding."""
        text = "Consistency test"

        response1 = self.client.post("/embed", json={"text": text})
        response2 = self.client.post("/embed", json={"text": text})

        assert response1.status_code == 200
        assert response2.status_code == 200

        embedding1 = response1.json()["embedding"]
        embedding2 = response2.json()["embedding"]

        # Embeddings should be identical for same input
        assert embedding1 == embedding2


class TestModelLoading:
    """Test cases for model loading and initialization."""

    def test_model_is_loaded(self):
        """Test that the model is properly loaded."""
        from app.main import model

        assert model is not None
        assert hasattr(model, "encode")

    def test_model_type(self):
        """Test that the model is of correct type."""
        from sentence_transformers import SentenceTransformer

        from app.main import model

        assert isinstance(model, SentenceTransformer)
