"""Pytest configuration and fixtures for the embedding service tests."""

from unittest.mock import Mock

import pytest
from fastapi.testclient import TestClient

from app.main import app


@pytest.fixture
def client():
    """Create a test client for the FastAPI app."""
    return TestClient(app)


@pytest.fixture
def mock_model():
    """Create a mock model for testing."""
    mock = Mock()
    mock.encode.return_value = [[0.1, 0.2, 0.3, 0.4, 0.5]]
    return mock


@pytest.fixture
def sample_text():
    """Provide sample text for testing."""
    return "This is a sample text for testing."


@pytest.fixture
def sample_embedding():
    """Provide sample embedding for testing."""
    return [0.1, 0.2, 0.3, 0.4, 0.5]


@pytest.fixture
def long_text():
    """Provide long text for testing."""
    return "A" * 1000


@pytest.fixture
def special_text():
    """Provide text with special characters for testing."""
    return "Hello! @#$%^&*()_+ 世界 🌍"
