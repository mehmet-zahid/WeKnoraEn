#!/usr/bin/env python3
"""
WeKnora MCP Server Convenience Startup Script

This is a simplified startup script that provides basic functionality.
For more options, please use main.py
"""

import os
import sys
from pathlib import Path


def main():
    """Simple startup function"""
    # Add current directory to Python path
    current_dir = Path(__file__).parent.absolute()
    if str(current_dir) not in sys.path:
        sys.path.insert(0, str(current_dir))

    # Check environment variables
    base_url = os.getenv("WEKNORA_BASE_URL", "http://localhost:8080/api/v1")
    api_key = os.getenv("WEKNORA_API_KEY", "")

    print("WeKnora MCP Server")
    print(f"Base URL: {base_url}")
    print(f"API Key: {'Set' if api_key else 'Not set'}")
    print("-" * 40)

    try:
        # Import and run
        from main import sync_main

        sync_main()
    except ImportError:
        print("Error: Unable to import required modules")
        print("Please ensure to run: pip install -r requirements.txt")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\nServer stopped")
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
